import { useState, useEffect } from 'react'
import { useWailsAPI } from '../hooks/useWailsAPI'
import KLineChart from './KLineChart'
import type { StockData, KLineData } from '../types'

interface WatchlistDetailProps {
  stock: StockData
}

type Period = 'daily' | 'week' | 'month'

function WatchlistDetail({ stock }: WatchlistDetailProps) {
  const [klines, setKlines] = useState<KLineData[]>([])
  const [loading, setLoading] = useState(true)
  const [period, setPeriod] = useState<Period>('daily')
  const [techAnalysis, setTechAnalysis] = useState<string>('')
  const [analyzing, setAnalyzing] = useState(false)
  const [indicators, setIndicators] = useState({
    macd: false,
    kdj: false,
    rsi: false
  })
  const { getKLineData, analyzeTechnical } = useWailsAPI()

  useEffect(() => {
    loadKLines()
    setTechAnalysis('') // 切换股票或周期时清空分析
  }, [stock.code, period])

  const loadKLines = async () => {
    setLoading(true)
    try {
      const data = await getKLineData(stock.code, 150, period)
      setKlines(data)
    } catch (err) {
      console.error('Failed to load K-lines:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleTechnicalAnalysis = async () => {
    setAnalyzing(true)
    try {
      const result = await analyzeTechnical(stock.code, period)
      setTechAnalysis(result)
    } catch (err) {
      console.error('Technical analysis failed:', err)
    } finally {
      setAnalyzing(false)
    }
  }

  const toggleIndicator = (key: keyof typeof indicators) => {
    setIndicators(prev => ({ ...prev, [key]: !prev[key] }))
  }

  const periodLabels: Record<Period, string> = {
    daily: '日线',
    week: '周线',
    month: '月线'
  }

  return (
    <div className="space-y-6 pb-12">
      {/* 头部概览 */}
      <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 flex items-center justify-between">
        <div>
          <div className="flex items-center space-x-3">
            <h2 className="text-2xl font-bold text-gray-900">{stock.name}</h2>
            <span className="px-2 py-1 bg-gray-100 text-gray-500 rounded text-xs font-mono">{stock.code}</span>
          </div>
          <div className="flex items-center space-x-4 mt-2">
            <span className="text-3xl font-mono font-bold text-gray-900">{stock.price.toFixed(2)}</span>
            <div className={`flex flex-col ${stock.changeRate >= 0 ? 'text-red-500' : 'text-green-500'}`}>
              <span className="text-sm font-bold">{stock.changeRate >= 0 ? '+' : ''}{stock.changeRate.toFixed(2)}%</span>
              <span className="text-xs">{stock.change >= 0 ? '+' : ''}{stock.change.toFixed(2)}</span>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-3 gap-8 text-right">
          <div>
            <p className="text-[10px] text-gray-400 uppercase font-bold">最高</p>
            <p className="text-sm font-mono font-bold text-red-500">{stock.high.toFixed(2)}</p>
          </div>
          <div>
            <p className="text-[10px] text-gray-400 uppercase font-bold">最低</p>
            <p className="text-sm font-mono font-bold text-green-500">{stock.low.toFixed(2)}</p>
          </div>
          <div>
            <p className="text-[10px] text-gray-400 uppercase font-bold">成交量</p>
            <p className="text-sm font-mono font-bold text-gray-700">{(stock.volume / 10000).toFixed(2)}万</p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-6">
        {/* 左侧：K线图 (占2列) */}
        <div className="col-span-2 space-y-6">
          <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
            <div className="flex items-center justify-between mb-6">
              <div className="flex items-center space-x-4">
                <h3 className="font-bold text-gray-800 flex items-center">
                  <span className="mr-2">📈</span> 行情图表
                </h3>
                
                <select 
                  value={period}
                  onChange={(e) => setPeriod(e.target.value as Period)}
                  className="block pl-3 pr-10 py-1 text-xs font-bold border-gray-200 focus:outline-none focus:ring-blue-500 focus:border-blue-500 rounded-md bg-gray-50 text-gray-700 cursor-pointer"
                >
                  <option value="daily">日线</option>
                  <option value="week">周线</option>
                  <option value="month">月线</option>
                </select>

                <div className="flex bg-gray-100 p-1 rounded-lg">
                  {['macd', 'kdj', 'rsi'].map((key) => (
                    <button 
                      key={key}
                      onClick={() => toggleIndicator(key as any)}
                      className={`px-3 py-1 rounded-md text-xs font-bold transition-all ${indicators[key as keyof typeof indicators] ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-400 hover:text-gray-600'}`}
                    >
                      {key.toUpperCase()}
                    </button>
                  ))}
                </div>
              </div>
            </div>
            
            {loading ? (
              <div className="h-[500px] flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
              </div>
            ) : (
              <KLineChart 
                data={klines} 
                height={500}
                showMACD={indicators.macd}
                showKDJ={indicators.kdj}
                showRSI={indicators.rsi}
              />
            )}
          </div>

          {/* 更多指标 */}
          <div className="grid grid-cols-4 gap-4">
            {[
              { label: '换手率', value: `${stock.turnover.toFixed(2)}%` },
              { label: '市盈率(动)', value: stock.pe.toFixed(2) },
              { label: '市净率', value: stock.pb.toFixed(2) },
              { label: '总市值', value: `${(stock.totalMV / 100000000).toFixed(2)}亿` },
            ].map((item, i) => (
              <div key={i} className="bg-white rounded-xl border border-gray-100 p-4">
                <p className="text-[10px] text-gray-400 uppercase font-bold mb-1">{item.label}</p>
                <p className="text-lg font-mono font-bold text-gray-800">{item.value}</p>
              </div>
            ))}
          </div>
        </div>

        {/* 右侧：技术分析师面板 (占1列) */}
        <div className="col-span-1">
          <div className="bg-slate-900 rounded-2xl shadow-xl border border-slate-800 p-6 h-full flex flex-col">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-white font-bold flex items-center">
                <span className="mr-2 text-xl">👨‍💻</span> 技术分析师
              </h3>
              <span className="px-2 py-0.5 bg-blue-500/20 text-blue-400 rounded text-[10px] font-bold border border-blue-500/30">PRO</span>
            </div>

            <div className="flex-1 overflow-y-auto custom-scrollbar">
              {analyzing ? (
                <div className="flex flex-col items-center justify-center h-full space-y-4">
                  <div className="w-12 h-12 border-4 border-blue-500/20 border-t-blue-500 rounded-full animate-spin"></div>
                  <p className="text-slate-400 text-xs animate-pulse">正在深度复盘量价形态...</p>
                </div>
              ) : techAnalysis ? (
                <div className="prose prose-invert prose-sm max-w-none">
                  <div className="text-slate-300 leading-relaxed whitespace-pre-wrap text-xs">
                    {techAnalysis}
                  </div>
                </div>
              ) : (
                <div className="flex flex-col items-center justify-center h-full text-center space-y-4 px-4">
                  <div className="w-16 h-16 bg-slate-800 rounded-full flex items-center justify-center text-3xl">📊</div>
                  <div>
                    <p className="text-slate-300 font-bold text-sm">需要深度技术解读吗？</p>
                    <p className="text-slate-500 text-xs mt-1">我将结合当前 {periodLabels[period]} 的 K 线形态和技术指标为您提供操盘建议。</p>
                  </div>
                </div>
              )}
            </div>

            <button
              onClick={handleTechnicalAnalysis}
              disabled={analyzing || loading}
              className={`mt-6 w-full py-3 rounded-xl font-bold text-sm transition-all flex items-center justify-center space-x-2 ${
                analyzing || loading 
                ? 'bg-slate-800 text-slate-500 cursor-not-allowed' 
                : 'bg-blue-600 hover:bg-blue-500 text-white shadow-lg shadow-blue-900/20 active:scale-[0.98]'
              }`}
            >
              {analyzing ? '分析中...' : '开始深度技术分析'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export default WatchlistDetail
