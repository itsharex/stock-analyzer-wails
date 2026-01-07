package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"stock-analyzer-wails/internal/logger"

	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

// DBService 数据库服务
type DBService struct {
	db *sql.DB
}

// NewDBService 初始化数据库连接并创建表
func NewDBService() (*DBService, error) {
	dbPath := filepath.Join(GetAppDataDir(), "stock_analyzer.db")

	// 确保目录存在（GetAppDataDir 通常已创建，但这里做一次兜底）
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %w", err)
	}

	// 检查数据库文件是否存在，如果不存在则创建
	_, err := os.Stat(dbPath)
	isNewDB := os.IsNotExist(err)

	// 连接数据库
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(1) // SQLite 建议最大连接数为 1
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// 立即验证连接可用性（避免延迟到首次 Query/Exec 才报错）
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("数据库连接不可用: %w", err)
	}

	svc := &DBService{db: db}

	if isNewDB {
		logger.Info("数据库文件不存在，开始初始化表结构", zap.String("path", dbPath))
		if err := svc.initTables(); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("初始化数据库表失败: %w", err)
		}
		logger.Info("数据库表结构初始化完成")
	} else {
		logger.Info("数据库连接成功", zap.String("path", dbPath))
	}

	return svc, nil
}

// initTables 初始化数据库表结构
func (s *DBService) initTables() error {
	// 开启事务
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()

	// 1. Watchlist 表
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS watchlist (
			code TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			data TEXT NOT NULL, -- 存储 StockData 的 JSON 字符串
			added_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("创建 watchlist 表失败: %w", err)
	}

	// 2. Alerts 表
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			stock_code TEXT NOT NULL,
			stock_name TEXT NOT NULL,
			price REAL NOT NULL,
			type TEXT NOT NULL, -- 'above' or 'below'
			is_active BOOLEAN NOT NULL,
			last_triggered DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(stock_code, price, type)
		);
	`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("创建 alerts 表失败: %w", err)
	}

	// 3. Alert History 表
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS alert_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			stock_code TEXT NOT NULL,
			stock_name TEXT NOT NULL,
			triggered_price REAL NOT NULL,
			message TEXT NOT NULL,
			triggered_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("创建 alert_history 表失败: %w", err)
	}

	// 4. Positions 表
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS positions (
			stock_code TEXT PRIMARY KEY,
			stock_name TEXT NOT NULL,
			entry_price REAL NOT NULL,
			entry_time DATETIME NOT NULL,
			current_status TEXT NOT NULL, -- 'holding', 'closed'
			logic_status TEXT NOT NULL, -- 'valid', 'violated'
			strategy_json TEXT NOT NULL, -- 存储 EntryStrategyResult 的 JSON 字符串
			trailing_config_json TEXT NOT NULL, -- 存储 TrailingStopConfig 的 JSON 字符串
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("创建 positions 表失败: %w", err)
	}

	// 5. Config 表 (用于存储全局配置，如 AI 配置、Alert 配置等)
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS config (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);
	`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("创建 config 表失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return err
	}

	// 插入默认配置
	return s.insertDefaultConfigs()
}

// GetDB 返回数据库连接对象
func (s *DBService) GetDB() *sql.DB {
	return s.db
}

// Close 关闭数据库连接
func (s *DBService) Close() {
	if s.db != nil {
		_ = s.db.Close()
		logger.Info("数据库连接已关闭")
	}
}

// insertDefaultConfigs 插入默认配置项
func (s *DBService) insertDefaultConfigs() error {
	// 默认配置值
	defaults := map[string]string{
		"trailing_stop_default_activation": "0.05", // 默认盈利 5% 启动
		"trailing_stop_default_callback":   "0.03", // 默认回撤 3% 止盈
	}

	for key, value := range defaults {
		_, err := s.db.Exec(`
			INSERT OR IGNORE INTO config (key, value) VALUES (?, ?)
		`, key, value)
		if err != nil {
			return fmt.Errorf("插入默认配置 (%s) 失败: %w", key, err)
		}
	}
	logger.Info("默认配置项插入完成")
	return nil
}


// CreateKLineCacheTable 为指定股票创建 K 线缓存表
// 表名格式：kline_{code}，例如 kline_600519
func (s *DBService) CreateKLineCacheTable(code string) error {
	tableName := fmt.Sprintf("kline_%s", code)
	
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL UNIQUE,
			open REAL NOT NULL,
			high REAL NOT NULL,
			low REAL NOT NULL,
			close REAL NOT NULL,
			volume INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`, tableName)
	
	_, err := s.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("创建 K 线缓存表 %s 失败: %w", tableName, err)
	}
	
	// 创建日期索引以加速查询
	indexSQL := fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS idx_%s_date ON %s(date);
	`, tableName, tableName)
	
	_, err = s.db.Exec(indexSQL)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}
	
	return nil
}

