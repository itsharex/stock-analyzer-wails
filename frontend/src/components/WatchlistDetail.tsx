import { useState, useEffect } from 'react'
import { useWailsAPI } from '../hooks/useWailsAPI'
import KLineChart from './KLineChart'
import type { StockData, KLineData } from '../types'

interface WatchlistDetailProps {
  stock: StockData
}

function WatchlistDetail({ stock }: WatchlistDetailProps) {
  const [klines, setKlines] = useState<KLineData[]>([])
  const [loading, setLoading] = useState(true)
  const { getKLineData } = useWailsAPI()

  useEffect(() => {
    loadKLines()
  }, [stock.code])

  const loadKLines = async () => {
    setLoading(true)
    try {
      const data = await getKLineData(stock.code, 100)
      setKlines(data)
    } catch (err) {
      console.error('Failed to load K-lines:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
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

      {/* K线图区域 */}
      <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
        <div className="flex items-center justify-between mb-6">
          <h3 className="font-bold text-gray-800 flex items-center">
            <span className="mr-2">📈</span> 日K线图
          </h3>
          <div className="flex space-x-2">
            <span className="px-2 py-1 bg-blue-50 text-blue-600 rounded text-[10px] font-bold">前复权</span>
            <span className="px-2 py-1 bg-gray-50 text-gray-400 rounded text-[10px] font-bold">日线</span>
          </div>
        </div>
        
        {loading ? (
          <div className="h-[400px] flex items-center justify-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
          </div>
        ) : (
          <KLineChart data={klines} />
        )}
      </div>

      {/* 更多指标 */}
      <div className="grid grid-cols-4 gap-6">
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
  )
}

export default WatchlistDetail
