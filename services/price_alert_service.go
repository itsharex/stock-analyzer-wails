package services

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/repositories"

	"go.uber.org/zap"
)

// PriceAlertService 价格预警服务
type PriceAlertService struct {
	repo *repositories.PriceAlertRepository
}

// NewPriceAlertService 创建价格预警服务
func NewPriceAlertService(repo *repositories.PriceAlertRepository) *PriceAlertService {
	return &PriceAlertService{repo: repo}
}

// GetRepository 获取 Repository（供 AlertMonitor 使用）
func (s *PriceAlertService) GetRepository() *repositories.PriceAlertRepository {
	return s.repo
}

// CreateAlertRequest 创建预警请求
type CreateAlertRequest struct {
	StockCode         string  `json:"stockCode"`
	StockName         string  `json:"stockName"`
	AlertType         string  `json:"alertType"`
	Conditions        string  `json:"conditions"`
	Sensitivity       float64 `json:"sensitivity"`
	CooldownHours     int     `json:"cooldownHours"`
	PostTriggerAction string  `json:"postTriggerAction"`
	EnableSound       bool    `json:"enableSound"`
	EnableDesktop     bool    `json:"enableDesktop"`
}

// UpdateAlertRequest 更新预警请求
type UpdateAlertRequest struct {
	ID                int64   `json:"id"`
	StockCode         string  `json:"stockCode"`
	StockName         string  `json:"stockName"`
	AlertType         string  `json:"alertType"`
	Conditions        string  `json:"conditions"`
	IsActive          bool    `json:"isActive"`
	Sensitivity       float64 `json:"sensitivity"`
	CooldownHours     int     `json:"cooldownHours"`
	PostTriggerAction string  `json:"postTriggerAction"`
	EnableSound       bool    `json:"enableSound"`
	EnableDesktop     bool    `json:"enableDesktop"`
}

// StockDataForAlert 用于预警检测的股票数据
type StockDataForAlert struct {
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	ClosePrice         float64 `json:"closePrice"`
	OpenPrice          float64 `json:"openPrice"`
	HighPrice          float64 `json:"highPrice"`
	LowPrice           float64 `json:"lowPrice"`
	PreClosePrice      float64 `json:"preClosePrice"`
	PriceChangePercent float64 `json:"priceChangePercent"` // 涨跌幅百分比
	Volume             int64   `json:"volume"`
	VolumeRatio        float64 `json:"volumeRatio"` // 量比
	MA5                float64 `json:"ma5"`         // 5日均线
	MA10               float64 `json:"ma10"`        // 10日均线
	MA20               float64 `json:"ma20"`        // 20日均线
	HistoricalHigh     float64 `json:"historicalHigh"` // 历史最高价
	HistoricalLow      float64 `json:"historicalLow"`  // 历史最低价
}

// CreateAlert 创建价格预警
func (s *PriceAlertService) CreateAlert(req *CreateAlertRequest) error {
	// 验证请求
	if err := s.validateAlertRequest(req); err != nil {
		return err
	}

	// 验证预警条件JSON格式
	if !s.isValidConditionsJSON(req.Conditions) {
		return fmt.Errorf("预警条件JSON格式无效")
	}

	alert := &repositories.PriceThresholdAlert{
		StockCode:         req.StockCode,
		StockName:         req.StockName,
		AlertType:         req.AlertType,
		Conditions:        req.Conditions,
		IsActive:          true,
		Sensitivity:       req.Sensitivity,
		CooldownHours:     req.CooldownHours,
		PostTriggerAction: req.PostTriggerAction,
		EnableSound:       req.EnableSound,
		EnableDesktop:     req.EnableDesktop,
	}

	return s.repo.CreateAlert(alert)
}

// UpdateAlert 更新价格预警
func (s *PriceAlertService) UpdateAlert(req *UpdateAlertRequest) error {
	// 验证请求
	if req.ID <= 0 {
		return fmt.Errorf("预警ID无效")
	}

	// 验证预警条件JSON格式
	if !s.isValidConditionsJSON(req.Conditions) {
		return fmt.Errorf("预警条件JSON格式无效")
	}

	alert, err := s.repo.GetAlertByID(req.ID)
	if err != nil {
		return fmt.Errorf("预警不存在")
	}

	alert.StockCode = req.StockCode
	alert.StockName = req.StockName
	alert.AlertType = req.AlertType
	alert.Conditions = req.Conditions
	alert.IsActive = req.IsActive
	alert.Sensitivity = req.Sensitivity
	alert.CooldownHours = req.CooldownHours
	alert.PostTriggerAction = req.PostTriggerAction
	alert.EnableSound = req.EnableSound
	alert.EnableDesktop = req.EnableDesktop

	return s.repo.UpdateAlert(alert)
}

