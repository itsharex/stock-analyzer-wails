package models

// MoneyFlowData 资金流向数据
type MoneyFlowData struct {
	Code       string  `json:"code"`
	TradeDate  string  `json:"tradeDate"`
	MainNet    float64 `json:"mainNet"`  // 主力净额
	SuperNet   float64 `json:"superNet"` // 超大单净额
	BigNet     float64 `json:"bigNet"`   // 大单
	MidNet     float64 `json:"midNet"`   // 中单
	SmallNet   float64 `json:"smallNet"` // 小单净额
	ClosePrice float64 `json:"closePrice"`
	ChgPct     float64 `json:"chgPct"`   // 涨跌幅
	Amount     float64 `json:"amount"`   // 成交金额
	MainRate   float64 `json:"mainRate"` // 主力强度 (百分比)
	Turnover   float64 `json:"turnover"` // 换手率
}

// StrategySignal 策略信号
type StrategySignal struct {
	ID           int64   `json:"id"`
	Code         string  `json:"code"`
	StockName    string  `json:"stockName"` // 股票名称
	TradeDate    string  `json:"tradeDate"`
	SignalType   string  `json:"signalType"` // 'B' for Buy, 'S' for Sell
	StrategyName string  `json:"strategyName"`
	Score        float64 `json:"score"`
	Details      string  `json:"details"` // JSON string
	AIScore      int     `json:"aiScore"`
	AIReason     string  `json:"aiReason"`
	CreatedAt    string  `json:"createdAt"`
}

// AIVerificationResult AI 验证结果
type AIVerificationResult struct {
	Score     int    `json:"score"`
	Opinion   string `json:"opinion"`
	RiskLevel string `json:"risk_level"`
}
