package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"runtime/debug"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/services"
	"strings"
	"sync"
	"time"

	"stock-analyzer-wails/internal/logger"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

// App 应用程序结构
type App struct {
	ctx              context.Context
	stockService     *services.StockService
	aiService        *services.AIService
	watchlistService *services.WatchlistService
		alertStorage     *services.AlertStorage
		positionStorage  *services.PositionStorageService
		dbService        *services.DBService // 新增 DBService
		aiInitErr        error
	alerts           []*models.PriceAlert
	alertMutex       sync.Mutex
	alertConfig      models.AlertConfig
}

// NewApp 创建新的App应用程序
	func NewApp() *App {
		// 初始化数据库服务
		dbSvc, err := services.NewDBService()
		if err != nil {
			// 数据库初始化失败是致命错误，这里直接 panic 或返回 nil
			// 但由于 NewApp 不返回 error，我们先记录错误并返回一个 App 实例
			logger.Error("初始化数据库服务失败", zap.Error(err))
		}

		// 使用 DBService 初始化其他服务
		// 注意：这里暂时使用旧的 NewXXXService()，后续需要修改这些服务的构造函数以接受 dbSvc
		watchlistSvc := services.NewWatchlistService(dbSvc)
		alertSvc := services.NewAlertStorage(dbSvc)
			positionSvc := services.NewPositionStorageService(dbSvc)

		return &App{
			stockService:     services.NewStockService(),
			aiService:        nil,
			dbService:        dbSvc, // 存储 DBService
			watchlistService: watchlistSvc,
			alertStorage:     alertSvc,
			positionStorage:  positionSvc,
			alertConfig: models.AlertConfig{
			Sensitivity: 0.005, // 默认 0.5%
			Cooldown:    1,     // 默认 1 小时
			Enabled:     true,
		},
	}
}

// startup 在应用程序启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.initAIService()
	
	// 加载持久化的预警订阅
	if a.alertStorage != nil {
		alerts, err := a.alertStorage.LoadActiveAlerts()
		if err == nil {
			a.alertMutex.Lock()
			a.alerts = alerts
			a.alertMutex.Unlock()
			logger.Info("成功加载持久化预警订阅", zap.Int("count", len(alerts)))
		}
	}
	
	go a.startAlertMonitor()
	go a.startPositionMonitor()
}

// startPositionMonitor 启动持仓逻辑监控引擎
func (a *App) startPositionMonitor() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.checkPositionLogics()
		}
	}
}

// checkPositionLogics 校验所有活跃持仓的建仓逻辑
func (a *App) checkPositionLogics() {
	positions, err := a.positionStorage.GetPositions()
	if err != nil || len(positions) == 0 {
		return
	}

	for _, pos := range positions {
		if pos.CurrentStatus != "holding" {
			continue
		}

		// 1. 获取最新实时数据
		stock, err := a.stockService.GetStockByCode(pos.StockCode)
		if err != nil {
			continue
		}

		moneyFlow, err := a.stockService.GetMoneyFlowData(pos.StockCode)
		if err != nil {
			continue
		}

		// 2. 执行移动止损算法 (M4)
		a.updateTrailingStop(pos, stock.Price)

		violatedReasons := make([]string, 0)
		
		// 3. 校验止损位 (硬性逻辑)
		if stock.Price < pos.Strategy.StopLossPrice {
			violatedReasons = append(violatedReasons, fmt.Sprintf("股价(%.2f)已跌破止损位(%.2f)", stock.Price, pos.Strategy.StopLossPrice))
		}

		// 4. 校验资金流逻辑 (基于 AI 设定的核心理由)
		for _, reason := range pos.Strategy.CoreReasons {
			if reason.Type == "money_flow" {
				// 简单启发式：如果今日主力净流出超过 5000 万，且理由中包含主力流入
				if moneyFlow.TodayMain < -50000000 {
					violatedReasons = append(violatedReasons, "主力资金出现大幅流出，背离建仓逻辑")
				}
			}
		}

		// 5. 如果逻辑失效，触发预警并更新状态
		if len(violatedReasons) > 0 && pos.LogicStatus != "violated" {
			pos.LogicStatus = "violated"
			pos.UpdatedAt = time.Now()
			a.positionStorage.SavePosition(pos)

			// 发送 Wails 事件通知前端
			runtime.EventsEmit(a.ctx, "logic_violation", map[string]interface{}{
				"code":    pos.StockCode,
				"name":    pos.StockName,
				"reasons": violatedReasons,
				"price":   stock.Price,
			})
			
			logger.Warn("持仓逻辑失效预警", 
				zap.String("code", pos.StockCode), 
				zap.Strings("reasons", violatedReasons))
		}
	}
}

