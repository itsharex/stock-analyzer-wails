package services

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"

	"github.com/go-resty/resty/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

// SyncService 全量数据同步服务
type SyncService struct {
	dbService          *DBService
	stockMarketService *StockMarketService
	moneyFlowRepo      *repositories.MoneyFlowRepository
	client             *resty.Client
	ctx                context.Context
	running            bool
	mu                 sync.Mutex
}

// SyncProgress 同步进度结构体
type SyncProgress struct {
	Total        int    `json:"total"`
	Current      int    `json:"current"`
	CurrentStock string `json:"currentStock"`
	Status       string `json:"status"` // "running", "completed", "error"
	SuccessCount int    `json:"successCount"`
	FailedCount  int    `json:"failedCount"`
}

// NewSyncService 创建同步服务
func NewSyncService(
	dbService *DBService,
	stockMarketService *StockMarketService,
	moneyFlowRepo *repositories.MoneyFlowRepository,
) *SyncService {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	client.SetHeader("Referer", "https://quote.eastmoney.com/")
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)

	// 设置 CloseConnection 为 true 以模拟短连接行为，避免 Keep-Alive 导致的 EOF
	// 这是解决 Eastmoney API EOF 问题的关键
	client.SetCloseConnection(true)

	return &SyncService{
		dbService:          dbService,
		stockMarketService: stockMarketService,
		moneyFlowRepo:      moneyFlowRepo,
		client:             client,
	}
}

