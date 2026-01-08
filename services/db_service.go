package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"stock-analyzer-wails/internal/logger"

	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

// DBService 数据库服务
type DBService struct {
	db     *sql.DB
	dbPath string // 数据库文件的绝对路径
}

// NewDBService 初始化数据库连接并创建表
func NewDBService() (*DBService, error) {
	appDir := GetAppDataDir()
	dbPath := filepath.Join(appDir, "stock_analyzer.db")
	logger.Info("开始初始化 SQLite 数据库服务",
		zap.String("module", "services.db"),
		zap.String("op", "NewDBService"),
		zap.String("appDir", appDir),
		zap.String("dbPath", dbPath),
	)

	// 确保目录存在（GetAppDataDir 通常已创建，但这里做一次兜底）
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		logger.Error("创建数据库目录失败",
			zap.String("module", "services.db"),
			zap.String("op", "NewDBService"),
			zap.String("step", "mkdir"),
			zap.String("dbDir", filepath.Dir(dbPath)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("创建数据库目录失败: %w", err)
	}

	// 检查数据库文件是否存在，如果不存在则创建
	_, err := os.Stat(dbPath)
	isNewDB := os.IsNotExist(err)
	logger.Info("数据库文件状态",
		zap.String("module", "services.db"),
		zap.String("op", "NewDBService"),
		zap.Bool("isNewDB", isNewDB),
		zap.String("dbPath", dbPath),
	)

	// 连接数据库
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		logger.Error("无法连接到数据库（首次 Open）",
			zap.String("module", "services.db"),
			zap.String("op", "NewDBService"),
			zap.String("step", "open_1"),
			zap.String("dbPath", dbPath),
			zap.Error(err),
		)
		return nil, fmt.Errorf("无法连接到数据库: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(1) // SQLite 建议最大连接数为 1
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// 设置 SQLite 忙等待超时时间（5秒），避免并发写入时立即失败
	// 在 DSN 中添加 _busy_timeout 参数（单位：毫秒）
	dbPathWithTimeout := fmt.Sprintf("%s?_busy_timeout=5000&_journal_mode=WAL", dbPath)
	db.Close()
	db, err = sql.Open("sqlite", dbPathWithTimeout)
	if err != nil {
		logger.Error("无法连接到数据库（带 busy_timeout/WAL）",
			zap.String("module", "services.db"),
			zap.String("op", "NewDBService"),
			zap.String("step", "open_2"),
			zap.String("dsn", dbPathWithTimeout),
			zap.Error(err),
		)
		return nil, fmt.Errorf("无法连接到数据库: %w", err)
	}

	// 立即验证连接可用性（避免延迟到首次 Query/Exec 才报错）
	if err := db.Ping(); err != nil {
		_ = db.Close()
		logger.Error("数据库连接不可用（Ping 失败）",
			zap.String("module", "services.db"),
			zap.String("op", "NewDBService"),
			zap.String("step", "ping"),
			zap.String("dbPath", dbPath),
			zap.Error(err),
		)
		return nil, fmt.Errorf("数据库连接不可用: %w", err)
	}

	svc := &DBService{db: db, dbPath: dbPath}

	// 始终调用 initTables() 确保所有表都存在（使用 CREATE TABLE IF NOT EXISTS）
	// 这样即使数据库文件存在但缺少某些表，也能自动创建
	logger.Info("开始初始化数据库表结构", zap.String("path", dbPath))
	if err := svc.initTables(); err != nil {
		_ = db.Close()
		logger.Error("初始化数据库表结构失败",
			zap.String("module", "services.db"),
			zap.String("op", "NewDBService"),
			zap.String("step", "initTables"),
			zap.String("dbPath", dbPath),
			zap.Error(err),
		)
		return nil, fmt.Errorf("初始化数据库表失败: %w", err)
	}

	if isNewDB {
		logger.Info("数据库文件不存在，已完成初始化")
	} else {
		logger.Info("数据库表结构检查完成")
	}

	return svc, nil
}

// GetDBPath 返回数据库文件的绝对路径
func (s *DBService) GetDBPath() string {
	return s.dbPath
}

// initTables 初始化数据库表结构
func (s *DBService) initTables() error {
	// 开启事务
	tx, err := s.db.Begin()
	if err != nil {
		logger.Error("开启初始化表结构事务失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("step", "begin_tx"),
			zap.Error(err),
		)
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
		logger.Error("创建 watchlist 表失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("table", "watchlist"),
			zap.Error(err),
		)
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
		logger.Error("创建 alerts 表失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("table", "alerts"),
			zap.Error(err),
		)
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
		logger.Error("创建 alert_history 表失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("table", "alert_history"),
			zap.Error(err),
		)
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
		logger.Error("创建 positions 表失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("table", "positions"),
			zap.Error(err),
		)
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
		logger.Error("创建 config 表失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("table", "config"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 config 表失败: %w", err)
	}

	// 6. Sync History 表
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS sync_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			stock_code TEXT NOT NULL,
			stock_name TEXT NOT NULL,
			sync_type TEXT NOT NULL, -- 'single' or 'batch'
			start_date TEXT NOT NULL,
			end_date TEXT NOT NULL,
			status TEXT NOT NULL, -- 'success' or 'failed'
			records_added INTEGER DEFAULT 0,
			records_updated INTEGER DEFAULT 0,
			duration INTEGER DEFAULT 0, -- 耗时（秒）
			error_msg TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 sync_history 表失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("table", "sync_history"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 sync_history 表失败: %w", err)
	}

	// 7. Strategy Config 表
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS strategy_config (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			strategy_type TEXT NOT NULL, -- 'simple_ma', 'macd', etc.
			parameters TEXT NOT NULL, -- JSON 格式的策略参数
			last_backtest_result TEXT, -- 最后一次回测结果（JSON）
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 strategy_config 表失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("table", "strategy_config"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 strategy_config 表失败: %w", err)
	}

	// 8. Stocks 表 - 市场全量股票列表
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS stocks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL,
			name TEXT NOT NULL,
			market TEXT NOT NULL, -- 'SH', 'SZ', 'BJ'
			full_code TEXT NOT NULL UNIQUE, -- 如 'SH600519'
			type TEXT, -- 板块类型：主板, 创业板, 科创板, 北交所
			is_active INTEGER DEFAULT 1, -- 1: 正常, 0: 退市/停牌
			price REAL, -- 最新价
			change_rate REAL, -- 涨跌幅(%)
			change_amount REAL, -- 涨跌额
			volume REAL, -- 成交量(手)
			amount REAL, -- 成交额
			amplitude REAL, -- 振幅(%)
			high REAL, -- 最高价
			low REAL, -- 最低价
			open REAL, -- 开盘价
			pre_close REAL, -- 昨收
			turnover REAL, -- 换手率(%)
			volume_ratio REAL, -- 量比
			pe REAL, -- 市盈率
			warrant_ratio REAL, -- 委比(%)
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(code)
		);
	`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 stocks 表失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("table", "stocks"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 stocks 表失败: %w", err)
	}

	// 创建 stocks 表的索引
	_, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_stocks_name ON stocks(name);`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 stocks.name 索引失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("index", "idx_stocks_name"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 stocks.name 索引失败: %w", err)
	}

	_, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_stocks_code ON stocks(code);`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 stocks.code 索引失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("index", "idx_stocks_code"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 stocks.code 索引失败: %w", err)
	}

	_, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_stocks_market ON stocks(market);`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 stocks.market 索引失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("index", "idx_stocks_market"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 stocks.market 索引失败: %w", err)
	}

	_, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_stocks_full_code ON stocks(full_code);`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 stocks.full_code 索引失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("index", "idx_stocks_full_code"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 stocks.full_code 索引失败: %w", err)
	}

	// 9. Price Threshold Alerts 表 - 价格预警配置
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS price_threshold_alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			stock_code TEXT NOT NULL,
			stock_name TEXT NOT NULL,
			alert_type TEXT NOT NULL, -- 'price_change', 'target_price', 'stop_loss', 'high_low', 'price_range', 'ma_deviation', 'combined'
			conditions TEXT NOT NULL, -- JSON格式: [{"field": "...", "operator": "...", "value": "..."}]
			is_active BOOLEAN DEFAULT 1,
			sensitivity REAL DEFAULT 0.001, -- 价格波动容差
			cooldown_hours INTEGER DEFAULT 1, -- 冷却时间（小时）
			post_trigger_action TEXT DEFAULT 'continue', -- 'continue', 'disable', 'once'
			enable_sound BOOLEAN DEFAULT 1,
			enable_desktop BOOLEAN DEFAULT 1,
			template_id TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_triggered_at DATETIME
		);
	`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 price_threshold_alerts 表失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("table", "price_threshold_alerts"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 price_threshold_alerts 表失败: %w", err)
	}

	// 创建 price_threshold_alerts 表的索引
	_, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_price_alerts_stock_code ON price_threshold_alerts(stock_code);`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 price_threshold_alerts.stock_code 索引失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("index", "idx_price_alerts_stock_code"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 price_threshold_alerts.stock_code 索引失败: %w", err)
	}

	_, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_price_alerts_is_active ON price_threshold_alerts(is_active);`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 price_threshold_alerts.is_active 索引失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("index", "idx_price_alerts_is_active"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 price_threshold_alerts.is_active 索引失败: %w", err)
	}

	// 10. Price Alert Templates 表 - 预警模板
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS price_alert_templates (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			alert_type TEXT NOT NULL,
			conditions TEXT NOT NULL, -- JSON格式的模板条件
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 price_alert_templates 表失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("table", "price_alert_templates"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 price_alert_templates 表失败: %w", err)
	}

	// 11. Price Alert Trigger History 表 - 预警触发历史
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS price_alert_trigger_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			alert_id INTEGER NOT NULL,
			stock_code TEXT NOT NULL,
			stock_name TEXT NOT NULL,
			alert_type TEXT NOT NULL,
			trigger_price REAL,
			trigger_message TEXT,
			triggered_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 price_alert_trigger_history 表失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("table", "price_alert_trigger_history"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 price_alert_trigger_history 表失败: %w", err)
	}

	// 创建 price_alert_trigger_history 表的索引
	_, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_alert_history_alert_id ON price_alert_trigger_history(alert_id);`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 price_alert_trigger_history.alert_id 索引失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("index", "idx_alert_history_alert_id"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 price_alert_trigger_history.alert_id 索引失败: %w", err)
	}

	_, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_alert_history_stock_code ON price_alert_trigger_history(stock_code);`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 price_alert_trigger_history.stock_code 索引失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("index", "idx_alert_history_stock_code"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 price_alert_trigger_history.stock_code 索引失败: %w", err)
	}

	_, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_alert_history_triggered_at ON price_alert_trigger_history(triggered_at);`)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("创建 price_alert_trigger_history.triggered_at 索引失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("index", "idx_alert_history_triggered_at"),
			zap.Error(err),
		)
		return fmt.Errorf("创建 price_alert_trigger_history.triggered_at 索引失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交初始化表结构事务失败",
			zap.String("module", "services.db"),
			zap.String("op", "initTables"),
			zap.String("step", "commit"),
			zap.Error(err),
		)
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

	// 插入默认预警模板
	if err := s.insertDefaultAlertTemplates(); err != nil {
		logger.Warn("插入默认预警模板失败", zap.Error(err))
	}

	return nil
}

// insertDefaultAlertTemplates 插入默认预警模板
func (s *DBService) insertDefaultAlertTemplates() error {
	templates := []struct {
		ID          string
		Name        string
		Description string
		AlertType   string
		Conditions  string
	}{
		{
			ID:          "template_price_change_5",
			Name:        "涨跌幅预警（5%）",
			Description: "当日涨跌幅度超过5%时触发预警",
			AlertType:   "price_change",
			Conditions:  `[{"field": "price_change_percent", "operator": ">", "value": 5}]`,
		},
		{
			ID:          "template_price_change_neg_5",
			Name:        "涨跌幅预警（-5%）",
			Description: "当日涨跌幅度低于-5%时触发预警",
			AlertType:   "price_change",
			Conditions:  `[{"field": "price_change_percent", "operator": "<", "value": -5}]`,
		},
		{
			ID:          "template_target_price",
			Name:        "目标价预警",
			Description: "价格达到目标价时触发预警",
			AlertType:   "target_price",
			Conditions:  `[{"field": "close_price", "operator": ">=", "value": 0.0}]`,
		},
		{
			ID:          "template_stop_loss",
			Name:        "止损价预警",
			Description: "价格跌破止损价时触发预警",
			AlertType:   "stop_loss",
			Conditions:  `[{"field": "close_price", "operator": "<=", "value": 0.0}]`,
		},
		{
			ID:          "template_high_new",
			Name:        "突破历史新高",
			Description: "价格突破历史最高价时触发预警",
			AlertType:   "high_low",
			Conditions:  `[{"field": "high_price", "operator": ">", "value": 0.0, "reference": "historical_high"}]`,
		},
		{
			ID:          "template_low_new",
			Name:        "跌破历史新低",
			Description: "价格跌破历史最低价时触发预警",
			AlertType:   "high_low",
			Conditions:  `[{"field": "low_price", "operator": "<", "value": 0.0, "reference": "historical_low"}]`,
		},
		{
			ID:          "template_ma5_golden_cross",
			Name:        "MA5金叉MA20",
			Description: "5日均线上穿20日均线时触发预警",
			AlertType:   "ma_deviation",
			Conditions:  `[{"field": "ma5", "operator": ">", "value": 0.0, "reference": "ma20"}]`,
		},
		{
			ID:          "template_volume_surge",
			Name:        "放量预警",
			Description: "成交量超过平均2倍时触发预警",
			AlertType:   "combined",
			Conditions:  `[{"field": "volume_ratio", "operator": ">", "value": 2}]`,
		},
	}

	for _, tmpl := range templates {
		_, err := s.db.Exec(`
			INSERT OR IGNORE INTO price_alert_templates (id, name, description, alert_type, conditions)
			VALUES (?, ?, ?, ?, ?)
		`, tmpl.ID, tmpl.Name, tmpl.Description, tmpl.AlertType, tmpl.Conditions)
		if err != nil {
			logger.Warn("插入预警模板失败",
				zap.String("template_id", tmpl.ID),
				zap.Error(err),
			)
		}
	}

	logger.Info("默认预警模板插入完成")
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
// 使用事务和重试机制来处理数据库锁定
func (s *DBService) InsertOrUpdateKLineData(code string, klines []map[string]interface{}) (int64, int64, error) {
	tableName := fmt.Sprintf("kline_%s", code)

	// 先确保表存在
	if err := s.CreateKLineCacheTable(code); err != nil {
		return 0, 0, err
	}

	// 使用重试机制处理数据库锁定
	maxRetries := 3
	var addedCount, updatedCount int64
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		// 开启事务
		tx, err := s.db.Begin()
		if err != nil {
			lastErr = fmt.Errorf("开启事务失败: %w", err)
			time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
			continue
		}

		addedCount = 0
		updatedCount = 0

		// 准备批量插入语句
		insertSQL := fmt.Sprintf(`
			INSERT INTO %s (date, open, high, low, close, volume, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			ON CONFLICT(date) DO UPDATE SET
				open = excluded.open,
				high = excluded.high,
				low = excluded.low,
				close = excluded.close,
				volume = excluded.volume,
				updated_at = CURRENT_TIMESTAMP
		`, tableName)

		stmt, err := tx.Prepare(insertSQL)
		if err != nil {
			_ = tx.Rollback()
			lastErr = fmt.Errorf("准备插入语句失败: %w", err)
			time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
			continue
		}
		defer stmt.Close()

		// 批量插入数据
		for _, kline := range klines {
			date := kline["date"].(string)
			open := kline["open"].(float64)
			high := kline["high"].(float64)
			low := kline["low"].(float64)
			close := kline["close"].(float64)
			volume := kline["volume"].(int64)

			_, err := stmt.Exec(date, open, high, low, close, volume)
			if err != nil {
				_ = tx.Rollback()
				lastErr = fmt.Errorf("插入 K 线数据失败: %w", err)
				time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
				break // 退出内层循环，继续重试
			}
			addedCount++
		}

		// 检查是否所有数据都插入成功
		if int(addedCount) != len(klines) {
			continue // 继续重试
		}

		// 提交事务
		if err := tx.Commit(); err != nil {
			lastErr = fmt.Errorf("提交事务失败: %w", err)
			time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
			continue
		}

		// 成功
		return addedCount, updatedCount, nil
	}

	// 所有重试都失败
	return 0, 0, fmt.Errorf("插入/更新 K 线数据失败（已重试 %d 次）: %w", maxRetries, lastErr)
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

// GetKLineDataWithPagination 获取指定股票的 K 线数据（支持分页和日期筛选）
func (s *DBService) GetKLineDataWithPagination(code string, startDate string, endDate string, page int, pageSize int) ([]map[string]interface{}, int, error) {
	tableName := fmt.Sprintf("kline_%s", code)

	// 检查表是否存在，如果不存在则返回空数据
	var tableExists bool
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM sqlite_master
		WHERE type='table' AND name=?
	`, tableName).Scan(&tableExists)

	if err != nil {
		logger.Warn("查询表是否存在失败，返回空数组", zap.String("tableName", tableName), zap.Error(err))
		return []map[string]interface{}{}, 0, nil
	}

	if !tableExists {
		// 表不存在，返回初始化的空数组而不是 nil
		logger.Info("K线表不存在，返回空数组", zap.String("tableName", tableName))
		return []map[string]interface{}{}, 0, nil
	}

	// 构建基础查询和参数
	var conditions []string
	var params []interface{}

	// 日期筛选
	if startDate != "" {
		conditions = append(conditions, "date >= ?")
		params = append(params, startDate)
	}
	if endDate != "" {
		conditions = append(conditions, "date <= ?")
		params = append(params, endDate)
	}

	// 构建 WHERE 子句
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", tableName, whereClause)
	var totalCount int
	paramsForCount := make([]interface{}, len(params))
	copy(paramsForCount, params)
	err = s.db.QueryRow(countQuery, paramsForCount...).Scan(&totalCount)
	if err != nil {
		// 如果查询失败（可能是表不存在等），返回空数组
		return []map[string]interface{}{}, 0, nil
	}

	// 查询分页数据
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT date, open, high, low, close, volume FROM %s
		%s
		ORDER BY date DESC
		LIMIT ? OFFSET ?
	`, tableName, whereClause)

	params = append(params, pageSize, offset)

	rows, err := s.db.Query(query, params...)
	if err != nil {
		// 如果查询失败，返回空数组
		return []map[string]interface{}{}, 0, nil
	}
	defer rows.Close()

	var klines []map[string]interface{}
	for rows.Next() {
		var date string
		var open, high, low, close float64
		var volume int64

		if err := rows.Scan(&date, &open, &high, &low, &close, &volume); err != nil {
			return []map[string]interface{}{}, 0, nil
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

	return klines, totalCount, nil
}
