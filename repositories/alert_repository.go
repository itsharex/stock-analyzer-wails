package repositories

import (
	"database/sql"
	"fmt"
	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"
	"time"

	"go.uber.org/zap"
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
	db *sql.DB
}

// NewSQLiteAlertRepository 构造函数
func NewSQLiteAlertRepository(db *sql.DB) *SQLiteAlertRepository {
	return &SQLiteAlertRepository{db: db}
}

// SaveAlertHistory 保存告警记录到 alert_history 表
func (r *SQLiteAlertRepository) SaveAlertHistory(alert *models.PriceAlert, message string) error {
	query := `
		INSERT INTO alert_history (stock_code, stock_name, triggered_price, message)
		VALUES (?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, alert.StockCode, alert.StockName, alert.Price, message)
	if err != nil {
		logger.Error("保存告警历史失败", zap.Error(err))
		return fmt.Errorf("保存告警历史失败: %w", err)
	}
	return nil
}

// SaveActiveAlerts 保存当前活跃的预警订阅到 alerts 表
func (r *SQLiteAlertRepository) SaveActiveAlerts(alerts []*models.PriceAlert) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 清空 alerts 表
	if _, err := tx.Exec("DELETE FROM alerts"); err != nil {
		return fmt.Errorf("清空 alerts 表失败: %w", err)
	}

	// 2. 批量插入新的活跃预警
	stmt, err := tx.Prepare(`
		INSERT INTO alerts (stock_code, stock_name, price, type, is_active, last_triggered)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("准备插入语句失败: %w", err)
	}
	defer stmt.Close()

	for _, alert := range alerts {
		_, err := stmt.Exec(
			alert.StockCode,
			alert.StockName,
			alert.Price,
			alert.Type,
			alert.IsActive,
			alert.LastTriggered,
		)
		if err != nil {
			return fmt.Errorf("插入预警失败 (%s): %w", alert.StockCode, err)
		}
	}

	return tx.Commit()
}

// LoadActiveAlerts 加载保存的活跃预警订阅
func (r *SQLiteAlertRepository) LoadActiveAlerts() ([]*models.PriceAlert, error) {
	rows, err := r.db.Query(`
		SELECT stock_code, stock_name, price, type, is_active, last_triggered
		FROM alerts
		WHERE is_active = TRUE
	`)
	if err != nil {
		return nil, fmt.Errorf("查询活跃预警失败: %w", err)
	}
	defer rows.Close()

	var alerts []*models.PriceAlert
	for rows.Next() {
		alert := &models.PriceAlert{}
		var lastTriggered time.Time
		
		err := rows.Scan(
			&alert.StockCode,
			&alert.StockName,
			&alert.Price,
			&alert.Type,
			&alert.IsActive,
			&lastTriggered,
		)
		if err != nil {
			logger.Error("扫描活跃预警数据失败", zap.Error(err))
			continue
		}
		alert.LastTriggered = lastTriggered
		alerts = append(alerts, alert)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历活跃预警结果集失败: %w", err)
	}

	return alerts, nil
}

// GetAlertHistory 获取告警历史，支持分页和股票代码筛选
func (r *SQLiteAlertRepository) GetAlertHistory(stockCode string, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT stock_code, stock_name, triggered_price, message, triggered_at
		FROM alert_history
	`
	args := []interface{}{}
	whereClauses := []string{}

	if stockCode != "" {
		whereClauses = append(whereClauses, "stock_code = ?")
		args = append(args, stockCode)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + whereClauses[0]
	}

	query += " ORDER BY triggered_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询告警历史失败: %w", err)
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var code, name, message string
		var price float64
		var triggeredAt time.Time

		if err := rows.Scan(&code, &name, &price, &message, &triggeredAt); err != nil {
			logger.Error("扫描告警历史数据失败", zap.Error(err))
			continue
		}

		history = append(history, map[string]interface{}{
			"stockCode": code,
			"stockName": name,
			"triggeredPrice": price,
			"message": message,
			"triggeredAt": triggeredAt.Format(time.RFC3339),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历告警历史结果集失败: %w", err)
	}

	return history, nil
}
