package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WatchlistRepository 定义自选股持久化接口
type WatchlistRepository interface {
	Add(stock *models.StockData) error
	Remove(code string) error
	GetAll() ([]*models.StockData, error)
}

// SQLiteWatchlistRepository 基于 SQLite 的实现
type SQLiteWatchlistRepository struct {
	db *gorm.DB
}

// NewSQLiteWatchlistRepository 构造函数
func NewSQLiteWatchlistRepository(db *gorm.DB) *SQLiteWatchlistRepository {
	return &SQLiteWatchlistRepository{db: db}
}

func (r *SQLiteWatchlistRepository) GetAll() ([]*models.StockData, error) {
	start := time.Now()

	var entities []models.WatchlistEntity
	if err := r.db.Order("added_at DESC").Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询自选股列表失败: %w", err)
	}

	stocks := make([]*models.StockData, 0)
	for _, entity := range entities {
		var stock models.StockData
		if entity.Data != "" {
			if err := json.Unmarshal([]byte(entity.Data), &stock); err != nil {
				logger.Error("解析自选股 JSON 数据失败", zap.String("code", entity.Code), zap.Error(err))
				continue
			}
		}
		stocks = append(stocks, &stock)
	}

	logger.Debug("成功获取自选股列表",
		zap.String("module", "repositories.watchlist"),
		zap.String("op", "get_all_watchlist"),
		zap.Int("stocks_count", len(stocks)),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return stocks, nil
}

func (r *SQLiteWatchlistRepository) Add(stock *models.StockData) error {
	start := time.Now()

	// 序列化 StockData
	data, err := json.Marshal(stock)
	if err != nil {
		return fmt.Errorf("序列化 StockData 失败: %w", err)
	}

	entity := models.WatchlistEntity{
		Code:    stock.Code,
		Name:    stock.Name,
		Data:    string(data),
		AddedAt: time.Now(),
	}

	// 使用 Clauses 实现 upsert (INSERT OR REPLACE)
	// 在 SQLite 中，ON CONFLICT(id) DO UPDATE SET ... 等价于 UPSERT
	// 这里我们更新除主键外的所有字段
	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "data", "added_at"}),
	}).Create(&entity).Error; err != nil {
		return fmt.Errorf("添加/更新自选股失败: %w", err)
	}

	logger.Info("成功添加/更新股票到自选股",
		zap.String("module", "repositories.watchlist"),
		zap.String("op", "add_to_watchlist"),
		zap.String("stock_code", stock.Code),
		zap.String("stock_name", stock.Name),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return nil
}

func (r *SQLiteWatchlistRepository) Remove(code string) error {
	start := time.Now()

	result := r.db.Delete(&models.WatchlistEntity{}, "code = ?", code)
	if result.Error != nil {
		return fmt.Errorf("从自选股移除股票失败: %w", result.Error)
	}

	logger.Info("成功从自选股移除股票",
		zap.String("module", "repositories.watchlist"),
		zap.String("op", "remove_from_watchlist"),
		zap.String("stock_code", code),
		zap.Int64("rows_affected", result.RowsAffected),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return nil
}
