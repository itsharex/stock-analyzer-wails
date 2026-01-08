import React, { useState, useEffect } from 'react';
import { BacktestResult } from '../types';
import { useWailsAPI } from '../hooks/useWailsAPI';
import { parseError } from '../utils/errorHandler';
import { Save, BookOpen } from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface StrategyConfig {
  id: number;
  name: string;
  description: string;
  strategyType: string;
  parameters: Record<string, any>;
}

interface BacktestPanelProps {
  stockCode: string;
}

const BacktestPanel: React.FC<BacktestPanelProps> = ({ stockCode }) => {
  const { BacktestSimpleMA, BacktestMACD, GetAllStrategies, CreateStrategy, UpdateStrategyBacktestResult } = useWailsAPI();

  const [shortPeriod, setShortPeriod] = useState<number>(5);
  const [longPeriod, setLongPeriod] = useState<number>(20);
  const [signalPeriod, setSignalPeriod] = useState<number>(9);
  const [initialCapital, setInitialCapital] = useState<number>(100000);
  const [startDate, setStartDate] = useState<string>('2023-01-01');
  const [endDate, setEndDate] = useState<string>('2023-12-31');

  // 策略相关状态
  const [strategies, setStrategies] = useState<StrategyConfig[]>([]);
  const [selectedStrategy, setSelectedStrategy] = useState<StrategyConfig | null>(null);
  const [showSaveModal, setShowSaveModal] = useState(false);
  const [strategyName, setStrategyName] = useState('');
  const [strategyDescription, setStrategyDescription] = useState('');

  const [backtestResult, setBacktestResult] = useState<BacktestResult | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // 加载策略列表
  useEffect(() => {
    loadStrategies();
  }, []);

  const loadStrategies = async () => {
    try {
      const result = await GetAllStrategies();
      if (Array.isArray(result)) {
        // 显示所有策略
        setStrategies(result);
      }
    } catch (err) {
      console.error('加载策略列表失败:', err);
    }
  };

  const handleBacktest = async () => {
    setLoading(true);
    setError(null);
    setBacktestResult(null);
    try {
      let result;
      
      if (selectedStrategy && selectedStrategy.strategyType === 'macd') {
        // MACD策略回测
        result = await BacktestMACD(stockCode, shortPeriod, longPeriod, signalPeriod, initialCapital, startDate, endDate);
      } else {
        // 默认使用双均线策略回测
        result = await BacktestSimpleMA(stockCode, shortPeriod, longPeriod, initialCapital, startDate, endDate);
      }
      
      setBacktestResult(result);

      // 如果是从策略库执行的，更新策略的回测结果
      if (selectedStrategy && (selectedStrategy.strategyType === 'simple_ma' || selectedStrategy.strategyType === 'macd')) {
        try {
          await UpdateStrategyBacktestResult(selectedStrategy.id, {
            totalReturn: result.totalReturn,
            annualizedReturn: result.annualizedReturn,
            maxDrawdown: result.maxDrawdown,
            winRate: result.winRate,
            tradeCount: result.tradeCount,
            initialCapital: result.initialCapital,
            finalCapital: result.finalCapital,
          });
          // 重新加载策略列表以更新显示的回测结果
          await loadStrategies();
        } catch (err) {
          console.error('更新策略回测结果失败:', err);
        }
      }
    } catch (err) {
      const errorResult = parseError(err);
      setError(errorResult.message);
    } finally {
      setLoading(false);
    }
  };

  const handleSelectStrategy = (strategyId: string) => {
    if (strategyId === '') {
      setSelectedStrategy(null);
      return;
    }

    const strategy = strategies.find(s => s.id === parseInt(strategyId));
    if (strategy) {
      setSelectedStrategy(strategy);

      // 根据策略类型应用不同的参数
      if (strategy.strategyType === 'simple_ma') {
        // 双均线策略参数
        if (strategy.parameters.shortPeriod) {
          setShortPeriod(Number(strategy.parameters.shortPeriod));
        }
        if (strategy.parameters.longPeriod) {
          setLongPeriod(Number(strategy.parameters.longPeriod));
        }
      } else if (strategy.strategyType === 'macd') {
        // MACD 策略参数
        if (strategy.parameters.fastPeriod) {
          setShortPeriod(Number(strategy.parameters.fastPeriod));
        }
        if (strategy.parameters.slowPeriod) {
          setLongPeriod(Number(strategy.parameters.slowPeriod));
        }
        if (strategy.parameters.signalPeriod) {
          setSignalPeriod(Number(strategy.parameters.signalPeriod));
        }
      }

      // 通用参数
      if (strategy.parameters.initialCapital) {
        setInitialCapital(Number(strategy.parameters.initialCapital));
      }
    }
  };

  const handleSaveAsStrategy = () => {
    setShowSaveModal(true);
    setStrategyName('');
    setStrategyDescription('');
  };

  const handleSaveStrategy = async () => {
    if (!strategyName.trim()) {
      alert('请输入策略名称');
      return;
    }

    try {
      // 根据当前选择或默认的策略类型来确定保存类型
      const strategyType = selectedStrategy ? selectedStrategy.strategyType : 'simple_ma';
      
      // 根据策略类型构建不同的参数
      let parameters: any = {
        initialCapital,
      };

      if (strategyType === 'macd') {
        parameters = {
          ...parameters,
          fastPeriod: shortPeriod,
          slowPeriod: longPeriod,
          signalPeriod: signalPeriod,
        };
      } else {
        parameters = {
          ...parameters,
          shortPeriod,
          longPeriod,
        };
      }

      await CreateStrategy(
        strategyName,
        strategyDescription,
        strategyType,
        parameters
      );
      setShowSaveModal(false);
      setStrategyName('');
      setStrategyDescription('');
      // 重新加载策略列表
      await loadStrategies();
      alert('策略保存成功！');
    } catch (err: any) {
      alert(`保存策略失败: ${err.message}`);
    }
  };

  const formatCurrency = (value: number) => `¥${value.toFixed(2)}`;
  const formatPercentage = (value: number) => `${value.toFixed(2)}%`;

  const equityChartData = backtestResult?.equityCurve.map((value, index) => ({
    date: backtestResult.equityDates[index],
    capital: value,
  })) || [];

  // 根据策略类型获取参数标签
  const getShortPeriodLabel = () => {
    if (selectedStrategy) {
      return selectedStrategy.strategyType === 'macd' ? '快线周期' : '短周期均线 (MA)';
    }
    return '短周期均线 (MA)';
  };

  const getLongPeriodLabel = () => {
    if (selectedStrategy) {
      return selectedStrategy.strategyType === 'macd' ? '慢线周期' : '长周期均线 (MA)';
    }
    return '长周期均线 (MA)';
  };

  return (
    <div className="p-4 bg-gray-800 text-gray-100 rounded-lg shadow-lg">
      <h2 className="text-2xl font-bold mb-4">回测面板 - {stockCode}</h2>

      {/* 策略选择和保存区域 */}
      <div className="mb-6 p-4 bg-gray-700 rounded-md">
        <div className="flex items-center gap-4 mb-3">
          <div className="flex items-center gap-2">
            <BookOpen className="w-5 h-5 text-blue-400" />
            <label className="text-sm font-medium text-gray-300">从策略库选择:</label>
          </div>
          <select
            value={selectedStrategy ? selectedStrategy.id.toString() : ''}
            onChange={(e) => handleSelectStrategy(e.target.value)}
            className="flex-1 px-3 py-2 rounded-md bg-gray-600 border border-gray-500 text-gray-100 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
          >
            <option value="">手动输入参数</option>
            {strategies.map((strategy) => {
              const strategyTypeName = strategy.strategyType === 'simple_ma' ? '双均线' :
                                      strategy.strategyType === 'macd' ? 'MACD' :
                                      strategy.strategyType;
              const isSupported = strategy.strategyType === 'simple_ma' || strategy.strategyType === 'macd';
              return (
                <option key={strategy.id} value={strategy.id.toString()} disabled={!isSupported}>
                  {strategy.name} ({strategyTypeName}) {!isSupported && '[暂不支持回测]'}
                </option>
              );
            })}
          </select>
          <button
            onClick={handleSaveAsStrategy}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md transition-colors"
            title="保存为新策略"
          >
            <Save className="w-4 h-4" />
            保存为新策略
          </button>
        </div>
        {selectedStrategy && (
          <div className="text-xs text-gray-400">
            已加载策略: <span className="font-semibold text-blue-400">{selectedStrategy.name}</span>
            {selectedStrategy.description && ` - ${selectedStrategy.description}`}
          </div>
        )}
      </div>

      {/* 参数配置 */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <div>
          <label htmlFor="shortPeriod" className="block text-sm font-medium text-gray-300">{getShortPeriodLabel()}:</label>
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
          <label htmlFor="longPeriod" className="block text-sm font-medium text-gray-300">{getLongPeriodLabel()}:</label>
          <input
            type="number"
            id="longPeriod"
            value={longPeriod}
            onChange={(e) => setLongPeriod(parseInt(e.target.value))}
            className="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-gray-100 shadow-sm focus:border-blue-500 focus:ring-blue-500"
            min="1"
          />
        </div>
        {/* MACD策略的信号线周期参数 */}
        {(!selectedStrategy || selectedStrategy.strategyType === 'macd') && (
          <div>
            <label htmlFor="signalPeriod" className="block text-sm font-medium text-gray-300">信号线周期 (DEA):</label>
            <input
              type="number"
              id="signalPeriod"
              value={signalPeriod}
              onChange={(e) => setSignalPeriod(parseInt(e.target.value))}
              className="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-gray-100 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              min="1"
            />
          </div>
        )}
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
                </tr>
              </thead>
              <tbody className="bg-gray-700 divide-y divide-gray-600">
                {backtestResult.trades.map((trade, index) => (
                  <tr key={index} className={trade.type === 'BUY' ? 'bg-green-900/20' : 'bg-red-900/20'}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{trade.time}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                        trade.type === 'BUY' ? 'bg-green-900/50 text-green-300' : 'bg-red-900/50 text-red-300'
                      }`}>
                        {trade.type === 'BUY' ? '买入' : '卖出'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{trade.price.toFixed(2)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{trade.volume}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{trade.amount.toFixed(2)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{trade.commission.toFixed(2)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* 保存策略 Modal */}
      {showSaveModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-lg shadow-2xl w-full max-w-md">
            <div className="p-6">
              <h3 className="text-xl font-bold mb-4">保存为新策略</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-2">策略名称 *</label>
                  <input
                    type="text"
                    value={strategyName}
                    onChange={(e) => setStrategyName(e.target.value)}
                    className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                    placeholder="输入策略名称"
                    autoFocus
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-2">策略描述</label>
                  <textarea
                    value={strategyDescription}
                    onChange={(e) => setStrategyDescription(e.target.value)}
                    rows={3}
                    className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                    placeholder="输入策略描述（可选）"
                  />
                </div>
              </div>
              <div className="flex justify-end gap-3 mt-6">
                <button
                  onClick={() => setShowSaveModal(false)}
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white font-medium rounded-md transition-colors"
                >
                  取消
                </button>
                <button
                  onClick={handleSaveStrategy}
                  className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md transition-colors"
                >
                  <Save className="w-4 h-4" />
                  保存
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default BacktestPanel;
