import { useState, useEffect, useCallback } from 'react'
import SmartSignalsModal from '../components/SmartSignalsModal'
import type { SmartSignal } from '../hooks/useSmartSignals'
import { useWailsAPI } from './useWailsAPI'

export function useSmartSignalsModal(stock: { code: string; name: string; price: number }) {
  const [open, setOpen] = useState(false)
  const [klines, setKlines] = useState<any[]>([])
  const [drawings, setDrawings] = useState<any[]>([])
  const [moneyFlow, setMoneyFlow] = useState<any[]>([])
  const [loading, setLoading] = useState(false)

  const { getKLineData, getMoneyFlowData } = useWailsAPI()

  const loadData = useCallback(async () => {
    if (!open) return
    setLoading(true)
    try {
      const [k, mf] = await Promise.all([
        getKLineData(stock.code, 120, 'daily'),
        getMoneyFlowData(stock.code)
      ])
      // Ensure klines are sorted ascending by time
      if (Array.isArray(k) && k.length > 0) {
        console.log('SmartSignals: KLines before sort (first 2):', k.slice(0, 2))
      }
      const sortedK = Array.isArray(k) ? [...k].sort((a: any, b: any) => {
        if (!a.time || !b.time) return 0
        return a.time.localeCompare(b.time)
      }) : []
      if (sortedK.length > 0) {
        console.log('SmartSignals: KLines after sort (first 2):', sortedK.slice(0, 2))
      }
      setKlines(sortedK)
      setMoneyFlow(mf.data || [])
      // 这里可以集成 AI drawings，先留空
      setDrawings([])
    } catch (e) {
      console.error(e)
    } finally {
      setLoading(false)
    }
  }, [open, stock.code, getKLineData, getMoneyFlowData])

  useEffect(() => { loadData() }, [loadData])

  const show = useCallback(() => setOpen(true), [])
  const hide = useCallback(() => setOpen(false), [])

  const onAddAlert = useCallback(async (signal: SmartSignal) => {
    try {
      await (window as any).go.main.App.PriceAlertCreateAlert(JSON.stringify({
        stockCode: stock.code,
        stockName: stock.name,
        type: signal.type === 'buy' ? 'above' : 'below',
        price: signal.price,
        note: `智能信号：${signal.aiReason}`
      }))
      alert('预警已创建')
    } catch (e: any) {
      alert('创建预警失败：' + e.message)
    }
  }, [stock.code, stock.name])

  const onCreatePosition = useCallback(async (signal: SmartSignal) => {
    try {
      const entryStrategy = {
        entryPrice: signal.price,
        stopLossPrice: signal.type === 'buy' ? signal.price * 0.97 : signal.price * 1.03,
        takeProfitPrice: signal.type === 'buy' ? signal.price * 1.08 : signal.price * 0.92,
        coreReasons: [{ type: 'technical' as const, description: signal.aiReason }]
      }
      await (window as any).go.main.App.AddPosition({
        stockCode: stock.code,
        stockName: stock.name,
        entryPrice: signal.price,
        entryTime: new Date().toISOString(),
        strategy: entryStrategy,
        currentStatus: 'holding',
        logicStatus: 'valid'
      })
      alert('建仓已提交')
    } catch (e: any) {
      alert('建仓失败：' + e.message)
    }
  }, [stock.code, stock.name])

  return {
    open,
    show,
    hide,
    klines,
    drawings,
    moneyFlow,
    loading,
    onAddAlert,
    onCreatePosition,
    SmartSignalsModal: open ? SmartSignalsModal : null
  }
}