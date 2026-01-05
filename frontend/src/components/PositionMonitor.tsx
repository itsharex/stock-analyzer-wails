import React, { useEffect, useState } from 'react'
import { Shield, AlertTriangle, CheckCircle, TrendingDown, Activity } from 'lucide-react'

interface Position {
  stockCode: string
  stockName: string
  entryPrice: number
  entryTime: string
  logicStatus: 'valid' | 'violated' | 'warning'
  strategy: {
    stopLossPrice: number
    takeProfitPrice: number
    coreReasons: Array<{
      type: string
      description: string
    }>
  }
}

interface PositionMonitorProps {
  positions: Record<string, Position>
  onRefresh: () => void
}

const PositionMonitor: React.FC<PositionMonitorProps> = ({ positions, onRefresh }) => {
  const [violations, setViolations] = useState<any[]>([])

  useEffect(() => {
    // 监听来自后端的逻辑失效事件
    // @ts-ignore
    const unsubscribe = window.runtime.EventsOn('logic_violation', (data: any) => {
      setViolations(prev => [data, ...prev].slice(0, 5))
      onRefresh() // 刷新持仓列表状态
    })

    return () => unsubscribe()
  }, [onRefresh])

  const activePositions = Object.values(positions)

  if (activePositions.length === 0) return null

  return (
    <div className="space-y-4">
      {/* 实时预警横幅 */}
      {violations.length > 0 && (
        <div className="animate-pulse">
          {violations.map((v, i) => (
            <div key={i} className="bg-red-500/10 border border-red-500/50 rounded-lg p-3 mb-2 flex items-start space-x-3">
              <AlertTriangle className="w-5 h-5 text-red-500 shrink-0 mt-0.5" />
              <div>
                <div className="text-red-500 font-bold text-sm">逻辑失效预警: {v.name} ({v.code})</div>
                <div className="text-red-400/80 text-xs mt-1">
                  {v.reasons.map((r: string, idx: number) => (
                    <div key={idx}>• {r}</div>
                  ))}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* 持仓逻辑看板 */}
      <div className="bg-slate-900/50 border border-slate-800 rounded-xl overflow-hidden">
        <div className="px-4 py-3 border-b border-slate-800 flex items-center justify-between bg-slate-800/30">
          <div className="flex items-center space-x-2">
            <Shield className="w-4 h-4 text-indigo-400" />
            <span className="text-sm font-bold text-slate-200">AI 持仓逻辑监控</span>
          </div>
          <span className="text-[10px] text-slate-500 uppercase tracking-wider">实时扫描中</span>
        </div>

        <div className="divide-y divide-slate-800">
          {activePositions.map((pos) => (
            <div key={pos.stockCode} className="p-4 hover:bg-slate-800/20 transition-colors">
              <div className="flex items-center justify-between mb-3">
                <div>
                  <div className="text-sm font-bold text-white">{pos.stockName}</div>
                  <div className="text-[10px] text-slate-500">{pos.stockCode}</div>
                </div>
                <div className={`flex items-center space-x-1.5 px-2 py-1 rounded-full text-[10px] font-bold ${
                  pos.logicStatus === 'valid' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-red-500/10 text-red-500'
                }`}>
                  {pos.logicStatus === 'valid' ? (
                    <><CheckCircle className="w-3 h-3" /><span>逻辑有效</span></>
                  ) : (
                    <><AlertTriangle className="w-3 h-3" /><span>逻辑失效</span></>
                  )}
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-3">
                <div className="bg-slate-800/40 rounded-lg p-2">
                  <div className="text-[10px] text-slate-500 mb-1">建仓价 / 止损价</div>
                  <div className="text-xs font-mono text-slate-300">
                    {pos.entryPrice.toFixed(2)} / <span className="text-red-400">{pos.strategy.stopLossPrice.toFixed(2)}</span>
                  </div>
                </div>
                <div className="bg-slate-800/40 rounded-lg p-2">
                  <div className="text-[10px] text-slate-500 mb-1">监控理由数</div>
                  <div className="text-xs font-mono text-slate-300 flex items-center space-x-1">
                    <Activity className="w-3 h-3 text-indigo-400" />
                    <span>{pos.strategy.coreReasons.length} 条核心逻辑</span>
                  </div>
                </div>
              </div>

              {pos.logicStatus === 'violated' && (
                <div className="mt-2 p-2 bg-red-500/5 border border-red-500/20 rounded-lg">
                  <div className="text-[10px] text-red-400 font-bold mb-1 flex items-center space-x-1">
                    <TrendingDown className="w-3 h-3" />
                    <span>建议操作: 立即评估减仓或止损</span>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export default PositionMonitor
