import React, { useState, useEffect } from 'react';
import { parseError } from '../utils/errorHandler';
import { Download, Trash2, RefreshCw, CheckCircle, AlertCircle, Clock } from 'lucide-react';

interface SyncResult {
  stock_code: string;
  success: boolean;
  records_added: number;
  records_updated: number;
  message: string;
  error_message?: string;
}

interface SyncProgress {
  stock_code: string;
  status: 'syncing' | 'completed' | 'failed';
  current_index: number;
  total_count: number;
  progress: number;
  records_added?: number;
  message?: string;
  error_message?: string;
}

// DataSyncStats 与后端 models.DataSyncStats 对齐（snake_case）。
// 注意：后端在“数据库未初始化”降级时可能返回空结构或字段为 null，前端必须兜底防止白屏。
interface DataSyncStats {
  total_stocks?: number;
  synced_stocks?: number;
  total_records?: number;
  stock_list?: string[] | null;
  last_sync_time?: string;
}

const DataSyncPage: React.FC = () => {
  const [stockCodes, setStockCodes] = useState<string>('600519,000858');
  const [startDate, setStartDate] = useState<string>('2023-01-01');
  const [endDate, setEndDate] = useState<string>('2024-12-31');
  const [syncStats, setSyncStats] = useState<DataSyncStats | null>(null);
  const [syncResults, setSyncResults] = useState<SyncResult[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedStock, setSelectedStock] = useState<string>('');
  const [syncProgress, setSyncProgress] = useState<SyncProgress | null>(null);
  const [syncLog, setSyncLog] = useState<string[]>([]);

  useEffect(() => {
    loadSyncStats();

    // 监听同步进度事件
    const EventsOn = (window as any).EventsOn;
    if (!EventsOn) return;

    const unsubscribeBatchStart = EventsOn('dataSyncBatchStart', (data: any) => {
      setSyncLog([`开始批量同步，总共 ${data.total_count} 个股票`]);
      setSyncProgress(null);
    });

    const unsubscribeProgress = EventsOn('dataSyncProgress', (data: any) => {
      setSyncProgress(data);
      const logMessage = `[${data.current_index}/${data.total_count}] ${data.stock_code}: ${data.status} (${data.progress.toFixed(1)}%)`;
      setSyncLog((prev) => [...prev, logMessage]);
    });

    const unsubscribeBatchComplete = EventsOn('dataSyncBatchComplete', () => {
      setSyncLog((prev) => [...prev, '批量同步完成']);
      setSyncProgress(null);
      loadSyncStats();
    });

    return () => {
      if (unsubscribeBatchStart) unsubscribeBatchStart();
      if (unsubscribeProgress) unsubscribeProgress();
      if (unsubscribeBatchComplete) unsubscribeBatchComplete();
    };
  }, []);

  const loadSyncStats = async () => {
    try {
      // @ts-ignore
      const stats = await window.go.main.App.GetDataSyncStats();

      // 兜底：保证 syncStats 至少是一个可渲染的结构，避免 stock_list 为 null 时 join 崩溃。
      const safeStats: DataSyncStats = {
        total_stocks: stats?.total_stocks ?? 0,
        synced_stocks: stats?.synced_stocks ?? 0,
        total_records: stats?.total_records ?? 0,
        stock_list: Array.isArray(stats?.stock_list) ? stats.stock_list : [],
        last_sync_time: stats?.last_sync_time ?? '',
      };

      setSyncStats(safeStats);
    } catch (err) {
      console.error('加载同步统计失败:', err);
      // 如果统计加载失败，也不要让页面崩溃；显示空统计即可。
      setSyncStats({
        total_stocks: 0,
        synced_stocks: 0,
        total_records: 0,
        stock_list: [],
        last_sync_time: '',
      });
    }
  };

  const handleSingleSync = async () => {
    if (!selectedStock.trim()) {
      setError('请输入股票代码');
      return;
    }

    setLoading(true);
    setError(null);
    setSyncResults([]);
    setSyncLog([`开始同步 ${selectedStock}`]);

    try {
      // @ts-ignore
      const result = await window.go.main.App.SyncStockData(selectedStock.trim(), startDate, endDate);
      setSyncResults([result]);
      setSyncLog((prev) => [...prev, `${selectedStock} 同步完成: ${result.message}`]);
      await loadSyncStats();
    } catch (err) {
      const errorResult = parseError(err);
      setError(errorResult.message);
      setSyncLog((prev) => [...prev, `错误: ${errorResult.message}`]);
    } finally {
      setLoading(false);
    }
  };

  const handleBatchSync = async () => {
    const codes = stockCodes
      .split(',')
      .map((code) => code.trim())
      .filter((code) => code.length > 0);

    if (codes.length === 0) {
      setError('请输入至少一个股票代码');
      return;
    }

    setLoading(true);
    setError(null);
    setSyncResults([]);

    try {
      // @ts-ignore
      await window.go.main.App.BatchSyncStockData(codes, startDate, endDate);
      // 批量同步完成后，逐个获取结果
      const results: SyncResult[] = [];
      for (const code of codes) {
        try {
          // @ts-ignore
          const result = await window.go.main.App.SyncStockData(code, startDate, endDate);
          results.push(result);
        } catch (err) {
          results.push({
            stock_code: code,
            success: false,
            records_added: 0,
            records_updated: 0,
            message: '同步失败',
            error_message: String(err),
          });
        }
      }
      setSyncResults(results);
      await loadSyncStats();
    } catch (err) {
      const errorResult = parseError(err);
      setError(errorResult.message);
    } finally {
      setLoading(false);
    }
  };

  const handleClearCache = async (code: string) => {
    if (!window.confirm(`确定要清除 ${code} 的本地缓存数据吗？`)) {
      return;
    }

    try {
      // @ts-ignore
      await window.go.main.App.ClearStockCache(code);
      alert('缓存已清除');
      await loadSyncStats();
    } catch (err) {
      const errorResult = parseError(err);
      alert(`清除失败: ${errorResult.message}`);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* 页面标题 */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2">数据同步中心</h1>
          <p className="text-gray-400">将股票历史数据同步到本地 SQLite 数据库，加速回测和分析</p>
        </div>

        {/* 同步统计信息 */}
        {syncStats && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-gray-800 p-4 rounded-lg shadow-lg">
              <p className="text-sm text-gray-400">已同步股票数</p>
              <p className="text-2xl font-bold text-blue-400">{syncStats.synced_stocks ?? 0}</p>
            </div>
            <div className="bg-gray-800 p-4 rounded-lg shadow-lg">
              <p className="text-sm text-gray-400">总数据记录数</p>
              <p className="text-2xl font-bold text-green-400">{syncStats.total_records ?? 0}</p>
            </div>
            <div className="bg-gray-800 p-4 rounded-lg shadow-lg">
              <p className="text-sm text-gray-400">最后同步时间</p>
              <p className="text-sm font-mono text-gray-300">{syncStats.last_sync_time || '-'}</p>
            </div>
            <div className="bg-gray-800 p-4 rounded-lg shadow-lg">
              <p className="text-sm text-gray-400">已同步股票列表</p>
              <p className="text-sm font-mono text-gray-300">{(syncStats.stock_list ?? []).join(', ')}</p>
            </div>
          </div>
        )}

        {/* 实时进度显示 */}
        {syncProgress && (
          <div className="bg-gray-800 rounded-lg shadow-lg p-6 mb-6">
            <h2 className="text-xl font-bold mb-4">同步进度</h2>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between mb-2">
                  <span className="text-sm font-medium text-gray-300">
                    {syncProgress.stock_code}
                  </span>
                  <span className="text-sm font-medium text-gray-300">
                    {syncProgress.current_index}/{syncProgress.total_count}
                  </span>
                </div>
                <div className="w-full bg-gray-700 rounded-full h-2">
                  <div
                    className="bg-blue-500 h-2 rounded-full transition-all duration-300"
                    style={{ width: `${syncProgress.progress}%` }}
                  ></div>
                </div>
                <p className="text-xs text-gray-400 mt-2">
                  {syncProgress.progress.toFixed(1)}% - {syncProgress.status}
                </p>
              </div>
            </div>
          </div>
        )}

        {/* 单个同步面板 */}
        <div className="bg-gray-800 rounded-lg shadow-lg p-6 mb-6">
          <h2 className="text-xl font-bold mb-4">单个股票同步</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
            <div>
              <label htmlFor="singleStock" className="block text-sm font-medium text-gray-300 mb-2">
                股票代码
              </label>
              <input
                type="text"
                id="singleStock"
                value={selectedStock}
                onChange={(e) => setSelectedStock(e.target.value)}
                placeholder="如: 600519"
                className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
              />
            </div>
            <div>
              <label htmlFor="startDate" className="block text-sm font-medium text-gray-300 mb-2">
                开始日期
              </label>
              <input
                type="date"
                id="startDate"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
              />
            </div>
            <div>
              <label htmlFor="endDate" className="block text-sm font-medium text-gray-300 mb-2">
                结束日期
              </label>
              <input
                type="date"
                id="endDate"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
                className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
              />
            </div>
            <div className="flex items-end">
              <button
                onClick={handleSingleSync}
                className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                disabled={loading}
              >
                <Download className="w-4 h-4" />
                {loading ? '同步中...' : '开始同步'}
              </button>
            </div>
          </div>
        </div>

        {/* 批量同步面板 */}
        <div className="bg-gray-800 rounded-lg shadow-lg p-6 mb-6">
          <h2 className="text-xl font-bold mb-4">批量同步</h2>
          <div className="mb-4">
            <label htmlFor="stockCodes" className="block text-sm font-medium text-gray-300 mb-2">
              股票代码列表（用逗号分隔，如: 600519,000858,600000）
            </label>
            <textarea
              id="stockCodes"
              value={stockCodes}
              onChange={(e) => setStockCodes(e.target.value)}
              rows={3}
              className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
              placeholder="输入多个股票代码，用逗号分隔"
            />
          </div>
          <button
            onClick={handleBatchSync}
            className="w-full px-4 py-2 bg-green-600 hover:bg-green-700 text-white font-medium rounded-md transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
            disabled={loading}
          >
            <RefreshCw className="w-4 h-4" />
            {loading ? '同步中...' : '批量同步'}
          </button>
        </div>

        {/* 同步日志 */}
        {syncLog.length > 0 && (
          <div className="bg-gray-800 rounded-lg shadow-lg p-6 mb-6">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
              <Clock className="w-5 h-5" />
              同步日志
            </h2>
            <div className="bg-gray-900 rounded-md p-4 max-h-64 overflow-y-auto font-mono text-sm text-gray-300 space-y-1">
              {syncLog.map((log, index) => (
                <div key={index} className="text-gray-400">
                  {log}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* 错误提示 */}
        {error && (
          <div className="mb-6 bg-red-900/30 border border-red-600 rounded-lg p-4 flex items-start gap-3">
            <AlertCircle className="w-5 h-5 text-red-400 flex-shrink-0 mt-0.5" />
            <div>
              <p className="text-red-400 font-medium">错误</p>
              <p className="text-red-300 text-sm">{error}</p>
            </div>
          </div>
        )}

        {/* 同步结果 */}
        {syncResults.length > 0 && (
          <div className="bg-gray-800 rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-bold mb-4">同步结果</h2>
            <div className="space-y-3">
              {syncResults.map((result) => (
                <div
                  key={result.stock_code}
                  className={`flex items-start gap-4 p-4 rounded-lg ${
                    result.success ? 'bg-green-900/20 border border-green-600' : 'bg-red-900/20 border border-red-600'
                  }`}
                >
                  {result.success ? (
                    <CheckCircle className="w-5 h-5 text-green-400 flex-shrink-0 mt-0.5" />
                  ) : (
                    <AlertCircle className="w-5 h-5 text-red-400 flex-shrink-0 mt-0.5" />
                  )}
                  <div className="flex-1">
                    <p className={`font-medium ${result.success ? 'text-green-400' : 'text-red-400'}`}>
                      {result.stock_code}
                    </p>
                    <p className="text-sm text-gray-300">{result.message}</p>
                    {result.success && (
                      <p className="text-xs text-gray-400 mt-1">
                        新增: {result.records_added} | 更新: {result.records_updated}
                      </p>
                    )}
                    {result.error_message && (
                      <p className="text-xs text-red-300 mt-1">{result.error_message}</p>
                    )}
                  </div>
                  {result.success && (
                    <button
                      onClick={() => handleClearCache(result.stock_code)}
                      className="px-3 py-1 bg-gray-700 hover:bg-gray-600 text-gray-300 text-sm rounded transition-colors flex items-center gap-1 flex-shrink-0"
                    >
                      <Trash2 className="w-4 h-4" />
                      清除
                    </button>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default DataSyncPage;
