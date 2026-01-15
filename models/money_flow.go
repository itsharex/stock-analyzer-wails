package models

// MoneyFlowData 资金流向数据
type MoneyFlowData struct {
	Code       string  `json:"code"`
	TradeDate  string  `json:"tradeDate"`
	MainNet    float64 `json:"mainNet"`
	SuperNet   float64 `json:"superNet"`
	BigNet     float64 `json:"bigNet"`
	MidNet     float64 `json:"midNet"`
	SmallNet   float64 `json:"smallNet"`
	ClosePrice float64 `json:"closePrice"`
	ChgPct     float64 `json:"chgPct"`
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