// updateTrailingStop 动态更新移动止损位 (M4)
func (a *App) updateTrailingStop(pos *models.Position, currentPrice float64) {
	config := pos.TrailingConfig
	// 如果未启用移动止损，或者股价低于买入价，则跳过
	if !config.Enabled || currentPrice <= pos.EntryPrice {
		return
	}

	// 计算当前盈利比例
	profitRate := (currentPrice - pos.EntryPrice) / pos.EntryPrice
	
	// 动态调整止损位 (参数化移动止损)
	newStopLoss := pos.Strategy.StopLossPrice
	
	// 只有当盈利超过触发阈值时才启动
	if profitRate > config.ActivationThreshold {
		// 使用配置的回撤比例计算新止损位
		potentialStop := currentPrice * (1 - config.CallbackRate)
		if potentialStop > pos.Strategy.StopLossPrice {
			newStopLoss = potentialStop
		}
	}

	// 只有当新的止损位高于旧的止损位时才更新 (止损位只能上移，不能下移)
	if newStopLoss > pos.Strategy.StopLossPrice {
		oldStop := pos.Strategy.StopLossPrice
		pos.Strategy.StopLossPrice = newStopLoss
		pos.UpdatedAt = time.Now()
		a.positionStorage.SavePosition(pos)

		// 发送 Wails 事件通知用户止损位已上移
		runtime.EventsEmit(a.ctx, "stop_loss_raised", map[string]interface{}{
			"code":    pos.StockCode,
			"name":    pos.StockName,
			"oldStop": oldStop,
			"newStop": newStopLoss,
			"price":   currentPrice,
		})
		
		logger.Info("移动止损位上移", 
			zap.String("code", pos.StockCode), 
			zap.Float64("old", oldStop), 
			zap.Float64("new", newStopLoss))
	}
}

// startAlertMonitor 启动价格预警监控引擎
func (a *App) startAlertMonitor() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.checkAlerts()
		}
	}
}

// checkAlerts 检查所有激活的预警
func (a *App) checkAlerts() {
	a.alertMutex.Lock()
	activeAlerts := make([]*models.PriceAlert, 0)
	for _, alert := range a.alerts {
		if alert.IsActive {
			activeAlerts = append(activeAlerts, alert)
		}
	}
	a.alertMutex.Unlock()

	if len(activeAlerts) == 0 {
		return
	}

	// 批量获取最新价进行比对
	for _, alert := range activeAlerts {
		stock, err := a.stockService.GetStockByCode(alert.StockCode)
		if err != nil {
			continue
		}

		triggered := false
		if alert.Type == "above" && stock.Price >= alert.Price {
			triggered = true
		} else if alert.Type == "below" && stock.Price <= alert.Price {
			triggered = true
		}

		if triggered {
			// 检查冷却时间
			if time.Since(alert.LastTriggered) < time.Duration(a.alertConfig.Cooldown)*time.Hour {
				continue
			}

			// 触发告警
			alert.LastTriggered = time.Now()
			
			// 发送 Wails 事件
			runtime.EventsEmit(a.ctx, "price_alert", map[string]interface{}{
				"code":  alert.StockCode,
				"name":  alert.StockName,
				"price": stock.Price,
				"type":  alert.Type,
				"target": alert.Price,
			})

			// 记录到历史
			if a.alertStorage != nil {
				a.alertStorage.SaveAlert(alert, fmt.Sprintf("股价已%s预警位 %.2f", map[string]string{"above":"突破","below":"跌破"}[alert.Type], alert.Price))
			}
		}
	}
}

// AddAlert 添加新的价格预警
func (a *App) AddAlert(alert models.PriceAlert) error {
	a.alertMutex.Lock()
	defer a.alertMutex.Unlock()

	alert.IsActive = true
	alert.LastTriggered = time.Unix(0, 0)
	a.alerts = append(a.alerts, &alert)

	// 持久化保存
	if a.alertStorage != nil {
		return a.alertStorage.SaveActiveAlerts(a.alerts)
	}
	return nil
}

