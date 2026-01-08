package controllers

import (
	"encoding/json"
	"stock-analyzer-wails/services"

	"go.uber.org/zap"
)

// PriceAlertController 价格预警控制器
type PriceAlertController struct {
	service *services.PriceAlertService
}

// NewPriceAlertController 创建价格预警控制器
func NewPriceAlertController(svc *services.PriceAlertService) *PriceAlertController {
	return &PriceAlertController{service: svc}
}

// CreateAlertResponse 创建预警响应
type CreateAlertResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	AlertID int64  `json:"alertId,omitempty"`
}

// CreateAlert Wails绑定方法：创建价格预警
func (c *PriceAlertController) CreateAlert(jsonData string) *CreateAlertResponse {
	var req services.CreateAlertRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		return &CreateAlertResponse{
			Success: false,
			Message: "请求数据格式错误",
		}
	}

	if err := c.service.CreateAlert(&req); err != nil {
		return &CreateAlertResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	return &CreateAlertResponse{
		Success: true,
		Message: "预警创建成功",
	}
}

// UpdateAlert Wails绑定方法：更新价格预警
func (c *PriceAlertController) UpdateAlert(jsonData string) *CreateAlertResponse {
	var req services.UpdateAlertRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		return &CreateAlertResponse{
			Success: false,
			Message: "请求数据格式错误",
		}
	}

	if err := c.service.UpdateAlert(&req); err != nil {
		return &CreateAlertResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	return &CreateAlertResponse{
		Success: true,
		Message: "预警更新成功",
	}
}

