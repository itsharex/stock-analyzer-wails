package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"stock-analyzer-wails/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"go.uber.org/zap"
)

// KLineSyncService K线数据同步服务
type KLineSyncService struct {
	dbService *DBService
	ctx       context.Context
	running   bool
	mu        sync.Mutex
}

// NewKLineSyncService 创建K线同步服务
func NewKLineSyncService(dbService *DBService) *KLineSyncService {
	return &KLineSyncService{
		dbService: dbService,
		running:   false,
	}
}

// SetContext 设置上下文（用于发送事件）
func (s *KLineSyncService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// KLineSyncProgress K线同步进度
type KLineSyncProgress struct {
	IsRunning       bool    `json:"isRunning"`       // 是否正在运行
	CurrentIndex    int     `json:"currentIndex"`    // 当前处理索引
	TotalCount      int     `json:"totalCount"`      // 总数
	CurrentCode     string  `json:"currentCode"`     // 当前股票代码
	CurrentName     string  `json:"currentName"`     // 当前股票名称
	SuccessCount    int     `json:"successCount"`    // 成功数量
	FailedCount     int     `json:"failedCount"`     // 失败数量
	TotalRecords    int     `json:"totalRecords"`    // 总记录数
	RecordsPerSec   float64 `json:"recordsPerSec"`   // 每秒记录数
	StartTime       string  `json:"startTime"`       // 开始时间
	ElapsedSeconds  int     `json:"elapsedSeconds"`  // 已用时间（秒）
	EstimatedSeconds int    `json:"estimatedSeconds"` // 预计剩余时间（秒）
}

// KLineSyncResult K线同步结果
type KLineSyncResult struct {
	Success       bool   `json:"success"`
	TotalCount    int    `json:"totalCount"`
	SuccessCount  int    `json:"successCount"`
	FailedCount   int    `json:"failedCount"`
	TotalRecords  int    `json:"totalRecords"`
	Duration      int    `json:"duration"`      // 耗时（秒）
	Message       string `json:"message"`
}

// KLineSyncTask K线同步任务
type KLineSyncTask struct {
	Code   string
	Name   string
	Market string
}

// StartKLineSync 开始K线数据同步
func (s *KLineSyncService) StartKLineSync(days int) (*KLineSyncResult, error) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return &KLineSyncResult{
			Success: false,
			Message: "同步任务已在运行中",
		}, fmt.Errorf("同步任务已在运行中")
	}
	s.running = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
	}()

	// 验证日期范围
	if days <= 0 || days > 2000 {
		return &KLineSyncResult{
			Success: false,
			Message: "日期范围无效（1-2000天）",
		}, fmt.Errorf("日期范围无效")
	}

	startTime := time.Now()

	// 1. 获取所有活跃股票
	tasks, err := s.getActiveStocks()
	if err != nil {
		logger.Error("获取活跃股票失败", zap.Error(err))
		return &KLineSyncResult{
			Success: false,
			Message: fmt.Sprintf("获取活跃股票失败: %v", err),
		}, err
	}

	if len(tasks) == 0 {
		return &KLineSyncResult{
			Success: false,
			Message: "没有需要同步的股票",
		}, fmt.Errorf("没有需要同步的股票")
	}

	logger.Info("开始K线数据同步",
		zap.Int("stock_count", len(tasks)),
		zap.Int("days", days),
	)

	// 2. 顺序同步每只股票，避免SQLite并发锁库问题
	// 统计结果
	var successCount, failedCount, totalRecords int

	// 初始化进度
	progress := &KLineSyncProgress{
		IsRunning:     true,
		TotalCount:    len(tasks),
		StartTime:     startTime.Format("2006-01-02 15:04:05"),
	}

	// 顺序处理每只股票
	for i, task := range tasks {
		// 随机延迟模拟真人行为（200-500ms）
		// 延迟可以防止被反爬虫机制识别
		delay := time.Duration(rand.Intn(300)+200) * time.Millisecond
		time.Sleep(delay)

		// 获取K线数据
		klines, err := s.fetchKLineData(task, days)
		if err != nil {
			failedCount++
			logger.Error("获取K线数据失败",
				zap.String("code", task.Code),
				zap.Error(err),
			)

			// 更新进度
			s.updateProgress(progress, i+1, len(tasks), task.Code, task.Name, successCount, failedCount, totalRecords, startTime)

			// 记录同步历史
			if recordErr := s.recordSyncHistory(task, days, 0, 0, 0, false, err.Error()); recordErr != nil {
				logger.Error("记录同步历史失败",
					zap.String("code", task.Code),
					zap.Error(recordErr),
				)
			}
			continue
		}

		// 存储K线数据
		added, updated, err := s.saveKLineData(task.Code, klines)
		if err != nil {
			failedCount++
			logger.Error("保存K线数据失败",
				zap.String("code", task.Code),
				zap.Error(err),
			)

			// 更新进度
			s.updateProgress(progress, i+1, len(tasks), task.Code, task.Name, successCount, failedCount, totalRecords, startTime)

			// 记录同步历史
			if recordErr := s.recordSyncHistory(task, days, len(klines), 0, 0, false, err.Error()); recordErr != nil {
				logger.Error("记录同步历史失败",
					zap.String("code", task.Code),
					zap.Error(recordErr),
				)
			}
			continue
		}

		successCount++
		totalRecords += int(added + updated)

		// 记录同步历史
		if err := s.recordSyncHistory(task, days, len(klines), int(added), int(updated), true, ""); err != nil {
			logger.Error("记录同步历史失败",
				zap.String("code", task.Code),
				zap.Error(err),
			)
		}

		logger.Debug("K线数据同步成功",
			zap.String("code", task.Code),
			zap.Int("records", len(klines)),
			zap.Int("added", int(added)),
			zap.Int("updated", int(updated)),
		)

		// 更新进度
		s.updateProgress(progress, i+1, len(tasks), task.Code, task.Name, successCount, failedCount, totalRecords, startTime)
	}

	// 发送最终进度
	progress.IsRunning = false
	s.emitProgress(progress)

	duration := int(time.Since(startTime).Seconds())

	logger.Info("K线数据同步完成",
		zap.Int("total_count", len(tasks)),
		zap.Int("success_count", successCount),
		zap.Int("failed_count", failedCount),
		zap.Int("total_records", totalRecords),
		zap.Int("duration", duration),
	)

	return &KLineSyncResult{
		Success:      true,
		TotalCount:   len(tasks),
		SuccessCount: successCount,
		FailedCount:  failedCount,
		TotalRecords: totalRecords,
		Duration:     duration,
		Message:      fmt.Sprintf("同步完成：成功 %d 只，失败 %d 只，总记录数 %d 条", successCount, failedCount, totalRecords),
	}, nil
}

