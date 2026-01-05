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
    try {
      console.log('开始添加股票到自选股:', stock.code, stock.name)
      // @ts-ignore
      const result = await window.go.main.App.AddToWatchlist(stock)
      console.log('成功添加股票到自选股:', stock.code)
      return result
    } catch (error) {
      console.error('添加股票到自选股失败:', error)
      throw error
    }
  }, [])

  const removeFromWatchlist = useCallback(async (code: string): Promise<void> => {
    try {
      console.log('开始从自选股移除股票:', code)
      // @ts-ignore
      const result = await window.go.main.App.RemoveFromWatchlist(code)
      console.log('成功从自选股移除股票:', code)
      return result
    } catch (error) {
      console.error('从自选股移除股票失败:', error)
      throw error
    }
  }, [])

  const getWatchlist = useCallback(async (): Promise<StockData[]> => {
    try {
      console.log('开始获取自选股列表')
      // @ts-ignore
      const result = await window.go.main.App.GetWatchlist()
      console.log('成功获取自选股列表:', result.length, '只股票')
      return result
    } catch (error) {
      console.error('获取自选股列表失败:', error)
      throw error
    }
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
