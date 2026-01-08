package services

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

// StockService 股票数据服务
type StockService struct {
	ctx       context.Context // Wails Context
	exactURL  string
	listURL   string
	klineURL  string
	client    *http.Client
	sseClient *http.Client

	streamMu     sync.Mutex
	streams      map[string]context.CancelFunc
	emitIntraday func(ctx context.Context, code string, trends []string)
	dbService    *DBService // 数据库服务
}

// SetDBService 注入数据库服务。
//
// 为什么需要：StockService 的“历史数据同步/本地缓存”等功能依赖 SQLite。
// 如果不注入 dbService，会导致 SyncStockData/GetDataSyncStats 直接报“数据库服务未初始化”，
// 前端通常会把这类错误汇总展示为“自选股功能暂不可用”。
func (s *StockService) SetDBService(db *DBService) {
	s.dbService = db
}

// NewStockService 创建股票数据服务实例
func NewStockService() *StockService {
	s := &StockService{
		exactURL: "https://push2.eastmoney.com/api/qt/stock/get",
		listURL:  "http://78.push2.eastmoney.com/api/qt/clist/get",
		klineURL: "https://push2his.eastmoney.com/api/qt/stock/kline/get",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		// SSE 是长连接：不能用短 Timeout，否则读 body 会被强制 cancel 导致 context deadline exceeded
		sseClient: &http.Client{
			Timeout: 0,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				// 只限制“建连+拿到响应头”的时间，避免永不返回
				ResponseHeaderTimeout: 10 * time.Second,
			},
		},
		streams: make(map[string]context.CancelFunc),
	}
	// 默认事件推送实现（生产环境）
	s.emitIntraday = func(ctx context.Context, code string, trends []string) {
		runtime.EventsEmit(ctx, "intradayDataUpdate:"+code, trends)
	}
	return s
}

// Startup is called at application startup
func (s *StockService) Startup(ctx context.Context) {
	s.ctx = ctx
}

