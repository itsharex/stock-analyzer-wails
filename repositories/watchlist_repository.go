package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"


	"go.uber.org/zap"
)

// WatchlistRepository 定义自选股持久化接口
type WatchlistRepository interface {
	Add(stock *models.StockData) error
	Remove(code string) error
	GetAll() ([]*models.StockData, error)
}

// SQLiteWatchlistRepository 基于 SQLite 的实现
type SQLiteWatchlistRepository struct {
	db *sql.DB
}

// NewSQLiteWatchlistRepository 构造函数
func NewSQLiteWatchlistRepository(db *sql.DB) *SQLiteWatchlistRepository {
	return &SQLiteWatchlistRepository{db: db}
}

func (r *SQLiteWatchlistRepository) GetAll() ([]*models.StockData, error) {
	start := time.Now()
	
	rows, err := r.db.Query("SELECT code, name, data FROM watchlist ORDER BY added_at DESC")
	if err != nil {
		return nil, fmt.Errorf("查询自选股列表失败: %w", err)
	}
	defer rows.Close()

	stocks := make([]*models.StockData, 0)
	for rows.Next() {
		var code, name, data string
		if err := rows.Scan(&code, &name, &data); err != nil {
			logger.Error("扫描自选股数据失败", zap.Error(err))
			continue
		}

		var stock models.StockData
		if err := json.Unmarshal([]byte(data), &stock); err != nil {
			logger.Error("解析自选股 JSON 数据失败", zap.String("code", code), zap.Error(err))
			continue
		}
		stocks = append(stocks, &stock)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历自选股结果集失败: %w", err)
	}

	logger.Debug("成功获取自选股列表",
		zap.String("module", "repositories.watchlist"),
		zap.String("op", "get_all_watchlist"),
		zap.Int("stocks_count", len(stocks)),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return stocks, nil
}

func (r *SQLiteWatchlistRepository) Add(stock *models.StockData) error {
	start := time.Now()
	
	// 序列化 StockData
	data, err := json.Marshal(stock)
	if err != nil {
		return fmt.Errorf("序列化 StockData 失败: %w", err)
	}

	// 使用 INSERT OR REPLACE 实现 upsert 逻辑
	query := `
		INSERT OR REPLACE INTO watchlist (code, name, data) 
		VALUES (?, ?, ?)
	`
	_, err = r.db.Exec(query, stock.Code, stock.Name, string(data))
	if err != nil {
		return fmt.Errorf("添加/更新自选股失败: %w", err)
	}

	logger.Info("成功添加/更新股票到自选股",
		zap.String("module", "repositories.watchlist"),
		zap.String("op", "add_to_watchlist"),
		zap.String("stock_code", stock.Code),
		zap.String("stock_name", stock.Name),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return nil
}

func (r *SQLiteWatchlistRepository) Remove(code string) error {
	start := time.Now()
	
	query := `DELETE FROM watchlist WHERE code = ?`
	result, err := r.db.Exec(query, code)
	if err != nil {
		return fmt.Errorf("从自选股移除股票失败: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	
	logger.Info("成功从自选股移除股票",
		zap.String("module", "repositories.watchlist"),
		zap.String("op", "remove_from_watchlist"),
		zap.String("stock_code", code),
		zap.Int64("rows_affected", rowsAffected),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return nil
}
