package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"runtime/debug"
	"stock-analyzer-wails/controllers"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"
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
	ctx               context.Context
	stockService      *services.StockService
	aiService         *services.AIService
	alertStorage      *services.AlertService
	positionStorage   *services.PositionService
	configService     *services.ConfigService            // 新增 ConfigService
	dbService         *services.DBService                // 新增 DBService
	syncHistoryCtrl   *controllers.SyncHistoryController // 同步历史控制器
	aiInitErr         error
	alerts            []*models.PriceAlert
	alertMutex        sync.Mutex
	alertConfigMutex  sync.RWMutex
	alertConfig       models.AlertConfig
	klineSyncService  *services.KLineSyncService // K线同步服务
	syncService       *services.SyncService      // 全量同步服务
	priceAlertMonitor *services.AlertMonitor     // 价格预警监控引擎

	// Controllers (Wails Bindings)
	WatchlistController   *controllers.WatchlistController
	AlertController       *controllers.AlertController
	PositionController    *controllers.PositionController
	ConfigController      *controllers.ConfigController
	SyncHistoryController *controllers.SyncHistoryController
	StrategyController    *controllers.StrategyController    // 策略控制器
	StockMarketController *controllers.StockMarketController // 市场股票控制器
	PriceAlertController  *controllers.PriceAlertController  // 价格预警控制器

	// Services (for internal use)
	watchlistService  *services.WatchlistService  // 保持，用于内部逻辑调用
	priceAlertService *services.PriceAlertService // 价格预警服务（内部使用）
	backtestService   *services.BacktestService   // 回测服务
	strategyService   *services.StrategyService   // 策略服务
}

