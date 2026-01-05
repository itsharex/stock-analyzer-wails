import { useEffect, useRef } from 'react'
import { createChart, ColorType, IChartApi, PaneSize } from 'lightweight-charts'
import type { KLineData } from '../types'

interface KLineChartProps {
  data: KLineData[]
  height?: number
  showMACD?: boolean
  showKDJ?: boolean
  showRSI?: boolean
}

function KLineChart({ data, height = 600, showMACD = false, showKDJ = false, showRSI = false }: KLineChartProps) {
  const chartContainerRef = useRef<HTMLDivElement>(null)
  const chartRef = useRef<IChartApi | null>(null)

  useEffect(() => {
    if (!chartContainerRef.current) return

    const chart = createChart(chartContainerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: 'transparent' },
        textColor: '#64748b',
      },
      grid: {
        vertLines: { color: '#f1f5f9' },
        horzLines: { color: '#f1f5f9' },
      },
      width: chartContainerRef.current.clientWidth,
      height: height,
      timeScale: {
        borderColor: '#f1f5f9',
        timeVisible: true,
      },
    })

    // 1. 主图：K线
    const candlestickSeries = chart.addCandlestickSeries({
      upColor: '#ef4444',
      downColor: '#22c55e',
      borderVisible: false,
      wickUpColor: '#ef4444',
      wickDownColor: '#22c55e',
      pane: 0,
    })

    // 2. 成交量 (Pane 1)
    const volumeSeries = chart.addHistogramSeries({
      color: '#94a3b8',
      priceFormat: { type: 'volume' },
      pane: 1,
    })

    // 3. MACD (Pane 2)
    let macdSeries: any = null
    if (showMACD) {
      const difSeries = chart.addLineSeries({ color: '#2196F3', lineWidth: 1, pane: 2, title: 'DIF' })
      const deaSeries = chart.addLineSeries({ color: '#FF9800', lineWidth: 1, pane: 2, title: 'DEA' })
      const barSeries = chart.addHistogramSeries({ pane: 2, title: 'MACD' })
      
      difSeries.setData(data.map(d => ({ time: d.time, value: d.macd?.dif || 0 })))
      deaSeries.setData(data.map(d => ({ time: d.time, value: d.macd?.dea || 0 })))
      barSeries.setData(data.map(d => ({ 
        time: d.time, 
        value: d.macd?.bar || 0,
        color: (d.macd?.bar || 0) >= 0 ? '#ef444480' : '#22c55e80'
      })))
      macdSeries = { difSeries, deaSeries, barSeries }
    }

    // 4. KDJ (Pane 3)
    if (showKDJ) {
      const kSeries = chart.addLineSeries({ color: '#9C27B0', lineWidth: 1, pane: 3, title: 'K' })
      const dSeries = chart.addLineSeries({ color: '#FFEB3B', lineWidth: 1, pane: 3, title: 'D' })
      const jSeries = chart.addLineSeries({ color: '#E91E63', lineWidth: 1, pane: 3, title: 'J' })
      
      kSeries.setData(data.map(d => ({ time: d.time, value: d.kdj?.k || 0 })))
      dSeries.setData(data.map(d => ({ time: d.time, value: d.kdj?.d || 0 })))
      jSeries.setData(data.map(d => ({ time: d.time, value: d.kdj?.j || 0 })))
    }

    // 5. RSI (Pane 4)
    if (showRSI) {
      const rsiSeries = chart.addLineSeries({ color: '#00BCD4', lineWidth: 1, pane: 4, title: 'RSI' })
      rsiSeries.setData(data.map(d => ({ time: d.time, value: d.rsi || 0 })))
    }

    // 设置 Pane 比例
    const panesCount = 2 + (showMACD ? 1 : 0) + (showKDJ ? 1 : 0) + (showRSI ? 1 : 0)
    // 简单分配比例：主图占 50%，其余平分
    // 注意：lightweight-charts 的 pane 比例设置较为复杂，这里通过 height 间接控制或使用默认分配

    const formattedData = data.map(d => ({
      time: d.time,
      open: d.open,
      high: d.high,
      low: d.low,
      close: d.close,
    }))

    const volumeData = data.map(d => ({
      time: d.time,
      value: d.volume,
      color: d.close >= d.open ? '#ef444480' : '#22c55e80',
    }))

    candlestickSeries.setData(formattedData)
    volumeSeries.setData(volumeData)

    chart.timeScale().fitContent()
    chartRef.current = chart

    const handleResize = () => {
      if (chartContainerRef.current) {
        chart.applyOptions({ width: chartContainerRef.current.clientWidth })
      }
    }

    window.addEventListener('resize', handleResize)

    return () => {
      window.removeEventListener('resize', handleResize)
      chart.remove()
    }
  }, [data, height, showMACD, showKDJ, showRSI])

  return <div ref={chartContainerRef} className="w-full" />
}

export default KLineChart
