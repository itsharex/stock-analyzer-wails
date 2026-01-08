package services

import (
	"fmt"

	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"
)

// StrategyService 策略服务
type StrategyService struct {
	repo *repositories.StrategyRepository
}

// NewStrategyService 创建策略服务
func NewStrategyService(repo *repositories.StrategyRepository) *StrategyService {
	return &StrategyService{repo: repo}
}

// CreateStrategy 创建策略
func (s *StrategyService) CreateStrategy(name string, description string, strategyType string, parameters map[string]interface{}) (*models.StrategyConfig, error) {
	// 验证策略类型是否存在
	validStrategyType := false
	for _, st := range models.StrategyTypes {
		if st.Type == strategyType {
			validStrategyType = true
			break
		}
	}
	if !validStrategyType {
		return nil, fmt.Errorf("无效的策略类型: %s", strategyType)
	}

	strategy := &models.StrategyConfig{
		Name:         name,
		Description:  description,
		StrategyType: strategyType,
		Parameters:   parameters,
	}

	if err := s.repo.Create(strategy); err != nil {
		return nil, err
	}

	return strategy, nil
}

// UpdateStrategy 更新策略
func (s *StrategyService) UpdateStrategy(id int64, name string, description string, strategyType string, parameters map[string]interface{}) error {
	// 验证策略是否存在
	strategy, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if strategy == nil {
		return fmt.Errorf("策略不存在")
	}

	// 验证策略类型是否存在
	validStrategyType := false
	for _, st := range models.StrategyTypes {
		if st.Type == strategyType {
			validStrategyType = true
			break
		}
	}
	if !validStrategyType {
		return fmt.Errorf("无效的策略类型: %s", strategyType)
	}

	strategy.Name = name
	strategy.Description = description
	strategy.StrategyType = strategyType
	strategy.Parameters = parameters

	return s.repo.Update(strategy)
}

// DeleteStrategy 删除策略
func (s *StrategyService) DeleteStrategy(id int64) error {
	return s.repo.Delete(id)
}

// GetStrategy 获取策略
func (s *StrategyService) GetStrategy(id int64) (*models.StrategyConfig, error) {
	return s.repo.GetByID(id)
}

// GetAllStrategies 获取所有策略
func (s *StrategyService) GetAllStrategies() ([]models.StrategyConfig, error) {
	return s.repo.GetAll()
}

// GetStrategyTypes 获取所有策略类型定义
func (s *StrategyService) GetStrategyTypes() []models.StrategyTypeDefinition {
	return models.StrategyTypes
}

// UpdateStrategyBacktestResult 更新策略的回测结果
func (s *StrategyService) UpdateStrategyBacktestResult(id int64, backtestResult map[string]interface{}) error {
	return s.repo.UpdateBacktestResult(id, backtestResult)
}
