import { useCallback } from 'react'
import type { StockData, AnalysisReport, AppConfig, KLineData, TechnicalAnalysisResult, IntradayResponse, MoneyFlowResponse, HealthCheckResult, EntryStrategyResult, StockDetail } from '../types'
import { StreamIntradayData } from '../../wailsjs/go/main/App'

export const useWailsAPI = () => {
  const getStockData = useCallback(async (code: string): Promise<StockData> => {
    // @ts-ignore
    return window.go.main.App.GetStockData(code)
  }, [])

	  const getKLineData = useCallback(async (code: string, limit: number, period: string = 'daily'): Promise<KLineData[]> => {
	    // @ts-ignore
	    return window.go.main.App.GetKLineData(code, limit, period)
	  }, [])
	
const getIntradayData = useCallback(async (code: string): Promise<IntradayResponse> => {
    // @ts-ignore
    return window.go.main.App.GetIntradayData(code)
  }, [])

	const getMoneyFlowData = useCallback(async (code: string): Promise<MoneyFlowResponse> => {
	    // @ts-ignore
	    return window.go.main.App.GetMoneyFlowData(code)
	  }, [])
	
	  const streamIntradayData = useCallback(async (code: string): Promise<void> => {
	    return StreamIntradayData(code)
	  }, [])

  const getStockDetail = useCallback(async (code: string): Promise<StockDetail> => {
    // @ts-ignore
    return window.go.main.App.GetStockDetail(code)
  }, [])

  const getStockHealthCheck = useCallback(async (code: string): Promise<HealthCheckResult> => {
    // @ts-ignore
    return window.go.main.App.GetStockHealthCheck(code)
  }, [])

  const batchAnalyzeStocks = useCallback(async (codes: string[], role: string): Promise<void> => {
    // @ts-ignore
    return window.go.main.App.BatchAnalyzeStocks(codes, role)
  }, [])

  const analyzeStock = useCallback(async (code: string): Promise<AnalysisReport> => {
    // @ts-ignore
    return window.go.main.App.AnalyzeStock(code)
  }, [])

  const analyzeTechnical = useCallback(async (code: string, period: string, role: string = 'technical'): Promise<TechnicalAnalysisResult> => {
    // @ts-ignore
    return window.go.main.App.AnalyzeTechnical(code, period, role)
  }, [])

  const searchStock = useCallback(async (keyword: string): Promise<StockData[]> => {
    // @ts-ignore
    return window.go.main.App.SearchStock(keyword)
  }, [])

  const getConfig = useCallback(async (): Promise<AppConfig> => {
    // @ts-ignore
    return window.go.main.App.GetConfig()
  }, [])

  const saveConfig = useCallback(async (config: AppConfig): Promise<void> => {
    // @ts-ignore
    return window.go.main.App.SaveConfig(config)
  }, [])

  // Watchlist API
  const addToWatchlist = useCallback(async (stock: StockData): Promise<void> => {
    // @ts-ignore
    return window.go.main.App.AddToWatchlist(stock)
  }, [])

  const removeFromWatchlist = useCallback(async (code: string): Promise<void> => {
    // @ts-ignore
    return window.go.main.App.RemoveFromWatchlist(code)
  }, [])

  const getWatchlist = useCallback(async (): Promise<StockData[]> => {
    // @ts-ignore
    return window.go.main.App.GetWatchlist()
  }, [])

  const getAlertHistory = useCallback(async (stockCode: string, limit: number): Promise<any[]> => {
    // @ts-ignore
    return window.go.main.App.GetAlertHistory(stockCode, limit)
  }, [])

  const getAlertConfig = useCallback(async (): Promise<any> => {
    // @ts-ignore
    return window.go.main.App.GetAlertConfig()
  }, [])

  const updateAlertConfig = useCallback(async (config: any): Promise<void> => {
    // @ts-ignore
    return window.go.main.App.UpdateAlertConfig(config)
  }, [])

  const getActiveAlerts = useCallback(async (): Promise<any[]> => {
    // @ts-ignore
    return window.go.main.App.GetActiveAlerts()
  }, [])

  const removeAlert = useCallback(async (stockCode: string, alertType: string, price: number): Promise<void> => {
    // @ts-ignore
    return window.go.main.App.RemoveAlert(stockCode, alertType, price)
  }, [])

  const analyzeEntryStrategy = useCallback(async (code: string): Promise<EntryStrategyResult> => {
    // @ts-ignore
    return window.go.main.App.AnalyzeEntryStrategy(code)
  }, [])

  const addPosition = useCallback(async (pos: any): Promise<void> => {
    // @ts-ignore
    return window.go.main.App.AddPosition(pos)
  }, [])

  const getPositions = useCallback(async (): Promise<Record<string, any>> => {
    // @ts-ignore
    return window.go.main.App.GetPositions()
  }, [])

  const removePosition = useCallback(async (code: string): Promise<void> => {
    // @ts-ignore
    return window.go.main.App.RemovePosition(code)
  }, [])

		return {
		    getStockData,
		    getIntradayData,
		    getMoneyFlowData,
		    streamIntradayData,
		    getStockDetail,
	    getStockHealthCheck,
    batchAnalyzeStocks,
    getKLineData,
	    analyzeStock,
	    analyzeTechnical,
	    searchStock,
	    getConfig,
	    saveConfig,
	    addToWatchlist,
	    removeFromWatchlist,
	    getWatchlist,
	    getAlertHistory,
	    getAlertConfig,
	    updateAlertConfig,
	    getActiveAlerts,
	    removeAlert,
	    analyzeEntryStrategy,
	    addPosition,
	    getPositions,
	    removePosition,
	  }
}

export default useWailsAPI
