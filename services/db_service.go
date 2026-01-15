package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gorm_logger "gorm.io/gorm/logger"
)

// DBService 数据库服务
type DBService struct {
	db     *gorm.DB
	dbPath string // 数据库文件的绝对路径
}

// NewDBService 初始化数据库连接并创建表
func NewDBService() (*DBService, error) {
	appDir := GetAppDataDir()
	dbPath := filepath.Join(appDir, "stock_analyzer.db")
	logger.Info("开始初始化 SQLite 数据库服务 (GORM)",
		zap.String("module", "services.db"),
		zap.String("op", "NewDBService"),
		zap.String("appDir", appDir),
		zap.String("dbPath", dbPath),
	)

	// 确保目录存在
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

	// 检查数据库文件是否存在
	_, err := os.Stat(dbPath)
	isNewDB := os.IsNotExist(err)

	// 配置 GORM 连接
	// 使用 glebarez/sqlite (pure go)
	// 设置 busy_timeout 和 WAL 模式
	dsn := fmt.Sprintf("%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", dbPath)

	// 自定义 GORM Logger 以集成到应用的 logger 系统
	newLogger := gorm_logger.New(
		zapLogger{logger.L()}, // 适配器
		gorm_logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gorm_logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logger.Error("无法连接到数据库",
			zap.String("module", "services.db"),
			zap.String("op", "NewDBService"),
			zap.String("dbPath", dbPath),
			zap.Error(err),
		)
		return nil, fmt.Errorf("无法连接到数据库: %w", err)
	}

	// 获取底层的 sql.DB 对象以设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(1) // SQLite 建议最大连接数为 1 (即使是 WAL 模式，写入也需要锁)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 立即验证连接可用性
	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
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

	// 初始化表结构
	logger.Info("开始初始化数据库表结构", zap.String("path", dbPath))
	if err := svc.initTables(); err != nil {
		_ = sqlDB.Close()
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

// zapLogger 适配 GORM 的 Writer 接口
type zapLogger struct {
	*zap.Logger
}

func (l zapLogger) Printf(format string, args ...interface{}) {
	l.Sugar().Infof(format, args...)
}

// GetDBPath 返回数据库文件的绝对路径
func (s *DBService) GetDBPath() string {
	return s.dbPath
}

// initTables 初始化数据库表结构
func (s *DBService) initTables() error {
	// 使用 AutoMigrate 自动迁移表结构
	err := s.db.AutoMigrate(
		&models.WatchlistEntity{},
		&models.AlertEntity{},
		&models.AlertHistoryEntity{},
		&models.PositionEntity{},
		&models.ConfigEntity{},
		&models.SyncHistoryEntity{},
		&models.StrategyConfigEntity{},
		&models.StockEntity{},
		&models.PriceThresholdAlertEntity{},
		&models.PriceAlertTemplateEntity{},
		&models.PriceAlertTriggerHistoryEntity{},
		&models.StockMoneyFlowHistEntity{},
		&models.StockStrategySignalEntity{},
	)
	if err != nil {
		return fmt.Errorf("自动迁移表结构失败: %w", err)
	}

	// 插入默认配置
	return s.insertDefaultConfigs()
}

// GetDB 返回数据库连接对象
func (s *DBService) GetDB() *gorm.DB {
	return s.db
}

// Close 关闭数据库连接
func (s *DBService) Close() {
	sqlDB, err := s.db.DB()
	if err == nil && sqlDB != nil {
		_ = sqlDB.Close()
		logger.Info("数据库连接已关闭")
	}
}

// insertDefaultConfigs 插入默认配置项
func (s *DBService) insertDefaultConfigs() error {
	// 默认配置值
	defaults := []models.ConfigEntity{
		{Key: "trailing_stop_default_activation", Value: "0.05"}, // 默认盈利 5% 启动
		{Key: "trailing_stop_default_callback", Value: "0.03"},   // 默认回撤 3% 止盈
	}

	for _, config := range defaults {
		// 使用 Clauses(clause.OnConflict{DoNothing: true}) 实现 INSERT OR IGNORE
		if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&config).Error; err != nil {
			return fmt.Errorf("插入默认配置 (%s) 失败: %w", config.Key, err)
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
	templates := []models.PriceAlertTemplateEntity{
		{
			ID:          "template_price_change_5",
			Name:        "涨跌幅预警（5%）",
			Description: "当日涨跌幅度超过5%时触发预警",
			AlertType:   "price_change",
			Conditions:  `{"logic":"AND","conditions":[{"field":"price_change_percent","operator":">","value":5}]}`,
		},
		{
			ID:          "template_price_change_neg_5",
			Name:        "涨跌幅预警（-5%）",
			Description: "当日涨跌幅度低于-5%时触发预警",
			AlertType:   "price_change",
			Conditions:  `{"logic":"AND","conditions":[{"field":"price_change_percent","operator":"<","value":-5}]}`,
		},
		{
			ID:          "template_target_price",
			Name:        "目标价预警",
			Description: "价格达到目标价时触发预警",
			AlertType:   "target_price",
			Conditions:  `{"logic":"AND","conditions":[{"field":"close_price","operator":">=","value":0.0}]}`,
		},
		{
			ID:          "template_stop_loss",
			Name:        "止损价预警",
			Description: "价格跌破止损价时触发预警",
			AlertType:   "stop_loss",
			Conditions:  `{"logic":"AND","conditions":[{"field":"close_price","operator":"<=","value":0.0}]}`,
		},
		{
			ID:          "template_high_new",
			Name:        "突破历史新高",
			Description: "价格突破历史最高价时触发预警",
			AlertType:   "high_low",
			Conditions:  `{"logic":"AND","conditions":[{"field":"high_price","operator":">","value":0.0,"reference":"historical_high"}]}`,
		},
		{
			ID:          "template_low_new",
			Name:        "跌破历史新低",
			Description: "价格跌破历史最低价时触发预警",
			AlertType:   "high_low",
			Conditions:  `{"logic":"AND","conditions":[{"field":"low_price","operator":"<","value":0.0,"reference":"historical_low"}]}`,
		},
		{
			ID:          "template_ma5_golden_cross",
			Name:        "MA5金叉MA20",
			Description: "5日均线上穿20日均线时触发预警",
			AlertType:   "ma_deviation",
			Conditions:  `{"logic":"AND","conditions":[{"field":"ma5","operator":">","value":0.0,"reference":"ma20"}]}`,
		},
		{
			ID:          "template_volume_surge",
			Name:        "放量预警",
			Description: "成交量超过平均2倍时触发预警",
			AlertType:   "combined",
			Conditions:  `{"logic":"AND","conditions":[{"field":"volume_ratio","operator":">","value":2}]}`,
		},
	}

	for _, tmpl := range templates {
		if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&tmpl).Error; err != nil {
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
	// 使用 GORM 的 Table 方法指定表名，并使用 AutoMigrate 创建表
	if err := s.db.Table(tableName).AutoMigrate(&models.KLineEntity{}); err != nil {
		return fmt.Errorf("创建 K 线缓存表 %s 失败: %w", tableName, err)
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

	// 转换 map 到 KLineEntity 结构体切片
	var entities []models.KLineEntity
	for _, kline := range klines {
		entity := models.KLineEntity{
			Date:   kline["date"].(string),
			Open:   kline["open"].(float64),
			High:   kline["high"].(float64),
			Low:    kline["low"].(float64),
			Close:  kline["close"].(float64),
			Volume: kline["volume"].(int64),
		}
		entities = append(entities, entity)
	}

	if len(entities) == 0 {
		return 0, 0, nil
	}

	// 使用 GORM 的 CreateInBatches 和 Clauses 进行 upsert
	result := s.db.Table(tableName).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"open", "high", "low", "close", "volume", "updated_at"}),
	}).CreateInBatches(entities, 100) // 每批 100 条

	if result.Error != nil {
		return 0, 0, fmt.Errorf("插入/更新 K 线数据失败: %w", result.Error)
	}

	// 注意：GORM 的 RowsAffected 在 upsert 时可能不准确反映“新增”和“更新”的分别计数
	// 这里简单返回受影响的行数作为 addedCount (实际上是 added + updated)
	return result.RowsAffected, 0, nil
}

// GetLatestKLineDate 获取指定股票在本地缓存中的最新 K 线日期
func (s *DBService) GetLatestKLineDate(code string) (string, error) {
	tableName := fmt.Sprintf("kline_%s", code)

	var date string
	err := s.db.Table(tableName).Select("date").Order("date DESC").Limit(1).Scan(&date).Error

	if err != nil {
		return "", fmt.Errorf("查询最新 K 线日期失败: %w", err)
	}

	// 如果没有找到记录，date 会是空字符串，符合预期
	return date, nil
}

// GetKLineDataFromCache 从本地缓存获取 K 线数据
func (s *DBService) GetKLineDataFromCache(code string, limit int) ([]map[string]interface{}, error) {
	tableName := fmt.Sprintf("kline_%s", code)

	var entities []models.KLineEntity
	err := s.db.Table(tableName).Order("date DESC").Limit(limit).Find(&entities).Error
	if err != nil {
		// 检查表是否存在
		if strings.Contains(err.Error(), "no such table") {
			return nil, nil // 表不存在视为空
		}
		return nil, fmt.Errorf("查询 K 线数据失败: %w", err)
	}

	var klines []map[string]interface{}
	for _, e := range entities {
		klines = append(klines, map[string]interface{}{
			"date":   e.Date,
			"open":   e.Open,
			"high":   e.High,
			"low":    e.Low,
			"close":  e.Close,
			"volume": e.Volume,
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

	var count int64
	err := s.db.Table(tableName).Count(&count).Error
	if err != nil {
		if strings.Contains(err.Error(), "no such table") {
			return 0, nil
		}
		return 0, fmt.Errorf("查询 K 线数据总数失败: %w", err)
	}

	return int(count), nil
}

// GetAllSyncedStocks 获取所有已同步的股票列表
func (s *DBService) GetAllSyncedStocks() ([]string, error) {
	// 查询 sqlite_master 表
	var tables []string
	err := s.db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name LIKE 'kline_%'").Scan(&tables).Error
	if err != nil {
		return nil, fmt.Errorf("查询已同步股票列表失败: %w", err)
	}

	var stocks []string
	for _, tableName := range tables {
		code := strings.TrimPrefix(tableName, "kline_")
		stocks = append(stocks, code)
	}

	return stocks, nil
}

// ClearKLineCacheTable 清除指定股票的 K 线缓存表
func (s *DBService) ClearKLineCacheTable(code string) error {
	tableName := fmt.Sprintf("kline_%s", code)
	// 使用 Migrator 删除表
	if err := s.db.Migrator().DropTable(tableName); err != nil {
		return fmt.Errorf("清除 K 线缓存表失败: %w", err)
	}
	return nil
}

// GetKLineDataWithPagination 获取指定股票的 K 线数据（支持分页和日期筛选）
func (s *DBService) GetKLineDataWithPagination(code string, startDate string, endDate string, page int, pageSize int) ([]map[string]interface{}, int, error) {
	tableName := fmt.Sprintf("kline_%s", code)

	// 检查表是否存在
	if !s.db.Migrator().HasTable(tableName) {
		logger.Info("K线表不存在，返回空数组", zap.String("tableName", tableName))
		return []map[string]interface{}{}, 0, nil
	}

	tx := s.db.Table(tableName)

	if startDate != "" {
		tx = tx.Where("date >= ?", startDate)
	}
	if endDate != "" {
		tx = tx.Where("date <= ?", endDate)
	}

	var totalCount int64
	if err := tx.Count(&totalCount).Error; err != nil {
		return []map[string]interface{}{}, 0, nil
	}

	var entities []models.KLineEntity
	offset := (page - 1) * pageSize
	if err := tx.Order("date DESC").Limit(pageSize).Offset(offset).Find(&entities).Error; err != nil {
		return []map[string]interface{}{}, 0, nil
	}

	var klines []map[string]interface{}
	for _, e := range entities {
		klines = append(klines, map[string]interface{}{
			"date":   e.Date,
			"open":   e.Open,
			"high":   e.High,
			"low":    e.Low,
			"close":  e.Close,
			"volume": e.Volume,
		})
	}

	// 反转切片，使其按日期升序排列
	for i, j := 0, len(klines)-1; i < j; i, j = i+1, j-1 {
		klines[i], klines[j] = klines[j], klines[i]
	}

	return klines, int(totalCount), nil
}
