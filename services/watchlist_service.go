package services

import (
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"
)

// WatchlistService 业务逻辑层
type WatchlistService struct {
	repo repositories.WatchlistRepository
}

// NewWatchlistService 构造函数，接受 WatchlistRepository 接口
func NewWatchlistService(repo repositories.WatchlistRepository) *WatchlistService {
	return &WatchlistService{repo: repo}
}

func (s *WatchlistService) AddToWatchlist(stock *models.StockData) error {
	// 业务逻辑：这里可以添加如“自选股数量限制”等业务规则
	return s.repo.Add(stock)
}

func (s *WatchlistService) RemoveFromWatchlist(code string) error {
	// 业务逻辑：这里可以添加如“移除前检查是否有预警”等业务规则
	return s.repo.Remove(code)
}

func (s *WatchlistService) GetWatchlist() ([]*models.StockData, error) {
	// 业务逻辑：这里可以添加如“按用户偏好排序”等业务规则
	return s.repo.GetAll()
}
