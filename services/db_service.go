package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"stock-analyzer-wails/internal/logger"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// DBService 数据库服务
type DBService struct {
	db *sql.DB
}

// NewDBService 初始化数据库连接并创建表
func NewDBService() (*DBService, error) {
	dbPath := filepath.Join(GetAppDataDir(), "stock_analyzer.db")
	
	// 检查数据库文件是否存在，如果不存在则创建
	_, err := os.Stat(dbPath)
	isNewDB := os.IsNotExist(err)

	// 连接数据库
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(1) // SQLite 建议最大连接数为 1
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	svc := &DBService{db: db}

	if isNewDB {
		logger.Info("数据库文件不存在，开始初始化表结构", zap.String("path", dbPath))
		if err := svc.initTables(); err != nil {
			db.Close()
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
			tx.Rollback()
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
		tx.Rollback()
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
		tx.Rollback()
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
		tx.Rollback()
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
		tx.Rollback()
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
		tx.Rollback()
		return fmt.Errorf("创建 config 表失败: %w", err)
	}

	// 提交事务
	return tx.Commit()
}

// GetDB 返回数据库连接对象
func (s *DBService) GetDB() *sql.DB {
	return s.db
}

// Close 关闭数据库连接
func (s *DBService) Close() {
	if s.db != nil {
		s.db.Close()
		logger.Info("数据库连接已关闭")
	}
}
