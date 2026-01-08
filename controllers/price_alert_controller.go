package controllers

import (
	"encoding/json"
	"fmt"
	"stock-analyzer-wails/services"
	"time"

	"stock-analyzer-wails/internal/logger"

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

func newTraceID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func formatRFC3339(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	// 使用带时区的 RFC3339，避免前端 new Date("YYYY-MM-DD HH:mm:ss") 解析不一致导致时区偏移
	return t.Format(time.RFC3339Nano)
}

// CreateAlert Wails绑定方法：创建价格预警
func (c *PriceAlertController) CreateAlert(jsonData string) *CreateAlertResponse {
	traceId := newTraceID()
	logger.Info("PriceAlertController.CreateAlert 调用",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "CreateAlert"),
		zap.String("traceId", traceId),
		zap.Int("payload_len", len(jsonData)),
	)

	var req services.CreateAlertRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		logger.Warn("CreateAlert JSON 解析失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "CreateAlert"),
			zap.String("traceId", traceId),
			zap.Error(err),
		)
		return &CreateAlertResponse{
			Success: false,
			Message: "请求数据格式错误",
		}
	}

	if err := c.service.CreateAlert(&req); err != nil {
		logger.Warn("CreateAlert 服务层失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "CreateAlert"),
			zap.String("traceId", traceId),
			zap.String("stockCode", req.StockCode),
			zap.String("alertType", req.AlertType),
			zap.Error(err),
		)
		return &CreateAlertResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	logger.Info("CreateAlert 成功",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "CreateAlert"),
		zap.String("traceId", traceId),
		zap.String("stockCode", req.StockCode),
		zap.String("alertType", req.AlertType),
	)
	return &CreateAlertResponse{
		Success: true,
		Message: "预警创建成功",
	}
}

// UpdateAlert Wails绑定方法：更新价格预警
func (c *PriceAlertController) UpdateAlert(jsonData string) *CreateAlertResponse {
	traceId := newTraceID()
	logger.Info("PriceAlertController.UpdateAlert 调用",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "UpdateAlert"),
		zap.String("traceId", traceId),
		zap.Int("payload_len", len(jsonData)),
	)

	var req services.UpdateAlertRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		logger.Warn("UpdateAlert JSON 解析失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "UpdateAlert"),
			zap.String("traceId", traceId),
			zap.Error(err),
		)
		return &CreateAlertResponse{
			Success: false,
			Message: "请求数据格式错误",
		}
	}

	if err := c.service.UpdateAlert(&req); err != nil {
		logger.Warn("UpdateAlert 服务层失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "UpdateAlert"),
			zap.String("traceId", traceId),
			zap.Int64("id", req.ID),
			zap.String("stockCode", req.StockCode),
			zap.String("alertType", req.AlertType),
			zap.Error(err),
		)
		return &CreateAlertResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	logger.Info("UpdateAlert 成功",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "UpdateAlert"),
		zap.String("traceId", traceId),
		zap.Int64("id", req.ID),
		zap.String("stockCode", req.StockCode),
		zap.String("alertType", req.AlertType),
	)
	return &CreateAlertResponse{
		Success: true,
		Message: "预警更新成功",
	}
}

// DeleteAlert Wails绑定方法：删除价格预警
func (c *PriceAlertController) DeleteAlert(id int64) *CreateAlertResponse {
	traceId := newTraceID()
	logger.Info("PriceAlertController.DeleteAlert 调用",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "DeleteAlert"),
		zap.String("traceId", traceId),
		zap.Int64("id", id),
	)

	if err := c.service.DeleteAlert(id); err != nil {
		logger.Warn("DeleteAlert 服务层失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "DeleteAlert"),
			zap.String("traceId", traceId),
			zap.Int64("id", id),
			zap.Error(err),
		)
		return &CreateAlertResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	logger.Info("DeleteAlert 成功",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "DeleteAlert"),
		zap.String("traceId", traceId),
		zap.Int64("id", id),
	)
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
	traceId := newTraceID()
	logger.Info("PriceAlertController.ToggleAlertStatus 调用",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "ToggleAlertStatus"),
		zap.String("traceId", traceId),
		zap.Int64("id", id),
		zap.Bool("isActive", isActive),
	)

	if err := c.service.ToggleAlertStatus(id, isActive); err != nil {
		logger.Warn("ToggleAlertStatus 服务层失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "ToggleAlertStatus"),
			zap.String("traceId", traceId),
			zap.Int64("id", id),
			zap.Bool("isActive", isActive),
			zap.Error(err),
		)
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
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Alerts  []map[string]interface{} `json:"alerts"`
}

