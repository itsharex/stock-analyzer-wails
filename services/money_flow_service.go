package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// MoneyFlowService 资金流向服务
type MoneyFlowService struct {
	repo   *repositories.MoneyFlowRepository
	client *http.Client
}

// NewMoneyFlowService 创建资金流向服务
func NewMoneyFlowService(repo *repositories.MoneyFlowRepository) *MoneyFlowService {
	return &MoneyFlowService{
		repo: repo,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchAndSaveHistory 获取并保存近100日资金流向数据
func (s *MoneyFlowService) FetchAndSaveHistory(code string) error {
	logger.Info("开始抓取资金流向数据", zap.String("code", code))

	// 1. 生成 secid
	secid := s.generateSecid(code)

	// 2. 构建 API URL
	url := fmt.Sprintf("https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get?lmt=100&klt=101&fields1=f1,f2,f3,f7&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65&secid=%s", secid)

	// 3. 发起请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Referer", "https://quote.eastmoney.com/")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("请求 API 失败: %w", err)
	}
	defer resp.Body.Close()

	// 4. 解析响应
	var apiResp struct {
		RC   int `json:"rc"`
		Data struct {
			Klines []string `json:"klines"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if apiResp.Data.Klines == nil {
		logger.Warn("未获取到资金流向数据", zap.String("code", code))
		return nil
	}

	// 5. 解析数据
	var flows []models.MoneyFlowData
	for _, line := range apiResp.Data.Klines {
		parts := strings.Split(line, ",")
		if len(parts) < 13 {
			continue
		}

		flows = append(flows, models.MoneyFlowData{
			Code:       code,
			TradeDate:  parts[0],
			MainNet:    s.parseFloat(parts[1]),
			SmallNet:   s.parseFloat(parts[2]),
			MidNet:     s.parseFloat(parts[3]),
			BigNet:     s.parseFloat(parts[4]),
			SuperNet:   s.parseFloat(parts[5]),
			ClosePrice: s.parseFloat(parts[11]),
			ChgPct:     s.parseFloat(parts[12]),
		})
	}

	if len(flows) == 0 {
		return nil
	}

	// 6. 调用 Repository 保存
	if err := s.repo.SaveMoneyFlows(flows); err != nil {
		return fmt.Errorf("保存数据失败: %w", err)
	}

	logger.Info("成功保存资金流向数据", zap.Int("count", len(flows)))
	return nil
}

func (s *MoneyFlowService) generateSecid(code string) string {
	if strings.HasPrefix(code, "6") {
		return "1." + code
	} else if strings.HasPrefix(code, "0") || strings.HasPrefix(code, "3") {
		return "0." + code
	} else if strings.HasPrefix(code, "8") || strings.HasPrefix(code, "4") {
		return "2." + code
	}
	return "0." + code // Default fallback
}

func (s *MoneyFlowService) parseFloat(val string) float64 {
	if val == "-" {
		return 0.0
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0.0
	}
	return f
}
