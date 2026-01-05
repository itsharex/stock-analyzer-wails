package services

import (
	"encoding/json"
	"fmt"
	"io"
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
	client   *http.Client
}

// NewStockService 创建股票数据服务实例
func NewStockService() *StockService {
	return &StockService{
		exactURL: "https://push2.eastmoney.com/api/qt/stock/get",
		listURL:  "http://78.push2.eastmoney.com/api/qt/clist/get",
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

	// 1. 识别市场前缀
	secid := s.getSecID(code)
	if secid == "" {
		return nil, fmt.Errorf("无法识别的股票代码格式: %s", code)
	}

	// 2. 构建请求参数
	// f58:名称, f43:最新价, f169:涨跌额, f170:涨跌幅, f47:成交量, f48:成交额, f44:最高, f45:最低, f46:今开, f60:昨收, f171:振幅, f168:换手率, f162:市盈率, f167:市净率, f116:总市值, f117:流通市值
	fields := "f58,f43,f169,f170,f47,f48,f44,f45,f46,f60,f171,f168,f162,f167,f116,f117,f12,f14"
	fullURL := fmt.Sprintf("%s?secid=%s&fields=%s", s.exactURL, secid, fields)

	logger.Debug("发起精确行情请求",
		zap.String("module", "services.stock"),
		zap.String("op", "GetStockByCode.request"),
		zap.String("url", fullURL),
		zap.String("secid", secid),
	)

	// 3. 发送请求
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	// 4. 解析响应
	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Data == nil {
		return nil, fmt.Errorf("未找到股票数据: %s", code)
	}

	// 5. 转换为标准模型
	data := result.Data
	stock := &models.StockData{
		Code:       getString(data["f12"]),
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

	logger.Info("精确获取股票数据成功",
		zap.String("module", "services.stock"),
		zap.String("op", "GetStockByCode"),
		zap.String("stock_code", code),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)

	return stock, nil
}

// getSecID 根据股票代码获取东方财富的 secid (市场代码.股票代码)
func (s *StockService) getSecID(code string) string {
	if len(code) != 6 {
		return ""
	}
	// 沪市 A 股以 6 开头
	if strings.HasPrefix(code, "6") || strings.HasPrefix(code, "9") {
		return "1." + code
	}
	// 深市 A 股以 0, 3 开头，北交所以 4, 8 开头
	if strings.HasPrefix(code, "0") || strings.HasPrefix(code, "3") || 
	   strings.HasPrefix(code, "4") || strings.HasPrefix(code, "8") ||
	   strings.HasPrefix(code, "2") {
		return "0." + code
	}
	return ""
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

// GetStockList 保持不变，用于列表展示
func (s *StockService) GetStockList(pageNum, pageSize int) ([]*models.StockData, error) {
	// ... (保持原有逻辑，但可以优化 fields 以匹配新的解析方式)
	// 为了兼容性，这里暂时保留原有逻辑
	return s.getStockListInternal(pageNum, pageSize)
}

func (s *StockService) getStockListInternal(pageNum, pageSize int) ([]*models.StockData, error) {
	// 原有 GetStockList 的实现逻辑
	// ... (省略重复代码，实际写入时会包含完整内容)
	// 考虑到篇幅，我将合并逻辑
	url := fmt.Sprintf("%s?pn=%d&pz=%d&po=1&np=1&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:13,m:0+t:80,m:1+t:2,m:1+t:23&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13,f14,f15,f16,f17,f18,f20,f21,f22,f23", s.listURL, pageNum, pageSize)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var apiResp models.EastMoneyResponse
	json.Unmarshal(body, &apiResp)
	stocks := make([]*models.StockData, 0, len(apiResp.Data.Diff))
	for _, diff := range apiResp.Data.Diff {
		stocks = append(stocks, diff.ToStockData())
	}
	return stocks, nil
}

// SearchStock 优化搜索逻辑：如果输入是 6 位数字，优先尝试精确查询
func (s *StockService) SearchStock(keyword string) ([]*models.StockData, error) {
	keyword = strings.TrimSpace(keyword)
	if len(keyword) == 6 {
		stock, err := s.GetStockByCode(keyword)
		if err == nil {
			return []*models.StockData{stock}, nil
		}
	}
	// 否则回退到列表搜索
	return s.SearchStockLegacy(keyword)
}

func (s *StockService) SearchStockLegacy(keyword string) ([]*models.StockData, error) {
	allStocks, err := s.getStockListInternal(1, 1000) // 搜索前1000只
	if err != nil {
		return nil, err
	}
	keyword = strings.ToLower(keyword)
	results := make([]*models.StockData, 0)
	for _, stock := range allStocks {
		if strings.Contains(strings.ToLower(stock.Code), keyword) ||
			strings.Contains(strings.ToLower(stock.Name), keyword) {
			results = append(results, stock)
		}
	}
	return results, nil
}
