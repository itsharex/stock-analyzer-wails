package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"stock-analyzer-wails/models"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"stock-analyzer-wails/internal/logger"

	"github.com/wailsapp/wails/v2/pkg/runtime"
		"go.uber.org/zap"
)

// StockService 股票数据服务
type StockService struct {
	exactURL string
	listURL  string
	klineURL string
	client   *http.Client
}

// GetStockDetail 获取个股详情页所需的所有数据
func (s *StockService) GetStockDetail(code string) (*models.StockDetail, error) {
	// 1. 获取基础行情数据 (包含 PE/PB)
	baseStock, err := s.GetStockByCode(code)
	if err != nil {
		return nil, err
	}

	detail := &models.StockDetail{
		StockData: *baseStock,
	}

	// 2. 获取实时盘口数据 (五档)
	orderBook, err := s.getOrderBook(code)
	if err == nil {
		detail.OrderBook = *orderBook
	} else {
		logger.Error("获取盘口数据失败", zap.Error(err))
	}

	// 3. 获取财务数据 (ROE, 净利润增长率等)
	financial, err := s.getFinancialSummary(code)
	if err == nil {
		detail.Financial = *financial
	} else {
		logger.Error("获取财务数据失败", zap.Error(err))
	}

	// 4. 获取行业数据
	industry, err := s.getIndustryInfo(code)
	if err == nil {
		detail.Industry = *industry
	} else {
		logger.Error("获取行业数据失败", zap.Error(err))
	}

	return detail, nil
}

// getOrderBook 获取五档盘口数据
func (s *StockService) getOrderBook(code string) (*models.OrderBook, error) {
	secid := s.getSecID(code)
	if secid == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	// 东方财富 qt/stock/get 接口的 fields 字段中包含盘口数据
	// 盘口字段: f68-f85 (买卖五档)
	fields := "f68,f69,f70,f71,f72,f73,f74,f75,f76,f77,f78,f79,f80,f81,f82,f83,f84,f85"
	fullURL := fmt.Sprintf("%s?secid=%s&fields=%s", s.exactURL, secid, fields)

	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求盘口数据失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析盘口响应失败: %w", err)
	}

	if result.Data == nil {
		return nil, fmt.Errorf("未找到盘口数据: %s", code)
	}

	data := result.Data
	
	orderBook := &models.OrderBook{
		Buy5: make([]models.OrderBookEntry, 5),
		Sell5: make([]models.OrderBookEntry, 5),
	}

	// 卖盘 (Sell1-Sell5: f73-f82)
	// f73: 卖一价, f74: 卖一量, f75: 卖二价, f76: 卖二量, ..., f81: 卖五价, f82: 卖五量
	// 注意：东方财富返回的价格是整数，需要除以 100
	for i := 0; i < 5; i++ {
		priceField := fmt.Sprintf("f%d", 73 + i*2)
		volumeField := fmt.Sprintf("f%d", 74 + i*2)
		orderBook.Sell5[i] = models.OrderBookEntry{
			Price: getFloat(data[priceField]) / 100,
			Volume: getInt64(data[volumeField]),
		}
	}

	// 买盘 (Buy1-Buy5: f68-f77)
	// f68: 买一价, f69: 买一量, f70: 买二价, f71: 买二量, ..., f77: 买五价, f78: 买五量
	for i := 0; i < 5; i++ {
		priceField := fmt.Sprintf("f%d", 68 + i*2)
		volumeField := fmt.Sprintf("f%d", 69 + i*2)
		orderBook.Buy5[i] = models.OrderBookEntry{
			Price: getFloat(data[priceField]) / 100,
			Volume: getInt64(data[volumeField]),
		}
	}

	return orderBook, nil
}

// getFinancialSummary 获取核心财务数据 (Mock 数据)
func (s *StockService) getFinancialSummary(code string) (*models.FinancialSummary, error) {
	// TODO: 后续可在此处集成真实的财务数据 API。
	// 建议使用如东方财富数据中心 (datacenter-web.eastmoney.com) 的 API，
	// 需找到正确的 reportName (例如 RPT_F10_MAIN_INDICATOR_DATA 或类似名称)
	// 并解析返回的 JSON 数据。
	// 
	// 示例 API 结构 (需要自行调研):
	// url := "https://datacenter-web.eastmoney.com/api/data/v1/get?reportName=RPT_FINANCE_MAIN_INDICATOR&columns=REPORT_DATE,ROE,NETPROFIT_GROWTH_RATE,GROSS_PROFIT_MARGIN,TOTAL_MARKET_VALUE,CIRCULATING_MARKET_VALUE,DIVIDEND_YIELD&filter=(SECURITY_CODE=\"" + code + "\")&pageNumber=1&pageSize=1"

	// 使用 Mock 数据
	return &models.FinancialSummary{
		ROE: 15.8,
		NetProfitGrowthRate: 22.5,
		GrossProfitMargin: 45.1,
		TotalMarketValue: 12000.5, // 亿元
		CirculatingMarketValue: 8000.2, // 亿元
		DividendYield: 1.5,
		ReportDate: time.Now().AddDate(0, -3, 0),
	}, nil
}

