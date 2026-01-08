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

// 同步历史数据模型
interface SyncHistoryItem {
  id: number;
  stock_code: string;
  stock_name: string;
  sync_type: string;
  start_date: string;
  end_date: string;
  status: string;
  records_added: number;
  records_updated: number;
  duration: number;
  error_msg: string;
  created_at: string;
}

// window.go 全局类型声明
declare global {
  interface Window {
    go: {
      main: {
        App: {
          // 原有方法
          GetStockData(code: string): Promise<StockData>;
          AnalyzeStock(code: string): Promise<AnalysisReport>;
          QuickAnalyze(code: string): Promise<string>;
          SearchStock(keyword: string): Promise<StockData[]>;
          GetStockList(pageNum: number, pageSize: number): Promise<StockData[]>;
          Greet(name: string): Promise<string>;
          GetDataSyncStats(): Promise<any>;
          SyncStockData(code: string, startDate: string, endDate: string): Promise<any>;
          BatchSyncStockData(codes: string[], startDate: string, endDate: string): Promise<void>;
          ClearStockCache(code: string): Promise<void>;
          // 同步历史方法
          GetAllSyncHistory(limit: number, offset: number): Promise<SyncHistoryItem[]>;
          GetSyncHistoryByCode(code: string, limit: number): Promise<SyncHistoryItem[]>;
          GetSyncHistoryCount(): Promise<number>;
          ClearAllSyncHistory(): Promise<void>;
          // 获取已同步的K线数据
          GetSyncedKLineData(code: string, startDate: string, endDate: string, page: number, pageSize: number): Promise<{ data: any[], total: number }>;
        };
      };
    };
  }
}