// GetAllAlerts Wails绑定方法：获取所有预警
func (c *PriceAlertController) GetAllAlerts() *GetAlertsResponse {
	traceId := newTraceID()
	logger.Info("PriceAlertController.GetAllAlerts 调用",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "GetAllAlerts"),
		zap.String("traceId", traceId),
	)

	alerts, err := c.service.GetAllAlerts()
	if err != nil {
		logger.Warn("GetAllAlerts 服务层失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "GetAllAlerts"),
			zap.String("traceId", traceId),
			zap.Error(err),
		)
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
			"id":                alert.ID,
			"stockCode":         alert.StockCode,
			"stockName":         alert.StockName,
			"alertType":         alert.AlertType,
			"conditions":        alert.Conditions,
			"isActive":          alert.IsActive,
			"sensitivity":       alert.Sensitivity,
			"cooldownHours":     alert.CooldownHours,
			"postTriggerAction": alert.PostTriggerAction,
			"enableSound":       alert.EnableSound,
			"enableDesktop":     alert.EnableDesktop,
			"templateId":        alert.TemplateID,
			"createdAt":         formatRFC3339(alert.CreatedAt),
			"updatedAt":         formatRFC3339(alert.UpdatedAt),
			"lastTriggeredAt":   formatRFC3339(alert.LastTriggeredAt),
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
	traceId := newTraceID()
	logger.Info("PriceAlertController.GetActiveAlerts 调用",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "GetActiveAlerts"),
		zap.String("traceId", traceId),
	)

	alerts, err := c.service.GetActiveAlerts()
	if err != nil {
		logger.Warn("GetActiveAlerts 服务层失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "GetActiveAlerts"),
			zap.String("traceId", traceId),
			zap.Error(err),
		)
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
			"id":                alert.ID,
			"stockCode":         alert.StockCode,
			"stockName":         alert.StockName,
			"alertType":         alert.AlertType,
			"conditions":        alert.Conditions,
			"isActive":          alert.IsActive,
			"sensitivity":       alert.Sensitivity,
			"cooldownHours":     alert.CooldownHours,
			"postTriggerAction": alert.PostTriggerAction,
			"enableSound":       alert.EnableSound,
			"enableDesktop":     alert.EnableDesktop,
			"templateId":        alert.TemplateID,
			"createdAt":         formatRFC3339(alert.CreatedAt),
			"updatedAt":         formatRFC3339(alert.UpdatedAt),
			"lastTriggeredAt":   formatRFC3339(alert.LastTriggeredAt),
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
	traceId := newTraceID()
	logger.Info("PriceAlertController.GetAlertsByStockCode 调用",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "GetAlertsByStockCode"),
		zap.String("traceId", traceId),
		zap.String("stockCode", stockCode),
	)

	alerts, err := c.service.GetAlertsByStockCode(stockCode)
	if err != nil {
		logger.Warn("GetAlertsByStockCode 服务层失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "GetAlertsByStockCode"),
			zap.String("traceId", traceId),
			zap.String("stockCode", stockCode),
			zap.Error(err),
		)
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
			"id":                alert.ID,
			"stockCode":         alert.StockCode,
			"stockName":         alert.StockName,
			"alertType":         alert.AlertType,
			"conditions":        alert.Conditions,
			"isActive":          alert.IsActive,
			"sensitivity":       alert.Sensitivity,
			"cooldownHours":     alert.CooldownHours,
			"postTriggerAction": alert.PostTriggerAction,
			"enableSound":       alert.EnableSound,
			"enableDesktop":     alert.EnableDesktop,
			"templateId":        alert.TemplateID,
			"createdAt":         formatRFC3339(alert.CreatedAt),
			"updatedAt":         formatRFC3339(alert.UpdatedAt),
			"lastTriggeredAt":   formatRFC3339(alert.LastTriggeredAt),
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
	Success   bool                     `json:"success"`
	Message   string                   `json:"message"`
	Histories []map[string]interface{} `json:"histories"`
}

