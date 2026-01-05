package models

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
}

// RadarData 多维度评分雷达图数据
type RadarData struct {
	Scores  map[string]int    `json:"scores"`  // 维度名称 -> 分数 (0-100)
	Reasons map[string]string `json:"reasons"` // 维度名称 -> 评分理由
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
