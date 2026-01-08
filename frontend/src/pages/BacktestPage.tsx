import React, { useState } from 'react';
import BacktestPanel from '../components/BacktestPanelEnhanced';
import { Search } from 'lucide-react';

const BacktestPage: React.FC = () => {
  const [selectedStock, setSelectedStock] = useState<string>('600519');
  const [stockInput, setStockInput] = useState<string>('');

  const handleSearchStock = () => {
    if (stockInput.trim()) {
      setSelectedStock(stockInput.trim());
      setStockInput('');
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleSearchStock();
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* 页面标题 */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2">策略回测</h1>
          <p className="text-gray-400">使用历史数据验证交易策略的有效性</p>
        </div>

        {/* 股票搜索栏 */}
        <div className="mb-6 bg-gray-800 p-4 rounded-lg shadow-lg">
          <div className="flex gap-2">
            <div className="flex-1 relative">
              <input
                type="text"
                value={stockInput}
                onChange={(e) => setStockInput(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="输入股票代码 (如: 600519, SH600519)"
                className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
              />
              <Search className="absolute right-3 top-2.5 w-5 h-5 text-gray-400" />
            </div>
            <button
              onClick={handleSearchStock}
              className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md transition-colors"
            >
              搜索
            </button>
          </div>
          <div className="mt-3 text-sm text-gray-400">
            当前股票: <span className="font-semibold text-blue-400">{selectedStock}</span>
          </div>
        </div>

        {/* 回测面板 */}
        <div className="bg-gray-800 rounded-lg shadow-lg p-6">
          <BacktestPanel stockCode={selectedStock} />
        </div>
      </div>
    </div>
  );
};

export default BacktestPage;
