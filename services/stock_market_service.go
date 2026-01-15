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

// 东方财富API字段定义
const (
	FieldCode             = "f12"  // f12:股票代码
	FieldMarket           = "f13"  // f13:市场
	FieldName             = "f14"  // f14:股票名称
	FieldPrice            = "f2"   // f2:最新价
	FieldChangeRate       = "f3"   // f3:涨跌幅
	FieldChangeAmount     = "f4"   // f4:涨跌额
	FieldVolume           = "f5"   // f5:总手（VOL）/成交量
	FieldAmount           = "f6"   // f6:成交额
	FieldAmplitude        = "f7"   // f7:振幅
	FieldTurnover         = "f8"   // f8:换手率
	FieldPE               = "f9"   // f9:市盈率(动态)
	FieldVolumeRatio      = "f10"  // f10:量比
	Field5MinChangeRate   = "f11"  // f11:5分钟涨跌幅
	FieldHigh             = "f15"  // f15:今日最高
	FieldLow              = "f16"  // f16:今日最低
	FieldOpen             = "f17"  // f17:今开
	FieldPreClose         = "f18"  // f18:昨收
	FieldMarketCap        = "f20"  // f20:总市值
	FieldCircCap          = "f21"  // f21:流通市值
	FieldRiseSpeed        = "f22"  // f22:涨速
	FieldPB               = "f23"  // f23:市净率
	Field60DayChange      = "f24"  // f24:60日涨跌幅
	FieldYTDChange        = "f25"  // f25:年初至今涨跌幅
	FieldListDate         = "f26"  // f26:上市时间
	FieldPreSettlePrice   = "f28"  // f28:昨日结算价
	FieldLastVol          = "f30"  // f30:每天最后一笔交易的成交量
	FieldBuyPrice         = "f31"  // f31:现汇买入价
	FieldSellPrice        = "f32"  // f32:现汇卖出价
	FieldWarrantRatio     = "f33"  // f33:委比
	FieldBuyVol           = "f34"  // f34:外盘
	FieldSellVol          = "f35"  // f35:内盘
	FieldAOE              = "f37"  // f37:净资产收益率加权（AOE）最近季度
	FieldTotalShare       = "f38"  // f38:总股本
	FieldCircAShare       = "f39"  // f39:流通A股（万股）
	FieldTotalRevenue     = "f40"  // f40:总营收（最近季度）
	FieldRevenueYOY       = "f41"  // f41:总营收同比
	FieldTotalProfit      = "f44"  // f44:总利润（最近季度）
	FieldNetProfit        = "f45"  // f45:净利润（最近季度）
	FieldProfitGrowth     = "f46"  // f46:净利润增长率（%）（同比）（最近季度）
	FieldUndividedProfit  = "f48"  // f48:每股未分配利润
	FieldGrossMargin      = "f49"  // f49:毛利率（最近季度）
	FieldTotalAssets      = "f50"  // f50:总资产（最近季度）
	FieldDebtRatio        = "f57"  // f57:负债率
	FieldEquity           = "f58"  // f58:股东权益
	FieldMainNetInflow    = "f62"  // f62:今日主力净流入
	FieldSuperBuy         = "f64"  // f64:超大单流入
	FieldSuperSell        = "f65"  // f65:超大单流出
	FieldSuperNetInflow   = "f66"  // f66:今日超大单净流入
	FieldSuperNetRatio    = "f69"  // f69:超大单净比
	FieldBigBuy           = "f70"  // f70:大单流入
	FieldBigSell          = "f71"  // f71:大单流出
	FieldBigNetInflow     = "f72"  // f72:今日大单净流入
	FieldBigNetRatio      = "f75"  // f75:大单净比
	FieldMidBuy           = "f76"  // f76:中单流入
	FieldMidSell          = "f77"  // f77:中单流出
	FieldMidNetInflow     = "f78"  // f78:今日中单净流入
	FieldMidNetRatio      = "f81"  // f81:中单净比（%）
	FieldSmallBuy         = "f82"  // f82:小单流入
	FieldSmallSell        = "f83"  // f83:小单流出
	FieldSmallNetInflow   = "f84"  // f84:进入小单净流入
	FieldSmallNetRatio    = "f87"  // f87:小单净比
	FieldIndustry         = "f100" // f100:行业
	FieldRegion           = "f102" // f102:地区板块
	FieldRemark           = "f103" // f103:备注
	FieldRiseCount        = "f104" // f104:上涨家数
	FieldFallCount        = "f105" // f105:下跌家数
	FieldFlatCount        = "f106" // f106:平盘家数
	FieldEPS1             = "f112" // f112:每股收益（一）
	FieldNetAssetPerShare = "f113" // f113:每股净资产
	FieldPEStatic         = "f114" // f114:市盈率（静）
	FieldPETTM            = "f115" // f115:市盈率（TTM）
	FieldTradeTime        = "f124" // f124:交易时间
	FieldLeaderStock      = "f128" // f128:板块领涨股
	FieldNetProfitTTM     = "f129" // f129:净利润TTM
	FieldPSTTM            = "f130" // f130:市销率TTM
	FieldPCTTM            = "f131" // f131:市现率TTM
	FieldRevenueTTM       = "f132" // f132:总营业收入TTM
	FieldDividendRate     = "f133" // f133:股息率
	FieldIndustryCount    = "f134" // f134:行业板块的成分股数
	FieldNetAssets        = "f135" // f135:净资产
	FieldNetProfitTTM2    = "f138" // f138:净利润TTM
	FieldMain5DayNet      = "f164" // f164:5日主力净额
	FieldSuper5DayNet     = "f166" // f166:5日超大单净额
	FieldBig5DayNet       = "f168" // f168:5日大单净额
	FieldMid5DayNet       = "f170" // f170:5日中单净额
	FieldSmall5DayNet     = "f172" // f172:5日小单净额
	FieldMain10DayNet     = "f174" // f174:10日主力净额
	FieldSuper10DayNet    = "f176" // f176:10日超大单净额
	FieldBig10DayNet      = "f178" // f178:10日大单净额
	FieldMid10DayNet      = "f180" // f180:10日中单净额
	FieldSmall10DayNet    = "f182" // f182:10日小单净额
	FieldBondBuyCode      = "f348" // f348:可转债申购代码
	FieldBondBuyDate      = "f243" // f243:可转债申购日期
	FieldLimitUpPrice     = "f350" // f350:涨停价
	FieldLimitDownPrice   = "f351" // f351:跌停价
	FieldAvgPrice         = "f352" // f352:均价
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

// IndustryInfo 行业信息
type IndustryInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// GetIndustries 获取行业列表
func (s *StockMarketService) GetIndustries() ([]IndustryInfo, error) {
	url := "https://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=200&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f3&fs=m:90+t:2&fields=f12,f14"

	resp, err := s.client.Get(url)
	if err != nil {
		logger.Error("请求行业列表失败", zap.Error(err))
		return nil, fmt.Errorf("请求行业列表失败: %w", err)
	}
	defer resp.Body.Close()

	// 定义专用响应结构体
	type IndustryListResponse struct {
		RC   int `json:"rc"`
		Data struct {
			Diff []struct {
				F12 string `json:"f12"` // 代码
				F14 string `json:"f14"` // 名称
			} `json:"diff"`
		} `json:"data"`
	}

	var apiResp IndustryListResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		logger.Error("解析行业列表失败", zap.Error(err))
		return nil, fmt.Errorf("解析行业列表失败: %w", err)
	}

	if apiResp.RC != 0 {
		return nil, fmt.Errorf("API返回错误: rc=%d", apiResp.RC)
	}

	industries := make([]IndustryInfo, 0)
	for _, item := range apiResp.Data.Diff {
		industries = append(industries, IndustryInfo{
			Code: item.F12,
			Name: item.F14,
		})
	}

	return industries, nil
}

