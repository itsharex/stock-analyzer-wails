package controllers

import (
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/services"
)

// AlertController 负责处理前端对预警的操作请求
type AlertController struct {
	service *services.AlertService
}

// NewAlertController 构造函数
func NewAlertController(svc *services.AlertService) *AlertController {
	return &AlertController{service: svc}
}

// GetAlerts Wails 绑定方法：获取当前活跃的预警订阅
func (c *AlertController) GetAlerts() ([]*models.PriceAlert, error) {
	return c.service.GetAlertsForWails()
}

// SaveAlerts Wails 绑定方法：保存当前活跃的预警订阅
func (c *AlertController) SaveAlerts(alerts []*models.PriceAlert) error {
	return c.service.SaveAlertsForWails(alerts)
}

// GetAlertHistory Wails 绑定方法：获取告警历史
func (c *AlertController) GetAlertHistory(code string, limit int) ([]map[string]interface{}, error) {
	return c.service.GetAlertHistoryForWails(code, limit)
}
