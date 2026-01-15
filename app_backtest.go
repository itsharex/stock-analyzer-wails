package main

import (
	"fmt"
	"stock-analyzer-wails/models"
)

// --- 回测功能 ---

// BacktestSimpleMA 使用简单双均线策略
func (a *App) BacktestSimpleMA(code string, shortPeriod int, longPeriod int, initialCapital float64, startDate string, endDate string) (*models.BacktestResult, error) {
	if a.backtestService == nil {
		return nil, fmt.Errorf("回测服务未初始化")
	}
	return a.backtestService.BacktestSimpleMA(code, shortPeriod, longPeriod, initialCapital, startDate, endDate)
}

// BacktestMACD 使用MACD策略
func (a *App) BacktestMACD(code string, fastPeriod int, slowPeriod int, signalPeriod int, initialCapital float64, startDate string, endDate string) (*models.BacktestResult, error) {
	if a.backtestService == nil {
		return nil, fmt.Errorf("回测服务未初始化")
	}
	return a.backtestService.BacktestMACD(code, fastPeriod, slowPeriod, signalPeriod, initialCapital, startDate, endDate)
}

// BacktestRSI 使用RSI策略
func (a *App) BacktestRSI(code string, period int, buyThreshold float64, sellThreshold float64, initialCapital float64, startDate string, endDate string) (*models.BacktestResult, error) {
	if a.backtestService == nil {
		return nil, fmt.Errorf("回测服务未初始化")
	}
	return a.backtestService.BacktestRSI(code, period, buyThreshold, sellThreshold, initialCapital, startDate, endDate)
}
