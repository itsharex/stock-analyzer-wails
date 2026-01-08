import { useCallback } from 'react'
import type { StockData, AnalysisReport, AppConfig, KLineData, TechnicalAnalysisResult, IntradayResponse, MoneyFlowResponse, HealthCheckResult, EntryStrategyResult, StockDetail, BacktestResult } from '../types'
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

  const BacktestSimpleMA = useCallback(async (code: string, shortPeriod: number, longPeriod: number, initialCapital: number, startDate: string, endDate: string): Promise<BacktestResult> => {
    // @ts-ignore
    return window.go.main.App.BacktestSimpleMA(code, shortPeriod, longPeriod, initialCapital, startDate, endDate)
  }, [])

  const BacktestMACD = useCallback(async (code: string, fastPeriod: number, slowPeriod: number, signalPeriod: number, initialCapital: number, startDate: string, endDate: string): Promise<BacktestResult> => {
    // @ts-ignore
    return window.go.main.App.BacktestMACD(code, fastPeriod, slowPeriod, signalPeriod, initialCapital, startDate, endDate)
  }, [])

  const SyncStockData = useCallback(async (code: string, startDate: string, endDate: string): Promise<any> => {
    // @ts-ignore
    return window.go.main.App.SyncStockData(code, startDate, endDate)
  }, [])

  const GetDataSyncStats = useCallback(async (): Promise<any> => {
    // @ts-ignore
    return window.go.main.App.GetDataSyncStats()
  }, [])

  const BatchSyncStockData = useCallback(async (codes: string[], startDate: string, endDate: string): Promise<void> => {
    // @ts-ignore
    return window.go.main.App.BatchSyncStockData(codes, startDate, endDate)
  }, [])

  const ClearStockCache = useCallback(async (code: string): Promise<void> => {
    // @ts-ignore
    return window.go.main.App.ClearStockCache(code)
  }, [])

  // Sync History API
  const getAllSyncHistory = useCallback(async (limit: number, offset: number): Promise<any[]> => {
    // @ts-ignore
    return window.go.main.App.GetAllSyncHistory(limit, offset)
  }, [])

  const getSyncHistoryCount = useCallback(async (): Promise<number> => {
    // @ts-ignore
    return window.go.main.App.GetSyncHistoryCount()
  }, [])

  const clearAllSyncHistory = useCallback(async (): Promise<void> => {
    // @ts-ignore
    return window.go.main.App.ClearAllSyncHistory()
  }, [])

  const getSyncHistoryByCode = useCallback(async (code: string, limit: number): Promise<any[]> => {
    // @ts-ignore
    return window.go.main.App.GetSyncHistoryByCode(code, limit)
  }, [])

  const getSyncedKLineData = useCallback(async (code: string, startDate: string, endDate: string, page: number, pageSize: number): Promise<{ data: any[], total: number }> => {
    // @ts-ignore
    if (!window.go?.main?.App?.GetSyncedKLineData) {
      throw new Error('GetSyncedKLineData 方法不可用，请确保已运行 wails dev 重新生成绑定文件')
    }
    try {
      console.log('开始调用 GetSyncedKLineData, 参数:', { code, startDate, endDate, page, pageSize })

      // @ts-ignore
      const result = await window.go.main.App.GetSyncedKLineData(code, startDate, endDate, page, pageSize)

      console.log('GetSyncedKLineData 返回的原始 result:', result)
      console.log('result 的类型:', typeof result)
      console.log('result 是否为 null:', result === null)
      console.log('result 是否为 undefined:', result === undefined)

      // 检查 result 是否为 null 或 undefined
      if (result == null) {
        console.error('GetSyncedKLineData 返回了 null 或 undefined')
        throw new Error('GetSyncedKLineData 返回了 null 或 undefined，请检查后端实现')
      }

      // 新的返回格式是对象 { data: [...], total: number }
      const data = result.data || []
      const total = result.total || 0

      console.log('解析后的 data:', data)
      console.log('解析后的 total:', total)
      console.log('data 的类型:', typeof data)
      console.log('data 是否为数组:', Array.isArray(data))

      // 验证数据有效性
      if (!Array.isArray(data)) {
        console.error('返回的 data 字段不是数组:', data)
        throw new Error('返回的数据格式错误，data 字段应该是数组')
      }

      return { data, total }
    } catch (error: any) {
      console.error('调用 GetSyncedKLineData 失败:', error)
      throw error
    }
  }, [])

  // 市场股票管理 API
  const syncAllStocks = useCallback(async (): Promise<any> => {
    // @ts-ignore
    return window.go.main.App.SyncAllStocks()
  }, [])

  const getStocksList = useCallback(async (page: number, pageSize: number, search: string): Promise<any> => {
    // @ts-ignore
    return window.go.main.App.GetStocksList(page, pageSize, search)
  }, [])

  const getSyncStats = useCallback(async (): Promise<any> => {
    // @ts-ignore
    return window.go.main.App.GetSyncStats()
  }, [])

  // 策略管理 API
  const CreateStrategy = useCallback(async (name: string, description: string, strategyType: string, parameters: Record<string, any>) => {
    // @ts-ignore
    return window.go.main.App.CreateStrategy(name, description, strategyType, parameters)
  }, [])

  const UpdateStrategy = useCallback(async (id: number, name: string, description: string, strategyType: string, parameters: Record<string, any>) => {
    // @ts-ignore
    return window.go.main.App.UpdateStrategy(id, name, description, strategyType, parameters)
  }, [])

  const DeleteStrategy = useCallback(async (id: number) => {
    // @ts-ignore
    return window.go.main.App.DeleteStrategy(id)
  }, [])

  const GetStrategy = useCallback(async (id: number) => {
    // @ts-ignore
    return window.go.main.App.GetStrategy(id)
  }, [])

  const GetAllStrategies = useCallback(async () => {
    try {
      // @ts-ignore
      if (!window.go?.main?.App?.GetAllStrategies) {
        throw new Error('GetAllStrategies 方法不可用，请确保已运行 wails dev 重新生成绑定文件')
      }
      console.log('调用 GetAllStrategies...')
      // @ts-ignore
      const result = await window.go.main.App.GetAllStrategies()
      console.log('GetAllStrategies 返回结果:', result)
      return result
    } catch (error: any) {
      console.error('调用 GetAllStrategies 失败:', error)
      throw error
    }
  }, [])

  const GetStrategyTypes = useCallback(async () => {
    // @ts-ignore
    return window.go.main.App.GetStrategyTypes()
  }, [])

  const UpdateStrategyBacktestResult = useCallback(async (id: number, backtestResult: Record<string, any>) => {
    // @ts-ignore
    return window.go.main.App.UpdateStrategyBacktestResult(id, backtestResult)
  }, [])

  // Price Alert API
  const getAllPriceAlerts = useCallback(async () => {
    // @ts-ignore
    return window.go.main.App.PriceAlertController?.GetAllAlerts()
  }, [])

  const getActivePriceAlerts = useCallback(async () => {
    // @ts-ignore
    return window.go.main.App.PriceAlertController?.GetActiveAlerts()
  }, [])

  const getPriceAlertsByCode = useCallback(async (code: string) => {
    // @ts-ignore
    return window.go.main.App.PriceAlertController?.GetAlertsByStockCode(code)
  }, [])

  const getPriceAlertTemplates = useCallback(async () => {
    // @ts-ignore
    return window.go.main.App.PriceAlertController?.GetAllTemplates()
  }, [])

  const getPriceAlertHistory = useCallback(async (code: string, limit: number) => {
    // @ts-ignore
    return window.go.main.App.PriceAlertController?.GetTriggerHistory(code, limit)
  }, [])

  const createPriceAlert = useCallback(async (jsonData: string) => {
    // @ts-ignore
    return window.go.main.App.PriceAlertController?.CreateAlert(jsonData)
  }, [])

  const updatePriceAlert = useCallback(async (jsonData: string) => {
    // @ts-ignore
    return window.go.main.App.PriceAlertController?.UpdateAlert(jsonData)
  }, [])

  const deletePriceAlert = useCallback(async (id: number) => {
    // @ts-ignore
    return window.go.main.App.PriceAlertController?.DeleteAlert(id)
  }, [])

  const togglePriceAlert = useCallback(async (id: number, isActive: boolean) => {
    // @ts-ignore
    return window.go.main.App.PriceAlertController?.ToggleAlertStatus(id, isActive)
  }, [])

  const createPriceAlertFromTemplate = useCallback(async (templateId: string, stockCode: string, stockName: string, paramsJson: string) => {
    // @ts-ignore
    return window.go.main.App.PriceAlertController?.CreateAlertFromTemplate(templateId, stockCode, stockName, paramsJson)
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
    BacktestSimpleMA,
    BacktestMACD,
    SyncStockData,
    GetDataSyncStats,
    BatchSyncStockData,
    ClearStockCache,
    // Sync History API
    getAllSyncHistory,
    getSyncHistoryCount,
    clearAllSyncHistory,
    getSyncHistoryByCode,
    getSyncedKLineData,
    // Market Stock API
    syncAllStocks,
    getStocksList,
    getSyncStats,
    // Strategy Management API
    CreateStrategy,
    UpdateStrategy,
    DeleteStrategy,
    GetStrategy,
    GetAllStrategies,
    GetStrategyTypes,
    UpdateStrategyBacktestResult,
    // Price Alert API
    getAllPriceAlerts,
    getActivePriceAlerts,
    getPriceAlertsByCode,
    getPriceAlertTemplates,
    getPriceAlertHistory,
    createPriceAlert,
    updatePriceAlert,
    deletePriceAlert,
    togglePriceAlert,
    createPriceAlertFromTemplate,
  }
}

export default useWailsAPI