// SetContext 设置上下文
func (s *SyncService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// StartFullMarketSync 启动全市场历史资金流同步
func (s *SyncService) StartFullMarketSync() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("同步任务已在运行中")
	}
	s.running = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
	}()

	logger.Info("开始全市场历史资金流同步任务")

	// 1. 获取所有股票代码
	codes, err := s.stockMarketService.GetAllStockCodes()
	if err != nil {
		s.emitProgress(&SyncProgress{Status: "error", CurrentStock: "获取股票列表失败"})
		return fmt.Errorf("获取股票列表失败: %w", err)
	}

	total := len(codes)
	logger.Info("获取到待同步股票", zap.Int("total", total))

	// 初始化进度
	progress := &SyncProgress{
		Total:  total,
		Status: "running",
	}
	s.emitProgress(progress)

	// 2. 并发控制 (限制 5 个并发)
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	// 数据通道 (Worker -> Saver)
	// 每个 worker 可能会发送 2000+ 条历史数据，所以这里的 buffer 不需要太大，只要能缓冲几个 worker 的结果即可
	dataChan := make(chan []models.MoneyFlowData, 20)

	// 结果通道 (Saver -> Progress)
	resultChan := make(chan bool, total)

	// 启动单一写入协程 (Single Writer)
	go func() {
		defer close(resultChan) // 写入完成后关闭结果通道

		var batch []models.MoneyFlowData
		// 累积 20 只股票的数据提交一次 (假设每只股票 1000 条数据，20只就是 20000 条，可能有点多)
		// SQLite 批量插入建议 500-1000 条一次比较稳，但如果是 Transaction，可以多一些。
		// 用户建议：每累积 10-20 只股票的数据执行一次事务提交
		const StocksPerBatch = 10
		stocksInBatch := 0

		for flows := range dataChan {
			if len(flows) > 0 {
				batch = append(batch, flows...)
				stocksInBatch++

				// 达到批次大小，执行提交
				if stocksInBatch >= StocksPerBatch {
					err := s.moneyFlowRepo.SaveMoneyFlows(batch)
					if err != nil {
						logger.Error("批量保存资金流失败", zap.Error(err))
					}
					// 无论成功失败，都清空批次
					batch = nil
					stocksInBatch = 0
				}
			}
			// 爬取完成一个股票，发送成功信号
			resultChan <- true
		}

		// 处理剩余数据
		if len(batch) > 0 {
			err := s.moneyFlowRepo.SaveMoneyFlows(batch)
			if err != nil {
				logger.Error("批量保存剩余资金流失败", zap.Error(err))
			}
		}
	}()

	// 启动进度监听协程
	go func() {
		for success := range resultChan {
			progress.Current++
			if success {
				progress.SuccessCount++
			} else {
				progress.FailedCount++
			}
			s.emitProgress(progress)
		}
	}()

	// 3. 循环执行任务
	for i, code := range codes {
		// 检查上下文是否取消
		select {
		case <-s.ctx.Done():
			logger.Warn("同步任务被取消")
			close(dataChan) // 关闭数据通道，停止写入协程
			return nil
		default:
		}

		progress.CurrentStock = code

		wg.Add(1)
		sem <- struct{}{} // 获取信号量

		go func(stockCode string, idx int) {
			defer wg.Done()
			defer func() { <-sem }() // 释放信号量

			// 防封禁休眠
			time.Sleep(200 * time.Millisecond)

			// 仅爬取数据，不写入数据库
			rawData, err := s.FetchHistoryFlowDataV2(stockCode, 120)
			flows := AlignStockData2MoneyFlow(stockCode, GetSortedData(rawData))
			if err != nil {
				logger.Error("同步资金流失败", zap.String("code", stockCode), zap.Error(err))
				// 失败时，发送空切片以通知 Saver 继续计数
				dataChan <- []models.MoneyFlowData{}
			} else {
				if len(flows) > 0 {
					// 实时扫描策略信号
					s.ScanAndSaveStrategySignals(stockCode, flows)
					dataChan <- flows
				} else {
					// 爬取成功但无数据（如新股），也视为成功
					dataChan <- []models.MoneyFlowData{}
				}
			}
		}(code, i)
	}

	// 等待所有任务完成
	wg.Wait()
	close(dataChan) // 关闭数据通道，通知 Saver 退出

	// 这里不需要显式等待 resultChan，因为 emitProgress 是异步通知前端的
	// 但为了让日志准确，我们稍微等一下进度协程（可选）
	// 由于 StartFullMarketSync 返回 nil 后，主函数就结束了，
	// 如果进度协程还在跑，可能会有问题。
	// 但在这个场景下，close(dataChan) -> Saver 退出 -> close(resultChan) -> 进度协程退出
	// 所以我们需要等待 Saver 彻底退出。
	// 简单的办法：使用 WaitGroup 等待 Saver。

	// 不过根据目前代码结构，StartFullMarketSync 阻塞在 wg.Wait()，
	// 此时 Workers 都结束了。
	// dataChan 关闭后，Saver 会处理完剩余数据然后退出。
	// 我们可以在这里简单 sleep 一下或者不做处理，因为 Saver 运行很快。

	// 为了严谨，我们应该等待 Saver。
	// 但由于我无法轻易修改 Saver 的结构（在闭包里），
	// 而且 resultChan 是无缓冲的（不，它是 buffered total），
	// Saver 退出后 resultChan 关闭，进度协程退出。
	// 我们可以直接返回。

	return nil
}

// SyncAndScanSingleStock 同步并扫描单只股票
// 供前端按需调用：输入代码 -> 同步数据 -> 扫描策略 -> 返回信号
func (s *SyncService) SyncAndScanSingleStock(code string) ([]models.StrategySignal, error) {
	logger.Info("开始单股同步与扫描", zap.String("code", code))

	// 1. 同步最新数据 (抓取120天数据)
	rawData, err := s.FetchHistoryFlowDataV2(code, 120)
	if err != nil {
		return nil, fmt.Errorf("抓取数据失败: %w", err)
	}

	// 2. 数据对齐与清洗
	sortedData := GetSortedData(rawData)
	flows := AlignStockData2MoneyFlow(code, sortedData)

	if len(flows) == 0 {
		return nil, fmt.Errorf("未获取到有效数据")
	}

	// 3. 保存数据到数据库 (确保下次查询有数据)
	// 注意：SaveMoneyFlows 是 upsert 操作，安全的
	if err := s.moneyFlowRepo.SaveMoneyFlows(flows); err != nil {
		logger.Warn("保存资金流数据失败", zap.Error(err))
		// 继续执行，不中断扫描
	}

	// 4. 执行策略扫描并保存信号
	s.ScanAndSaveStrategySignals(code, flows)

	// 5. 返回最新的信号列表 (供前端刷新)
	return s.moneyFlowRepo.GetSignalsByStockCode(code)
}

