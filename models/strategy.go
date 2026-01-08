package models

// StrategyConfig 策略配置模型
type StrategyConfig struct {
	ID                  int64                  `json:"id"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	StrategyType        string                 `json:"strategyType"`
	Parameters          map[string]interface{} `json:"parameters"`
	LastBacktestResult  map[string]interface{} `json:"lastBacktestResult,omitempty"`
	CreatedAt           string                 `json:"createdAt"`
	UpdatedAt           string                 `json:"updatedAt"`
}

// StrategyParameter 策略参数定义
type StrategyParameter struct {
	Name        string      `json:"name"`
	Label       string      `json:"label"`
	Type        string      `json:"type"` // "number", "select", "boolean"
	MinValue    *float64    `json:"minValue,omitempty"`
	MaxValue    *float64    `json:"maxValue,omitempty"`
	DefaultValue interface{} `json:"defaultValue"`
	Options     []string    `json:"options,omitempty"` // 用于 select 类型
}

// StrategyTypeDefinition 策略类型定义
type StrategyTypeDefinition struct {
	Type       string              `json:"type"`
	Name       string              `json:"name"`
	Parameters []StrategyParameter `json:"parameters"`
}

// 内置策略类型定义
var StrategyTypes = []StrategyTypeDefinition{
	{
		Type: "simple_ma",
		Name: "双均线策略",
		Parameters: []StrategyParameter{
			{
				Name:         "shortPeriod",
				Label:        "短周期均线",
				Type:         "number",
				MinValue:     float64Ptr(1),
				MaxValue:     float64Ptr(100),
				DefaultValue: 5,
			},
			{
				Name:         "longPeriod",
				Label:        "长周期均线",
				Type:         "number",
				MinValue:     float64Ptr(1),
				MaxValue:     float64Ptr(500),
				DefaultValue: 20,
			},
			{
				Name:         "initialCapital",
				Label:        "初始资金",
				Type:         "number",
				MinValue:     float64Ptr(1000),
				MaxValue:     float64Ptr(10000000),
				DefaultValue: 100000,
			},
		},
	},
	{
		Type: "macd",
		Name: "MACD 策略",
		Parameters: []StrategyParameter{
			{
				Name:         "fastPeriod",
				Label:        "快线周期",
				Type:         "number",
				MinValue:     float64Ptr(1),
				MaxValue:     float64Ptr(50),
				DefaultValue: 12,
			},
			{
				Name:         "slowPeriod",
				Label:        "慢线周期",
				Type:         "number",
				MinValue:     float64Ptr(1),
				MaxValue:     float64Ptr(100),
				DefaultValue: 26,
			},
			{
				Name:         "signalPeriod",
				Label:        "信号线周期",
				Type:         "number",
				MinValue:     float64Ptr(1),
				MaxValue:     float64Ptr(50),
				DefaultValue: 9,
			},
			{
				Name:         "initialCapital",
				Label:        "初始资金",
				Type:         "number",
				MinValue:     float64Ptr(1000),
				MaxValue:     float64Ptr(10000000),
				DefaultValue: 100000,
			},
		},
	},
}

func float64Ptr(v float64) *float64 {
	return &v
}
