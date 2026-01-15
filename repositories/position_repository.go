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

// PositionRepository 定义持仓记录持久化接口
type PositionRepository interface {
	SavePosition(pos *models.Position) error
	GetPositions() (map[string]*models.Position, error)
	RemovePosition(code string) error
}

// SQLitePositionRepository 基于 SQLite 的实现
type SQLitePositionRepository struct {
	db *gorm.DB
}

// NewSQLitePositionRepository 构造函数
func NewSQLitePositionRepository(db *gorm.DB) *SQLitePositionRepository {
	return &SQLitePositionRepository{db: db}
}

// SavePosition 保存或更新持仓记录
func (r *SQLitePositionRepository) SavePosition(pos *models.Position) error {
	// 序列化嵌套的 JSON 字段
	strategyJSON, err := json.Marshal(pos.Strategy)
	if err != nil {
		return fmt.Errorf("序列化 Strategy 失败: %w", err)
	}
	trailingConfigJSON, err := json.Marshal(pos.TrailingConfig)
	if err != nil {
		return fmt.Errorf("序列化 TrailingConfig 失败: %w", err)
	}

	pos.UpdatedAt = time.Now()

	entity := models.PositionEntity{
		StockCode:          pos.StockCode,
		StockName:          pos.StockName,
		EntryPrice:         pos.EntryPrice,
		EntryTime:          pos.EntryTime,
		CurrentStatus:      pos.CurrentStatus,
		LogicStatus:        pos.LogicStatus,
		StrategyJSON:       string(strategyJSON),
		TrailingConfigJSON: string(trailingConfigJSON),
		UpdatedAt:          pos.UpdatedAt,
	}

	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "stock_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"stock_name", "entry_price", "entry_time", "current_status", "logic_status", "strategy_json", "trailing_config_json", "updated_at"}),
	}).Create(&entity).Error; err != nil {
		logger.Error("保存持仓数据失败", zap.Error(err), zap.String("code", pos.StockCode))
		return fmt.Errorf("保存持仓数据失败: %w", err)
	}

	return nil
}

// GetPositions 获取所有持仓记录
func (r *SQLitePositionRepository) GetPositions() (map[string]*models.Position, error) {
	var entities []models.PositionEntity
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询持仓记录失败: %w", err)
	}

	positions := make(map[string]*models.Position)
	for _, entity := range entities {
		pos := &models.Position{
			StockCode:     entity.StockCode,
			StockName:     entity.StockName,
			EntryPrice:    entity.EntryPrice,
			EntryTime:     entity.EntryTime,
			CurrentStatus: entity.CurrentStatus,
			LogicStatus:   entity.LogicStatus,
			UpdatedAt:     entity.UpdatedAt,
		}

		// 反序列化嵌套的 JSON 字段
		if err := json.Unmarshal([]byte(entity.StrategyJSON), &pos.Strategy); err != nil {
			logger.Error("反序列化 Strategy 失败", zap.String("code", pos.StockCode), zap.Error(err))
			continue
		}
		if err := json.Unmarshal([]byte(entity.TrailingConfigJSON), &pos.TrailingConfig); err != nil {
			logger.Error("反序列化 TrailingConfig 失败", zap.String("code", pos.StockCode), zap.Error(err))
			continue
		}

		positions[pos.StockCode] = pos
	}

	return positions, nil
}

// RemovePosition 移除持仓记录
func (r *SQLitePositionRepository) RemovePosition(code string) error {
	if err := r.db.Delete(&models.PositionEntity{}, "stock_code = ?", code).Error; err != nil {
		return fmt.Errorf("移除持仓记录失败: %w", err)
	}
	return nil
}
