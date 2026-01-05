import { useCallback } from 'react'
import type { StockData, AnalysisReport } from '../types'

/**
 * Wails API 调用 hooks
 * 用于封装与后端 Go 代码的通信
 */
export const useWailsAPI = () => {
  /**
   * 获取股票数据
   */
  const getStockData = useCallback(async (code: string): Promise<StockData> => {
    try {
      // @ts-ignore - Wails 运行时会自动生成这些方法
      const result = await window.go.main.App.GetStockData(code)
      return result
    } catch (error) {
      throw new Error(`获取股票数据失败: ${error instanceof Error ? error.message : String(error)}`)
    }
  }, [])

  /**
   * 分析股票
   */
  const analyzeStock = useCallback(async (code: string): Promise<AnalysisReport> => {
    try {
      // @ts-ignore - Wails 运行时会自动生成这些方法
      const result = await window.go.main.App.AnalyzeStock(code)
      return result
    } catch (error) {
      throw new Error(`AI分析失败: ${error instanceof Error ? error.message : String(error)}`)
    }
  }, [])

  /**
   * 快速分析
   */
  const quickAnalyze = useCallback(async (code: string): Promise<string> => {
    try {
      // @ts-ignore - Wails 运行时会自动生成这些方法
      const result = await window.go.main.App.QuickAnalyze(code)
      return result
    } catch (error) {
      throw new Error(`快速分析失败: ${error instanceof Error ? error.message : String(error)}`)
    }
  }, [])

  /**
   * 搜索股票
   */
  const searchStock = useCallback(async (keyword: string): Promise<StockData[]> => {
    try {
      // @ts-ignore - Wails 运行时会自动生成这些方法
      const result = await window.go.main.App.SearchStock(keyword)
      return result || []
    } catch (error) {
      throw new Error(`搜索失败: ${error instanceof Error ? error.message : String(error)}`)
    }
  }, [])

  /**
   * 获取股票列表
   */
  const getStockList = useCallback(async (pageNum: number = 1, pageSize: number = 20): Promise<StockData[]> => {
    try {
      // @ts-ignore - Wails 运行时会自动生成这些方法
      const result = await window.go.main.App.GetStockList(pageNum, pageSize)
      return result || []
    } catch (error) {
      throw new Error(`获取股票列表失败: ${error instanceof Error ? error.message : String(error)}`)
    }
  }, [])

  return {
    getStockData,
    analyzeStock,
    quickAnalyze,
    searchStock,
    getStockList,
  }
}
