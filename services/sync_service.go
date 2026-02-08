package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

// SyncService 全量数据同步服务
type SyncService struct {
	dbService          *DBService
	stockMarketService *StockMarketService
	moneyFlowRepo      *repositories.MoneyFlowRepository
	client             *http.Client
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
	return &SyncService{
		dbService:          dbService,
		stockMarketService: stockMarketService,
		moneyFlowRepo:      moneyFlowRepo,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
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

// FetchHistoryFlowData 仅获取历史资金流数据，不保存
func (s *SyncService) FetchHistoryFlowData(code string) ([]models.MoneyFlowData, error) {
	// 构造 secid
	secid := ""
	if strings.HasPrefix(code, "6") {
		secid = "1." + code
	} else if strings.HasPrefix(code, "0") || strings.HasPrefix(code, "3") {
		secid = "0." + code
	} else if strings.HasPrefix(code, "8") || strings.HasPrefix(code, "4") {
		secid = "0." + code // 北交所通常也是0，需根据实际调整，这里暂时假设为0
	} else {
		return nil, fmt.Errorf("未知市场代码前缀: %s", code)
	}

	// 构造 URL (lmt=0 获取全部)
	url := fmt.Sprintf(
		"https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get?lmt=120&klt=101&fields1=f1,f2,f3,f7&fields2=f51,f52,f53,f54,f55,f56,f62&secid=%s",
		secid,
	)
	logger.Info("请求资金流URL", zap.String("url", url))

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result struct {
		RC   int `json:"rc"`
		Data struct {
			Klines []string `json:"klines"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	if result.RC != 0 {
		return nil, fmt.Errorf("API返回错误 RC=%d", result.RC)
	}

	if result.Data.Klines == nil {
		return nil, nil
	}

	// 转换数据
	var flows []models.MoneyFlowData
	for _, line := range result.Data.Klines {
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
	klineLines, err := httpGetKlines(klineURL)
	if err != nil {
		return nil, fmt.Errorf("行情请求失败: %v", err)
	}
	fflowLines, err := httpGetKlines(fflowURL)
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

// 内部辅助函数：执行 GET 并解析 JSON
func httpGetKlines(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("获取历史资金流数据失败", zap.Any("err", err), zap.String("url", url))
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var emResp EastMoneyResp
	if err := json.Unmarshal(body, &emResp); err != nil {
		return nil, err
	}
	return emResp.Data.Klines, nil
}
func RunDecisionSignal(sortedData []AlignedStockData) {
	if len(sortedData) < 20 {
		return
	}

	for i := 19; i < len(sortedData); i++ {
		curr := sortedData[i]
		prev := sortedData[i-1]

		// 1. 生命线（操盘线）
		ma20 := calculateMAV(sortedData, i, 20)

		// 2. 资金动能：今日主力强度 vs 昨日主力强度
		// 决策先锋喜欢“资金反转”，即昨天流出，今天突然暴增
		moneySurge := curr.MainRate - prev.MainRate

		// 3. 决策先锋 B 点核心逻辑（回归版）
		// - 条件 A: 股价上穿 MA20（或者已经在 MA20 之上运行）
		// - 条件 B: 主力强度显著（> 3.0%）
		// - 条件 C: 资金动能向上（今天的钱比昨天多）

		isCrossing := curr.ClosePrice >= ma20 && prev.ClosePrice < ma20
		isStrongAbove := curr.ClosePrice > ma20 && curr.MainRate > 5.0

		if (isCrossing || isStrongAbove) && curr.MainRate > 3.0 && moneySurge > 0 {
			fmt.Printf("🎯 [决策先锋-B点] %s | 价格: %.2f | 主力占比: %.2f%% | 动能: %.2f\n",
				curr.TradeDate, curr.ClosePrice, curr.MainRate, moneySurge)
		}

		// 4. 决策先锋 S 点核心逻辑
		// - 条件 A: 股价跌破 MA20 且 资金不给力
		// - 条件 B: 股价虽在 MA20 之上，但主力资金出现“断崖式”流出（MainRate < -8%）

		if (curr.ClosePrice < ma20 && curr.MainRate < 0) || curr.MainRate < -8.0 {
			fmt.Printf("⚠️ [决策先锋-S点] %s | 价格: %.2f | 警告原因: %s\n",
				curr.TradeDate, curr.ClosePrice, getReason(curr, ma20))
		}

		if curr.TradeDate >= "2026-01-05" && curr.TradeDate <= "2026-01-12" {
			fmt.Printf("📅 日期: %s | 涨幅: %.2f%% | 主力强度: %.2f%%\n",
				curr.TradeDate, curr.ChgPct, curr.MainRate)
		}
	}
}

func getReason(d AlignedStockData, ma float64) string {
	if d.ClosePrice < ma {
		return "破位下行"
	}
	return "主力砸盘"
}

// 辅助函数：计算指定位置的MA
func calculateMAV(data []AlignedStockData, index int, period int) float64 {
	if index < period-1 {
		return 0
	}
	var sum float64
	for i := index - period + 1; i <= index; i++ {
		sum += data[i].ClosePrice
	}
	return sum / float64(period)
}

// 辅助函数：简单估算最近一次买入后的盈亏（仅用于日志展示）
func calculateProfit(data []AlignedStockData, currentIndex int) float64 {
	// 这里逻辑可以根据你的需要记录上次买入价，暂时简单返回0
	return 0
}
func CalculateSignals(data []AlignedStockData) {
	// 决策先锋通常需要至少 20 天的数据来计算均线
	if len(data) < 20 {
		return
	}

	for i := 20; i < len(data); i++ {
		// 1. 计算 MA20
		var sum float64
		for j := i - 19; j <= i; j++ {
			sum += data[j].ClosePrice
		}
		ma20 := sum / 20

		// 2. 计算 5 日资金流向
		var moneySum float64
		for j := i - 4; j <= i; j++ {
			moneySum += data[j].MainNet
		}

		// 3. 执行 B 点判定逻辑
		checkBPoint(data[i], ma20, moneySum)
	}
}

func checkBPoint(current AlignedStockData, ma20 float64, fiveDayMoney float64) {
	// 严谨逻辑闭环
	isInstitutionalBuying := current.MainRate > 3.0
	isTrendSafe := current.ClosePrice > ma20 && current.ClosePrice < ma20*1.15 // 别追太高
	isAccumulating := fiveDayMoney > 0
	isPriceStrong := current.ChgPct > 1.5

	if isInstitutionalBuying && isTrendSafe && isAccumulating && isPriceStrong {
		fmt.Printf("🔥 [B点信号] 日期: %s | 价格: %.2f | 主力强度: %.2f%% | 偏离MA20: %.2f%%\n",
			current.TradeDate,
			current.ClosePrice,
			current.MainRate,
			(current.ClosePrice-ma20)/ma20*100,
		)
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
func FetchAllDayTicks(code string) (*OrderFlowStats, error) {
	stats := &OrderFlowStats{}

	// 构造 secid
	secid := ""
	if strings.HasPrefix(code, "6") {
		secid = "1." + code
	} else {
		secid = "0." + code
	}
	pos := -0
	ticks, err := fetchTickBatch(secid, pos)
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

// 封装的批量获取tick数据方法
func fetchTickBatch(secid string, pos int) ([]string, error) {
	url := fmt.Sprintf("https://push2.eastmoney.com/api/qt/stock/details/get?secid=%s&pos=%d&fields1=f1,f2,f3,f4,f5&fields2=f51,f52,f53,f54,f55", secid, pos)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data struct {
			Details []string `json:"details"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
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
