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
	// 清理股票代码（去除空格等）
	code = strings.TrimSpace(code)
	
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
	
	// 创建请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	
	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "http://quote.eastmoney.com/")
	
	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	
	// 解析JSON响应
	var apiResp models.EastMoneyResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	
	// 查找匹配的股票
	for _, diff := range apiResp.Data.Diff {
		if diff.F12 == code {
			return diff.ToStockData(), nil
		}
	}
	
	return nil, fmt.Errorf("未找到股票代码: %s", code)
}

// GetStockList 获取股票列表（可选功能）
func (s *StockService) GetStockList(pageNum, pageSize int) ([]*models.StockData, error) {
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
	
	// 创建请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	
	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "http://quote.eastmoney.com/")
	
	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	
	// 解析JSON响应
	var apiResp models.EastMoneyResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	
	// 转换为标准股票数据
	stocks := make([]*models.StockData, 0, len(apiResp.Data.Diff))
	for _, diff := range apiResp.Data.Diff {
		stocks = append(stocks, diff.ToStockData())
	}
	
	return stocks, nil
}

// SearchStock 搜索股票（支持代码和名称模糊搜索）
func (s *StockService) SearchStock(keyword string) ([]*models.StockData, error) {
	// 获取所有股票数据
	allStocks, err := s.GetStockList(1, 5000)
	if err != nil {
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
	
	return results, nil
}
