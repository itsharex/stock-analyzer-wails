package main

import (
	"context"
	"fmt"
	"stock-analyzer-wails/models"
	"stock-analyzer-wails/services"
)

// App 应用程序结构
type App struct {
	ctx          context.Context
	stockService *services.StockService
	aiService    *services.AIService
}

// NewApp 创建新的App应用程序
func NewApp() *App {
	return &App{
		stockService: services.NewStockService(),
		aiService:    services.NewAIService(),
	}
}

// startup 在应用程序启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetStockData 获取股票数据（暴露给前端的方法）
func (a *App) GetStockData(code string) (*models.StockData, error) {
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	
	stock, err := a.stockService.GetStockByCode(code)
	if err != nil {
		return nil, fmt.Errorf("获取股票数据失败: %w", err)
	}
	
	return stock, nil
}

// AnalyzeStock 分析股票（暴露给前端的方法）
func (a *App) AnalyzeStock(code string) (*models.AnalysisReport, error) {
	if code == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	
	// 1. 获取股票数据
	stock, err := a.stockService.GetStockByCode(code)
	if err != nil {
		return nil, fmt.Errorf("获取股票数据失败: %w", err)
	}
	
	// 2. 使用AI分析股票
	report, err := a.aiService.AnalyzeStock(stock)
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}
	
	return report, nil
}

// QuickAnalyze 快速分析（暴露给前端的方法）
func (a *App) QuickAnalyze(code string) (string, error) {
	if code == "" {
		return "", fmt.Errorf("股票代码不能为空")
	}
	
	// 1. 获取股票数据
	stock, err := a.stockService.GetStockByCode(code)
	if err != nil {
		return "", fmt.Errorf("获取股票数据失败: %w", err)
	}
	
	// 2. 快速分析
	analysis, err := a.aiService.QuickAnalyze(stock)
	if err != nil {
		return "", fmt.Errorf("快速分析失败: %w", err)
	}
	
	return analysis, nil
}

// SearchStock 搜索股票（暴露给前端的方法）
func (a *App) SearchStock(keyword string) ([]*models.StockData, error) {
	if keyword == "" {
		return nil, fmt.Errorf("搜索关键词不能为空")
	}
	
	results, err := a.stockService.SearchStock(keyword)
	if err != nil {
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
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	stocks, err := a.stockService.GetStockList(pageNum, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取股票列表失败: %w", err)
	}
	
	return stocks, nil
}

// Greet 示例方法：返回问候语
func (a *App) Greet(name string) string {
	return fmt.Sprintf("你好 %s, 欢迎使用A股股票分析AI-Agent！", name)
}
