package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var ErrInvalidInput = errors.New("输入参数无效")

// StockData 股票数据结构
type StockData struct {
	Code      string  `json:"code"`      // 股票代码
	Name      string  `json:"name"`      // 股票名称
	Price     float64 `json:"price"`     // 最新价
	Change    float64 `json:"change"`    // 涨跌额
	ChangeRate float64 `json:"changeRate"` // 涨跌幅(%)
	Volume    int64   `json:"volume"`    // 成交量(手)
	Amount    float64 `json:"amount"`    // 成交额
	High      float64 `json:"high"`      // 最高价
	Low       float64 `json:"low"`       // 最低价
	Open      float64 `json:"open"`      // 今开
	PreClose  float64 `json:"preClose"`  // 昨收
	Amplitude float64 `json:"amplitude"` // 振幅(%)
	Turnover  float64 `json:"turnover"`  // 换手率(%)
	PE        float64 `json:"pe"`        // 市盈率
	PB        float64 `json:"pb"`        // 市净率
	TotalMV   float64 `json:"totalMV"`   // 总市值
	CircMV    float64 `json:"circMV"`    // 流通市值
	VolumeRatio float64 `json:"volumeRatio"` // 量比
	WarrantRatio float64 `json:"warrantRatio"` // 委比
}

// KLineData K线数据点
type KLineData struct {
	Time   string  `json:"time"`   // 时间 (YYYY-MM-DD)
	Open   float64 `json:"open"`   // 开盘价
	High   float64 `json:"high"`   // 最高价
	Low    float64 `json:"low"`    // 最低价
	Close  float64 `json:"close"`  // 收盘价
	Volume int64   `json:"volume"` // 成交量
	MACD   *MACD   `json:"macd,omitempty"`
	KDJ    *KDJ    `json:"kdj,omitempty"`
	RSI    float64 `json:"rsi,omitempty"`
}

// IntradayData 分时数据点
type IntradayData struct {
	Time      string  `json:"time"`      // 时间 (HH:MM)
	Price     float64 `json:"price"`     // 价格
	AvgPrice  float64 `json:"avgPrice"`  // 均价
	Volume    int64   `json:"volume"`    // 成交量
	PreClose  float64 `json:"preClose"`  // 昨收价
}

// IntradayResponse 分时数据响应结构
type IntradayResponse struct {
	Data      []IntradayData `json:"data"`
	PreClose  float64        `json:"preClose"`
}

// MoneyFlowData 资金流向数据点
type MoneyFlowData struct {
	Time       string  `json:"time"`       // 时间 (HH:MM)
		Main    float64 `json:"main"`    // 主力净流入
		Retail  float64 `json:"retail"`  // 散户净流入
		Super   float64 `json:"super"`   // 超大单净流入
		Big     float64 `json:"big"`     // 大单净流入
		Medium  float64 `json:"medium"`  // 中单净流入
		Small   float64 `json:"small"`   // 小单净流入
}

	// MoneyFlowResponse 资金流向响应结构
	type MoneyFlowResponse struct {
		Data        []MoneyFlowData `json:"data"`
		TodayMain   float64         `json:"todayMain"`   // 今日主力净流入总额
		TodayRetail float64         `json:"todayRetail"` // 今日散户净流入总额
		Status      string          `json:"status"`      // 智能识别状态: "主力建仓", "散户追高", "机构洗盘", "平稳运行"
		Description string          `json:"description"` // 状态详细描述
	}

// HealthCheckResult 股票深度体检结果
type HealthCheckResult struct {
	Score       int            `json:"score"`       // 综合评分 (0-100)
	Status      string         `json:"status"`      // "健康", "亚健康", "风险"
	Items       []HealthItem   `json:"items"`       // 体检项
	Summary     string         `json:"summary"`     // AI 总结
	RiskLevel   string         `json:"riskLevel"`   // 风险等级: "低", "中", "高"
	UpdatedAt   string         `json:"updatedAt"`   // 更新时间
}

// HealthItem 体检子项
type HealthItem struct {
	Category    string `json:"category"`    // 类别: "财务", "资金", "技术", "舆情"
	Name        string `json:"name"`        // 项目名称
	Value       string `json:"value"`       // 检测值
	Status      string `json:"status"`      // "正常", "警告", "异常"
	Description string `json:"description"` // 详细解释
}

type MACD struct {
	DIF float64 `json:"dif"`
	DEA float64 `json:"dea"`
	Bar float64 `json:"bar"`
}

type KDJ struct {
	K float64 `json:"k"`
	D float64 `json:"d"`
	J float64 `json:"j"`
}