// NewApp 创建新的App应用程序
func NewApp() *App {
	// 初始化数据库服务
	dbSvc, err := services.NewDBService()
	if err != nil {
		// 数据库初始化失败会直接导致：自选股/预警/持仓/配置等 SQLite 相关功能不可用。
		// 这里不再继续使用 dbSvc.GetDB() 做 DI（否则 dbSvc 可能为 nil 导致 panic）。
		logger.Error("初始化数据库服务失败（SQLite 功能将不可用）",
			zap.String("module", "app"),
			zap.String("op", "NewApp"),
			zap.Error(err),
		)
		dbSvc = nil
	}

	// StockService 一定要创建，并在有 DB 时注入（否则 SyncStockData 会报“数据库服务未初始化”）。
	stockSvc := services.NewStockService()
	if dbSvc != nil {
		stockSvc.SetDBService(dbSvc)
	}

	// 如果数据库不可用，相关 controller/service 置空，避免启动阶段 panic。
	if dbSvc == nil {
		logger.Warn("SQLite 功能已降级：依赖数据库的模块将不可用（包括价格预警/自选股/配置/策略等）",
			zap.String("module", "app"),
			zap.String("op", "NewApp"),
		)
		return &App{
			stockService: stockSvc,
			aiService:    nil,
			dbService:    nil,
			alertConfig: models.AlertConfig{
				Sensitivity: 0.005,
				Cooldown:    1,
				Enabled:     true,
			},
		}
	}

	// --- 依赖注入 (DI) ---
	// 1. Repository 层
	watchlistRepo := repositories.NewSQLiteWatchlistRepository(dbSvc.GetDB())
	alertRepo := repositories.NewSQLiteAlertRepository(dbSvc.GetDB())
	positionRepo := repositories.NewSQLitePositionRepository(dbSvc.GetDB())
	configRepo := repositories.NewSQLiteConfigRepository(dbSvc.GetDB())
	syncHistoryRepo := repositories.NewSQLiteSyncHistoryRepository(dbSvc.GetDB())
	strategyRepo := repositories.NewStrategyRepository(dbSvc.GetDB())
	priceAlertRepo := repositories.NewPriceAlertRepository(dbSvc.GetDB())
	moneyFlowRepo := repositories.NewMoneyFlowRepository(dbSvc.GetDB()) // 新增 MoneyFlowRepository

	// 2. Service 层
	watchlistSvc := services.NewWatchlistService(watchlistRepo)
	alertSvc := services.NewAlertService(alertRepo)
	positionSvc := services.NewPositionService(positionRepo)
	configSvc := services.NewConfigService(configRepo)
	strategySvc := services.NewStrategyService(strategyRepo, moneyFlowRepo) // 注入 MoneyFlowRepository
	stockMarketSvc := services.NewStockMarketService(dbSvc)
	priceAlertSvc := services.NewPriceAlertService(priceAlertRepo)

	var klineSyncSvc *services.KLineSyncService
	var syncSvc *services.SyncService
	if dbSvc != nil {
		klineSyncSvc = services.NewKLineSyncService(dbSvc)
		syncSvc = services.NewSyncService(dbSvc, stockMarketSvc, moneyFlowRepo)
	}

	// 3. Controller 层 (Wails 绑定)
	watchlistCtrl := controllers.NewWatchlistController(watchlistSvc)
	alertCtrl := controllers.NewAlertController(alertSvc)
	positionCtrl := controllers.NewPositionController(positionSvc)
	configCtrl := controllers.NewConfigController(configSvc)
	syncHistoryCtrl := controllers.NewSyncHistoryController(syncHistoryRepo)
	strategyCtrl := controllers.NewStrategyController(strategySvc)
	stockMarketCtrl := controllers.NewStockMarketController(stockMarketSvc)
	priceAlertCtrl := controllers.NewPriceAlertController(priceAlertSvc)

	logger.Info("SQLite 初始化成功，已创建控制器绑定",
		zap.String("module", "app"),
		zap.String("op", "NewApp"),
		zap.String("dbPath", dbSvc.GetDBPath()),
		zap.Bool("hasPriceAlertController", priceAlertCtrl != nil),
		zap.Bool("hasAlertController", alertCtrl != nil),
		zap.Bool("hasWatchlistController", watchlistCtrl != nil),
		zap.Bool("hasConfigController", configCtrl != nil),
	)

	// 创建价格预警监控引擎（在 startup 中启动）
	//var alertMonitor *services.AlertMonitor
	// 注意：AlertMonitor 需要传入 context，所以在 startup 中创建

	// 4. 回测服务
	backtestSvc := services.NewBacktestService(stockSvc, strategySvc)

	return &App{
		stockService:     stockSvc,
		aiService:        nil,
		dbService:        dbSvc,        // 存储 DBService
		klineSyncService: klineSyncSvc, // K线同步服务
		syncService:      syncSvc,      // 全量同步服务
		backtestService:  backtestSvc,  // 回测服务

		// Controllers
		WatchlistController:   watchlistCtrl,
		AlertController:       alertCtrl,
		PositionController:    positionCtrl,
		ConfigController:      configCtrl,
		SyncHistoryController: syncHistoryCtrl,
		StrategyController:    strategyCtrl,
		StockMarketController: stockMarketCtrl, // 市场股票控制器
		PriceAlertController:  priceAlertCtrl,  // 价格预警控制器

		// Services (for internal use)
		watchlistService:  watchlistSvc,
		alertStorage:      alertSvc,
		positionStorage:   positionSvc,
		configService:     configSvc,
		syncHistoryCtrl:   syncHistoryCtrl, // 内部引用
		priceAlertService: priceAlertSvc,   // 价格预警服务
		strategyService:   strategySvc,     // 策略服务
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

	// 注册需要上下文的服务的 Startup 方法
	a.stockService.Startup(ctx)
	if a.klineSyncService != nil {
		a.klineSyncService.SetContext(ctx)
	}
	if a.syncService != nil {
		a.syncService.SetContext(ctx)
	}

	// 迁移旧的 AI 配置（数据库不可用时 configService 为空，需要安全跳过）
	if a.configService != nil {
		if err := a.configService.MigrateAIConfigFromYAML(); err != nil {
			logger.Error("执行 YAML 配置迁移失败", zap.Error(err))
		}
	} else {
		logger.Warn("跳过 YAML 配置迁移：ConfigService 未初始化（可能数据库不可用）")
	}

	if err := a.initAIService(); err != nil {
		// initAIService 内部会设置 a.aiInitErr，这里额外打日志方便定位
		logger.Warn("AI 服务初始化未完成", zap.Error(err))
	}

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

	// 数据库不可用时，这两个监控不启动（它们依赖 alertStorage/positionStorage）
	if a.alertStorage != nil {
		go a.startAlertMonitor()
	}
	if a.positionStorage != nil {
		go a.startPositionMonitor()
	}

	// 启动价格预警监控引擎
	if a.priceAlertService != nil && a.stockService != nil {
		a.priceAlertMonitor = services.NewAlertMonitor(
			a.ctx,
			a.priceAlertService,
			a.stockService,
			a.stockService, // StockService 实现了 KLineDataService 接口（通过 GetKLineData 方法）
		)
		a.priceAlertMonitor.Start()
		logger.Info("价格预警监控引擎已启动")
	}
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
	if a.positionStorage == nil {
		return
	}
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
	// 全局预警关闭时直接跳过，避免无意义轮询/触发
	a.alertConfigMutex.RLock()
	cfg := a.alertConfig
	a.alertConfigMutex.RUnlock()
	if !cfg.Enabled {
		return
	}

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
			if time.Since(alert.LastTriggered) < time.Duration(cfg.Cooldown)*time.Hour {
				continue
			}

			// 触发告警
			alert.LastTriggered = time.Now()

			// 发送 Wails 事件
			runtime.EventsEmit(a.ctx, "price_alert", map[string]interface{}{
				"code":   alert.StockCode,
				"name":   alert.StockName,
				"price":  stock.Price,
				"type":   alert.Type,
				"target": alert.Price,
			})

			// 记录到历史
			if a.alertStorage != nil {
				a.alertStorage.SaveAlert(alert, fmt.Sprintf("股价已%s预警位 %.2f", map[string]string{"above": "突破", "below": "跌破"}[alert.Type], alert.Price))
			}
		}
	}
}

// --- Alert 转发器 ---
// AddAlert 添加新的价格预警
func (a *App) AddAlert(alert models.PriceAlert) error {
	// 这是一个复杂逻辑，需要先获取所有预警，添加新的，再保存
	// 考虑到 AlertController 已经有 SaveAlerts 方法，我们直接调用它
	alerts, err := a.AlertController.GetAlerts()
	if err != nil {
		return err
	}
	alerts = append(alerts, &alert)
	return a.AlertController.SaveAlerts(alerts)
}

