import { useEffect, useMemo, useState, useRef } from 'react'
import { createChart, ColorType, IChartApi } from 'lightweight-charts'
import type { KLineData, StockData } from '../types'
import { useWailsAPI } from '../hooks/useWailsAPI'

interface VolumePriceAnalysisProps {
  stock: StockData
}

function VolumePriceAnalysis({ stock }: VolumePriceAnalysisProps) {
  const { getKLineData } = useWailsAPI()
  const [data, setData] = useState<KLineData[]>([])
  const chartContainerRef = useRef<HTMLDivElement>(null)
  const chartRef = useRef<IChartApi | null>(null)

  const fetchData = async () => {
    const klines = await getKLineData(stock.code, 120, 'daily')
    setData(klines)
  }

  useEffect(() => {
    fetchData()
  }, [stock.code])

  const volMA20 = useMemo(() => {
    const res: number[] = []
    let sum = 0
    const p = 20
    for (let i = 0; i < data.length; i++) {
      sum += data[i].volume
      if (i >= p) sum -= data[i - p].volume
      res.push(i >= p - 1 ? sum / p : NaN)
    }
    return res
  }, [data])

  const obv = useMemo(() => {
    const res: number[] = []
    let cur = 0
    for (let i = 0; i < data.length; i++) {
      if (i === 0) {
        res.push(0)
        continue
      }
      const prevClose = data[i - 1].close
      const c = data[i].close
      if (c > prevClose) cur += data[i].volume
      else if (c < prevClose) cur -= data[i].volume
      res.push(cur)
    }
    return res
  }, [data])

  const analysis = useMemo(() => {
    if (data.length < 30) {
      return { label: '数据不足', detail: '近期数据不足，暂无法给出量价分析', advice: '请稍后再试或换更长周期' }
    }
    const n = 5
    const m = data.length
    let priceChange = 0
    for (let i = m - n; i < m; i++) {
      priceChange += data[i].close - data[i - 1].close
    }
    let recentVol = 0
    for (let i = m - n; i < m; i++) recentVol += data[i].volume
    let prevVol = 0
    for (let i = m - 2 * n; i < m - n; i++) prevVol += data[i].volume
    const volUp = recentVol > prevVol * 1.1
    const volDown = recentVol < prevVol * 0.9
    const priceUp = priceChange > 0
    const priceDown = priceChange < 0
    let label = ''
    let detail = ''
    let advice = ''
    if (priceUp && volUp) {
      label = '价升量增'
      detail = '上涨伴随量能放大，买盘积极，趋势延续的概率较高'
      advice = '顺势为主，关注突破后回踩不破的低风险介入点，止损放在最近低点下方'
    } else if (priceUp && volDown) {
      label = '价升量缩'
      detail = '上涨但量能未跟随，动能偏弱，需警惕假突破'
      advice = '谨慎追高，等待放量确认或在支撑位出现缩量回踩后再考虑介入'
    } else if (priceDown && volUp) {
      label = '价跌量增'
      detail = '下跌伴随放量，抛压显著，风险加大'
      advice = '以防守为主，严格止损，等待企稳与量能回落后再评估'
    } else if (priceDown && volDown) {
      label = '价跌量缩'
      detail = '下跌量能缩小，抛压减弱，更像整理或盘整'
      advice = '耐心等待方向选择，关注缩量止跌后的小阳线与放量突破'
    } else {
      label = '量价平衡'
      detail = '价格与量能变化不明显，市场在观察或换手'
      advice = '保持耐心，结合支撑阻力与资金流再作判断'
    }
    const spike = m >= 1 && volMA20[m - 1] && data[m - 1].volume > (volMA20[m - 1] || 0) * 1.5
    if (spike) {
      advice += '；今日有显著放量，关注是否为突破或异动'
    }
    const obvUp = obv[m - 1] > obv[m - 6]
    if (obvUp && priceUp) {
      advice += '；OBV同步走高，资金配合度较好'
    }
    return { label, detail, advice }
  }, [data, volMA20, obv])

  useEffect(() => {
    if (!chartContainerRef.current) return
    if (data.length === 0) return
    const chart = createChart(chartContainerRef.current, {
      layout: { background: { type: ColorType.Solid, color: 'transparent' }, textColor: '#64748b' },
      grid: { vertLines: { color: '#f1f5f9' }, horzLines: { color: '#f1f5f9' } },
      width: chartContainerRef.current.clientWidth,
      height: 340,
      timeScale: { borderColor: '#f1f5f9', timeVisible: true },
    })
    const candle = chart.addCandlestickSeries({
      upColor: '#ef4444',
      downColor: '#22c55e',
      borderVisible: false,
      wickUpColor: '#ef4444',
      wickDownColor: '#22c55e',
    })
    const volume = chart.addHistogramSeries({ color: '#94a3b8', priceFormat: { type: 'volume' }, priceScaleId: 'volume' })
    chart.priceScale('volume').applyOptions({ scaleMargins: { top: 0.8, bottom: 0 } })
    const volMASeries = chart.addLineSeries({ color: '#3b82f6', lineWidth: 2, priceScaleId: 'volumeMA', title: 'VOL MA20' })
    chart.priceScale('volumeMA').applyOptions({ scaleMargins: { top: 0.8, bottom: 0 } })
    const obvSeries = chart.addLineSeries({ color: '#8b5cf6', lineWidth: 2, priceScaleId: 'obv', title: 'OBV' })
    chart.priceScale('obv').applyOptions({ scaleMargins: { top: 0.6, bottom: 0.2 } })
    const cdata = data.map(d => ({ time: d.time, open: d.open, high: d.high, low: d.low, close: d.close }))
    const vdata = data.map(d => ({ time: d.time, value: d.volume, color: d.close >= d.open ? '#ef444480' : '#22c55e80' }))
    const vma = data.map((d, i) => ({ time: d.time, value: Number.isFinite(volMA20[i]) ? volMA20[i] : 0 }))
    const obvData = data.map((d, i) => ({ time: d.time, value: obv[i] }))
    candle.setData(cdata)
    volume.setData(vdata)
    volMASeries.setData(vma)
    obvSeries.setData(obvData)
    chart.timeScale().fitContent()
    chartRef.current = chart
    const handleResize = () => {
      if (chartContainerRef.current) chart.applyOptions({ width: chartContainerRef.current.clientWidth })
    }
    window.addEventListener('resize', handleResize)
    return () => {
      window.removeEventListener('resize', handleResize)
      chart.remove()
    }
  }, [data, volMA20, obv])

  return (
    <div className="space-y-3">
      <div className="text-sm text-slate-600">
        <div className="font-medium">{analysis.label}</div>
        <div>{analysis.detail}</div>
        <div className="text-slate-700">{analysis.advice}</div>
      </div>
      <div ref={chartContainerRef} className="w-full" />
    </div>
  )
}

export default VolumePriceAnalysis
