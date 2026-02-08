package repositories

import (
	"fmt"
	"stock-analyzer-wails/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MoneyFlowRepository 资金流向仓库
type MoneyFlowRepository struct {
	db *gorm.DB
}

// NewMoneyFlowRepository 创建资金流向仓库
func NewMoneyFlowRepository(db *gorm.DB) *MoneyFlowRepository {
	return &MoneyFlowRepository{db: db}
}

// SaveMoneyFlows 批量保存资金流向数据
func (r *MoneyFlowRepository) SaveMoneyFlows(flows []models.MoneyFlowData) error {
	var entities []models.StockMoneyFlowHistEntity
	for _, flow := range flows {
		entities = append(entities, models.StockMoneyFlowHistEntity{
			Code:       flow.Code,
			TradeDate:  flow.TradeDate,
			MainNet:    flow.MainNet,
			SuperNet:   flow.SuperNet,
			BigNet:     flow.BigNet,
			MidNet:     flow.MidNet,
			SmallNet:   flow.SmallNet,
			ClosePrice: flow.ClosePrice,
			ChgPct:     flow.ChgPct,
			Amount:     flow.Amount,
			MainRate:   flow.MainRate,
			Turnover:   flow.Turnover,
		})
	}

	if len(entities) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}, {Name: "trade_date"}},
		DoUpdates: clause.AssignmentColumns([]string{"main_net", "super_net", "big_net", "mid_net", "small_net", "close_price", "chg_pct"}),
	}).CreateInBatches(entities, 100).Error
}

// GetMoneyFlowHistory 获取最近的资金流向历史数据
func (r *MoneyFlowRepository) GetMoneyFlowHistory(code string, limit int) ([]models.MoneyFlowData, error) {
	var entities []models.StockMoneyFlowHistEntity
	if err := r.db.Where("code = ?", code).Order("trade_date DESC").Limit(limit).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询资金流向数据失败: %w", err)
	}

	var flows []models.MoneyFlowData
	for _, entity := range entities {
		flows = append(flows, models.MoneyFlowData{
			Code:       entity.Code,
			TradeDate:  entity.TradeDate,
			MainNet:    entity.MainNet,
			SuperNet:   entity.SuperNet,
			BigNet:     entity.BigNet,
			MidNet:     entity.MidNet,
			SmallNet:   entity.SmallNet,
			ClosePrice: entity.ClosePrice,
			ChgPct:     entity.ChgPct,
			Amount:     entity.Amount,   // 新增：成交金额
			MainRate:   entity.MainRate, // 新增：主力强度
			Turnover:   entity.Turnover, // 新增：换手率
		})
	}
	return flows, nil
}

// GetAllMoneyFlowHistory 获取所有的资金流向历史数据
func (r *MoneyFlowRepository) GetAllMoneyFlowHistory(code string) ([]models.MoneyFlowData, error) {
	var entities []models.StockMoneyFlowHistEntity
	if err := r.db.Where("code = ?", code).Order("trade_date ASC").Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询所有资金流向数据失败: %w", err)
	}

	var flows []models.MoneyFlowData
	for _, entity := range entities {
		flows = append(flows, models.MoneyFlowData{
			Code:       entity.Code,
			TradeDate:  entity.TradeDate,
			MainNet:    entity.MainNet,
			SuperNet:   entity.SuperNet,
			BigNet:     entity.BigNet,
			MidNet:     entity.MidNet,
			SmallNet:   entity.SmallNet,
			ClosePrice: entity.ClosePrice,
			ChgPct:     entity.ChgPct,
			Amount:     entity.Amount,   // 新增：成交金额
			MainRate:   entity.MainRate, // 新增：主力强度
			Turnover:   entity.Turnover, // 新增：换手率
		})
	}
	return flows, nil
}

// GetStockCircMV 获取股票流通市值
func (r *MoneyFlowRepository) GetStockCircMV(code string) (float64, error) {
	var circMV float64
	// 直接查询 stocks 表
	err := r.db.Table("stocks").Select("circ_mv").Where("code = ?", code).Scan(&circMV).Error
	if err != nil {
		return 0, fmt.Errorf("查询流通市值失败: %w", err)
	}
	return circMV, nil
}

// SaveStrategySignal 保存策略信号
func (r *MoneyFlowRepository) SaveStrategySignal(signal *models.StrategySignal) error {
	entity := models.StockStrategySignalEntity{
		Code:         signal.Code,
		TradeDate:    signal.TradeDate,
		SignalType:   signal.SignalType,
		StrategyName: signal.StrategyName,
		Score:        signal.Score,
		Details:      signal.Details,
		AIScore:      signal.AIScore,
		AIReason:     signal.AIReason,
	}

	// INSERT OR IGNORE
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&entity).Error; err != nil {
		return fmt.Errorf("保存策略信号失败: %w", err)
	}

	// 如果插入成功，获取ID (如果是 IGNORE 且未插入，ID 可能是 0，但 GORM 通常会回填 ID 如果是新插入的)
	// 如果是 DoNothing 且冲突，entity.ID 不会被赋值为现有记录的 ID
	if entity.ID != 0 {
		signal.ID = int64(entity.ID)
	}

	return nil
}