// GetActiveAlerts 获取所有激活的预警
func (a *App) GetActiveAlerts() ([]*models.PriceAlert, error) {
	return a.AlertController.GetAlerts()
}

// RemoveAlert 移除预警
func (a *App) RemoveAlert(stockCode string, alertType string, price float64) error {
	// 这是一个复杂逻辑，需要先获取所有预警，移除匹配的，再保存
	alerts, err := a.AlertController.GetAlerts()
	if err != nil {
		return err
	}

	newAlerts := make([]*models.PriceAlert, 0)
	for _, alert := range alerts {
		if alert.StockCode != stockCode || alert.Type != alertType || alert.Price != price {
			newAlerts = append(newAlerts, alert)
		}
	}
	return a.AlertController.SaveAlerts(newAlerts)
}

// GetAlertHistory 获取告警历史
func (a *App) GetAlertHistory(stockCode string, limit int) ([]map[string]interface{}, error) {
	return a.AlertController.GetAlertHistory(stockCode, limit)
}

// UpdateAlertConfig 更新告警全局配置
func (a *App) UpdateAlertConfig(config models.AlertConfig) error {
	// 即使数据库不可用/AlertController 未初始化，也应允许更新全局配置（内存态）
	a.alertConfigMutex.Lock()
	a.alertConfig = config
	a.alertConfigMutex.Unlock()

	// 如果控制器可用，顺带同步到 service（未来可扩展为持久化）
	if a.AlertController != nil {
		return a.AlertController.UpdateAlertConfig(config)
	}
	return nil
}

// GetAlertConfig 获取告警全局配置
func (a *App) GetAlertConfig() (models.AlertConfig, error) {
	// 优先走 controller（若未来做了持久化），失败则回退到内存配置
	if a.AlertController != nil {
		cfg, err := a.AlertController.GetAlertConfig()
		if err == nil {
			a.alertConfigMutex.Lock()
			a.alertConfig = cfg
			a.alertConfigMutex.Unlock()
			return cfg, nil
		}
	}

	a.alertConfigMutex.RLock()
	cfg := a.alertConfig
	a.alertConfigMutex.RUnlock()
	return cfg, nil
}

// SetAlertsFromAI 接收 AI 识别的支撑位和压力位并自动设置预警
func (a *App) SetAlertsFromAI(code string, name string, drawings []models.TechnicalDrawing) {
	a.AlertController.SetAlertsFromAI(code, name, drawings)
}

// --- Alert 转发器 结束 ---

// --- PriceAlertController 转发器 ---
// 注意：Wails 只会绑定传入 Bind 列表的结构体“方法”，不会把 App 的字段（如 PriceAlertController 指针）自动暴露给前端。
// 因此前端应调用这些转发方法：window.go.main.App.PriceAlertGetAllAlerts() 等。

func (a *App) PriceAlertGetAllAlerts() *controllers.GetAlertsResponse {
	if a.PriceAlertController == nil {
		return &controllers.GetAlertsResponse{Success: false, Message: "价格预警模块未初始化（PriceAlertController=nil）", Alerts: nil}
	}
	return a.PriceAlertController.GetAllAlerts()
}

func (a *App) PriceAlertGetActiveAlerts() *controllers.GetAlertsResponse {
	if a.PriceAlertController == nil {
		return &controllers.GetAlertsResponse{Success: false, Message: "价格预警模块未初始化（PriceAlertController=nil）", Alerts: nil}
	}
	return a.PriceAlertController.GetActiveAlerts()
}

func (a *App) PriceAlertGetAlertsByStockCode(stockCode string) *controllers.GetAlertsResponse {
	if a.PriceAlertController == nil {
		return &controllers.GetAlertsResponse{Success: false, Message: "价格预警模块未初始化（PriceAlertController=nil）", Alerts: nil}
	}
	return a.PriceAlertController.GetAlertsByStockCode(stockCode)
}

func (a *App) PriceAlertGetTriggerHistory(stockCode string, limit int) *controllers.GetTriggerHistoryResponse {
	if a.PriceAlertController == nil {
		return &controllers.GetTriggerHistoryResponse{Success: false, Message: "价格预警模块未初始化（PriceAlertController=nil）", Histories: nil}
	}
	return a.PriceAlertController.GetTriggerHistory(stockCode, limit)
}

func (a *App) PriceAlertGetAllTemplates() *controllers.GetTemplatesResponse {
	if a.PriceAlertController == nil {
		return &controllers.GetTemplatesResponse{Success: false, Message: "价格预警模块未初始化（PriceAlertController=nil）", Templates: nil}
	}
	return a.PriceAlertController.GetAllTemplates()
}

func (a *App) PriceAlertCreateAlert(jsonData string) *controllers.CreateAlertResponse {
	if a.PriceAlertController == nil {
		return &controllers.CreateAlertResponse{Success: false, Message: "价格预警模块未初始化（PriceAlertController=nil）"}
	}
	return a.PriceAlertController.CreateAlert(jsonData)
}

func (a *App) PriceAlertUpdateAlert(jsonData string) *controllers.CreateAlertResponse {
	if a.PriceAlertController == nil {
		return &controllers.CreateAlertResponse{Success: false, Message: "价格预警模块未初始化（PriceAlertController=nil）"}
	}
	return a.PriceAlertController.UpdateAlert(jsonData)
}

