import { useEffect, useState } from 'react'
import { MoneyFlowData } from '../types'
import { Zap, TrendingUp, TrendingDown } from 'lucide-react'

interface SignalTickerProps {
  data: MoneyFlowData[]
}

function SignalTicker({ data }: SignalTickerProps) {
  const [signals, setSignals] = useState<MoneyFlowData[]>([])

  useEffect(() => {
    // 提取最近的 5 条异动信号
    const filtered = data
      .filter(d => d.signal === '扫货' || d.signal === '砸盘')
      .slice(-5)
      .reverse()
    setSignals(filtered)
  }, [data])

  if (signals.length === 0) return null

  return (
    <div className="bg-slate-900 text-white py-1.5 px-4 flex items-center space-x-6 overflow-hidden whitespace-nowrap">
      <div className="flex items-center space-x-2 text-amber-400 shrink-0">
        <Zap className="w-3.5 h-3.5 fill-current" />
        <span className="text-[10px] font-black uppercase tracking-tighter">异动快报</span>
      </div>
      
      <div className="flex items-center space-x-8 animate-marquee">
        {signals.map((s, i) => (
          <div key={i} className="flex items-center space-x-2">
            <span className="text-slate-400 font-mono text-[10px]">{s.time}</span>
            <div className={`flex items-center space-x-1 px-1.5 py-0.5 rounded text-[10px] font-bold ${
              s.signal === '扫货' ? 'bg-red-500/20 text-red-400' : 'bg-green-500/20 text-green-400'
            }`}>
              {s.signal === '扫货' ? <TrendingUp className="w-3 h-3" /> : <TrendingDown className="w-3 h-3" />}
              <span>主力{s.signal}</span>
            </div>
            <span className="text-slate-300 text-[10px]">强度: {(Math.abs(s.mainNet) / 10000).toFixed(1)}万</span>
          </div>
        ))}
      </div>
    </div>
  )
}

export default SignalTicker
