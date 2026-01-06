package controllers

import (
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/services"
)

// WatchlistController 负责处理前端对自选股的操作请求
type WatchlistController struct {
	service *services.WatchlistService
}

// NewWatchlistController 构造函数
func NewWatchlistController(svc *services.WatchlistService) *WatchlistController {
	return &WatchlistController{service: svc}
}

// AddToWatchlist Wails 绑定方法：添加股票到自选股
func (c *WatchlistController) AddToWatchlist(stock models.StockData) error {
	// Controller 职责：参数校验（如果需要）
	if stock.Code == "" {
		return models.ErrInvalidInput
	}
	// 调用 Service 层
	return c.service.AddToWatchlist(&stock)
}

// RemoveFromWatchlist Wails 绑定方法：从自选股移除股票
func (c *WatchlistController) RemoveFromWatchlist(code string) error {
	if code == "" {
		return models.ErrInvalidInput
	}
	return c.service.RemoveFromWatchlist(code)
}

// GetWatchlist Wails 绑定方法：获取自选股列表
func (c *WatchlistController) GetWatchlist() ([]*models.StockData, error) {
	return c.service.GetWatchlist()
}
