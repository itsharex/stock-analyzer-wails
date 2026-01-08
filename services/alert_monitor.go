package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

// AlertMonitor 价格预警监控引擎
type AlertMonitor struct {
	ctx              context.Context
	priceAlertSvc    *PriceAlertService
	repo             *repositories.PriceAlertRepository
	stockService     StockDataService
	klineService     KLineDataService
	ticker           *time.Ticker
	mu               sync.Mutex
	running          bool
	checkInterval    time.Duration
	onAlertTriggered func(alert *repositories.PriceThresholdAlert, stockData *StockDataForAlert, message string)
}

// StockDataService 股票数据服务接口
type StockDataService interface {
	GetStockByCode(code string) (*models.StockData, error)
}

// KLineDataService K线数据服务接口（用于获取MA数据）
type KLineDataService interface {
	GetKLineData(code string, count int, period string) ([]*models.KLineData, error)
}

// NewAlertMonitor 创建价格预警监控引擎
func NewAlertMonitor(ctx context.Context, priceAlertSvc *PriceAlertService, stockService StockDataService, klineService KLineDataService) *AlertMonitor {
	return &AlertMonitor{
		ctx:              ctx,
		priceAlertSvc:    priceAlertSvc,
		repo:             priceAlertSvc.GetRepository(),
		stockService:     stockService,
		klineService:     klineService,
		checkInterval:    10 * time.Second, // 默认每10秒检查一次
		onAlertTriggered: nil,
	}
}

// SetCheckInterval 设置检查间隔
func (m *AlertMonitor) SetCheckInterval(interval time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checkInterval = interval
}

// SetAlertTriggerCallback 设置预警触发回调
func (m *AlertMonitor) SetAlertTriggerCallback(callback func(alert *repositories.PriceThresholdAlert, stockData *StockDataForAlert, message string)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onAlertTriggered = callback
}

// convertToStockDataForAlert 将 models.StockData 转换为 StockDataForAlert
func (m *AlertMonitor) convertToStockDataForAlert(stockData *models.StockData, klineData []*models.KLineData) (*StockDataForAlert, error) {
	if stockData == nil {
		return nil, fmt.Errorf("股票数据为空")
	}

	result := &StockDataForAlert{
		Code:               stockData.Code,
		Name:               stockData.Name,
		ClosePrice:         stockData.Price,
		OpenPrice:          stockData.Open,
		HighPrice:          stockData.High,
		LowPrice:           stockData.Low,
		PreClosePrice:      stockData.PreClose,
		PriceChangePercent: stockData.ChangeRate,
		Volume:             stockData.Volume,
		VolumeRatio:        stockData.VolumeRatio,
		MA5:                0,
		MA10:               0,
		MA20:               0,
		HistoricalHigh:     0,
		HistoricalLow:      0,
	}

	// 计算均线和历史高低点（如果有K线数据）
	if len(klineData) > 0 {
		// 计算历史最高价和最低价
		for _, k := range klineData {
			if k.High > result.HistoricalHigh {
				result.HistoricalHigh = k.High
			}
			if result.HistoricalLow == 0 || k.Low < result.HistoricalLow {
				result.HistoricalLow = k.Low
			}
		}

		// 计算MA5、MA10、MA20
		result.MA5 = calculateMA(klineData, 5)
		result.MA10 = calculateMA(klineData, 10)
		result.MA20 = calculateMA(klineData, 20)
	}

	return result, nil
}

// calculateMA 计算移动平均线
func calculateMA(klineData []*models.KLineData, period int) float64 {
	if len(klineData) < period {
		return 0
	}

	sum := 0.0
	count := 0
	// 从最新的数据开始计算
	for i := len(klineData) - 1; i >= 0 && count < period; i-- {
		sum += klineData[i].Close
		count++
	}

	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// Start 启动预警监控
func (m *AlertMonitor) Start() {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return
	}
	m.running = true
	m.mu.Unlock()

	logger.Info("价格预警监控引擎已启动", zap.Duration("interval", m.checkInterval))

	m.ticker = time.NewTicker(m.checkInterval)
	go m.run()
}

