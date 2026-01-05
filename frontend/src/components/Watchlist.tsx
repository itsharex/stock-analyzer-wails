import { useState, useEffect } from 'react'
import { useWailsAPI } from '../hooks/useWailsAPI'
import type { StockData } from '../types'
import BatchAnalyzeModal from './BatchAnalyzeModal'
import { Brain } from 'lucide-react'

interface WatchlistProps {
  onSelectStock: (code: string) => void
  refreshTrigger: number
}

function Watchlist({ onSelectStock, refreshTrigger }: WatchlistProps) {
  const [stocks, setStocks] = useState<StockData[]>([])
  const [loading, setLoading] = useState(true)
  const [isBatchModalOpen, setIsBatchModalOpen] = useState(false)
  const { getWatchlist, removeFromWatchlist, batchAnalyzeStocks } = useWailsAPI()

  useEffect(() => {
    loadWatchlist()
  }, [refreshTrigger])

  const loadWatchlist = async () => {
    try {
      const data = await getWatchlist()
      setStocks(data)
    } catch (err) {
      console.error('Failed to load watchlist:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleRemove = async (e: React.MouseEvent, code: string) => {
    e.stopPropagation()
    try {
      await removeFromWatchlist(code)
      loadWatchlist()
    } catch (err) {
      alert('删除失败')
    }
  }

  if (loading) return <div className="p-4 text-center text-gray-500 text-sm">加载中...</div>

  if (stocks.length === 0) {
    return (
      <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 text-center">
        <p className="text-gray-400 text-sm">暂无自选股</p>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
      <div className="p-4 border-b border-gray-50 bg-gray-50/50 flex justify-between items-center">
        <div className="flex items-center space-x-2">
          <h3 className="font-bold text-gray-800 flex items-center">
            <span className="mr-2">⭐</span> 我的自选
          </h3>
          <span className="text-[10px] bg-blue-100 text-blue-600 px-2 py-0.5 rounded-full font-bold">
            {stocks.length}
          </span>
        </div>
        <button
          onClick={() => setIsBatchModalOpen(true)}
          className="flex items-center space-x-1 px-2 py-1 bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 transition-colors text-xs font-bold"
          title="批量 AI 分析"
        >
          <Brain className="w-3 h-3" />
          <span>批量分析</span>
        </button>
      </div>
      <div className="divide-y divide-gray-50 max-h-[400px] overflow-y-auto">
        {stocks.map((stock) => (
          <div
            key={stock.code}
            onClick={() => onSelectStock(stock.code)}
            className="p-4 hover:bg-blue-50 cursor-pointer transition-colors group flex items-center justify-between"
          >
            <div className="flex-1 min-w-0">
              <div className="flex items-center space-x-2">
                <span className="font-bold text-gray-900 truncate">{stock.name}</span>
                <span className="text-[10px] font-mono text-gray-400">{stock.code}</span>
              </div>
              <div className="flex items-center space-x-3 mt-1">
                <span className="text-sm font-mono font-semibold text-gray-700">
                  {stock.price.toFixed(2)}
                </span>
                <span className={`text-xs font-medium ${stock.changeRate >= 0 ? 'text-red-500' : 'text-green-500'}`}>
                  {stock.changeRate >= 0 ? '+' : ''}{stock.changeRate.toFixed(2)}%
                </span>
              </div>
            </div>
            <button
              onClick={(e) => handleRemove(e, stock.code)}
              className="opacity-0 group-hover:opacity-100 p-2 hover:bg-red-100 text-red-400 hover:text-red-600 rounded-lg transition-all"
              title="移出自选"
            >
              <span className="text-lg">🗑️</span>
            </button>
          </div>
        ))}
      </div>

      <BatchAnalyzeModal
        isOpen={isBatchModalOpen}
        onClose={() => setIsBatchModalOpen(false)}
        stocks={stocks.map(s => ({ code: s.code, name: s.name }))}
        onStart={(codes) => batchAnalyzeStocks(codes, 'technical_master')}
      />
    </div>
  )
}

export default Watchlist