func (a *App) PriceAlertDeleteAlert(id int64) *controllers.CreateAlertResponse {
	if a.PriceAlertController == nil {
		return &controllers.CreateAlertResponse{Success: false, Message: "价格预警模块未初始化（PriceAlertController=nil）"}
	}
	return a.PriceAlertController.DeleteAlert(id)
}

func (a *App) PriceAlertToggleAlertStatus(id int64, isActive bool) *controllers.ToggleAlertStatusResponse {
	if a.PriceAlertController == nil {
		return &controllers.ToggleAlertStatusResponse{Success: false, Message: "价格预警模块未初始化（PriceAlertController=nil）"}
	}
	return a.PriceAlertController.ToggleAlertStatus(id, isActive)
}

func (a *App) PriceAlertCreateAlertFromTemplate(templateID, stockCode, stockName string, paramsJSON string) *controllers.CreateAlertFromTemplateResponse {
	if a.PriceAlertController == nil {
		return &controllers.CreateAlertFromTemplateResponse{Success: false, Message: "价格预警模块未初始化（PriceAlertController=nil）"}
	}
	return a.PriceAlertController.CreateAlertFromTemplate(templateID, stockCode, stockName, paramsJSON)
}

// --- PriceAlertController 转发器结束 ---

// MathAbs 辅助函数
func MathAbs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