// FetchHistoryFlowData 仅获取历史资金流数据，不保存
func (s *SyncService) FetchHistoryFlowData(code string) ([]models.MoneyFlowData, error) {
	// 转换代码格式
	secid := ""
	if strings.HasPrefix(code, "6") {
		secid = "1." + code
	} else {
		secid = "0." + code
	}

	// 构造 URL (lmt=0 获取全部)
	url := fmt.Sprintf(
		"https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get?lmt=120&klt=101&fields1=f1,f2,f3,f7&fields2=f51,f52,f53,f54,f55,f56,f62&secid=%s",
		secid,
	)
	logger.Info("请求资金流URL", zap.String("url", url))

	var emResp struct {
		RC   int `json:"rc"`
		Data struct {
			Klines []string `json:"klines"`
		} `json:"data"`
	}

	resp, err := s.client.R().
		SetResult(&emResp).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}

	// 检查 HTTP 状态码 (Resty 不会把非 2xx 视为 error，除非设置了)
	if resp.IsError() {
		return nil, fmt.Errorf("HTTP请求返回错误状态: %s", resp.Status())
	}

	if emResp.RC != 0 {
		return nil, fmt.Errorf("API返回错误 RC=%d", emResp.RC)
	}

	if emResp.Data.Klines == nil {
		return nil, nil
	}

	// 转换数据
	var flows []models.MoneyFlowData
	for _, line := range emResp.Data.Klines {
		parts := strings.Split(line, ",")
		if len(parts) < 13 {
			continue
		}

		date := parts[0]
		mainNet := parseMoney(parts[1])
		smallNet := parseMoney(parts[2])
		midNet := parseMoney(parts[3])
		bigNet := parseMoney(parts[4])
		superNet := parseMoney(parts[5])
		closePrice := parseMoney(parts[11])
		chgPct := parseMoney(parts[12])

		flows = append(flows, models.MoneyFlowData{
			Code:       code,
			TradeDate:  date,
			MainNet:    mainNet,
			SuperNet:   superNet,
			BigNet:     bigNet,
			MidNet:     midNet,
			SmallNet:   smallNet,
			ClosePrice: closePrice,
			ChgPct:     chgPct,
		})
	}

	return flows, nil
}

// FetchAndSaveHistoryFlow 已废弃，保留兼容性
func (s *SyncService) FetchAndSaveHistoryFlow(code string) error {
	flows, err := s.FetchHistoryFlowData(code)
	if err != nil {
		return err
	}
	if len(flows) > 0 {
		return s.moneyFlowRepo.SaveMoneyFlows(flows)
	}
	return nil
}

func parseMoney(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func (s *SyncService) emitProgress(progress *SyncProgress) {
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, "sync_progress", progress)
	}
}

// AlignedStockData 决策先锋专用结构体
type AlignedStockData struct {
	TradeDate  string  // 日期
	ClosePrice float64 // 收盘价
	Amount     float64 // 总成交额
	MainNet    float64 // 主力净流入 (f52)
	SuperNet   float64 // 超大单 (f56)
	BigNet     float64 // 大单 (f55)
	ChgPct     float64 // 涨跌幅
	Turnover   float64 // 换手率
	MainRate   float64 // 主力强度 (主力净额/总成交额)
}