// StockMarketData 市场股票数据
type StockMarketData struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Market   string `json:"market"`   // SH, SZ, BJ
	FullCode string `json:"fullCode"` // 市场代码 + 股票代码，如 SH600000
	Type     string `json:"type"`     // 主板, 创业板, 科创板, 北交所
	// 实时行情数据
	Price        float64 `json:"price"`        // 最新价
	ChangeRate   float64 `json:"changeRate"`   // 涨跌幅(%)
	ChangeAmount float64 `json:"changeAmount"` // 涨跌额
	Volume       float64 `json:"volume"`       // 成交量(手)
	Amount       float64 `json:"amount"`       // 成交额
	Amplitude    float64 `json:"amplitude"`    // 振幅(%)
	High         float64 `json:"high"`         // 最高价
	Low          float64 `json:"low"`          // 最低价
	Open         float64 `json:"open"`         // 开盘价
	PreClose     float64 `json:"preClose"`     // 昨收
	Turnover     float64 `json:"turnover"`     // 换手率(%)
	VolumeRatio  float64 `json:"volumeRatio"`  // 量比
	PE           float64 `json:"pe"`           // 市盈率
	WarrantRatio float64 `json:"warrantRatio"` // 委比(%)
	Industry     string  `json:"industry"`     // 所属行业
	Region       string  `json:"region"`       // 地区
	Board        string  `json:"board"`        // 板块
	TotalMV      float64 `json:"totalMV"`      // 总市值
	CircMV       float64 `json:"circMV"`       // 流通市值
	IsActive     int     `json:"isActive"`     // 是否在交易
	UpdatedAt    string  `json:"updatedAt"`    // 最后更新时间
}

