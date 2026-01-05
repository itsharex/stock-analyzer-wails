package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"stock-analyzer-wails/models"
	"strings"
	"time"

	"stock-analyzer-wails/internal/logger"

	"go.uber.org/zap"
)

// StockService 股票数据服务
type StockService struct {
	baseURL string
	client  *http.Client
}

// NewStockService 创建股票数据服务实例
func NewStockService() *StockService {
	return &StockService{
		baseURL: "http://78.push2.eastmoney.com/api/qt/clist/get",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetStockByCode 根据股票代码获取股票数据
func (s *StockService) GetStockByCode(code string) (*models.StockData, error) {
	start := time.Now()
	// 清理股票代码（去除空格等）
	code = strings.TrimSpace(code)
	if code == "" {
		logger.Warn("股票代码为空",
			zap.String("module", "services.stock"),
			zap.String("op", "GetStockByCode"),
		)
		return nil, fmt.Errorf("股票代码不能为空")
	}

	// 构建请求参数
	params := url.Values{}
	params.Add("pn", "1")
	params.Add("pz", "5000")
	params.Add("po", "1")
	params.Add("np", "1")
	params.Add("fltt", "2")
	params.Add("invt", "2")
	params.Add("fid", "f3")
	params.Add("fs", "m:0+t:6,m:0+t:13,m:0+t:80,m:1+t:2,m:1+t:23") // 沪深A股
	params.Add("fields", "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13,f14,f15,f16,f17,f18,f20,f21,f22,f23")

	// 构建完整URL
	fullURL := fmt.Sprintf("%s?%s", s.baseURL, params.Encode())
	logger.Debug("发起行情请求",
		zap.String("module", "services.stock"),
		zap.String("op", "GetStockByCode.request"),
		zap.String("url", fullURL),
		zap.String("stock_code", code),
	)

	// 创建请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		logger.Error("创建请求失败",
			zap.String("module", "services.stock"),
			zap.String("op", "GetStockByCode.http.NewRequest"),
			zap.String("url", fullURL),
			zap.String("stock_code", code),
			zap.Error(err),
		)
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "http://quote.eastmoney.com/")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		logger.Error("请求失败",
			zap.String("module", "services.stock"),
			zap.String("op", "GetStockByCode.client.Do"),
			zap.String("url", fullURL),
			zap.String("stock_code", code),
			zap.Int64("timeout_ms", s.client.Timeout.Milliseconds()),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取响应失败",
			zap.String("module", "services.stock"),
			zap.String("op", "GetStockByCode.io.ReadAll"),
			zap.String("url", fullURL),
			zap.String("stock_code", code),
			zap.Int("http_status", resp.StatusCode),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		logger.Warn("行情接口返回非 200",
			zap.String("module", "services.stock"),
			zap.String("op", "GetStockByCode.httpStatus"),
			zap.String("url", fullURL),
			zap.String("stock_code", code),
			zap.Int("http_status", resp.StatusCode),
			zap.Int("body_size", len(body)),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
	}

	// 解析JSON响应
	var apiResp models.EastMoneyResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		logger.Error("解析响应失败",
			zap.String("module", "services.stock"),
			zap.String("op", "GetStockByCode.json.Unmarshal"),
			zap.String("url", fullURL),
			zap.String("stock_code", code),
			zap.Int("http_status", resp.StatusCode),
			zap.Int("body_size", len(body)),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 查找匹配的股票
	for _, diff := range apiResp.Data.Diff {
		if diff.F12 == code {
			logger.Info("获取股票数据成功",
				zap.String("module", "services.stock"),
				zap.String("op", "GetStockByCode"),
				zap.String("stock_code", code),
				zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			)
			return diff.ToStockData(), nil
		}
	}

	logger.Warn("未找到股票代码",
		zap.String("module", "services.stock"),
		zap.String("op", "GetStockByCode.notFound"),
		zap.String("stock_code", code),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return nil, fmt.Errorf("未找到股票代码: %s", code)
}

// GetStockList 获取股票列表（可选功能）
func (s *StockService) GetStockList(pageNum, pageSize int) ([]*models.StockData, error) {
	start := time.Now()
	// 构建请求参数
	params := url.Values{}
	params.Add("pn", fmt.Sprintf("%d", pageNum))
	params.Add("pz", fmt.Sprintf("%d", pageSize))
	params.Add("po", "1")
	params.Add("np", "1")
	params.Add("fltt", "2")
	params.Add("invt", "2")
	params.Add("fid", "f3")
	params.Add("fs", "m:0+t:6,m:0+t:13,m:0+t:80,m:1+t:2,m:1+t:23")
	params.Add("fields", "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13,f14,f15,f16,f17,f18,f20,f21,f22,f23")

	// 构建完整URL
	fullURL := fmt.Sprintf("%s?%s", s.baseURL, params.Encode())
	logger.Debug("发起股票列表请求",
		zap.String("module", "services.stock"),
		zap.String("op", "GetStockList.request"),
		zap.String("url", fullURL),
		zap.Int("page_num", pageNum),
		zap.Int("page_size", pageSize),
	)

	// 创建请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		logger.Error("创建请求失败",
			zap.String("module", "services.stock"),
			zap.String("op", "GetStockList.http.NewRequest"),
			zap.String("url", fullURL),
			zap.Int("page_num", pageNum),
			zap.Int("page_size", pageSize),
			zap.Error(err),
		)
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "http://quote.eastmoney.com/")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		logger.Error("请求失败",
			zap.String("module", "services.stock"),
			zap.String("op", "GetStockList.client.Do"),
			zap.String("url", fullURL),
			zap.Int("page_num", pageNum),
			zap.Int("page_size", pageSize),
			zap.Int64("timeout_ms", s.client.Timeout.Milliseconds()),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取响应失败",
			zap.String("module", "services.stock"),
			zap.String("op", "GetStockList.io.ReadAll"),
			zap.String("url", fullURL),
			zap.Int("http_status", resp.StatusCode),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		logger.Warn("股票列表接口返回非 200",
			zap.String("module", "services.stock"),
			zap.String("op", "GetStockList.httpStatus"),
			zap.String("url", fullURL),
			zap.Int("http_status", resp.StatusCode),
			zap.Int("body_size", len(body)),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
	}

	// 解析JSON响应
	var apiResp models.EastMoneyResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		logger.Error("解析响应失败",
			zap.String("module", "services.stock"),
			zap.String("op", "GetStockList.json.Unmarshal"),
			zap.String("url", fullURL),
			zap.Int("http_status", resp.StatusCode),
			zap.Int("body_size", len(body)),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 转换为标准股票数据
	stocks := make([]*models.StockData, 0, len(apiResp.Data.Diff))
	for _, diff := range apiResp.Data.Diff {
		stocks = append(stocks, diff.ToStockData())
	}

	logger.Info("获取股票列表成功",
		zap.String("module", "services.stock"),
		zap.String("op", "GetStockList"),
		zap.Int("count", len(stocks)),
		zap.Int("page_num", pageNum),
		zap.Int("page_size", pageSize),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return stocks, nil
}

// SearchStock 搜索股票（支持代码和名称模糊搜索）
func (s *StockService) SearchStock(keyword string) ([]*models.StockData, error) {
	start := time.Now()
	// 获取所有股票数据
	allStocks, err := s.GetStockList(1, 5000)
	if err != nil {
		logger.Error("获取股票列表失败（用于搜索）",
			zap.String("module", "services.stock"),
			zap.String("op", "SearchStock.GetStockList"),
			zap.String("keyword", keyword),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, err
	}

	// 过滤匹配的股票
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	results := make([]*models.StockData, 0)

	for _, stock := range allStocks {
		if strings.Contains(strings.ToLower(stock.Code), keyword) ||
			strings.Contains(strings.ToLower(stock.Name), keyword) {
			results = append(results, stock)
		}
	}

	logger.Info("搜索股票完成",
		zap.String("module", "services.stock"),
		zap.String("op", "SearchStock"),
		zap.String("keyword", keyword),
		zap.Int("result_count", len(results)),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return results, nil
}
