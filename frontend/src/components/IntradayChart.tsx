import { useEffect, useRef } from 'react'
import { createChart, ColorType, IChartApi, LineStyle } from 'lightweight-charts'
import type { IntradayData } from '../types'

interface IntradayChartProps {
  data: IntradayData[]
  preClose: number
  height?: number
}

function IntradayChart({ data, preClose, height = 400 }: IntradayChartProps) {
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
    // 东方财富分时数据时间格式通常是 "YYYY-MM-DD HH:mm"
    // lightweight-charts 需要 unix timestamp 或 "YYYY-MM-DD"
    // 对于分时图，我们需要将其转换为 timestamp
    const formattedPriceData = data.map(d => {
      const timestamp = Math.floor(new Date(d.time.replace(/-/g, '/')).getTime() / 1000)
      return {
        time: timestamp as any,
        value: d.price,
      }
    })

    const formattedAvgPriceData = data.map(d => {
      const timestamp = Math.floor(new Date(d.time.replace(/-/g, '/')).getTime() / 1000)
      return {
        time: timestamp as any,
        value: d.avgPrice,
      }
    })

    const formattedVolumeData = data.map((d, i) => {
      const timestamp = Math.floor(new Date(d.time.replace(/-/g, '/')).getTime() / 1000)
      // 分时图成交量通常是累计的，或者是每分钟的。东方财富 trends2 接口返回的是每分钟成交量
      const color = i === 0 
        ? (d.price >= preClose ? '#ef444480' : '#22c55e80')
        : (d.price >= data[i-1].price ? '#ef444480' : '#22c55e80')
      
      return {
        time: timestamp as any,
        value: d.volume,
        color: color,
      }
    })

    priceSeries.setData(formattedPriceData)
    avgPriceSeries.setData(formattedAvgPriceData)
    volumeSeries.setData(formattedVolumeData)

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
  }, [data, preClose, height])

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
