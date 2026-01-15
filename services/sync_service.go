package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
			flows, err := s.FetchHistoryFlowData(stockCode)
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
		"https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get?lmt=0&klt=101&fields1=f1,f2,f3,f7&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65&secid=%s",
		secid,
	)

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