// DeleteAlert 删除价格预警
func (s *PriceAlertService) DeleteAlert(id int64) error {
	if id <= 0 {
		return fmt.Errorf("预警ID无效")
	}
	return s.repo.DeleteAlert(id)
}

// ToggleAlertStatus 切换预警状态
func (s *PriceAlertService) ToggleAlertStatus(id int64, isActive bool) error {
	if id <= 0 {
		return fmt.Errorf("预警ID无效")
	}
	return s.repo.ToggleAlertStatus(id, isActive)
}

// GetAllAlerts 获取所有预警
func (s *PriceAlertService) GetAllAlerts() ([]*repositories.PriceThresholdAlert, error) {
	return s.repo.GetAllAlerts()
}

// GetActiveAlerts 获取活跃的预警
func (s *PriceAlertService) GetActiveAlerts() ([]*repositories.PriceThresholdAlert, error) {
	return s.repo.GetActiveAlerts()
}

// GetAlertsByStockCode 根据股票代码获取预警
func (s *PriceAlertService) GetAlertsByStockCode(stockCode string) ([]*repositories.PriceThresholdAlert, error) {
	return s.repo.GetAlertsByStockCode(stockCode)
}

// GetAllTemplates 获取所有预警模板
func (s *PriceAlertService) GetAllTemplates() ([]*repositories.PriceAlertTemplate, error) {
	return s.repo.GetAllTemplates()
}

// CreateAlertFromTemplate 从模板创建预警
func (s *PriceAlertService) CreateAlertFromTemplate(templateID string, stockCode, stockName string, params map[string]interface{}) error {
	template, err := s.repo.GetTemplateByID(templateID)
	if err != nil {
		return fmt.Errorf("模板不存在: %w", err)
	}

	// 解析模板条件
	var conditions repositories.PriceAlertConditions
	if err := json.Unmarshal([]byte(template.Conditions), &conditions); err != nil {
		return fmt.Errorf("解析模板条件失败: %w", err)
	}

	// 替换模板参数
	s.replaceTemplateParams(&conditions, params)

	// 序列化回JSON
	conditionsJSON, err := json.Marshal(conditions)
	if err != nil {
		return fmt.Errorf("序列化条件失败: %w", err)
	}

	req := &CreateAlertRequest{
		StockCode:         stockCode,
		StockName:         stockName,
		AlertType:         template.AlertType,
		Conditions:        string(conditionsJSON),
		Sensitivity:       0.001,
		CooldownHours:     1,
		PostTriggerAction: "continue",
		EnableSound:       true,
		EnableDesktop:     true,
	}

	return s.CreateAlert(req)
}

// replaceTemplateParams 替换模板中的参数
func (s *PriceAlertService) replaceTemplateParams(conditions *repositories.PriceAlertConditions, params map[string]interface{}) {
	for i := range conditions.Conditions {
		condition := &conditions.Conditions[i]

		// 如果参数中有对应的值，则替换
		if val, ok := params[condition.Field]; ok {
			switch v := val.(type) {
			case float64:
				condition.Value = v
			case int:
				condition.Value = float64(v)
			case int64:
				condition.Value = float64(v)
			}
		}

		// 替换特定的值参数
		if val, ok := params["value"]; ok {
			switch v := val.(type) {
			case float64:
				condition.Value = v
			case int:
				condition.Value = float64(v)
			case int64:
				condition.Value = float64(v)
			}
		}
	}
}

// GetTriggerHistory 获取预警触发历史
func (s *PriceAlertService) GetTriggerHistory(stockCode string, limit int) ([]*repositories.PriceAlertTriggerHistory, error) {
	return s.repo.GetTriggerHistory(stockCode, limit)
}

// CheckAlert 检测预警是否触发
func (s *PriceAlertService) CheckAlert(alert *repositories.PriceThresholdAlert, stockData *StockDataForAlert) (bool, string, error) {
	// 检查是否在冷却期内
	inCooldown, err := s.repo.IsInCooldown(alert.ID, alert.CooldownHours)
	if err != nil {
		return false, "", err
	}

	if inCooldown {
		return false, "在冷却期内", nil
	}

	// 解析预警条件
	var conditions repositories.PriceAlertConditions
	if err := json.Unmarshal([]byte(alert.Conditions), &conditions); err != nil {
		return false, "", fmt.Errorf("解析预警条件失败: %w", err)
	}

	// 检测条件是否满足
	triggered, message := s.evaluateConditions(&conditions, stockData, alert.Sensitivity)

	if triggered {
		return true, message, nil
	}

	return false, "条件不满足", nil
}

