import { useState, useCallback } from 'react'
import { useWailsAPI } from '../hooks/useWailsAPI'
import type { StockData, AnalysisReport } from '../types'

interface StockSearchProps {
  onStockDataLoaded: (data: StockData) => void
  onAnalysisComplete: (report: AnalysisReport) => void
  onError: (error: string) => void
  onLoadingChange: (loading: boolean) => void
  onWatchlistChanged?: () => void
}

function StockSearch({ onStockDataLoaded, onAnalysisComplete, onError, onLoadingChange, onWatchlistChanged }: StockSearchProps) {
  const [stockCode, setStockCode] = useState('')
  const [isSearching, setIsSearching] = useState(false)
  const [searchSuggestions, setSearchSuggestions] = useState<StockData[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [currentStock, setCurrentStock] = useState<StockData | null>(null)
  
  const { getStockData, analyzeStock, searchStock, addToWatchlist } = useWailsAPI()

  const handleGetStockData = useCallback(async () => {
    if (!stockCode.trim()) {
      onError('请输入股票代码')
      return
    }

    setIsSearching(true)
    onLoadingChange(true)
    onError('')
    setShowSuggestions(false)

    try {
      const data = await getStockData(stockCode.trim())
      setCurrentStock(data)
      onStockDataLoaded(data)
      onLoadingChange(false)
    } catch (err: any) {
      onError(err.message || '获取股票数据失败')
      onLoadingChange(false)
    } finally {
      setIsSearching(false)
    }
  }, [stockCode, getStockData, onStockDataLoaded, onLoadingChange, onError])

  const handleAnalyzeStock = useCallback(async () => {
    if (!stockCode.trim()) {
      onError('请输入股票代码')
      return
    }

    onLoadingChange(true)
    onError('')
    setShowSuggestions(false)

    try {
      const report = await analyzeStock(stockCode.trim())
      onAnalysisComplete(report)
      onLoadingChange(false)
    } catch (err: any) {
      onError(err.message || 'AI分析失败')
      onLoadingChange(false)
    }
  }, [stockCode, analyzeStock, onAnalysisComplete, onLoadingChange, onError])

  const handleAddToWatchlist = async () => {
    if (!currentStock) return
    try {
      await addToWatchlist(currentStock)
      if (onWatchlistChanged) onWatchlistChanged()
    } catch (err: any) {
      alert(err.message || '添加失败')
    }
  }

  const handleSearchSuggestions = useCallback(async (keyword: string) => {
    if (!keyword.trim()) {
      setSearchSuggestions([])
      setShowSuggestions(false)
      return
    }

    try {
      const results = await searchStock(keyword)
      setSearchSuggestions(results.slice(0, 5))
      setShowSuggestions(true)
    } catch (err) {
      setSearchSuggestions([])
    }
  }, [searchStock])

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleGetStockData()
    }
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setStockCode(value)
    handleSearchSuggestions(value)
  }

  const handleSuggestionClick = (code: string) => {
    setStockCode(code)
    setShowSuggestions(false)
  }

  const quickStocks = [
    { code: '600519', name: '贵州茅台' },
    { code: '000001', name: '平安银行' },
    { code: '600036', name: '招商银行' },
    { code: '000858', name: '五粮液' },
  ]

  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <h2 className="text-lg font-semibold text-gray-800 mb-4">股票查询</h2>
      
      <div className="space-y-4">
        <div className="relative">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            股票代码
          </label>
          <input
            type="text"
            value={stockCode}
            onChange={handleInputChange}
            onKeyPress={handleKeyPress}
            placeholder="例如: 600519"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
            disabled={isSearching}
            autoComplete="off"
          />
          <p className="mt-1 text-xs text-gray-500">
            支持沪深A股代码，如：600519（茅台）、000001（平安）
          </p>

          {showSuggestions && searchSuggestions.length > 0 && (
            <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-300 rounded-lg shadow-lg z-10">
              {searchSuggestions.map((stock) => (
                <button
                  key={stock.code}
                  onClick={() => handleSuggestionClick(stock.code)}
                  className="w-full text-left px-4 py-2 hover:bg-gray-50 transition border-b border-gray-100 last:border-b-0"
                >
                  <div className="font-medium text-gray-800">{stock.code}</div>
                  <div className="text-xs text-gray-500">{stock.name}</div>
                </button>
              ))}
            </div>
          )}
        </div>

        <div className="flex space-x-3">
          <button
            onClick={handleGetStockData}
            disabled={isSearching}
            className="flex-1 bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded-lg transition disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isSearching ? '查询中...' : '查询数据'}
          </button>
          
          <button
            onClick={handleAnalyzeStock}
            disabled={isSearching}
            className="flex-1 bg-gradient-to-r from-purple-500 to-indigo-600 hover:from-purple-600 hover:to-indigo-700 text-white font-medium py-2 px-4 rounded-lg transition disabled:opacity-50 disabled:cursor-not-allowed"
          >
            AI分析
          </button>
        </div>

        {currentStock && (
          <button
            type="button"
            onClick={handleAddToWatchlist}
            className="w-full py-2 border-2 border-dashed border-gray-200 text-gray-500 hover:border-blue-300 hover:text-blue-500 rounded-lg transition-all flex items-center justify-center space-x-2 text-sm font-medium"
          >
            <span>⭐ 加入自选股</span>
          </button>
        )}

        <div className="pt-4 border-t border-gray-200">
          <h3 className="text-sm font-medium text-gray-700 mb-2">常用股票</h3>
          <div className="grid grid-cols-2 gap-2">
            {quickStocks.map((stock) => (
              <button
                key={stock.code}
                onClick={() => handleSuggestionClick(stock.code)}
                className="text-left px-3 py-2 text-sm bg-gray-50 hover:bg-gray-100 rounded-lg transition"
              >
                <div className="font-medium text-gray-800">{stock.code}</div>
                <div className="text-xs text-gray-500">{stock.name}</div>
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

export default StockSearch
