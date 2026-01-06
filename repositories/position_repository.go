package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"stock-analyzer-wails/models"
	"time"

	"stock-analyzer-wails/internal/logger"
	"go.uber.org/zap"
)

// PositionRepository 定义持仓记录持久化接口
type PositionRepository interface {
	SavePosition(pos *models.Position) error
	GetPositions() (map[string]*models.Position, error)
	RemovePosition(code string) error
}

// SQLitePositionRepository 基于 SQLite 的实现
type SQLitePositionRepository struct {
	db *sql.DB
}

// NewSQLitePositionRepository 构造函数
func NewSQLitePositionRepository(db *sql.DB) *SQLitePositionRepository {
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

	query := `
		INSERT OR REPLACE INTO positions (
			stock_code, stock_name, entry_price, entry_time, current_status, logic_status, 
			strategy_json, trailing_config_json, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.Exec(query,
		pos.StockCode,
		pos.StockName,
		pos.EntryPrice,
		pos.EntryTime,
		pos.CurrentStatus,
		pos.LogicStatus,
		string(strategyJSON),
		string(trailingConfigJSON),
		pos.UpdatedAt,
	)
	if err != nil {
		logger.Error("保存持仓数据失败", zap.Error(err), zap.String("code", pos.StockCode))
		return fmt.Errorf("保存持仓数据失败: %w", err)
	}

	return nil
}

// GetPositions 获取所有持仓记录
func (r *SQLitePositionRepository) GetPositions() (map[string]*models.Position, error) {
	rows, err := r.db.Query(`
		SELECT 
			stock_code, stock_name, entry_price, entry_time, current_status, logic_status, 
			strategy_json, trailing_config_json, updated_at
		FROM positions
	`)
	if err != nil {
		return nil, fmt.Errorf("查询持仓记录失败: %w", err)
	}
	defer rows.Close()

	positions := make(map[string]*models.Position)
	for rows.Next() {
		pos := &models.Position{}
		var strategyJSON, trailingConfigJSON string

		err := rows.Scan(
			&pos.StockCode,
			&pos.StockName,
			&pos.EntryPrice,
			&pos.EntryTime,
			&pos.CurrentStatus,
			&pos.LogicStatus,
			&strategyJSON,
			&trailingConfigJSON,
			&pos.UpdatedAt,
		)
		if err != nil {
			logger.Error("扫描持仓数据失败", zap.Error(err))
			continue
		}

		// 反序列化嵌套的 JSON 字段
		if err := json.Unmarshal([]byte(strategyJSON), &pos.Strategy); err != nil {
			logger.Error("反序列化 Strategy 失败", zap.String("code", pos.StockCode), zap.Error(err))
			continue
		}
		if err := json.Unmarshal([]byte(trailingConfigJSON), &pos.TrailingConfig); err != nil {
			logger.Error("反序列化 TrailingConfig 失败", zap.String("code", pos.StockCode), zap.Error(err))
			continue
		}

		positions[pos.StockCode] = pos
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历持仓记录结果集失败: %w", err)
	}

	return positions, nil
}

// RemovePosition 移除持仓记录
func (r *SQLitePositionRepository) RemovePosition(code string) error {
	query := `DELETE FROM positions WHERE stock_code = ?`
	_, err := r.db.Exec(query, code)
	if err != nil {
		return fmt.Errorf("移除持仓记录失败: %w", err)
	}
	return nil
}
