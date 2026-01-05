import { useState } from 'react'

// 模拟Wails运行时导入（实际使用时会由Wails自动生成）
// import { GetStockData, AnalyzeStock } from '../../wailsjs/go/main/App'

interface StockSearchProps {
  onStockDataLoaded: (data: StockData) => void
  onAnalysisComplete: (report: AnalysisReport) => void
  onError: (error: string) => void
  onLoadingChange: (loading: boolean) => void
}

// 模拟Wails API调用（开发时使用）
const mockGetStockData = async (code: string): Promise<StockData> => {
  // 模拟网络延迟
  await new Promise(resolve => setTimeout(resolve, 500))
  
  return {
    code: code,
    name: '贵州茅台',
    price: 1688.88,
    change: 15.88,
    changeRate: 0.95,
    volume: 1234567,
    amount: 2088888888,
    high: 1699.99,
    low: 1680.00,
    open: 1685.00,
    preClose: 1673.00,
    amplitude: 1.19,
    turnover: 0.85,
    pe: 35.6,
    pb: 12.8,
    totalMV: 2120000000000,
    circMV: 2120000000000,
  }
}

const mockAnalyzeStock = async (code: string): Promise<AnalysisReport> => {
  // 模拟AI分析延迟
  await new Promise(resolve => setTimeout(resolve, 2000))
  
  return {
    stockCode: code,
    stockName: '贵州茅台',
    summary: '贵州茅台作为A股白酒龙头，基本面稳健，品牌价值突出。当前估值处于合理区间，短期走势偏强。',
    fundamentals: '公司市盈率35.6倍，市净率12.8倍，处于行业中等水平。总市值2.12万亿元，流通市值充足。作为白酒行业龙头，公司具有强大的品牌护城河和定价能力，盈利能力持续稳定。',
    technical: '股价近期呈现上涨趋势，今日涨幅0.95%，成交量适中，换手率0.85%显示市场活跃度良好。振幅1.19%表明波动较小，多头趋势明显。',
    recommendation: '建议：持有或适量买入。理由：1) 基本面优秀，业绩稳定增长；2) 估值合理，具有长期投资价值；3) 短期技术面良好，上涨趋势明确。适合中长期价值投资者配置。',
    riskLevel: '中等风险',
    targetPrice: '目标价位区间：1750-1850元',
    generatedAt: new Date().toLocaleString('zh-CN'),
  }
}

function StockSearch({ onStockDataLoaded, onAnalysisComplete, onError, onLoadingChange }: StockSearchProps) {
  const [stockCode, setStockCode] = useState('')
  const [isSearching, setIsSearching] = useState(false)

  const handleGetStockData = async () => {
    if (!stockCode.trim()) {
      onError('请输入股票代码')
      return
    }

    setIsSearching(true)
    onLoadingChange(true)
    onError('')

    try {
      // 实际使用时取消注释以下行，注释掉mock调用
      // const data = await GetStockData(stockCode.trim())
      const data = await mockGetStockData(stockCode.trim())
      
      onStockDataLoaded(data)
      onLoadingChange(false)
    } catch (err: any) {
      onError(err.message || '获取股票数据失败')
      onLoadingChange(false)
    } finally {
      setIsSearching(false)
    }
  }

  const handleAnalyzeStock = async () => {
    if (!stockCode.trim()) {
      onError('请输入股票代码')
      return
    }

    onLoadingChange(true)
    onError('')

    try {
      // 实际使用时取消注释以下行，注释掉mock调用
      // const report = await AnalyzeStock(stockCode.trim())
      const report = await mockAnalyzeStock(stockCode.trim())
      
      onAnalysisComplete(report)
      onLoadingChange(false)
    } catch (err: any) {
      onError(err.message || 'AI分析失败')
      onLoadingChange(false)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleGetStockData()
    }
  }

  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <h2 className="text-lg font-semibold text-gray-800 mb-4">股票查询</h2>
      
      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            股票代码
          </label>
          <input
            type="text"
            value={stockCode}
            onChange={(e) => setStockCode(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="例如: 600519"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
            disabled={isSearching}
          />
          <p className="mt-1 text-xs text-gray-500">
            支持沪深A股代码，如：600519（茅台）、000001（平安）
          </p>
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

        <div className="pt-4 border-t border-gray-200">
          <h3 className="text-sm font-medium text-gray-700 mb-2">常用股票</h3>
          <div className="grid grid-cols-2 gap-2">
            {[
              { code: '600519', name: '贵州茅台' },
              { code: '000001', name: '平安银行' },
              { code: '600036', name: '招商银行' },
              { code: '000858', name: '五粮液' },
            ].map((stock) => (
              <button
                key={stock.code}
                onClick={() => setStockCode(stock.code)}
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