// getIndustryInfo 获取行业与宏观信息 (Mock 数据)
func (s *StockService) getIndustryInfo(code string) (*models.IndustryInfo, error) {
	// TODO: 后续可在此处集成真实的行业和概念板块 API。
	// 建议使用如东方财富数据中心 (datacenter-web.eastmoney.com) 的 API，
	// 需找到正确的 reportName (例如 RPT_F10_CONCEPT 或类似名称)
	// 并解析返回的 JSON 数据。
	// 
	// 示例 API 结构 (需要自行调研):
	// url := "https://datacenter-web.eastmoney.com/api/data/v1/get?reportName=RPT_F10_CONCEPT&columns=CONCEPT_NAME,CONCEPT_CODE&filter=(SECURITY_CODE=\"" + code + "\")&pageNumber=1&pageSize=50"

	// 使用 Mock 数据
	return &models.IndustryInfo{
		IndustryName: "软件开发",
		ConceptNames: []string{"人工智能", "云计算", "数字经济"},
		IndustryPE: 45.8,
	}, nil
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

	// 增加盘口相关字段 f19(委比), f20(量比)
	fields := "f58,f43,f169,f170,f47,f48,f44,f45,f46,f60,f171,f168,f162,f167,f116,f117,f12,f14,f19,f20"
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
			VolumeRatio: getFloat(data["f20"]) / 100,
			WarrantRatio: getFloat(data["f19"]) / 100,
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

func (s *StockService) GetIntradayData(code string) (*models.IntradayResponse, error) {
	secid := s.getSecID(code)
	if secid == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	url := fmt.Sprintf("https://push2.eastmoney.com/api/qt/stock/trends2/get?secid=%s&fields1=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13&fields2=f51,f52,f53,f54,f55,f56,f57,f58&ndays=1&iscr=0&ut=fa5fd1943c7b386f172d6893dbfba10b", secid)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data struct {
			PreClose float64  `json:"preClose"`
			Trends   []string `json:"trends"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	intradayData := make([]models.IntradayData, 0)
	for _, line := range result.Data.Trends {
		parts := strings.Split(line, ",")
		if len(parts) < 5 {
			continue
		}
		intradayData = append(intradayData, models.IntradayData{
			Time:     parts[0],
			Price:    s.parsePrice(parts[2]),
			AvgPrice: s.parsePrice(parts[7]),
			Volume:   int64(s.parsePrice(parts[5])),
			PreClose: result.Data.PreClose,
		})
	}

	return &models.IntradayResponse{
		Data:     intradayData,
		PreClose: result.Data.PreClose,
	}, nil
}

func (s *StockService) GetMoneyFlowData(code string) (*models.MoneyFlowResponse, error) {
	secid := s.getSecID(code)
	if secid == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	url := fmt.Sprintf("https://push2.eastmoney.com/api/qt/stock/fflow/kline/get?secid=%s&fields1=f1,f2,f3,f7&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65&klt=1&lmt=240", secid)

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

	moneyFlows := make([]models.MoneyFlowData, 0)
	var totalMain, totalRetail float64

	for _, line := range result.Data.Klines {
		parts := strings.Split(line, ",")
		if len(parts) < 13 {
			continue
		}
		
		mainNet := s.parsePrice(parts[1])
		small := s.parsePrice(parts[2])
		medium := s.parsePrice(parts[3])
		large := s.parsePrice(parts[4])
		superLarge := s.parsePrice(parts[5])

		moneyFlows = append(moneyFlows, models.MoneyFlowData{
			Time:       parts[0],
			SuperLarge: superLarge,
			Large:      large,
			Medium:     medium,
			Small:      small,
			MainNet:    mainNet,
		})
		
		totalMain += mainNet
		totalRetail += small
	}

	if len(moneyFlows) > 10 {
		for i := 10; i < len(moneyFlows); i++ {
			var sumAbsMain float64
			for j := i - 10; j < i; j++ {
				val := moneyFlows[j].MainNet
				if val < 0 {
					val = -val
				}
				sumAbsMain += val
			}
			avgAbsMain := sumAbsMain / 10
			if avgAbsMain < 10000 {
				avgAbsMain = 10000
			}

			currentMain := moneyFlows[i].MainNet
			if currentMain > avgAbsMain*5 {
				moneyFlows[i].Signal = "扫货"
			} else if currentMain < -avgAbsMain*5 {
				moneyFlows[i].Signal = "砸盘"
			}
		}
	}

	status := "平稳运行"
	description := "当前资金进出相对平衡，建议关注趋势确认。"

	if len(moneyFlows) > 0 {
		lastFlow := moneyFlows[len(moneyFlows)-1]
		if totalMain > 0 && lastFlow.MainNet > 0 {
			status = "主力建仓"
			description = "主力资金正在持续流入，且单笔成交金额较大，说明机构看好后市，正在积极吸筹。"
		}
		if totalMain < 0 && totalRetail > 0 {
			status = "散户追高"
			description = "当前股价上涨主要由散户情绪推动，主力资金正在趁高点派发筹码，请警惕冲高回落风险。"
		}
		if totalMain > 0 && totalRetail < 0 {
			status = "机构洗盘"
			description = "主力资金在股价回调时默默接盘，散户因恐慌抛售，这通常是拉升前的洗盘行为。"
		}
	}

	return &models.MoneyFlowResponse{
		Data:        moneyFlows,
		TodayMain:   totalMain,
		TodayRetail: totalRetail,
		Status:      status,
		Description: description,
	}, nil
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

func (s *StockService) GetStockHealthCheck(code string) (*models.HealthCheckResult, error) {
	stock, err := s.GetStockByCode(code)
	if err != nil {
		return nil, err
	}

	items := make([]models.HealthItem, 0)
	score := 100

	peStatus := "正常"
	peDesc := fmt.Sprintf("当前市盈率 %.2f，处于行业合理区间。", stock.PE)
	if stock.PE < 0 {
		peStatus = "异常"
		peDesc = "公司目前处于亏损状态，财务基本面存在较大不确定性。"
		score -= 20
	} else if stock.PE > 100 {
		peStatus = "警告"
		peDesc = "市盈率过高，估值存在泡沫风险，需警惕回调。"
		score -= 10
	}
	items = append(items, models.HealthItem{
		Category: "财务",
		Name:     "估值水平",
		Value:    fmt.Sprintf("%.2f PE", stock.PE),
		Status:   peStatus,
		Description: peDesc,
	})

	turnoverStatus := "正常"
	turnoverDesc := "换手率适中，交投活跃度正常。"
	if stock.Turnover > 15 {
		turnoverStatus = "警告"
		turnoverDesc = "换手率极高，说明筹码松动，主力可能在进行高位派发。"
		score -= 15
	} else if stock.Turnover < 0.5 {
		turnoverStatus = "警告"
		turnoverDesc = "成交极其低迷，属于“僵尸股”，流动性风险较大。"
		score -= 10
	}
	items = append(items, models.HealthItem{
		Category: "资金",
		Name:     "流动性检测",
		Value:    fmt.Sprintf("%.2f%%", stock.Turnover),
		Status:   turnoverStatus,
		Description: turnoverDesc,
	})

	ampStatus := "正常"
	ampDesc := "日内波动在正常范围内。"
	if stock.Amplitude > 10 {
		ampStatus = "警告"
		ampDesc = "日内振幅巨大，多空分歧严重，短期波动风险极高。"
		score -= 10
	}
	items = append(items, models.HealthItem{
		Category: "技术",
		Name:     "波动风险",
		Value:    fmt.Sprintf("%.2f%%", stock.Amplitude),
		Status:   ampStatus,
		Description: ampDesc,
	})

	status, riskLevel := "健康", "低"
	if score < 60 {
		status, riskLevel = "风险", "高"
	} else if score < 85 {
		status, riskLevel = "亚健康", "中"
	}

	return &models.HealthCheckResult{
		Score:     score,
		Status:    status,
		RiskLevel: riskLevel,
		Summary:   fmt.Sprintf("%s目前综合评分为 %d 分，%s", stock.Name, score, getHealthSummary(score)),
		Items:     items,
		UpdatedAt: time.Now().Format("15:04:05"),
	}, nil
}

func getHealthSummary(score int) string {
	if score >= 85 {
		return "基本面和技术面表现稳健，适合中长期关注。"
	} else if score >= 60 {
		return "存在部分指标异常，建议控制仓位，关注关键支撑位。"
	}
	return "多项指标触发预警，建议新手坚决回避，等待风险释放。"
}

func (s *StockService) BatchAnalyzeStocks(ctx context.Context, codes []string, role string, aiSvc *AIService) error {
	total := len(codes)
	var completed int32 = 0
	concurrency := 2
	jobs := make(chan string, total)
	var wg sync.WaitGroup

	for w := 1; w <= concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for code := range jobs {
				stock, err := s.GetStockByCode(code)
				if err != nil {
					continue
				}
				klines, err := s.GetKLineData(code, 100, "daily")
				if err != nil {
					continue
				}
				_, _ = aiSvc.AnalyzeTechnical(stock, klines, "daily", role)
				newCompleted := atomic.AddInt32(&completed, 1)
				runtime.EventsEmit(ctx, "batch_analyze_progress", map[string]interface{}{
					"code":      code,
					"name":      stock.Name,
					"completed": newCompleted,
					"total":     total,
					"percent":   float64(newCompleted) / float64(total) * 100,
				})
			}
		}()
	}

	for _, code := range codes {
		jobs <- code
	}
	close(jobs)
	wg.Wait()
	return nil
}

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