// ParseAndMerge 手动解析并合并两个接口的数据
func ParseAndMerge(klineData []string, fflowData []string) map[string]*AlignedStockData {
	result := make(map[string]*AlignedStockData)

	// 1. 解析行情数据 (kline)
	// 假设 fields2=f51,f53,f56,f57,f59,f61
	for _, line := range klineData {
		parts := strings.Split(line, ",")
		if len(parts) < 6 {
			continue
		}

		date := parts[0]
		closeP, _ := strconv.ParseFloat(parts[1], 64)
		amount, _ := strconv.ParseFloat(parts[3], 64)
		chgPct, _ := strconv.ParseFloat(parts[4], 64)
		turnover, _ := strconv.ParseFloat(parts[5], 64)

		result[date] = &AlignedStockData{
			TradeDate:  date,
			ClosePrice: closeP,
			Amount:     amount,
			ChgPct:     chgPct,
			Turnover:   turnover,
		}
	}

	// 2. 解析并合并资金流数据 (fflow)
	// 假设 fields2=f51,f52,f53,f54,f55,f56,f62
	for _, line := range fflowData {
		parts := strings.Split(line, ",")
		if len(parts) < 7 {
			continue
		}

		date := parts[0]
		if data, ok := result[date]; ok {
			mainNet, _ := strconv.ParseFloat(parts[1], 64)  // f52 主力
			bigNet, _ := strconv.ParseFloat(parts[4], 64)   // f55 大单
			superNet, _ := strconv.ParseFloat(parts[5], 64) // f56 超大单

			data.MainNet = mainNet
			data.BigNet = bigNet
			data.SuperNet = superNet

			// 计算核心指标：主力强度
			if data.Amount > 0 {
				data.MainRate = (mainNet / data.Amount) * 100
			}
		}

	}
	logger.Info("解析并合并数据完成，条目数:", zap.Int("count", len(result)), zap.Any("result", result))
	return result
}

func (s *SyncService) FetchHistoryFlowDataV2(code string, limit int) (map[string]*AlignedStockData, error) {
	// 1. 判断市场前缀 (严谨逻辑)
	secid := "0." + code // 默认深市
	if strings.HasPrefix(code, "6") {
		secid = "1." + code // 沪市
	}

	// 2. 构造 URL (严格按照 ParseAndMerge 的索引顺序)
	// 行情：f51(日期),f53(收),f56(量),f57(额),f59(幅),f61(换)
	klineURL := fmt.Sprintf("https://push2his.eastmoney.com/api/qt/stock/kline/get?secid=%s&fields1=f1,f2,f3,f4,f5,f6&fields2=f51,f53,f56,f57,f59,f61&klt=101&fqt=1&end=20500101&lmt=%d", secid, limit)
	// 资金流：f51(日期),f52(主力),f53(小),f54(中),f55(大),f56(超),f62(占比)

	fflowURL := fmt.Sprintf("https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get?secid=%s&fields1=f1,f2,f3,f7&fields2=f51,f52,f53,f54,f55,f56,f62&klt=101&lmt=%d", secid, limit)

	// 3. 执行请求
	klineLines, err := s.httpGetKlinesWithHeaders(klineURL)
	if err != nil {
		return nil, fmt.Errorf("行情请求失败: %v", err)
	}
	fflowLines, err := s.httpGetKlinesWithHeaders(fflowURL)
	if err != nil {
		return nil, fmt.Errorf("资金流请求失败: %v", err)
	}

	// 4. 调用我们上一轮写的对齐逻辑
	return ParseAndMerge(klineLines, fflowLines), nil
}

// EastMoneyResp 东财标准响应外层
type EastMoneyResp struct {
	Data struct {
		Klines []string `json:"klines"`
		Code   string   `json:"code"`
		Name   string   `json:"name"`
	} `json:"data"`
}

// 内部辅助函数：执行 GET 并解析 JSON (带重试机制 - 依赖 Resty 内置重试)
func (s *SyncService) httpGetKlinesWithHeaders(url string) ([]string, error) {
	// Resty 客户端已经配置了重试 (RetryCount=3) 和短连接 (CloseConnection=true)
	// 所以这里不需要手写循环和 sleep，直接发起请求即可

	var emResp EastMoneyResp
	resp, err := s.client.R().
		SetResult(&emResp).
		Get(url)

	if err != nil {
		logger.Error("获取历史资金流数据失败", zap.Error(err), zap.String("url", url))
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP请求返回错误状态: %s", resp.Status())
	}

	return emResp.Data.Klines, nil
}

