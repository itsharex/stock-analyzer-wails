package services

import (
	"fmt"
	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// MoneyFlowService 资金流向服务
type MoneyFlowService struct {
	repo   *repositories.MoneyFlowRepository
	client *resty.Client
}

// NewMoneyFlowService 创建资金流向服务
func NewMoneyFlowService(repo *repositories.MoneyFlowRepository) *MoneyFlowService {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	client.SetHeader("Referer", "https://quote.eastmoney.com/")
	client.SetRetryCount(3)

	return &MoneyFlowService{
		repo:   repo,
		client: client,
	}
}

// FetchAndSaveHistory 获取并保存近100日资金流向数据
func (s *MoneyFlowService) FetchAndSaveHistory(code string) error {
	logger.Info("开始抓取资金流向数据", zap.String("code", code))

	// 1. 生成 secid
	secid := s.generateSecid(code)

	url := fmt.Sprintf(
		"https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get?lmt=0&klt=101&fields1=f1,f2,f3,f7&fields2=f51,f52,f53,f54,f55,f56,f62&secid=%s",
		secid,
	)

	var emResp struct {
		RC   int `json:"rc"`
		Data struct {
			Klines []string `json:"klines"`
		} `json:"data"`
	}

	resp, err := s.client.R().
		SetResult(&emResp).
		Get(url)

	if err != nil {
		return fmt.Errorf("HTTP请求失败: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("HTTP请求返回错误状态: %s", resp.Status())
	}

	if emResp.RC != 0 {
		return fmt.Errorf("API返回错误 RC=%d", emResp.RC)
	}

	if emResp.Data.Klines == nil {
		return nil
	}

	var flows []models.MoneyFlowData
	for _, line := range emResp.Data.Klines {
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
