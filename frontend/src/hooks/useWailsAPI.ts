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

  return {
    getStockData,
    analyzeStock,
    searchStock,
    getConfig,
    saveConfig,
  }
}
