package controllers

import (
	"stock-analyzer-wails/services"
)

// StrategyController 策略控制器
type StrategyController struct {
	strategyService *services.StrategyService
}

// NewStrategyController 创建策略控制器
func NewStrategyController(strategyService *services.StrategyService) *StrategyController {
	return &StrategyController{
		strategyService: strategyService,
	}
}

// CreateStrategyRequest 创建策略请求
type CreateStrategyRequest struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	StrategyType string                 `json:"strategyType"`
	Parameters   map[string]interface{} `json:"parameters"`
}

// UpdateStrategyRequest 更新策略请求
type UpdateStrategyRequest struct {
	ID           int64                  `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	StrategyType string                 `json:"strategyType"`
	Parameters   map[string]interface{} `json:"parameters"`
}

// CreateStrategy 创建策略
func (c *StrategyController) CreateStrategy(name string, description string, strategyType string, parameters map[string]interface{}) error {
	_, err := c.strategyService.CreateStrategy(name, description, strategyType, parameters)
	return err
}

// UpdateStrategy 更新策略
func (c *StrategyController) UpdateStrategy(id int64, name string, description string, strategyType string, parameters map[string]interface{}) error {
	return c.strategyService.UpdateStrategy(id, name, description, strategyType, parameters)
}

// DeleteStrategy 删除策略
func (c *StrategyController) DeleteStrategy(id int64) error {
	return c.strategyService.DeleteStrategy(id)
}

// GetStrategy 获取策略
func (c *StrategyController) GetStrategy(id int64) (interface{}, error) {
	strategy, err := c.strategyService.GetStrategy(id)
	if err != nil {
		return nil, err
	}
	if strategy == nil {
		return nil, nil
	}
	return strategy, nil
}

// GetAllStrategies 获取所有策略
func (c *StrategyController) GetAllStrategies() (interface{}, error) {
	return c.strategyService.GetAllStrategies()
}

// GetStrategyTypes 获取所有策略类型
func (c *StrategyController) GetStrategyTypes() interface{} {
	return c.strategyService.GetStrategyTypes()
}

// UpdateStrategyBacktestResult 更新策略回测结果
func (c *StrategyController) UpdateStrategyBacktestResult(id int64, backtestResult map[string]interface{}) error {
	return c.strategyService.UpdateStrategyBacktestResult(id, backtestResult)
}
