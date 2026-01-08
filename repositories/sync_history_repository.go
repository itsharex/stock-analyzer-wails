package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"

	"go.uber.org/zap"
)

// SyncHistoryRepository 定义同步历史持久化接口
type SyncHistoryRepository interface {
	Add(history *models.SyncHistory) error
	GetAll(limit int, offset int) ([]*models.SyncHistory, error)
	GetByStockCode(code string, limit int) ([]*models.SyncHistory, error)
	GetCount() (int, error)
	ClearAll() error
}

// SQLiteSyncHistoryRepository 基于 SQLite 的实现
type SQLiteSyncHistoryRepository struct {
	db *sql.DB
}

// NewSQLiteSyncHistoryRepository 构造函数
func NewSQLiteSyncHistoryRepository(db *sql.DB) *SQLiteSyncHistoryRepository {
	return &SQLiteSyncHistoryRepository{db: db}
}

// Add 添加同步历史记录
func (r *SQLiteSyncHistoryRepository) Add(history *models.SyncHistory) error {
	start := time.Now()

	query := `
		INSERT INTO sync_history (stock_code, stock_name, sync_type, start_date, end_date, 
			status, records_added, records_updated, duration, error_msg, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		history.StockCode,
		history.StockName,
		history.SyncType,
		history.StartDate,
		history.EndDate,
		history.Status,
		history.RecordsAdded,
		history.RecordsUpdated,
		history.Duration,
		history.ErrorMsg,
		history.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("添加同步历史记录失败: %w", err)
	}

	logger.Info("成功添加同步历史记录",
		zap.String("module", "repositories.sync_history"),
		zap.String("op", "add_sync_history"),
		zap.String("stock_code", history.StockCode),
		zap.String("status", history.Status),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return nil
}

// GetAll 获取所有同步历史记录（分页）
func (r *SQLiteSyncHistoryRepository) GetAll(limit int, offset int) ([]*models.SyncHistory, error) {
	start := time.Now()

	query := `
		SELECT id, stock_code, stock_name, sync_type, start_date, end_date, 
			status, records_added, records_updated, duration, error_msg, created_at
		FROM sync_history
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("查询同步历史记录失败: %w", err)
	}
	defer rows.Close()

	histories := make([]*models.SyncHistory, 0)
	for rows.Next() {
		var history models.SyncHistory
		err := rows.Scan(
			&history.ID,
			&history.StockCode,
			&history.StockName,
			&history.SyncType,
			&history.StartDate,
			&history.EndDate,
			&history.Status,
			&history.RecordsAdded,
			&history.RecordsUpdated,
			&history.Duration,
			&history.ErrorMsg,
			&history.CreatedAt,
		)
		if err != nil {
			logger.Error("扫描同步历史数据失败", zap.Error(err))
			continue
		}
		histories = append(histories, &history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历同步历史结果集失败: %w", err)
	}

	logger.Debug("成功获取同步历史记录",
		zap.String("module", "repositories.sync_history"),
		zap.String("op", "get_all_sync_history"),
		zap.Int("histories_count", len(histories)),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return histories, nil
}

// GetByStockCode 根据股票代码获取同步历史记录
func (r *SQLiteSyncHistoryRepository) GetByStockCode(code string, limit int) ([]*models.SyncHistory, error) {
	start := time.Now()

	query := `
		SELECT id, stock_code, stock_name, sync_type, start_date, end_date, 
			status, records_added, records_updated, duration, error_msg, created_at
		FROM sync_history
		WHERE stock_code = ?
		ORDER BY created_at DESC
		LIMIT ?
	`
	rows, err := r.db.Query(query, code, limit)
	if err != nil {
		return nil, fmt.Errorf("查询股票同步历史记录失败: %w", err)
	}
	defer rows.Close()

	histories := make([]*models.SyncHistory, 0)
	for rows.Next() {
		var history models.SyncHistory
		err := rows.Scan(
			&history.ID,
			&history.StockCode,
			&history.StockName,
			&history.SyncType,
			&history.StartDate,
			&history.EndDate,
			&history.Status,
			&history.RecordsAdded,
			&history.RecordsUpdated,
			&history.Duration,
			&history.ErrorMsg,
			&history.CreatedAt,
		)
		if err != nil {
			logger.Error("扫描同步历史数据失败", zap.Error(err))
			continue
		}
		histories = append(histories, &history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历同步历史结果集失败: %w", err)
	}

	logger.Debug("成功获取股票同步历史记录",
		zap.String("module", "repositories.sync_history"),
		zap.String("op", "get_sync_history_by_code"),
		zap.String("stock_code", code),
		zap.Int("histories_count", len(histories)),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return histories, nil
}

// GetCount 获取同步历史记录总数
func (r *SQLiteSyncHistoryRepository) GetCount() (int, error) {
	start := time.Now()

	var count int
	query := `SELECT COUNT(*) FROM sync_history`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询同步历史记录总数失败: %w", err)
	}

	logger.Debug("成功获取同步历史记录总数",
		zap.String("module", "repositories.sync_history"),
		zap.String("op", "get_sync_history_count"),
		zap.Int("count", count),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return count, nil
}

// ClearAll 清除所有同步历史记录
func (r *SQLiteSyncHistoryRepository) ClearAll() error {
	start := time.Now()

	query := `DELETE FROM sync_history`
	result, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("清除同步历史记录失败: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()

	logger.Info("成功清除同步历史记录",
		zap.String("module", "repositories.sync_history"),
		zap.String("op", "clear_all_sync_history"),
		zap.Int64("rows_affected", rowsAffected),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return nil
}
