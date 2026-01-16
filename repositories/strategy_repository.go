package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"stock-analyzer-wails/models"

	"gorm.io/gorm"
)

// StrategyRepository 策略仓库
type StrategyRepository struct {
	db *gorm.DB
}

// NewStrategyRepository 创建策略仓库
func NewStrategyRepository(db *gorm.DB) *StrategyRepository {
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

	now := time.Now()

	entity := models.StrategyConfigEntity{
		Name:             strategy.Name,
		Description:      strategy.Description,
		StrategyType:     strategy.StrategyType,
		Parameters:       string(parametersJSON),
		LastBacktestResult: string(backtestJSON),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := r.db.Create(&entity).Error; err != nil {
		return fmt.Errorf("插入策略失败: %w", err)
	}

	strategy.ID = int64(entity.ID)
	strategy.CreatedAt = now.Format("2006-01-02 15:04:05")
	strategy.UpdatedAt = now.Format("2006-01-02 15:04:05")

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

	now := time.Now()

	updates := map[string]interface{}{
		"name":               strategy.Name,
		"description":         strategy.Description,
		"strategy_type":       strategy.StrategyType,
		"parameters":          string(parametersJSON),
		"last_backtest_result": string(backtestJSON),
		"updated_at":          now,
	}

	if err := r.db.Model(&models.StrategyConfigEntity{}).Where("id = ?", strategy.ID).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新策略失败: %w", err)
	}

	strategy.UpdatedAt = now.Format("2006-01-02 15:04:05")
	return nil
}

// Delete 删除策略
func (r *StrategyRepository) Delete(id int64) error {
	if err := r.db.Delete(&models.StrategyConfigEntity{}, id).Error; err != nil {
		return fmt.Errorf("删除策略失败: %w", err)
	}
	return nil
}

// GetByID 根据 ID 获取策略
func (r *StrategyRepository) GetByID(id int64) (*models.StrategyConfig, error) {
	var entity models.StrategyConfigEntity
	if err := r.db.First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询策略失败: %w", err)
	}

	strategy := &models.StrategyConfig{
		ID:           int64(entity.ID),
		Name:         entity.Name,
		Description:  entity.Description,
		StrategyType: entity.StrategyType,
		CreatedAt:    entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    entity.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if entity.Parameters != "" {
		if err := json.Unmarshal([]byte(entity.Parameters), &strategy.Parameters); err != nil {
			return nil, fmt.Errorf("反序列化参数失败: %w", err)
		}
	}

	if entity.LastBacktestResult != "" {
		if err := json.Unmarshal([]byte(entity.LastBacktestResult), &strategy.LastBacktestResult); err != nil {
			return nil, fmt.Errorf("反序列化回测结果失败: %w", err)
		}
	}

	return strategy, nil
}

// GetAll 获取所有策略
func (r *StrategyRepository) GetAll() ([]models.StrategyConfig, error) {
	var entities []models.StrategyConfigEntity
	if err := r.db.Order("updated_at DESC").Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询策略列表失败: %w", err)
	}

	strategies := make([]models.StrategyConfig, 0)
	for _, entity := range entities {
		strategy := models.StrategyConfig{
			ID:           int64(entity.ID),
			Name:         entity.Name,
			Description:  entity.Description,
			StrategyType: entity.StrategyType,
			CreatedAt:    entity.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if entity.Parameters != "" {
			if err := json.Unmarshal([]byte(entity.Parameters), &strategy.Parameters); err != nil {
				return nil, fmt.Errorf("反序列化参数失败: %w", err)
			}
		}

		if entity.LastBacktestResult != "" {
			if err := json.Unmarshal([]byte(entity.LastBacktestResult), &strategy.LastBacktestResult); err != nil {
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

	now := time.Now()

	if err := r.db.Model(&models.StrategyConfigEntity{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_backtest_result": string(backtestJSON),
			"updated_at":           now,
		}).Error; err != nil {
		return fmt.Errorf("更新回测结果失败: %w", err)
	}

	return nil
}