// AnalysisReport AI分析报告结构
type AnalysisReport struct {
	StockCode      string `json:"stockCode"`      // 股票代码
	StockName      string `json:"stockName"`      // 股票名称
	Summary        string `json:"summary"`        // 分析摘要
	Fundamentals   string `json:"fundamentals"`   // 基本面分析
	Technical      string `json:"technical"`      // 技术面分析
	Recommendation string `json:"recommendation"` // 投资建议
	RiskLevel      string `json:"riskLevel"`      // 风险等级
	TargetPrice    string `json:"targetPrice"`    // 目标价位
	GeneratedAt    string `json:"generatedAt"`    // 生成时间
}

// TechnicalAnalysisResult 深度技术分析结果（包含绘图数据和风险评估）
type TechnicalAnalysisResult struct {
	Analysis     string            `json:"analysis"`
	Drawings     []TechnicalDrawing `json:"drawings"`
	RiskScore    int               `json:"riskScore"`    // 0-100
	ActionAdvice string            `json:"actionAdvice"` // "买入", "卖出", "观望", "减持", "增持"
	RadarData    *RadarData        `json:"radarData,omitempty"` // 多维度评分雷达图数据
	TradePlan    *TradePlan        `json:"tradePlan,omitempty"` // 智能交易计划
}

// RadarData 多维度评分雷达图数据
type RadarData struct {
	Scores  map[string]int    `json:"scores"`  // 维度名称 -> 分数 (0-100)
	Reasons map[string]string `json:"reasons"` // 维度名称 -> 评分理由
}

// TradePlan 智能交易计划
type TradePlan struct {
	SuggestedPosition string  `json:"suggestedPosition"` // 建议仓位 (如 "30%")
	StopLoss          float64 `json:"stopLoss"`          // 止损价
	TakeProfit        float64 `json:"takeProfit"`        // 止盈价
	RiskRewardRatio   float64 `json:"riskRewardRatio"`   // 盈亏比
	Strategy          string  `json:"strategy"`          // 操作策略描述
}

// TechnicalDrawing AI识别的绘图数据
type TechnicalDrawing struct {
	Type       string  `json:"type"`       // "support", "resistance", "trendline"
	Price      float64 `json:"price"`      // 用于支撑/阻力位
	Start      string  `json:"start"`      // 用于趋势线起点时间
	End        string  `json:"end"`        // 用于趋势线终点时间
	StartPrice float64 `json:"startPrice"` // 趋势线起点价格
	EndPrice   float64 `json:"endPrice"`   // 趋势线终点价格
	Label      string  `json:"label"`      // 标签
}

// PriceAlert 价格预警配置
type PriceAlert struct {
	StockCode     string    `json:"stockCode"`
	StockName     string    `json:"stockName"`
	Type          string    `json:"type"`      // support (支撑), resistance (压力)
	Price         float64   `json:"price"`     // 触发价格
	Label         string    `json:"label"`     // 标签描述
	Role          string    `json:"role"`      // 触发时的AI角色
	IsActive      bool      `json:"isActive"`  // 是否激活
	LastTriggered time.Time `json:"-"`         // 上次触发时间，避免频繁骚扰
}

// AlertConfig 预警系统配置
type AlertConfig struct {
	Sensitivity float64 `json:"sensitivity"` // 灵敏度 (0.001 - 0.02)
	Cooldown    int     `json:"cooldown"`    // 冷却时间 (小时)
	Enabled     bool    `json:"enabled"`     // 是否开启全局预警
}

// AlertHistory 告警历史记录
type AlertHistory struct {
	ID        string    `json:"id"`
	StockCode string    `json:"stockCode"`
	StockName string    `json:"stockName"`
	Type      string    `json:"type"`      // support, resistance
	Price     float64   `json:"price"`     // 触发时的价格
	AlertPrice float64  `json:"alertPrice"` // 预警设定的价格
	Message   string    `json:"message"`
	Time      time.Time `json:"time"`
}

// EntryStrategyResult 智能建仓方案结构
type EntryStrategyResult struct {
		Recommendation    string        `json:"recommendation"`    // 总体建议
		EntryPriceRange   string        `json:"entryPriceRange"`   // 建议买入价格区间
		InitialPosition   string        `json:"initialPosition"`   // 建议首仓比例
		StopLossPrice     float64       `json:"stopLossPrice"`     // 止损价
		TakeProfitPrice   float64       `json:"takeProfitPrice"`   // 目标止盈价
		CoreReasons       []CoreReason  `json:"coreReasons"`       // 核心建仓理由
		RiskRewardRatio   float64       `json:"riskRewardRatio"`   // 预估盈亏比
		ActionPlan        string        `json:"actionPlan"`        // 具体操作步骤
		TrailingStopConfig *TrailingStopConfig `json:"trailingStopConfig,omitempty"` // 移动止损配置
	}

