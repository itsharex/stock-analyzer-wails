import { useEffect, useRef } from 'react'
import { createChart, ColorType, IChartApi } from 'lightweight-charts'
import type { MoneyFlowData } from '../types'

interface MoneyFlowChartProps {
  data: MoneyFlowData[]
  height?: number
}

function MoneyFlowChart({ data, height = 200 }: MoneyFlowChartProps) {
  const chartContainerRef = useRef<HTMLDivElement>(null)
  const chartRef = useRef<IChartApi | null>(null)

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
    })

    // 主力净流入 (特大单 + 大单)
    const mainSeries = chart.addHistogramSeries({
      color: '#ef4444',
      priceFormat: { type: 'volume' },
      title: '主力净流入',
    })

    // 散户净流入 (小单)
    const retailSeries = chart.addHistogramSeries({
      color: '#22c55e',
      priceFormat: { type: 'volume' },
      title: '散户净流入',
    })

    // 准备数据
    const formattedMainData = data.map(d => {
      const timestamp = Math.floor(new Date(d.time.replace(/-/g, '/')).getTime() / 1000)
      return {
        time: timestamp as any,
        value: d.mainNet,
        color: d.mainNet >= 0 ? '#ef444480' : '#22c55e80',
      }
    })

    const formattedRetailData = data.map(d => {
      const timestamp = Math.floor(new Date(d.time.replace(/-/g, '/')).getTime() / 1000)
      return {
        time: timestamp as any,
        value: d.small,
        color: d.small >= 0 ? '#ef444440' : '#22c55e40',
      }
    })

    mainSeries.setData(formattedMainData)
    retailSeries.setData(formattedRetailData)

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
  }, [data, height])

  return (
    <div className="relative">
      <div ref={chartContainerRef} className="w-full" />
      <div className="absolute top-2 left-2 flex gap-4 text-[10px] font-bold uppercase tracking-wider">
        <div className="flex items-center gap-1">
          <div className="w-2 h-2 bg-red-500/50 rounded-sm"></div>
          <span className="text-slate-500">主力净流入</span>
        </div>
        <div className="flex items-center gap-1">
          <div className="w-2 h-2 bg-green-500/50 rounded-sm"></div>
          <span className="text-slate-500">散户净流入</span>
        </div>
      </div>
    </div>
  )
}

export default MoneyFlowChart