// evaluateConditions 评估预警条件
func (s *PriceAlertService) evaluateConditions(conditions *repositories.PriceAlertConditions, stockData *StockDataForAlert, sensitivity float64) (bool, string) {
	var results []bool
	var messages []string

	for _, condition := range conditions.Conditions {
		triggered, msg := s.evaluateCondition(&condition, stockData, sensitivity)
		results = append(results, triggered)
		if msg != "" {
			messages = append(messages, msg)
		}
	}

	// 根据逻辑关系判断
	if strings.ToUpper(conditions.Logic) == "OR" {
		// OR逻辑：任一条件满足即可
		for i, result := range results {
			if result {
				return true, s.buildTriggerMessage(conditions, messages, i)
			}
		}
		return false, ""
	} else {
		// AND逻辑（默认）：所有条件都必须满足
		for i, result := range results {
			if !result {
				return false, ""
			}
		}
		return true, s.buildTriggerMessage(conditions, messages, 0)
	}
}

// evaluateCondition 评估单个条件
func (s *PriceAlertService) evaluateCondition(condition *repositories.PriceAlertCondition, stockData *StockDataForAlert, sensitivity float64) (bool, string) {
	var actualValue float64
	var fieldName string

	// 根据字段名获取实际值
	switch condition.Field {
	case "price_change_percent":
		actualValue = stockData.PriceChangePercent
		fieldName = "涨跌幅"
	case "close_price":
		actualValue = stockData.ClosePrice
		fieldName = "收盘价"
	case "high_price":
		actualValue = stockData.HighPrice
		fieldName = "最高价"
	case "low_price":
		actualValue = stockData.LowPrice
		fieldName = "最低价"
	case "open_price":
		actualValue = stockData.OpenPrice
		fieldName = "开盘价"
	case "volume_ratio":
		actualValue = stockData.VolumeRatio
		fieldName = "量比"
	case "ma5":
		actualValue = stockData.MA5
		fieldName = "MA5"
	case "ma10":
		actualValue = stockData.MA10
		fieldName = "MA10"
	case "ma20":
		actualValue = stockData.MA20
		fieldName = "MA20"
	default:
		return false, fmt.Sprintf("未知字段: %s", condition.Field)
	}

	// 比较操作
	var triggered bool
	switch condition.Operator {
	case ">":
		triggered = actualValue > condition.Value
	case ">=":
		triggered = actualValue >= condition.Value
	case "<":
		triggered = actualValue < condition.Value
	case "<=":
		triggered = actualValue <= condition.Value
	case "==":
		// 考虑灵敏度容差
		triggered = math.Abs(actualValue-condition.Value) <= sensitivity
	case "!=":
		triggered = math.Abs(actualValue-condition.Value) > sensitivity
	default:
		return false, fmt.Sprintf("未知操作符: %s", condition.Operator)
	}

	message := fmt.Sprintf("%s %.2f %s %.2f", fieldName, actualValue, condition.Operator, condition.Value)

	// 特殊处理：历史高低点
	if condition.Reference == "historical_high" && condition.Field == "high_price" {
		if stockData.HistoricalHigh > 0 {
			triggered = stockData.HighPrice > stockData.HistoricalHigh*(1-sensitivity)
			message = fmt.Sprintf("最高价 %.2f 突破历史新高 %.2f", stockData.HighPrice, stockData.HistoricalHigh)
		}
	} else if condition.Reference == "historical_low" && condition.Field == "low_price" {
		if stockData.HistoricalLow > 0 {
			triggered = stockData.LowPrice < stockData.HistoricalLow*(1+sensitivity)
			message = fmt.Sprintf("最低价 %.2f 跌破历史新低 %.2f", stockData.LowPrice, stockData.HistoricalLow)
		}
	}

	// 均线金叉/死叉特殊处理
	if condition.Reference == "ma20" && condition.Field == "ma5" {
		triggered = stockData.MA5 > stockData.MA20*(1-sensitivity)
		message = fmt.Sprintf("MA5 %.2f 上穿 MA20 %.2f (金叉)", stockData.MA5, stockData.MA20)
	}

	return triggered, message
}

