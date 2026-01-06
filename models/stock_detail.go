package models

import "time"

// OrderBookEntry 盘口数据项
type OrderBookEntry struct {
	Price float64 `json:"price"` // 价格
	Volume int64 `json:"volume"` // 数量 (手)
}

// OrderBook 实时五档盘口数据
type OrderBook struct {
	Buy5  []OrderBookEntry `json:"buy5"`  // 买五档
	Sell5 []OrderBookEntry `json:"sell5"` // 卖五档
}

// FinancialSummary 核心财务摘要
type FinancialSummary struct {
	// 净资产收益率 (ROE)
	ROE float64 `json:"roe"`
	// 净利润增长率 (%)
	NetProfitGrowthRate float64 `json:"net_profit_growth_rate"`
	// 毛利率 (%)
	GrossProfitMargin float64 `json:"gross_profit_margin"`
	// 总市值 (亿元)
	TotalMarketValue float64 `json:"total_market_value"`
	// 流通市值 (亿元)
	CirculatingMarketValue float64 `json:"circulating_market_value"`
	// 股息率 (%)
	DividendYield float64 `json:"dividend_yield"`
	// 报告期
	ReportDate time.Time `json:"report_date"`
}

// IndustryInfo 行业与宏观信息
type IndustryInfo struct {
	// 行业名称
	IndustryName string `json:"industry_name"`
	// 概念板块名称列表
	ConceptNames []string `json:"concept_names"`
	// 行业平均市盈率 (P/E)
	IndustryPE float64 `json:"industry_pe"`
}

// StockDetail 个股详情页聚合数据
type StockDetail struct {
	StockData // 继承基础行情数据

	OrderBook OrderBook `json:"orderBook"` // 实时盘口数据
	Financial FinancialSummary `json:"financial_summary"` // 核心财务数据
	Industry IndustryInfo `json:"industry_info"` // 行业与宏观信息
}
