package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stock-analyzer-wails/internal/logger"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// StockMarketService 市场股票服务
type StockMarketService struct {
	dbService *DBService
	client    *http.Client
}

// NewStockMarketService 创建市场股票服务
func NewStockMarketService(dbService *DBService) *StockMarketService {
	return &StockMarketService{
		dbService: dbService,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SyncStocksRequest 同步股票请求
type SyncStocksRequest struct {
	PN      int      `json:"pn"`      // 页码
	PZ      int      `json:"pz"`      // 每页数量
	FS      string   `json:"fs"`      // 市场范围
	Fields  []string `json:"fields"`  // 返回字段
}

// StockMarketData 市场股票数据
type StockMarketData struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Market string `json:"market"`   // SH, SZ, BJ
	Type   string `json:"type"`     // 主板, 创业板, 科创板, 北交所
	// 实时行情数据
	Price         float64 `json:"price"`         // 最新价
	ChangeRate    float64 `json:"changeRate"`    // 涨跌幅(%)
	ChangeAmount  float64 `json:"changeAmount"`  // 涨跌额
	Volume        float64 `json:"volume"`        // 成交量(手)
	Amount        float64 `json:"amount"`        // 成交额
	Amplitude     float64 `json:"amplitude"`     // 振幅(%)
	High          float64 `json:"high"`         // 最高价
	Low           float64 `json:"low"`          // 最低价
	Open          float64 `json:"open"`         // 开盘价
	PreClose      float64 `json:"preClose"`      // 昨收
	Turnover      float64 `json:"turnover"`      // 换手率(%)
	VolumeRatio   float64 `json:"volumeRatio"`   // 量比
	PE            float64 `json:"pe"`            // 市盈率
	WarrantRatio  float64 `json:"warrantRatio"`  // 委比(%)
	IsActive      int     `json:"isActive"`      // 是否在交易
	UpdatedAt     string  `json:"updatedAt"`     // 最后更新时间
}

// SyncStocksResult 同步结果
type SyncStocksResult struct {
	Total      int     `json:"total"`      // 总记录数
	Processed  int     `json:"processed"`  // 已处理记录数
	Inserted   int     `json:"inserted"`   // 新增记录数
	Updated    int     `json:"updated"`    // 更新记录数
	Duration   float64 `json:"duration"`   // 耗时（秒）
	Message    string  `json:"message"`    // 消息
}

// 东方财富API响应结构
type EastMoneyResponse struct {
	RC   int    `json:"rc"`
	RT   int    `json:"rt"`
	SVR  int    `json:"svr"`
	LT   int    `json:"lt"`
	Full int    `json:"full"`
	Data struct {
		Total int                    `json:"total"`
		Diff  map[string]interface{} `json:"diff"`
	} `json:"data"`
}

// SyncAllStocks 同步所有市场股票
func (s *StockMarketService) SyncAllStocks() (*SyncStocksResult, error) {
	startTime := time.Now()

	// 默认参数
	pn := 1
	pz := 5000
	fs := "m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23,m:0+t:81+s:2048"
	fields := "f12,f14,f2,f3,f62,f184,f66,f69,f72,f75,f78,f81,f84,f87,f64,f65"

	url := fmt.Sprintf(
		"https://push2.eastmoney.com/api/qt/clist/get?pn=%d&pz=%d&fs=%s&fields=%s",
		pn, pz, fs, fields,
	)

	logger.Info("开始同步市场股票", zap.String("url", url))

	// 请求接口
	resp, err := s.client.Get(url)
	if err != nil {
		logger.Error("请求东方财富API失败", zap.Error(err))
		return nil, fmt.Errorf("请求东方财富API失败: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var apiResp EastMoneyResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		logger.Error("解析API响应失败", zap.Error(err))
		return nil, fmt.Errorf("解析API响应失败: %w", err)
	}

	if apiResp.RC != 0 {
		return nil, fmt.Errorf("API返回错误: rc=%d", apiResp.RC)
	}

	// 获取数据库连接
	db := s.dbService.GetDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("开启事务失败: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()

	// 准备upsert语句
	upsertSQL := `
		INSERT INTO stocks (
			code, name, market, full_code, type, is_active,
			price, change_rate, change_amount, volume, amount,
			amplitude, high, low, open, pre_close, turnover,
			volume_ratio, pe, warrant_ratio, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(code) DO UPDATE SET
			name = excluded.name,
			market = excluded.market,
			full_code = excluded.full_code,
			type = excluded.type,
			is_active = excluded.is_active,
			price = excluded.price,
			change_rate = excluded.change_rate,
			change_amount = excluded.change_amount,
			volume = excluded.volume,
			amount = excluded.amount,
			amplitude = excluded.amplitude,
			high = excluded.high,
			low = excluded.low,
			open = excluded.open,
			pre_close = excluded.pre_close,
			turnover = excluded.turnover,
			volume_ratio = excluded.volume_ratio,
			pe = excluded.pe,
			warrant_ratio = excluded.warrant_ratio,
			updated_at = CURRENT_TIMESTAMP
	`

	stmt, err := tx.Prepare(upsertSQL)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("准备upsert语句失败: %w", err)
	}
	defer stmt.Close()

	// 遍历所有股票数据
	processed := 0
	inserted := 0
	updated := 0
	now := time.Now().Format("2006-01-02 15:04:05")

	for _, item := range apiResp.Data.Diff {
		// 解析股票数据
		stockData := s.parseStockItem(item, now)
		if stockData == nil {
			continue
		}

		// 执行upsert
		_, err := stmt.Exec(
			stockData.Code,
			stockData.Name,
			stockData.Market,
			stockData.Market+stockData.Code,
			stockData.Type,
			stockData.IsActive,
			stockData.Price,
			stockData.ChangeRate,
			stockData.ChangeAmount,
			stockData.Volume,
			stockData.Amount,
			stockData.Amplitude,
			stockData.High,
			stockData.Low,
			stockData.Open,
			stockData.PreClose,
			stockData.Turnover,
			stockData.VolumeRatio,
			stockData.PE,
			stockData.WarrantRatio,
		)

		if err != nil {
			logger.Error("插入股票数据失败", zap.String("code", stockData.Code), zap.Error(err))
			continue
		}

		processed++
		if processed == 1 {
			inserted++
		} else {
			updated++
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	duration := time.Since(startTime).Seconds()

	logger.Info("同步市场股票完成",
		zap.Int("total", apiResp.Data.Total),
		zap.Int("processed", processed),
		zap.Int("inserted", inserted),
		zap.Int("updated", updated),
		zap.Float64("duration", duration),
	)

	return &SyncStocksResult{
		Total:     apiResp.Data.Total,
		Processed: processed,
		Inserted:  inserted,
		Updated:   updated,
		Duration:  duration,
		Message:   "同步成功",
	}, nil
}

// parseStockItem 解析股票数据项
func (s *StockMarketService) parseStockItem(item interface{}, updatedAt string) *StockMarketData {
	data, ok := item.(map[string]interface{})
	if !ok {
		return nil
	}

	// 获取code
	code, ok := data["f12"].(string)
	if !ok || code == "" {
		return nil
	}

	// 获取name
	name, ok := data["f14"].(string)
	if !ok {
		name = ""
	}

	// 判断市场
	market := "SH"
	if strings.HasPrefix(code, "6") {
		market = "SH"
	} else if strings.HasPrefix(code, "0") || strings.HasPrefix(code, "3") {
		market = "SZ"
	} else if strings.HasPrefix(code, "8") || strings.HasPrefix(code, "4") {
		market = "BJ"
	}

	// 判断板块类型
	stockType := "主板"
	if strings.HasPrefix(code, "6") {
		stockType = "主板"
	} else if strings.HasPrefix(code, "0") {
		stockType = "主板"
	} else if strings.HasPrefix(code, "3") {
		stockType = "创业板"
	} else if strings.HasPrefix(code, "6") && len(code) == 6 {
		stockType = "科创板" // 实际科创板以688开头
	} else if strings.HasPrefix(code, "8") || strings.HasPrefix(code, "4") {
		stockType = "北交所"
	}

	// 修正科创板判断
	if strings.HasPrefix(code, "688") {
		stockType = "科创板"
	}

	// 解析价格相关字段（接口返回的值通常是×100）
	parseFloat := func(key string) float64 {
		if val, ok := data[key]; ok {
			switch v := val.(type) {
			case float64:
				return v / 100.0 // 除以100转换为真实值
			case string:
				f, _ := strconv.ParseFloat(v, 64)
				return f / 100.0
			}
		}
		return 0
	}

	// 解析数值字段
	parseInt := func(key string) int {
		if val, ok := data[key]; ok {
			switch v := val.(type) {
			case float64:
				return int(v)
			case string:
				i, _ := strconv.Atoi(v)
				return i
			}
		}
		return 0
	}

	// 解析成交量（手）
	volume := float64(parseInt("f69"))

	// 解析成交额
	amount := parseFloat("f62")

	stock := &StockMarketData{
		Code:         code,
		Name:         name,
		Market:       market,
		FullCode:     market + code,
		Type:         stockType,
		IsActive:     1,
		Price:        parseFloat("f2"),
		ChangeRate:   parseFloat("f3"),
		ChangeAmount: parseFloat("f72"),
		Volume:       volume,
		Amount:       amount,
		Amplitude:    parseFloat("f64"),
		High:         parseFloat("f65"),
		Low:          parseFloat("f66"),
		Open:         parseFloat("f81"),
		PreClose:     parseFloat("f78"),
		Turnover:     parseFloat("f75"),
		VolumeRatio:  parseFloat("f84"),
		PE:           parseFloat("f184"),
		WarrantRatio: parseFloat("f87"),
		UpdatedAt:    updatedAt,
	}

	// 判断是否在交易（价格不为0表示在交易）
	if stock.Price == 0 {
		stock.IsActive = 0
	}

	return stock
}

// GetStocksList 获取股票列表（分页）
func (s *StockMarketService) GetStocksList(page int, pageSize int, search string) ([]StockMarketData, int, error) {
	offset := (page - 1) * pageSize

	db := s.dbService.GetDB()

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if search != "" {
		whereClause += " AND (code LIKE ? OR name LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	// 查询总数
	var total int
	countSQL := "SELECT COUNT(*) FROM stocks " + whereClause
	err := db.QueryRow(countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("查询总数失败: %w", err)
	}

	// 查询列表
	querySQL := `
		SELECT id, code, name, market, full_code, type, is_active,
		       price, change_rate, change_amount, volume, amount,
		       amplitude, high, low, open, pre_close, turnover,
		       volume_ratio, pe, warrant_ratio, updated_at
		FROM stocks ` + whereClause + `
		ORDER BY code ASC
		LIMIT ? OFFSET ?
	`

	args = append(args, pageSize, offset)

	rows, err := db.Query(querySQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询股票列表失败: %w", err)
	}
	defer rows.Close()

	stocks := []StockMarketData{}
	for rows.Next() {
		var stock StockMarketData
		err := rows.Scan(
			&stock.ID,
			&stock.Code,
			&stock.Name,
			&stock.Market,
			&stock.FullCode,
			&stock.Type,
			&stock.IsActive,
			&stock.Price,
			&stock.ChangeRate,
			&stock.ChangeAmount,
			&stock.Volume,
			&stock.Amount,
			&stock.Amplitude,
			&stock.High,
			&stock.Low,
			&stock.Open,
			&stock.PreClose,
			&stock.Turnover,
			&stock.VolumeRatio,
			&stock.PE,
			&stock.WarrantRatio,
			&stock.UpdatedAt,
		)
		if err != nil {
			logger.Error("扫描股票数据失败", zap.Error(err))
			continue
		}
		stocks = append(stocks, stock)
	}

	return stocks, total, nil
}

// GetSyncStats 获取同步统计信息
func (s *StockMarketService) GetSyncStats() (map[string]interface{}, error) {
	db := s.dbService.GetDB()

	// 查询总数量
	var totalCount int
	err := db.QueryRow("SELECT COUNT(*) FROM stocks").Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("查询总数量失败: %w", err)
	}

	// 查询最近更新时间
	var lastUpdate string
	err = db.QueryRow("SELECT updated_at FROM stocks ORDER BY updated_at DESC LIMIT 1").Scan(&lastUpdate)
	if err != nil {
		lastUpdate = "未同步"
	}

	// 按市场统计
	marketStats := make(map[string]int)
	rows, err := db.Query("SELECT market, COUNT(*) as count FROM stocks GROUP BY market")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var market string
			var count int
			rows.Scan(&market, &count)
			marketStats[market] = count
		}
	}

	return map[string]interface{}{
		"totalCount":  totalCount,
		"lastUpdate":   lastUpdate,
		"marketStats":  marketStats,
	}, nil
}
