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
	ctx          context.Context
	stockService *services.StockService
	aiService    *services.AIService
	aiInitErr    error
}

// NewApp 创建新的App应用程序
func NewApp() *App {
	return &App{
		stockService: services.NewStockService(),
		aiService:    nil,
	}
}

// startup 在应用程序启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	start := time.Now()
	cfg, err := services.LoadDashscopeConfig()
	if err != nil {
		a.aiInitErr = err
		logger.Error("AI 配置加载失败",
			zap.String("module", "app"),
			zap.String("op", "startup.LoadDashscopeConfig"),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return
	}

	aiSvc, err := services.NewAIService(cfg)
	if err != nil {
		a.aiInitErr = err
		logger.Error("AI 服务初始化失败",
			zap.String("module", "app"),
			zap.String("op", "startup.NewAIService"),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.String("model", cfg.Model),
			zap.String("base_url", cfg.BaseURL),
			zap.Error(err),
		)
		return
	}

	a.aiService = aiSvc
	a.aiInitErr = nil

	logger.Info("应用启动完成",
		zap.String("module", "app"),
		zap.String("op", "startup"),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		zap.Bool("ai_ready", true),
	)
}

// GetStockData 获取股票数据（暴露给前端的方法）
func (a *App) GetStockData(code string) (*models.StockData, error) {
	start := time.Now()
	if code == "" {
		logger.Warn("股票代码不能为空",
			zap.String("module", "app"),
			zap.String("op", "GetStockData"),
		)
		return nil, fmt.Errorf("股票代码不能为空")
	}
	
	stock, err := a.stockService.GetStockByCode(code)
	if err != nil {
		logger.Error("获取股票数据失败",
			zap.String("module", "app"),
			zap.String("op", "GetStockData"),
			zap.String("stock_code", code),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("获取股票数据失败: %w", err)
	}
	
	return stock, nil
}

// AnalyzeStock 分析股票（暴露给前端的方法）
func (a *App) AnalyzeStock(code string) (*models.AnalysisReport, error) {
	start := time.Now()
	if code == "" {
		logger.Warn("股票代码不能为空",
			zap.String("module", "app"),
			zap.String("op", "AnalyzeStock"),
		)
		return nil, fmt.Errorf("股票代码不能为空")
	}
	if a.aiService == nil {
		if a.aiInitErr != nil {
			logger.Error("AI服务未正确初始化",
				zap.String("module", "app"),
				zap.String("op", "AnalyzeStock"),
				zap.String("stock_code", code),
				zap.Error(a.aiInitErr),
			)
			return nil, fmt.Errorf("AI服务未正确初始化: %w", a.aiInitErr)
		}
		logger.Error("AI服务未正确初始化",
			zap.String("module", "app"),
			zap.String("op", "AnalyzeStock"),
			zap.String("stock_code", code),
		)
		return nil, fmt.Errorf("AI服务未正确初始化")
	}
	
	// 1. 获取股票数据
	stock, err := a.stockService.GetStockByCode(code)
	if err != nil {
		logger.Error("获取股票数据失败",
			zap.String("module", "app"),
			zap.String("op", "AnalyzeStock.GetStockByCode"),
			zap.String("stock_code", code),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("获取股票数据失败: %w", err)
	}
	
	// 2. 使用AI分析股票
	report, err := a.aiService.AnalyzeStock(stock)
	if err != nil {
		logger.Error("AI分析失败",
			zap.String("module", "app"),
			zap.String("op", "AnalyzeStock.AnalyzeStock"),
			zap.String("stock_code", code),
			zap.String("stock_name", stock.Name),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}
	
	return report, nil
}

// QuickAnalyze 快速分析（暴露给前端的方法）
func (a *App) QuickAnalyze(code string) (string, error) {
	start := time.Now()
	if code == "" {
		logger.Warn("股票代码不能为空",
			zap.String("module", "app"),
			zap.String("op", "QuickAnalyze"),
		)
		return "", fmt.Errorf("股票代码不能为空")
	}
	if a.aiService == nil {
		if a.aiInitErr != nil {
			logger.Error("AI服务未正确初始化",
				zap.String("module", "app"),
				zap.String("op", "QuickAnalyze"),
				zap.String("stock_code", code),
				zap.Error(a.aiInitErr),
			)
			return "", fmt.Errorf("AI服务未正确初始化: %w", a.aiInitErr)
		}
		logger.Error("AI服务未正确初始化",
			zap.String("module", "app"),
			zap.String("op", "QuickAnalyze"),
			zap.String("stock_code", code),
		)
		return "", fmt.Errorf("AI服务未正确初始化")
	}
	
	// 1. 获取股票数据
	stock, err := a.stockService.GetStockByCode(code)
	if err != nil {
		logger.Error("获取股票数据失败",
			zap.String("module", "app"),
			zap.String("op", "QuickAnalyze.GetStockByCode"),
			zap.String("stock_code", code),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return "", fmt.Errorf("获取股票数据失败: %w", err)
	}
	
	// 2. 快速分析
	analysis, err := a.aiService.QuickAnalyze(stock)
	if err != nil {
		logger.Error("快速分析失败",
			zap.String("module", "app"),
			zap.String("op", "QuickAnalyze.QuickAnalyze"),
			zap.String("stock_code", code),
			zap.String("stock_name", stock.Name),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return "", fmt.Errorf("快速分析失败: %w", err)
	}
	
	return analysis, nil
}

// SearchStock 搜索股票（暴露给前端的方法）
func (a *App) SearchStock(keyword string) ([]*models.StockData, error) {
	start := time.Now()
	if keyword == "" {
		logger.Warn("搜索关键词不能为空",
			zap.String("module", "app"),
			zap.String("op", "SearchStock"),
		)
		return nil, fmt.Errorf("搜索关键词不能为空")
	}
	
	results, err := a.stockService.SearchStock(keyword)
	if err != nil {
		logger.Error("搜索失败",
			zap.String("module", "app"),
			zap.String("op", "SearchStock"),
			zap.String("keyword", keyword),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("搜索失败: %w", err)
	}
	
	// 限制返回结果数量
	if len(results) > 20 {
		results = results[:20]
	}
	
	return results, nil
}

// GetStockList 获取股票列表（暴露给前端的方法）
func (a *App) GetStockList(pageNum, pageSize int) ([]*models.StockData, error) {
	start := time.Now()
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	stocks, err := a.stockService.GetStockList(pageNum, pageSize)
	if err != nil {
		logger.Error("获取股票列表失败",
			zap.String("module", "app"),
			zap.String("op", "GetStockList"),
			zap.Int("page_num", pageNum),
			zap.Int("page_size", pageSize),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("获取股票列表失败: %w", err)
	}
	
	return stocks, nil
}

// Greet 示例方法：返回问候语
func (a *App) Greet(name string) string {
	return fmt.Sprintf("你好 %s, 欢迎使用A股股票分析AI-Agent！", name)
}
