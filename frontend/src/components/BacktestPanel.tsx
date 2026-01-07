
import React, { useState, useEffect } from 'react';
import { BacktestResult } from '../types';
import { useWailsAPI } from '../hooks/useWailsAPI';
import { parseError } from '../utils/errorHandler';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface BacktestPanelProps {
  stockCode: string;
}

const BacktestPanel: React.FC<BacktestPanelProps> = ({ stockCode }) => {
  const { BacktestSimpleMA } = useWailsAPI();

  const [shortPeriod, setShortPeriod] = useState<number>(5);
  const [longPeriod, setLongPeriod] = useState<number>(20);
  const [initialCapital, setInitialCapital] = useState<number>(100000);
  const [startDate, setStartDate] = useState<string>('2023-01-01');
  const [endDate, setEndDate] = useState<string>('2023-12-31');
  const [backtestResult, setBacktestResult] = useState<BacktestResult | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const handleBacktest = async () => {
    setLoading(true);
    setError(null);
    setBacktestResult(null);
    try {
      const result = await BacktestSimpleMA(stockCode, shortPeriod, longPeriod, initialCapital, startDate, endDate);
      setBacktestResult(result);
    } catch (err) {
      const errorResult = parseError(err);
      setError(errorResult.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // Optionally trigger backtest on stockCode change or initial load
    if (stockCode) {
      // handleBacktest(); // Uncomment if you want to auto-run backtest
    }
  }, [stockCode]);

  const formatCurrency = (value: number) => `¥${value.toFixed(2)}`;
  const formatPercentage = (value: number) => `${value.toFixed(2)}%`;

  const equityChartData = backtestResult?.equityCurve.map((value, index) => ({
    date: backtestResult.equityDates[index],
    capital: value,
  })) || [];

  return (
    <div className="p-4 bg-gray-800 text-gray-100 rounded-lg shadow-lg">
      <h2 className="text-2xl font-bold mb-4">回测面板 - {stockCode}</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <div>
          <label htmlFor="shortPeriod" className="block text-sm font-medium text-gray-300">短周期均线 (MA):</label>
          <input
            type="number"
            id="shortPeriod"
            value={shortPeriod}
            onChange={(e) => setShortPeriod(parseInt(e.target.value))}
            className="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-gray-100 shadow-sm focus:border-blue-500 focus:ring-blue-500"
            min="1"
          />
        </div>
        <div>
          <label htmlFor="longPeriod" className="block text-sm font-medium text-gray-300">长周期均线 (MA):</label>
          <input
            type="number"
            id="longPeriod"
            value={longPeriod}
            onChange={(e) => setLongPeriod(parseInt(e.target.value))}
            className="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-gray-100 shadow-sm focus:border-blue-500 focus:ring-blue-500"
            min="1"
          />
        </div>
        <div>
          <label htmlFor="initialCapital" className="block text-sm font-medium text-gray-300">初始资金:</label>
          <input
            type="number"
            id="initialCapital"
            value={initialCapital}
            onChange={(e) => setInitialCapital(parseFloat(e.target.value))}
            className="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-gray-100 shadow-sm focus:border-blue-500 focus:ring-blue-500"
            min="1000"
          />
        </div>
        <div>
          <label htmlFor="startDate" className="block text-sm font-medium text-gray-300">开始日期:</label>
          <input
            type="date"
            id="startDate"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
            className="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-gray-100 shadow-sm focus:border-blue-500 focus:ring-blue-500"
          />
        </div>
        <div>
          <label htmlFor="endDate" className="block text-sm font-medium text-gray-300">结束日期:</label>
          <input
            type="date"
            id="endDate"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
            className="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-gray-100 shadow-sm focus:border-blue-500 focus:ring-blue-500"
          />
        </div>
      </div>

      <button
        onClick={handleBacktest}
        className="w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        disabled={loading}
      >
        {loading ? '回测中...' : '开始回测'}
      </button>

      {error && <div className="mt-4 text-red-400">错误: {error}</div>}

      {backtestResult && (
        <div className="mt-6">
          <h3 className="text-xl font-bold mb-4">回测结果 ({backtestResult.strategyName})</h3>

          {/* 结果摘要 */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
            <div className="bg-gray-700 p-4 rounded-md">
              <p className="text-sm text-gray-400">总收益率</p>
              <p className={`text-lg font-semibold ${backtestResult.totalReturn >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                {formatPercentage(backtestResult.totalReturn)}
              </p>
            </div>
            <div className="bg-gray-700 p-4 rounded-md">
              <p className="text-sm text-gray-400">年化收益率</p>
              <p className={`text-lg font-semibold ${backtestResult.annualizedReturn >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                {formatPercentage(backtestResult.annualizedReturn)}
              </p>
            </div>
            <div className="bg-gray-700 p-4 rounded-md">
              <p className="text-sm text-gray-400">最大回撤</p>
              <p className="text-red-400 text-lg font-semibold">
                {formatPercentage(backtestResult.maxDrawdown)}
              </p>
            </div>
            <div className="bg-gray-700 p-4 rounded-md">
              <p className="text-sm text-gray-400">胜率</p>
              <p className="text-lg font-semibold text-blue-400">
                {formatPercentage(backtestResult.winRate)}
              </p>
            </div>
            <div className="bg-gray-700 p-4 rounded-md">
              <p className="text-sm text-gray-400">交易次数</p>
              <p className="text-lg font-semibold">
                {backtestResult.tradeCount}
              </p>
            </div>
            <div className="bg-gray-700 p-4 rounded-md">
              <p className="text-sm text-gray-400">初始资金</p>
              <p className="text-lg font-semibold">
                {formatCurrency(backtestResult.initialCapital)}
              </p>
            </div>
            <div className="bg-gray-700 p-4 rounded-md">
              <p className="text-sm text-gray-400">最终资金</p>
              <p className={`text-lg font-semibold ${backtestResult.finalCapital >= backtestResult.initialCapital ? 'text-green-400' : 'text-red-400'}`}>
                {formatCurrency(backtestResult.finalCapital)}
              </p>
            </div>
          </div>

          {/* 净值曲线图 */}
          <h4 className="text-lg font-bold mb-2">净值曲线</h4>
          <div className="bg-gray-700 p-4 rounded-md mb-6" style={{ height: '400px' }}>
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={equityChartData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#4a5568" />
                <XAxis dataKey="date" stroke="#cbd5e0" tickFormatter={(tick) => new Date(tick).toLocaleDateString()} />
                <YAxis stroke="#cbd5e0" tickFormatter={(value: number) => formatCurrency(value)} />
                <Tooltip formatter={(value: number | undefined) => value !== undefined ? formatCurrency(value) : ''} labelFormatter={(label: string) => `日期: ${label}`} />
                <Legend />
                <Line type="monotone" dataKey="capital" stroke="#4299e1" dot={false} name="总资产" />
              </LineChart>
            </ResponsiveContainer>
          </div>

          {/* 交易明细 */}
          <h4 className="text-lg font-bold mb-2">交易明细</h4>
          <div className="overflow-x-auto bg-gray-700 rounded-md">
            <table className="min-w-full divide-y divide-gray-600">
              <thead className="bg-gray-700">
                <tr>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">日期</th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">类型</th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">价格</th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">数量</th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">金额</th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">佣金</th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">印花税</th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">盈亏</th>
                </tr>
              </thead>
              <tbody className="bg-gray-800 divide-y divide-gray-700">
                {backtestResult.trades.map((trade, index) => (
                  <tr key={index}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-200">{trade.time}</td>
                    <td className={`px-6 py-4 whitespace-nowrap text-sm font-medium ${trade.type === 'BUY' ? 'text-green-400' : 'text-red-400'}`}>{trade.type}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-200">{formatCurrency(trade.price)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-200">{trade.volume}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-200">{formatCurrency(trade.amount)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-200">{formatCurrency(trade.commission)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-200">{formatCurrency(trade.tax)}</td>
                    <td className={`px-6 py-4 whitespace-nowrap text-sm ${trade.profit >= 0 ? 'text-green-400' : 'text-red-400'}`}>{formatCurrency(trade.profit)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
};

export default BacktestPanel;