// InsertOrUpdateKLineData 批量插入或更新 K 线数据
func (s *DBService) InsertOrUpdateKLineData(code string, klines []map[string]interface{}) (int64, int64, error) {
	tableName := fmt.Sprintf("kline_%s", code)
	
	// 先确保表存在
	if err := s.CreateKLineCacheTable(code); err != nil {
		return 0, 0, err
	}
	
	var addedCount int64
	var updatedCount int64
	
	for _, kline := range klines {
		date := kline["date"].(string)
		open := kline["open"].(float64)
		high := kline["high"].(float64)
		low := kline["low"].(float64)
		close := kline["close"].(float64)
		volume := kline["volume"].(int64)
		
		// 尝试插入，如果日期已存在则更新
		result, err := s.db.Exec(fmt.Sprintf(`
			INSERT INTO %s (date, open, high, low, close, volume, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			ON CONFLICT(date) DO UPDATE SET
				open = excluded.open,
				high = excluded.high,
				low = excluded.low,
				close = excluded.close,
				volume = excluded.volume,
				updated_at = CURRENT_TIMESTAMP
		`, tableName), date, open, high, low, close, volume)
		
		if err != nil {
			return addedCount, updatedCount, fmt.Errorf("插入/更新 K 线数据失败: %w", err)
		}
		
		// 检查是否是新插入还是更新
		rowsAffected, err := result.RowsAffected()
		if err == nil && rowsAffected > 0 {
			// 这里无法直接区分是插入还是更新，所以我们计数所有操作
			addedCount++
		}
	}
	
	return addedCount, updatedCount, nil
}

// GetLatestKLineDate 获取指定股票在本地缓存中的最新 K 线日期
func (s *DBService) GetLatestKLineDate(code string) (string, error) {
	tableName := fmt.Sprintf("kline_%s", code)
	
	var latestDate string
	err := s.db.QueryRow(fmt.Sprintf(`
		SELECT date FROM %s ORDER BY date DESC LIMIT 1
	`, tableName)).Scan(&latestDate)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // 表为空或不存在
		}
		return "", fmt.Errorf("查询最新 K 线日期失败: %w", err)
	}
	
	return latestDate, nil
}

// GetKLineDataFromCache 从本地缓存获取 K 线数据
func (s *DBService) GetKLineDataFromCache(code string, limit int) ([]map[string]interface{}, error) {
	tableName := fmt.Sprintf("kline_%s", code)
	
	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT date, open, high, low, close, volume FROM %s
		ORDER BY date DESC LIMIT ?
	`, tableName), limit)
	
	if err != nil {
		return nil, fmt.Errorf("查询 K 线数据失败: %w", err)
	}
	defer rows.Close()
	
	var klines []map[string]interface{}
	for rows.Next() {
		var date string
		var open, high, low, close float64
		var volume int64
		
		if err := rows.Scan(&date, &open, &high, &low, &close, &volume); err != nil {
			return nil, fmt.Errorf("扫描 K 线数据失败: %w", err)
		}
		
		klines = append(klines, map[string]interface{}{
			"date":   date,
			"open":   open,
			"high":   high,
			"low":    low,
			"close":  close,
			"volume": volume,
		})
	}
	
	// 反转切片，使其按日期升序排列
	for i, j := 0, len(klines)-1; i < j; i, j = i+1, j-1 {
		klines[i], klines[j] = klines[j], klines[i]
	}
	
	return klines, nil
}

// GetKLineCountByCode 获取指定股票的 K 线数据总数
func (s *DBService) GetKLineCountByCode(code string) (int, error) {
	tableName := fmt.Sprintf("kline_%s", code)
	
	var count int
	err := s.db.QueryRow(fmt.Sprintf(`
		SELECT COUNT(*) FROM %s
	`, tableName)).Scan(&count)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("查询 K 线数据总数失败: %w", err)
	}
	
	return count, nil
}

// GetAllSyncedStocks 获取所有已同步的股票列表
func (s *DBService) GetAllSyncedStocks() ([]string, error) {
	// 查询所有 kline_* 表
	rows, err := s.db.Query(`
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name LIKE 'kline_%'
	`)
	
	if err != nil {
		return nil, fmt.Errorf("查询已同步股票列表失败: %w", err)
	}
	defer rows.Close()
	
	var stocks []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("扫描表名失败: %w", err)
		}
		
		// 从表名中提取股票代码
		code := tableName[6:] // 移除 "kline_" 前缀
		stocks = append(stocks, code)
	}
	
	return stocks, nil
}

// ClearKLineCacheTable 清除指定股票的 K 线缓存表
func (s *DBService) ClearKLineCacheTable(code string) error {
	tableName := fmt.Sprintf("kline_%s", code)
	
	_, err := s.db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tableName))
	if err != nil {
		return fmt.Errorf("清除 K 线缓存表失败: %w", err)
	}
	
	return nil
}
