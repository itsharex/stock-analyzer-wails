package services

import (
	"context"
	"reflect"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/repositories"
	"sync"
	"testing"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestAIService_AnalyzeEntryStrategy(t *testing.T) {
	type fields struct {
		chatModel    model.ChatModel
		config       AIResolvedConfig
		cacheService *AnalysisCacheService
		semaphore    chan struct{}
		enableMock   bool
	}
	type args struct {
		stock     *models.StockData
		klines    []*models.KLineData
		moneyFlow *models.MoneyFlowResponse
		health    *models.HealthCheckResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.EntryStrategyResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AIService{
				chatModel:    tt.fields.chatModel,
				config:       tt.fields.config,
				cacheService: tt.fields.cacheService,
				semaphore:    tt.fields.semaphore,
				enableMock:   tt.fields.enableMock,
			}
			got, err := s.AnalyzeEntryStrategy(tt.args.stock, tt.args.klines, tt.args.moneyFlow, tt.args.health)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeEntryStrategy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnalyzeEntryStrategy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIService_AnalyzeStock(t *testing.T) {
	type fields struct {
		chatModel    model.ChatModel
		config       AIResolvedConfig
		cacheService *AnalysisCacheService
		semaphore    chan struct{}
		enableMock   bool
	}
	type args struct {
		stock *models.StockData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.AnalysisReport
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AIService{
				chatModel:    tt.fields.chatModel,
				config:       tt.fields.config,
				cacheService: tt.fields.cacheService,
				semaphore:    tt.fields.semaphore,
				enableMock:   tt.fields.enableMock,
			}
			got, err := s.AnalyzeStock(tt.args.stock)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeStock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnalyzeStock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIService_AnalyzeTechnical(t *testing.T) {
	type fields struct {
		chatModel    model.ChatModel
		config       AIResolvedConfig
		cacheService *AnalysisCacheService
		semaphore    chan struct{}
		enableMock   bool
	}
	type args struct {
		stock  *models.StockData
		klines []*models.KLineData
		period string
		role   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.TechnicalAnalysisResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AIService{
				chatModel:    tt.fields.chatModel,
				config:       tt.fields.config,
				cacheService: tt.fields.cacheService,
				semaphore:    tt.fields.semaphore,
				enableMock:   tt.fields.enableMock,
			}
			got, err := s.AnalyzeTechnical(tt.args.stock, tt.args.klines, tt.args.period, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeTechnical() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnalyzeTechnical() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIService_GenerateAlertAdvice(t *testing.T) {
	type fields struct {
		chatModel    model.ChatModel
		config       AIResolvedConfig
		cacheService *AnalysisCacheService
		semaphore    chan struct{}
		enableMock   bool
	}
	type args struct {
		stockName    string
		alertType    string
		label        string
		role         string
		currentPrice float64
		alertPrice   float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AIService{
				chatModel:    tt.fields.chatModel,
				config:       tt.fields.config,
				cacheService: tt.fields.cacheService,
				semaphore:    tt.fields.semaphore,
				enableMock:   tt.fields.enableMock,
			}
			got, err := s.GenerateAlertAdvice(tt.args.stockName, tt.args.alertType, tt.args.label, tt.args.role, tt.args.currentPrice, tt.args.alertPrice)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAlertAdvice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateAlertAdvice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIService_SetEnableMock(t *testing.T) {
	type fields struct {
		chatModel    model.ChatModel
		config       AIResolvedConfig
		cacheService *AnalysisCacheService
		semaphore    chan struct{}
		enableMock   bool
	}
	type args struct {
		enable bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AIService{
				chatModel:    tt.fields.chatModel,
				config:       tt.fields.config,
				cacheService: tt.fields.cacheService,
				semaphore:    tt.fields.semaphore,
				enableMock:   tt.fields.enableMock,
			}
			s.SetEnableMock(tt.args.enable)
		})
	}
}

func TestAIService_VerifySignal(t *testing.T) {
	type fields struct {
		chatModel    model.ChatModel
		config       AIResolvedConfig
		cacheService *AnalysisCacheService
		semaphore    chan struct{}
		enableMock   bool
	}
	type args struct {
		stock       *models.StockData
		recentFlows []models.MoneyFlowData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.AIVerificationResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AIService{
				chatModel:    tt.fields.chatModel,
				config:       tt.fields.config,
				cacheService: tt.fields.cacheService,
				semaphore:    tt.fields.semaphore,
				enableMock:   tt.fields.enableMock,
			}
			got, err := s.VerifySignal(tt.args.stock, tt.args.recentFlows)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifySignal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifySignal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIService_VerifySignalAsync(t *testing.T) {
	type fields struct {
		chatModel    model.ChatModel
		config       AIResolvedConfig
		cacheService *AnalysisCacheService
		semaphore    chan struct{}
		enableMock   bool
	}
	type args struct {
		stock       *models.StockData
		recentFlows []models.MoneyFlowData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan *models.AIVerificationResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AIService{
				chatModel:    tt.fields.chatModel,
				config:       tt.fields.config,
				cacheService: tt.fields.cacheService,
				semaphore:    tt.fields.semaphore,
				enableMock:   tt.fields.enableMock,
			}
			if got := s.VerifySignalAsync(tt.args.stock, tt.args.recentFlows); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifySignalAsync() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIService_extractJSON(t *testing.T) {
	type fields struct {
		chatModel    model.ChatModel
		config       AIResolvedConfig
		cacheService *AnalysisCacheService
		semaphore    chan struct{}
		enableMock   bool
	}
	type args struct {
		content string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AIService{
				chatModel:    tt.fields.chatModel,
				config:       tt.fields.config,
				cacheService: tt.fields.cacheService,
				semaphore:    tt.fields.semaphore,
				enableMock:   tt.fields.enableMock,
			}
			if got := s.extractJSON(tt.args.content); got != tt.want {
				t.Errorf("extractJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlertMonitor_CheckStockAlerts(t *testing.T) {
	type fields struct {
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
	type args struct {
		stockCode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AlertMonitor{
				ctx:              tt.fields.ctx,
				priceAlertSvc:    tt.fields.priceAlertSvc,
				repo:             tt.fields.repo,
				stockService:     tt.fields.stockService,
				klineService:     tt.fields.klineService,
				ticker:           tt.fields.ticker,
				mu:               tt.fields.mu,
				running:          tt.fields.running,
				checkInterval:    tt.fields.checkInterval,
				onAlertTriggered: tt.fields.onAlertTriggered,
			}
			if err := m.CheckStockAlerts(tt.args.stockCode); (err != nil) != tt.wantErr {
				t.Errorf("CheckStockAlerts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlertMonitor_IsRunning(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AlertMonitor{
				ctx:              tt.fields.ctx,
				priceAlertSvc:    tt.fields.priceAlertSvc,
				repo:             tt.fields.repo,
				stockService:     tt.fields.stockService,
				klineService:     tt.fields.klineService,
				ticker:           tt.fields.ticker,
				mu:               tt.fields.mu,
				running:          tt.fields.running,
				checkInterval:    tt.fields.checkInterval,
				onAlertTriggered: tt.fields.onAlertTriggered,
			}
			if got := m.IsRunning(); got != tt.want {
				t.Errorf("IsRunning() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlertMonitor_SetAlertTriggerCallback(t *testing.T) {
	type fields struct {
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
	type args struct {
		callback func(alert *repositories.PriceThresholdAlert, stockData *StockDataForAlert, message string)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AlertMonitor{
				ctx:              tt.fields.ctx,
				priceAlertSvc:    tt.fields.priceAlertSvc,
				repo:             tt.fields.repo,
				stockService:     tt.fields.stockService,
				klineService:     tt.fields.klineService,
				ticker:           tt.fields.ticker,
				mu:               tt.fields.mu,
				running:          tt.fields.running,
				checkInterval:    tt.fields.checkInterval,
				onAlertTriggered: tt.fields.onAlertTriggered,
			}
			m.SetAlertTriggerCallback(tt.args.callback)
		})
	}
}

func TestAlertMonitor_SetCheckInterval(t *testing.T) {
	type fields struct {
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
	type args struct {
		interval time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AlertMonitor{
				ctx:              tt.fields.ctx,
				priceAlertSvc:    tt.fields.priceAlertSvc,
				repo:             tt.fields.repo,
				stockService:     tt.fields.stockService,
				klineService:     tt.fields.klineService,
				ticker:           tt.fields.ticker,
				mu:               tt.fields.mu,
				running:          tt.fields.running,
				checkInterval:    tt.fields.checkInterval,
				onAlertTriggered: tt.fields.onAlertTriggered,
			}
			m.SetCheckInterval(tt.args.interval)
		})
	}
}

func TestAlertMonitor_Start(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AlertMonitor{
				ctx:              tt.fields.ctx,
				priceAlertSvc:    tt.fields.priceAlertSvc,
				repo:             tt.fields.repo,
				stockService:     tt.fields.stockService,
				klineService:     tt.fields.klineService,
				ticker:           tt.fields.ticker,
				mu:               tt.fields.mu,
				running:          tt.fields.running,
				checkInterval:    tt.fields.checkInterval,
				onAlertTriggered: tt.fields.onAlertTriggered,
			}
			m.Start()
		})
	}
}

func TestAlertMonitor_Stop(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AlertMonitor{
				ctx:              tt.fields.ctx,
				priceAlertSvc:    tt.fields.priceAlertSvc,
				repo:             tt.fields.repo,
				stockService:     tt.fields.stockService,
				klineService:     tt.fields.klineService,
				ticker:           tt.fields.ticker,
				mu:               tt.fields.mu,
				running:          tt.fields.running,
				checkInterval:    tt.fields.checkInterval,
				onAlertTriggered: tt.fields.onAlertTriggered,
			}
			m.Stop()
		})
	}
}

func TestAlertMonitor_checkAllAlerts(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AlertMonitor{
				ctx:              tt.fields.ctx,
				priceAlertSvc:    tt.fields.priceAlertSvc,
				repo:             tt.fields.repo,
				stockService:     tt.fields.stockService,
				klineService:     tt.fields.klineService,
				ticker:           tt.fields.ticker,
				mu:               tt.fields.mu,
				running:          tt.fields.running,
				checkInterval:    tt.fields.checkInterval,
				onAlertTriggered: tt.fields.onAlertTriggered,
			}
			m.checkAllAlerts()
		})
	}
}

func TestAlertMonitor_convertToStockDataForAlert(t *testing.T) {
	type fields struct {
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
	type args struct {
		stockData *models.StockData
		klineData []*models.KLineData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *StockDataForAlert
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AlertMonitor{
				ctx:              tt.fields.ctx,
				priceAlertSvc:    tt.fields.priceAlertSvc,
				repo:             tt.fields.repo,
				stockService:     tt.fields.stockService,
				klineService:     tt.fields.klineService,
				ticker:           tt.fields.ticker,
				mu:               tt.fields.mu,
				running:          tt.fields.running,
				checkInterval:    tt.fields.checkInterval,
				onAlertTriggered: tt.fields.onAlertTriggered,
			}
			got, err := m.convertToStockDataForAlert(tt.args.stockData, tt.args.klineData)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToStockDataForAlert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToStockDataForAlert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlertMonitor_handleTriggeredAlert(t *testing.T) {
	type fields struct {
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
	type args struct {
		alert     *repositories.PriceThresholdAlert
		stockData *StockDataForAlert
		message   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AlertMonitor{
				ctx:              tt.fields.ctx,
				priceAlertSvc:    tt.fields.priceAlertSvc,
				repo:             tt.fields.repo,
				stockService:     tt.fields.stockService,
				klineService:     tt.fields.klineService,
				ticker:           tt.fields.ticker,
				mu:               tt.fields.mu,
				running:          tt.fields.running,
				checkInterval:    tt.fields.checkInterval,
				onAlertTriggered: tt.fields.onAlertTriggered,
			}
			m.handleTriggeredAlert(tt.args.alert, tt.args.stockData, tt.args.message)
		})
	}
}

func TestAlertMonitor_run(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AlertMonitor{
				ctx:              tt.fields.ctx,
				priceAlertSvc:    tt.fields.priceAlertSvc,
				repo:             tt.fields.repo,
				stockService:     tt.fields.stockService,
				klineService:     tt.fields.klineService,
				ticker:           tt.fields.ticker,
				mu:               tt.fields.mu,
				running:          tt.fields.running,
				checkInterval:    tt.fields.checkInterval,
				onAlertTriggered: tt.fields.onAlertTriggered,
			}
			m.run()
		})
	}
}

func TestAlertMonitor_sendNotification(t *testing.T) {
	type fields struct {
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
	type args struct {
		alert     *repositories.PriceThresholdAlert
		stockData *StockDataForAlert
		message   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AlertMonitor{
				ctx:              tt.fields.ctx,
				priceAlertSvc:    tt.fields.priceAlertSvc,
				repo:             tt.fields.repo,
				stockService:     tt.fields.stockService,
				klineService:     tt.fields.klineService,
				ticker:           tt.fields.ticker,
				mu:               tt.fields.mu,
				running:          tt.fields.running,
				checkInterval:    tt.fields.checkInterval,
				onAlertTriggered: tt.fields.onAlertTriggered,
			}
			m.sendNotification(tt.args.alert, tt.args.stockData, tt.args.message)
		})
	}
}

func TestAlertService_GetAlertConfig(t *testing.T) {
	type fields struct {
		repo repositories.AlertRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    models.AlertConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AlertService{
				repo: tt.fields.repo,
			}
			got, err := s.GetAlertConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAlertConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlertConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlertService_GetAlertHistory(t *testing.T) {
	type fields struct {
		repo repositories.AlertRepository
	}
	type args struct {
		stockCode string
		limit     int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AlertService{
				repo: tt.fields.repo,
			}
			got, err := s.GetAlertHistory(tt.args.stockCode, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAlertHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlertHistory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlertService_GetAlertHistoryForWails(t *testing.T) {
	type fields struct {
		repo repositories.AlertRepository
	}
	type args struct {
		stockCode string
		limit     int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AlertService{
				repo: tt.fields.repo,
			}
			got, err := s.GetAlertHistoryForWails(tt.args.stockCode, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAlertHistoryForWails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlertHistoryForWails() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlertService_GetAlertsForWails(t *testing.T) {
	type fields struct {
		repo repositories.AlertRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*models.PriceAlert
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AlertService{
				repo: tt.fields.repo,
			}
			got, err := s.GetAlertsForWails()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAlertsForWails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlertsForWails() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlertService_LoadActiveAlerts(t *testing.T) {
	type fields struct {
		repo repositories.AlertRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*models.PriceAlert
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AlertService{
				repo: tt.fields.repo,
			}
			got, err := s.LoadActiveAlerts()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadActiveAlerts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadActiveAlerts() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlertService_SaveActiveAlerts(t *testing.T) {
	type fields struct {
		repo repositories.AlertRepository
	}
	type args struct {
		alerts []*models.PriceAlert
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AlertService{
				repo: tt.fields.repo,
			}
			if err := s.SaveActiveAlerts(tt.args.alerts); (err != nil) != tt.wantErr {
				t.Errorf("SaveActiveAlerts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlertService_SaveAlert(t *testing.T) {
	type fields struct {
		repo repositories.AlertRepository
	}
	type args struct {
		alert   *models.PriceAlert
		message string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AlertService{
				repo: tt.fields.repo,
			}
			if err := s.SaveAlert(tt.args.alert, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("SaveAlert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlertService_SaveAlertsForWails(t *testing.T) {
	type fields struct {
		repo repositories.AlertRepository
	}
	type args struct {
		alerts []*models.PriceAlert
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AlertService{
				repo: tt.fields.repo,
			}
			if err := s.SaveAlertsForWails(tt.args.alerts); (err != nil) != tt.wantErr {
				t.Errorf("SaveAlertsForWails() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlertService_SetAlertsFromAI(t *testing.T) {
	type fields struct {
		repo repositories.AlertRepository
	}
	type args struct {
		code     string
		name     string
		drawings []models.TechnicalDrawing
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AlertService{
				repo: tt.fields.repo,
			}
			s.SetAlertsFromAI(tt.args.code, tt.args.name, tt.args.drawings)
		})
	}
}

func TestAlertService_UpdateAlertConfig(t *testing.T) {
	type fields struct {
		repo repositories.AlertRepository
	}
	type args struct {
		config models.AlertConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AlertService{
				repo: tt.fields.repo,
			}
			if err := s.UpdateAlertConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAlertConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlignStockData2MoneyFlow(t *testing.T) {
	type args struct {
		stockCode string
		data      []AlignedStockData
	}
	tests := []struct {
		name string
		args args
		want []models.MoneyFlowData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AlignStockData2MoneyFlow(tt.args.stockCode, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AlignStockData2MoneyFlow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnalysisCacheService_Get(t *testing.T) {
	type fields struct {
		cachePath string
		cache     models.AnalysisCache
		mutex     sync.RWMutex
	}
	type args struct {
		code   string
		role   string
		period string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *models.TechnicalAnalysisResult
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AnalysisCacheService{
				cachePath: tt.fields.cachePath,
				cache:     tt.fields.cache,
				mutex:     tt.fields.mutex,
			}
			got, got1 := s.Get(tt.args.code, tt.args.role, tt.args.period)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAnalysisCacheService_Set(t *testing.T) {
	type fields struct {
		cachePath string
		cache     models.AnalysisCache
		mutex     sync.RWMutex
	}
	type args struct {
		code   string
		role   string
		period string
		result models.TechnicalAnalysisResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AnalysisCacheService{
				cachePath: tt.fields.cachePath,
				cache:     tt.fields.cache,
				mutex:     tt.fields.mutex,
			}
			if err := s.Set(tt.args.code, tt.args.role, tt.args.period, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnalysisCacheService_load(t *testing.T) {
	type fields struct {
		cachePath string
		cache     models.AnalysisCache
		mutex     sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AnalysisCacheService{
				cachePath: tt.fields.cachePath,
				cache:     tt.fields.cache,
				mutex:     tt.fields.mutex,
			}
			s.load()
		})
	}
}

func TestAnalysisCacheService_save(t *testing.T) {
	type fields struct {
		cachePath string
		cache     models.AnalysisCache
		mutex     sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AnalysisCacheService{
				cachePath: tt.fields.cachePath,
				cache:     tt.fields.cache,
				mutex:     tt.fields.mutex,
			}
			if err := s.save(); (err != nil) != tt.wantErr {
				t.Errorf("save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnalyzeL2Market(t *testing.T) {
	type args struct {
		ticks []TickData
	}
	tests := []struct {
		name string
		args args
		want *OrderFlowStats
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnalyzeL2Market(tt.args.ticks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnalyzeL2Market() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacktestService_AnalyzePastSignals(t *testing.T) {
	type fields struct {
		stockService    *StockService
		strategyService *StrategyService
	}
	type args struct {
		days int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.SignalAnalysisResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BacktestService{
				stockService:    tt.fields.stockService,
				strategyService: tt.fields.strategyService,
			}
			got, err := s.AnalyzePastSignals(tt.args.days)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzePastSignals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnalyzePastSignals() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacktestService_BacktestDecisionPioneer(t *testing.T) {
	type fields struct {
		stockService    *StockService
		strategyService *StrategyService
	}
	type args struct {
		code           string
		initialCapital float64
		startDate      string
		endDate        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.BacktestResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BacktestService{
				stockService:    tt.fields.stockService,
				strategyService: tt.fields.strategyService,
			}
			got, err := s.BacktestDecisionPioneer(tt.args.code, tt.args.initialCapital, tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("BacktestDecisionPioneer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BacktestDecisionPioneer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacktestService_BacktestMACD(t *testing.T) {
	type fields struct {
		stockService    *StockService
		strategyService *StrategyService
	}
	type args struct {
		code           string
		fastPeriod     int
		slowPeriod     int
		signalPeriod   int
		initialCapital float64
		startDate      string
		endDate        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.BacktestResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BacktestService{
				stockService:    tt.fields.stockService,
				strategyService: tt.fields.strategyService,
			}
			got, err := s.BacktestMACD(tt.args.code, tt.args.fastPeriod, tt.args.slowPeriod, tt.args.signalPeriod, tt.args.initialCapital, tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("BacktestMACD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BacktestMACD() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacktestService_BacktestRSI(t *testing.T) {
	type fields struct {
		stockService    *StockService
		strategyService *StrategyService
	}
	type args struct {
		code           string
		period         int
		buyThreshold   float64
		sellThreshold  float64
		initialCapital float64
		startDate      string
		endDate        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.BacktestResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BacktestService{
				stockService:    tt.fields.stockService,
				strategyService: tt.fields.strategyService,
			}
			got, err := s.BacktestRSI(tt.args.code, tt.args.period, tt.args.buyThreshold, tt.args.sellThreshold, tt.args.initialCapital, tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("BacktestRSI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BacktestRSI() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacktestService_BacktestSimpleMA(t *testing.T) {
	type fields struct {
		stockService    *StockService
		strategyService *StrategyService
	}
	type args struct {
		code           string
		shortPeriod    int
		longPeriod     int
		initialCapital float64
		startDate      string
		endDate        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.BacktestResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BacktestService{
				stockService:    tt.fields.stockService,
				strategyService: tt.fields.strategyService,
			}
			got, err := s.BacktestSimpleMA(tt.args.code, tt.args.shortPeriod, tt.args.longPeriod, tt.args.initialCapital, tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("BacktestSimpleMA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BacktestSimpleMA() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacktestService_runBacktest(t *testing.T) {
	type fields struct {
		stockService    *StockService
		strategyService *StrategyService
	}
	type args struct {
		code           string
		strategyName   string
		initialCapital float64
		startDate      string
		endDate        string
		limit          int
		signalGen      SignalGenerator
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.BacktestResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BacktestService{
				stockService:    tt.fields.stockService,
				strategyService: tt.fields.strategyService,
			}
			got, err := s.runBacktest(tt.args.code, tt.args.strategyName, tt.args.initialCapital, tt.args.startDate, tt.args.endDate, tt.args.limit, tt.args.signalGen)
			if (err != nil) != tt.wantErr {
				t.Errorf("runBacktest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runBacktest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigService_GetGlobalStrategyConfig(t *testing.T) {
	type fields struct {
		repo repositories.ConfigRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    GlobalStrategyConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ConfigService{
				repo: tt.fields.repo,
			}
			got, err := s.GetGlobalStrategyConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGlobalStrategyConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGlobalStrategyConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigService_LoadAIConfig(t *testing.T) {
	type fields struct {
		repo repositories.ConfigRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    AIResolvedConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ConfigService{
				repo: tt.fields.repo,
			}
			got, err := s.LoadAIConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAIConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadAIConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigService_MigrateAIConfigFromYAML(t *testing.T) {
	type fields struct {
		repo repositories.ConfigRepository
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ConfigService{
				repo: tt.fields.repo,
			}
			if err := s.MigrateAIConfigFromYAML(); (err != nil) != tt.wantErr {
				t.Errorf("MigrateAIConfigFromYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigService_SaveAIConfig(t *testing.T) {
	type fields struct {
		repo repositories.ConfigRepository
	}
	type args struct {
		config AIResolvedConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ConfigService{
				repo: tt.fields.repo,
			}
			if err := s.SaveAIConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("SaveAIConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigService_UpdateGlobalStrategyConfig(t *testing.T) {
	type fields struct {
		repo repositories.ConfigRepository
	}
	type args struct {
		config GlobalStrategyConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ConfigService{
				repo: tt.fields.repo,
			}
			if err := s.UpdateGlobalStrategyConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("UpdateGlobalStrategyConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigService_getConfigValue(t *testing.T) {
	type fields struct {
		repo repositories.ConfigRepository
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ConfigService{
				repo: tt.fields.repo,
			}
			got, err := s.getConfigValue(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getConfigValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getConfigValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigService_setConfigValue(t *testing.T) {
	type fields struct {
		repo repositories.ConfigRepository
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ConfigService{
				repo: tt.fields.repo,
			}
			if err := s.setConfigValue(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("setConfigValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCookieManager_GetAutoCookie(t *testing.T) {
	cookie, err := NewCookieManager().GetAutoCookie()
	if err != nil {
		t.Logf("GetAutoCookie error = %v", err)
		return
	}
	t.Log(cookie)
}

func TestCookieManager_GetStockCookie(t *testing.T) {
	cookie, err := NewCookieManager().GetStockCookie("https://quote.eastmoney.com/center/gridlist.html")
	if err != nil {
		t.Logf("GetStockCookie error = %v", err)
		return
	}
	t.Log(cookie)
}

func TestDBService_ClearKLineCacheTable(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			if err := s.ClearKLineCacheTable(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("ClearKLineCacheTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBService_Close(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			s.Close()
		})
	}
}

func TestDBService_CreateKLineCacheTable(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			if err := s.CreateKLineCacheTable(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("CreateKLineCacheTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBService_GetAllSyncedStocks(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			got, err := s.GetAllSyncedStocks()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllSyncedStocks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllSyncedStocks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBService_GetDB(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	tests := []struct {
		name   string
		fields fields
		want   *gorm.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			if got := s.GetDB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBService_GetDBPath(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			if got := s.GetDBPath(); got != tt.want {
				t.Errorf("GetDBPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBService_GetKLineCountByCode(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			got, err := s.GetKLineCountByCode(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKLineCountByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetKLineCountByCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBService_GetKLineDataFromCache(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	type args struct {
		code  string
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			got, err := s.GetKLineDataFromCache(tt.args.code, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKLineDataFromCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKLineDataFromCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBService_GetKLineDataWithPagination(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	type args struct {
		code      string
		startDate string
		endDate   string
		page      int
		pageSize  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			got, got1, err := s.GetKLineDataWithPagination(tt.args.code, tt.args.startDate, tt.args.endDate, tt.args.page, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKLineDataWithPagination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKLineDataWithPagination() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetKLineDataWithPagination() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDBService_GetLatestKLineDate(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			got, err := s.GetLatestKLineDate(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestKLineDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLatestKLineDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBService_InsertOrUpdateKLineData(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	type args struct {
		code   string
		klines []map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			got, got1, err := s.InsertOrUpdateKLineData(tt.args.code, tt.args.klines)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertOrUpdateKLineData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InsertOrUpdateKLineData() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("InsertOrUpdateKLineData() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDBService_fixWatchlistNullValues(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			if err := s.fixWatchlistNullValues(); (err != nil) != tt.wantErr {
				t.Errorf("fixWatchlistNullValues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBService_initTables(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			if err := s.initTables(); (err != nil) != tt.wantErr {
				t.Errorf("initTables() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBService_insertDefaultAlertTemplates(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			if err := s.insertDefaultAlertTemplates(); (err != nil) != tt.wantErr {
				t.Errorf("insertDefaultAlertTemplates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBService_insertDefaultConfigs(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			if err := s.insertDefaultConfigs(); (err != nil) != tt.wantErr {
				t.Errorf("insertDefaultConfigs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBService_manualMigrateWatchlistTable(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			if err := s.manualMigrateWatchlistTable(); (err != nil) != tt.wantErr {
				t.Errorf("manualMigrateWatchlistTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBService_prepareWatchlistTable(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			if err := s.prepareWatchlistTable(); (err != nil) != tt.wantErr {
				t.Errorf("prepareWatchlistTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBService_rebuildWatchlistTable(t *testing.T) {
	type fields struct {
		db     *gorm.DB
		dbPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DBService{
				db:     tt.fields.db,
				dbPath: tt.fields.dbPath,
			}
			if err := s.rebuildWatchlistTable(); (err != nil) != tt.wantErr {
				t.Errorf("rebuildWatchlistTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetSortedData(t *testing.T) {
	type args struct {
		dataMap map[string]*AlignedStockData
	}
	tests := []struct {
		name string
		args args
		want []AlignedStockData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSortedData(tt.args.dataMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSortedData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLineSyncService_GetKLineSyncHistory(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
		ctx       context.Context
		running   bool
		mu        sync.Mutex
	}
	type args struct {
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KLineSyncService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
				ctx:       tt.fields.ctx,
				running:   tt.fields.running,
				mu:        tt.fields.mu,
			}
			got, err := s.GetKLineSyncHistory(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKLineSyncHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKLineSyncHistory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLineSyncService_GetSyncProgress(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
		ctx       context.Context
		running   bool
		mu        sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    *KLineSyncProgress
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KLineSyncService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
				ctx:       tt.fields.ctx,
				running:   tt.fields.running,
				mu:        tt.fields.mu,
			}
			got, err := s.GetSyncProgress()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSyncProgress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSyncProgress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLineSyncService_SetContext(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
		ctx       context.Context
		running   bool
		mu        sync.Mutex
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KLineSyncService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
				ctx:       tt.fields.ctx,
				running:   tt.fields.running,
				mu:        tt.fields.mu,
			}
			s.SetContext(tt.args.ctx)
		})
	}
}

func TestKLineSyncService_StartKLineSync(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
		ctx       context.Context
		running   bool
		mu        sync.Mutex
	}
	type args struct {
		days int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *KLineSyncResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KLineSyncService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
				ctx:       tt.fields.ctx,
				running:   tt.fields.running,
				mu:        tt.fields.mu,
			}
			got, err := s.StartKLineSync(tt.args.days)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartKLineSync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartKLineSync() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLineSyncService_emitProgress(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
		ctx       context.Context
		running   bool
		mu        sync.Mutex
	}
	type args struct {
		progress *KLineSyncProgress
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KLineSyncService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
				ctx:       tt.fields.ctx,
				running:   tt.fields.running,
				mu:        tt.fields.mu,
			}
			s.emitProgress(tt.args.progress)
		})
	}
}

func TestKLineSyncService_fetchKLineData(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
		ctx       context.Context
		running   bool
		mu        sync.Mutex
	}
	type args struct {
		task *KLineSyncTask
		days int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KLineSyncService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
				ctx:       tt.fields.ctx,
				running:   tt.fields.running,
				mu:        tt.fields.mu,
			}
			got, err := s.fetchKLineData(tt.args.task, tt.args.days)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchKLineData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchKLineData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLineSyncService_getActiveStocks(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
		ctx       context.Context
		running   bool
		mu        sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*KLineSyncTask
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KLineSyncService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
				ctx:       tt.fields.ctx,
				running:   tt.fields.running,
				mu:        tt.fields.mu,
			}
			got, err := s.getActiveStocks()
			if (err != nil) != tt.wantErr {
				t.Errorf("getActiveStocks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getActiveStocks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLineSyncService_recordSyncHistory(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
		ctx       context.Context
		running   bool
		mu        sync.Mutex
	}
	type args struct {
		task         *KLineSyncTask
		days         int
		totalRecords int
		added        int
		updated      int
		success      bool
		errorMsg     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KLineSyncService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
				ctx:       tt.fields.ctx,
				running:   tt.fields.running,
				mu:        tt.fields.mu,
			}
			if err := s.recordSyncHistory(tt.args.task, tt.args.days, tt.args.totalRecords, tt.args.added, tt.args.updated, tt.args.success, tt.args.errorMsg); (err != nil) != tt.wantErr {
				t.Errorf("recordSyncHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKLineSyncService_saveKLineData(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
		ctx       context.Context
		running   bool
		mu        sync.Mutex
	}
	type args struct {
		code   string
		klines []map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KLineSyncService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
				ctx:       tt.fields.ctx,
				running:   tt.fields.running,
				mu:        tt.fields.mu,
			}
			got, got1, err := s.saveKLineData(tt.args.code, tt.args.klines)
			if (err != nil) != tt.wantErr {
				t.Errorf("saveKLineData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("saveKLineData() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("saveKLineData() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestKLineSyncService_updateProgress(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
		ctx       context.Context
		running   bool
		mu        sync.Mutex
	}
	type args struct {
		progress     *KLineSyncProgress
		currentIndex int
		totalCount   int
		code         string
		name         string
		successCount int
		failedCount  int
		totalRecords int
		startTime    time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KLineSyncService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
				ctx:       tt.fields.ctx,
				running:   tt.fields.running,
				mu:        tt.fields.mu,
			}
			s.updateProgress(tt.args.progress, tt.args.currentIndex, tt.args.totalCount, tt.args.code, tt.args.name, tt.args.successCount, tt.args.failedCount, tt.args.totalRecords, tt.args.startTime)
		})
	}
}

func TestLoadAIConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    AIResolvedConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadAIConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAIConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadAIConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyFlowService_FetchAndSaveHistory(t *testing.T) {
	type fields struct {
		repo   *repositories.MoneyFlowRepository
		client *resty.Client
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MoneyFlowService{
				repo:   tt.fields.repo,
				client: tt.fields.client,
			}
			if err := s.FetchAndSaveHistory(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("FetchAndSaveHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMoneyFlowService_generateSecid(t *testing.T) {
	type fields struct {
		repo   *repositories.MoneyFlowRepository
		client *resty.Client
	}
	type args struct {
		code string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MoneyFlowService{
				repo:   tt.fields.repo,
				client: tt.fields.client,
			}
			if got := s.generateSecid(tt.args.code); got != tt.want {
				t.Errorf("generateSecid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyFlowService_parseFloat(t *testing.T) {
	type fields struct {
		repo   *repositories.MoneyFlowRepository
		client *resty.Client
	}
	type args struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MoneyFlowService{
				repo:   tt.fields.repo,
				client: tt.fields.client,
			}
			if got := s.parseFloat(tt.args.val); got != tt.want {
				t.Errorf("parseFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAIService(t *testing.T) {
	type args struct {
		cfg AIResolvedConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *AIService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAIService(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAIService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAIService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAlertMonitor(t *testing.T) {
	type args struct {
		ctx           context.Context
		priceAlertSvc *PriceAlertService
		stockService  StockDataService
		klineService  KLineDataService
	}
	tests := []struct {
		name string
		args args
		want *AlertMonitor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAlertMonitor(tt.args.ctx, tt.args.priceAlertSvc, tt.args.stockService, tt.args.klineService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAlertMonitor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAlertService(t *testing.T) {
	type args struct {
		repo repositories.AlertRepository
	}
	tests := []struct {
		name string
		args args
		want *AlertService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAlertService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAlertService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAlertStorage(t *testing.T) {
	type args struct {
		dbSvc *DBService
	}
	tests := []struct {
		name string
		args args
		want *AlertService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAlertStorage(tt.args.dbSvc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAlertStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAnalysisCacheService(t *testing.T) {
	tests := []struct {
		name    string
		want    *AnalysisCacheService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAnalysisCacheService()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAnalysisCacheService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnalysisCacheService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBacktestService(t *testing.T) {
	type args struct {
		stockService    *StockService
		strategyService *StrategyService
	}
	tests := []struct {
		name string
		args args
		want *BacktestService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBacktestService(tt.args.stockService, tt.args.strategyService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBacktestService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConfigService(t *testing.T) {
	type args struct {
		repo repositories.ConfigRepository
	}
	tests := []struct {
		name string
		args args
		want *ConfigService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfigService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCookieManager(t *testing.T) {
	tests := []struct {
		name string
		want *CookieManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCookieManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCookieManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDBService(t *testing.T) {
	tests := []struct {
		name    string
		want    *DBService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDBService()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDBService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDBService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKLineSyncService(t *testing.T) {
	type args struct {
		dbService *DBService
	}
	tests := []struct {
		name string
		args args
		want *KLineSyncService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKLineSyncService(tt.args.dbService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKLineSyncService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMoneyFlowService(t *testing.T) {
	type args struct {
		repo *repositories.MoneyFlowRepository
	}
	tests := []struct {
		name string
		args args
		want *MoneyFlowService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMoneyFlowService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMoneyFlowService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPositionService(t *testing.T) {
	type args struct {
		repo repositories.PositionRepository
	}
	tests := []struct {
		name string
		args args
		want *PositionService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPositionService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPositionService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPositionStorageService(t *testing.T) {
	type args struct {
		dbSvc *DBService
	}
	tests := []struct {
		name string
		args args
		want *PositionService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPositionStorageService(tt.args.dbSvc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPositionStorageService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPriceAlertService(t *testing.T) {
	type args struct {
		repo *repositories.PriceAlertRepository
	}
	tests := []struct {
		name string
		args args
		want *PriceAlertService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPriceAlertService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPriceAlertService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStockMarketService(t *testing.T) {
	type args struct {
		dbService *DBService
	}
	tests := []struct {
		name string
		args args
		want *StockMarketService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStockMarketService(tt.args.dbService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStockMarketService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStockService(t *testing.T) {
	tests := []struct {
		name string
		want *StockService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStockService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStockService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStrategyService(t *testing.T) {
	type args struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	tests := []struct {
		name string
		args args
		want *StrategyService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStrategyService(tt.args.repo, tt.args.moneyFlowRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStrategyService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSyncService(t *testing.T) {
	type args struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
	}
	tests := []struct {
		name string
		args args
		want *SyncService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSyncService(tt.args.dbService, tt.args.stockMarketService, tt.args.moneyFlowRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSyncService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWatchlistService(t *testing.T) {
	type args struct {
		repo repositories.WatchlistRepository
	}
	tests := []struct {
		name string
		args args
		want *WatchlistService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWatchlistService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWatchlistService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAndMerge(t *testing.T) {
	type args struct {
		klineData []string
		fflowData []string
	}
	tests := []struct {
		name string
		args args
		want map[string]*AlignedStockData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseAndMerge(tt.args.klineData, tt.args.fflowData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAndMerge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPositionService_GetPositions(t *testing.T) {
	type fields struct {
		repo repositories.PositionRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]*models.Position
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PositionService{
				repo: tt.fields.repo,
			}
			got, err := s.GetPositions()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPositions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPositions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPositionService_RemovePosition(t *testing.T) {
	type fields struct {
		repo repositories.PositionRepository
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PositionService{
				repo: tt.fields.repo,
			}
			if err := s.RemovePosition(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("RemovePosition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPositionService_SavePosition(t *testing.T) {
	type fields struct {
		repo repositories.PositionRepository
	}
	type args struct {
		pos *models.Position
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PositionService{
				repo: tt.fields.repo,
			}
			if err := s.SavePosition(tt.args.pos); (err != nil) != tt.wantErr {
				t.Errorf("SavePosition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPriceAlertService_CheckAlert(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		alert     *repositories.PriceThresholdAlert
		stockData *StockDataForAlert
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			got, got1, err := s.CheckAlert(tt.args.alert, tt.args.stockData)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAlert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAlert() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CheckAlert() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPriceAlertService_CreateAlert(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		req *CreateAlertRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			if err := s.CreateAlert(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("CreateAlert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPriceAlertService_CreateAlertFromTemplate(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		templateID string
		stockCode  string
		stockName  string
		params     map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			if err := s.CreateAlertFromTemplate(tt.args.templateID, tt.args.stockCode, tt.args.stockName, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("CreateAlertFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPriceAlertService_DeleteAlert(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			if err := s.DeleteAlert(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAlert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPriceAlertService_GetActiveAlerts(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*repositories.PriceThresholdAlert
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			got, err := s.GetActiveAlerts()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActiveAlerts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActiveAlerts() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceAlertService_GetAlertsByStockCode(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		stockCode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*repositories.PriceThresholdAlert
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			got, err := s.GetAlertsByStockCode(tt.args.stockCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAlertsByStockCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlertsByStockCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceAlertService_GetAllAlerts(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*repositories.PriceThresholdAlert
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			got, err := s.GetAllAlerts()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllAlerts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllAlerts() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceAlertService_GetAllTemplates(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*repositories.PriceAlertTemplate
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			got, err := s.GetAllTemplates()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTemplates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllTemplates() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceAlertService_GetRepository(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	tests := []struct {
		name   string
		fields fields
		want   *repositories.PriceAlertRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			if got := s.GetRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceAlertService_GetTriggerHistory(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		stockCode string
		limit     int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*repositories.PriceAlertTriggerHistory
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			got, err := s.GetTriggerHistory(tt.args.stockCode, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTriggerHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTriggerHistory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceAlertService_ToggleAlertStatus(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		id       int64
		isActive bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			if err := s.ToggleAlertStatus(tt.args.id, tt.args.isActive); (err != nil) != tt.wantErr {
				t.Errorf("ToggleAlertStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPriceAlertService_TriggerAlert(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		alert     *repositories.PriceThresholdAlert
		stockData *StockDataForAlert
		message   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			if err := s.TriggerAlert(tt.args.alert, tt.args.stockData, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("TriggerAlert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPriceAlertService_UpdateAlert(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		req *UpdateAlertRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			if err := s.UpdateAlert(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAlert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPriceAlertService_buildTriggerMessage(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		conditions *repositories.PriceAlertConditions
		messages   []string
		index      int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			if got := s.buildTriggerMessage(tt.args.conditions, tt.args.messages, tt.args.index); got != tt.want {
				t.Errorf("buildTriggerMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceAlertService_evaluateCondition(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		condition   *repositories.PriceAlertCondition
		stockData   *StockDataForAlert
		sensitivity float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			got, got1 := s.evaluateCondition(tt.args.condition, tt.args.stockData, tt.args.sensitivity)
			if got != tt.want {
				t.Errorf("evaluateCondition() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("evaluateCondition() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPriceAlertService_evaluateConditions(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		conditions  *repositories.PriceAlertConditions
		stockData   *StockDataForAlert
		sensitivity float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			got, got1 := s.evaluateConditions(tt.args.conditions, tt.args.stockData, tt.args.sensitivity)
			if got != tt.want {
				t.Errorf("evaluateConditions() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("evaluateConditions() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPriceAlertService_isValidConditionsJSON(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		jsonStr string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			if got := s.isValidConditionsJSON(tt.args.jsonStr); got != tt.want {
				t.Errorf("isValidConditionsJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceAlertService_replaceTemplateParams(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		conditions *repositories.PriceAlertConditions
		params     map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			s.replaceTemplateParams(tt.args.conditions, tt.args.params)
		})
	}
}

func TestPriceAlertService_validateAlertRequest(t *testing.T) {
	type fields struct {
		repo *repositories.PriceAlertRepository
	}
	type args struct {
		req *CreateAlertRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceAlertService{
				repo: tt.fields.repo,
			}
			if err := s.validateAlertRequest(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("validateAlertRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStockMarketService_GetAllStockCodes(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockMarketService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
			}
			got, err := s.GetAllStockCodes()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllStockCodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllStockCodes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockMarketService_GetIndustries(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    []IndustryInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockMarketService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
			}
			got, err := s.GetIndustries()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIndustries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIndustries() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockMarketService_GetStocksList(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
	}
	type args struct {
		page     int
		pageSize int
		search   string
		industry string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []StockMarketData
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockMarketService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
			}
			got, got1, err := s.GetStocksList(tt.args.page, tt.args.pageSize, tt.args.search, tt.args.industry)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStocksList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStocksList() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetStocksList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestStockMarketService_GetSyncStats(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockMarketService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
			}
			got, err := s.GetSyncStats()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSyncStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSyncStats() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockMarketService_SyncAllStocks(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    *SyncStocksResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockMarketService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
			}
			got, err := s.SyncAllStocks()
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncAllStocks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SyncAllStocks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockMarketService_parseStockItemToEntity(t *testing.T) {
	type fields struct {
		dbService *DBService
		client    *resty.Client
	}
	type args struct {
		item      interface{}
		updatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *models.StockEntity
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockMarketService{
				dbService: tt.fields.dbService,
				client:    tt.fields.client,
			}
			if got := s.parseStockItemToEntity(tt.args.item, tt.args.updatedAt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseStockItemToEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_BatchAnalyzeStocks(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		ctx   context.Context
		codes []string
		role  string
		aiSvc *AIService
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			if err := s.BatchAnalyzeStocks(tt.args.ctx, tt.args.codes, tt.args.role, tt.args.aiSvc); (err != nil) != tt.wantErr {
				t.Errorf("BatchAnalyzeStocks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStockService_BatchSyncStockData(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		codes     []string
		startDate string
		endDate   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			if err := s.BatchSyncStockData(tt.args.codes, tt.args.startDate, tt.args.endDate); (err != nil) != tt.wantErr {
				t.Errorf("BatchSyncStockData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStockService_ClearStockCache(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			if err := s.ClearStockCache(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("ClearStockCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStockService_GetDataSyncStats(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *models.DataSyncStats
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.GetDataSyncStats()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataSyncStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDataSyncStats() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_GetIntradayData(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.IntradayResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.GetIntradayData(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIntradayData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIntradayData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_GetKLineData(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code   string
		limit  int
		period string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*models.KLineData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.GetKLineData(tt.args.code, tt.args.limit, tt.args.period)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKLineData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKLineData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_GetKLineFromCache(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code  string
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*models.KLineData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.GetKLineFromCache(tt.args.code, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKLineFromCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKLineFromCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_GetMoneyFlowData(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.MoneyFlowResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.GetMoneyFlowData(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMoneyFlowData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMoneyFlowData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_GetStockByCode(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.StockData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.GetStockByCode(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStockByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStockByCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_GetStockDetail(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.StockDetail
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.GetStockDetail(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStockDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStockDetail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_GetStockHealthCheck(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.HealthCheckResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.GetStockHealthCheck(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStockHealthCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStockHealthCheck() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_SearchStock(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		keyword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*models.StockData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.SearchStock(tt.args.keyword)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchStock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchStock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_SearchStockLegacy(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		keyword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*models.StockData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.SearchStockLegacy(tt.args.keyword)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchStockLegacy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchStockLegacy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_SetDBService(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		db *DBService
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			s.SetDBService(tt.args.db)
		})
	}
}

func TestStockService_Startup(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			s.Startup(tt.args.ctx)
		})
	}
}

func TestStockService_StopIntradayStream(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			s.StopIntradayStream(tt.args.code)
		})
	}
}

func TestStockService_StreamIntradayData(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			s.StreamIntradayData(tt.args.code)
		})
	}
}

func TestStockService_SyncStockData(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code      string
		startDate string
		endDate   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.SyncResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.SyncStockData(tt.args.code, tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncStockData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SyncStockData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_calculateIndicators(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		klines []*models.KLineData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			s.calculateIndicators(tt.args.klines)
		})
	}
}

func TestStockService_getFinancialSummary(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.FinancialSummary
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.getFinancialSummary(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFinancialSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFinancialSummary() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_getIndustryInfo(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.IndustryInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.getIndustryInfo(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIndustryInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIndustryInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_getOrderBook(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.OrderBook
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			got, err := s.getOrderBook(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOrderBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOrderBook() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_getSecID(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			if got := s.getSecID(tt.args.code); got != tt.want {
				t.Errorf("getSecID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_logWarnThrottled(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		code   string
		msg    string
		fields []zap.Field
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			s.logWarnThrottled(tt.args.code, tt.args.msg, tt.args.fields...)
		})
	}
}

func TestStockService_parsePrice(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		p string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			if got := s.parsePrice(tt.args.p); got != tt.want {
				t.Errorf("parsePrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockService_sleepBackoff(t *testing.T) {
	type fields struct {
		client       *resty.Client
		sseClient    *resty.Client
		ctx          context.Context
		cancel       context.CancelFunc
		streamMu     sync.Mutex
		streams      map[string]context.CancelFunc
		emitIntraday func(ctx context.Context, code string, trends []string)
		dbService    *DBService
		warnMu       sync.Mutex
		lastWarnAt   map[string]time.Time
		exactURL     string
		listURL      string
		klineURL     string
	}
	type args struct {
		ctx   context.Context
		retry int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				client:       tt.fields.client,
				sseClient:    tt.fields.sseClient,
				ctx:          tt.fields.ctx,
				cancel:       tt.fields.cancel,
				streamMu:     tt.fields.streamMu,
				streams:      tt.fields.streams,
				emitIntraday: tt.fields.emitIntraday,
				dbService:    tt.fields.dbService,
				warnMu:       tt.fields.warnMu,
				lastWarnAt:   tt.fields.lastWarnAt,
				exactURL:     tt.fields.exactURL,
				listURL:      tt.fields.listURL,
				klineURL:     tt.fields.klineURL,
			}
			if got := s.sleepBackoff(tt.args.ctx, tt.args.retry); got != tt.want {
				t.Errorf("sleepBackoff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_CalculateBuildSignals(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.StrategySignal
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			got, err := s.CalculateBuildSignals(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateBuildSignals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateBuildSignals() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_CheckDecisionPioneerSellSignal(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		data []models.MoneyFlowData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *models.StrategySignal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			if got := s.CheckDecisionPioneerSellSignal(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckDecisionPioneerSellSignal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_CheckDecisionPioneerSignal(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		data   []models.MoneyFlowData
		circMV float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *models.StrategySignal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			if got := s.CheckDecisionPioneerSignal(tt.args.data, tt.args.circMV); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckDecisionPioneerSignal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_CheckMoneySurgeSignal(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		data []models.MoneyFlowData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *models.StrategySignal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			if got := s.CheckMoneySurgeSignal(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckMoneySurgeSignal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_CreateStrategy(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		name         string
		description  string
		strategyType string
		parameters   map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.StrategyConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			got, err := s.CreateStrategy(tt.args.name, tt.args.description, tt.args.strategyType, tt.args.parameters)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStrategy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateStrategy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_DeleteStrategy(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			if err := s.DeleteStrategy(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteStrategy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStrategyService_GetAllMoneyFlowHistory(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.MoneyFlowData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			got, err := s.GetAllMoneyFlowHistory(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMoneyFlowHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllMoneyFlowHistory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_GetAllStrategies(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []models.StrategyConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			got, err := s.GetAllStrategies()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllStrategies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllStrategies() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_GetLatestSignals(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.StrategySignal
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			got, err := s.GetLatestSignals(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestSignals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLatestSignals() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_GetRecentMoneyFlows(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		code  string
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.MoneyFlowData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			got, err := s.GetRecentMoneyFlows(tt.args.code, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecentMoneyFlows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRecentMoneyFlows() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_GetSignalsByDateRange(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		startDate string
		endDate   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.StrategySignal
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			got, err := s.GetSignalsByDateRange(tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSignalsByDateRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignalsByDateRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_GetSignalsByStockCode(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.StrategySignal
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			got, err := s.GetSignalsByStockCode(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSignalsByStockCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignalsByStockCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_GetStockCircMV(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			got, err := s.GetStockCircMV(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStockCircMV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetStockCircMV() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_GetStrategy(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.StrategyConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			got, err := s.GetStrategy(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStrategy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStrategy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_GetStrategyTypes(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	tests := []struct {
		name   string
		fields fields
		want   []models.StrategyTypeDefinition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			if got := s.GetStrategyTypes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStrategyTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrategyService_UpdateSignalAIResult(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		code         string
		tradeDate    string
		strategyName string
		aiScore      int
		aiReason     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			if err := s.UpdateSignalAIResult(tt.args.code, tt.args.tradeDate, tt.args.strategyName, tt.args.aiScore, tt.args.aiReason); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSignalAIResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStrategyService_UpdateStrategy(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		id           int64
		name         string
		description  string
		strategyType string
		parameters   map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			if err := s.UpdateStrategy(tt.args.id, tt.args.name, tt.args.description, tt.args.strategyType, tt.args.parameters); (err != nil) != tt.wantErr {
				t.Errorf("UpdateStrategy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStrategyService_UpdateStrategyBacktestResult(t *testing.T) {
	type fields struct {
		repo          *repositories.StrategyRepository
		moneyFlowRepo *repositories.MoneyFlowRepository
	}
	type args struct {
		id             int64
		backtestResult map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StrategyService{
				repo:          tt.fields.repo,
				moneyFlowRepo: tt.fields.moneyFlowRepo,
			}
			if err := s.UpdateStrategyBacktestResult(tt.args.id, tt.args.backtestResult); (err != nil) != tt.wantErr {
				t.Errorf("UpdateStrategyBacktestResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSyncService_FetchAllDayTicks(t *testing.T) {
	type fields struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
		client             *resty.Client
		ctx                context.Context
		running            bool
		mu                 sync.Mutex
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *OrderFlowStats
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				dbService:          tt.fields.dbService,
				stockMarketService: tt.fields.stockMarketService,
				moneyFlowRepo:      tt.fields.moneyFlowRepo,
				client:             tt.fields.client,
				ctx:                tt.fields.ctx,
				running:            tt.fields.running,
				mu:                 tt.fields.mu,
			}
			got, err := s.FetchAllDayTicks(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchAllDayTicks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchAllDayTicks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncService_FetchAndSaveHistoryFlow(t *testing.T) {
	type fields struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
		client             *resty.Client
		ctx                context.Context
		running            bool
		mu                 sync.Mutex
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				dbService:          tt.fields.dbService,
				stockMarketService: tt.fields.stockMarketService,
				moneyFlowRepo:      tt.fields.moneyFlowRepo,
				client:             tt.fields.client,
				ctx:                tt.fields.ctx,
				running:            tt.fields.running,
				mu:                 tt.fields.mu,
			}
			if err := s.FetchAndSaveHistoryFlow(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("FetchAndSaveHistoryFlow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSyncService_FetchHistoryFlowData(t *testing.T) {
	type fields struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
		client             *resty.Client
		ctx                context.Context
		running            bool
		mu                 sync.Mutex
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.MoneyFlowData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				dbService:          tt.fields.dbService,
				stockMarketService: tt.fields.stockMarketService,
				moneyFlowRepo:      tt.fields.moneyFlowRepo,
				client:             tt.fields.client,
				ctx:                tt.fields.ctx,
				running:            tt.fields.running,
				mu:                 tt.fields.mu,
			}
			got, err := s.FetchHistoryFlowData(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchHistoryFlowData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchHistoryFlowData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncService_FetchHistoryFlowDataV2(t *testing.T) {
	type fields struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
		client             *resty.Client
		ctx                context.Context
		running            bool
		mu                 sync.Mutex
	}
	type args struct {
		code  string
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*AlignedStockData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				dbService:          tt.fields.dbService,
				stockMarketService: tt.fields.stockMarketService,
				moneyFlowRepo:      tt.fields.moneyFlowRepo,
				client:             tt.fields.client,
				ctx:                tt.fields.ctx,
				running:            tt.fields.running,
				mu:                 tt.fields.mu,
			}
			got, err := s.FetchHistoryFlowDataV2(tt.args.code, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchHistoryFlowDataV2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchHistoryFlowDataV2() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncService_ScanAndSaveStrategySignals(t *testing.T) {
	type fields struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
		client             *resty.Client
		ctx                context.Context
		running            bool
		mu                 sync.Mutex
	}
	type args struct {
		code  string
		flows []models.MoneyFlowData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				dbService:          tt.fields.dbService,
				stockMarketService: tt.fields.stockMarketService,
				moneyFlowRepo:      tt.fields.moneyFlowRepo,
				client:             tt.fields.client,
				ctx:                tt.fields.ctx,
				running:            tt.fields.running,
				mu:                 tt.fields.mu,
			}
			s.ScanAndSaveStrategySignals(tt.args.code, tt.args.flows)
		})
	}
}

func TestSyncService_SetContext(t *testing.T) {
	type fields struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
		client             *resty.Client
		ctx                context.Context
		running            bool
		mu                 sync.Mutex
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				dbService:          tt.fields.dbService,
				stockMarketService: tt.fields.stockMarketService,
				moneyFlowRepo:      tt.fields.moneyFlowRepo,
				client:             tt.fields.client,
				ctx:                tt.fields.ctx,
				running:            tt.fields.running,
				mu:                 tt.fields.mu,
			}
			s.SetContext(tt.args.ctx)
		})
	}
}

func TestSyncService_StartFullMarketSync(t *testing.T) {
	type fields struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
		client             *resty.Client
		ctx                context.Context
		running            bool
		mu                 sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				dbService:          tt.fields.dbService,
				stockMarketService: tt.fields.stockMarketService,
				moneyFlowRepo:      tt.fields.moneyFlowRepo,
				client:             tt.fields.client,
				ctx:                tt.fields.ctx,
				running:            tt.fields.running,
				mu:                 tt.fields.mu,
			}
			if err := s.StartFullMarketSync(); (err != nil) != tt.wantErr {
				t.Errorf("StartFullMarketSync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSyncService_SyncAndScanSingleStock(t *testing.T) {
	type fields struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
		client             *resty.Client
		ctx                context.Context
		running            bool
		mu                 sync.Mutex
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.StrategySignal
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				dbService:          tt.fields.dbService,
				stockMarketService: tt.fields.stockMarketService,
				moneyFlowRepo:      tt.fields.moneyFlowRepo,
				client:             tt.fields.client,
				ctx:                tt.fields.ctx,
				running:            tt.fields.running,
				mu:                 tt.fields.mu,
			}
			got, err := s.SyncAndScanSingleStock(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncAndScanSingleStock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SyncAndScanSingleStock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncService_emitProgress(t *testing.T) {
	type fields struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
		client             *resty.Client
		ctx                context.Context
		running            bool
		mu                 sync.Mutex
	}
	type args struct {
		progress *SyncProgress
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				dbService:          tt.fields.dbService,
				stockMarketService: tt.fields.stockMarketService,
				moneyFlowRepo:      tt.fields.moneyFlowRepo,
				client:             tt.fields.client,
				ctx:                tt.fields.ctx,
				running:            tt.fields.running,
				mu:                 tt.fields.mu,
			}
			s.emitProgress(tt.args.progress)
		})
	}
}

func TestSyncService_fetchTickBatch(t *testing.T) {
	type fields struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
		client             *resty.Client
		ctx                context.Context
		running            bool
		mu                 sync.Mutex
	}
	type args struct {
		secid string
		pos   int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				dbService:          tt.fields.dbService,
				stockMarketService: tt.fields.stockMarketService,
				moneyFlowRepo:      tt.fields.moneyFlowRepo,
				client:             tt.fields.client,
				ctx:                tt.fields.ctx,
				running:            tt.fields.running,
				mu:                 tt.fields.mu,
			}
			got, err := s.fetchTickBatch(tt.args.secid, tt.args.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchTickBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchTickBatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncService_httpGetKlinesWithHeaders(t *testing.T) {
	type fields struct {
		dbService          *DBService
		stockMarketService *StockMarketService
		moneyFlowRepo      *repositories.MoneyFlowRepository
		client             *resty.Client
		ctx                context.Context
		running            bool
		mu                 sync.Mutex
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				dbService:          tt.fields.dbService,
				stockMarketService: tt.fields.stockMarketService,
				moneyFlowRepo:      tt.fields.moneyFlowRepo,
				client:             tt.fields.client,
				ctx:                tt.fields.ctx,
				running:            tt.fields.running,
				mu:                 tt.fields.mu,
			}
			got, err := s.httpGetKlinesWithHeaders(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("httpGetKlinesWithHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpGetKlinesWithHeaders() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWatchlistService_AddToWatchlist(t *testing.T) {
	type fields struct {
		repo repositories.WatchlistRepository
	}
	type args struct {
		stock *models.StockData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &WatchlistService{
				repo: tt.fields.repo,
			}
			if err := s.AddToWatchlist(tt.args.stock); (err != nil) != tt.wantErr {
				t.Errorf("AddToWatchlist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWatchlistService_GetWatchlist(t *testing.T) {
	type fields struct {
		repo repositories.WatchlistRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*models.StockData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &WatchlistService{
				repo: tt.fields.repo,
			}
			got, err := s.GetWatchlist()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWatchlist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWatchlist() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWatchlistService_RemoveFromWatchlist(t *testing.T) {
	type fields struct {
		repo repositories.WatchlistRepository
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &WatchlistService{
				repo: tt.fields.repo,
			}
			if err := s.RemoveFromWatchlist(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("RemoveFromWatchlist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_calculateEMA(t *testing.T) {
	type args struct {
		data   []float64
		period int
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateEMA(tt.args.data, tt.args.period); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateEMA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateMA(t *testing.T) {
	type args struct {
		klineData []*models.KLineData
		period    int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateMA(tt.args.klineData, tt.args.period); got != tt.want {
				t.Errorf("calculateMA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateMACD(t *testing.T) {
	type args struct {
		data         []float64
		fastPeriod   int
		slowPeriod   int
		signalPeriod int
	}
	tests := []struct {
		name  string
		args  args
		want  []float64
		want1 []float64
		want2 []float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := calculateMACD(tt.args.data, tt.args.fastPeriod, tt.args.slowPeriod, tt.args.signalPeriod)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateMACD() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("calculateMACD() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("calculateMACD() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_calculateRSI(t *testing.T) {
	type args struct {
		data   []float64
		period int
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateRSI(tt.args.data, tt.args.period); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateRSI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateSMA(t *testing.T) {
	type args struct {
		data   []float64
		period int
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateSMA(tt.args.data, tt.args.period); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateSMA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractSectionImpl(t *testing.T) {
	type args struct {
		text        string
		startMarker string
		endMarker   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractSectionImpl(tt.args.text, tt.args.startMarker, tt.args.endMarker); got != tt.want {
				t.Errorf("extractSectionImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFloat(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFloat(tt.args.v); got != tt.want {
				t.Errorf("getFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getHealthSummary(t *testing.T) {
	type args struct {
		score int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHealthSummary(tt.args.score); got != tt.want {
				t.Errorf("getHealthSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getInt64(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInt64(tt.args.v); got != tt.want {
				t.Errorf("getInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getString(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getString(tt.args.v); got != tt.want {
				t.Errorf("getString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_minInt(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minInt(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("minInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeDashscopeBaseURL(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := normalizeDashscopeBaseURL(tt.args.in)
			if got != tt.want {
				t.Errorf("normalizeDashscopeBaseURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("normalizeDashscopeBaseURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_parseAlertConditions(t *testing.T) {
	type args struct {
		jsonStr string
	}
	tests := []struct {
		name    string
		args    args
		want    repositories.PriceAlertConditions
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAlertConditions(tt.args.jsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAlertConditions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAlertConditions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseMoney(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseMoney(tt.args.s); got != tt.want {
				t.Errorf("parseMoney() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePrice(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePrice(tt.args.s); got != tt.want {
				t.Errorf("parsePrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseRawTicks(t *testing.T) {
	type args struct {
		rawTicks []string
	}
	tests := []struct {
		name string
		args args
		want []TickData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRawTicks(tt.args.rawTicks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRawTicks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_robustParseDrawings(t *testing.T) {
	type args struct {
		jsonStr string
	}
	tests := []struct {
		name string
		args args
		want []models.TechnicalDrawing
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := robustParseDrawings(tt.args.jsonStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("robustParseDrawings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_truncateString(t *testing.T) {
	type args struct {
		s   string
		max int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateString(tt.args.s, tt.args.max); got != tt.want {
				t.Errorf("truncateString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_zapLogger_Printf(t *testing.T) {
	type fields struct {
		Logger *zap.Logger
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := zapLogger{
				Logger: tt.fields.Logger,
			}
			l.Printf(tt.args.format, tt.args.args...)
		})
	}
}