// GetActiveAlerts 获取所有激活的预警
func (a *App) GetActiveAlerts() ([]*models.PriceAlert, error) {
	a.alertMutex.Lock()
	defer a.alertMutex.Unlock()
	return a.alerts, nil
}

// RemoveAlert 移除预警
func (a *App) RemoveAlert(stockCode string, alertType string, price float64) error {
	a.alertMutex.Lock()
	defer a.alertMutex.Unlock()

	newAlerts := make([]*models.PriceAlert, 0)
	for _, alert := range a.alerts {
		if alert.StockCode == stockCode && alert.Type == alertType && alert.Price == price {
			continue
		}
		newAlerts = append(newAlerts, alert)
	}
	a.alerts = newAlerts

	// 持久化保存
	if a.alertStorage != nil {
		return a.alertStorage.SaveActiveAlerts(a.alerts)
	}
	return nil
}

// GetAlertHistory 获取告警历史
func (a *App) GetAlertHistory(stockCode string, limit int) ([]map[string]interface{}, error) {
	if a.alertStorage == nil {
		return nil, fmt.Errorf("告警存储服务未就绪")
	}
	return a.alertStorage.GetAlertHistory(stockCode, limit)
}

// UpdateAlertConfig 更新告警全局配置
func (a *App) UpdateAlertConfig(config models.AlertConfig) error {
	a.alertConfig = config
	return nil
}

// GetAlertConfig 获取告警全局配置
func (a *App) GetAlertConfig() (models.AlertConfig, error) {
	return a.alertConfig, nil
}

// SetAlertsFromAI 接收 AI 识别的支撑位和压力位并自动设置预警
func (a *App) SetAlertsFromAI(code string, name string, drawings []models.TechnicalDrawing) {
	a.alertMutex.Lock()
	
	addedCount := 0
	for _, d := range drawings {
		if d.Price == 0 {
			continue
		}

		// 检查是否已存在相同的预警
		exists := false
		alertType := "above"
		if d.Type == "support" {
			alertType = "below"
		}

		for _, existing := range a.alerts {
			if existing.StockCode == code && existing.Type == alertType && MathAbs(existing.Price-d.Price) < 0.01 {
				exists = true
				break
			}
		}

		if !exists {
			a.alerts = append(a.alerts, &models.PriceAlert{
				StockCode:     code,
				StockName:     name,
				Price:         d.Price,
				Type:          alertType,
				IsActive:      true,
				LastTriggered: time.Unix(0, 0),
			})
			addedCount++
		}
	}
	
	newAlerts := a.alerts
	a.alertMutex.Unlock()
	
	// 持久化保存活跃预警
	if a.alertStorage != nil {
		a.alertStorage.SaveActiveAlerts(newAlerts)
	}
	
	logger.Info("预警位更新完成", zap.Int("added_count", addedCount), zap.Int("total_active", len(newAlerts)))
}

// MathAbs 辅助函数
func MathAbs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