func (r *EntryStrategyResult) UnmarshalJSON(data []byte) error {
	type Alias EntryStrategyResult
	aux := struct {
		EntryPriceRange json.RawMessage `json:"entryPriceRange"`
		InitialPosition json.RawMessage `json:"initialPosition"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// entryPriceRange: 兼容 string / [number, number] / number
	if len(aux.EntryPriceRange) > 0 {
		var s string
		if err := json.Unmarshal(aux.EntryPriceRange, &s); err == nil {
			r.EntryPriceRange = s
		} else {
			var arr []float64
			if err := json.Unmarshal(aux.EntryPriceRange, &arr); err == nil && len(arr) > 0 {
				if len(arr) == 1 {
					r.EntryPriceRange = fmt.Sprintf("%.2f", arr[0])
				} else {
					r.EntryPriceRange = fmt.Sprintf("%.2f-%.2f", arr[0], arr[1])
				}
			} else {
				var n float64
				if err := json.Unmarshal(aux.EntryPriceRange, &n); err == nil {
					r.EntryPriceRange = fmt.Sprintf("%.2f", n)
				}
			}
		}
	}

	// initialPosition: 兼容 string / number（number 统一转成百分比字符串）
	if len(aux.InitialPosition) > 0 {
		var s string
		if err := json.Unmarshal(aux.InitialPosition, &s); err == nil {
			r.InitialPosition = s
		} else {
			var n float64
			if err := json.Unmarshal(aux.InitialPosition, &n); err == nil {
				r.InitialPosition = fmt.Sprintf("%.0f%%", n)
			}
		}
	}

	return nil
}

// CoreReason 核心建仓理由
type CoreReason struct {
	Type        string `json:"type"`        // fundamental, technical, money_flow
	Description string `json:"description"`
	Threshold   string `json:"threshold"`   // 逻辑失效的触发阈值
}

// TrailingStopConfig 移动止损个性化配置
type TrailingStopConfig struct {
	Enabled             bool    `json:"enabled"`             // 是否启用
	ActivationThreshold float64 `json:"activationThreshold"` // 触发阈值 (如 0.05 代表盈利 5% 启动)
	CallbackRate        float64 `json:"callbackRate"`        // 跟踪回撤比例 (如 0.03 代表回撤 3% 止盈)
}

// Position 持仓记录（用于逻辑跟踪）
type Position struct {
	StockCode      string              `json:"stockCode"`
	StockName      string              `json:"stockName"`
	EntryPrice     float64             `json:"entryPrice"`
	EntryTime      time.Time           `json:"entryTime"`
	Strategy       EntryStrategyResult `json:"strategy"`
	TrailingConfig TrailingStopConfig  `json:"trailingConfig"` // 移动止损配置
	CurrentStatus  string              `json:"currentStatus"`  // "holding", "closed"
	LogicStatus    string              `json:"logicStatus"`    // "valid", "violated", "warning"
	UpdatedAt      time.Time           `json:"updatedAt"`
}

// EastMoneyResponse 东方财富API响应结构
type EastMoneyResponse struct {
	RC   int    `json:"rc"`
	RT   int    `json:"rt"`
	SVRT int    `json:"svrt"`
	LT   int    `json:"lt"`
	Full int    `json:"full"`
	Data struct {
		Total int           `json:"total"`
		Diff  []StockDiff   `json:"diff"`
	} `json:"data"`
}

// StockDiff 东方财富股票数据差异结构
type StockDiff struct {
	F1  int     `json:"f1"`  // 未知
	F2  float64 `json:"f2"`  // 最新价
	F3  float64 `json:"f3"`  // 涨跌幅
	F4  float64 `json:"f4"`  // 涨跌额
	F5  int64   `json:"f5"`  // 成交量(手)
	F6  float64 `json:"f6"`  // 成交额
	F7  float64 `json:"f7"`  // 振幅
	F8  float64 `json:"f8"`  // 换手率
	F9  float64 `json:"f9"`  // 市盈率
	F10 float64 `json:"f10"` // 量比
	F11 float64 `json:"f11"` // 未知
	F12 string  `json:"f12"` // 股票代码
	F13 int     `json:"f13"` // 市场编号
	F14 string  `json:"f14"` // 股票名称
	F15 float64 `json:"f15"` // 最高价
	F16 float64 `json:"f16"` // 最低价
	F17 float64 `json:"f17"` // 今开
	F18 float64 `json:"f18"` // 昨收
	F20 float64 `json:"f20"` // 总市值
	F21 float64 `json:"f21"` // 流通市值
	F22 float64 `json:"f22"` // 未知
	F23 float64 `json:"f23"` // 市净率
}

// ToStockData 将东方财富数据转换为标准股票数据
func (sd *StockDiff) ToStockData() *StockData {
	return &StockData{
		Code:       sd.F12,
		Name:       sd.F14,
		Price:      sd.F2,
		Change:     sd.F4,
		ChangeRate: sd.F3,
		Volume:     sd.F5,
		Amount:     sd.F6,
		High:       sd.F15,
		Low:        sd.F16,
		Open:       sd.F17,
		PreClose:   sd.F18,
		Amplitude:  sd.F7,
		Turnover:   sd.F8,
		PE:         sd.F9,
		PB:         sd.F23,
		TotalMV:    sd.F20,
		CircMV:     sd.F21,
	}
}


// TradeRecord 单笔交易记录
type TradeRecord struct {
	Time      string  `json:"time"`      // 交易时间
	Type      string  `json:"type"`      // 交易类型: "BUY" 或 "SELL"
	Price     float64 `json:"price"`     // 交易价格
	Volume    int64   `json:"volume"`    // 交易数量
	Amount    float64 `json:"amount"`    // 交易金额
	Commission float64 `json:"commission"` // 佣金
	Tax       float64 `json:"tax"`       // 印花税 (仅卖出)
	Profit    float64 `json:"profit"`    // 单笔交易盈亏
}

// BacktestResult 回测结果结构
type BacktestResult struct {
	StrategyName    string        `json:"strategyName"`    // 策略名称
	StockCode       string        `json:"stockCode"`       // 股票代码
	StartDate       string        `json:"startDate"`       // 回测开始日期
	EndDate         string        `json:"endDate"`         // 回测结束日期
	InitialCapital  float64       `json:"initialCapital"`  // 初始资金
	FinalCapital    float64       `json:"finalCapital"`    // 最终资金
	TotalReturn     float64       `json:"totalReturn"`     // 总收益率
	AnnualizedReturn float64       `json:"annualizedReturn"` // 年化收益率
	MaxDrawdown     float64       `json:"maxDrawdown"`     // 最大回撤
	WinRate         float64       `json:"winRate"`         // 胜率
	TradeCount      int           `json:"tradeCount"`      // 交易次数
	Trades          []TradeRecord `json:"trades"`          // 交易记录
	EquityCurve     []float64     `json:"equityCurve"`     // 净值曲线 (每日资产总值)
	EquityDates     []string      `json:"equityDates"`     // 净值曲线对应的日期
}


// KLineCacheRecord 用于存储到 SQLite 的 K 线缓存记录
type KLineCacheRecord struct {
	ID        int64   `db:"id" json:"id"`
	StockCode string  `db:"stock_code" json:"stock_code"`
	Date      string  `db:"date" json:"date"` // YYYY-MM-DD
	Open      float64 `db:"open" json:"open"`
	High      float64 `db:"high" json:"high"`
	Low       float64 `db:"low" json:"low"`
	Close     float64 `db:"close" json:"close"`
	Volume    int64   `db:"volume" json:"volume"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt string  `db:"updated_at" json:"updated_at"`
}

// SyncProgress 数据同步进度信息
type SyncProgress struct {
	StockCode      string `json:"stock_code"`
	Status         string `json:"status"` // "pending", "syncing", "completed", "failed"
	TotalRecords   int    `json:"total_records"`
	SyncedRecords  int    `json:"synced_records"`
	Message        string `json:"message"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	ErrorMessage   string `json:"error_message,omitempty"`
}

// SyncResult 数据同步结果
type SyncResult struct {
	StockCode     string `json:"stock_code"`
	Success       bool   `json:"success"`
	RecordsAdded  int    `json:"records_added"`
	RecordsUpdated int   `json:"records_updated"`
	Message       string `json:"message"`
	ErrorMessage  string `json:"error_message,omitempty"`
}

// DataSyncStats 数据同步统计信息
type DataSyncStats struct {
	TotalStocks    int      `json:"total_stocks"`
	SyncedStocks   int      `json:"synced_stocks"`
	TotalRecords   int64    `json:"total_records"`
	StockList      []string `json:"stock_list"`
	LastSyncTime   string   `json:"last_sync_time"`
}
