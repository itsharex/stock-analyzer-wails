import React, { useEffect, useState } from 'react';
import { useWailsAPI } from '../hooks/useWailsAPI';
import { SignalAnalysisResult } from '../types';

interface StatsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const StatsModal: React.FC<StatsModalProps> = ({ isOpen, onClose }) => {
  const { analyzePastSignals } = useWailsAPI();
  const [result, setResult] = useState<SignalAnalysisResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen) {
      loadStats();
    }
  }, [isOpen]);

  const loadStats = async () => {
    setLoading(true);
    setError(null);
    try {
      // 扫描过去 60 天的信号
      const res = await analyzePastSignals(60);
      setResult(res);
    } catch (err: any) {
      console.error('获取统计数据失败详情:', err);
      // 检查错误对象是否为 string
      const errorMessage = typeof err === 'string' ? err : (err.message || '获取统计数据失败');
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="bg-[#161b22] w-[600px] rounded-lg shadow-2xl border border-gray-700 flex flex-col max-h-[90vh]">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-700">
          <h3 className="text-xl font-bold text-white flex items-center gap-2">
            <span className="text-blue-500">📊</span> 历史信号统计 (近60天)
          </h3>
          <button 
            onClick={onClose}
            className="text-gray-400 hover:text-white transition-colors"
          >
            ✕
          </button>
        </div>

        {/* Content */}
        <div className="p-6 flex-1 overflow-y-auto">
          {loading ? (
            <div className="flex flex-col items-center justify-center py-12 text-gray-400">
              <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500 mb-4"></div>
              <p>正在分析历史数据...</p>
            </div>
          ) : error ? (
            <div className="p-4 bg-red-900/20 border border-red-700 text-red-200 rounded-md">
              错误: {error}
            </div>
          ) : result ? (
            <div className="space-y-6">
              {/* Key Metrics Grid */}
              <div className="grid grid-cols-2 gap-4">
                <div className="bg-[#0D1117] p-4 rounded-lg border border-gray-800">
                  <div className="text-gray-400 text-sm mb-1">总信号数</div>
                  <div className="text-2xl font-bold text-white">{result.totalSignals}</div>
                </div>
                <div className="bg-[#0D1117] p-4 rounded-lg border border-gray-800">
                  <div className="text-gray-400 text-sm mb-1">胜率 (T+5)</div>
                  <div className={`text-2xl font-bold ${result.winRate >= 0.5 ? 'text-red-400' : 'text-green-400'}`}>
                    {(result.winRate * 100).toFixed(1)}%
                  </div>
                </div>
                <div className="bg-[#0D1117] p-4 rounded-lg border border-gray-800">
                  <div className="text-gray-400 text-sm mb-1">平均收益率</div>
                  <div className={`text-2xl font-bold ${result.avgReturn >= 0 ? 'text-red-400' : 'text-green-400'}`}>
                    {(result.avgReturn * 100).toFixed(2)}%
                  </div>
                </div>
                <div className="bg-[#0D1117] p-4 rounded-lg border border-gray-800">
                  <div className="text-gray-400 text-sm mb-1">最大亏损</div>
                  <div className="text-2xl font-bold text-green-400">
                    {(result.maxLoss * 100).toFixed(2)}%
                  </div>
                </div>
              </div>

              {/* Best/Worst Performers */}
              <div className="bg-[#0D1117] rounded-lg border border-gray-800 overflow-hidden">
                <div className="px-4 py-3 bg-gray-800/50 text-sm font-semibold text-gray-300">
                  个股表现
                </div>
                <div className="p-4 space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-gray-400">最佳表现</span>
                    <span className="text-red-400 font-medium">{result.bestStock || '-'}</span>
                  </div>
                  <div className="w-full h-px bg-gray-800"></div>
                  <div className="flex justify-between items-center">
                    <span className="text-gray-400">最差表现</span>
                    <span className="text-green-400 font-medium">{result.worstStock || '-'}</span>
                  </div>
                </div>
              </div>

              <div className="text-xs text-gray-500 text-center">
                统计时间: {result.analysisDate} <br/>
                * 收益率计算基于信号发出日收盘价至T+5日收盘价
              </div>
            </div>
          ) : null}
        </div>
      </div>
    </div>
  );
};

export default StatsModal;
