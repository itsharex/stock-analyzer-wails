import React, { useState } from 'react';
import SignalList from '../components/SignalList';
import EnhancedKLineChart from '../components/EnhancedKLineChart';
import AIAnalysisPanel from '../components/AIAnalysisPanel';
import ScanButton from '../components/ScanButton';
import MoneyFlowSyncButton from '../components/MoneyFlowSyncButton';
import StatsModal from '../components/StatsModal';
import { StrategySignal } from '../types';
import { useWailsAPI } from '../hooks/useWailsAPI';
import { toast } from 'react-hot-toast';

const DecisionPioneerPage: React.FC = () => {
  const [selectedSignal, setSelectedSignal] = useState<StrategySignal | null>(null);
  const [showStats, setShowStats] = useState(false);
  const [searchCode, setSearchCode] = useState('');
  const [isScanning, setIsScanning] = useState(false);
  const { ScanSingleStock } = useWailsAPI();

  const handleSingleScan = async () => {
    if (!searchCode || searchCode.length < 6) {
      toast.error('请输入正确的股票代码');
      return;
    }

    setIsScanning(true);
    try {
      const signals = await ScanSingleStock(searchCode);
      if (signals && signals.length > 0) {
        toast.success(`扫描完成，发现 ${signals.length} 个信号`);
        // 自动选中最新的一个信号
        setSelectedSignal(signals[0]);
      } else {
        toast.success('扫描完成，未发现符合策略的信号');
      }
    } catch (error) {
      console.error('扫描失败:', error);
      toast.error('扫描失败: ' + (error as any).message);
    } finally {
      setIsScanning(false);
    }
  };

  return (
    <div className="w-full h-[calc(100vh-4rem)] bg-[#0D1117] text-white flex items-stretch justify-center overflow-hidden">
      <StatsModal isOpen={showStats} onClose={() => setShowStats(false)} />
      
      <div className="w-full h-full flex bg-[#0D1117]">
        {/* Right Main Area */}
        <div className="flex-1 flex flex-col min-w-0">
          {/* Top Bar (Figma: Top Bar) */}
          <div className="h-16 w-full bg-[#161b22]/90 backdrop-blur border-b border-white/5 flex items-center px-6 shrink-0 justify-between">
            <h1 className="text-lg font-semibold tracking-wide flex items-center gap-2">
              <span className="w-2 h-6 bg-blue-500 rounded-full"></span>
              决策先锋
            </h1>
            
            {/* 顶部搜索栏 */}
            <div className="flex items-center gap-2 mx-4 bg-gray-800/50 rounded-lg p-1 border border-white/10">
              <input
                type="text"
                placeholder="输入代码 (如 000001)"
                className="bg-transparent border-none outline-none px-2 text-sm w-32 placeholder-gray-500"
                value={searchCode}
                onChange={(e) => setSearchCode(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleSingleScan()}
              />
              <button
                onClick={handleSingleScan}
                disabled={isScanning}
                className={`px-3 py-1 text-xs rounded transition-colors ${
                  isScanning 
                    ? 'bg-blue-500/20 text-blue-300 cursor-wait' 
                    : 'bg-blue-600 hover:bg-blue-500 text-white'
                }`}
              >
                {isScanning ? '扫描中...' : '单股扫描'}
              </button>
            </div>

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
                  请选择一个信号或在顶部搜索股票查看详情
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
