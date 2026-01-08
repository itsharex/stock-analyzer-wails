package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"stock-analyzer-wails/models"

	_ "modernc.org/sqlite"
)

// StrategyRepository 策略仓库
type StrategyRepository struct {
	db *sql.DB
}

// NewStrategyRepository 创建策略仓库
func NewStrategyRepository(db *sql.DB) *StrategyRepository {
	return &StrategyRepository{db: db}
}

// Create 创建策略
func (r *StrategyRepository) Create(strategy *models.StrategyConfig) error {
	parametersJSON, err := json.Marshal(strategy.Parameters)
	if err != nil {
		return fmt.Errorf("序列化参数失败: %w", err)
	}

	var backtestJSON []byte
	if strategy.LastBacktestResult != nil {
		backtestJSON, err = json.Marshal(strategy.LastBacktestResult)
		if err != nil {
			return fmt.Errorf("序列化回测结果失败: %w", err)
		}
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	result, err := r.db.Exec(`
		INSERT INTO strategy_config (name, description, strategy_type, parameters, last_backtest_result, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, strategy.Name, strategy.Description, strategy.StrategyType, string(parametersJSON), string(backtestJSON), now, now)

	if err != nil {
		return fmt.Errorf("插入策略失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取插入ID失败: %w", err)
	}

	strategy.ID = id
	strategy.CreatedAt = now
	strategy.UpdatedAt = now

	return nil
}

// Update 更新策略
func (r *StrategyRepository) Update(strategy *models.StrategyConfig) error {
	parametersJSON, err := json.Marshal(strategy.Parameters)
	if err != nil {
		return fmt.Errorf("序列化参数失败: %w", err)
	}

	var backtestJSON []byte
	if strategy.LastBacktestResult != nil {
		backtestJSON, err = json.Marshal(strategy.LastBacktestResult)
		if err != nil {
			return fmt.Errorf("序列化回测结果失败: %w", err)
		}
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	_, err = r.db.Exec(`
		UPDATE strategy_config
		SET name = ?, description = ?, strategy_type = ?, parameters = ?, last_backtest_result = ?, updated_at = ?
		WHERE id = ?
	`, strategy.Name, strategy.Description, strategy.StrategyType, string(parametersJSON), string(backtestJSON), now, strategy.ID)

	if err != nil {
		return fmt.Errorf("更新策略失败: %w", err)
	}

	strategy.UpdatedAt = now
	return nil
}

// Delete 删除策略
func (r *StrategyRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM strategy_config WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除策略失败: %w", err)
	}
	return nil
}

// GetByID 根据 ID 获取策略
func (r *StrategyRepository) GetByID(id int64) (*models.StrategyConfig, error) {
	var strategy models.StrategyConfig
	var parametersJSON, backtestJSON sql.NullString

	err := r.db.QueryRow(`
		SELECT id, name, description, strategy_type, parameters, last_backtest_result, created_at, updated_at
		FROM strategy_config
		WHERE id = ?
	`, id).Scan(&strategy.ID, &strategy.Name, &strategy.Description, &strategy.StrategyType, &parametersJSON, &backtestJSON, &strategy.CreatedAt, &strategy.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("查询策略失败: %w", err)
	}

	if parametersJSON.Valid {
		if err := json.Unmarshal([]byte(parametersJSON.String), &strategy.Parameters); err != nil {
			return nil, fmt.Errorf("反序列化参数失败: %w", err)
		}
	}

	if backtestJSON.Valid {
		if err := json.Unmarshal([]byte(backtestJSON.String), &strategy.LastBacktestResult); err != nil {
			return nil, fmt.Errorf("反序列化回测结果失败: %w", err)
		}
	}

	return &strategy, nil
}

// GetAll 获取所有策略
func (r *StrategyRepository) GetAll() ([]models.StrategyConfig, error) {
	rows, err := r.db.Query(`
		SELECT id, name, description, strategy_type, parameters, last_backtest_result, created_at, updated_at
		FROM strategy_config
		ORDER BY updated_at DESC
	`)

	if err != nil {
		return nil, fmt.Errorf("查询策略列表失败: %w", err)
	}
	defer rows.Close()

	var strategies []models.StrategyConfig
	for rows.Next() {
		var strategy models.StrategyConfig
		var parametersJSON, backtestJSON sql.NullString

		if err := rows.Scan(&strategy.ID, &strategy.Name, &strategy.Description, &strategy.StrategyType, &parametersJSON, &backtestJSON, &strategy.CreatedAt, &strategy.UpdatedAt); err != nil {
			return nil, fmt.Errorf("扫描策略数据失败: %w", err)
		}

		if parametersJSON.Valid {
			if err := json.Unmarshal([]byte(parametersJSON.String), &strategy.Parameters); err != nil {
				return nil, fmt.Errorf("反序列化参数失败: %w", err)
			}
		}

		if backtestJSON.Valid {
			if err := json.Unmarshal([]byte(backtestJSON.String), &strategy.LastBacktestResult); err != nil {
				return nil, fmt.Errorf("反序列化回测结果失败: %w", err)
			}
		}

		strategies = append(strategies, strategy)
	}

	return strategies, nil
}

// UpdateBacktestResult 更新策略的最后回测结果
func (r *StrategyRepository) UpdateBacktestResult(id int64, backtestResult map[string]interface{}) error {
	backtestJSON, err := json.Marshal(backtestResult)
	if err != nil {
		return fmt.Errorf("序列化回测结果失败: %w", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	_, err = r.db.Exec(`
		UPDATE strategy_config
		SET last_backtest_result = ?, updated_at = ?
		WHERE id = ?
	`, string(backtestJSON), now, id)

	if err != nil {
		return fmt.Errorf("更新回测结果失败: %w", err)
	}

	return nil
}
