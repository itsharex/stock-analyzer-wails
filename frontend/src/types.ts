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
 * 系统配置类型定义
 */
export interface AppConfig {
  apiKey: string
  baseUrl: string
  model: string
}

/**
 * 导航菜单项
 */
export type NavItem = 'analysis' | 'settings'
