package controllers

import (
	"stock-analyzer-wails/services"
)

// StockMarketController 市场股票控制器
type StockMarketController struct {
	stockMarketService *services.StockMarketService
}

// NewStockMarketController 创建市场股票控制器
func NewStockMarketController(stockMarketService *services.StockMarketService) *StockMarketController {
	return &StockMarketController{
		stockMarketService: stockMarketService,
	}
}

// SyncAllStocks 同步所有市场股票
func (c *StockMarketController) SyncAllStocks() (*services.SyncStocksResult, error) {
	return c.stockMarketService.SyncAllStocks()
}

// GetStocksListRequest 获取股票列表请求
type GetStocksListRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Search   string `json:"search"`
}

// GetStocksList 获取股票列表
func (c *StockMarketController) GetStocksList(page int, pageSize int, search string, industry string) (interface{}, error) {
	stocks, total, err := c.stockMarketService.GetStocksList(page, pageSize, search, industry)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"stocks": stocks,
		"total":  total,
		"page":   page,
		"pageSize": pageSize,
	}, nil
}

// GetIndustries 获取行业列表
func (c *StockMarketController) GetIndustries() ([]services.IndustryInfo, error) {
	return c.stockMarketService.GetIndustries()
}

// GetSyncStats 获取同步统计信息
func (c *StockMarketController) GetSyncStats() (interface{}, error) {
	return c.stockMarketService.GetSyncStats()
}
