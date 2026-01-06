package controllers

import (
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/services"
)

// PositionController 负责处理前端对持仓的操作请求
type PositionController struct {
	service *services.PositionService
}

// NewPositionController 构造函数
func NewPositionController(svc *services.PositionService) *PositionController {
	return &PositionController{service: svc}
}

// AddPosition Wails 绑定方法：添加持仓记录
func (c *PositionController) AddPosition(pos models.Position) error {
	// Controller 职责：参数校验（如果需要）
	if pos.StockCode == "" {
		return models.ErrInvalidInput
	}
	// 调用 Service 层
	return c.service.SavePosition(&pos)
}

// GetPositions Wails 绑定方法：获取所有活跃持仓
func (c *PositionController) GetPositions() (map[string]*models.Position, error) {
	return c.service.GetPositions()
}

// RemovePosition Wails 绑定方法：移除持仓记录
func (c *PositionController) RemovePosition(code string) error {
	if code == "" {
		return models.ErrInvalidInput
	}
	return c.service.RemovePosition(code)
}