// Stop 停止预警监控
func (m *AlertMonitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return
	}

	m.running = false
	if m.ticker != nil {
		m.ticker.Stop()
		m.ticker = nil
	}

	logger.Info("价格预警监控引擎已停止")
}

// run 运行监控循环
func (m *AlertMonitor) run() {
	for {
		select {
		case <-m.ctx.Done():
			m.Stop()
			return
		case <-m.ticker.C:
			m.checkAllAlerts()
		}
	}
}

// checkAllAlerts 检查所有活跃预警
func (m *AlertMonitor) checkAllAlerts() {
	// 获取所有活跃预警
	alerts, err := m.repo.GetActiveAlerts()
	if err != nil {
		logger.Error("获取活跃预警失败", zap.Error(err))
		return
	}

	if len(alerts) == 0 {
		return
	}

	logger.Debug("开始检查活跃预警", zap.Int("count", len(alerts)))

	// 按股票代码分组，批量获取股票数据
	stockCodes := make(map[string]bool)
	for _, alert := range alerts {
		stockCodes[alert.StockCode] = true
	}

	// 获取股票数据并转换
	stockDataMap := make(map[string]*StockDataForAlert)
	for code := range stockCodes {
		// 获取实时行情数据
		stockData, err := m.stockService.GetStockByCode(code)
		if err != nil {
			logger.Warn("获取股票数据失败", zap.String("code", code), zap.Error(err))
			continue
		}

		// 获取K线数据（用于计算MA和历史高低点）
		var klineData []*models.KLineData
		if m.klineService != nil {
			klineData, err = m.klineService.GetKLineData(code, 100, "daily")
			if err != nil {
				logger.Warn("获取K线数据失败", zap.String("code", code), zap.Error(err))
				// K线数据获取失败不影响预警检测，继续使用实时数据
			}
		}

		// 转换为预警所需的数据格式
		alertData, err := m.convertToStockDataForAlert(stockData, klineData)
		if err != nil {
			logger.Error("转换股票数据失败", zap.String("code", code), zap.Error(err))
			continue
		}

		stockDataMap[code] = alertData
	}

	// 检查每个预警
	for _, alert := range alerts {
		stockData, exists := stockDataMap[alert.StockCode]
		if !exists {
			continue
		}

		// 检查预警是否触发
		triggered, message, err := m.priceAlertSvc.CheckAlert(alert, stockData)
		if err != nil {
			logger.Error("检查预警失败",
				zap.Int64("alertId", alert.ID),
				zap.String("stockCode", alert.StockCode),
				zap.Error(err))
			continue
		}

		if triggered {
			logger.Info("价格预警触发",
				zap.Int64("alertId", alert.ID),
				zap.String("stockCode", alert.StockCode),
				zap.String("stockName", alert.StockName),
				zap.String("message", message))

			// 处理触发的预警
			m.handleTriggeredAlert(alert, stockData, message)
		}
	}
}

// handleTriggeredAlert 处理触发的预警
func (m *AlertMonitor) handleTriggeredAlert(alert *repositories.PriceThresholdAlert, stockData *StockDataForAlert, message string) {
	// 1. 更新最后触发时间
	if err := m.repo.UpdateLastTriggeredTime(alert.ID); err != nil {
		logger.Error("更新预警触发时间失败", zap.Error(err))
	}

	// 2. 记录触发历史
	triggerHistory := &repositories.PriceAlertTriggerHistory{
		AlertID:        alert.ID,
		StockCode:      alert.StockCode,
		StockName:      alert.StockName,
		AlertType:      alert.AlertType,
		TriggerPrice:   stockData.ClosePrice,
		TriggerMessage: message,
		TriggeredAt:    time.Now(),
	}
	if err := m.repo.SaveTriggerHistory(triggerHistory); err != nil {
		logger.Error("保存预警触发历史失败", zap.Error(err))
	}

	// 3. 根据触发后行为处理
	switch alert.PostTriggerAction {
	case "disable":
		// 触发后禁用预警
		if err := m.repo.ToggleAlertStatus(alert.ID, false); err != nil {
			logger.Error("禁用预警失败", zap.Error(err))
		} else {
			logger.Info("预警触发后已禁用", zap.Int64("alertId", alert.ID))
		}
	case "once":
		// 仅触发一次（等同于禁用）
		if err := m.repo.ToggleAlertStatus(alert.ID, false); err != nil {
			logger.Error("禁用预警失败", zap.Error(err))
		} else {
			logger.Info("预警触发后已禁用（仅触发一次）", zap.Int64("alertId", alert.ID))
		}
	case "continue":
		// 继续监控（不执行额外操作）
		logger.Info("预警触发后继续监控", zap.Int64("alertId", alert.ID))
	default:
		logger.Warn("未知的触发后行为", zap.String("action", alert.PostTriggerAction))
	}

	// 4. 发送通知
	m.sendNotification(alert, stockData, message)
}

