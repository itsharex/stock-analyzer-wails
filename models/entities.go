package models

import (
	"time"
)

// WatchlistEntity 对应 watchlist 表
type WatchlistEntity struct {
	Code    string    `gorm:"primaryKey;column:code" json:"code"`
	Name    string    `gorm:"column:name;not null" json:"name"`
	Data    string    `gorm:"column:data;not null" json:"data"` // 存储 StockData 的 JSON 字符串
	AddedAt time.Time `gorm:"column:added_at;default:CURRENT_TIMESTAMP" json:"addedAt"`
}

func (WatchlistEntity) TableName() string {
	return "watchlist"
}

// AlertEntity 对应 alerts 表
type AlertEntity struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	StockCode     string    `gorm:"column:stock_code;not null;uniqueIndex:idx_alert_unique" json:"stockCode"`
	StockName     string    `gorm:"column:stock_name;not null" json:"stockName"`
	Price         float64   `gorm:"column:price;not null;uniqueIndex:idx_alert_unique" json:"price"`
	Type          string    `gorm:"column:type;not null;uniqueIndex:idx_alert_unique" json:"type"` // 'above' or 'below'
	IsActive      bool      `gorm:"column:is_active;not null" json:"isActive"`
	LastTriggered time.Time `gorm:"column:last_triggered" json:"lastTriggered"`
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (AlertEntity) TableName() string {
	return "alerts"
}

// AlertHistoryEntity 对应 alert_history 表
type AlertHistoryEntity struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	StockCode      string    `gorm:"column:stock_code;not null" json:"stockCode"`
	StockName      string    `gorm:"column:stock_name;not null" json:"stockName"`
	TriggeredPrice float64   `gorm:"column:triggered_price;not null" json:"triggeredPrice"`
	Message        string    `gorm:"column:message;not null" json:"message"`
	TriggeredAt    time.Time `gorm:"column:triggered_at;default:CURRENT_TIMESTAMP" json:"triggeredAt"`
}

func (AlertHistoryEntity) TableName() string {
	return "alert_history"
}