// getActiveStocks 获取所有活跃股票
func (s *KLineSyncService) getActiveStocks() ([]*KLineSyncTask, error) {
	db := s.dbService.GetDB()

	rows, err := db.Query(`
		SELECT code, name, market FROM stocks WHERE is_active = 1 ORDER BY code
	`)
	if err != nil {
		return nil, fmt.Errorf("查询活跃股票失败: %w", err)
	}
	defer rows.Close()

	var tasks []*KLineSyncTask
	for rows.Next() {
		var code, name, market string
		if err := rows.Scan(&code, &name, &market); err != nil {
			logger.Error("扫描股票记录失败", zap.Error(err))
			continue
		}
		tasks = append(tasks, &KLineSyncTask{
			Code:   code,
			Name:   name,
			Market: market,
		})
	}

	return tasks, nil
}

// fetchKLineData 获取K线数据
func (s *KLineSyncService) fetchKLineData(task *KLineSyncTask, days int) ([]map[string]interface{}, error) {
	// 构造secid
	secid := ""
	if task.Market == "SH" {
		secid = "1." + task.Code
	} else if task.Market == "SZ" {
		secid = "0." + task.Code
	} else if task.Market == "BJ" {
		secid = "2." + task.Code
	} else {
		// 默认根据代码判断
		if task.Code[0] == '6' {
			secid = "1." + task.Code
		} else if task.Code[0] == '0' || task.Code[0] == '3' {
			secid = "0." + task.Code
		} else {
			secid = "2." + task.Code
		}
	}

	// 计算开始日期
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// 构造请求URL
	// klt=101: 日K
	// fqt=1: 前复权
	url := fmt.Sprintf(
		"https://push2his.eastmoney.com/api/qt/stock/kline/get?secid=%s&fields1=f1,f2,f3,f4,f5,f6&fields2=f51,f52,f53,f54,f55,f56&klt=101&fqt=1&end=%s&lmt=%d",
		secid,
		endDate.Format("20060102"),
		days,
	)

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求K线API失败: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result struct {
		RC   int `json:"rc"`
		RT   int `json:"rt"`
		Data struct {
			Klines []string `json:"klines"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析K线响应失败: %w", err)
	}

	if result.RC != 0 {
		return nil, fmt.Errorf("API返回错误: rc=%d", result.RC)
	}

	// 转换K线数据格式
	klines := make([]map[string]interface{}, 0)
	for _, line := range result.Data.Klines {
		parts := strings.Split(line, ",")
		if len(parts) < 6 {
			continue
		}

		date := parts[0]
		// 过滤日期范围
		klineDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			continue
		}

		if klineDate.Before(startDate) || klineDate.After(endDate) {
			continue
		}

		klines = append(klines, map[string]interface{}{
			"date":   date,
			"open":   parsePrice(parts[1]),
			"high":   parsePrice(parts[2]),
			"low":    parsePrice(parts[3]),
			"close":  parsePrice(parts[4]),
			"volume": int64(parsePrice(parts[5])),
		})
	}

	return klines, nil
}

// parsePrice 解析价格
func parsePrice(s string) float64 {
	if s == "" || s == "-" {
		return 0
	}
	price, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return price
}

// saveKLineData 保存K线数据
func (s *KLineSyncService) saveKLineData(code string, klines []map[string]interface{}) (int64, int64, error) {
	return s.dbService.InsertOrUpdateKLineData(code, klines)
}

// recordSyncHistory 记录同步历史
func (s *KLineSyncService) recordSyncHistory(task *KLineSyncTask, days int, totalRecords, added, updated int, success bool, errorMsg string) error {
	db := s.dbService.GetDB()

	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	status := "success"
	if !success {
		status = "failed"
	}

	_, err := db.Exec(`
		INSERT INTO sync_history (
			stock_code, stock_name, sync_type, start_date, end_date, status,
			records_added, records_updated, duration, error_msg
		) VALUES (?, ?, 'kline', ?, ?, ?, ?, ?, 0, ?)
	`, task.Code, task.Name, startDate, endDate, status, added, updated, errorMsg)

	return err
}

// updateProgress 更新进度信息
func (s *KLineSyncService) updateProgress(progress *KLineSyncProgress, currentIndex, totalCount int, code, name string, successCount, failedCount, totalRecords int, startTime time.Time) {
	// 计算速率
	elapsed := time.Since(startTime).Seconds()
	recordsPerSec := float64(totalRecords) / elapsed
	estimatedSeconds := 0
	if currentIndex > 0 && elapsed > 0 {
		estimatedSeconds = int((elapsed / float64(currentIndex)) * float64(totalCount-currentIndex))
	}

	// 更新进度
	progress.CurrentIndex = currentIndex
	progress.CurrentCode = code
	progress.CurrentName = name
	progress.SuccessCount = successCount
	progress.FailedCount = failedCount
	progress.TotalRecords = totalRecords
	progress.RecordsPerSec = recordsPerSec
	progress.ElapsedSeconds = int(elapsed)
	progress.EstimatedSeconds = estimatedSeconds

	// 发送进度事件
	s.emitProgress(progress)
}

// emitProgress 发送进度事件
func (s *KLineSyncService) emitProgress(progress *KLineSyncProgress) {
	if s.ctx == nil {
		return
	}

	runtime.EventsEmit(s.ctx, "klineSyncProgress", progress)
}

// GetSyncProgress 获取当前同步进度
func (s *KLineSyncService) GetSyncProgress() (*KLineSyncProgress, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return &KLineSyncProgress{
			IsRunning: false,
		}, nil
	}

	// 这里返回默认的进度信息
	// 实际的进度信息通过事件发送到前端
	return &KLineSyncProgress{
		IsRunning: true,
	}, nil
}

// GetKLineSyncHistory 获取K线同步历史记录
func (s *KLineSyncService) GetKLineSyncHistory(limit int) ([]map[string]interface{}, error) {
	db := s.dbService.GetDB()

	rows, err := db.Query(`
		SELECT id, stock_code, stock_name, sync_type, start_date, end_date,
		       status, records_added, records_updated, duration, error_msg, created_at
		FROM sync_history
		WHERE sync_type = 'kline'
		ORDER BY created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("查询K线同步历史失败: %w", err)
	}
	defer rows.Close()

	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int
		var stockCode, stockName, syncType, startDate, endDate, status string
		var recordsAdded, recordsUpdated, duration int
		var errorMsg sql.NullString
		var createdAt string

		if err := rows.Scan(&id, &stockCode, &stockName, &syncType, &startDate, &endDate,
			&status, &recordsAdded, &recordsUpdated, &duration, &errorMsg, &createdAt); err != nil {
			logger.Error("扫描同步历史记录失败", zap.Error(err))
			continue
		}

		record := map[string]interface{}{
			"id":             id,
			"stockCode":      stockCode,
			"stockName":      stockName,
			"syncType":       syncType,
			"startDate":      startDate,
			"endDate":        endDate,
			"status":         status,
			"recordsAdded":   recordsAdded,
			"recordsUpdated": recordsUpdated,
			"duration":       duration,
			"errorMsg":       errorMsg.String,
			"createdAt":      createdAt,
		}
		records = append(records, record)
	}

	return records, nil
}
