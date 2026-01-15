import React, { useEffect, useState } from 'react';
import { useWailsAPI } from '../hooks/useWailsAPI';
import { StrategySignal } from '../types';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import { Clock, Activity } from 'lucide-react';

interface SignalListProps {
  onSelect: (signal: StrategySignal) => void;
  selectedId?: number;
}

const SignalList: React.FC<SignalListProps> = ({ onSelect, selectedId }) => {
  const { GetLatestSignals } = useWailsAPI();
  const [signals, setSignals] = useState<StrategySignal[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchSignals = async () => {
    try {
      const data = await GetLatestSignals(20);
      setSignals(data || []);
      if (data && data.length > 0 && !selectedId) {
        onSelect(data[0]);
      }
    } catch (error) {
      console.error('Failed to fetch signals:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSignals();

    const handler = (event: any) => {
      console.log('New signal received:', event);
      // Ensure the event data matches StrategySignal structure
      // Wails events might return the data directly or wrapped.
      // Assuming event is the signal object as emitted in Go.
      const newSignal = event as StrategySignal;
      setSignals(prev => [newSignal, ...prev]);
    };

    EventsOn('new_signal', handler);

    return () => {
      EventsOff('new_signal');
    };
  }, []);

  const getScoreColor = (score: number) => {
    if (score >= 80) return 'text-green-400 border-green-400';
    if (score < 60) return 'text-orange-400 border-orange-400';
    return 'text-blue-400 border-blue-400';
  };

  const getScoreRingColor = (score: number) => {
    if (score >= 80) return '#4ade80'; // green-400
    if (score < 60) return '#fb923c'; // orange-400
    return '#60a5fa'; // blue-400
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-full bg-[#0D1117] border-r border-white/10 w-[320px] shrink-0 text-gray-400">
        <div className="relative w-12 h-12 mb-3">
          <div className="absolute inset-0 border-4 border-blue-500/30 rounded-full"></div>
          <div className="absolute inset-0 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        </div>
        <p className="text-sm font-medium">正在全 A 股扫描中...</p>
        <p className="text-xs text-gray-600 mt-1">AI 引擎正在分析市场数据</p>
      </div>
    );
  }

  if (signals.length === 0) {
    return (
      <div className="flex flex-col h-full bg-[#0D1117] border-r border-white/10 w-[320px] shrink-0">
        <div className="p-4 border-b border-white/10 bg-[#161b22]/80 backdrop-blur-md sticky top-0 z-10">
          <div className="flex items-center gap-2 mb-1">
            <Activity className="w-5 h-5 text-blue-400" />
            <h2 className="text-lg font-semibold text-white">量化信号捕捉</h2>
          </div>
          <div className="flex justify-between items-center text-xs text-gray-400">
            <div className="flex items-center gap-1">
              <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
              实时监控中
            </div>
            <span>0 个信号</span>
          </div>
        </div>
        <div className="flex-1 flex flex-col items-center justify-center text-gray-500 p-4 text-center">
          <Activity className="w-12 h-12 mb-4 opacity-20" />
          <p className="text-sm">暂无新信号</p>
          <p className="text-xs text-gray-600 mt-2">AI 正在持续监控市场动态...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-[#0D1117] border-r border-white/10 w-[320px] shrink-0">
      <div className="p-4 border-b border-white/10 bg-[#161b22]/80 backdrop-blur-md sticky top-0 z-10">
        <div className="flex items-center gap-2 mb-1">
          <Activity className="w-5 h-5 text-blue-400" />
          <h2 className="text-lg font-semibold text-white">量化信号捕捉</h2>
        </div>
        <div className="flex justify-between items-center text-xs text-gray-400">
          <div className="flex items-center gap-1">
            <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
            实时监控中
          </div>
          <span>{signals.length} 个信号</span>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-3 space-y-3 custom-scrollbar">
        {signals.map((signal, index) => (
          <div
            key={signal.id || index}
            onClick={() => onSelect(signal)}
            className={`group relative p-3 rounded-xl border transition-all cursor-pointer ${
              selectedId === signal.id
                ? 'bg-blue-600/20 border-blue-500/50 shadow-[0_0_15px_rgba(37,99,235,0.2)]'
                : 'bg-gray-900/50 border-white/5 hover:bg-gray-800/60 hover:border-white/10'
            } backdrop-blur-md`}
          >
            <div className="flex justify-between items-start mb-2">
              <div>
                <div className="flex items-center gap-2">
                  <span className="font-bold text-white text-base">{signal.stockName || signal.code}</span>
                  <span className="text-xs text-gray-500">{signal.code}</span>
                </div>
                <div className="flex items-center gap-2 mt-1">
                  <span className={`text-xs px-1.5 py-0.5 rounded border ${
                    signal.signalType === 'B' 
                      ? 'bg-red-500/10 text-red-400 border-red-500/20' 
                      : 'bg-green-500/10 text-green-400 border-green-500/20'
                  }`}>
                    {signal.strategyName || 'AI量化'}
                  </span>
                </div>
              </div>

              {/* AI Score Ring */}
              <div className="relative w-12 h-12 flex items-center justify-center">
                <svg className="w-full h-full transform -rotate-90">
                  <circle
                    cx="24"
                    cy="24"
                    r="20"
                    stroke="#1f2937"
                    strokeWidth="4"
                    fill="none"
                  />
                  <circle
                    cx="24"
                    cy="24"
                    r="20"
                    stroke={getScoreRingColor(signal.aiScore)}
                    strokeWidth="4"
                    fill="none"
                    strokeDasharray={125.6}
                    strokeDashoffset={125.6 - (125.6 * signal.aiScore) / 100}
                    className="transition-all duration-1000 ease-out"
                  />
                </svg>
                <div className="absolute inset-0 flex flex-col items-center justify-center">
                  <span className="text-[10px] text-gray-400 scale-75">AI</span>
                  <span className={`text-sm font-bold leading-none ${getScoreColor(signal.aiScore).split(' ')[0]}`}>
                    {signal.aiScore}
                  </span>
                </div>
              </div>
            </div>

            <div className="flex justify-between items-center text-xs text-gray-500 mt-2 pt-2 border-t border-white/5">
              <div className="flex items-center gap-1">
                <Clock className="w-3 h-3" />
                {signal.tradeDate}
              </div>
              <div className="flex items-center gap-1">
                {/* Change percent isn't in StrategySignal. 
                    It is in MoneyFlowData details JSON. 
                    I'll skip it for now or try to parse details.
                */}
                <span className="text-white/40">详细 &gt;</span>
              </div>
            </div>
            
            {/* Selection Indicator */}
            {selectedId === signal.id && (
              <div className="absolute left-0 top-3 bottom-3 w-1 bg-blue-500 rounded-r-full shadow-[0_0_8px_rgba(59,130,246,0.6)]"></div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default SignalList;