// GetStockDetail 获取个股详情页所需的所有数据
func (s *StockService) GetStockDetail(code string) (*models.StockDetail, error) {
	baseStock, err := s.GetStockByCode(code)
	if err != nil {
		return nil, err
	}

	detail := &models.StockDetail{
		StockData: *baseStock,
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		orderBook, err := s.getOrderBook(code)
		if err == nil {
			detail.OrderBook = *orderBook
		} else {
			logger.Error("获取盘口数据失败", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		financial, err := s.getFinancialSummary(code)
		if err == nil {
			detail.Financial = *financial
		} else {
			logger.Error("获取财务数据失败", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		industry, err := s.getIndustryInfo(code)
		if err == nil {
			detail.Industry = *industry
		} else {
			logger.Error("获取行业数据失败", zap.Error(err))
		}
	}()

	wg.Wait()

	return detail, nil
}

// getOrderBook 获取五档盘口数据
func (s *StockService) getOrderBook(code string) (*models.OrderBook, error) {
	secid := s.getSecID(code)
	if secid == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	fields := "f68,f69,f70,f71,f72,f73,f74,f75,f76,f77,f78,f79,f80,f81,f82,f83,f84,f85"
	fullURL := fmt.Sprintf("%s?secid=%s&fields=%s", s.exactURL, secid, fields)

	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求盘口数据失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析盘口响应失败: %w", err)
	}

	if result.Data == nil {
		return nil, fmt.Errorf("未找到盘口数据: %s", code)
	}

	data := result.Data

	orderBook := &models.OrderBook{
		Buy5:  make([]models.OrderBookEntry, 5),
		Sell5: make([]models.OrderBookEntry, 5),
	}

	for i := 0; i < 5; i++ {
		priceField := fmt.Sprintf("f%d", 73+i*2)
		volumeField := fmt.Sprintf("f%d", 74+i*2)
		orderBook.Sell5[i] = models.OrderBookEntry{
			Price:  getFloat(data[priceField]) / 100,
			Volume: getInt64(data[volumeField]),
		}
	}

	for i := 0; i < 5; i++ {
		priceField := fmt.Sprintf("f%d", 68+i*2)
		volumeField := fmt.Sprintf("f%d", 69+i*2)
		orderBook.Buy5[i] = models.OrderBookEntry{
			Price:  getFloat(data[priceField]) / 100,
			Volume: getInt64(data[volumeField]),
		}
	}

	return orderBook, nil
}

// getFinancialSummary 获取核心财务数据 (Mock 数据)
func (s *StockService) getFinancialSummary(code string) (*models.FinancialSummary, error) {
	_ = code // mock data: keep signature for future real implementation
	// MOCK DATA
	return &models.FinancialSummary{
		ROE:                    15.8,
		NetProfitGrowthRate:    22.5,
		GrossProfitMargin:      45.1,
		TotalMarketValue:       12000.5, // 亿元
		CirculatingMarketValue: 8000.2,  // 亿元
		DividendYield:          1.5,
		ReportDate:             time.Now().AddDate(0, -3, 0),
	}, nil
}

// getIndustryInfo 获取行业与宏观信息 (Mock 数据)
func (s *StockService) getIndustryInfo(code string) (*models.IndustryInfo, error) {
	_ = code // mock data: keep signature for future real implementation
	// MOCK DATA
	return &models.IndustryInfo{
		IndustryName: "软件开发",
		ConceptNames: []string{"人工智能", "云计算", "数字经济"},
		IndustryPE:   45.8,
	}, nil
}

// GetStockByCode 根据股票代码获取股票数据（精确查询）
func (s *StockService) GetStockByCode(code string) (*models.StockData, error) {
	start := time.Now()
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}

	secid := s.getSecID(code)
	if secid == "" {
		return nil, fmt.Errorf("无法识别的股票代码格式: %s", code)
	}

	fields := "f58,f43,f169,f170,f47,f48,f44,f45,f46,f60,f171,f168,f162,f167,f116,f117,f12,f14,f19,f20"
	fullURL := fmt.Sprintf("%s?secid=%s&fields=%s", s.exactURL, secid, fields)

	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Data == nil {
		return nil, fmt.Errorf("未找到股票数据: %s", code)
	}

	data := result.Data
	stockCode := getString(data["f12"])
	if stockCode == "" {
		stockCode = code
	}

	stock := &models.StockData{
		Code:         stockCode,
		Name:         getString(data["f58"]),
		Price:        getFloat(data["f43"]) / 100,
		Change:       getFloat(data["f169"]) / 100,
		ChangeRate:   getFloat(data["f170"]) / 100,
		Volume:       getInt64(data["f47"]),
		Amount:       getFloat(data["f48"]),
		High:         getFloat(data["f44"]) / 100,
		Low:          getFloat(data["f45"]) / 100,
		Open:         getFloat(data["f46"]) / 100,
		PreClose:     getFloat(data["f60"]) / 100,
		Amplitude:    getFloat(data["f171"]) / 100,
		Turnover:     getFloat(data["f168"]) / 100,
		PE:           getFloat(data["f162"]) / 100,
		PB:           getFloat(data["f167"]) / 100,
		TotalMV:      getFloat(data["f116"]),
		CircMV:       getFloat(data["f117"]),
		VolumeRatio:  getFloat(data["f20"]) / 100,
		WarrantRatio: getFloat(data["f19"]) / 100,
	}

	logger.Info("精确获取股票数据成功", zap.String("code", code), zap.Int64("ms", time.Since(start).Milliseconds()))
	return stock, nil
}

// GetKLineData 获取历史K线数据并计算技术指标
func (s *StockService) GetKLineData(code string, limit int, period string) ([]*models.KLineData, error) {
	secid := s.getSecID(code)
	if secid == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	klt := "101"
	switch period {
	case "week":
		klt = "102"
	case "month":
		klt = "103"
	}

	fetchLimit := limit + 50
	url := fmt.Sprintf("%s?secid=%s&fields1=f1,f2,f3,f4,f5,f6&fields2=f51,f52,f53,f54,f55,f56&klt=%s&fqt=1&end=20500101&lmt=%d", s.klineURL, secid, klt, fetchLimit)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data struct {
			Klines []string `json:"klines"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	klines := make([]*models.KLineData, 0)
	for _, line := range result.Data.Klines {
		parts := strings.Split(line, ",")
		if len(parts) < 6 {
			continue
		}
		klines = append(klines, &models.KLineData{
			Time:   parts[0],
			Open:   s.parsePrice(parts[1]),
			Close:  s.parsePrice(parts[2]),
			High:   s.parsePrice(parts[3]),
			Low:    s.parsePrice(parts[4]),
			Volume: int64(s.parsePrice(parts[5])),
		})
	}

	s.calculateIndicators(klines)

	if len(klines) > limit {
		return klines[len(klines)-limit:], nil
	}
	return klines, nil
}

func (s *StockService) calculateIndicators(klines []*models.KLineData) {
	if len(klines) == 0 {
		return
	}

	ema12 := klines[0].Close
	ema26 := klines[0].Close
	dea := 0.0
	for i, k := range klines {
		ema12 = ema12*11/13 + k.Close*2/13
		ema26 = ema26*25/27 + k.Close*2/27
		dif := ema12 - ema26
		if i == 0 {
			dea = dif
		} else {
			dea = dea*8/10 + dif*2/10
		}
		k.MACD = &models.MACD{
			DIF: dif,
			DEA: dea,
			Bar: (dif - dea) * 2,
		}
	}

	for i := 0; i < len(klines); i++ {
		if i < 8 {
			klines[i].KDJ = &models.KDJ{K: 50, D: 50, J: 50}
			continue
		}
		low := klines[i].Low
		high := klines[i].High
		for j := i - 8; j < i; j++ {
			low = math.Min(low, klines[j].Low)
			high = math.Max(high, klines[j].High)
		}
		rsv := 0.0
		if high != low {
			rsv = (klines[i].Close - low) / (high - low) * 100
		}
		prevK := klines[i-1].KDJ.K
		prevD := klines[i-1].KDJ.D
		k := prevK*2/3 + rsv/3
		d := prevD*2/3 + k/3
		klines[i].KDJ = &models.KDJ{
			K: k,
			D: d,
			J: 3*k - 2*d,
		}
	}

	if len(klines) > 14 {
		for i := 14; i < len(klines); i++ {
			upSum := 0.0
			downSum := 0.0
			for j := i - 13; j <= i; j++ {
				diff := klines[j].Close - klines[j-1].Close
				if diff > 0 {
					upSum += diff
				} else {
					downSum -= diff
				}
			}
			if upSum+downSum == 0 {
				klines[i].RSI = 50
			} else {
				klines[i].RSI = upSum / (upSum + downSum) * 100
			}
		}
	}
}

// GetIntradayData 获取分时数据 (非实时快照，用于初始化)
func (s *StockService) GetIntradayData(code string) (*models.IntradayResponse, error) {
	secid := s.getSecID(code)
	if secid == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	url := fmt.Sprintf("https://push2.eastmoney.com/api/qt/stock/trends2/get?secid=%s&fields1=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13&fields2=f51,f52,f53,f54,f55,f56,f57,f58&ndays=1&iscr=0&ut=fa5fd1943c7b386f172d6893dbfba10b", secid)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data struct {
			PreClose float64  `json:"preClose"`
			Trends   []string `json:"trends"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	intradayData := make([]models.IntradayData, 0)
	for _, line := range result.Data.Trends {
		parts := strings.Split(line, ",")
		if len(parts) < 8 {
			continue
		}
		intradayData = append(intradayData, models.IntradayData{
			Time:     parts[0],
			Price:    s.parsePrice(parts[2]),
			AvgPrice: s.parsePrice(parts[7]),
			Volume:   int64(s.parsePrice(parts[5])),
			PreClose: result.Data.PreClose,
		})
	}

	return &models.IntradayResponse{
		Data:     intradayData,
		PreClose: result.Data.PreClose,
	}, nil
}

// StopIntradayStream 停止指定股票的分时 SSE 流。
//
// 设计说明：
// - 前端可能会重复调用 StreamIntradayData（例如切换 Tab/刷新页面/重复进入详情页）。
// - SSE 是长连接，如果不显式 stop，会造成 goroutine 泄漏、重复推送、网络连接占用。
// - 这里按股票 code 维度管理 cancel func：同一 code 只允许存在一个活跃 SSE 流。
func (s *StockService) StopIntradayStream(code string) {
	code = strings.TrimSpace(code)
	if code == "" {
		return
	}
	s.streamMu.Lock()
	cancel := s.streams[code]
	delete(s.streams, code)
	s.streamMu.Unlock()
	if cancel != nil {
		cancel()
	}
}

// StreamIntradayData 实时流式获取分时数据 (SSE 代理)。
//
// 核心特性：
// 1) 使用专用 sseClient（Timeout=0）避免默认 10s 超时导致的 "context deadline exceeded"。
// 2) 同一 code 重复启动会先 Stop 旧流，避免重复 goroutine & 重复推送。
// 3) 自动重连：连接失败/非 200/读流错误/EOF 会进入重试（指数退避 + 抖动）。
func (s *StockService) StreamIntradayData(code string) {
	code = strings.TrimSpace(code)
	if code == "" {
		return
	}

	// 如果同一个 code 已经在流式推送，先停掉旧的，避免重复推送/泄漏
	s.StopIntradayStream(code)

	// 子 context：用于单个 code 的生命周期控制（StopIntradayStream 会 cancel）
	sseCtx, cancel := context.WithCancel(s.ctx)
	s.streamMu.Lock()
	s.streams[code] = cancel
	s.streamMu.Unlock()

	go func() {
		defer func() {
			// goroutine 退出时清理 map，避免残留
			s.streamMu.Lock()
			delete(s.streams, code)
			s.streamMu.Unlock()
		}()

		secid := s.getSecID(code)
		if secid == "" {
			logger.Error("无效的股票代码（无法推送 SSE）",
				zap.String("module", "services.stock"),
				zap.String("op", "StreamIntradayData"),
				zap.String("code", code),
			)
			return
		}

		sseURL := fmt.Sprintf("https://42.push2.eastmoney.com/api/qt/stock/trends2/sse?secid=%s&fields1=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13&fields2=f51,f52,f53,f54,f55,f56,f57,f58&ndays=1&iscr=0&ut=fa5fd1943c7b386f172d6893dbfba10b", secid)

		retry := 0
		for {
			select {
			case <-sseCtx.Done():
				logger.Info("SSE 流已停止",
					zap.String("module", "services.stock"),
					zap.String("op", "StreamIntradayData"),
					zap.String("code", code),
					zap.Error(sseCtx.Err()),
				)
				return
			default:
			}

			attemptStart := time.Now()
			req, err := http.NewRequestWithContext(sseCtx, "GET", sseURL, nil)
			if err != nil {
				logger.Error("创建 SSE 请求失败",
					zap.String("module", "services.stock"),
					zap.String("op", "StreamIntradayData"),
					zap.String("code", code),
					zap.String("url", sseURL),
					zap.Error(err),
				)
				return
			}
			req.Header.Set("Accept", "text/event-stream")
			req.Header.Set("Cache-Control", "no-cache")
			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

			resp, err := s.sseClient.Do(req)
			if err != nil {
				// 取消（页面离开、应用退出等）属于正常退出，不需要报错
				if sseCtx.Err() != nil {
					logger.Info("SSE 连接结束（已取消）",
						zap.String("module", "services.stock"),
						zap.String("op", "StreamIntradayData"),
						zap.String("code", code),
						zap.Error(sseCtx.Err()),
						zap.Int64("duration_ms", time.Since(attemptStart).Milliseconds()),
					)
					return
				}
				logger.Warn("连接 SSE 接口失败，准备重试",
					zap.String("module", "services.stock"),
					zap.String("op", "StreamIntradayData"),
					zap.String("code", code),
					zap.String("url", sseURL),
					zap.Int("retry", retry),
					zap.Error(err),
					zap.Int64("duration_ms", time.Since(attemptStart).Milliseconds()),
				)
				if !s.sleepBackoff(sseCtx, retry) {
					return
				}
				retry++
				continue
			}

			func() {
				defer func() { _ = resp.Body.Close() }()

				if resp.StatusCode != http.StatusOK {
					// 尽量读一点 body 帮助定位（例如被限流、被 WAF 拦截、返回错误 JSON 等）
					b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
					logger.Warn("SSE 接口返回非 200，准备重试",
						zap.String("module", "services.stock"),
						zap.String("op", "StreamIntradayData"),
						zap.String("code", code),
						zap.String("url", sseURL),
						zap.Int("status", resp.StatusCode),
						zap.ByteString("body", b),
						zap.Int("retry", retry),
						zap.Int64("duration_ms", time.Since(attemptStart).Milliseconds()),
					)
					return
				}

				scanner := bufio.NewScanner(resp.Body)
				// SSE 的单行 data 可能超过默认 64KB，必须放大 buffer，否则会 ErrTooLong
				buf := make([]byte, 0, 64*1024)
				scanner.Buffer(buf, 2*1024*1024)

				for scanner.Scan() {
					select {
					case <-sseCtx.Done():
						logger.Info("SSE 流被取消",
							zap.String("module", "services.stock"),
							zap.String("op", "StreamIntradayData"),
							zap.String("code", code),
							zap.Int64("duration_ms", time.Since(attemptStart).Milliseconds()),
						)
						return
					default:
					}

					line := scanner.Text()
					// SSE 允许空行/注释/心跳行（例如 ": keep-alive"）
					if line == "" || strings.HasPrefix(line, ":") {
						continue
					}
					if strings.HasPrefix(line, "data:") {
						jsonStr := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
						var sseData struct {
							Data struct {
								Trends []string `json:"trends"`
							} `json:"data"`
						}
						if err := json.Unmarshal([]byte(jsonStr), &sseData); err != nil {
							// 解析失败通常是服务端偶发推送了非预期格式；这里用 debug 级别避免刷屏。
							logger.Debug("SSE data JSON 解析失败（将忽略该条）",
								zap.String("module", "services.stock"),
								zap.String("op", "StreamIntradayData"),
								zap.String("code", code),
								zap.Int("json_len", len(jsonStr)),
								zap.Error(err),
							)
							continue
						}
						if sseData.Data.Trends != nil {
							s.emitIntraday(s.ctx, code, sseData.Data.Trends)
						}
					}
				}

				if err := scanner.Err(); err != nil && err != io.EOF {
					if sseCtx.Err() != nil {
						logger.Info("SSE 读取结束（已取消）",
							zap.String("module", "services.stock"),
							zap.String("op", "StreamIntradayData"),
							zap.String("code", code),
							zap.Error(sseCtx.Err()),
							zap.Int64("duration_ms", time.Since(attemptStart).Milliseconds()),
						)
						return
					}
					logger.Warn("读取 SSE 流发生错误，准备重试",
						zap.String("module", "services.stock"),
						zap.String("op", "StreamIntradayData"),
						zap.String("code", code),
						zap.String("url", sseURL),
						zap.Int("retry", retry),
						zap.Error(err),
						zap.Int64("duration_ms", time.Since(attemptStart).Milliseconds()),
					)
					return
				}
			}()

			// 正常结束（EOF/非200）也进入重试，除非已取消
			if sseCtx.Err() != nil {
				return
			}
			if !s.sleepBackoff(sseCtx, retry) {
				return
			}
			retry++
		}
	}()
}

// sleepBackoff sleeps with exponential backoff + jitter. Returns false if ctx cancelled.
func (s *StockService) sleepBackoff(ctx context.Context, retry int) bool {
	// 500ms, 1s, 2s, 4s, 8s... capped at 15s
	base := 500 * time.Millisecond
	d := base * time.Duration(1<<minInt(retry, 5))
	if d > 15*time.Second {
		d = 15 * time.Second
	}
	// jitter: 0~250ms
	d += time.Duration(rand.Intn(250)) * time.Millisecond

	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetStockHealthCheck 获取股票健康状况
func (s *StockService) GetStockHealthCheck(code string) (*models.HealthCheckResult, error) {
	stock, err := s.GetStockByCode(code)
	if err != nil {
		return nil, err
	}

	items := make([]models.HealthItem, 0)
	score := 100

	peStatus := "正常"
	peDesc := fmt.Sprintf("当前市盈率 %.2f，处于行业合理区间。", stock.PE)
	if stock.PE < 0 {
		peStatus = "异常"
		peDesc = "公司目前处于亏损状态，财务基本面存在较大不确定性。"
		score -= 20
	} else if stock.PE > 100 {
		peStatus = "警告"
		peDesc = "市盈率过高，估值存在泡沫风险，需警惕回调。"
		score -= 10
	}
	items = append(items, models.HealthItem{
		Category:    "财务",
		Name:        "估值水平",
		Value:       fmt.Sprintf("%.2f PE", stock.PE),
		Status:      peStatus,
		Description: peDesc,
	})

	turnoverStatus := "正常"
	turnoverDesc := "换手率适中，交投活跃度正常。"
	if stock.Turnover > 15 {
		turnoverStatus = "警告"
		turnoverDesc = "换手率极高，说明筹码松动，主力可能在进行高位派发。"
		score -= 15
	} else if stock.Turnover < 0.5 {
		turnoverStatus = "警告"
		turnoverDesc = "成交极其低迷，属于“僵尸股”，流动性风险较大。"
		score -= 10
	}
	items = append(items, models.HealthItem{
		Category:    "资金",
		Name:        "流动性检测",
		Value:       fmt.Sprintf("%.2f%%", stock.Turnover),
		Status:      turnoverStatus,
		Description: turnoverDesc,
	})

	ampStatus := "正常"
	ampDesc := "日内波动在正常范围内。"
	if stock.Amplitude > 10 {
		ampStatus = "警告"
		ampDesc = "日内振幅巨大，多空分歧严重，短期波动风险极高。"
		score -= 10
	}
	items = append(items, models.HealthItem{
		Category:    "技术",
		Name:        "波动风险",
		Value:       fmt.Sprintf("%.2f%%", stock.Amplitude),
		Status:      ampStatus,
		Description: ampDesc,
	})

	status, riskLevel := "健康", "低"
	if score < 60 {
		status, riskLevel = "风险", "高"
	} else if score < 85 {
		status, riskLevel = "亚健康", "中"
	}

	return &models.HealthCheckResult{
		Score:     score,
		Status:    status,
		RiskLevel: riskLevel,
		Summary:   fmt.Sprintf("%s目前综合评分为 %d 分，%s", stock.Name, score, getHealthSummary(score)),
		Items:     items,
		UpdatedAt: time.Now().Format("15:04:05"),
	}, nil
}

func getHealthSummary(score int) string {
	if score >= 85 {
		return "基本面和技术面表现稳健，适合中长期关注。"
	} else if score >= 60 {
		return "存在部分指标异常，建议控制仓位，关注关键支撑位。"
	}
	return "多项指标触发预警，建议新手坚决回避，等待风险释放。"
}

func (s *StockService) BatchAnalyzeStocks(ctx context.Context, codes []string, role string, aiSvc *AIService) error {
	total := len(codes)
	var completed int32 = 0
	concurrency := 2
	jobs := make(chan string, total)
	var wg sync.WaitGroup

	for w := 1; w <= concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for code := range jobs {
				stock, err := s.GetStockByCode(code)
				if err != nil {
					continue
				}
				klines, err := s.GetKLineData(code, 100, "daily")
				if err != nil {
					continue
				}
				_, _ = aiSvc.AnalyzeTechnical(stock, klines, "daily", role)
				newCompleted := atomic.AddInt32(&completed, 1)
				runtime.EventsEmit(ctx, "batch_analyze_progress", map[string]interface{}{
					"code":      code,
					"name":      stock.Name,
					"completed": newCompleted,
					"total":     total,
					"percent":   float64(newCompleted) / float64(total) * 100,
				})
			}
		}()
	}

	for _, code := range codes {
		jobs <- code
	}
	close(jobs)
	wg.Wait()
	return nil
}

func getString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getFloat(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

func getInt64(v interface{}) int64 {
	if f, ok := v.(float64); ok {
		return int64(f)
	}
	return 0
}

func (s *StockService) parsePrice(p string) float64 {
	if val, err := strconv.ParseFloat(p, 64); err == nil {
		return val
	}
	return 0
}

func (s *StockService) getSecID(code string) string {
	if len(code) != 6 {
		return ""
	}
	if code[0] == '6' {
		return "1." + code
	} else if code[0] == '0' || code[0] == '3' {
		return "0." + code
	}
	return ""
}

// GetMoneyFlowData 获取资金流数据
func (s *StockService) GetMoneyFlowData(code string) (*models.MoneyFlowResponse, error) {
	secid := s.getSecID(code)
	if secid == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	// 使用东方财富的资金流接口
	url := fmt.Sprintf("http://push2.eastmoney.com/api/qt/stock/fflow/daykline/get?secid=%s&fields1=f1,f2,f3,f7&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65&lmt=0&klt=101&fqt=1", secid)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data struct {
			Flows []string `json:"klines"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	flowData := make([]models.MoneyFlowData, 0)
	for _, line := range result.Data.Flows {
		parts := strings.Split(line, ",")
		if len(parts) < 11 {
			continue
		}
		flowData = append(flowData, models.MoneyFlowData{
			Time:   parts[0],
			Main:   s.parsePrice(parts[1]),
			Retail: s.parsePrice(parts[4]),
			Super:  s.parsePrice(parts[7]),
			Big:    s.parsePrice(parts[8]),
			Medium: s.parsePrice(parts[9]),
			Small:  s.parsePrice(parts[10]),
		})
	}

	return &models.MoneyFlowResponse{
		Data: flowData,
	}, nil
}

// SearchStock 搜索股票
func (s *StockService) SearchStock(keyword string) ([]*models.StockData, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []*models.StockData{}, nil
	}

	// 兼容：如果输入是 6 位代码，优先走精确查询
	if len(keyword) == 6 {
		stock, err := s.GetStockByCode(keyword)
		if err == nil {
			return []*models.StockData{stock}, nil
		}
	}

	return s.SearchStockLegacy(keyword)
}

// SearchStockLegacy 基于列表接口的模糊搜索
func (s *StockService) SearchStockLegacy(keyword string) ([]*models.StockData, error) {
	url := fmt.Sprintf("%s?pn=1&pz=1000&po=1&np=1&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:13,m:0+t:80,m:1+t:2,m:1+t:23&fields=f12,f14", s.listURL)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	var apiResp models.EastMoneyResponse
	_ = json.Unmarshal(body, &apiResp)

	kw := strings.ToLower(strings.TrimSpace(keyword))
	results := make([]*models.StockData, 0)
	for _, diff := range apiResp.Data.Diff {
		if strings.Contains(strings.ToLower(diff.F12), kw) ||
			strings.Contains(strings.ToLower(diff.F14), kw) {
			results = append(results, diff.ToStockData())
		}
	}
	return results, nil
}

// SyncStockData 同步单个股票的历史数据到本地 SQLite
// 该方法会为每个股票创建一个独立的表（如 kline_600519），并存储历史 K 线数据
func (s *StockService) SyncStockData(code string, startDate string, endDate string) (*models.SyncResult, error) {
	result := &models.SyncResult{
		StockCode: code,
		Success:   false,
	}

	// 获取数据库服务实例
	db := s.dbService
	if db == nil {
		// 这里必须打日志：否则用户只看到“自选股暂不可用”，不知道实际原因是 dbService 未注入。
		err := fmt.Errorf("数据库服务未初始化")
		result.ErrorMessage = err.Error()
		logger.Error("同步历史数据失败：数据库服务未初始化（可能未注入或初始化失败）",
			zap.String("module", "services.stock"),
			zap.String("op", "SyncStockData"),
			zap.String("code", code),
			zap.String("startDate", startDate),
			zap.String("endDate", endDate),
			zap.Error(err),
		)
		return result, err
	}

	// 1. 获取股票的历史 K 线数据
	klines, err := s.GetKLineData(code, 5000, "daily")
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("获取 K 线数据失败: %v", err)
		return result, err
	}

	if len(klines) == 0 {
		result.ErrorMessage = "未获取到 K 线数据"
		return result, fmt.Errorf("未获取到 K 线数据")
	}

	// 2. 创建数据库表
	if err := db.CreateKLineCacheTable(code); err != nil {
		result.ErrorMessage = fmt.Sprintf("创建数据表失败: %v", err)
		return result, err
	}

	// 3. 转换 K 线数据格式并过滤日期范围
	var klineRecords []map[string]interface{}
	for _, kline := range klines {
		// 过滤日期范围
		if kline.Time >= startDate && kline.Time <= endDate {
			klineRecords = append(klineRecords, map[string]interface{}{
				"date":   kline.Time,
				"open":   kline.Open,
				"high":   kline.High,
				"low":    kline.Low,
				"close":  kline.Close,
				"volume": kline.Volume,
			})
		}
	}

	if len(klineRecords) == 0 {
		result.ErrorMessage = "指定日期范围内没有数据"
		return result, fmt.Errorf("指定日期范围内没有数据")
	}

	// 4. 批量插入或更新数据
	addedCount, updatedCount, err := db.InsertOrUpdateKLineData(code, klineRecords)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("插入数据失败: %v", err)
		return result, err
	}

	// 5. 返回成功结果
	result.Success = true
	result.RecordsAdded = int(addedCount)
	result.RecordsUpdated = int(updatedCount)
	result.Message = fmt.Sprintf("成功同步 %s 的历史数据，新增 %d 条记录，更新 %d 条记录", code, addedCount, updatedCount)

	return result, nil
}

// GetDataSyncStats 获取数据同步统计信息
// 返回已同步的股票列表、总记录数等信息
func (s *StockService) GetDataSyncStats() (*models.DataSyncStats, error) {
	stats := &models.DataSyncStats{
		TotalStocks:  0,
		SyncedStocks: 0,
		TotalRecords: 0,
		StockList:    []string{},
		LastSyncTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 获取数据库服务实例
	db := s.dbService
	if db == nil {
		// 统计接口更适合降级：返回空数据，但保留 warn 日志用于排查。
		logger.Warn("获取同步统计失败：数据库服务未初始化（返回空统计）",
			zap.String("module", "services.stock"),
			zap.String("op", "GetDataSyncStats"),
		)
		return stats, nil
	}

	// 1. 查询所有已同步的股票
	stocks, err := db.GetAllSyncedStocks()
	if err != nil {
		return stats, err
	}

	stats.StockList = stocks
	stats.SyncedStocks = len(stocks)
	stats.TotalStocks = len(stocks)

	// 2. 统计总记录数
	var totalRecords int64
	for _, code := range stocks {
		count, err := db.GetKLineCountByCode(code)
		if err == nil {
			totalRecords += int64(count)
		}
	}

	stats.TotalRecords = totalRecords

	return stats, nil
}

// GetKLineFromCache 从本地缓存获取 K 线数据
// 优先从 SQLite 读取，如果数据不足则从 API 补充
func (s *StockService) GetKLineFromCache(code string, limit int) ([]*models.KLineData, error) {
	// TODO: 实现缓存读取逻辑
	// 1. 从数据库表 kline_{code} 中读取最新的 limit 条记录
	// 2. 如果记录数不足，从 API 获取补充数据
	// 3. 将补充的数据存储到数据库
	// 4. 返回合并后的 K 线数据

	// 当前直接调用 API（后续优化为缓存优先）
	return s.GetKLineData(code, limit, "daily")
}

// ClearStockCache 清除指定股票的本地缓存数据
func (s *StockService) ClearStockCache(code string) error {
	// TODO: 实现缓存清除逻辑
	// 1. 删除数据库表 kline_{code}
	// 2. 返回删除结果

	return nil
}

// BatchSyncStockData 批量同步多个股票的历史数据
// 该方法会为每个股票创建独立的表，并通过 Wails 事件发送同步进度
func (s *StockService) BatchSyncStockData(codes []string, startDate string, endDate string) error {
	if len(codes) == 0 {
		return fmt.Errorf("股票代码列表为空")
	}

	logger.Info("开始批量同步股票数据",
		zap.String("module", "services.stock"),
		zap.String("op", "BatchSyncStockData"),
		zap.Int("stock_count", len(codes)),
		zap.String("start_date", startDate),
		zap.String("end_date", endDate),
	)

	startTime := time.Now()
	totalAdded := 0
	totalUpdated := 0
	failedCodes := []string{}
	successCodes := []string{}

	// 遍历 codes 列表
	for i, code := range codes {
		// 调用 SyncStockData
		result, err := s.SyncStockData(code, startDate, endDate)

		// 发送进度事件
		if s.ctx != nil {
			runtime.EventsEmit(s.ctx, "dataSyncProgress", map[string]interface{}{
				"currentIndex": i + 1,
				"totalCount":   len(codes),
				"currentCode":  code,
				"stockName":    result.StockCode, // 这里可以用实际股票名称
				"success":      result.Success,
				"message":      result.Message,
			})
		}

		if err != nil || !result.Success {
			failedCodes = append(failedCodes, code)
			logger.Error("同步单个股票数据失败",
				zap.String("code", code),
				zap.String("error", result.ErrorMessage),
			)
		} else {
			totalAdded += result.RecordsAdded
			totalUpdated += result.RecordsUpdated
			successCodes = append(successCodes, code)
		}

		// 避免 API 限流，短暂延迟
		if i < len(codes)-1 {
			time.Sleep(200 * time.Millisecond)
		}
	}

	duration := int(time.Since(startTime).Seconds())

	logger.Info("批量同步完成",
		zap.Int("success_count", len(successCodes)),
		zap.Int("failed_count", len(failedCodes)),
		zap.Int("total_added", totalAdded),
		zap.Int("total_updated", totalUpdated),
		zap.Int("duration", duration),
	)

	// 如果所有都失败了，返回错误
	if len(successCodes) == 0 {
		return fmt.Errorf("批量同步失败，所有股票同步都失败了")
	}

	return nil
}
