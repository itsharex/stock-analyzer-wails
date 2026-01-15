import React, { useState } from 'react';
import SignalList from '../components/SignalList';
import EnhancedKLineChart from '../components/EnhancedKLineChart';
import AIAnalysisPanel from '../components/AIAnalysisPanel';
import ScanButton from '../components/ScanButton';
import MoneyFlowSyncButton from '../components/MoneyFlowSyncButton';
import StatsModal from '../components/StatsModal';
import { StrategySignal } from '../types';

const DecisionPioneerPage: React.FC = () => {
  const [selectedSignal, setSelectedSignal] = useState<StrategySignal | null>(null);
  const [showStats, setShowStats] = useState(false);

  return (
    <div className="w-full h-[calc(100vh-4rem)] bg-[#0D1117] text-white flex items-stretch justify-center overflow-hidden">
      <StatsModal isOpen={showStats} onClose={() => setShowStats(false)} />
      
      <div className="w-full h-full flex bg-[#0D1117]">
        {/* Left 80px Rail (Figma: Dw) */}
        {/*<div className="w-20 h-full flex flex-col border-r border-white/10 shrink-0 items-center py-4 gap-6">*/}
        {/*  <div className="w-10 h-10 rounded-xl bg-blue-600/20 flex items-center justify-center text-blue-400">*/}
        {/*     <BarChart2 className="w-5 h-5" />*/}
        {/*  </div>*/}
        {/*  <div className="w-10 h-10 rounded-xl hover:bg-white/5 flex items-center justify-center text-gray-500 cursor-pointer">*/}
        {/*     <Layers className="w-5 h-5" />*/}
        {/*  </div>*/}
        {/*  <div className="mt-auto w-10 h-10 rounded-xl hover:bg-white/5 flex items-center justify-center text-gray-500 cursor-pointer">*/}
        {/*     <Settings className="w-5 h-5" />*/}
        {/*  </div>*/}
        {/*</div>*/}

        {/* Right Main Area */}
        <div className="flex-1 flex flex-col min-w-0">
          {/* Top Bar (Figma: Top Bar) */}
          <div className="h-16 w-full bg-[#161b22]/90 backdrop-blur border-b border-white/5 flex items-center px-6 shrink-0 justify-between">
            <h1 className="text-lg font-semibold tracking-wide flex items-center gap-2">
              <span className="w-2 h-6 bg-blue-500 rounded-full"></span>
              决策先锋
            </h1>
            <div className="flex items-center gap-4 text-sm text-gray-400">
              <span className="flex items-center gap-1.5">
                 <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                 AI 引擎运行中
              </span>
              <button
                onClick={() => setShowStats(true)}
                className="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded text-xs flex items-center gap-1 transition-colors border border-gray-700"
              >
                <span>📊</span> 历史统计
              </button>
              <div className="flex items-center gap-3">
                <MoneyFlowSyncButton />
                <ScanButton />
              </div>
            </div>
          </div>

          {/* Main Content (3-Column Layout) */}
          <div className="flex-1 w-full flex overflow-hidden">
            {/* 1. Left: Signal List */}
            <SignalList 
              onSelect={setSelectedSignal} 
              selectedId={selectedSignal?.id} 
            />

            {/* 2. Middle: Chart */}
            <div className="flex-1 flex flex-col min-w-0 border-r border-white/10 bg-[#0D1117]">
              {selectedSignal ? (
                <EnhancedKLineChart 
                  stockCode={selectedSignal.code} 
                  signal={selectedSignal} 
                  onSignalClick={setSelectedSignal}
                />
              ) : (
                <div className="flex-1 flex items-center justify-center text-gray-500">
                  请选择一个信号查看详情
                </div>
              )}
            </div>

            {/* 3. Right: AI Analysis Panel */}
            <AIAnalysisPanel signal={selectedSignal} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default DecisionPioneerPage;
