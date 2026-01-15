import { useState, useEffect, useRef, useCallback, useMemo } from 'react'
import { createChart, ColorType, IChartApi } from 'lightweight-charts'
import { X, Target, Bell, Plus } from 'lucide-react'
import type { KLineData, TechnicalDrawing, MoneyFlowData } from '../types'
import { useSmartSignals, type SmartSignal, type SmartSignalsConfig } from '../hooks/useSmartSignals'
import { useResizeObserver } from '../hooks/useResizeObserver'

interface SmartSignalsModalProps {
  stock: { code: string; name: string; price: number }
  klines: KLineData[]
  drawings: TechnicalDrawing[]
  moneyFlow: MoneyFlowData[]
  onClose: () => void
  onAddAlert: (signal: SmartSignal) => void
  onCreatePosition: (signal: SmartSignal) => void
}

const DEFAULT_CONFIG: SmartSignalsConfig = {
  volumeMult: 1.2,
  breakPct: 1.5,
  moneyFlowDays: 3,
  obvLookback: 5,
  weights: { volume: 25, breakout: 30, moneyFlow: 25, obv: 20 }
}

export default function SmartSignalsModal({
  stock,
  klines,
  drawings,
  moneyFlow,
  onClose,
  onAddAlert,
  onCreatePosition
}: SmartSignalsModalProps) {
  const chartContainerRef = useRef<HTMLDivElement>(null)
  const chartRef = useRef<IChartApi | null>(null)
  const [config, setConfig] = useState(DEFAULT_CONFIG)
  const [showConfig, setShowConfig] = useState(false)

  // 确保数据按时间升序排列，避免图表渲染和索引计算不一致
  const sortedKlines = useMemo(() => {
    return [...klines].sort((a, b) => a.time.localeCompare(b.time))
  }, [klines])

  const signals = useSmartSignals(sortedKlines, drawings, moneyFlow, config)

  // 图表初始化
  useEffect(() => {
    if (!chartContainerRef.current || sortedKlines.length === 0) return
    const chart = createChart(chartContainerRef.current, {
      layout: { background: { type: ColorType.Solid, color: '#ffffff' }, textColor: '#334155' },
      grid: { vertLines: { color: '#e2e8f0' }, horzLines: { color: '#e2e8f0' } },
      width: chartContainerRef.current.clientWidth,
      height: chartContainerRef.current.clientHeight,
      timeScale: { borderColor: '#e2e8f0', timeVisible: true }
    })

    const candle = chart.addCandlestickSeries({
      upColor: '#ef4444',
      downColor: '#22c55e',
      borderVisible: false,
      wickUpColor: '#ef4444',
      wickDownColor: '#22c55e'
    })

    const volume = chart.addHistogramSeries({
      color: '#94a3b8',
      priceFormat: { type: 'volume' },
      priceScaleId: 'volume'
    })
    chart.priceScale('volume').applyOptions({ scaleMargins: { top: 0.8, bottom: 0 } })

    const data = sortedKlines.map(k => ({ time: k.time, open: k.open, high: k.high, low: k.low, close: k.close }))
    const volData = sortedKlines.map(k => ({ time: k.time, value: k.volume, color: k.close >= k.open ? '#ef444480' : '#22c55e80' }))
    candle.setData(data)
    volume.setData(volData)

    // 信号标记
    const markers = signals
      .map(s => ({
        time: s.time,
        position: (s.type === 'buy' ? 'belowBar' : 'aboveBar') as any,
        color: s.type === 'buy' ? '#16a34a' : '#dc2626',
        shape: (s.type === 'buy' ? 'arrowUp' : 'arrowDown') as any,
        text: `${s.type === 'buy' ? 'B' : 'S'}${s.score}`
      }))
      .sort((a, b) => a.time.localeCompare(b.time))
    
    candle.setMarkers(markers)

    chart.timeScale().fitContent()
    chartRef.current = chart

    return () => chart.remove()
  }, [sortedKlines, signals])

  // 容器尺寸监听
  useResizeObserver(chartContainerRef, (w, h) => {
    chartRef.current?.applyOptions({ width: w, height: h })
  })

  // 点击信号定位
  const handleSignalClick = useCallback((s: SmartSignal) => {
    if (!chartRef.current) return
    
    // Find index of the signal time in klines
    // 使用 sortedKlines 确保索引与图表数据一致
    const index = sortedKlines.findIndex(k => k.time === s.time)
    if (index === -1) {
      console.warn('[SmartSignalsModal] Signal time not found in klines:', s.time)
      return
    }

    // 0 = oldest bar in array (if sorted asc), but Lightweight charts usually treats rightmost as newest.
    // klines are sorted ASC (old -> new).
    // index 0 is oldest. index length-1 is newest.
    // scrollToPosition(0) scrolls to the NEWEST bar (rightmost).
    // scrollToPosition(pos > 0) scrolls to the left (history).
    // We want to scroll to `index`.
    // Distance from newest (rightmost) is: (length - 1) - index.
    const distFromNewest = sortedKlines.length - 1 - index

    const visibleLogicalRange = chartRef.current.timeScale().getVisibleLogicalRange()
    const visibleBars = visibleLogicalRange ? visibleLogicalRange.to - visibleLogicalRange.from : 50
    
    // Scroll to center the signal
    // position = distance from right edge.
    // We want the target bar (at distFromNewest) to be at the center (visibleBars / 2).
    // So the Right Edge should be at (distFromNewest - visibleBars / 2).
    const position = distFromNewest - (visibleBars / 2)
    
    // Ensure we don't scroll into weird negative space if not needed, though negative is valid (future)
    chartRef.current.timeScale().scrollToPosition(position, true)
  }, [sortedKlines])

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-xl shadow-2xl w-[90vw] h-[90vh] max-w-7xl max-h-[900px] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-slate-200">
          <div className="flex items-center gap-3">
            <Target className="w-6 h-6 text-indigo-600" />
            <div>
              <h2 className="text-xl font-bold text-slate-800">AI 智能信号</h2>
              <p className="text-sm text-slate-500">{stock.name} ({stock.code})</p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={() => setShowConfig(!showConfig)}
              className="px-3 py-2 text-sm font-medium text-slate-600 hover:text-slate-800 border border-slate-300 rounded-lg hover:bg-slate-50"
            >
              配置规则
            </button>
            <button onClick={onClose} className="p-2 text-slate-400 hover:text-slate-600 rounded-lg hover:bg-slate-100">
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Config Panel */}
        {showConfig && (
          <div className="px-6 py-4 bg-slate-50 border-b border-slate-200 grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">放量倍数</label>
              <input
                type="number"
                step={0.1}
                min={1}
                value={config.volumeMult}
                onChange={e => setConfig({ ...config, volumeMult: parseFloat(e.target.value) })}
                className="w-full px-2 py-1 text-sm border border-slate-300 rounded"
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">突破幅度(%)</label>
              <input
                type="number"
                step={0.1}
                min={0.5}
                value={config.breakPct}
                onChange={e => setConfig({ ...config, breakPct: parseFloat(e.target.value) })}
                className="w-full px-2 py-1 text-sm border border-slate-300 rounded"
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">资金流天数</label>
              <input
                type="number"
                min={1}
                max={10}
                value={config.moneyFlowDays}
                onChange={e => setConfig({ ...config, moneyFlowDays: parseInt(e.target.value) })}
                className="w-full px-2 py-1 text-sm border border-slate-300 rounded"
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">OBV 回溯</label>
              <input
                type="number"
                min={3}
                max={20}
                value={config.obvLookback}
                onChange={e => setConfig({ ...config, obvLookback: parseInt(e.target.value) })}
                className="w-full px-2 py-1 text-sm border border-slate-300 rounded"
              />
            </div>
          </div>
        )}

        {/* Body */}
        <div className="flex-1 flex overflow-hidden">
          {/* Left - Chart */}
          <div className="flex-1 p-4">
            <div ref={chartContainerRef} className="w-full h-full rounded-lg border border-slate-200" />
          </div>

          {/* Right - Signals List */}
          <div className="w-[380px] border-l border-slate-200 flex flex-col">
            <div className="px-4 py-3 border-b border-slate-200">
              <div className="flex items-center justify-between">
                <h3 className="font-semibold text-slate-800">信号列表</h3>
                <span className="text-xs text-slate-500">{signals.length} 条</span>
              </div>
            </div>
            <div className="flex-1 overflow-y-auto px-3 py-2 space-y-2">
              {signals.length === 0 ? (
                <div className="text-center text-slate-500 text-sm py-8">暂无信号</div>
              ) : (
                signals.map((s, idx) => (
                  <div key={idx} onClick={() => handleSignalClick(s)} className="bg-slate-50 rounded-lg p-3 border border-slate-200 hover:bg-slate-100 transition-colors cursor-pointer">
                    <div className="flex items-center justify-between mb-2">
                      <div className="flex items-center gap-2">
                        <div className={`w-6 h-6 rounded-full flex items-center justify-center text-white text-xs font-bold ${s.type === 'buy' ? 'bg-green-600' : 'bg-red-600'}`}>
                          {s.type === 'buy' ? 'B' : 'S'}
                        </div>
                        <div>
                          <div className="text-sm font-semibold text-slate-800">{s.price.toFixed(2)}</div>
                          <div className="text-xs text-slate-500">{s.time}</div>
                        </div>
                      </div>
                      <div className="text-xs font-medium text-slate-600">{s.score}分</div>
                    </div>
                    <div className="text-xs text-slate-600 mb-2">{s.aiReason}</div>
                    <div className="flex gap-2">
                      <button
                        onClick={() => onAddAlert(s)}
                        className="flex-1 flex items-center justify-center gap-1 px-2 py-1 text-xs font-medium text-blue-700 bg-blue-100 hover:bg-blue-200 rounded"
                      >
                        <Bell className="w-3 h-3" /> 预警
                      </button>
                      <button
                        onClick={() => onCreatePosition(s)}
                        className="flex-1 flex items-center justify-center gap-1 px-2 py-1 text-xs font-medium text-green-700 bg-green-100 hover:bg-green-200 rounded"
                      >
                        <Plus className="w-3 h-3" /> 建仓
                      </button>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}