// initAIService 初始化或重新初始化 AI 服务
func (a *App) initAIService() error {
	if a.dbService == nil {
		a.aiInitErr = fmt.Errorf("数据库服务未就绪，无法加载 AI 配置")
		return a.aiInitErr
	}

	start := time.Now()
	configSvc := services.NewConfigService(repositories.NewSQLiteConfigRepository(a.dbService.GetDB()))
	cfg, err := configSvc.LoadAIConfig()
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

// GetConfig Wails 绑定方法：获取配置 (已废弃，由 ConfigController 替代)
// func (a *App) GetConfig() (services.AIResolvedConfig, error) {
// 	return a.configService.LoadAIConfig()
// }

// SaveConfig Wails 绑定方法：保存配置 (已废弃，由 ConfigController 替代)
// func (a *App) SaveConfig(config services.AIResolvedConfig) error {
// 	err := a.configService.SaveAIConfig(config)
// 	if err != nil {
// 		return err
// 	}
// 	return a.initAIService()
// }

// AnalyzePastSignals 分析历史信号表现
func (a *App) AnalyzePastSignals(days int) (*models.SignalAnalysisResult, error) {
	if a.backtestService == nil {
		return nil, fmt.Errorf("回测服务未初始化")
	}
	return a.backtestService.AnalyzePastSignals(days)
}

// GetStockData 获取股票数据
func (a *App) GetStockData(code string) (*models.StockData, error) {
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	return a.stockService.GetStockByCode(code)
}

// GetStockDetail 获取个股详情页所需的所有数据
func (a *App) GetStockDetail(code string) (*models.StockDetail, error) {
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	return a.stockService.GetStockDetail(code)
}

// GetIntradayData 获取分时数据
func (a *App) GetIntradayData(code string) (*models.IntradayResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	return a.stockService.GetIntradayData(code)
}

// StreamIntradayData 启动分时 SSE 流并通过 Wails Events 推送到前端
// 前端监听事件名：intradayDataUpdate:{code}
func (a *App) StreamIntradayData(code string) {
	if code == "" {
		return
	}
	a.stockService.StreamIntradayData(code)
}

// StopIntradayStream 停止分时 SSE 流
func (a *App) StopIntradayStream(code string) {
	if code == "" {
		return
	}
	a.stockService.StopIntradayStream(code)
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

// --- Position 转发器 ---
// AddPosition 添加持仓记录
func (a *App) AddPosition(pos models.Position) error {
	return a.PositionController.AddPosition(pos)
}

// GetPositions 获取所有活跃持仓
func (a *App) GetPositions() (map[string]*models.Position, error) {
	return a.PositionController.GetPositions()
}

// RemovePosition 移除持仓记录
func (a *App) RemovePosition(code string) error {
	return a.PositionController.RemovePosition(code)
}

// --- Position 转发器 结束 ---

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

	// 注入全局默认配置
	if res != nil && a.dbService != nil {
		configSvc := services.NewConfigService(repositories.NewSQLiteConfigRepository(a.dbService.GetDB()))
		globalConfig, err := configSvc.GetGlobalStrategyConfig()
		if err == nil {
			// 将全局配置注入到结果中
			if res.TrailingStopConfig == nil {
				res.TrailingStopConfig = &models.TrailingStopConfig{}
			}
			// 如果用户没有手动指定，使用全局默认值
			if res.TrailingStopConfig.ActivationThreshold == 0 {
				res.TrailingStopConfig.ActivationThreshold = globalConfig.TrailingStopActivation
			}
			if res.TrailingStopConfig.CallbackRate == 0 {
				res.TrailingStopConfig.CallbackRate = globalConfig.TrailingStopCallback
			}
			logger.Debug("已注入全局默认配置",
				zap.String("module", "app.entry_strategy"),
				zap.Float64("activation_threshold", res.TrailingStopConfig.ActivationThreshold),
				zap.Float64("callback_rate", res.TrailingStopConfig.CallbackRate),
			)
		}
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

// --- Watchlist 转发器 ---
// AddToWatchlist 添加到自选股
func (a *App) AddToWatchlist(stock models.StockData) error {
	return a.WatchlistController.AddToWatchlist(stock)
}

// RemoveFromWatchlist 从自选股移除
func (a *App) RemoveFromWatchlist(code string) error {
	return a.WatchlistController.RemoveFromWatchlist(code)
}

// GetWatchlist 获取自选股列表
func (a *App) GetWatchlist() ([]*models.StockData, error) {
	return a.WatchlistController.GetWatchlist()
}

// --- Watchlist 转发器 结束 ---

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

// --- Config 转发器 ---
// GetConfig 获取 AI 配置
func (a *App) GetConfig() (services.AIResolvedConfig, error) {
	return a.ConfigController.GetAIConfig()
}

// SaveConfig 保存 AI 配置
func (a *App) SaveConfig(config services.AIResolvedConfig) error {
	return a.ConfigController.SaveAIConfig(config)
}

// GetGlobalStrategyConfig 获取全局交易策略配置
func (a *App) GetGlobalStrategyConfig() (services.GlobalStrategyConfig, error) {
	return a.ConfigController.GetGlobalStrategyConfig()
}

// UpdateGlobalStrategyConfig 更新全局交易策略配置
func (a *App) UpdateGlobalStrategyConfig(config services.GlobalStrategyConfig) error {
	return a.ConfigController.UpdateGlobalStrategyConfig(config)
}

// --- Config 转发器 结束 ---

// RunStrategyScan 运行策略扫描
func (a *App) RunStrategyScan(codes []string) []map[string]interface{} {
	if a.strategyService == nil {
		logger.Error("策略服务未初始化")
		return nil
	}

	var results []map[string]interface{}

	for _, code := range codes {
		signal, err := a.strategyService.CalculateBuildSignals(code)
		if err != nil {
			logger.Error("策略计算失败", zap.String("code", code), zap.Error(err))
			continue
		}

		if signal != nil {
			// 触发 AI 验证 (异步)
			if a.aiService != nil {
				go func(sig *models.StrategySignal) {
					// 1. 获取最近 7 天资金流向
					flows, err := a.strategyService.GetRecentMoneyFlows(sig.Code, 7)
					if err != nil {
						logger.Error("获取近期资金流向失败", zap.String("code", sig.Code), zap.Error(err))
						return
					}

					// 2. 获取股票基本信息
					stock, err := a.stockService.GetStockByCode(sig.Code)
					if err != nil {
						logger.Error("获取股票信息失败", zap.String("code", sig.Code), zap.Error(err))
						// 构建一个临时的 StockData
						stock = &models.StockData{Code: sig.Code, Name: sig.Code}
					}

					// 3. 调用 AI 验证
					verifyChan := a.aiService.VerifySignalAsync(stock, flows)
					res := <-verifyChan

					if res != nil {
						// 4. 更新数据库
						err := a.strategyService.UpdateSignalAIResult(sig.Code, sig.TradeDate, sig.StrategyName, res.Score, res.Opinion)
						if err != nil {
							logger.Error("更新 AI 结果失败", zap.Error(err))
						}

						// 5. 通知前端
						signalData := map[string]interface{}{
							"code":         sig.Code,
							"tradeDate":    sig.TradeDate,
							"signalType":   sig.SignalType,
							"score":        sig.Score,
							"strategyName": sig.StrategyName,
							"aiScore":      res.Score,
							"aiReason":     res.Opinion,
							"riskLevel":    res.RiskLevel,
						}
						runtime.EventsEmit(a.ctx, "signal_verified", signalData)
						runtime.EventsEmit(a.ctx, "new_signal", signalData)
					}
				}(signal)
			}

			results = append(results, map[string]interface{}{
				"code":         signal.Code,
				"tradeDate":    signal.TradeDate,
				"signalType":   signal.SignalType,
				"score":        signal.Score,
				"details":      signal.Details,
				"strategyName": signal.StrategyName,
			})
		}
	}

	return results
}

// StartMassScan 启动全市场策略扫描
// 前端调用此方法后会立即返回，扫描过程在后台进行，通过事件推送进度
func (a *App) StartMassScan() {
	go func() {
		logger.Info("启动全市场扫描任务")

		// 1. 获取所有待扫描股票
		// 这里我们使用 stockMarketService 获取股票代码，而不是 stockService
		// 因为 stockMarketService 包含完整的市场股票列表
		if a.StockMarketController == nil || a.stockService == nil {
			runtime.EventsEmit(a.ctx, "scan_error", "市场股票服务未初始化")
			return
		}

		// 通过 service 直接获取，避免 controller 的封装
		// 我们需要在 App 结构体中添加 stockMarketService 字段的直接访问，或者通过 Controller 获取
		// 这里假设 App 结构体中有 stockMarketService 字段，但实际上是通过 NewApp 注入到了 Controller
		// 我们需要修改 App 结构体或者直接使用 stockMarketCtrl 对应的 service
		// 查看 App 结构体，发现没有直接保存 stockMarketService 的引用，只在 NewApp 局部变量里
		// 所以我们需要先解决这个问题，或者暂时通过数据库直接查询

		// 修正：App 结构体实际上没有保存 stockMarketService，只保存了 StockMarketController
		// 但 NewApp 中确实初始化了 stockMarketSvc 并传给了 StockMarketController
		// 为了简单起见，我们可以在 App 中增加一个 stockMarketService 字段，
		// 或者直接在 NewApp 中把 stockMarketSvc 赋值给 App 的新字段。
		// 不过，既然我们已经在 StockMarketService 中添加了 GetAllStockCodes，
		// 我们最好是在 App 结构体中添加 stockMarketService 字段。
		// 考虑到无法修改结构体定义（需要修改文件头部），我们尝试通过 StockMarketController 调用，
		// 但 StockMarketController 可能没有暴露这个方法。

		// 既然我们已经在前面的步骤中修改了 services/stock_market_service.go，
		// 我们可以尝试通过 a.dbService 直接查询，但这重复了逻辑。
		// 最好的办法是修改 App 结构体，添加 stockMarketService *services.StockMarketService 字段。
		// 但由于我只能通过 SearchReplace 修改文件，添加字段比较麻烦。

		// 替代方案：在 RunStrategyScan 中我们已经有现成的逻辑。
		// 我们可以通过 dbService 获取所有股票代码。

		// 让我们先尝试使用 dbService 获取所有代码，模拟 stockMarketService.GetAllStockCodes 的逻辑
		if a.dbService == nil {
			runtime.EventsEmit(a.ctx, "scan_error", "数据库服务未初始化")
			return
		}

		db := a.dbService.GetDB()
		rows, err := db.Query("SELECT code FROM stocks WHERE is_active = 1 ORDER BY code ASC")
		if err != nil {
			logger.Error("获取股票代码失败", zap.Error(err))
			runtime.EventsEmit(a.ctx, "scan_error", fmt.Sprintf("获取股票列表失败: %v", err))
			return
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

		total := len(codes)
		foundCount := 0

		// 发送扫描开始事件
		runtime.EventsEmit(a.ctx, "scan_start", map[string]interface{}{
			"total": total,
		})

		logger.Info("开始扫描股票", zap.Int("total", total))

		// 2. 遍历扫描
		for i, code := range codes {
			// 检查上下文是否已取消（程序退出）
			select {
			case <-a.ctx.Done():
				return
			default:
			}

			// 计算策略信号
			signal, err := a.strategyService.CalculateBuildSignals(code)
			if err != nil {
				// 单个失败不中断整体扫描，仅记录日志
				// logger.Debug("策略计算跳过", zap.String("code", code), zap.Error(err))
			}

			// 3. 发现信号后的处理
			if signal != nil {
				foundCount++

				// 立即发送基础信号发现事件
				runtime.EventsEmit(a.ctx, "scan_signal_found", signal)

				// 触发 AI 深度验证 (异步)
				if a.aiService != nil {
					go func(sig *models.StrategySignal) {
						// 获取辅助数据
						flows, _ := a.strategyService.GetRecentMoneyFlows(sig.Code, 7)
						stock, err := a.stockService.GetStockByCode(sig.Code)
						if err != nil {
							stock = &models.StockData{Code: sig.Code, Name: sig.Code}
						}

						// 调用 AI 分析
						verifyChan := a.aiService.VerifySignalAsync(stock, flows)
						res := <-verifyChan

						if res != nil {
							// 更新 AI 评分结果
							_ = a.strategyService.UpdateSignalAIResult(sig.Code, sig.TradeDate, sig.StrategyName, res.Score, res.Opinion)

							// 组装完整数据推送到前端
							signalData := map[string]interface{}{
								"code":         sig.Code,
								"tradeDate":    sig.TradeDate,
								"signalType":   sig.SignalType,
								"score":        sig.Score,
								"strategyName": sig.StrategyName,
								"aiScore":      res.Score,
								"aiReason":     res.Opinion,
								"riskLevel":    res.RiskLevel,
								"details":      sig.Details,
							}

							// 推送 AI 验证完成事件
							runtime.EventsEmit(a.ctx, "signal_verified", signalData)
							// 兼容旧的信号事件
							runtime.EventsEmit(a.ctx, "new_signal", signalData)
						}
					}(signal)
				} else {
					// 无 AI 服务时，直接推送原始信号
					runtime.EventsEmit(a.ctx, "new_signal", map[string]interface{}{
						"code":         signal.Code,
						"tradeDate":    signal.TradeDate,
						"signalType":   signal.SignalType,
						"score":        signal.Score,
						"strategyName": signal.StrategyName,
						"details":      signal.Details,
					})
				}
			}

			// 4. 发送进度事件 (每 20 个或最后一个发送一次，减少前端渲染压力)
			if (i+1)%20 == 0 || i == total-1 {
				runtime.EventsEmit(a.ctx, "scan_progress", map[string]interface{}{
					"current":  i + 1,
					"total":    total,
					"found":    foundCount,
					"lastCode": code,
				})
			}

			// 简单的限流，防止瞬间 CPU 占用过高
			if i%100 == 0 {
				time.Sleep(5 * time.Millisecond)
			}
		}

		// 5. 扫描完成
		logger.Info("全市场扫描完成", zap.Int("total", total), zap.Int("found", foundCount))
		runtime.EventsEmit(a.ctx, "scan_complete", map[string]interface{}{
			"total": total,
			"found": foundCount,
		})
	}()
}

// GetLatestSignals 获取最新的策略信号
func (a *App) GetLatestSignals(limit int) ([]models.StrategySignal, error) {
	if a.strategyService == nil {
		return nil, fmt.Errorf("策略服务未初始化")
	}
	return a.strategyService.GetLatestSignals(limit)
}

// GetSignalsByStockCode 根据股票代码获取历史信号
func (a *App) GetSignalsByStockCode(code string) ([]models.StrategySignal, error) {
	if a.strategyService == nil {
		return nil, fmt.Errorf("策略服务未初始化")
	}
	return a.strategyService.GetSignalsByStockCode(code)
}

// shutdown 在应用程序退出时调用
func (a *App) shutdown(ctx context.Context) {
	// 关闭数据库连接
	if a.dbService != nil {
		a.dbService.Close()
	}
}

// --- 数据同步功能 开始 ---

// SyncStockData 同步单个股票的历史数据到本地 SQLite
func (a *App) SyncStockData(code string, startDate string, endDate string) (*models.SyncResult, error) {
	startTime := time.Now()

	// 调用 stockService 同步数据
	result, err := a.stockService.SyncStockData(code, startDate, endDate)

	// 保存同步历史记录
	duration := int(time.Since(startTime).Seconds())
	stockName := ""

	// 尝试获取股票名称
	if stock, err := a.stockService.GetStockByCode(code); err == nil {
		stockName = stock.Name
	}

	history := models.SyncHistory{
		StockCode:      code,
		StockName:      stockName,
		SyncType:       "single",
		StartDate:      startDate,
		EndDate:        endDate,
		Status:         "success",
		RecordsAdded:   result.RecordsAdded,
		RecordsUpdated: result.RecordsUpdated,
		Duration:       duration,
		CreatedAt:      time.Now(),
	}

	if err != nil {
		history.Status = "failed"
		history.ErrorMsg = result.ErrorMessage
	}

	// 异步保存历史记录，避免影响同步性能
	if a.syncHistoryCtrl != nil {
		go func() {
			if saveErr := a.syncHistoryCtrl.AddSyncHistory(history); saveErr != nil {
				logger.Error("保存同步历史记录失败",
					zap.String("stock_code", code),
					zap.Error(saveErr),
				)
			}
		}()
	}

	return result, err
}

// GetDataSyncStats 获取数据同步统计信息
func (a *App) GetDataSyncStats() (*models.DataSyncStats, error) {
	return a.stockService.GetDataSyncStats()
}

// BatchSyncStockData 批量同步多个股票的历史数据
func (a *App) BatchSyncStockData(codes []string, startDate string, endDate string) error {
	// 直接调用 stockService 批量同步
	// 每个股票的同步历史记录会在 SyncStockData 中单独记录
	return a.stockService.BatchSyncStockData(codes, startDate, endDate)
}

// ClearStockCache 清除指定股票的本地缓存数据
func (a *App) ClearStockCache(code string) error {
	return a.stockService.ClearStockCache(code)
}

// --- 数据同步功能 结束 ---

// --- 回测功能 结束 ---

// --- SyncHistoryController 转发方法 ---

// GetAllSyncHistory 获取所有同步历史记录（分页）
func (a *App) GetAllSyncHistory(limit int, offset int) ([]*models.SyncHistory, error) {
	return a.SyncHistoryController.GetAllSyncHistory(limit, offset)
}

// GetSyncHistoryByCode 根据股票代码获取同步历史记录
func (a *App) GetSyncHistoryByCode(code string, limit int) ([]*models.SyncHistory, error) {
	return a.SyncHistoryController.GetSyncHistoryByCode(code, limit)
}

// GetSyncHistoryCount 获取同步历史记录总数
func (a *App) GetSyncHistoryCount() (int, error) {
	return a.SyncHistoryController.GetSyncHistoryCount()
}

// ClearAllSyncHistory 清除所有同步历史记录
func (a *App) ClearAllSyncHistory() error {
	return a.SyncHistoryController.ClearAllSyncHistory()
}

// GetSyncedKLineDataResponse GetSyncedKLineData 的响应结构
type GetSyncedKLineDataResponse struct {
	Data  []map[string]interface{} `json:"data"`
	Total int                      `json:"total"`
}

// GetSyncedKLineData 获取指定股票已同步的K线数据（支持分页和日期筛选）
func (a *App) GetSyncedKLineData(code string, startDate string, endDate string, page int, pageSize int) GetSyncedKLineDataResponse {
	// 打印调用参数
	logger.Info("GetSyncedKLineData 被调用",
		zap.String("code", code),
		zap.String("startDate", startDate),
		zap.String("endDate", endDate),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
	)

	// 初始化返回值
	response := GetSyncedKLineDataResponse{
		Data:  []map[string]interface{}{},
		Total: 0,
	}

	// 获取数据库服务实例
	db := a.dbService
	if db == nil {
		logger.Error("数据库服务未初始化")
		return response
	}

	// 调用数据库服务查询K线数据
	data, total, err := db.GetKLineDataWithPagination(code, startDate, endDate, page, pageSize)

	if err != nil {
		logger.Error("获取K线数据失败", zap.Error(err))
		return response
	}

	// 确保 data 不是 nil
	if data == nil {
		logger.Warn("GetKLineDataWithPagination 返回了 nil，初始化为空数组")
		data = []map[string]interface{}{}
	}

	// 打印返回结果
	logger.Info("GetSyncedKLineData 返回结果",
		zap.Int("dataLength", len(data)),
		zap.Int("total", total),
	)

	// 如果有数据，打印第一条数据用于调试
	if len(data) > 0 {
		logger.Debug("第一条K线数据", zap.Any("firstItem", data[0]))
	}

	response.Data = data
	response.Total = total
	return response
}

// ============ 策略管理 API ============

// CreateStrategy 创建策略
func (a *App) CreateStrategy(name string, description string, strategyType string, parameters map[string]interface{}) error {
	if a.StrategyController == nil {
		return fmt.Errorf("策略控制器未初始化")
	}
	return a.StrategyController.CreateStrategy(name, description, strategyType, parameters)
}

// UpdateStrategy 更新策略
func (a *App) UpdateStrategy(id int64, name string, description string, strategyType string, parameters map[string]interface{}) error {
	if a.StrategyController == nil {
		return fmt.Errorf("策略控制器未初始化")
	}
	return a.StrategyController.UpdateStrategy(id, name, description, strategyType, parameters)
}

// DeleteStrategy 删除策略
func (a *App) DeleteStrategy(id int64) error {
	if a.StrategyController == nil {
		return fmt.Errorf("策略控制器未初始化")
	}
	return a.StrategyController.DeleteStrategy(id)
}

// GetStrategy 获取策略
func (a *App) GetStrategy(id int64) (interface{}, error) {
	if a.StrategyController == nil {
		return nil, fmt.Errorf("策略控制器未初始化")
	}
	return a.StrategyController.GetStrategy(id)
}

// GetAllStrategies 获取所有策略
func (a *App) GetAllStrategies() (interface{}, error) {
	if a.StrategyController == nil {
		return nil, fmt.Errorf("策略控制器未初始化")
	}
	return a.StrategyController.GetAllStrategies()
}

// GetStrategyTypes 获取所有策略类型
func (a *App) GetStrategyTypes() interface{} {
	if a.StrategyController == nil {
		return []interface{}{}
	}
	return a.StrategyController.GetStrategyTypes()
}

// UpdateStrategyBacktestResult 更新策略回测结果
func (a *App) UpdateStrategyBacktestResult(id int64, backtestResult map[string]interface{}) error {
	if a.StrategyController == nil {
		return fmt.Errorf("策略控制器未初始化")
	}
	return a.StrategyController.UpdateStrategyBacktestResult(id, backtestResult)
}

// ============ 市场股票管理 API ============

// SyncAllStocks 同步所有市场股票
func (a *App) SyncAllStocks() (interface{}, error) {
	if a.StockMarketController == nil {
		return nil, fmt.Errorf("市场股票控制器未初始化")
	}
	return a.StockMarketController.SyncAllStocks()
}

// GetStocksList 获取股票列表
func (a *App) GetStocksList(page int, pageSize int, search string, industry string) (interface{}, error) {
	if a.StockMarketController == nil {
		return nil, fmt.Errorf("市场股票控制器未初始化")
	}
	return a.StockMarketController.GetStocksList(page, pageSize, search, industry)
}

// GetIndustries 获取行业列表
func (a *App) GetIndustries() (interface{}, error) {
	if a.StockMarketController == nil {
		return nil, fmt.Errorf("市场股票控制器未初始化")
	}
	return a.StockMarketController.GetIndustries()
}

// GetSyncStats 获取同步统计信息
func (a *App) GetSyncStats() (interface{}, error) {
	if a.StockMarketController == nil {
		return nil, fmt.Errorf("市场股票控制器未初始化")
	}
	return a.StockMarketController.GetSyncStats()
}

// ============ K线数据同步 API ============

// StartKLineSync 开始K线数据同步
func (a *App) StartKLineSync(days int) (interface{}, error) {
	if a.klineSyncService == nil {
		return nil, fmt.Errorf("K线同步服务未初始化")
	}
	return a.klineSyncService.StartKLineSync(days)
}

// GetKLineSyncProgress 获取K线同步进度
func (a *App) GetKLineSyncProgress() (interface{}, error) {
	if a.klineSyncService == nil {
		return nil, fmt.Errorf("K线同步服务未初始化")
	}
	return a.klineSyncService.GetSyncProgress()
}

// GetKLineSyncHistory 获取K线同步历史记录
func (a *App) GetKLineSyncHistory(limit int) (interface{}, error) {
	if a.klineSyncService == nil {
		return nil, fmt.Errorf("K线同步服务未初始化")
	}
	return a.klineSyncService.GetKLineSyncHistory(limit)
}

// StartFullMarketSync 启动全市场资金流同步
func (a *App) StartFullMarketSync() error {
	if a.syncService == nil {
		return fmt.Errorf("全量同步服务未初始化")
	}
	return a.syncService.StartFullMarketSync()
}