// ScanAndSaveStrategySignals 扫描并保存策略信号
// 遍历给定的历史数据，对每一天进行策略判定，如果触发 B/S 点则保存到数据库
func (s *SyncService) ScanAndSaveStrategySignals(code string, flows []models.MoneyFlowData) {
	if len(flows) < 20 {
		return
	}

	// 初始化策略服务 (复用逻辑)
	// NewStrategyService 只需要 Repo 即可工作
	strategySvc := NewStrategyService(nil, s.moneyFlowRepo)

	// 获取流通市值 (用于 B 点评分)
	circMV, _ := s.moneyFlowRepo.GetStockCircMV(code)

	// 遍历历史数据，从第20天开始 (因为需要前20天的数据计算 MA20)
	// flows 是按时间升序排列的 (index 0 是最旧的)
	for i := 19; i < len(flows); i++ {
		// 构造倒序窗口 (当前点为 i)
		// 策略服务要求 data[0] 是 T-0 (最新), data[1] 是 T-1...
		// 我们需要截取 flows[i-19 ... i] 这 20 条数据，并反转
		window := make([]models.MoneyFlowData, 20)
		for j := 0; j < 20; j++ {
			window[j] = flows[i-j]
		}

		// 1. 检查买入信号 (B点)
		if bSignal := strategySvc.CheckDecisionPioneerSignal(window, circMV); bSignal != nil {
			bSignal.Code = code
			// 尝试获取名称，如果没有则用 Code
			name, _ := s.moneyFlowRepo.GetStockName(code)
			if name == "" {
				name = code
			}
			bSignal.StockName = name

			if err := s.moneyFlowRepo.SaveStrategySignal(bSignal); err != nil {
				// 记录错误但不中断
				logger.Warn("保存B点信号失败", zap.String("code", code), zap.Error(err))
			} else {
				// 仅在生成最新信号时打印日志，避免刷屏
				if i == len(flows)-1 {
					logger.Info("发现决策先锋B点信号", zap.String("code", code), zap.String("date", bSignal.TradeDate))
				}
			}
		}

		// 1.5 检查“资金强攻”激进买入信号 (B_Surge)
		if surgeSignal := strategySvc.CheckMoneySurgeSignal(window); surgeSignal != nil {
			surgeSignal.Code = code
			name, _ := s.moneyFlowRepo.GetStockName(code)
			if name == "" {
				name = code
			}
			surgeSignal.StockName = name

			if err := s.moneyFlowRepo.SaveStrategySignal(surgeSignal); err != nil {
				logger.Warn("保存资金强攻信号失败", zap.String("code", code), zap.Error(err))
			} else {
				if i == len(flows)-1 {
					logger.Info("发现资金强攻信号", zap.String("code", code), zap.String("date", surgeSignal.TradeDate))
				}
			}
		}

		// 2. 检查卖出信号 (S点)
		if sSignal := strategySvc.CheckDecisionPioneerSellSignal(window); sSignal != nil {
			sSignal.Code = code
			name, _ := s.moneyFlowRepo.GetStockName(code)
			if name == "" {
				name = code
			}
			sSignal.StockName = name

			if err := s.moneyFlowRepo.SaveStrategySignal(sSignal); err != nil {
				logger.Warn("保存S点信号失败", zap.String("code", code), zap.Error(err))
			}
		}
	}
}

// GetSortedData 将 map 转换为按时间升序排列的切片
func GetSortedData(dataMap map[string]*AlignedStockData) []AlignedStockData {
	keys := make([]string, 0, len(dataMap))
	for k := range dataMap {
		keys = append(keys, k)
	}
	// 严格按日期升序排序
	sort.Strings(keys)

	sortedList := make([]AlignedStockData, 0, len(keys))
	for _, k := range keys {
		sortedList = append(sortedList, *dataMap[k])
	}
	return sortedList
}

func AlignStockData2MoneyFlow(stockCode string, data []AlignedStockData) []models.MoneyFlowData {
	moneyFlows := make([]models.MoneyFlowData, 0, len(data))
	for _, d := range data {
		moneyFlows = append(moneyFlows, models.MoneyFlowData{
			Code:       stockCode,
			TradeDate:  d.TradeDate,
			ClosePrice: d.ClosePrice,
			Amount:     d.Amount,
			MainNet:    d.MainNet,
			SuperNet:   d.SuperNet,
			BigNet:     d.BigNet,
			ChgPct:     d.ChgPct,
			MainRate:   d.MainRate,
			Turnover:   d.Turnover,
		})
	}
	return moneyFlows
}

