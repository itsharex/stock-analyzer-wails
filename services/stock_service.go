package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"stock-analyzer-wails/models"
	"strings"
	"time"

	"stock-analyzer-wails/internal/logger"

	"go.uber.org/zap"
)

// StockService 股票数据服务
type StockService struct {
	exactURL string
	listURL  string
	klineURL string
	client   *http.Client
}

// NewStockService 创建股票数据服务实例
func NewStockService() *StockService {
	return &StockService{
		exactURL: "https://push2.eastmoney.com/api/qt/stock/get",
		listURL:  "http://78.push2.eastmoney.com/api/qt/clist/get",
		klineURL: "https://push2his.eastmoney.com/api/qt/stock/kline/get",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetStockByCode 根据股票代码获取股票数据（精确查询）
func (s *StockService) GetStockByCode(code string) (*models.StockData, error) {
	start := time.Now()
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}

	secid := s.getSecID(code)
	if secid == "" {
		return nil, fmt.Errorf("无法识别的股票代码格式: %s", code)
	}

	fields := "f58,f43,f169,f170,f47,f48,f44,f45,f46,f60,f171,f168,f162,f167,f116,f117,f12,f14"
	fullURL := fmt.Sprintf("%s?secid=%s&fields=%s", s.exactURL, secid, fields)

	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Data == nil {
		return nil, fmt.Errorf("未找到股票数据: %s", code)
	}

	data := result.Data
	stockCode := getString(data["f12"])
	if stockCode == "" {
		stockCode = code
	}
	
	stock := &models.StockData{
		Code:       stockCode,
		Name:       getString(data["f58"]),
		Price:      getFloat(data["f43"]) / 100,
		Change:     getFloat(data["f169"]) / 100,
		ChangeRate: getFloat(data["f170"]) / 100,
		Volume:     getInt64(data["f47"]),
		Amount:     getFloat(data["f48"]),
		High:       getFloat(data["f44"]) / 100,
		Low:        getFloat(data["f45"]) / 100,
		Open:       getFloat(data["f46"]) / 100,
		PreClose:   getFloat(data["f60"]) / 100,
		Amplitude:  getFloat(data["f171"]) / 100,
		Turnover:   getFloat(data["f168"]) / 100,
		PE:         getFloat(data["f162"]) / 100,
		PB:         getFloat(data["f167"]) / 100,
		TotalMV:    getFloat(data["f116"]),
		CircMV:     getFloat(data["f117"]),
	}

	logger.Info("精确获取股票数据成功", zap.String("code", code), zap.Int64("ms", time.Since(start).Milliseconds()))
	return stock, nil
}

// GetKLineData 获取历史K线数据并计算技术指标，支持周期选择
func (s *StockService) GetKLineData(code string, limit int, period string) ([]*models.KLineData, error) {
	secid := s.getSecID(code)
	if secid == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	// 映射周期到东方财富的 klt 参数
	// 101: 日线, 102: 周线, 103: 月线
	klt := "101"
	switch period {
	case "week":
		klt = "102"
	case "month":
		klt = "103"
	}

	fetchLimit := limit + 50
	url := fmt.Sprintf("%s?secid=%s&fields1=f1,f2,f3,f4,f5,f6&fields2=f51,f52,f53,f54,f55,f56&klt=%s&fqt=1&end=20500101&lmt=%d", s.klineURL, secid, klt, fetchLimit)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data struct {
			Klines []string `json:"klines"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	klines := make([]*models.KLineData, 0)
	for _, line := range result.Data.Klines {
		parts := strings.Split(line, ",")
		if len(parts) < 6 {
			continue
		}
		klines = append(klines, &models.KLineData{
			Time:   parts[0],
			Open:   s.parsePrice(parts[1]),
			Close:  s.parsePrice(parts[2]),
			High:   s.parsePrice(parts[3]),
			Low:    s.parsePrice(parts[4]),
			Volume: int64(s.parsePrice(parts[5])),
		})
	}

	s.calculateIndicators(klines)

	if len(klines) > limit {
		return klines[len(klines)-limit:], nil
	}
	return klines, nil
}

func (s *StockService) calculateIndicators(klines []*models.KLineData) {
	if len(klines) == 0 {
		return
	}

	ema12 := klines[0].Close
	ema26 := klines[0].Close
	dea := 0.0
	for i, k := range klines {
		ema12 = ema12*11/13 + k.Close*2/13
		ema26 = ema26*25/27 + k.Close*2/27
		dif := ema12 - ema26
		if i == 0 {
			dea = dif
		} else {
			dea = dea*8/10 + dif*2/10
		}
		k.MACD = &models.MACD{
			DIF: dif,
			DEA: dea,
			Bar: (dif - dea) * 2,
		}
	}

	for i := 0; i < len(klines); i++ {
		if i < 8 {
			klines[i].KDJ = &models.KDJ{K: 50, D: 50, J: 50}
			continue
		}
		low := klines[i].Low
		high := klines[i].High
		for j := i - 8; j < i; j++ {
			low = math.Min(low, klines[j].Low)
			high = math.Max(high, klines[j].High)
		}
		rsv := 0.0
		if high != low {
			rsv = (klines[i].Close - low) / (high - low) * 100
		}
		prevK := klines[i-1].KDJ.K
		prevD := klines[i-1].KDJ.D
		k := prevK*2/3 + rsv/3
		d := prevD*2/3 + k/3
		klines[i].KDJ = &models.KDJ{
			K: k,
			D: d,
			J: 3*k - 2*d,
		}
	}

	if len(klines) > 14 {
		for i := 14; i < len(klines); i++ {
			upSum := 0.0
			downSum := 0.0
			for j := i - 13; j <= i; j++ {
				diff := klines[j].Close - klines[j-1].Close
				if diff > 0 {
					upSum += diff
				} else {
					downSum -= diff
				}
			}
			if upSum+downSum == 0 {
				klines[i].RSI = 50
			} else {
				klines[i].RSI = upSum / (upSum + downSum) * 100
			}
		}
	}
}

func (s *StockService) getSecID(code string) string {
	if len(code) != 6 {
		return ""
	}
	if strings.HasPrefix(code, "6") || strings.HasPrefix(code, "9") {
		return "1." + code
	}
	if strings.HasPrefix(code, "0") || strings.HasPrefix(code, "3") || 
	   strings.HasPrefix(code, "4") || strings.HasPrefix(code, "8") ||
	   strings.HasPrefix(code, "2") {
		return "0." + code
	}
	return ""
}

func (s *StockService) parsePrice(sVal string) float64 {
	var f float64
	fmt.Sscanf(sVal, "%f", &f)
	return f
}

func (s *StockService) SearchStock(keyword string) ([]*models.StockData, error) {
	keyword = strings.TrimSpace(keyword)
	if len(keyword) == 6 {
		stock, err := s.GetStockByCode(keyword)
		if err == nil {
			return []*models.StockData{stock}, nil
		}
	}
	return s.SearchStockLegacy(keyword)
}

func (s *StockService) SearchStockLegacy(keyword string) ([]*models.StockData, error) {
	url := fmt.Sprintf("%s?pn=1&pz=1000&po=1&np=1&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:13,m:0+t:80,m:1+t:2,m:1+t:23&fields=f12,f14", s.listURL)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var apiResp models.EastMoneyResponse
	json.Unmarshal(body, &apiResp)
	
	keyword = strings.ToLower(keyword)
	results := make([]*models.StockData, 0)
	for _, diff := range apiResp.Data.Diff {
		if strings.Contains(strings.ToLower(diff.F12), keyword) ||
			strings.Contains(strings.ToLower(diff.F14), keyword) {
			results = append(results, diff.ToStockData())
		}
	}
	return results, nil
}

// 辅助解析函数
func getString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getFloat(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

func getInt64(v interface{}) int64 {
	if f, ok := v.(float64); ok {
		return int64(f)
	}
	return 0
}
