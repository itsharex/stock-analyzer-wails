import React, { useState, useEffect } from 'react';
import { Bell, History, Settings, Shield, Search, Calendar, Clock, Anchor, Sword, Cpu } from 'lucide-react';
import { useWailsAPI } from '../hooks/useWailsAPI';

export const AlertCenter: React.FC = () => {
  const { getAlertHistory, getAlertConfig, updateAlertConfig } = useWailsAPI();
  const [history, setHistory] = useState<any[]>([]);
  const [config, setConfig] = useState<any>({ sensitivity: 0.005, cooldown: 1, enabled: true });
  const [loading, setLoading] = useState(true);
  const [filterCode, setFilterCode] = useState('');

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [hist, cfg] = await Promise.all([
        getAlertHistory('', 50),
        getAlertConfig()
      ]);
      setHistory(Array.isArray(hist) ? hist : []);
      if (cfg) setConfig(cfg);
    } catch (err) {
      console.error('Failed to fetch alert data:', err);
      setHistory([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSaveConfig = async (newConfig: any) => {
    try {
      await updateAlertConfig(newConfig);
      setConfig(newConfig);
    } catch (err) {
      console.error('Failed to save config:', err);
    }
  };

  const getRoleIcon = (role: string) => {
    switch (role) {
      case 'conservative': return <Anchor className="w-3 h-3 text-blue-600" />;
      case 'aggressive': return <Sword className="w-3 h-3 text-red-600" />;
      default: return <Cpu className="w-3 h-3 text-slate-600" />;
    }
  };

  return (
    <div className="space-y-6">
      {/* 预警配置卡片 */}
      <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="p-6 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-blue-600 rounded-xl shadow-lg shadow-blue-200">
              <Settings className="w-5 h-5 text-white" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-slate-800">预警系统配置</h3>
              <p className="text-xs text-slate-500">个性化您的 AI 盯盘策略</p>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <span className="text-sm font-medium text-slate-600">{config.enabled ? '已开启全局预警' : '已关闭全局预警'}</span>
            <button 
              onClick={() => handleSaveConfig({ ...config, enabled: !config.enabled })}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none ${config.enabled ? 'bg-blue-600' : 'bg-slate-300'}`}
            >
              <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${config.enabled ? 'translate-x-6' : 'translate-x-1'}`} />
            </button>
          </div>
        </div>
        
        <div className="p-6 grid grid-cols-1 md:grid-cols-2 gap-8">
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <label className="text-sm font-bold text-slate-700 flex items-center">
                <Shield className="w-4 h-4 mr-2 text-blue-500" />
                预警灵敏度
              </label>
              <span className="text-xs font-mono bg-blue-50 text-blue-700 px-2 py-1 rounded">{(config.sensitivity * 100).toFixed(1)}%</span>
            </div>
            <input 
              type="range" 
              min="0.001" 
              max="0.02" 
              step="0.001"
              value={config.sensitivity}
              onChange={(e) => handleSaveConfig({ ...config, sensitivity: parseFloat(e.target.value) })}
              className="w-full h-2 bg-slate-100 rounded-lg appearance-none cursor-pointer accent-blue-600"
            />
            <p className="text-[10px] text-slate-400">数值越小越灵敏。0.5% 表示价格距离关键位 0.5% 时即触发预警。</p>
          </div>

          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <label className="text-sm font-bold text-slate-700 flex items-center">
                <Clock className="w-4 h-4 mr-2 text-blue-500" />
                触发冷却时间
              </label>
              <span className="text-xs font-mono bg-blue-50 text-blue-700 px-2 py-1 rounded">{config.cooldown} 小时</span>
            </div>
            <select 
              value={config.cooldown}
              onChange={(e) => handleSaveConfig({ ...config, cooldown: parseInt(e.target.value) })}
              className="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            >
              <option value={1}>1 小时 (推荐)</option>
              <option value={4}>4 小时 (稳健)</option>
              <option value={24}>24 小时 (长线)</option>
            </select>
            <p className="text-[10px] text-slate-400">同一关键位在冷却时间内不会重复触发，避免频繁骚扰。</p>
          </div>
        </div>
      </div>

      {/* 告警历史列表 */}
      <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden flex flex-col min-h-[500px]">
        <div className="p-6 border-b border-slate-100 flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-slate-800 rounded-xl shadow-lg">
              <History className="w-5 h-5 text-white" />
            </div>
            <h3 className="text-lg font-bold text-slate-800">告警历史记录</h3>
          </div>
          <div className="relative">
            <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
            <input 
              type="text" 
              placeholder="搜索股票代码..."
              value={filterCode}
              onChange={(e) => setFilterCode(e.target.value)}
              className="pl-10 pr-4 py-2 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 w-48"
            />
          </div>
        </div>

        <div className="flex-1 overflow-y-auto">
          {loading ? (
            <div className="flex flex-col items-center justify-center h-64">
              <div className="w-8 h-8 border-4 border-blue-100 border-t-blue-600 rounded-full animate-spin mb-4" />
              <p className="text-sm text-slate-500">正在加载历史数据...</p>
            </div>
          ) : history.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-64 text-slate-400">
              <Bell className="w-12 h-12 mb-4 opacity-20" />
              <p>暂无告警记录</p>
            </div>
          ) : (
            <div className="divide-y divide-slate-50">
              {history
                .filter(item => item && (!filterCode || (item.stockCode && item.stockCode.includes(filterCode))))
                .map((item, idx) => (
                <div key={idx} className="p-4 hover:bg-slate-50 transition-colors group">
                  <div className="flex items-start justify-between">
                    <div className="flex items-start space-x-4">
                      <div className={`mt-1 p-2 rounded-xl ${item.type === 'resistance' ? 'bg-red-50 text-red-600' : 'bg-green-50 text-green-600'}`}>
                        <Bell className="w-4 h-4" />
                      </div>
                      <div>
                        <div className="flex items-center space-x-2 mb-1">
                          <span className="font-bold text-slate-800">{item.stockName || '未知股票'}</span>
                          <span className="text-xs font-mono text-slate-400">{item.stockCode || '000000'}</span>
                          <span className={`text-[10px] px-1.5 py-0.5 rounded ${item.type === 'resistance' ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
                            {item.type === 'resistance' ? '压力位' : '支撑位'}
                          </span>
                        </div>
                        <p className="text-sm text-slate-600 mb-2">{item.message || '无预警消息'}</p>
                        <div className="flex items-center space-x-3">
                          <div className="flex items-center space-x-1 bg-white border border-slate-100 rounded-lg px-2 py-1 shadow-sm">
                            {getRoleIcon(item.role || 'technical')}
                            <span className="text-[10px] font-medium text-slate-500 italic">"{item.advice || '暂无建议'}"</span>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="flex items-center text-[10px] text-slate-400 mb-1">
                        <Calendar className="w-3 h-3 mr-1" />
                        {item.timestamp ? new Date(item.timestamp).toLocaleDateString() : '-'}
                      </div>
                      <div className="flex items-center text-[10px] text-slate-400">
                        <Clock className="w-3 h-3 mr-1" />
                        {item.timestamp ? new Date(item.timestamp).toLocaleTimeString() : '-'}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
