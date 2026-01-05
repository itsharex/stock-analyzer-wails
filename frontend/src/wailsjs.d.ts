// Wails Runtime类型定义
declare module '*/wailsjs/go/main/App' {
  export function GetStockData(code: string): Promise<StockData>;
  export function AnalyzeStock(code: string): Promise<AnalysisReport>;
  export function QuickAnalyze(code: string): Promise<string>;
  export function SearchStock(keyword: string): Promise<StockData[]>;
  export function GetStockList(pageNum: number, pageSize: number): Promise<StockData[]>;
  export function Greet(name: string): Promise<string>;
}

// 数据模型类型定义
interface StockData {
  code: string;
  name: string;
  price: number;
  change: number;
  changeRate: number;
  volume: number;
  amount: number;
  high: number;
  low: number;
  open: number;
  preClose: number;
  amplitude: number;
  turnover: number;
  pe: number;
  pb: number;
  totalMV: number;
  circMV: number;
}

interface AnalysisReport {
  stockCode: string;
  stockName: string;
  summary: string;
  fundamentals: string;
  technical: string;
  recommendation: string;
  riskLevel: string;
  targetPrice: string;
  generatedAt: string;
}
