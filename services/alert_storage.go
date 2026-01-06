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

// NewAlertStorage 兼容旧的命名，但返回新的 AlertService
func NewAlertStorage(dbSvc *DBService) *AlertService {
	repo := repositories.NewSQLiteAlertRepository(dbSvc.GetDB())
	return NewAlertService(repo)
}