// initAIService 初始化或重新初始化 AI 服务
func (a *App) initAIService() error {
	start := time.Now()
	cfg, err := services.LoadAIConfig()
	if err != nil {
		a.aiInitErr = err
		return err
	}

	if cfg.APIKey == "" {
		a.aiService = nil
		a.aiInitErr = fmt.Errorf("API Key 未配置")
		return nil
	}

	aiSvc, err := services.NewAIService(cfg)
	if err != nil {
		a.aiInitErr = err
		return err
	}

	a.aiService = aiSvc
	a.aiInitErr = nil

	logger.Info("AI 服务初始化成功",
		zap.String("module", "app"),
		zap.String("provider", string(cfg.Provider)),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return nil
}

// GetConfig 获取当前配置
func (a *App) GetConfig() (services.AIResolvedConfig, error) {
	return services.LoadAIConfig()
}

// SaveConfig 保存配置并重置 AI 服务
func (a *App) SaveConfig(config services.AIResolvedConfig) error {
	err := services.SaveAIConfig(config)
	if err != nil {
		return err
	}
	return a.initAIService()
}

// GetStockData 获取股票数据
func (a *App) GetStockData(code string) (*models.StockData, error) {
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	return a.stockService.GetStockByCode(code)
}

// GetIntradayData 获取分时数据
func (a *App) GetIntradayData(code string) (*models.IntradayResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	return a.stockService.GetIntradayData(code)
}

// GetMoneyFlowData 获取资金流向数据
func (a *App) GetMoneyFlowData(code string) (*models.MoneyFlowResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	return a.stockService.GetMoneyFlowData(code)
}

// GetStockHealthCheck 获取股票深度体检报告
func (a *App) GetStockHealthCheck(code string) (*models.HealthCheckResult, error) {
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	return a.stockService.GetStockHealthCheck(code)
}

// AddPosition 添加持仓记录
func (a *App) AddPosition(pos models.Position) error {
	return a.positionStorage.SavePosition(&pos)
}

// GetPositions 获取所有活跃持仓
func (a *App) GetPositions() (map[string]*models.Position, error) {
	return a.positionStorage.GetPositions()
}

// RemovePosition 移除持仓记录
func (a *App) RemovePosition(code string) error {
	return a.positionStorage.RemovePosition(code)
}

func newTraceID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// AnalyzeEntryStrategy 获取 AI 智能建仓方案
func (a *App) AnalyzeEntryStrategy(code string) (res *models.EntryStrategyResult, err error) {
	start := time.Now()
	traceId := newTraceID()
	code = strings.TrimSpace(code)

	defer func() {
		if r := recover(); r != nil {
			logger.Error("建仓分析发生 panic",
				zap.String("module", "app.entry_strategy"),
				zap.String("op", "AnalyzeEntryStrategy"),
				zap.String("step", "panic"),
				zap.String("stock_code", code),
				zap.String("traceId", traceId),
				zap.Any("panic", r),
				zap.ByteString("stack", debug.Stack()),
				zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			)
			res = nil
			err = fmt.Errorf("建仓分析失败(step=panic, code=ENTRY_PANIC, traceId=%s): 后端发生异常，请查看日志", traceId)
		}
	}()

	logger.Info("开始建仓分析（入口）",
		zap.String("module", "app.entry_strategy"),
		zap.String("op", "AnalyzeEntryStrategy"),
		zap.String("stock_code", code),
		zap.String("traceId", traceId),
	)

	if code == "" {
		return nil, fmt.Errorf("建仓分析失败(step=input, code=ENTRY_INPUT_INVALID, traceId=%s): 股票代码不能为空", traceId)
	}
	if a.aiService == nil {
		logger.Warn("建仓分析失败：AI 服务未就绪",
			zap.String("module", "app.entry_strategy"),
			zap.String("op", "AnalyzeEntryStrategy"),
			zap.String("step", "init"),
			zap.String("stock_code", code),
			zap.String("traceId", traceId),
			zap.Error(a.aiInitErr),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return nil, fmt.Errorf("建仓分析失败(step=init, code=ENTRY_AI_NOT_READY, traceId=%s): AI服务未就绪", traceId)
	}

	stepStart := time.Now()
	stock, e := a.stockService.GetStockByCode(code)
	if e != nil {
		logger.Error("建仓分析失败：获取股票数据失败",
			zap.String("module", "app.entry_strategy"),
			zap.String("op", "AnalyzeEntryStrategy"),
			zap.String("step", "stock_get"),
			zap.String("stock_code", code),
			zap.String("traceId", traceId),
			zap.Error(e),
			zap.Int64("duration_ms", time.Since(stepStart).Milliseconds()),
		)
		return nil, fmt.Errorf("建仓分析失败(step=stock_get, code=ENTRY_STOCK_FETCH_FAILED, traceId=%s): %w", traceId, e)
	}

	stepStart = time.Now()
	klines, e := a.stockService.GetKLineData(code, 100, "daily")
	if e != nil {
		logger.Error("建仓分析失败：获取K线数据失败",
			zap.String("module", "app.entry_strategy"),
			zap.String("op", "AnalyzeEntryStrategy"),
			zap.String("step", "kline_get"),
			zap.String("stock_code", code),
			zap.String("traceId", traceId),
			zap.Error(e),
			zap.Int64("duration_ms", time.Since(stepStart).Milliseconds()),
		)
		return nil, fmt.Errorf("建仓分析失败(step=kline_get, code=ENTRY_KLINE_FETCH_FAILED, traceId=%s): %w", traceId, e)
	}

	stepStart = time.Now()
	moneyFlow, e := a.stockService.GetMoneyFlowData(code)
	if e != nil {
		logger.Error("建仓分析失败：获取资金流向失败",
			zap.String("module", "app.entry_strategy"),
			zap.String("op", "AnalyzeEntryStrategy"),
			zap.String("step", "moneyflow_get"),
			zap.String("stock_code", code),
			zap.String("traceId", traceId),
			zap.Error(e),
			zap.Int64("duration_ms", time.Since(stepStart).Milliseconds()),
		)
		return nil, fmt.Errorf("建仓分析失败(step=moneyflow_get, code=ENTRY_MONEYFLOW_FETCH_FAILED, traceId=%s): %w", traceId, e)
	}

	stepStart = time.Now()
	health, e := a.stockService.GetStockHealthCheck(code)
	if e != nil {
		logger.Error("建仓分析失败：获取体检报告失败",
			zap.String("module", "app.entry_strategy"),
			zap.String("op", "AnalyzeEntryStrategy"),
			zap.String("step", "health_get"),
			zap.String("stock_code", code),
			zap.String("traceId", traceId),
			zap.Error(e),
			zap.Int64("duration_ms", time.Since(stepStart).Milliseconds()),
		)
		return nil, fmt.Errorf("建仓分析失败(step=health_get, code=ENTRY_HEALTH_FETCH_FAILED, traceId=%s): %w", traceId, e)
	}

	stepStart = time.Now()
	res, e = a.aiService.AnalyzeEntryStrategy(stock, klines, moneyFlow, health)
	if e != nil {
		logger.Error("建仓分析失败：AI 分析失败",
			zap.String("module", "app.entry_strategy"),
			zap.String("op", "AnalyzeEntryStrategy"),
			zap.String("step", "ai_analyze"),
			zap.String("stock_code", code),
			zap.String("traceId", traceId),
			zap.Error(e),
			zap.Int64("duration_ms", time.Since(stepStart).Milliseconds()),
		)
		// 不覆盖底层错误码（ENTRY_*），仅附加 traceId 方便排查
		return nil, fmt.Errorf("traceId=%s: %w", traceId, e)
	}

	logger.Info("建仓分析成功（入口）",
		zap.String("module", "app.entry_strategy"),
		zap.String("op", "AnalyzeEntryStrategy"),
		zap.String("stock_code", code),
		zap.String("traceId", traceId),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return res, nil
}

// BatchAnalyzeStocks 批量分析股票
func (a *App) BatchAnalyzeStocks(codes []string, role string) error {
	if a.aiService == nil {
		return fmt.Errorf("AI服务未就绪")
	}
	return a.stockService.BatchAnalyzeStocks(a.ctx, codes, role, a.aiService)
}

// GetKLineData 获取K线数据，支持周期参数
func (a *App) GetKLineData(code string, limit int, period string) ([]*models.KLineData, error) {
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	return a.stockService.GetKLineData(code, limit, period)
}

// AddToWatchlist 添加到自选股
func (a *App) AddToWatchlist(stock models.StockData) error {
	return a.watchlistService.AddToWatchlist(&stock)
}

// RemoveFromWatchlist 从自选股移除
func (a *App) RemoveFromWatchlist(code string) error {
	return a.watchlistService.RemoveFromWatchlist(code)
}

// GetWatchlist 获取自选股列表
func (a *App) GetWatchlist() ([]*models.StockData, error) {
	return a.watchlistService.GetWatchlist()
}

// SearchStock 搜索股票
func (a *App) SearchStock(keyword string) ([]*models.StockData, error) {
	return a.stockService.SearchStock(keyword)
}

// AnalyzeStock 分析股票
func (a *App) AnalyzeStock(code string) (*models.AnalysisReport, error) {
	if a.aiService == nil {
		return nil, fmt.Errorf("AI服务未就绪")
	}
	stock, err := a.stockService.GetStockByCode(code)
	if err != nil {
		return nil, err
	}
	return a.aiService.AnalyzeStock(stock)
}

// AnalyzeTechnical 深度技术分析
func (a *App) AnalyzeTechnical(code string, period string, role string) (*models.TechnicalAnalysisResult, error) {
	if a.aiService == nil {
		return nil, fmt.Errorf("AI服务未就绪")
	}
	stock, err := a.stockService.GetStockByCode(code)
	if err != nil {
		return nil, err
	}
	klines, err := a.stockService.GetKLineData(code, 100, period)
	if err != nil {
		return nil, err
	}
	return a.aiService.AnalyzeTechnical(stock, klines, period, role)
}


// shutdown 在应用程序退出时调用
func (a *App) shutdown(ctx context.Context) {
	// 关闭数据库连接
	if a.dbService != nil {
		a.dbService.Close()
	}
}