// buildTriggerMessage 构建触发消息
func (s *PriceAlertService) buildTriggerMessage(conditions *repositories.PriceAlertConditions, messages []string, index int) string {
	if len(messages) > 0 {
		return messages[index]
	}

	switch conditions.Conditions[0].Field {
	case "price_change_percent":
		if conditions.Conditions[0].Operator == ">" {
			return "涨幅达到预警阈值"
		} else if conditions.Conditions[0].Operator == "<" {
			return "跌幅达到预警阈值"
		}
	case "close_price":
		return "价格达到目标价"
	case "high_price":
		return "突破历史新高"
	case "low_price":
		return "跌破历史新低"
	case "volume_ratio":
		return "成交量激增"
	default:
		return "预警条件已触发"
	}
}

// TriggerAlert 触发预警
func (s *PriceAlertService) TriggerAlert(alert *repositories.PriceThresholdAlert, stockData *StockDataForAlert, message string) error {
	// 保存触发历史
	history := &repositories.PriceAlertTriggerHistory{
		AlertID:        alert.ID,
		StockCode:      alert.StockCode,
		StockName:      alert.StockName,
		AlertType:      alert.AlertType,
		TriggerPrice:   stockData.ClosePrice,
		TriggerMessage: message,
	}

	if err := s.repo.SaveTriggerHistory(history); err != nil {
		logger.Error("保存预警触发历史失败",
			zap.Int64("alert_id", alert.ID),
			zap.Error(err),
		)
		return err
	}

	// 更新最后触发时间
	if err := s.repo.UpdateLastTriggeredTime(alert.ID); err != nil {
		logger.Error("更新最后触发时间失败",
			zap.Int64("alert_id", alert.ID),
			zap.Error(err),
		)
		return err
	}

	// 根据触发后行为执行操作
	switch alert.PostTriggerAction {
	case "disable":
		// 禁用预警
		if err := s.repo.ToggleAlertStatus(alert.ID, false); err != nil {
			logger.Error("禁用预警失败",
				zap.Int64("alert_id", alert.ID),
				zap.Error(err),
			)
		}
		logger.Info("预警已禁用",
			zap.Int64("alert_id", alert.ID),
			zap.String("stock_code", alert.StockCode),
		)
	case "once":
		// 仅触发一次（自动禁用）
		if err := s.repo.ToggleAlertStatus(alert.ID, false); err != nil {
			logger.Error("禁用预警失败",
				zap.Int64("alert_id", alert.ID),
				zap.Error(err),
			)
		}
		logger.Info("预警已触发并禁用（仅触发一次）",
			zap.Int64("alert_id", alert.ID),
			zap.String("stock_code", alert.StockCode),
		)
	// "continue" 不做任何操作，预警继续监控
	}

	logger.Info("价格预警已触发",
		zap.Int64("alert_id", alert.ID),
		zap.String("stock_code", alert.StockCode),
		zap.String("stock_name", alert.StockName),
		zap.Float64("trigger_price", stockData.ClosePrice),
		zap.String("message", message),
	)

	return nil
}

// validateAlertRequest 验证预警请求
func (s *PriceAlertService) validateAlertRequest(req *CreateAlertRequest) error {
	if req.StockCode == "" {
		return fmt.Errorf("股票代码不能为空")
	}
	if req.StockName == "" {
		return fmt.Errorf("股票名称不能为空")
	}
	if req.AlertType == "" {
		return fmt.Errorf("预警类型不能为空")
	}
	if req.Conditions == "" {
		return fmt.Errorf("预警条件不能为空")
	}
	if req.Sensitivity < 0 || req.Sensitivity > 0.1 {
		return fmt.Errorf("灵敏度必须在0到0.1之间")
	}
	if req.CooldownHours < 0 || req.CooldownHours > 24 {
		return fmt.Errorf("冷却时间必须在0到24小时之间")
	}
	if req.PostTriggerAction != "continue" && req.PostTriggerAction != "disable" && req.PostTriggerAction != "once" {
		return fmt.Errorf("触发后行为无效")
	}
	return nil
}

// isValidConditionsJSON 验证预警条件JSON格式
func (s *PriceAlertService) isValidConditionsJSON(jsonStr string) bool {
	var conditions repositories.PriceAlertConditions
	return json.Unmarshal([]byte(jsonStr), &conditions) == nil
}
