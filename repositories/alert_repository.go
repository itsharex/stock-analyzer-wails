package repositories

import (
	"fmt"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AlertRepository 定义预警持久化接口
type AlertRepository interface {
	SaveAlertHistory(alert *models.PriceAlert, message string) error
	SaveActiveAlerts(alerts []*models.PriceAlert) error
	LoadActiveAlerts() ([]*models.PriceAlert, error)
	GetAlertHistory(stockCode string, limit int) ([]map[string]interface{}, error)
}

// SQLiteAlertRepository 基于 SQLite 的实现
type SQLiteAlertRepository struct {
	db *gorm.DB
}

// NewSQLiteAlertRepository 构造函数
func NewSQLiteAlertRepository(db *gorm.DB) *SQLiteAlertRepository {
	return &SQLiteAlertRepository{db: db}
}

// SaveAlertHistory 保存告警记录到 alert_history 表
func (r *SQLiteAlertRepository) SaveAlertHistory(alert *models.PriceAlert, message string) error {
	entity := models.AlertHistoryEntity{
		StockCode:      alert.StockCode,
		StockName:      alert.StockName,
		TriggeredPrice: alert.Price,
		Message:        message,
		TriggeredAt:    time.Now(),
	}

	if err := r.db.Create(&entity).Error; err != nil {
		logger.Error("保存告警历史失败", zap.Error(err))
		return fmt.Errorf("保存告警历史失败: %w", err)
	}
	return nil
}

// SaveActiveAlerts 保存当前活跃的预警订阅到 alerts 表
func (r *SQLiteAlertRepository) SaveActiveAlerts(alerts []*models.PriceAlert) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 清空 alerts 表
		// 使用 Where("1=1") 避免 GORM 的全局删除保护（如果有）
		if err := tx.Exec("DELETE FROM alerts").Error; err != nil {
			return fmt.Errorf("清空 alerts 表失败: %w", err)
		}

		if len(alerts) == 0 {
			return nil
		}

		// 2. 批量插入新的活跃预警
		var entities []models.AlertEntity
		for _, alert := range alerts {
			entities = append(entities, models.AlertEntity{
				StockCode:     alert.StockCode,
				StockName:     alert.StockName,
				Price:         alert.Price,
				Type:          alert.Type,
				IsActive:      alert.IsActive,
				LastTriggered: alert.LastTriggered,
				CreatedAt:     time.Now(),
			})
		}

		if err := tx.CreateInBatches(entities, 100).Error; err != nil {
			return fmt.Errorf("批量插入预警失败: %w", err)
		}

		return nil
	})
}

// LoadActiveAlerts 加载保存的活跃预警订阅
func (r *SQLiteAlertRepository) LoadActiveAlerts() ([]*models.PriceAlert, error) {
	var entities []models.AlertEntity
	if err := r.db.Where("is_active = ?", true).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询活跃预警失败: %w", err)
	}

	var alerts []*models.PriceAlert
	for _, entity := range entities {
		alert := &models.PriceAlert{
			StockCode:     entity.StockCode,
			StockName:     entity.StockName,
			Price:         entity.Price,
			Type:          entity.Type,
			IsActive:      entity.IsActive,
			LastTriggered: entity.LastTriggered,
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// GetAlertHistory 获取告警历史，支持分页和股票代码筛选
func (r *SQLiteAlertRepository) GetAlertHistory(stockCode string, limit int) ([]map[string]interface{}, error) {
	var entities []models.AlertHistoryEntity
	query := r.db.Order("triggered_at DESC").Limit(limit)

	if stockCode != "" {
		query = query.Where("stock_code = ?", stockCode)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询告警历史失败: %w", err)
	}

	var history []map[string]interface{}
	for _, entity := range entities {
		history = append(history, map[string]interface{}{
			"stockCode":      entity.StockCode,
			"stockName":      entity.StockName,
			"triggeredPrice": entity.TriggeredPrice,
			"message":        entity.Message,
			"triggeredAt":    entity.TriggeredAt.Format(time.RFC3339),
		})
	}

	return history, nil
}
