import { useEffect, useRef } from 'react'
import { createChart, ColorType, IChartApi } from 'lightweight-charts'
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
    })

    // 2. 成交量 (通过 priceScaleId 模拟副图)
    const volumeSeries = chart.addHistogramSeries({
      color: '#94a3b8',
      priceFormat: { type: 'volume' },
      priceScaleId: 'volume', // 使用独立的 price scale
    })

    chart.priceScale('volume').applyOptions({
      scaleMargins: {
        top: 0.8, // 放在底部 20% 区域
        bottom: 0,
      },
    })

    // 3. MACD (通过 priceScaleId 模拟副图)
    if (showMACD) {
      const difSeries = chart.addLineSeries({ color: '#2196F3', lineWidth: 1, priceScaleId: 'macd', title: 'DIF' })
      const deaSeries = chart.addLineSeries({ color: '#FF9800', lineWidth: 1, priceScaleId: 'macd', title: 'DEA' })
      const barSeries = chart.addHistogramSeries({ priceScaleId: 'macd', title: 'MACD' })
      
      chart.priceScale('macd').applyOptions({
        scaleMargins: { top: 0.6, bottom: 0.2 },
      })

      difSeries.setData(data.map(d => ({ time: d.time, value: d.macd?.dif || 0 })))
      deaSeries.setData(data.map(d => ({ time: d.time, value: d.macd?.dea || 0 })))
      barSeries.setData(data.map(d => ({ 
        time: d.time, 
        value: d.macd?.bar || 0,
        color: (d.macd?.bar || 0) >= 0 ? '#ef444480' : '#22c55e80'
      })))
    }

    // 4. KDJ (通过 priceScaleId 模拟副图)
    if (showKDJ) {
      const kSeries = chart.addLineSeries({ color: '#9C27B0', lineWidth: 1, priceScaleId: 'kdj', title: 'K' })
      const dSeries = chart.addLineSeries({ color: '#FFEB3B', lineWidth: 1, priceScaleId: 'kdj', title: 'D' })
      const jSeries = chart.addLineSeries({ color: '#E91E63', lineWidth: 1, priceScaleId: 'kdj', title: 'J' })
      
      chart.priceScale('kdj').applyOptions({
        scaleMargins: { top: 0.7, bottom: 0.1 },
      })

      kSeries.setData(data.map(d => ({ time: d.time, value: d.kdj?.k || 0 })))
      dSeries.setData(data.map(d => ({ time: d.time, value: d.kdj?.d || 0 })))
      jSeries.setData(data.map(d => ({ time: d.time, value: d.kdj?.j || 0 })))
    }

    // 5. RSI (通过 priceScaleId 模拟副图)
    if (showRSI) {
      const rsiSeries = chart.addLineSeries({ color: '#00BCD4', lineWidth: 1, priceScaleId: 'rsi', title: 'RSI' })
      
      chart.priceScale('rsi').applyOptions({
        scaleMargins: { top: 0.8, bottom: 0.05 },
      })

      rsiSeries.setData(data.map(d => ({ time: d.time, value: d.rsi || 0 })))
    }

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
