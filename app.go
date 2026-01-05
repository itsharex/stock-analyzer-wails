package main

import (
	"context"
	"fmt"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/services"
	"time"

	"stock-analyzer-wails/internal/logger"

	"go.uber.org/zap"
)

// App 应用程序结构
type App struct {
	ctx              context.Context
	stockService     *services.StockService
	aiService        *services.AIService
	watchlistService *services.WatchlistService
	aiInitErr        error
}

// NewApp 创建新的App应用程序
func NewApp() *App {
	watchlistSvc, _ := services.NewWatchlistService()
	return &App{
		stockService:     services.NewStockService(),
		aiService:        nil,
		watchlistService: watchlistSvc,
	}
}

// startup 在应用程序启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.initAIService()
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

// SearchStock 搜索股票
func (a *App) SearchStock(keyword string) ([]*models.StockData, error) {
	return a.stockService.SearchStock(keyword)
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