// DeleteAlert Wails绑定方法：删除价格预警
func (c *PriceAlertController) DeleteAlert(id int64) *CreateAlertResponse {
	if err := c.service.DeleteAlert(id); err != nil {
		return &CreateAlertResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	return &CreateAlertResponse{
		Success: true,
		Message: "预警删除成功",
	}
}

// ToggleAlertStatusResponse 切换预警状态响应
type ToggleAlertStatusResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ToggleAlertStatus Wails绑定方法：切换预警状态
func (c *PriceAlertController) ToggleAlertStatus(id int64, isActive bool) *ToggleAlertStatusResponse {
	if err := c.service.ToggleAlertStatus(id, isActive); err != nil {
		return &ToggleAlertStatusResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	status := "启用"
	if !isActive {
		status = "禁用"
	}

	return &ToggleAlertStatusResponse{
		Success: true,
		Message: "预警" + status + "成功",
	}
}

// GetAlertsResponse 获取预警列表响应
type GetAlertsResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Alerts  []map[string]interface{} `json:"alerts"`
}

// GetAllAlerts Wails绑定方法：获取所有预警
func (c *PriceAlertController) GetAllAlerts() *GetAlertsResponse {
	alerts, err := c.service.GetAllAlerts()
	if err != nil {
		return &GetAlertsResponse{
			Success: false,
			Message: err.Error(),
			Alerts:  nil,
		}
	}

	// 转换为map格式
	alertsMap := make([]map[string]interface{}, len(alerts))
	for i, alert := range alerts {
		alertsMap[i] = map[string]interface{}{
			"id":                  alert.ID,
			"stockCode":           alert.StockCode,
			"stockName":           alert.StockName,
			"alertType":           alert.AlertType,
			"conditions":          alert.Conditions,
			"isActive":            alert.IsActive,
			"sensitivity":         alert.Sensitivity,
			"cooldownHours":       alert.CooldownHours,
			"postTriggerAction":   alert.PostTriggerAction,
			"enableSound":         alert.EnableSound,
			"enableDesktop":       alert.EnableDesktop,
			"templateId":          alert.TemplateID,
			"createdAt":           alert.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt":           alert.UpdatedAt.Format("2006-01-02 15:04:05"),
			"lastTriggeredAt":     alert.LastTriggeredAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &GetAlertsResponse{
		Success: true,
		Message: "",
		Alerts:  alertsMap,
	}
}

// GetActiveAlerts Wails绑定方法：获取活跃的预警
func (c *PriceAlertController) GetActiveAlerts() *GetAlertsResponse {
	alerts, err := c.service.GetActiveAlerts()
	if err != nil {
		return &GetAlertsResponse{
			Success: false,
			Message: err.Error(),
			Alerts:  nil,
		}
	}

	// 转换为map格式
	alertsMap := make([]map[string]interface{}, len(alerts))
	for i, alert := range alerts {
		alertsMap[i] = map[string]interface{}{
			"id":                  alert.ID,
			"stockCode":           alert.StockCode,
			"stockName":           alert.StockName,
			"alertType":           alert.AlertType,
			"conditions":          alert.Conditions,
			"isActive":            alert.IsActive,
			"sensitivity":         alert.Sensitivity,
			"cooldownHours":       alert.CooldownHours,
			"postTriggerAction":   alert.PostTriggerAction,
			"enableSound":         alert.EnableSound,
			"enableDesktop":       alert.EnableDesktop,
			"templateId":          alert.TemplateID,
			"createdAt":           alert.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt":           alert.UpdatedAt.Format("2006-01-02 15:04:05"),
			"lastTriggeredAt":     alert.LastTriggeredAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &GetAlertsResponse{
		Success: true,
		Message: "",
		Alerts:  alertsMap,
	}
}

// GetAlertsByStockCode Wails绑定方法：根据股票代码获取预警
func (c *PriceAlertController) GetAlertsByStockCode(stockCode string) *GetAlertsResponse {
	alerts, err := c.service.GetAlertsByStockCode(stockCode)
	if err != nil {
		return &GetAlertsResponse{
			Success: false,
			Message: err.Error(),
			Alerts:  nil,
		}
	}

	// 转换为map格式
	alertsMap := make([]map[string]interface{}, len(alerts))
	for i, alert := range alerts {
		alertsMap[i] = map[string]interface{}{
			"id":                  alert.ID,
			"stockCode":           alert.StockCode,
			"stockName":           alert.StockName,
			"alertType":           alert.AlertType,
			"conditions":          alert.Conditions,
			"isActive":            alert.IsActive,
			"sensitivity":         alert.Sensitivity,
			"cooldownHours":       alert.CooldownHours,
			"postTriggerAction":   alert.PostTriggerAction,
			"enableSound":         alert.EnableSound,
			"enableDesktop":       alert.EnableDesktop,
			"templateId":          alert.TemplateID,
			"createdAt":           alert.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt":           alert.UpdatedAt.Format("2006-01-02 15:04:05"),
			"lastTriggeredAt":     alert.LastTriggeredAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &GetAlertsResponse{
		Success: true,
		Message: "",
		Alerts:  alertsMap,
	}
}

// GetTriggerHistoryResponse 获取预警触发历史响应
type GetTriggerHistoryResponse struct {
	Success   bool                    `json:"success"`
	Message   string                  `json:"message"`
	Histories []map[string]interface{} `json:"histories"`
}

// GetTriggerHistory Wails绑定方法：获取预警触发历史
func (c *PriceAlertController) GetTriggerHistory(stockCode string, limit int) *GetTriggerHistoryResponse {
	histories, err := c.service.GetTriggerHistory(stockCode, limit)
	if err != nil {
		return &GetTriggerHistoryResponse{
			Success:   false,
			Message:   err.Error(),
			Histories: nil,
		}
	}

	// 转换为map格式
	historiesMap := make([]map[string]interface{}, len(histories))
	for i, history := range histories {
		historiesMap[i] = map[string]interface{}{
			"id":              history.ID,
			"alertId":         history.AlertID,
			"stockCode":       history.StockCode,
			"stockName":       history.StockName,
			"alertType":       history.AlertType,
			"triggerPrice":    history.TriggerPrice,
			"triggerMessage":  history.TriggerMessage,
			"triggeredAt":     history.TriggeredAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &GetTriggerHistoryResponse{
		Success:   true,
		Message:   "",
		Histories: historiesMap,
	}
}

// GetTemplatesResponse 获取预警模板响应
type GetTemplatesResponse struct {
	Success   bool                    `json:"success"`
	Message   string                  `json:"message"`
	Templates []map[string]interface{} `json:"templates"`
}

// GetAllTemplates Wails绑定方法：获取所有预警模板
func (c *PriceAlertController) GetAllTemplates() *GetTemplatesResponse {
	templates, err := c.service.GetAllTemplates()
	if err != nil {
		return &GetTemplatesResponse{
			Success:   false,
			Message:   err.Error(),
			Templates: nil,
		}
	}

	// 转换为map格式
	templatesMap := make([]map[string]interface{}, len(templates))
	for i, template := range templates {
		templatesMap[i] = map[string]interface{}{
			"id":          template.ID,
			"name":        template.Name,
			"description": template.Description,
			"alertType":   template.AlertType,
			"conditions":  template.Conditions,
			"createdAt":   template.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &GetTemplatesResponse{
		Success:   true,
		Message:   "",
		Templates: templatesMap,
	}
}

// CreateAlertFromTemplateResponse 从模板创建预警响应
type CreateAlertFromTemplateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	AlertID int64  `json:"alertId,omitempty"`
}

// CreateAlertFromTemplate Wails绑定方法：从模板创建预警
func (c *PriceAlertController) CreateAlertFromTemplate(templateID, stockCode, stockName string, paramsJSON string) *CreateAlertFromTemplateResponse {
	var params map[string]interface{}
	if paramsJSON != "" {
		if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
			return &CreateAlertFromTemplateResponse{
				Success: false,
				Message: "参数格式错误",
			}
		}
	}

	if err := c.service.CreateAlertFromTemplate(templateID, stockCode, stockName, params); err != nil {
		return &CreateAlertFromTemplateResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	return &CreateAlertFromTemplateResponse{
		Success: true,
		Message: "从模板创建预警成功",
	}
}
