package models

import "time"

// AnalysisCacheEntry 单条分析缓存记录
type AnalysisCacheEntry struct {
	StockCode string                  `json:"stockCode"`
	Role      string                  `json:"role"`
	Period    string                  `json:"period"`
	Result    TechnicalAnalysisResult `json:"result"`
	Timestamp time.Time               `json:"timestamp"`
}

// AnalysisCache 缓存文件结构
type AnalysisCache struct {
	Entries map[string]AnalysisCacheEntry `json:"entries"`
}
