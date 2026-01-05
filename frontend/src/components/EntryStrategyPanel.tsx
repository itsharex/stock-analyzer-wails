import React, { useState } from 'react'
import { Target, Shield, TrendingUp, Info, CheckCircle2, AlertTriangle, Settings2 } from 'lucide-react'
import type { EntryStrategyResult, TrailingStopConfig } from '../types'

interface EntryStrategyPanelProps {
  strategy: EntryStrategyResult
  onConfirm?: (config: TrailingStopConfig) => void
}

const EntryStrategyPanel: React.FC<EntryStrategyPanelProps> = ({ strategy, onConfirm }) => {
  const [trailingConfig, setTrailingConfig] = useState<TrailingStopConfig>({
    enabled: true,
    activationThreshold: 0.05,
    callbackRate: 0.03
  })
  const [showAdvanced, setShowAdvanced] = useState(false)
  const getReasonIcon = (type: string) => {
    switch (type) {
      case 'fundamental': return <Info className="w-4 h-4 text-blue-500" />
      case 'technical': return <TrendingUp className="w-4 h-4 text-purple-500" />
      case 'money_flow': return <Target className="w-4 h-4 text-orange-500" />
      default: return <CheckCircle2 className="w-4 h-4 text-green-500" />
    }
  }

  return (
    <div className="bg-white rounded-xl border border-blue-100 overflow-hidden shadow-sm">
      {/* 头部摘要 */}
      <div className="bg-gradient-to-r from-blue-600 to-indigo-700 p-4 text-white">
        <div className="flex justify-between items-center mb-2">
          <h3 className="text-lg font-bold flex items-center">
            <Target className="w-5 h-5 mr-2" />
            AI 智能建仓方案
          </h3>
          <span className="px-3 py-1 bg-white/20 rounded-full text-xs font-bold backdrop-blur-sm">
            {strategy.recommendation}
          </span>
        </div>
        <p className="text-blue-50 text-sm leading-relaxed">
          {strategy.actionPlan}
        </p>
      </div>

      <div className="p-4 space-y-4">
        {/* 核心参数网格 */}
        <div className="grid grid-cols-2 gap-3">
          <div className="bg-gray-50 p-3 rounded-lg border border-gray-100">
            <div className="text-gray-500 text-[10px] font-bold uppercase mb-1">建议买入区间</div>
            <div className="text-blue-700 font-bold text-lg">{strategy.entryPriceRange}</div>
          </div>
          <div className="bg-gray-50 p-3 rounded-lg border border-gray-100">
            <div className="text-gray-500 text-[10px] font-bold uppercase mb-1">建议首仓比例</div>
            <div className="text-indigo-700 font-bold text-lg">{strategy.initialPosition}</div>
          </div>
          <div className="bg-red-50 p-3 rounded-lg border border-red-100">
            <div className="text-red-500 text-[10px] font-bold uppercase mb-1">止损价 (逻辑失效)</div>
            <div className="text-red-700 font-bold text-lg">¥{strategy.stopLossPrice.toFixed(2)}</div>
          </div>
          <div className="bg-green-50 p-3 rounded-lg border border-green-100">
            <div className="text-green-500 text-[10px] font-bold uppercase mb-1">目标止盈价</div>
            <div className="text-green-700 font-bold text-lg">¥{strategy.takeProfitPrice.toFixed(2)}</div>
          </div>
        </div>

        {/* 核心理由 */}
        <div>
          <h4 className="text-xs font-bold text-gray-400 uppercase mb-2 flex items-center">
            <Shield className="w-3 h-3 mr-1" />
            核心建仓理由与监控阈值
          </h4>
          <div className="space-y-2">
            {strategy.coreReasons.map((reason, idx) => (
              <div key={idx} className="flex items-start p-2 bg-gray-50 rounded-lg border border-gray-100">
                <div className="mt-0.5 mr-2">{getReasonIcon(reason.type)}</div>
                <div className="flex-1">
                  <div className="text-xs font-bold text-gray-700">{reason.description}</div>
                  <div className="text-[10px] text-red-500 mt-1 flex items-center">
                    <AlertTriangle className="w-2.5 h-2.5 mr-1" />
                    失效阈值: {reason.threshold}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* 移动止损高级配置 */}
        <div className="mt-2 pt-4 border-t border-gray-100">
          <button 
            onClick={() => setShowAdvanced(!showAdvanced)}
            className="flex items-center justify-between w-full text-xs font-bold text-gray-400 hover:text-gray-600 transition-colors"
          >
            <div className="flex items-center">
              <Settings2 className="w-3 h-3 mr-1" />
              移动止损策略 (Trailing Stop)
            </div>
            <span>{showAdvanced ? '收起' : '展开配置'}</span>
          </button>
          
          {showAdvanced && (
            <div className="mt-3 p-3 bg-blue-50/50 rounded-lg border border-blue-100 animate-in fade-in slide-in-from-top-1">
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs text-gray-600">启用自动移动止损</span>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input 
                    type="checkbox" 
                    className="sr-only peer" 
                    checked={trailingConfig.enabled}
                    onChange={(e) => setTrailingConfig({...trailingConfig, enabled: e.target.checked})}
                  />
                  <div className="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>
              
              {trailingConfig.enabled && (
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-[10px] text-gray-500 mb-1">启动阈值 (盈利 %)</label>
                    <input 
                      type="number" 
                      step="1"
                      value={trailingConfig.activationThreshold * 100}
                      onChange={(e) => setTrailingConfig({...trailingConfig, activationThreshold: Number(e.target.value) / 100})}
                      className="w-full bg-white border border-blue-200 rounded px-2 py-1 text-xs text-blue-700 focus:outline-none focus:border-blue-400"
                    />
                  </div>
                  <div>
                    <label className="block text-[10px] text-gray-500 mb-1">跟踪回撤 (%)</label>
                    <input 
                      type="number" 
                      step="0.5"
                      value={trailingConfig.callbackRate * 100}
                      onChange={(e) => setTrailingConfig({...trailingConfig, callbackRate: Number(e.target.value) / 100})}
                      className="w-full bg-white border border-blue-200 rounded px-2 py-1 text-xs text-blue-700 focus:outline-none focus:border-blue-400"
                    />
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        {/* 盈亏比与确认按钮 */}
        <div className="pt-2 flex items-center justify-between">
          <div className="text-xs text-gray-500">
            预估盈亏比: <span className="font-bold text-gray-800">{strategy.riskRewardRatio.toFixed(2)}</span>
          </div>
          {onConfirm && (
            <button
              onClick={() => onConfirm(trailingConfig)}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg font-bold text-sm hover:bg-blue-700 transition-all shadow-md shadow-blue-200"
            >
              确认按此方案建仓
            </button>
          )}
        </div>
      </div>
    </div>
  )
}

export default EntryStrategyPanel
