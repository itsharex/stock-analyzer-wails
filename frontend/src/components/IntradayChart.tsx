import { useEffect, useRef } from 'react'
import { createChart, ColorType, IChartApi, LineStyle } from 'lightweight-charts'
import type { IntradayData, MoneyFlowData } from '../types'

interface IntradayChartProps {
  data: IntradayData[]
  moneyFlowData?: MoneyFlowData[]
  preClose: number
  height?: number
}

function IntradayChart({ data, moneyFlowData, preClose, height = 400 }: IntradayChartProps) {
  const chartContainerRef = useRef<HTMLDivElement>(null)
  const chartRef = useRef<IChartApi | null>(null)

  const toUnixTimestampSeconds = (timeStr: string): number | null => {
    if (!timeStr) return null
    // 1) Full date/datetime like "2026-01-06 14:59" / "2026-01-06T14:59:00"
    if (timeStr.includes('-')) {
      const ms = new Date(timeStr.replace(/-/g, '/')).getTime()
      return Number.isFinite(ms) ? Math.floor(ms / 1000) : null
    }
    // 2) Intraday time like "14:59" / "14:59:30"
    const m = timeStr.match(/^(\d{1,2}):(\d{2})(?::(\d{2}))?$/)
    if (m) {
      const hh = Number(m[1])
      const mm = Number(m[2])
      const ss = m[3] ? Number(m[3]) : 0
      if (![hh, mm, ss].every(Number.isFinite)) return null
      if (hh < 0 || hh > 23 || mm < 0 || mm > 59 || ss < 0 || ss > 59) return null
      const base = new Date()
      base.setHours(0, 0, 0, 0)
      const ms = base.getTime() + (hh * 3600 + mm * 60 + ss) * 1000
      return Math.floor(ms / 1000)
    }
    return null
  }

  useEffect(() => {
    if (!chartContainerRef.current || data.length === 0) return

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
        secondsVisible: false,
      },
      rightPriceScale: {
        borderColor: '#f1f5f9',
        autoScale: true,
      },
    })

    // 1. 价格线
    const priceSeries = chart.addLineSeries({
      color: '#3b82f6',
      lineWidth: 2,
      priceFormat: {
        type: 'price',
        precision: 2,
        minMove: 0.01,
      },
    })

    // 2. 均价线
    const avgPriceSeries = chart.addLineSeries({
      color: '#f59e0b',
      lineWidth: 1,
      lineStyle: LineStyle.Solid,
      priceFormat: {
        type: 'price',
        precision: 2,
        minMove: 0.01,
      },
    })

    // 3. 成交量
    const volumeSeries = chart.addHistogramSeries({
      color: '#94a3b8',
      priceFormat: { type: 'volume' },
      priceScaleId: 'volume',
    })

    chart.priceScale('volume').applyOptions({
      scaleMargins: {
        top: 0.8,
        bottom: 0,
      },
    })

    // 4. 昨收参考线
    priceSeries.createPriceLine({
      price: preClose,
      color: '#94a3b8',
      lineWidth: 1,
      lineStyle: LineStyle.Dashed,
      axisLabelVisible: true,
      title: '昨收',
    })

    // 准备数据
    const formattedPriceData = data
      .map(d => {
        const timestamp = toUnixTimestampSeconds(d.time)
        if (timestamp == null) return null
        return { time: timestamp as any, value: d.price }
      })
      .filter(Boolean) as any[]

    const formattedAvgPriceData = data
      .map(d => {
        const timestamp = toUnixTimestampSeconds(d.time)
        if (timestamp == null) return null
        return { time: timestamp as any, value: d.avgPrice }
      })
      .filter(Boolean) as any[]

    const formattedVolumeData = data
      .map((d, i) => {
        const timestamp = toUnixTimestampSeconds(d.time)
        if (timestamp == null) return null
        const color =
          i === 0
            ? d.price >= preClose
              ? '#ef444480'
              : '#22c55e80'
            : d.price >= data[i - 1].price
              ? '#ef444480'
              : '#22c55e80'

        return { time: timestamp as any, value: d.volume, color }
      })
      .filter(Boolean) as any[]

    priceSeries.setData(formattedPriceData)
    avgPriceSeries.setData(formattedAvgPriceData)
    volumeSeries.setData(formattedVolumeData)

    // 添加异动标记 (Markers)
    if (moneyFlowData && moneyFlowData.length > 0) {
      const markers = moneyFlowData
        .filter(d => d.signal === '扫货' || d.signal === '砸盘')
        .map(d => {
          const timestamp = toUnixTimestampSeconds(d.time)
          if (timestamp == null) return null
          return {
            time: timestamp as any,
            position: d.signal === '扫货' ? 'belowBar' : 'aboveBar' as any,
            color: d.signal === '扫货' ? '#ef4444' : '#22c55e',
            shape: d.signal === '扫货' ? 'arrowUp' : 'arrowDown' as any,
            text: d.signal,
            size: 1,
          }
        })
        .filter(Boolean) as any[]
      priceSeries.setMarkers(markers)
    }

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
  }, [data, moneyFlowData, preClose, height])

  return (
    <div className="relative">
      <div ref={chartContainerRef} className="w-full" />
      <div className="absolute top-2 left-2 flex gap-4 text-xs font-medium">
        <div className="flex items-center gap-1">
          <div className="w-3 h-0.5 bg-blue-500"></div>
          <span className="text-slate-600">价格</span>
        </div>
        <div className="flex items-center gap-1">
          <div className="w-3 h-0.5 bg-amber-500"></div>
          <span className="text-slate-600">均价</span>
        </div>
      </div>
    </div>
  )
}

export default IntradayChart
