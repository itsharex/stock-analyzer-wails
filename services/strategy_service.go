package services

import (
	"encoding/json"
	"fmt"
	"math"
	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"

	"go.uber.org/zap"
)

// StrategyService 策略服务
type StrategyService struct {
	repo          *repositories.StrategyRepository
	moneyFlowRepo *repositories.MoneyFlowRepository
}

// NewStrategyService 创建策略服务
func NewStrategyService(repo *repositories.StrategyRepository, moneyFlowRepo *repositories.MoneyFlowRepository) *StrategyService {
	return &StrategyService{repo: repo, moneyFlowRepo: moneyFlowRepo}
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

// CalculateBuildSignals 计算建仓信号 (决策先锋算法)
func (s *StrategyService) CalculateBuildSignals(code string) (*models.StrategySignal, error) {
	if s.moneyFlowRepo == nil {
		return nil, fmt.Errorf("MoneyFlowRepository 未初始化")
	}

	// 使用 Repository 获取数据
	data, err := s.moneyFlowRepo.GetMoneyFlowHistory(code, 20)
	if err != nil {
		return nil, fmt.Errorf("查询资金流向数据失败: %w", err)
	}

	// 需要至少 20 天数据来计算 MA20
	if len(data) < 20 {
		return nil, nil // 数据不足，不报错，直接返回空
	}

	// data[0] 是最新一天 (T-0)
	// data[4] 是 T-4
	// data[19] 是 T-19

	// === A. 资金面：主力持续吸筹 (T-0 到 T-4) ===
	last5Days := data[0:5]
	positiveDays := 0
	netSum := 0.0
	absNetSum := 0.0

	for _, d := range last5Days {
		if d.MainNet > 0 {
			positiveDays++
		}
		netSum += d.MainNet
		absNetSum += math.Abs(d.MainNet)
	}

	// 1. 主力流入天数 >= 3
	if positiveDays < 3 {
		return nil, nil
	}

	// 2. 净额总和 > 0
	if netSum <= 0 {
		return nil, nil
	}

	// 3. 异动倍数 > 1.5
	avgAbsNet := absNetSum / 5.0
	currentAbsNet := math.Abs(data[0].MainNet)
	if currentAbsNet <= 1.5*avgAbsNet {
		return nil, nil
	}

	// === B. 技术面：趋势企稳回踩 (MA20) ===
	sumClose := 0.0
	for _, d := range data {
		sumClose += d.ClosePrice
	}
	ma20 := sumClose / 20.0

	currentClose := data[0].ClosePrice
	// 1. 站稳均线
	if currentClose < ma20 {
		return nil, nil
	}

	// 2. 回踩不追高 (0% <= 偏离度 <= 3%)
	deviation := (currentClose - ma20) / ma20
	if deviation < 0 || deviation > 0.03 {
		return nil, nil
	}

	// === C. 动能面：良性放量 ===
	// 涨幅在 0.5% 到 5% 之间
	currentChgPct := data[0].ChgPct
	if currentChgPct < 0.5 || currentChgPct > 5.0 {
		return nil, nil
	}

	// === 符合所有条件，生成信号 ===

	// 计算评分 Score = (过去5日主力流入总额 / 流通市值) * 100
	score := 0.0
	circMV, err := s.moneyFlowRepo.GetStockCircMV(code)
	if err == nil && circMV > 0 {
		score = (netSum / circMV) * 100
	} else {
		// 如果无法获取流通市值，使用备用评分逻辑
		score = deviation*100 + float64(positiveDays)*10
	}

	details := map[string]interface{}{
		"ma20":         ma20,
		"close":        currentClose,
		"deviation":    deviation,
		"positiveDays": positiveDays,
		"netSum":       netSum,
		"circMV":       circMV,
		"chgPct":       currentChgPct,
	}
	detailsJSON, _ := json.Marshal(details)

	signal := &models.StrategySignal{
		Code:         code,
		TradeDate:    data[0].TradeDate,
		SignalType:   "B",
		StrategyName: "决策先锋",
		Score:        score,
		Details:      string(detailsJSON),
	}

	// 持久化
	if err := s.moneyFlowRepo.SaveStrategySignal(signal); err != nil {
		logger.Error("保存策略信号失败", zap.Error(err))
		return nil, err
	}

	// 日志输出
	logger.Info(fmt.Sprintf("[信号发现] %s 均线回踩，主力近5日流入达%.2f万", code, netSum/10000),
		zap.String("code", code),
		zap.String("date", data[0].TradeDate),
		zap.Float64("score", score),
	)

	return signal, nil
}

// GetRecentMoneyFlows 获取近期资金流向数据
func (s *StrategyService) GetRecentMoneyFlows(code string, limit int) ([]models.MoneyFlowData, error) {
	if s.moneyFlowRepo == nil {
		return nil, fmt.Errorf("MoneyFlowRepository 未初始化")
	}
	return s.moneyFlowRepo.GetMoneyFlowHistory(code, limit)
}

// GetLatestSignals 获取最新的策略信号
func (s *StrategyService) GetLatestSignals(limit int) ([]models.StrategySignal, error) {
	if s.moneyFlowRepo == nil {
		return nil, fmt.Errorf("MoneyFlowRepository 未初始化")
	}
	return s.moneyFlowRepo.GetLatestSignals(limit)
}

// GetSignalsByStockCode 根据股票代码获取历史信号
func (s *StrategyService) GetSignalsByStockCode(code string) ([]models.StrategySignal, error) {
	if s.moneyFlowRepo == nil {
		return nil, fmt.Errorf("MoneyFlowRepository 未初始化")
	}
	return s.moneyFlowRepo.GetSignalsByStockCode(code)
}

// UpdateSignalAIResult 更新信号的 AI 分析结果
func (s *StrategyService) UpdateSignalAIResult(code, tradeDate, strategyName string, aiScore int, aiReason string) error {
	if s.moneyFlowRepo == nil {
		return fmt.Errorf("MoneyFlowRepository 未初始化")
	}
	return s.moneyFlowRepo.UpdateStrategySignalAI(code, tradeDate, strategyName, aiScore, aiReason)
}
