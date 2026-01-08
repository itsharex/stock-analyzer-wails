import React, { useState, useEffect } from 'react';
import { parseError } from '../utils/errorHandler';
import { RefreshCw, CheckCircle, AlertCircle, Clock, Zap, Database } from 'lucide-react';

interface KLineSyncResult {
  success: boolean;
  total_count: number;
  success_count: number;
  failed_count: number;
  total_records: number;
  duration: number;
  message: string;
}

interface KLineSyncProgress {
  is_running: boolean;
  current_index: number;
  total_count: number;
  current_code: string;
  current_name: string;
  success_count: number;
  failed_count: number;
  total_records: number;
  records_per_sec: number;
  start_time: string;
  elapsed_seconds: number;
  estimated_seconds: number;
}

interface KLineSyncHistory {
  id: number;
  stockCode: string;
  stockName: string;
  syncType: string;
  startDate: string;
  endDate: string;
  status: string;
  recordsAdded: number;
  recordsUpdated: number;
  duration: number;
  errorMsg: string;
  createdAt: string;
}

const KLineSyncPage: React.FC = () => {
  const [days, setDays] = useState<number>(200);
  const [syncProgress, setSyncProgress] = useState<KLineSyncProgress | null>(null);
  const [syncHistory, setSyncHistory] = useState<KLineSyncHistory[]>([]);
  const [syncLog, setSyncLog] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [syncResult, setSyncResult] = useState<KLineSyncResult | null>(null);

  useEffect(() => {
    loadSyncHistory();

    // 监听K线同步进度事件
    const EventsOn = (window as any).EventsOn;
    if (!EventsOn) return;

    const unsubscribeProgress = EventsOn('klineSyncProgress', (data: KLineSyncProgress) => {
      setSyncProgress(data);
      
      const progressPercent = ((data.current_index / data.total_count) * 100).toFixed(1);
      const logMessage = `[${data.current_index}/${data.total_count}] ${data.current_code} ${data.current_name}: 成功 ${data.success_count} | 失败 ${data.failed_count} | 进度 ${progressPercent}% | 速率 ${data.records_per_sec.toFixed(0)} 条/秒`;
      setSyncLog((prev) => {
        const newLog = [...prev, logMessage];
        // 只保留最近100条日志
        if (newLog.length > 100) {
          return newLog.slice(-100);
        }
        return newLog;
      });

      // 当同步完成时
      if (!data.is_running) {
        setLoading(false);
      }
    });

    return () => {
      if (unsubscribeProgress) unsubscribeProgress();
    };
  }, []);

  const loadSyncHistory = async () => {
    try {
      // @ts-ignore
      const history = await window.go.main.App.GetKLineSyncHistory(50);
      if (Array.isArray(history)) {
        setSyncHistory(history);
      }
    } catch (err) {
      console.error('加载同步历史失败:', err);
    }
  };

  const handleStartSync = async () => {
    if (days <= 0 || days > 2000) {
      setError('日期范围无效（1-2000天）');
      return;
    }

    setLoading(true);
    setError(null);
    setSyncResult(null);
    setSyncProgress(null);
    setSyncLog([`开始K线数据同步，同步最近 ${days} 天的数据`]);

    try {
      // @ts-ignore
      const result = await window.go.main.App.StartKLineSync(days);
      setSyncResult(result as KLineSyncResult);
      setSyncLog((prev) => [...prev, '同步任务已启动']);
      await loadSyncHistory();
    } catch (err) {
      const errorResult = parseError(err);
      setError(errorResult.message);
      setSyncLog((prev) => [...prev, `错误: ${errorResult.message}`]);
      setLoading(false);
    }
  };

  const formatDuration = (seconds: number): string => {
    if (seconds < 60) {
      return `${seconds}秒`;
    } else if (seconds < 3600) {
      const mins = Math.floor(seconds / 60);
      const secs = seconds % 60;
      return `${mins}分${secs}秒`;
    } else {
      const hours = Math.floor(seconds / 3600);
      const mins = Math.floor((seconds % 3600) / 60);
      return `${hours}小时${mins}分`;
    }
  };

  const progressPercent = syncProgress
    ? ((syncProgress.current_index / syncProgress.total_count) * 100).toFixed(1)
    : '0';

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* 页面标题 */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2 flex items-center gap-3">
            <Database className="w-10 h-10 text-blue-400" />
            K线数据同步
          </h1>
          <p className="text-gray-400">
            批量同步所有活跃股票的K线数据到本地数据库，顺序执行避免数据库锁定问题
          </p>
        </div>

        {/* 功能说明 */}
        <div className="bg-gray-800 rounded-lg shadow-lg p-6 mb-6 border border-gray-700">
          <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
            <Zap className="w-5 h-5 text-yellow-400" />
            功能说明
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-300">
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-400 mt-1 flex-shrink-0" />
              <span>自动同步数据库中 is_active=1 的所有股票</span>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-400 mt-1 flex-shrink-0" />
              <span>顺序执行每只股票同步，避免SQLite并发锁库问题</span>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-400 mt-1 flex-shrink-0" />
              <span>随机延迟（200-500ms）模拟真人行为，防止IP被封</span>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-400 mt-1 flex-shrink-0" />
              <span>实时显示同步进度、成功/失败数量、速率等信息</span>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-400 mt-1 flex-shrink-0" />
              <span>K线数据包含：开盘价、收盘价、最高价、最低价、成交量、成交额</span>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-400 mt-1 flex-shrink-0" />
              <span>每个股票使用独立表存储，数据存在则自动更新</span>
            </div>
          </div>
        </div>

        {/* 同步配置和启动 */}
        <div className="bg-gray-800 rounded-lg shadow-lg p-6 mb-6">
          <h2 className="text-xl font-bold mb-4">同步配置</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label htmlFor="days" className="block text-sm font-medium text-gray-300 mb-2">
                同步天数（最近几天）
              </label>
              <input
                type="number"
                id="days"
                value={days}
                onChange={(e) => setDays(parseInt(e.target.value) || 0)}
                min={1}
                max={2000}
                className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
              />
              <p className="text-xs text-gray-400 mt-1">建议：200天（约10个月数据）</p>
            </div>
            <div className="flex items-end">
              <button
                onClick={handleStartSync}
                disabled={loading || syncProgress?.is_running}
                className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
              >
                <RefreshCw className={`w-4 h-4 ${loading || syncProgress?.is_running ? 'animate-spin' : ''}`} />
                {syncProgress?.is_running ? '同步中...' : '开始同步'}
              </button>
            </div>
          </div>
        </div>

        {/* 实时进度显示 */}
        {syncProgress && (
          <div className="bg-gray-800 rounded-lg shadow-lg p-6 mb-6 border border-blue-600">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
              <Clock className="w-5 h-5 text-blue-400" />
              同步进度
            </h2>
            
            {/* 进度条 */}
            <div className="mb-6">
              <div className="flex justify-between mb-2">
                <span className="text-sm font-medium text-gray-300">
                  {syncProgress.current_code} {syncProgress.current_name}
                </span>
                <span className="text-sm font-medium text-gray-300">
                  {syncProgress.current_index} / {syncProgress.total_count} ({progressPercent}%)
                </span>
              </div>
              <div className="w-full bg-gray-700 rounded-full h-3">
                <div
                  className="bg-gradient-to-r from-blue-500 to-blue-600 h-3 rounded-full transition-all duration-300"
                  style={{ width: `${progressPercent}%` }}
                ></div>
              </div>
            </div>

            {/* 详细统计 */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
              <div className="bg-gray-700 p-3 rounded-lg">
                <p className="text-xs text-gray-400 mb-1">成功</p>
                <p className="text-lg font-bold text-green-400">{syncProgress.success_count}</p>
              </div>
              <div className="bg-gray-700 p-3 rounded-lg">
                <p className="text-xs text-gray-400 mb-1">失败</p>
                <p className="text-lg font-bold text-red-400">{syncProgress.failed_count}</p>
              </div>
              <div className="bg-gray-700 p-3 rounded-lg">
                <p className="text-xs text-gray-400 mb-1">总记录数</p>
                <p className="text-lg font-bold text-blue-400">{syncProgress.total_records}</p>
              </div>
              <div className="bg-gray-700 p-3 rounded-lg">
                <p className="text-xs text-gray-400 mb-1">同步速率</p>
                <p className="text-lg font-bold text-yellow-400">{syncProgress.records_per_sec.toFixed(0)} 条/秒</p>
              </div>
            </div>

            {/* 时间信息 */}
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div className="bg-gray-700 p-3 rounded-lg">
                <p className="text-gray-400">已用时间</p>
                <p className="font-medium text-gray-200">{formatDuration(syncProgress.elapsed_seconds)}</p>
              </div>
              <div className="bg-gray-700 p-3 rounded-lg">
                <p className="text-gray-400">预计剩余</p>
                <p className="font-medium text-gray-200">{formatDuration(syncProgress.estimated_seconds)}</p>
              </div>
            </div>
          </div>
        )}

        {/* 同步结果 */}
        {syncResult && (
          <div className="bg-gray-800 rounded-lg shadow-lg p-6 mb-6 border border-green-600">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
              <CheckCircle className="w-5 h-5 text-green-400" />
              同步完成
            </h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="bg-gray-700 p-4 rounded-lg">
                <p className="text-sm text-gray-400">总股票数</p>
                <p className="text-2xl font-bold text-blue-400">{syncResult.total_count}</p>
              </div>
              <div className="bg-gray-700 p-4 rounded-lg">
                <p className="text-sm text-gray-400">成功</p>
                <p className="text-2xl font-bold text-green-400">{syncResult.success_count}</p>
              </div>
              <div className="bg-gray-700 p-4 rounded-lg">
                <p className="text-sm text-gray-400">失败</p>
                <p className="text-2xl font-bold text-red-400">{syncResult.failed_count}</p>
              </div>
              <div className="bg-gray-700 p-4 rounded-lg">
                <p className="text-sm text-gray-400">总记录数</p>
                <p className="text-2xl font-bold text-yellow-400">{syncResult.total_records}</p>
              </div>
            </div>
            <div className="mt-4 text-center text-sm text-gray-300">
              <p>{syncResult.message}</p>
              <p className="text-gray-400 mt-1">总耗时: {formatDuration(syncResult.duration)}</p>
            </div>
          </div>
        )}

        {/* 同步日志 */}
        {syncLog.length > 0 && (
          <div className="bg-gray-800 rounded-lg shadow-lg p-6 mb-6">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
              <Clock className="w-5 h-5" />
              同步日志
            </h2>
            <div className="bg-gray-900 rounded-md p-4 max-h-96 overflow-y-auto font-mono text-sm space-y-1">
              {syncLog.map((log, index) => (
                <div key={index} className="text-gray-400 break-all">
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

        {/* 同步历史 */}
        <div className="bg-gray-800 rounded-lg shadow-lg p-6">
          <h2 className="text-xl font-bold mb-4">同步历史</h2>
          {syncHistory.length === 0 ? (
            <p className="text-gray-400 text-center py-4">暂无同步历史</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-700">
                    <th className="text-left py-2 px-3 text-gray-400 font-medium">时间</th>
                    <th className="text-left py-2 px-3 text-gray-400 font-medium">股票</th>
                    <th className="text-left py-2 px-3 text-gray-400 font-medium">日期范围</th>
                    <th className="text-left py-2 px-3 text-gray-400 font-medium">状态</th>
                    <th className="text-left py-2 px-3 text-gray-400 font-medium">新增</th>
                    <th className="text-left py-2 px-3 text-gray-400 font-medium">更新</th>
                  </tr>
                </thead>
                <tbody>
                  {syncHistory.map((history) => (
                    <tr key={history.id} className="border-b border-gray-700 hover:bg-gray-700/50">
                      <td className="py-2 px-3 text-gray-300">{history.createdAt}</td>
                      <td className="py-2 px-3">
                        <div>
                          <span className="text-gray-200">{history.stockCode}</span>
                          <span className="text-gray-400 ml-2">{history.stockName}</span>
                        </div>
                      </td>
                      <td className="py-2 px-3 text-gray-400">
                        {history.startDate} ~ {history.endDate}
                      </td>
                      <td className="py-2 px-3">
                        <span
                          className={`px-2 py-1 rounded text-xs font-medium ${
                            history.status === 'success'
                              ? 'bg-green-900/50 text-green-400'
                              : 'bg-red-900/50 text-red-400'
                          }`}
                        >
                          {history.status === 'success' ? '成功' : '失败'}
                        </span>
                      </td>
                      <td className="py-2 px-3 text-gray-300">{history.recordsAdded}</td>
                      <td className="py-2 px-3 text-gray-300">{history.recordsUpdated}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default KLineSyncPage;
