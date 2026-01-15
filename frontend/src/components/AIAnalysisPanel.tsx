import React, { useEffect, useState, useMemo } from 'react';
import { StrategySignal } from '../types';
import { useWailsAPI } from '../hooks/useWailsAPI';
import { Brain, Target, Users, PlusCircle, CheckCircle, AlertTriangle, Lightbulb } from 'lucide-react';

interface AIAnalysisPanelProps {
  signal: StrategySignal | null;
}

const AIAnalysisPanel: React.FC<AIAnalysisPanelProps> = ({ signal }) => {
  const { addToWatchlist } = useWailsAPI();
  const [animatedScore, setAnimatedScore] = useState(0);
  const [addedToWatchlist, setAddedToWatchlist] = useState(false);
  const [activeTab, setActiveTab] = useState<'chip' | 'risk' | 'logic'>('chip');

  useEffect(() => {
    if (signal) {
      setAnimatedScore(0);
      setAddedToWatchlist(false);
      let start = 0;
      const end = signal.aiScore;
      const duration = 1000;
      const startTime = performance.now();

      const animate = (currentTime: number) => {
        const elapsed = currentTime - startTime;
        const progress = Math.min(elapsed / duration, 1);
        const easeOutQuart = 1 - Math.pow(1 - progress, 4); // Easing function
        
        setAnimatedScore(Math.floor(start + (end - start) * easeOutQuart));

        if (progress < 1) {
          requestAnimationFrame(animate);
        }
      };

      requestAnimationFrame(animate);
    }
  }, [signal]);

  const handleAddToWatchlist = async () => {
    if (!signal) return;
    try {
      // Create a minimal StockData object for adding to watchlist
      // In a real scenario, we might want to fetch full stock data first
      // or ensure backend can handle partial data.
      await addToWatchlist({
        code: signal.code,
        name: signal.code, // Placeholder if name is missing
        price: 0, change: 0, changeRate: 0, volume: 0, amount: 0,
        high: 0, low: 0, open: 0, preClose: 0, amplitude: 0,
        turnover: 0, pe: 0, pb: 0, totalMV: 0, circMV: 0,
        volumeRatio: 0, warrantRatio: 0
      });
      setAddedToWatchlist(true);
    } catch (error) {
      console.error('Failed to add to watchlist', error);
    }
  };

  const details = useMemo(() => {
    if (!signal || !signal.details) return {};
    try {
      return JSON.parse(signal.details);
    } catch (e) {
      return {};
    }
  }, [signal]);

  if (!signal) {
    return (
      <div className="w-[360px] bg-[#0D1117] border-l border-white/10 flex flex-col items-center justify-center text-gray-500 shrink-0">
        <Brain className="w-12 h-12 mb-4 opacity-20" />
        <p>请选择一个信号查看 AI 解析</p>
      </div>
    );
  }

  // Parse keywords from aiReason (assuming it might be comma separated or just text)
  // If aiReason is a sentence, we might just display it.
  const keywords = signal.aiReason.split(/[,，]/).filter(k => k.length > 0);

  return (
    <div className="w-[360px] bg-[#0D1117] border-l border-white/10 flex flex-col h-full shrink-0">
      <div className="p-6 border-b border-white/10">
        <h2 className="text-xl font-bold text-white flex items-center gap-2">
          <Brain className="w-6 h-6 text-purple-400" />
          AI 深度解析
        </h2>
        <p className="text-gray-400 text-xs mt-1">基于大模型的多维度量化评估</p>
      </div>

      <div className="flex-1 overflow-y-auto p-6 space-y-6 custom-scrollbar">
        {/* Dashboard Animation (Semi-circle Gauge) */}
        <div className="flex flex-col items-center justify-center relative pt-4 pb-2">
          <div className="relative w-48 h-24 overflow-hidden">
             {/* Background Arc */}
             <svg viewBox="0 0 100 50" className="w-full h-full transform translate-y-1">
               <path d="M 10 50 A 40 40 0 0 1 90 50" fill="none" stroke="#1f2937" strokeWidth="8" strokeLinecap="round" />
             </svg>
             {/* Progress Arc */}
             <svg viewBox="0 0 100 50" className="absolute top-0 left-0 w-full h-full transform translate-y-1">
               <path 
                 d="M 10 50 A 40 40 0 0 1 90 50" 
                 fill="none" 
                 stroke={animatedScore >= 80 ? '#4ade80' : animatedScore < 60 ? '#fb923c' : '#a855f7'} 
                 strokeWidth="8" 
                 strokeLinecap="round"
                 strokeDasharray="126"
                 strokeDashoffset={126 - (126 * animatedScore / 100)}
                 className="transition-all duration-300 ease-out"
               />
             </svg>
          </div>
          
          <div className="absolute top-16 text-center">
            <span className="text-5xl font-black text-white tracking-tighter">{animatedScore}</span>
            <span className="text-sm text-gray-400 block -mt-1">AI 综合评分</span>
          </div>
        </div>

        {/* Tabs */}
        <div className="flex border-b border-white/10">
          <button
            onClick={() => setActiveTab('chip')}
            className={`flex-1 pb-3 text-sm font-medium transition-colors relative ${
              activeTab === 'chip' ? 'text-blue-400' : 'text-gray-400 hover:text-gray-300'
            }`}
          >
            筹码解析
            {activeTab === 'chip' && <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-500 rounded-t-full"></div>}
          </button>
          <button
            onClick={() => setActiveTab('risk')}
            className={`flex-1 pb-3 text-sm font-medium transition-colors relative ${
              activeTab === 'risk' ? 'text-orange-400' : 'text-gray-400 hover:text-gray-300'
            }`}
          >
            风险预警
            {activeTab === 'risk' && <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-orange-500 rounded-t-full"></div>}
          </button>
          <button
            onClick={() => setActiveTab('logic')}
            className={`flex-1 pb-3 text-sm font-medium transition-colors relative ${
              activeTab === 'logic' ? 'text-purple-400' : 'text-gray-400 hover:text-gray-300'
            }`}
          >
            逻辑推演
            {activeTab === 'logic' && <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-purple-500 rounded-t-full"></div>}
          </button>
        </div>

        {/* Tab Content */}
        <div className="min-h-[120px]">
          {activeTab === 'chip' && (
            <div className="space-y-3 animate-in fade-in slide-in-from-bottom-2 duration-300">
               <div className="flex flex-wrap gap-2">
                {keywords.map((k, i) => (
                  <span key={i} className="px-3 py-1 bg-blue-500/10 text-blue-300 rounded-full text-xs border border-blue-500/20">
                    {k}
                  </span>
                ))}
              </div>
              <p className="text-xs text-gray-400 leading-relaxed">
                当前主力资金介入明显，筹码集中度较高，上方套牢盘压力较小。
              </p>
            </div>
          )}

          {activeTab === 'risk' && (
            <div className="space-y-3 animate-in fade-in slide-in-from-bottom-2 duration-300">
              <div className="flex items-start gap-2 text-orange-300 bg-orange-500/10 p-3 rounded-lg border border-orange-500/20">
                <AlertTriangle className="w-4 h-4 shrink-0 mt-0.5" />
                <span className="text-xs">注意短期乖离率过大，若量能无法持续放大，可能面临回调风险。</span>
              </div>
              <div className="text-xs text-gray-400">
                风险等级: <span className="text-white font-mono">{signal.riskLevel || 'Medium'}</span>
              </div>
            </div>
          )}

          {activeTab === 'logic' && (
            <div className="space-y-3 animate-in fade-in slide-in-from-bottom-2 duration-300">
               <div className="flex items-start gap-2 text-purple-300 bg-purple-500/10 p-3 rounded-lg border border-purple-500/20">
                <Lightbulb className="w-4 h-4 shrink-0 mt-0.5" />
                <span className="text-xs">{signal.aiReason}</span>
              </div>
            </div>
          )}
        </div>

        {/* Chip Distribution (Cost) */}
        <div className="space-y-5 pt-4 border-t border-white/5">
          <h3 className="text-sm font-semibold text-white mb-2">筹码分布</h3>
          
          <div>
            <div className="flex justify-between text-xs mb-1.5">
              <span className="text-gray-400 flex items-center gap-1">
                <Target className="w-3 h-3" /> 成本区间 (获利比例)
              </span>
              <span className="text-green-400 font-mono">82%</span>
            </div>
            <div className="h-1.5 bg-gray-800 rounded-full overflow-hidden flex">
               {/* Mock distribution */}
               <div className="w-[15%] bg-gray-600"></div>
               <div className="w-[60%] bg-red-500/80"></div>
               <div className="w-[25%] bg-gray-600"></div>
            </div>
             <div className="flex justify-between text-[10px] text-gray-500 mt-1 font-mono">
               <span>9.50</span>
               <span>Current: {details.close || '-'}</span>
               <span>12.80</span>
            </div>
          </div>

          <div>
            <div className="flex justify-between text-xs mb-1.5">
              <span className="text-gray-400 flex items-center gap-1">
                <Users className="w-3 h-3" /> 主力成本 (预估)
              </span>
              <span className="text-blue-400 font-mono">{details.close ? (details.close * 0.92).toFixed(2) : '-'}</span>
            </div>
            <div className="relative h-6 bg-gray-800/50 rounded overflow-hidden flex items-center px-2">
               {/* Main cost bar relative to range */}
               <div className="absolute left-0 top-0 bottom-0 bg-blue-500/20 w-[70%] border-r border-blue-500/50"></div>
               <span className="relative text-[10px] text-blue-300 z-10">主力成本区支撑有效</span>
            </div>
          </div>
        </div>

        {/* Action Button */}
        <div className="pt-2">
          <button
            onClick={handleAddToWatchlist}
            disabled={addedToWatchlist}
            className={`w-full py-3 rounded-xl font-bold flex items-center justify-center gap-2 transition-all ${
              addedToWatchlist 
                ? 'bg-green-600/20 text-green-400 cursor-default'
                : 'bg-blue-600 hover:bg-blue-500 text-white shadow-lg shadow-blue-900/50'
            }`}
          >
            {addedToWatchlist ? (
              <>
                <CheckCircle className="w-5 h-5" />
                已加入自选
              </>
            ) : (
              <>
                <PlusCircle className="w-5 h-5" />
                加入自选监控
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default AIAnalysisPanel;
