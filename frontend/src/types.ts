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
  circMV: number,            // 流通市值
  volumeRatio: number,         // 量比
  warrantRatio: number,        // 委比
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

export interface OrderBookItem {
  price: number
  volume: number
}

export interface OrderBook {
  // 兼容后端当前返回的字段（Go models.OrderBook）：buy5/sell5
  buy5?: OrderBookItem[]
  sell5?: OrderBookItem[]
  // 兼容旧前端字段：buy/sell + volume/amount
  buy?: OrderBookItem[]
  sell?: OrderBookItem[]
  volume?: number
  amount?: number
}

export interface FinancialSummary {
  roe: number
  net_profit_growth_rate: number
  gross_profit_margin: number
  total_market_value: number
  circulating_market_value: number
  dividend_yield: number
  report_date: string // ISO date string
}

export interface IndustryInfo {
  industry_name: string
  concept_names: string[]
  industry_pe: number
}

export interface StockDetail {
  // 后端当前返回为“扁平字段 + orderBook/financial_summary/industry_info”
  // 这里做兼容，避免运行时字段缺失导致白屏
  stockData?: StockData
  orderBook?: OrderBook
  financial?: FinancialSummary
  industry?: IndustryInfo
  // 允许扁平行情字段存在
  code?: string
  name?: string
  price?: number
  change?: number
  changeRate?: number
  volume?: number
  amount?: number
  high?: number
  low?: number
  open?: number
  preClose?: number
  amplitude?: number
  turnover?: number
  pe?: number
  pb?: number
  totalMV?: number
  circMV?: number
  volumeRatio?: number
  warrantRatio?: number
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


export interface TradeRecord {
  time: string;
  type: "BUY" | "SELL";
  price: number;
  volume: number;
  amount: number;
  commission: number;
  tax: number;
  profit: number;
}

// 市场股票数据
export interface StockMarketData {
  id: number;
  code: string;
  name: string;
  market: string; // SH, SZ, BJ
  fullCode: string; // 如 SH600519
  type: string; // 主板, 创业板, 科创板, 北交所
  isActive: number; // 1: 正常, 0: 退市/停牌
  price: number; // 最新价
  changeRate: number; // 涨跌幅(%)
  changeAmount: number; // 涨跌额
  volume: number; // 成交量(手)
  amount: number; // 成交额
  amplitude: number; // 振幅(%)
  high: number; // 最高价
  low: number; // 最低价
  open: number; // 开盘价
  preClose: number; // 昨收
  turnover: number; // 换手率(%)
  volumeRatio: number; // 量比
  pe: number; // 市盈率
  warrantRatio: number; // 委比(%)
  industry: string; // 所属行业
  region: string; // 地区
  board: string; // 板块
  totalMV: number; // 总市值
  circMV: number; // 流通市值
  updatedAt: string; // 最后更新时间
}

// 同步结果
export interface SyncStocksResult {
  total: number;
  processed: number;
  inserted: number;
  updated: number;
  duration: number;
  message: string;
}

export interface BacktestResult {
  strategyName: string;
  stockCode: string;
  startDate: string;
  endDate: string;
  initialCapital: number;
  finalCapital: number;
  totalReturn: number;
  annualizedReturn: number;
  maxDrawdown: number;
  winRate: number;
  tradeCount: number;
  trades: TradeRecord[];
  equityCurve: number[];
  equityDates: string[];
}

/**
 * 策略信号
 */
export interface StrategySignal {
  id: number;
  code: string;
  stockName: string; // 股票名称
  tradeDate: string;
  signalType: string;
  strategyName: string;
  score: number;
  details: string;
  aiScore: number;
  aiReason: string;
  riskLevel?: string;
  createdAt: string;
}

/**
 * AI 分析详情 (用于 UI 展示)
 */
export interface AIAnalysisDetail {
  score: number;
  reason: string;
  keywords: string[];
  mainConcentration: number; // 主力持仓集中度 (0-100)
  retailProfit: number;      // 散户获利盘 (0-100)
}
