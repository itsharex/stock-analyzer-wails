/**
 * 股票数据类型定义
 */
export interface StockData {
  code: string              // 股票代码
  name: string              // 股票名称
  price: number             // 最新价
  change: number            // 涨跌额
  changeRate: number        // 涨跌幅 (%)
  volume: number            // 成交量 (手)
  amount: number            // 成交额
  high: number              // 最高价
  low: number               // 最低价
  open: number              // 今开
  preClose: number          // 昨收
  amplitude: number         // 振幅 (%)
  turnover: number          // 换手率 (%)
  pe: number                // 市盈率
  pb: number                // 市净率
  totalMV: number           // 总市值
  circMV: number            // 流通市值
}

/**
 * MACD 指标
 */
export interface MACD {
  dif: number
  dea: number
  bar: number
}

/**
 * KDJ 指标
 */
export interface KDJ {
  k: number
  d: number
  j: number
}

/**
 * K线数据点
 */
export interface KLineData {
  time: string              // 时间 (YYYY-MM-DD)
  open: number              // 开盘价
  high: number              // 最高价
  low: number               // 最低价
  close: number             // 收盘价
  volume: number            // 成交量
  macd?: MACD               // MACD 指标
  kdj?: KDJ                 // KDJ 指标
  rsi?: number              // RSI 指标
}

/**
 * AI 分析报告类型定义
 */
export interface AnalysisReport {
  stockCode: string         // 股票代码
  stockName: string         // 股票名称
  summary: string           // 分析摘要
  fundamentals: string      // 基本面分析
  technical: string         // 技术面分析
  recommendation: string    // 投资建议
  riskLevel: string         // 风险等级
  targetPrice: string       // 目标价位
  generatedAt: string       // 生成时间
}

/**
 * AI 识别的绘图数据
 */
export interface TechnicalDrawing {
  type: 'support' | 'resistance' | 'trendline'
  price?: number
  start?: string
  end?: string
  startPrice?: number
  endPrice?: number
  label: string
}

/**
 * 雷达图评分数据
 */
export interface RadarData {
  scores: Record<string, number>
  reasons: Record<string, string>
}

/**
 * 深度技术分析结果
 */
export interface TechnicalAnalysisResult {
  analysis: string
  drawings: TechnicalDrawing[]
  riskScore: number
  actionAdvice: string
  radarData?: RadarData
}

/**
 * 系统配置类型定义
 */
export interface AppConfig {
  provider: string
  apiKey: string
  baseUrl: string
  model: string
  providerModels: Record<string, string[]>
}

/**
 * 导航菜单项
 */
export type NavItem = 'analysis' | 'watchlist' | 'settings'