// UpdateStrategySignalAI 更新策略信号的 AI 评分和理由
func (r *MoneyFlowRepository) UpdateStrategySignalAI(code, tradeDate, strategyName string, aiScore int, aiReason string) error {
	result := r.db.Model(&models.StockStrategySignalEntity{}).
		Where("code = ? AND trade_date = ? AND strategy_name = ?", code, tradeDate, strategyName).
		Updates(map[string]interface{}{
			"ai_score":  aiScore,
			"ai_reason": aiReason,
		})

	if result.Error != nil {
		return fmt.Errorf("更新 AI 评分失败: %w", result.Error)
	}
	return nil
}

// GetLatestSignals 获取最新的策略信号
func (r *MoneyFlowRepository) GetLatestSignals(limit int) ([]models.StrategySignal, error) {
	var results []struct {
		models.StockStrategySignalEntity
		StockName string `gorm:"column:stock_name"`
	}

	// 联表查询
	err := r.db.Table("stock_strategy_signals s").
		Select("s.*, IFNULL(st.name, s.code) as stock_name").
		Joins("LEFT JOIN stocks st ON s.code = st.code").
		Order("s.created_at DESC, s.id DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("查询策略信号失败: %w", err)
	}

	var signals []models.StrategySignal
	for _, res := range results {
		signals = append(signals, models.StrategySignal{
			ID:           int64(res.ID),
			Code:         res.Code,
			StockName:    res.StockName,
			TradeDate:    res.TradeDate,
			SignalType:   res.SignalType,
			StrategyName: res.StrategyName,
			Score:        res.Score,
			Details:      res.Details,
			AIScore:      res.AIScore,
			AIReason:     res.AIReason,
			CreatedAt:    res.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return signals, nil
}

// GetStockName 根据股票代码查询股票名称
func (r *MoneyFlowRepository) GetStockName(code string) (string, error) {
	var result struct {
		Name string `gorm:"column:name"`
	}

	err := r.db.Table("stocks").
		Select("name").
		Where("code = ?", code).
		First(&result).Error

	if err != nil {
		return "", fmt.Errorf("查询股票名称失败: %w", err)
	}

	return result.Name, nil
}

// GetSignalsByDateRange 根据日期范围获取历史信号
func (r *MoneyFlowRepository) GetSignalsByDateRange(startDate, endDate string) ([]models.StrategySignal, error) {
	var results []struct {
		models.StockStrategySignalEntity
		StockName string `gorm:"column:stock_name"`
	}

	err := r.db.Table("stock_strategy_signals s").
		Select("s.*, IFNULL(st.name, s.code) as stock_name").
		Joins("LEFT JOIN stocks st ON s.code = st.code").
		Where("s.trade_date >= ? AND s.trade_date <= ? AND s.strategy_name = '决策先锋'", startDate, endDate).
		Order("s.trade_date DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("查询策略信号失败: %w", err)
	}

	var signals []models.StrategySignal
	for _, res := range results {
		signals = append(signals, models.StrategySignal{
			ID:           int64(res.ID),
			Code:         res.Code,
			StockName:    res.StockName,
			TradeDate:    res.TradeDate,
			SignalType:   res.SignalType,
			StrategyName: res.StrategyName,
			Score:        res.Score,
			Details:      res.Details,
			AIScore:      res.AIScore,
			AIReason:     res.AIReason,
			CreatedAt:    res.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return signals, nil
}

// GetSignalsByStockCode 根据股票代码获取历史信号
func (r *MoneyFlowRepository) GetSignalsByStockCode(code string) ([]models.StrategySignal, error) {
	var results []struct {
		models.StockStrategySignalEntity
		StockName string `gorm:"column:stock_name"`
	}

	err := r.db.Table("stock_strategy_signals s").
		Select("s.*, IFNULL(st.name, s.code) as stock_name").
		Joins("LEFT JOIN stocks st ON s.code = st.code").
		Where("s.code = ?", code).
		Order("s.trade_date DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("查询策略信号失败: %w", err)
	}

	var signals []models.StrategySignal
	for _, res := range results {
		signals = append(signals, models.StrategySignal{
			ID:           int64(res.ID),
			Code:         res.Code,
			StockName:    res.StockName,
			TradeDate:    res.TradeDate,
			SignalType:   res.SignalType,
			StrategyName: res.StrategyName,
			Score:        res.Score,
			Details:      res.Details,
			AIScore:      res.AIScore,
			AIReason:     res.AIReason,
			CreatedAt:    res.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return signals, nil
}
