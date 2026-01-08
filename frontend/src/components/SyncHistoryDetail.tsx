import React, { useState, useEffect } from 'react';
import { X, RefreshCw, ChevronLeft, ChevronRight, Filter, Search, Calendar } from 'lucide-react';
import KLineChart from './KLineChart';
import { useWailsAPI } from '../hooks/useWailsAPI';
import type { KLineData } from '../types';

interface SyncHistoryDetailProps {
  isOpen: boolean;
  onClose: () => void;
  stockCode: string;
  stockName: string;
  startDate: string;
  endDate: string;
  onResync?: (code: string) => void;
}

const SyncHistoryDetail: React.FC<SyncHistoryDetailProps> = ({
  isOpen,
  onClose,
  stockCode,
  stockName,
  startDate,
  endDate,
  onResync,
}) => {
  const { getSyncedKLineData, getStockData, SyncStockData } = useWailsAPI();
  const [klineData, setKlineData] = useState<KLineData[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [syncing, setSyncing] = useState<boolean>(false);

  // 表格数据
  const [tableData, setTableData] = useState<any[]>([]);
  const [filteredData, setFilteredData] = useState<any[]>([]);

  // 分页
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize] = useState<number>(20);
  const [totalCount, setTotalCount] = useState<number>(0);

  // 筛选
  const [filterStartDate, setFilterStartDate] = useState<string>('');
  const [filterEndDate, setFilterEndDate] = useState<string>('');

  useEffect(() => {
    if (isOpen) {
      // 初始化日期筛选
      setFilterStartDate(startDate);
      setFilterEndDate(endDate);
      setCurrentPage(1);
      loadKLineData(startDate, endDate);
    }
  }, [isOpen, stockCode, startDate, endDate]);

  useEffect(() => {
    applyFilter();
  }, [tableData, filterStartDate, filterEndDate]);

  const loadKLineData = async (start: string, end: string) => {
    setLoading(true);
    setError(null);
    try {
      const result = await getSyncedKLineData(stockCode, start, end, 1, 5000);

      // 验证返回的数据格式
      if (!result || !Array.isArray(result.data)) {
        throw new Error('返回数据格式错误：data 字段应该是数组')
      }

      const safeData = result.data || []

      setTableData(safeData);
      setTotalCount(result.total || 0);

      // 转换为K线图数据格式
      const klineChartData: KLineData[] = safeData.map((item: any) => ({
        time: item.date,
        open: item.open,
        high: item.high,
        low: item.low,
        close: item.close,
        volume: item.volume,
      }));

      setKlineData(klineChartData);
    } catch (err: any) {
      console.error('加载K线数据失败:', err)
      setError(err.message || '加载K线数据失败');
      setTableData([]);
      setKlineData([]);
      setTotalCount(0);
    } finally {
      setLoading(false);
    }
  };

  const applyFilter = () => {
    // 确保 tableData 是数组
    const safeTableData = Array.isArray(tableData) ? tableData : [];

    let filtered = [...safeTableData];

    if (filterStartDate) {
      filtered = filtered.filter((item) => item.date >= filterStartDate);
    }
    if (filterEndDate) {
      filtered = filtered.filter((item) => item.date <= filterEndDate);
    }

    setFilteredData(filtered);
    setCurrentPage(1);
  };

  const handleFilter = () => {
    applyFilter();
    // 更新K线图数据
    if (filterStartDate || filterEndDate) {
      loadKLineData(filterStartDate, filterEndDate);
    }
  };

  const handleResetFilter = () => {
    setFilterStartDate(startDate);
    setFilterEndDate(endDate);
    applyFilter();
    loadKLineData(startDate, endDate);
  };

  const handleResync = async () => {
    if (!window.confirm(`确定要重新同步 ${stockName}(${stockCode}) 的数据吗？`)) {
      return;
    }

    setSyncing(true);
    try {
      await SyncStockData(stockCode, startDate, endDate);
      alert('重新同步成功');
      loadKLineData(startDate, endDate);
      if (onResync) {
        onResync(stockCode);
      }
    } catch (err: any) {
      alert(`重新同步失败: ${err.message}`);
    } finally {
      setSyncing(false);
    }
  };

  const totalPages = Math.ceil(filteredData.length / pageSize);
  const startIndex = (currentPage - 1) * pageSize;
  const endIndex = startIndex + pageSize;
  const currentPageData = filteredData.slice(startIndex, endIndex);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-gray-800 rounded-xl w-full max-w-6xl max-h-[90vh] overflow-hidden flex flex-col">
        {/* 头部 */}
        <div className="flex items-center justify-between p-6 border-b border-gray-700">
          <div className="flex items-center space-x-4">
            <div className="w-12 h-12 bg-blue-600 rounded-lg flex items-center justify-center text-xl font-bold">
              {stockName.charAt(0)}
            </div>
            <div>
              <h2 className="text-2xl font-bold">{stockName}</h2>
              <p className="text-gray-400">{stockCode}</p>
            </div>
          </div>
          <div className="flex items-center space-x-4">
            <button
              onClick={handleResync}
              disabled={syncing}
              className="flex items-center space-x-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors disabled:opacity-50"
            >
              <RefreshCw className={`w-5 h-5 ${syncing ? 'animate-spin' : ''}`} />
              <span>重新同步</span>
            </button>
            <button
              onClick={onClose}
              className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
            >
              <X className="w-6 h-6" />
            </button>
          </div>
        </div>

        {/* 内容区 */}
        <div className="flex-1 overflow-y-auto p-6">
          {loading ? (
            <div className="flex items-center justify-center py-20">
              <RefreshCw className="w-8 h-8 text-blue-500 animate-spin" />
              <span className="ml-3 text-gray-400">加载中...</span>
            </div>
          ) : error ? (
            <div className="bg-red-900/30 border border-red-700 rounded-lg p-4">
              <p className="text-red-300">{error}</p>
            </div>
          ) : (
            <>
              {/* K线图 */}
              <div className="bg-gray-900 rounded-lg p-4 mb-6">
                <h3 className="text-lg font-semibold mb-4">K线图</h3>
                {klineData.length > 0 ? (
                  <KLineChart data={klineData} height={400} />
                ) : (
                  <div className="flex items-center justify-center h-64 text-gray-400">
                    暂无数据
                  </div>
                )}
              </div>

              {/* 数据表格 */}
              <div className="bg-gray-900 rounded-lg p-4">
                {/* 筛选栏 */}
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-semibold">数据明细</h3>
                  <div className="flex items-center space-x-4">
                    <div className="flex items-center space-x-2">
                      <Calendar className="w-4 h-4 text-gray-400" />
                      <input
                        type="date"
                        value={filterStartDate}
                        onChange={(e) => setFilterStartDate(e.target.value)}
                        className="px-3 py-1.5 bg-gray-700 border border-gray-600 rounded-lg text-sm focus:outline-none focus:border-blue-500"
                      />
                      <span className="text-gray-400">至</span>
                      <input
                        type="date"
                        value={filterEndDate}
                        onChange={(e) => setFilterEndDate(e.target.value)}
                        className="px-3 py-1.5 bg-gray-700 border border-gray-600 rounded-lg text-sm focus:outline-none focus:border-blue-500"
                      />
                    </div>
                    <button
                      onClick={handleFilter}
                      className="flex items-center space-x-2 px-3 py-1.5 bg-blue-600 hover:bg-blue-700 rounded-lg text-sm transition-colors"
                    >
                      <Filter className="w-4 h-4" />
                      <span>筛选</span>
                    </button>
                    <button
                      onClick={handleResetFilter}
                      className="px-3 py-1.5 bg-gray-700 hover:bg-gray-600 rounded-lg text-sm transition-colors"
                    >
                      重置
                    </button>
                  </div>
                </div>

                {/* 表格 */}
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b border-gray-700">
                        <th className="text-left py-3 px-4 text-sm font-medium text-gray-400">日期</th>
                        <th className="text-right py-3 px-4 text-sm font-medium text-gray-400">开盘价</th>
                        <th className="text-right py-3 px-4 text-sm font-medium text-gray-400">最高价</th>
                        <th className="text-right py-3 px-4 text-sm font-medium text-gray-400">最低价</th>
                        <th className="text-right py-3 px-4 text-sm font-medium text-gray-400">收盘价</th>
                        <th className="text-right py-3 px-4 text-sm font-medium text-gray-400">成交量</th>
                      </tr>
                    </thead>
                    <tbody>
                      {currentPageData.map((item, index) => (
                        <tr key={index} className="border-b border-gray-700 hover:bg-gray-800 transition-colors">
                          <td className="py-3 px-4 text-sm">{item.date}</td>
                          <td className="py-3 px-4 text-sm text-right">{item.open.toFixed(2)}</td>
                          <td className="py-3 px-4 text-sm text-right text-red-400">{item.high.toFixed(2)}</td>
                          <td className="py-3 px-4 text-sm text-right text-green-400">{item.low.toFixed(2)}</td>
                          <td className="py-3 px-4 text-sm text-right">{item.close.toFixed(2)}</td>
                          <td className="py-3 px-4 text-sm text-right">{item.volume.toLocaleString()}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                {/* 分页 */}
                {totalPages > 1 && (
                  <div className="flex items-center justify-between mt-4 pt-4 border-t border-gray-700">
                    <div className="text-sm text-gray-400">
                      显示 {startIndex + 1} 到 {Math.min(endIndex, filteredData.length)} 条，共 {filteredData.length} 条
                    </div>
                    <div className="flex items-center space-x-2">
                      <button
                        onClick={() => setCurrentPage((prev) => Math.max(1, prev - 1))}
                        disabled={currentPage === 1}
                        className="p-2 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition-colors"
                      >
                        <ChevronLeft className="w-5 h-5" />
                      </button>
                      <span className="px-4 py-2 bg-blue-600 rounded-lg text-sm font-medium">
                        {currentPage}
                      </span>
                      <button
                        onClick={() => setCurrentPage((prev) => Math.min(totalPages, prev + 1))}
                        disabled={currentPage === totalPages}
                        className="p-2 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition-colors"
                      >
                        <ChevronRight className="w-5 h-5" />
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default SyncHistoryDetail;
