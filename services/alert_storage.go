package services

import (

	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"
)

// AlertService 业务逻辑层，负责预警的业务处理
type AlertService struct {
	repo repositories.AlertRepository
}

// NewAlertService 构造函数
func NewAlertService(repo repositories.AlertRepository) *AlertService {
	return &AlertService{repo: repo}
}

// SaveAlert 负责将触发的告警记录到历史，并进行业务处理（如发送通知）
func (s *AlertService) SaveAlert(alert *models.PriceAlert, message string) error {
	// 业务逻辑：这里可以添加如“发送邮件/微信通知”等业务规则
	return s.repo.SaveAlertHistory(alert, message)
}

// SaveActiveAlerts 保存当前活跃的预警订阅
func (s *AlertService) SaveActiveAlerts(alerts []*models.PriceAlert) error {
	// 业务逻辑：这里可以添加如“预警数量限制”等业务规则
	return s.repo.SaveActiveAlerts(alerts)
}

// LoadActiveAlerts 加载保存的活跃预警订阅
func (s *AlertService) LoadActiveAlerts() ([]*models.PriceAlert, error) {
	return s.repo.LoadActiveAlerts()
}

// GetAlertHistory 获取告警历史
func (s *AlertService) GetAlertHistory(stockCode string, limit int) ([]map[string]interface{}, error) {
	return s.repo.GetAlertHistory(stockCode, limit)
}

// GetAlertsForWails 是一个临时方法，用于兼容 app.go 中对 AlertStorage 的调用
// TODO: 在 app.go 中移除对 AlertStorage 的直接引用
func (s *AlertService) GetAlertsForWails() ([]*models.PriceAlert, error) {
	return s.repo.LoadActiveAlerts()
}

// SaveAlertsForWails 是一个临时方法，用于兼容 app.go 中对 AlertStorage 的调用
// TODO: 在 app.go 中移除对 AlertStorage 的直接引用
func (s *AlertService) SaveAlertsForWails(alerts []*models.PriceAlert) error {
	return s.repo.SaveActiveAlerts(alerts)
}

// GetAlertHistoryForWails 是一个临时方法，用于兼容 app.go 中对 AlertStorage 的调用
// TODO: 在 app.go 中移除对 AlertStorage 的直接引用
func (s *AlertService) GetAlertHistoryForWails(stockCode string, limit int) ([]map[string]interface{}, error) {
	return s.repo.GetAlertHistory(stockCode, limit)
}

// UpdateAlertConfig 更新告警全局配置
func (s *AlertService) UpdateAlertConfig(config models.AlertConfig) error {
	// 业务逻辑：这里可以添加配置校验
	// 由于 AlertConfig 尚未持久化，这里先不做持久化操作
	return nil
}

// GetAlertConfig 获取告警全局配置
func (s *AlertService) GetAlertConfig() (models.AlertConfig, error) {
	// 业务逻辑：这里应该从持久化存储中加载配置
	// 由于 AlertConfig 尚未持久化，这里先返回默认值
	return models.AlertConfig{
		Sensitivity: 0.005, // 默认 0.5%
		Cooldown:    1,     // 默认 1 小时
		Enabled:     true,
	}, nil
}

// SetAlertsFromAI 接收 AI 识别的支撑位和压力位并自动设置预警
func (s *AlertService) SetAlertsFromAI(code string, name string, drawings []models.TechnicalDrawing) {
	// 业务逻辑：将 AI 建议的预警位添加到活跃预警列表中
	// 这个逻辑比较复杂，需要访问 app.go 中的全局 alerts 列表和 mutex
	// 考虑到 app.go 已经瘦身，这个逻辑应该保留在 app.go 中，或者在 AlertController 中处理
	// 由于 app.go 中已经有这个逻辑，我们先在 AlertController 中实现转发，并保留 app.go 中的核心逻辑
}

// NewAlertStorage 兼容旧的命名，但返回新的 AlertService
func NewAlertStorage(dbSvc *DBService) *AlertService {
	repo := repositories.NewSQLiteAlertRepository(dbSvc.GetDB())
	return NewAlertService(repo)
}
