import React, { useState, useEffect } from 'react';
import { parseError } from '../utils/errorHandler';
import { History, Trash2, RefreshCw, CheckCircle, AlertCircle, Clock, Filter, Search, ChevronLeft, ChevronRight } from 'lucide-react';
import { useWailsAPI } from '../hooks/useWailsAPI';

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

const SyncHistoryPage: React.FC = () => {
  const { getAllSyncHistory, getSyncHistoryCount, clearAllSyncHistory } = useWailsAPI();
  const [histories, setHistories] = useState<SyncHistoryItem[]>([]);
  const [filteredHistories, setFilteredHistories] = useState<SyncHistoryItem[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [totalCount, setTotalCount] = useState<number>(0);
  const pageSize = 20;

  useEffect(() => {
    loadSyncHistories();
  }, [currentPage]);

  useEffect(() => {
    // 过滤历史记录
    let filtered = histories;

    if (searchTerm) {
      filtered = filtered.filter(
        (item) =>
          item.stock_code.toLowerCase().includes(searchTerm.toLowerCase()) ||
          item.stock_name.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    if (statusFilter !== 'all') {
      filtered = filtered.filter((item) => item.status === statusFilter);
    }

    setFilteredHistories(filtered);
  }, [searchTerm, statusFilter, histories]);

  const loadSyncHistories = async () => {
    setLoading(true);
    setError(null);
    try {
      const [historyList, count] = await Promise.all([
        getAllSyncHistory(pageSize, (currentPage - 1) * pageSize),
        getSyncHistoryCount(),
      ]);
      setHistories(historyList);
      setTotalCount(count);
    } catch (err) {
      const errorResult = parseError(err);
      setError(errorResult.message);
    } finally {
      setLoading(false);
    }
  };

  const handleClearAll = async () => {
    if (!window.confirm('确定要清除所有同步历史记录吗？此操作不可恢复。')) {
      return;
    }

    try {
      await clearAllSyncHistory();
      alert('同步历史记录已清除');
      await loadSyncHistories();
    } catch (err) {
      const errorResult = parseError(err);
      alert(`清除失败: ${errorResult.message}`);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'success':
        return <CheckCircle className="w-5 h-5 text-green-500" />;
      case 'failed':
        return <AlertCircle className="w-5 h-5 text-red-500" />;
      default:
        return <Clock className="w-5 h-5 text-gray-500" />;
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'success':
        return '成功';
      case 'failed':
        return '失败';
      default:
        return status;
    }
  };

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'success':
        return 'bg-green-100 text-green-700 border-green-200';
      case 'failed':
        return 'bg-red-100 text-red-700 border-red-200';
      default:
        return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const totalPages = Math.ceil(totalCount / pageSize);

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* 页面标题和操作按钮 */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-4xl font-bold mb-2">同步历史记录</h1>
            <p className="text-gray-400">查看和管理所有的数据同步记录</p>
          </div>
          <div className="flex items-center space-x-4">
            <button
              onClick={loadSyncHistories}
              className="flex items-center space-x-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
              disabled={loading}
            >
              <RefreshCw className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
              <span>刷新</span>
            </button>
            <button
              onClick={handleClearAll}
              className="flex items-center space-x-2 px-4 py-2 bg-red-600 hover:bg-red-700 rounded-lg transition-colors"
            >
              <Trash2 className="w-5 h-5" />
              <span>清除全部</span>
            </button>
          </div>
        </div>

        {/* 统计信息 */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-gray-800 rounded-lg p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-400 text-sm">总记录数</p>
                <p className="text-2xl font-bold">{totalCount}</p>
              </div>
              <History className="w-8 h-8 text-blue-500" />
            </div>
          </div>
          <div className="bg-gray-800 rounded-lg p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-400 text-sm">成功次数</p>
                <p className="text-2xl font-bold text-green-500">
                  {histories.filter((h) => h.status === 'success').length}
                </p>
              </div>
              <CheckCircle className="w-8 h-8 text-green-500" />
            </div>
          </div>
          <div className="bg-gray-800 rounded-lg p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-400 text-sm">失败次数</p>
                <p className="text-2xl font-bold text-red-500">
                  {histories.filter((h) => h.status === 'failed').length}
                </p>
              </div>
              <AlertCircle className="w-8 h-8 text-red-500" />
            </div>
          </div>
          <div className="bg-gray-800 rounded-lg p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-400 text-sm">显示记录</p>
                <p className="text-2xl font-bold">{filteredHistories.length}</p>
              </div>
              <Filter className="w-8 h-8 text-purple-500" />
            </div>
          </div>
        </div>

        {/* 搜索和过滤 */}
        <div className="bg-gray-800 rounded-lg p-4 mb-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="搜索股票代码或名称..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:border-blue-500"
              />
            </div>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:border-blue-500"
            >
              <option value="all">全部状态</option>
              <option value="success">成功</option>
              <option value="failed">失败</option>
            </select>
          </div>
        </div>

        {/* 错误提示 */}
        {error && (
          <div className="bg-red-900/30 border border-red-700 rounded-lg p-4 mb-6 flex items-center">
            <AlertCircle className="w-5 h-5 text-red-500 mr-3" />
            <p className="text-red-300">{error}</p>
          </div>
        )}

        {/* 历史记录表格 */}
        <div className="bg-gray-800 rounded-lg overflow-hidden">
          {loading ? (
            <div className="flex items-center justify-center py-20">
              <RefreshCw className="w-8 h-8 text-blue-500 animate-spin" />
              <span className="ml-3 text-gray-400">加载中...</span>
            </div>
          ) : filteredHistories.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-20">
              <History className="w-16 h-16 text-gray-600 mb-4" />
              <p className="text-gray-400">暂无同步历史记录</p>
            </div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-700">
                      <th className="text-left py-4 px-6 text-sm font-medium text-gray-400">股票</th>
                      <th className="text-left py-4 px-6 text-sm font-medium text-gray-400">同步类型</th>
                      <th className="text-left py-4 px-6 text-sm font-medium text-gray-400">日期范围</th>
                      <th className="text-left py-4 px-6 text-sm font-medium text-gray-400">状态</th>
                      <th className="text-right py-4 px-6 text-sm font-medium text-gray-400">新增记录</th>
                      <th className="text-right py-4 px-6 text-sm font-medium text-gray-400">更新记录</th>
                      <th className="text-right py-4 px-6 text-sm font-medium text-gray-400">耗时</th>
                      <th className="text-left py-4 px-6 text-sm font-medium text-gray-400">同步时间</th>
                    </tr>
                  </thead>
                  <tbody>
                    {filteredHistories.map((item) => (
                      <tr key={item.id} className="border-b border-gray-700 hover:bg-gray-700/50 transition-colors">
                        <td className="py-4 px-6">
                          <div className="flex items-center space-x-3">
                            <div className="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center text-sm font-bold">
                              {item.stock_name ? item.stock_name.charAt(0) : item.stock_code}
                            </div>
                            <div>
                              <div className="font-medium">{item.stock_name || '未知'}</div>
                              <div className="text-sm text-gray-400">{item.stock_code}</div>
                            </div>
                          </div>
                        </td>
                        <td className="py-4 px-6">
                          <span className="text-sm">{item.sync_type === 'single' ? '单个同步' : '批量同步'}</span>
                        </td>
                        <td className="py-4 px-6">
                          <div className="text-sm text-gray-300">
                            {item.start_date} 至 {item.end_date}
                          </div>
                        </td>
                        <td className="py-4 px-6">
                          <div className={`inline-flex items-center space-x-1 px-2 py-1 rounded-full text-xs font-medium border ${getStatusBadgeClass(item.status)}`}>
                            {getStatusIcon(item.status)}
                            <span>{getStatusText(item.status)}</span>
                          </div>
                        </td>
                        <td className="py-4 px-6 text-right">
                          <span className="text-sm font-medium">{item.records_added}</span>
                        </td>
                        <td className="py-4 px-6 text-right">
                          <span className="text-sm font-medium">{item.records_updated}</span>
                        </td>
                        <td className="py-4 px-6 text-right">
                          <span className="text-sm text-gray-300">{item.duration}s</span>
                        </td>
                        <td className="py-4 px-6">
                          <div className="text-sm text-gray-300">{formatDate(item.created_at)}</div>
                          {item.error_msg && (
                            <div className="text-xs text-red-400 mt-1" title={item.error_msg}>
                              {item.error_msg.length > 30 ? item.error_msg.substring(0, 30) + '...' : item.error_msg}
                            </div>
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* 分页 */}
              {totalPages > 1 && (
                <div className="flex items-center justify-between px-6 py-4 border-t border-gray-700">
                  <div className="text-sm text-gray-400">
                    第 {currentPage} 页，共 {totalPages} 页
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
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default SyncHistoryPage;