// sendNotification 发送预警通知
func (m *AlertMonitor) sendNotification(alert *repositories.PriceThresholdAlert, stockData *StockDataForAlert, message string) {
	// 构建通知数据
	notificationData := map[string]interface{}{
		"alertId":          alert.ID,
		"stockCode":        alert.StockCode,
		"stockName":        alert.StockName,
		"alertType":        alert.AlertType,
		"triggerPrice":     stockData.ClosePrice,
		"priceChange":      stockData.PriceChangePercent,
		"message":          message,
		"triggeredAt":      time.Now().Format("2006-01-02 15:04:05"),
		"enableSound":      alert.EnableSound,
		"enableDesktop":    alert.EnableDesktop,
	}

	// 调用回调函数
	m.mu.Lock()
	callback := m.onAlertTriggered
	m.mu.Unlock()

	if callback != nil {
		go callback(alert, stockData, message)
	}

	// 发送 Wails 事件通知前端
	runtime.EventsEmit(m.ctx, "price_alert_triggered", notificationData)

	logger.Info("已发送价格预警通知",
		zap.Int64("alertId", alert.ID),
		zap.String("stockCode", alert.StockCode),
		zap.Bool("enableDesktop", alert.EnableDesktop),
		zap.Bool("enableSound", alert.EnableSound))
}

// CheckStockAlerts 检查单个股票的所有预警（供外部手动触发）
func (m *AlertMonitor) CheckStockAlerts(stockCode string) error {
	// 获取该股票的活跃预警
	alerts, err := m.repo.GetAlertsByStockCode(stockCode)
	if err != nil {
		return fmt.Errorf("获取股票预警失败: %w", err)
	}

	// 过滤出活跃的预警
	var activeAlerts []*repositories.PriceThresholdAlert
	for _, alert := range alerts {
		if alert.IsActive {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	if len(activeAlerts) == 0 {
		return nil
	}

	// 获取实时行情数据
	stockData, err := m.stockService.GetStockByCode(stockCode)
	if err != nil {
		return fmt.Errorf("获取股票数据失败: %w", err)
	}

	// 获取K线数据（用于计算MA和历史高低点）
	var klineData []*models.KLineData
	if m.klineService != nil {
		klineData, err = m.klineService.GetKLineData(stockCode, 100, "daily")
		if err != nil {
			logger.Warn("获取K线数据失败", zap.String("code", stockCode), zap.Error(err))
		}
	}

	// 转换为预警所需的数据格式
	alertData, err := m.convertToStockDataForAlert(stockData, klineData)
	if err != nil {
		return fmt.Errorf("转换股票数据失败: %w", err)
	}

	// 检查每个预警
	for _, alert := range activeAlerts {
		triggered, message, err := m.priceAlertSvc.CheckAlert(alert, alertData)
		if err != nil {
			logger.Error("检查预警失败",
				zap.Int64("alertId", alert.ID),
				zap.Error(err))
			continue
		}

		if triggered {
			m.handleTriggeredAlert(alert, alertData, message)
		}
	}

	return nil
}

// IsRunning 返回监控是否正在运行
func (m *AlertMonitor) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}