// GetTriggerHistory Wails绑定方法：获取预警触发历史
func (c *PriceAlertController) GetTriggerHistory(stockCode string, limit int) *GetTriggerHistoryResponse {
	traceId := newTraceID()
	logger.Info("PriceAlertController.GetTriggerHistory 调用",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "GetTriggerHistory"),
		zap.String("traceId", traceId),
		zap.String("stockCode", stockCode),
		zap.Int("limit", limit),
	)

	histories, err := c.service.GetTriggerHistory(stockCode, limit)
	if err != nil {
		logger.Warn("GetTriggerHistory 服务层失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "GetTriggerHistory"),
			zap.String("traceId", traceId),
			zap.String("stockCode", stockCode),
			zap.Int("limit", limit),
			zap.Error(err),
		)
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
			"id":             history.ID,
			"alertId":        history.AlertID,
			"stockCode":      history.StockCode,
			"stockName":      history.StockName,
			"alertType":      history.AlertType,
			"triggerPrice":   history.TriggerPrice,
			"triggerMessage": history.TriggerMessage,
			"triggeredAt":    formatRFC3339(history.TriggeredAt),
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
	Success   bool                     `json:"success"`
	Message   string                   `json:"message"`
	Templates []map[string]interface{} `json:"templates"`
}

// GetAllTemplates Wails绑定方法：获取所有预警模板
func (c *PriceAlertController) GetAllTemplates() *GetTemplatesResponse {
	traceId := newTraceID()
	logger.Info("PriceAlertController.GetAllTemplates 调用",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "GetAllTemplates"),
		zap.String("traceId", traceId),
	)

	templates, err := c.service.GetAllTemplates()
	if err != nil {
		logger.Warn("GetAllTemplates 服务层失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "GetAllTemplates"),
			zap.String("traceId", traceId),
			zap.Error(err),
		)
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
			"createdAt":   formatRFC3339(template.CreatedAt),
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
	traceId := newTraceID()
	logger.Info("PriceAlertController.CreateAlertFromTemplate 调用",
		zap.String("module", "controllers.price_alert"),
		zap.String("op", "CreateAlertFromTemplate"),
		zap.String("traceId", traceId),
		zap.String("templateID", templateID),
		zap.String("stockCode", stockCode),
		zap.String("stockName", stockName),
		zap.Int("params_len", len(paramsJSON)),
	)

	var params map[string]interface{}
	if paramsJSON != "" {
		if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
			logger.Warn("CreateAlertFromTemplate 参数解析失败",
				zap.String("module", "controllers.price_alert"),
				zap.String("op", "CreateAlertFromTemplate"),
				zap.String("traceId", traceId),
				zap.Error(err),
			)
			return &CreateAlertFromTemplateResponse{
				Success: false,
				Message: "参数格式错误",
			}
		}
	}

	if err := c.service.CreateAlertFromTemplate(templateID, stockCode, stockName, params); err != nil {
		logger.Warn("CreateAlertFromTemplate 服务层失败",
			zap.String("module", "controllers.price_alert"),
			zap.String("op", "CreateAlertFromTemplate"),
			zap.String("traceId", traceId),
			zap.String("templateID", templateID),
			zap.String("stockCode", stockCode),
			zap.Error(err),
		)
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
