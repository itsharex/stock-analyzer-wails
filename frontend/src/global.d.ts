import type { StockData, AnalysisReport } from './types'

declare global {
  interface StockData extends import('./types').StockData {}
  interface AnalysisReport extends import('./types').AnalysisReport {}
}

export {}