// PositionEntity 对应 positions 表
type PositionEntity struct {
	StockCode          string    `gorm:"primaryKey;column:stock_code" json:"stockCode"`
	StockName          string    `gorm:"column:stock_name;not null" json:"stockName"`
	EntryPrice         float64   `gorm:"column:entry_price;not null" json:"entryPrice"`
	EntryTime          time.Time `gorm:"column:entry_time;not null" json:"entryTime"`
	CurrentStatus      string    `gorm:"column:current_status;not null" json:"currentStatus"` // 'holding', 'closed'
	LogicStatus        string    `gorm:"column:logic_status;not null" json:"logicStatus"`     // 'valid', 'violated'
	StrategyJSON       string    `gorm:"column:strategy_json;not null" json:"strategyJson"`           // 存储 EntryStrategyResult 的 JSON 字符串
	TrailingConfigJSON string    `gorm:"column:trailing_config_json;not null" json:"trailingConfigJson"` // 存储 TrailingStopConfig 的 JSON 字符串
	UpdatedAt          time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (PositionEntity) TableName() string {
	return "positions"
}

// ConfigEntity 对应 config 表
type ConfigEntity struct {
	Key   string `gorm:"primaryKey;column:key" json:"key"`
	Value string `gorm:"column:value;not null" json:"value"`
}

func (ConfigEntity) TableName() string {
	return "config"
}

// SyncHistoryEntity 对应 sync_history 表
type SyncHistoryEntity struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	StockCode      string    `gorm:"column:stock_code;not null" json:"stockCode"`
	StockName      string    `gorm:"column:stock_name;not null" json:"stockName"`
	SyncType       string    `gorm:"column:sync_type;not null" json:"syncType"` // 'single' or 'batch'
	StartDate      string    `gorm:"column:start_date;not null" json:"startDate"`
	EndDate        string    `gorm:"column:end_date;not null" json:"endDate"`
	Status         string    `gorm:"column:status;not null" json:"status"` // 'success' or 'failed'
	RecordsAdded   int       `gorm:"column:records_added;default:0" json:"recordsAdded"`
	RecordsUpdated int       `gorm:"column:records_updated;default:0" json:"recordsUpdated"`
	Duration       int       `gorm:"column:duration;default:0" json:"duration"` // 耗时（秒）
	ErrorMsg       string    `gorm:"column:error_msg" json:"errorMsg"`
	CreatedAt      time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (SyncHistoryEntity) TableName() string {
	return "sync_history"
}

// StrategyConfigEntity 对应 strategy_config 表
type StrategyConfigEntity struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name               string    `gorm:"column:name;not null" json:"name"`
	Description        string    `gorm:"column:description" json:"description"`
	StrategyType       string    `gorm:"column:strategy_type;not null" json:"strategyType"`
	Parameters         string    `gorm:"column:parameters;not null" json:"parameters"` // JSON 格式
	LastBacktestResult string    `gorm:"column:last_backtest_result" json:"lastBacktestResult"` // JSON
	CreatedAt          time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (StrategyConfigEntity) TableName() string {
	return "strategy_config"
}

// StockEntity 对应 stocks 表
type StockEntity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code         string    `gorm:"column:code;not null;uniqueIndex" json:"code"`
	Name         string    `gorm:"column:name;not null;index" json:"name"`
	Market       string    `gorm:"column:market;not null;index" json:"market"`
	FullCode     string    `gorm:"column:full_code;not null;uniqueIndex" json:"fullCode"`
	Type         string    `gorm:"column:type" json:"type"`
	IsActive     int       `gorm:"column:is_active;default:1" json:"isActive"`
	Price        float64   `gorm:"column:price" json:"price"`
	ChangeRate   float64   `gorm:"column:change_rate" json:"changeRate"`
	ChangeAmount float64   `gorm:"column:change_amount" json:"changeAmount"`
	Volume       float64   `gorm:"column:volume" json:"volume"`
	Amount       float64   `gorm:"column:amount" json:"amount"`
	Amplitude    float64   `gorm:"column:amplitude" json:"amplitude"`
	High         float64   `gorm:"column:high" json:"high"`
	Low          float64   `gorm:"column:low" json:"low"`
	Open         float64   `gorm:"column:open" json:"open"`
	PreClose     float64   `gorm:"column:pre_close" json:"preClose"`
	Turnover     float64   `gorm:"column:turnover" json:"turnover"`
	VolumeRatio  float64   `gorm:"column:volume_ratio" json:"volumeRatio"`
	PE           float64   `gorm:"column:pe" json:"pe"`
	WarrantRatio float64   `gorm:"column:warrant_ratio" json:"warrantRatio"`
	Industry     string    `gorm:"column:industry" json:"industry"`
	Region       string    `gorm:"column:region" json:"region"`
	Board        string    `gorm:"column:board" json:"board"`
	TotalMV      float64   `gorm:"column:total_mv" json:"totalMV"`
	CircMV       float64   `gorm:"column:circ_mv" json:"circMV"`
	UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (StockEntity) TableName() string {
	return "stocks"
}

// PriceThresholdAlertEntity 对应 price_threshold_alerts 表
type PriceThresholdAlertEntity struct {
	ID                uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	StockCode         string     `gorm:"column:stock_code;not null;index" json:"stockCode"`
	StockName         string     `gorm:"column:stock_name;not null" json:"stockName"`
	AlertType         string     `gorm:"column:alert_type;not null" json:"alertType"`
	Conditions        string     `gorm:"column:conditions;not null" json:"conditions"` // JSON
	IsActive          bool       `gorm:"column:is_active;default:true;index" json:"isActive"`
	Sensitivity       float64    `gorm:"column:sensitivity;default:0.001" json:"sensitivity"`
	CooldownHours     int        `gorm:"column:cooldown_hours;default:1" json:"cooldownHours"`
	PostTriggerAction string     `gorm:"column:post_trigger_action;default:'continue'" json:"postTriggerAction"`
	EnableSound       bool       `gorm:"column:enable_sound;default:true" json:"enableSound"`
	EnableDesktop     bool       `gorm:"column:enable_desktop;default:true" json:"enableDesktop"`
	TemplateID        string     `gorm:"column:template_id" json:"templateId"`
	CreatedAt         time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	LastTriggeredAt   *time.Time `gorm:"column:last_triggered_at" json:"lastTriggeredAt"`
}

func (PriceThresholdAlertEntity) TableName() string {
	return "price_threshold_alerts"
}

// PriceAlertTemplateEntity 对应 price_alert_templates 表
type PriceAlertTemplateEntity struct {
	ID          string    `gorm:"primaryKey;column:id" json:"id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	AlertType   string    `gorm:"column:alert_type;not null" json:"alertType"`
	Conditions  string    `gorm:"column:conditions;not null" json:"conditions"` // JSON
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (PriceAlertTemplateEntity) TableName() string {
	return "price_alert_templates"
}

// PriceAlertTriggerHistoryEntity 对应 price_alert_trigger_history 表
type PriceAlertTriggerHistoryEntity struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	AlertID        uint      `gorm:"column:alert_id;not null;index" json:"alertId"`
	StockCode      string    `gorm:"column:stock_code;not null;index" json:"stockCode"`
	StockName      string    `gorm:"column:stock_name;not null" json:"stockName"`
	AlertType      string    `gorm:"column:alert_type;not null" json:"alertType"`
	TriggerPrice   float64   `gorm:"column:trigger_price" json:"triggerPrice"`
	TriggerMessage string    `gorm:"column:trigger_message" json:"triggerMessage"`
	TriggeredAt    time.Time `gorm:"column:triggered_at;default:CURRENT_TIMESTAMP;index" json:"triggeredAt"`
}

func (PriceAlertTriggerHistoryEntity) TableName() string {
	return "price_alert_trigger_history"
}

// StockMoneyFlowHistEntity 对应 stock_money_flow_hist 表
type StockMoneyFlowHistEntity struct {
	Code       string  `gorm:"primaryKey;column:code" json:"code"`
	TradeDate  string  `gorm:"primaryKey;column:trade_date" json:"tradeDate"` // Composite Key part 2
	MainNet    float64 `gorm:"column:main_net;default:0" json:"mainNet"`
	SuperNet   float64 `gorm:"column:super_net;default:0" json:"superNet"`
	BigNet     float64 `gorm:"column:big_net;default:0" json:"bigNet"`
	MidNet     float64 `gorm:"column:mid_net;default:0" json:"midNet"`
	SmallNet   float64 `gorm:"column:small_net;default:0" json:"smallNet"`
	ClosePrice float64 `gorm:"column:close_price;default:0" json:"closePrice"`
	ChgPct     float64 `gorm:"column:chg_pct;default:0" json:"chgPct"`
}

func (StockMoneyFlowHistEntity) TableName() string {
	return "stock_money_flow_hist"
}

// StockStrategySignalEntity 对应 stock_strategy_signals 表
type StockStrategySignalEntity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code         string    `gorm:"column:code;not null;uniqueIndex:idx_signal_unique" json:"code"`
	TradeDate    string    `gorm:"column:trade_date;not null;uniqueIndex:idx_signal_unique" json:"tradeDate"`
	SignalType   string    `gorm:"column:signal_type;not null" json:"signalType"` // 'B' or 'S'
	StrategyName string    `gorm:"column:strategy_name;not null;uniqueIndex:idx_signal_unique" json:"strategyName"`
	Score        float64   `gorm:"column:score;default:0" json:"score"`
	Details      string    `gorm:"column:details" json:"details"` // JSON or description
	AIScore      int       `gorm:"column:ai_score;default:0" json:"aiScore"`
	AIReason     string    `gorm:"column:ai_reason" json:"aiReason"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (StockStrategySignalEntity) TableName() string {
	return "stock_strategy_signals"
}

// KLineEntity 对应 kline_{code} 表 (动态表名)
type KLineEntity struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Date      string    `gorm:"column:date;not null;uniqueIndex" json:"date"`
	Open      float64   `gorm:"column:open;not null" json:"open"`
	High      float64   `gorm:"column:high;not null" json:"high"`
	Low       float64   `gorm:"column:low;not null" json:"low"`
	Close     float64   `gorm:"column:close;not null" json:"close"`
	Volume    int64     `gorm:"column:volume;not null" json:"volume"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}
