package controllers

import (
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"
)

// SyncHistoryController 负责处理前端对同步历史的操作请求
type SyncHistoryController struct {
	repository repositories.SyncHistoryRepository
}

// NewSyncHistoryController 构造函数
func NewSyncHistoryController(repo repositories.SyncHistoryRepository) *SyncHistoryController {
	return &SyncHistoryController{repository: repo}
}

// AddSyncHistory 添加同步历史记录
func (c *SyncHistoryController) AddSyncHistory(history models.SyncHistory) error {
	return c.repository.Add(&history)
}

// GetAllSyncHistory 获取所有同步历史记录（分页）
func (c *SyncHistoryController) GetAllSyncHistory(limit int, offset int) ([]*models.SyncHistory, error) {
	return c.repository.GetAll(limit, offset)
}

// GetSyncHistoryByCode 根据股票代码获取同步历史记录
func (c *SyncHistoryController) GetSyncHistoryByCode(code string, limit int) ([]*models.SyncHistory, error) {
	return c.repository.GetByStockCode(code, limit)
}

// GetSyncHistoryCount 获取同步历史记录总数
func (c *SyncHistoryController) GetSyncHistoryCount() (int, error) {
	return c.repository.GetCount()
}

// ClearAllSyncHistory 清除所有同步历史记录
func (c *SyncHistoryController) ClearAllSyncHistory() error {
	return c.repository.ClearAll()
}
