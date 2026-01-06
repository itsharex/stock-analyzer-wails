package services

import (
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"
)

// PositionService 业务逻辑层，负责持仓记录的业务处理
type PositionService struct {
	repo repositories.PositionRepository
}

// NewPositionService 构造函数
func NewPositionService(repo repositories.PositionRepository) *PositionService {
	return &PositionService{repo: repo}
}

// SavePosition 保存或更新持仓记录
func (s *PositionService) SavePosition(pos *models.Position) error {
	// 业务逻辑：这里可以添加如“持仓数量限制”等业务规则
	return s.repo.SavePosition(pos)
}

// GetPositions 获取所有持仓记录
func (s *PositionService) GetPositions() (map[string]*models.Position, error) {
	return s.repo.GetPositions()
}

// RemovePosition 移除持仓记录
func (s *PositionService) RemovePosition(code string) error {
	// 业务逻辑：这里可以添加如“移除前检查是否已平仓”等业务规则
	return s.repo.RemovePosition(code)
}

// NewPositionStorageService 兼容旧的命名，但返回新的 PositionService
func NewPositionStorageService(dbSvc *DBService) *PositionService {
	repo := repositories.NewSQLitePositionRepository(dbSvc.GetDB())
	return NewPositionService(repo)
}
