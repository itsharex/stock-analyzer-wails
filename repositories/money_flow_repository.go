package repositories

import (
	"database/sql"
	"fmt"
	"stock-analyzer-wails/models"
)

// MoneyFlowRepository 资金流向仓库
type MoneyFlowRepository struct {
	db *sql.DB
}

// NewMoneyFlowRepository 创建资金流向仓库
func NewMoneyFlowRepository(db *sql.DB) *MoneyFlowRepository {
	return &MoneyFlowRepository{db: db}
}

// SaveMoneyFlows 批量保存资金流向数据
func (r *MoneyFlowRepository) SaveMoneyFlows(flows []models.MoneyFlowData) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()

	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO stock_money_flow_hist 
		(code, trade_date, main_net, super_net, big_net, mid_net, small_net, close_price, chg_pct) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("准备 SQL 失败: %w", err)
	}
	defer stmt.Close()

	for _, flow := range flows {
		_, err := stmt.Exec(
			flow.Code,
			flow.TradeDate,
			flow.MainNet,
			flow.SuperNet,
			flow.BigNet,
			flow.MidNet,
			flow.SmallNet,
			flow.ClosePrice,
			flow.ChgPct,
		)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("写入数据失败: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// GetMoneyFlowHistory 获取最近的资金流向历史数据
func (r *MoneyFlowRepository) GetMoneyFlowHistory(code string, limit int) ([]models.MoneyFlowData, error) {
	rows, err := r.db.Query(`
		SELECT code, trade_date, main_net, super_net, big_net, mid_net, small_net, close_price, chg_pct 
		FROM stock_money_flow_hist 
		WHERE code = ? 
		ORDER BY trade_date DESC 
		LIMIT ?
	`, code, limit)
	if err != nil {
		return nil, fmt.Errorf("查询资金流向数据失败: %w", err)
	}
	defer rows.Close()

	var flows []models.MoneyFlowData
	for rows.Next() {
		var flow models.MoneyFlowData
		if err := rows.Scan(
			&flow.Code,
			&flow.TradeDate,
			&flow.MainNet,
			&flow.SuperNet,
			&flow.BigNet,
			&flow.MidNet,
			&flow.SmallNet,
			&flow.ClosePrice,
			&flow.ChgPct,
		); err != nil {
			return nil, fmt.Errorf("扫描数据失败: %w", err)
		}
		flows = append(flows, flow)
	}

	return flows, nil
}

// GetStockCircMV 获取股票流通市值
func (r *MoneyFlowRepository) GetStockCircMV(code string) (float64, error) {
	var circMV float64
	err := r.db.QueryRow("SELECT circ_mv FROM stocks WHERE code = ?", code).Scan(&circMV)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // 如果没有找到股票，返回0
		}
		return 0, fmt.Errorf("查询流通市值失败: %w", err)
	}
	return circMV, nil
}

// SaveStrategySignal 保存策略信号
func (r *MoneyFlowRepository) SaveStrategySignal(signal *models.StrategySignal) error {
	result, err := r.db.Exec(`
		INSERT OR IGNORE INTO stock_strategy_signals (code, trade_date, signal_type, strategy_name, score, details)
		VALUES (?, ?, ?, ?, ?, ?)
	`, signal.Code, signal.TradeDate, signal.SignalType, signal.StrategyName, signal.Score, signal.Details)

	if err != nil {
		return fmt.Errorf("保存策略信号失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		signal.ID = id
	}

	return nil
}

// UpdateStrategySignalAI 更新策略信号的 AI 评分和理由
func (r *MoneyFlowRepository) UpdateStrategySignalAI(code, tradeDate, strategyName string, aiScore int, aiReason string) error {
	_, err := r.db.Exec(`
		UPDATE stock_strategy_signals 
		SET ai_score = ?, ai_reason = ? 
		WHERE code = ? AND trade_date = ? AND strategy_name = ?
	`, aiScore, aiReason, code, tradeDate, strategyName)

	if err != nil {
		return fmt.Errorf("更新 AI 评分失败: %w", err)
	}
	return nil
}

// GetLatestSignals 获取最新的策略信号
func (r *MoneyFlowRepository) GetLatestSignals(limit int) ([]models.StrategySignal, error) {
	rows, err := r.db.Query(`
		SELECT s.id, s.code, IFNULL(st.name, s.code) as stock_name, s.trade_date, s.signal_type, s.strategy_name, s.score, s.details, s.ai_score, s.ai_reason, s.created_at
		FROM stock_strategy_signals s
		LEFT JOIN stocks st ON s.code = st.code
		ORDER BY s.created_at DESC, s.id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("查询策略信号失败: %w", err)
	}
	defer rows.Close()

	var signals []models.StrategySignal
	for rows.Next() {
		var sig models.StrategySignal
		var createdAt sql.NullString
		if err := rows.Scan(
			&sig.ID,
			&sig.Code,
			&sig.StockName,
			&sig.TradeDate,
			&sig.SignalType,
			&sig.StrategyName,
			&sig.Score,
			&sig.Details,
			&sig.AIScore,
			&sig.AIReason,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("扫描信号数据失败: %w", err)
		}
		if createdAt.Valid {
			sig.CreatedAt = createdAt.String
		}
		signals = append(signals, sig)
	}

	return signals, nil
}

// GetSignalsByStockCode 根据股票代码获取历史信号
func (r *MoneyFlowRepository) GetSignalsByStockCode(code string) ([]models.StrategySignal, error) {
	rows, err := r.db.Query(`
		SELECT s.id, s.code, IFNULL(st.name, s.code) as stock_name, s.trade_date, s.signal_type, s.strategy_name, s.score, s.details, s.ai_score, s.ai_reason, s.created_at
		FROM stock_strategy_signals s
		LEFT JOIN stocks st ON s.code = st.code
		WHERE s.code = ?
		ORDER BY s.trade_date DESC
	`, code)
	if err != nil {
		return nil, fmt.Errorf("查询策略信号失败: %w", err)
	}
	defer rows.Close()

	var signals []models.StrategySignal
	for rows.Next() {
		var sig models.StrategySignal
		var createdAt sql.NullString
		if err := rows.Scan(
			&sig.ID,
			&sig.Code,
			&sig.StockName,
			&sig.TradeDate,
			&sig.SignalType,
			&sig.StrategyName,
			&sig.Score,
			&sig.Details,
			&sig.AIScore,
			&sig.AIReason,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("扫描信号数据失败: %w", err)
		}
		if createdAt.Valid {
			sig.CreatedAt = createdAt.String
		}
		signals = append(signals, sig)
	}

	return signals, nil
}
