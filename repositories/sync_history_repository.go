package repositories

import (
	"fmt"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
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
	db *gorm.DB
}

// NewSQLiteSyncHistoryRepository 构造函数
func NewSQLiteSyncHistoryRepository(db *gorm.DB) *SQLiteSyncHistoryRepository {
	return &SQLiteSyncHistoryRepository{db: db}
}

// Add 添加同步历史记录
func (r *SQLiteSyncHistoryRepository) Add(history *models.SyncHistory) error {
	start := time.Now()

	entity := models.SyncHistoryEntity{
		StockCode:      history.StockCode,
		StockName:      history.StockName,
		SyncType:       history.SyncType,
		StartDate:      history.StartDate,
		EndDate:        history.EndDate,
		Status:         history.Status,
		RecordsAdded:   history.RecordsAdded,
		RecordsUpdated: history.RecordsUpdated,
		Duration:       history.Duration,
		ErrorMsg:       history.ErrorMsg,
		CreatedAt:      time.Now(),
	}

	if err := r.db.Create(&entity).Error; err != nil {
		return fmt.Errorf("添加同步历史记录失败: %w", err)
	}

	history.ID = int(entity.ID)

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

	var entities []models.SyncHistoryEntity
	if err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询同步历史记录失败: %w", err)
	}

	histories := make([]*models.SyncHistory, 0)
	for _, entity := range entities {
		histories = append(histories, &models.SyncHistory{
			ID:             int(entity.ID),
			StockCode:      entity.StockCode,
			StockName:      entity.StockName,
			SyncType:       entity.SyncType,
			StartDate:      entity.StartDate,
			EndDate:        entity.EndDate,
			Status:         entity.Status,
			RecordsAdded:   entity.RecordsAdded,
			RecordsUpdated: entity.RecordsUpdated,
			Duration:       entity.Duration,
			ErrorMsg:       entity.ErrorMsg,
			CreatedAt:      entity.CreatedAt,
		})
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

	var entities []models.SyncHistoryEntity
	if err := r.db.Where("stock_code = ?", code).Order("created_at DESC").Limit(limit).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询股票同步历史记录失败: %w", err)
	}

	histories := make([]*models.SyncHistory, 0)
	for _, entity := range entities {
		histories = append(histories, &models.SyncHistory{
			ID:             int(entity.ID),
			StockCode:      entity.StockCode,
			StockName:      entity.StockName,
			SyncType:       entity.SyncType,
			StartDate:      entity.StartDate,
			EndDate:        entity.EndDate,
			Status:         entity.Status,
			RecordsAdded:   entity.RecordsAdded,
			RecordsUpdated: entity.RecordsUpdated,
			Duration:       entity.Duration,
			ErrorMsg:       entity.ErrorMsg,
			CreatedAt:      entity.CreatedAt,
		})
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

	var count int64
	if err := r.db.Model(&models.SyncHistoryEntity{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("查询同步历史记录总数失败: %w", err)
	}

	logger.Debug("成功获取同步历史记录总数",
		zap.String("module", "repositories.sync_history"),
		zap.String("op", "get_sync_history_count"),
		zap.Int64("count", count),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return int(count), nil
}

// ClearAll 清除所有同步历史记录
func (r *SQLiteSyncHistoryRepository) ClearAll() error {
	start := time.Now()

	result := r.db.Exec("DELETE FROM sync_history")
	if result.Error != nil {
		return fmt.Errorf("清除同步历史记录失败: %w", result.Error)
	}

	logger.Info("成功清除同步历史记录",
		zap.String("module", "repositories.sync_history"),
		zap.String("op", "clear_all_sync_history"),
		zap.Int64("rows_affected", result.RowsAffected),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return nil
}