// SyncStocksResult 同步结果
type SyncStocksResult struct {
	Total     int     `json:"total"`     // 总记录数
	Processed int     `json:"processed"` // 已处理记录数
	Inserted  int     `json:"inserted"`  // 新增记录数
	Updated   int     `json:"updated"`   // 更新记录数
	Duration  float64 `json:"duration"`  // 耗时（秒）
	Message   string  `json:"message"`   // 消息
}

// 东方财富API响应结构
type EastMoneyResponse struct {
	RC   int `json:"rc"`
	RT   int `json:"rt"`
	SVR  int `json:"svr"`
	LT   int `json:"lt"`
	Full int `json:"full"`
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
	pz := 5000 // 每页5000条
	fs := "m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23,m:0+t:81+s:2048"
	// 只请求数据库表需要的字段
	fields := "f12,f13,f14,f2,f3,f4,f5,f6,f7,f8,f9,f10,f15,f16,f17,f18,f33,f100,f102,f103,f20,f21"

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
			volume_ratio, pe, warrant_ratio, 
			industry, region, board, total_mv, circ_mv,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
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
			industry = excluded.industry,
			region = excluded.region,
			board = excluded.board,
			total_mv = excluded.total_mv,
			circ_mv = excluded.circ_mv,
			updated_at = CURRENT_TIMESTAMP
	`

	stmt, err := tx.Prepare(upsertSQL)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("准备upsert语句失败: %w", err)
	}
	defer stmt.Close()

	// 循环分页获取所有股票
	totalProcessed := 0
	totalInserted := 0
	totalUpdated := 0
	totalCount := 0

	for {
		url := fmt.Sprintf(
			"https://push2.eastmoney.com/api/qt/clist/get?pn=%d&pz=%d&fs=%s&fields=%s",
			pn, pz, fs, fields,
		)

		logger.Info("开始同步市场股票", zap.String("url", url), zap.Int("page", pn))

		// 请求接口
		resp, err := s.client.Get(url)
		if err != nil {
			logger.Error("请求东方财富API失败", zap.Error(err))
			_ = tx.Rollback()
			return nil, fmt.Errorf("请求东方财富API失败: %w", err)
		}
		defer resp.Body.Close()

		// 解析响应
		var apiResp EastMoneyResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			logger.Error("解析API响应失败", zap.Error(err))
			_ = tx.Rollback()
			return nil, fmt.Errorf("解析API响应失败: %w", err)
		}

		if apiResp.RC != 0 {
			_ = tx.Rollback()
			return nil, fmt.Errorf("API返回错误: rc=%d", apiResp.RC)
		}

		// 第一次请求时获取总数
		if totalCount == 0 {
			totalCount = apiResp.Data.Total
			logger.Info("获取到股票总数", zap.Int("total", totalCount))
		}

		// 如果没有数据，退出循环
		if len(apiResp.Data.Diff) == 0 {
			logger.Info("本页无数据，同步完成", zap.Int("page", pn))
			break
		}

		// 处理当前页的数据
		processed := 0
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
				stockData.FullCode,
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
				stockData.Industry,
				stockData.Region,
				stockData.Board,
				stockData.TotalMV,
				stockData.CircMV,
			)

			if err != nil {
				logger.Error("插入股票数据失败", zap.String("code", stockData.Code), zap.Error(err))
				continue
			}

			processed++
			// 第一条记录算插入，其他算更新
			if totalProcessed == 0 && processed == 1 {
				totalInserted++
			} else {
				totalUpdated++
			}
		}

		totalProcessed += processed
		logger.Info("本页数据处理完成",
			zap.Int("page", pn),
			zap.Int("processed", processed),
			zap.Int("totalProcessed", totalProcessed),
		)

		// 检查是否已经获取完所有数据
		if totalProcessed >= totalCount {
			logger.Info("已获取全部股票数据", zap.Int("total", totalCount))
			break
		}

		// 下一页
		pn++
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	duration := time.Since(startTime).Seconds()

	logger.Info("同步市场股票完成",
		zap.Int("total", totalCount),
		zap.Int("processed", totalProcessed),
		zap.Int("inserted", totalInserted),
		zap.Int("updated", totalUpdated),
		zap.Float64("duration", duration),
	)

	return &SyncStocksResult{
		Total:     totalCount,
		Processed: totalProcessed,
		Inserted:  totalInserted,
		Updated:   totalUpdated,
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

	// 获取code (f12: 股票代码)
	code, ok := data["f12"].(string)
	if !ok || code == "" {
		return nil
	}

	// 获取name (f14: 股票名称)
	name, ok := data["f14"].(string)
	if !ok {
		name = ""
	}

	// 获取market (f13: 市场)
	market := "SH"
	if marketCode, ok := data["f13"].(string); ok && marketCode != "" {
		market = marketCode
	} else {
		// 根据代码前缀判断市场（兜底逻辑）
		if strings.HasPrefix(code, "6") {
			market = "SH"
		} else if strings.HasPrefix(code, "0") || strings.HasPrefix(code, "3") {
			market = "SZ"
		} else if strings.HasPrefix(code, "8") || strings.HasPrefix(code, "4") {
			market = "BJ"
		}
	}

	// 判断板块类型
	stockType := "主板"
	if strings.HasPrefix(code, "688") {
		stockType = "科创板"
	} else if strings.HasPrefix(code, "6") {
		stockType = "主板"
	} else if strings.HasPrefix(code, "0") {
		stockType = "主板"
	} else if strings.HasPrefix(code, "3") {
		stockType = "创业板"
	} else if strings.HasPrefix(code, "8") || strings.HasPrefix(code, "4") {
		stockType = "北交所"
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

	// 解析数值字段（不需要除以100）
	parseInt := func(key string) float64 {
		if val, ok := data[key]; ok {
			switch v := val.(type) {
			case float64:
				return v
			case string:
				f, _ := strconv.ParseFloat(v, 64)
				return f
			}
		}
		return 0
	}

	// 解析字符串字段
	parseString := func(key string) string {
		if val, ok := data[key]; ok {
			if s, ok := val.(string); ok {
				return s
			}
		}
		return "-"
	}

	stock := &StockMarketData{
		Code:         code,
		Name:         name,
		Market:       market,
		FullCode:     market + code,
		Type:         stockType,
		IsActive:     1,
		Price:        parseFloat("f2"),    // f2: 最新价
		ChangeRate:   parseFloat("f3"),    // f3: 涨跌幅
		ChangeAmount: parseFloat("f4"),    // f4: 涨跌额
		Volume:       parseInt("f5"),      // f5: 总手（VOL）/成交量
		Amount:       parseInt("f6"),      // f6: 成交额
		Amplitude:    parseFloat("f7"),    // f7: 振幅
		High:         parseFloat("f15"),   // f15: 今日最高
		Low:          parseFloat("f16"),   // f16: 今日最低
		Open:         parseFloat("f17"),   // f17: 今开
		PreClose:     parseFloat("f18"),   // f18: 昨收
		Turnover:     parseFloat("f8"),    // f8: 换手率
		VolumeRatio:  parseFloat("f10"),   // f10: 量比
		PE:           parseFloat("f9"),    // f9: 市盈率(动态)
		WarrantRatio: parseFloat("f33"),   // f33: 委比
		Industry:     parseString("f100"), // f100: 行业
		Region:       parseString("f102"), // f102: 地区
		Board:        parseString("f103"), // f103: 板块/备注
		TotalMV:      parseInt("f20"),     // f20: 总市值
		CircMV:       parseInt("f21"),     // f21: 流通市值
		UpdatedAt:    updatedAt,
	}

	// 判断是否在交易（价格不为0表示在交易）
	if stock.Price == 0 {
		stock.IsActive = 0
	}

	return stock
}

// GetStocksList 获取股票列表（分页）
func (s *StockMarketService) GetStocksList(page int, pageSize int, search string, industry string) ([]StockMarketData, int, error) {
	offset := (page - 1) * pageSize

	db := s.dbService.GetDB()

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if search != "" {
		whereClause += " AND (code LIKE ? OR name LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	if industry != "" {
		whereClause += " AND industry = ?"
		args = append(args, industry)
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
		       volume_ratio, pe, warrant_ratio, 
		       COALESCE(industry, ''), COALESCE(region, ''), COALESCE(board, ''), COALESCE(total_mv, 0), COALESCE(circ_mv, 0),
		       updated_at
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
			&stock.Industry,
			&stock.Region,
			&stock.Board,
			&stock.TotalMV,
			&stock.CircMV,
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

// GetAllStockCodes 获取所有活跃股票代码
func (s *StockMarketService) GetAllStockCodes() ([]string, error) {
	db := s.dbService.GetDB()
	// 只查询当前在交易的股票
	rows, err := db.Query("SELECT code FROM stocks WHERE is_active = 1 ORDER BY code ASC")
	if err != nil {
		return nil, fmt.Errorf("查询股票代码失败: %w", err)
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			continue
		}
		codes = append(codes, code)
	}
	return codes, nil
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
		"lastUpdate":  lastUpdate,
		"marketStats": marketStats,
	}, nil
}
