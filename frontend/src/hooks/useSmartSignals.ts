import { useMemo } from 'react'
import type { KLineData, TechnicalDrawing, MoneyFlowData } from '../types'

export interface SmartSignal {
  time: string
  price: number
  type: 'buy' | 'sell'
  score: number
  reasons: string[]
  aiReason: string
}

export interface SmartSignalsConfig {
  volumeMult: number        // 放量倍数 vs 均量
  breakPct: number          // 突破/跌破关键位幅度 %
  moneyFlowDays: number     // 资金流方向观察天数
  obvLookback: number       // OBV 方向过滤回溯
  weights: {
    volume: number
    breakout: number
    moneyFlow: number
    obv: number
  }
}

const DEFAULT_CONFIG: SmartSignalsConfig = {
  volumeMult: 1.2,
  breakPct: 1.5,
  moneyFlowDays: 3,
  obvLookback: 5,
  weights: { volume: 25, breakout: 30, moneyFlow: 25, obv: 20 }
}

export function useSmartSignals(
  klines: KLineData[],
  drawings: TechnicalDrawing[],
  moneyFlow: MoneyFlowData[],
  config = DEFAULT_CONFIG
): SmartSignal[] {
  return useMemo(() => {
    const n = klines.length
    if (n < 30) return []

    // 1. 计算均量
    const volMA20: number[] = []
    let sum = 0
    for (let i = 0; i < n; i++) {
      sum += klines[i].volume
      if (i >= 20) sum -= klines[i - 20].volume
      volMA20.push(i >= 19 ? sum / 20 : NaN)
    }

    // 2. OBV
    const obv: number[] = []
    let cur = 0
    for (let i = 0; i < n; i++) {
      if (i === 0) { obv.push(0); continue }
      const prev = klines[i - 1].close
      const now = klines[i].close
      if (now > prev) cur += klines[i].volume
      else if (now < prev) cur -= klines[i].volume
      obv.push(cur)
    }

    // 3. 资金流方向
    const mfDir: boolean[] = []
    for (let i = 0; i < n; i++) {
      const idx = moneyFlow.findIndex(m => m.time === klines[i].time)
      if (idx === -1) { mfDir.push(false); continue }
      mfDir.push(moneyFlow[idx].mainNet > 0)
    }

    // 4. 关键位映射
    const supports = drawings.filter(d => d.type === 'support').map(d => d.price || 0)
    const resistances = drawings.filter(d => d.type === 'resistance').map(d => d.price || 0)

    const signals: SmartSignal[] = []

    // 从第 20 根开始评估
    for (let i = 20; i < n; i++) {
      const k = klines[i]
      const volScore = Math.min(100, (k.volume / (volMA20[i] || 1) - 1) * 200)
      const priceChg = (k.close - klines[i - 1].close) / klines[i - 1].close * 100

      // 突破/跌破关键位
      let breakScore = 0
      let breakReason = ''
      for (const r of resistances) {
        if (k.close > r * (1 + config.breakPct / 100) && klines[i - 1].close <= r) {
          breakScore = 100
          breakReason = `突破阻力 ${r.toFixed(2)}`
          break
        }
      }
      for (const s of supports) {
        if (k.close < s * (1 - config.breakPct / 100) && klines[i - 1].close >= s) {
          breakScore = 100
          breakReason = `跌破支撑 ${s.toFixed(2)}`
          break
        }
      }

      // 资金流方向一致性
      let mfScore = 0
      let mfReason = ''
      const recentMF = mfDir.slice(Math.max(0, i - config.moneyFlowDays + 1), i + 1)
      const upCount = recentMF.filter(Boolean).length
      if (upCount >= config.moneyFlowDays * 0.7) {
        mfScore = 100
        mfReason = '连续主力流入'
      } else if (upCount <= config.moneyFlowDays * 0.3) {
        mfScore = 100
        mfReason = '连续主力流出'
      }

      // OBV 方向
      let obvScore = 0
      let obvReason = ''
      if (i >= config.obvLookback) {
        const prev = obv[i - config.obvLookback]
        const now = obv[i]
        const obvChg = (now - prev) / Math.abs(prev) * 100
        if (obvChg > 2) {
          obvScore = 100
          obvReason = 'OBV 上升'
        } else if (obvChg < -2) {
          obvScore = 100
          obvReason = 'OBV 下降'
        }
      }

      // 加权得分
      const total =
        volScore * config.weights.volume +
        breakScore * config.weights.breakout +
        mfScore * config.weights.moneyFlow +
        obvScore * config.weights.obv

      if (total < 60) continue

      const reasons = []
      if (volScore > 50) reasons.push('显著放量')
      if (breakScore > 0) reasons.push(breakReason)
      if (mfScore > 0) reasons.push(mfReason)
      if (obvScore > 0) reasons.push(obvReason)

      const type: 'buy' | 'sell' =
        (priceChg > 0 && total > 70) ? 'buy' : (priceChg < 0 && total > 70) ? 'sell' : 'buy'

      signals.push({
        time: k.time,
        price: k.close,
        type,
        score: Math.round(total),
        reasons,
        aiReason: `综合打分 ${Math.round(total)}，${reasons.join('，')}，建议${type === 'buy' ? '逢低买入' : '逢高减仓'}。`
      })
    }

    // 按时间倒序，取最近 20 条
    return signals.reverse().slice(0, 20)
  }, [klines, drawings, moneyFlow, config])
}