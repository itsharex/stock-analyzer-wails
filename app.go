package main

import (
	"context"
	"fmt"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/services"
	"sync"
	"time"

	"stock-analyzer-wails/internal/logger"

	"go.uber.org/zap"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 应用程序结构
type App struct {
	ctx              context.Context
	stockService     *services.StockService
	aiService        *services.AIService
	watchlistService *services.WatchlistService
	alertStorage     *services.AlertStorage
	aiInitErr        error
	alerts           []*models.PriceAlert
	alertMutex       sync.Mutex
	alertConfig      models.AlertConfig
}

// NewApp 创建新的App应用程序
func NewApp() *App {
	watchlistSvc, _ := services.NewWatchlistService()
	storage, _ := services.NewAlertStorage()
	return &App{
		stockService:     services.NewStockService(),
		aiService:        nil,
		watchlistService: watchlistSvc,
		alertStorage:     storage,
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

		// 碰撞检测逻辑
		triggered := false
		msg := ""
		
		// 阈值：距离关键位 0.5% 以内触发
		threshold := alert.Price * 0.005
		diff := stock.Price - alert.Price

		// 避免频繁触发：1小时冷却时间
		if time.Since(alert.LastTriggered) < 1*time.Hour {
			continue
		}

		if alert.Type == "resistance" && stock.Price >= alert.Price {
			triggered = true
			msg = fmt.Sprintf("突破压力位！%s 当前价 %.2f 已站上 AI 识别的压力位 %.2f", stock.Name, stock.Price, alert.Price)
		} else if alert.Type == "support" && stock.Price <= alert.Price {
			triggered = true
			msg = fmt.Sprintf("跌破支撑位！%s 当前价 %.2f 已跌穿 AI 识别的支撑位 %.2f", stock.Name, stock.Price, alert.Price)
		} else if diff > -threshold && diff < threshold {
			triggered = true
			msg = fmt.Sprintf("接近关键位！%s 当前价 %.2f 正在挑战 AI 识别的 %s %.2f", stock.Name, stock.Price, alert.Label, alert.Price)
		}

		if triggered {
			alert.LastTriggered = time.Now()
			
			// 异步生成 AI 建议，避免阻塞监控主循环
			go func(al *models.PriceAlert, st *models.StockData, baseMsg string) {
				advice, _ := a.aiService.GenerateAlertAdvice(st.Name, al.Type, al.Label, al.Role, st.Price, al.Price)
				
				// 发送系统通知
				runtime.EventsEmit(a.ctx, "price_alert", map[string]interface{}{
					"stockCode": al.StockCode,
					"stockName": al.StockName,
					"message":   baseMsg,
					"advice":    advice,
					"type":      al.Type,
					"price":     st.Price,
					"role":      al.Role,
				})
				
				logger.Info("触发价格预警",
					zap.String("stock", al.StockName),
					zap.Float64("price", st.Price),
					zap.String("msg", baseMsg),
					zap.String("advice", advice),
				)
			}(alert, stock, msg)
		}
	}
}

// UpdateAlertsFromAnalysis 从 AI 分析结果中更新预警位
func (a *App) UpdateAlertsFromAnalysis(code string, name string, result *models.TechnicalAnalysisResult, role string) {
	a.alertMutex.Lock()
	defer a.alertMutex.Unlock()

	logger.Info("开始更新预警位", 
		zap.String("code", code), 
		zap.Int("drawing_count", len(result.Drawings)),
		zap.String("role", role),
	)

	// 清除该股票旧的 AI 预警
	newAlerts := make([]*models.PriceAlert, 0)
	for _, alert := range a.alerts {
		if alert.StockCode != code {
			newAlerts = append(newAlerts, alert)
		}
	}

	// 添加新的预警位
	addedCount := 0
	for _, drawing := range result.Drawings {
		if (drawing.Type == "support" || drawing.Type == "resistance") && drawing.Price > 0 {
			newAlerts = append(newAlerts, &models.PriceAlert{
				StockCode: code,
				StockName: name,
				Type:      drawing.Type,
				Price:     drawing.Price,
				Label:     drawing.Label,
				Role:      role,
				IsActive:  true,
			})
			addedCount++
			logger.Info("成功添加预警位", 
				zap.String("type", drawing.Type), 
				zap.Float64("price", drawing.Price),
				zap.String("label", drawing.Label),
			)
		}
	}
	a.alerts = newAlerts
	
	// 持久化保存活跃预警
	if a.alertStorage != nil {
		a.alertStorage.SaveActiveAlerts(a.alerts)
	}
	
	logger.Info("预警位更新完成", zap.Int("added_count", addedCount), zap.Int("total_active", len(a.alerts)))
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

// GetKLineData 获取K线数据，支持周期参数
func (a *App) GetKLineData(code string, limit int, period string) ([]*models.KLineData, error) {
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	return a.stockService.GetKLineData(code, limit, period)
}

// AnalyzeStock 分析股票
func (a *App) AnalyzeStock(code string) (*models.AnalysisReport, error) {
	if a.aiService == nil {
		return nil, fmt.Errorf("AI服务未就绪，请检查配置: %v", a.aiInitErr)
	}
	
	stock, err := a.stockService.GetStockByCode(code)
	if err != nil {
		return nil, err
	}
	
	return a.aiService.AnalyzeStock(stock)
}

// AnalyzeTechnical 深度技术面分析（支持多角色切换和绘图数据）
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

	result, err := a.aiService.AnalyzeTechnical(stock, klines, role)
	if err != nil {
		return nil, err
	}

	// 自动更新该股票的预警位
	a.UpdateAlertsFromAnalysis(stock.Code, stock.Name, result, role)

	return result, nil
}

// SearchStock 搜索股票
func (a *App) SearchStock(keyword string) ([]*models.StockData, error) {
	return a.stockService.SearchStock(keyword)
}

// GetAlertHistory 获取告警历史
func (a *App) GetAlertHistory(stockCode string, limit int) ([]map[string]interface{}, error) {
	if a.alertStorage == nil {
		return []map[string]interface{}{}, nil
	}
	return a.alertStorage.GetAlertHistory(stockCode, limit)
}

// GetAlertConfig 获取预警配置
func (a *App) GetAlertConfig() models.AlertConfig {
	return a.alertConfig
}

// UpdateAlertConfig 更新预警配置
func (a *App) UpdateAlertConfig(config models.AlertConfig) {
	a.alertConfig = config
}

// GetActiveAlerts 获取当前所有活跃的预警订阅
func (a *App) GetActiveAlerts() []*models.PriceAlert {
	a.alertMutex.Lock()
	defer a.alertMutex.Unlock()
	return a.alerts
}

// RemoveAlert 移除指定的预警订阅
func (a *App) RemoveAlert(stockCode string, alertType string, price float64) {
	a.alertMutex.Lock()
	defer a.alertMutex.Unlock()
	
	newAlerts := []*models.PriceAlert{}
	for _, alert := range a.alerts {
		if alert.StockCode == stockCode && alert.Type == alertType && alert.Price == price {
			continue
		}
		newAlerts = append(newAlerts, alert)
	}
	a.alerts = newAlerts
	
	// 持久化保存更新后的预警列表
	if a.alertStorage != nil {
		a.alertStorage.SaveActiveAlerts(a.alerts)
	}
}

// Watchlist 相关接口

func (a *App) AddToWatchlist(stock *models.StockData) error {
	return a.watchlistService.AddToWatchlist(stock)
}

func (a *App) RemoveFromWatchlist(code string) error {
	return a.watchlistService.RemoveFromWatchlist(code)
}

func (a *App) GetWatchlist() ([]*models.StockData, error) {
	return a.watchlistService.GetWatchlist()
}