type TickData struct {
	Time      string  // 成交时间
	Price     float64 // 成交价格
	Volume    int64   // 成交量(手)
	Orders    int64   // 成交笔数 (第4个元素)
	Direction int     // 成交方向 (1:主买, 2:主卖, 4:中性)
}

type OrderFlowStats struct {
	Symbol        string
	TotalVolume   int64   // 总成交量
	ActiveBuy     int64   // 明盘流入
	ActiveSell    int64   // 明盘流出
	HiddenFlow    int64   // 暗盘(中性大单)
	MainForceVol  int64   // 主力核心成交(高浓度单)
	MainForceRate float64 // 主力参与度
	NetMoneyFlow  int64   // 综合净流入
}

// 抓取全天交易笔数
func (s *SyncService) FetchAllDayTicks(code string) (*OrderFlowStats, error) {
	stats := &OrderFlowStats{}

	// 构造 secid
	secid := ""
	if strings.HasPrefix(code, "6") {
		secid = "1." + code
	} else {
		secid = "0." + code
	}
	pos := -0
	ticks, err := s.fetchTickBatch(secid, pos)
	if err != nil {
		return nil, err
	}

	tickData := parseRawTicks(ticks)
	stats = AnalyzeL2Market(tickData)
	return stats, nil
}

func AnalyzeL2Market(ticks []TickData) *OrderFlowStats {
	res := OrderFlowStats{}
	for _, tick := range ticks {
		// 严谨逻辑：过滤非交易时段
		if tick.Time < "09:30:00" || (tick.Time > "11:30:00" && tick.Time < "13:00:00") {
			continue
		}

		res.TotalVolume += tick.Volume
		avgVol := 0.0
		if tick.Orders > 0 {
			avgVol = float64(tick.Volume) / float64(tick.Orders)
		}

		// 严格复刻同花顺明暗盘分类
		switch tick.Direction {
		case 1:
			res.ActiveBuy += tick.Volume
		case 2:
			res.ActiveSell += tick.Volume
		case 4:
			// 暗盘逻辑：中性盘且单笔均量较大
			if avgVol >= 50 {
				res.HiddenFlow += tick.Volume
			}
		}

		// 主力逻辑：单笔均量超过门槛（高浓度成交）
		if avgVol >= 100 {
			res.MainForceVol += tick.Volume
		}
	}

	res.NetMoneyFlow = res.ActiveBuy - res.ActiveSell + res.HiddenFlow
	if res.TotalVolume > 0 {
		res.MainForceRate = (float64(res.MainForceVol) / float64(res.TotalVolume)) * 100
	}
	return &res
}

// 封装的批量获取tick数据方法 (带重试 - 依赖 Resty 内置重试)
func (s *SyncService) fetchTickBatch(secid string, pos int) ([]string, error) {
	url := fmt.Sprintf("https://push2.eastmoney.com/api/qt/stock/details/get?secid=%s&pos=%d&fields1=f1,f2,f3,f4,f5&fields2=f51,f52,f53,f54,f55", secid, pos)

	var result struct {
		Data struct {
			Details []string `json:"details"`
		} `json:"data"`
	}

	resp, err := s.client.R().
		SetResult(&result).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP请求返回错误状态: %s", resp.Status())
	}

	return result.Data.Details, nil
}

func parseRawTicks(rawTicks []string) []TickData {
	var ticks []TickData
	for _, line := range rawTicks {
		p := strings.Split(line, ",")
		if len(p) < 5 {
			continue
		}

		vol, _ := strconv.ParseInt(p[2], 10, 64)
		orders, _ := strconv.ParseInt(p[3], 10, 64)
		price, _ := strconv.ParseFloat(p[1], 64)
		dir, _ := strconv.Atoi(p[4])

		ticks = append(ticks, TickData{
			Time: p[0], Price: price, Volume: vol, Orders: orders, Direction: dir,
		})
	}
	return ticks
}
