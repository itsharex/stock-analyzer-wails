import { useCallback } from 'react'
import type { StockData, AnalysisReport, AppConfig } from '../types'

export const useWailsAPI = () => {
  const getStockData = useCallback(async (code: string): Promise<StockData> => {
    // @ts-ignore
    return window.go.main.App.GetStockData(code)
  }, [])

  const analyzeStock = useCallback(async (code: string): Promise<AnalysisReport> => {
    // @ts-ignore
    return window.go.main.App.AnalyzeStock(code)
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

  return {
    getStockData,
    analyzeStock,
    searchStock,
    getConfig,
    saveConfig,
    addToWatchlist,
    removeFromWatchlist,
    getWatchlist,
  }
}
