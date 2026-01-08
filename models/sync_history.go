package models

import "time"

// SyncHistory 同步历史记录
type SyncHistory struct {
	ID             int       `json:"id"`
	StockCode      string    `json:"stockCode"`      // 股票代码
	StockName      string    `json:"stockName"`      // 股票名称
	SyncType       string    `json:"syncType"`       // 同步类型: "single" 单个同步, "batch" 批量同步
	StartDate      string    `json:"startDate"`      // 开始日期
	EndDate        string    `json:"endDate"`        // 结束日期
	Status         string    `json:"status"`         // 状态: "success" 成功, "failed" 失败
	RecordsAdded   int       `json:"recordsAdded"`   // 新增记录数
	RecordsUpdated int       `json:"recordsUpdated"` // 更新记录数
	Duration       int       `json:"duration"`       // 耗时（秒）
	ErrorMsg       string    `json:"errorMsg"`       // 错误信息
	CreatedAt      time.Time `json:"createdAt"`      // 创建时间
}
