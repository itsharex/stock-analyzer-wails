package controllers

import (

	"stock-analyzer-wails/services"
)

// ConfigController 负责处理前端对配置的操作请求
type ConfigController struct {
	service *services.ConfigService
}

// NewConfigController 构造函数
func NewConfigController(svc *services.ConfigService) *ConfigController {
	return &ConfigController{service: svc}
}

// GetAIConfig Wails 绑定方法：获取 AI 配置
func (c *ConfigController) GetAIConfig() (services.AIResolvedConfig, error) {
	return c.service.LoadAIConfig()
}

// SaveAIConfig Wails 绑定方法：保存 AI 配置
func (c *ConfigController) SaveAIConfig(config services.AIResolvedConfig) error {
	return c.service.SaveAIConfig(config)
}

// GetGlobalStrategyConfig Wails 绑定方法：获取全局策略配置
func (c *ConfigController) GetGlobalStrategyConfig() (services.GlobalStrategyConfig, error) {
	return c.service.GetGlobalStrategyConfig()
}

// UpdateGlobalStrategyConfig Wails 绑定方法：更新全局策略配置
func (c *ConfigController) UpdateGlobalStrategyConfig(config services.GlobalStrategyConfig) error {
	return c.service.UpdateGlobalStrategyConfig(config)
}
