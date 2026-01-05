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
 * KLineData K线数据点
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
 * 分时数据点
 */
export interface IntradayData {
  time: string              // 时间 (HH:MM)
  price: number             // 价格
  avgPrice: number          // 均价
  volume: number            // 成交量
  preClose: number          // 昨收价
}

/**
 * 分时数据响应结构
 */
export interface IntradayResponse {
  data: IntradayData[]
  preClose: number
}

/**
 * 资金流向数据点
 */
export interface MoneyFlowData {
  time: string
  superLarge: number
  large: number
  medium: number
  small: number
  mainNet: number
  signal?: string
}

/**
 * 资金流向响应结构
 */
export interface MoneyFlowResponse {
  data: MoneyFlowData[]
  todayMain: number
  todayRetail: number
  status: string
  description: string
}

/**
 * 股票深度体检结果
 */
export interface HealthItem {
  category: string
  name: string
  value: string
  status: '正常' | '警告' | '异常'
  description: string
}

export interface HealthCheckResult {
  score: number
  status: string
  items: HealthItem[]
  summary: string
  riskLevel: string
  updatedAt: string
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
 * 智能交易计划
 */
export interface TradePlan {
  suggestedPosition: string
  stopLoss: number
  takeProfit: number
  riskRewardRatio: number
  strategy: string
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
  tradePlan?: TradePlan
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
export type NavItem = 'analysis' | 'watchlist' | 'alerts' | 'settings'

export interface CoreReason {
  type: 'fundamental' | 'technical' | 'money_flow'
  description: string
  threshold: string
}

export interface EntryStrategyResult {
  recommendation: string
  entryPriceRange: string
  initialPosition: string
  stopLossPrice: number
  takeProfitPrice: number
  coreReasons: CoreReason[]
  riskRewardRatio: number
  actionPlan: string
}

export interface TrailingStopConfig {
  enabled: boolean
  activationThreshold: number
  callbackRate: number
}

export interface Position {
  stockCode: string
  stockName: string
  entryPrice: number
  entryTime: string
  strategy: EntryStrategyResult
  trailingConfig: TrailingStopConfig
  currentStatus: 'holding' | 'closed'
  logicStatus: 'valid' | 'violated' | 'warning'
}
