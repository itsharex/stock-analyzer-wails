import { useCallback } from 'react'
import type { StockData, AnalysisReport, AppConfig, KLineData, TechnicalAnalysisResult } from '../types'

export const useWailsAPI = () => {
  const getStockData = useCallback(async (code: string): Promise<StockData> => {
    // @ts-ignore
    return window.go.main.App.GetStockData(code)
  }, [])

  const getKLineData = useCallback(async (code: string, limit: number, period: string = 'daily'): Promise<KLineData[]> => {
    // @ts-ignore
    return window.go.main.App.GetKLineData(code, limit, period)
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

  return {
    getStockData,
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
  }
}

export default useWailsAPI
