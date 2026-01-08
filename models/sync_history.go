package models

import "time"

// SyncHistory 同步历史记录
type SyncHistory struct {
	ID             int       `json:"id"`
	StockCode      string    `json:"stock_code"`      // 股票代码
	StockName      string    `json:"stock_name"`      // 股票名称
	SyncType       string    `json:"sync_type"`       // 同步类型: "single" 单个同步, "batch" 批量同步
	StartDate      string    `json:"start_date"`      // 开始日期
	EndDate        string    `json:"end_date"`        // 结束日期
	Status         string    `json:"status"`         // 状态: "success" 成功, "failed" 失败
	RecordsAdded   int       `json:"records_added"`   // 新增记录数
	RecordsUpdated int       `json:"records_updated"` // 更新记录数
	Duration       int       `json:"duration"`       // 耗时（秒）
	ErrorMsg       string    `json:"error_msg"`       // 错误信息
	CreatedAt      time.Time `json:"created_at"`      // 创建时间
}
