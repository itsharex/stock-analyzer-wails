import React, { useState, useEffect } from 'react';
import { Search, RefreshCw, Activity } from 'lucide-react';
import { useWailsAPI } from '../hooks/useWailsAPI';
import { StockMarketData, SyncStocksResult } from '../types';

const StockListPage: React.FC = () => {
  const { getStocksList, syncAllStocks, getSyncStats } = useWailsAPI();

  const [stocks, setStocks] = useState<StockMarketData[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [page, setPage] = useState<number>(1);
  const [pageSize] = useState<number>(20);
  const [search, setSearch] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [syncing, setSyncing] = useState<boolean>(false);
  const [syncResult, setSyncResult] = useState<SyncStocksResult | null>(null);
  const [lastSyncTime, setLastSyncTime] = useState<string>('-');

  // 加载股票列表
  const loadStocks = async (currentPage: number, searchQuery: string) => {
    setLoading(true);
    try {
      const result = await getStocksList(currentPage, pageSize, searchQuery);
      if (result) {
        setStocks(result.stocks || []);
        setTotal(result.total || 0);
      }
    } catch (err) {
      console.error('加载股票列表失败:', err);
    } finally {
      setLoading(false);
    }
  };

  // 加载同步统计信息
  const loadSyncStats = async () => {
    try {
      const stats = await getSyncStats();
      if (stats && stats.lastUpdate) {
        setLastSyncTime(stats.lastUpdate);
      }
    } catch (err) {
      console.error('加载同步统计失败:', err);
    }
  };

  // 同步股票数据
  const handleSync = async () => {
    setSyncing(true);
    setSyncResult(null);
    try {
      const result = await syncAllStocks();
      setSyncResult(result);
      // 同步完成后重新加载列表
      await loadStocks(page, search);
      await loadSyncStats();
    } catch (err) {
      console.error('同步失败:', err);
      alert('同步失败，请查看控制台日志');
    } finally {
      setSyncing(false);
    }
  };

  // 处理搜索
  const handleSearch = () => {
    setPage(1);
    loadStocks(1, search);
  };

  // 处理回车搜索
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };

  // 格式化数值
  const formatNumber = (value: number | null | undefined, decimals: number = 2): string => {
    if (value == null || value === 0) return '-';
    return value.toFixed(decimals);
  };

  // 格式化百分比
  const formatPercent = (value: number | null | undefined): string => {
    if (value == null || value === 0) return '-';
    return `${value.toFixed(2)}%`;
  };

  // 格式化成交量
  const formatVolume = (value: number | null | undefined): string => {
    if (value == null || value === 0) return '-';
    if (value >= 100000000) return `${(value / 100000000).toFixed(2)}亿`;
    if (value >= 10000) return `${(value / 10000).toFixed(2)}万`;
    return value.toFixed(0);
  };

  // 获取涨跌颜色
  const getChangeColor = (value: number | null | undefined): string => {
    if (value == null || value === 0) return 'text-gray-400';
    if (value > 0) return 'text-red-400';
    if (value < 0) return 'text-green-400';
    return 'text-gray-400';
  };

  // 初始加载
  useEffect(() => {
    loadStocks(page, search);
    loadSyncStats();
  }, []);

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-100">市场股票</h1>
        <div className="flex items-center gap-4">
          <span className="text-sm text-gray-400">
            总数量: {total} | 最后同步: {lastSyncTime}
          </span>
          <button
            onClick={handleSync}
            disabled={syncing}
            className={`flex items-center gap-2 px-4 py-2 rounded-md font-medium transition-colors ${
              syncing
                ? 'bg-gray-600 text-gray-400 cursor-not-allowed'
                : 'bg-blue-600 hover:bg-blue-700 text-white'
            }`}
          >
            <RefreshCw className={`w-4 h-4 ${syncing ? 'animate-spin' : ''}`} />
            {syncing ? '同步中...' : '同步数据'}
          </button>
        </div>
      </div>

      {/* 同步结果提示 */}
      {syncResult && (
        <div className="mb-4 p-4 bg-gray-800 rounded-md border border-gray-700">
          <p className="text-sm text-gray-300">
            同步完成！总计 {syncResult.total} 只股票，新增 {syncResult.inserted} 只，更新 {syncResult.updated} 只，
            耗时 {syncResult.duration.toFixed(2)} 秒
          </p>
        </div>
      )}

      {/* 搜索框 */}
      <div className="mb-6 flex items-center gap-4">
        <div className="flex-1 relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
          <input
            type="text"
            placeholder="搜索股票代码或名称..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            onKeyDown={handleKeyDown}
            className="w-full pl-10 pr-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-gray-100 placeholder-gray-500 focus:outline-none focus:border-blue-500"
          />
        </div>
        <button
          onClick={handleSearch}
          className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md font-medium transition-colors"
        >
          搜索
        </button>
        {search && (
          <button
            onClick={() => {
              setSearch('');
              setPage(1);
              loadStocks(1, '');
            }}
            className="px-6 py-2 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-md font-medium transition-colors"
          >
            清空
          </button>
        )}
      </div>

      {/* 股票列表 */}
      {loading ? (
        <div className="flex items-center justify-center h-64 text-gray-400">
          <Activity className="w-6 h-6 mr-2 animate-spin" />
          加载中...
        </div>
      ) : (
        <>
          <div className="bg-gray-800 rounded-lg overflow-hidden">
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-700">
                <thead className="bg-gray-700">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">代码</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">名称</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">最新价</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">涨跌幅</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">涨跌额</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">成交量</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">成交额</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">振幅</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">最高价</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">最低价</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">市盈率</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">换手率</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">量比</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">市场</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">类型</th>
                  </tr>
                </thead>
                <tbody className="bg-gray-800 divide-y divide-gray-700">
                  {stocks.map((stock) => (
                    <tr key={stock.id} className="hover:bg-gray-700 transition-colors">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-400">{stock.code}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{stock.name}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-100">{formatNumber(stock.price)}</td>
                      <td className={`px-6 py-4 whitespace-nowrap text-sm font-medium ${getChangeColor(stock.changeRate)}`}>
                        {formatPercent(stock.changeRate)}
                      </td>
                      <td className={`px-6 py-4 whitespace-nowrap text-sm ${getChangeColor(stock.changeAmount)}`}>
                        {formatNumber(stock.changeAmount)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{formatVolume(stock.volume)}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{formatVolume(stock.amount)}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{formatPercent(stock.amplitude)}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{formatNumber(stock.high)}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{formatNumber(stock.low)}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{formatNumber(stock.pe)}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{formatPercent(stock.turnover)}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{formatNumber(stock.volumeRatio)}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{stock.market}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{stock.type}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* 分页 */}
          {total > pageSize && (
            <div className="mt-6 flex items-center justify-between">
              <div className="text-sm text-gray-400">
                显示 {((page - 1) * pageSize) + 1} 到 {Math.min(page * pageSize, total)} 条，共 {total} 条
              </div>
              <div className="flex items-center gap-2">
                <button
                  onClick={() => {
                    setPage(1);
                    loadStocks(1, search);
                  }}
                  disabled={page === 1}
                  className="px-3 py-1 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  首页
                </button>
                <button
                  onClick={() => {
                    setPage(page - 1);
                    loadStocks(page - 1, search);
                  }}
                  disabled={page === 1}
                  className="px-3 py-1 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  上一页
                </button>
                <span className="px-4 py-1 bg-gray-800 text-gray-300 rounded">
                  第 {page} / {Math.ceil(total / pageSize)} 页
                </span>
                <button
                  onClick={() => {
                    setPage(page + 1);
                    loadStocks(page + 1, search);
                  }}
                  disabled={page >= Math.ceil(total / pageSize)}
                  className="px-3 py-1 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  下一页
                </button>
                <button
                  onClick={() => {
                    setPage(Math.ceil(total / pageSize));
                    loadStocks(Math.ceil(total / pageSize), search);
                  }}
                  disabled={page >= Math.ceil(total / pageSize)}
                  className="px-3 py-1 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  末页
                </button>
              </div>
            </div>
          )}

          {stocks.length === 0 && !loading && (
            <div className="text-center py-12 text-gray-400">
              暂无数据
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default StockListPage;
