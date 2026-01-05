import { useState } from 'react'
import StockSearch from './components/StockSearch'
import StockInfo from './components/StockInfo'
import AnalysisReport from './components/AnalysisReport'

function App() {
  const [stockData, setStockData] = useState<StockData | null>(null)
  const [analysisReport, setAnalysisReport] = useState<AnalysisReport | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string>('')

  const handleStockDataLoaded = (data: StockData) => {
    setStockData(data)
    setAnalysisReport(null)
    setError('')
  }

  const handleAnalysisComplete = (report: AnalysisReport) => {
    setAnalysisReport(report)
    setError('')
  }

  const handleError = (errorMsg: string) => {
    setError(errorMsg)
  }

  const handleLoadingChange = (isLoading: boolean) => {
    setLoading(isLoading)
  }

  return (
    <div className="w-full h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex flex-col">
      {/* 顶部标题栏 */}
      <header className="bg-white shadow-md px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-lg flex items-center justify-center">
              <span className="text-white text-xl font-bold">股</span>
            </div>
            <div>
              <h1 className="text-2xl font-bold text-gray-800">A股股票分析AI-Agent</h1>
              <p className="text-sm text-gray-500">专业的AI驱动股票分析工具</p>
            </div>
          </div>
          <div className="text-sm text-gray-500">
            数据来源: 东方财富 | AI: OpenAI GPT
          </div>
        </div>
      </header>

      {/* 主内容区域 */}
      <main className="flex-1 overflow-hidden p-6">
        <div className="h-full max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* 左侧：搜索和股票信息 */}
          <div className="lg:col-span-1 flex flex-col space-y-6">
            <StockSearch
              onStockDataLoaded={handleStockDataLoaded}
              onAnalysisComplete={handleAnalysisComplete}
              onError={handleError}
              onLoadingChange={handleLoadingChange}
            />
            {stockData && <StockInfo stockData={stockData} />}
          </div>

          {/* 右侧：分析报告 */}
          <div className="lg:col-span-2">
            {error && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-4">
                <div className="flex items-center">
                  <svg className="w-5 h-5 text-red-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                  </svg>
                  <span className="text-red-700">{error}</span>
                </div>
              </div>
            )}

            {loading && (
              <div className="bg-white rounded-lg shadow-lg p-8 flex flex-col items-center justify-center">
                <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-blue-500 mb-4"></div>
                <p className="text-gray-600">AI正在分析中，请稍候...</p>
              </div>
            )}

            {!loading && analysisReport && (
              <AnalysisReport report={analysisReport} />
            )}

            {!loading && !analysisReport && !error && (
              <div className="bg-white rounded-lg shadow-lg p-12 flex flex-col items-center justify-center text-center h-full">
                <svg className="w-24 h-24 text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
                <h3 className="text-xl font-semibold text-gray-700 mb-2">开始分析</h3>
                <p className="text-gray-500 max-w-md">
                  输入股票代码并点击"AI分析"按钮，获取专业的投资分析报告
                </p>
              </div>
            )}
          </div>
        </div>
      </main>

      {/* 底部信息栏 */}
      <footer className="bg-white border-t border-gray-200 px-6 py-3">
        <div className="max-w-7xl mx-auto flex items-center justify-between text-sm text-gray-500">
          <div>
            © 2026 Stock Analyzer. 版本 1.0.0
          </div>
          <div className="flex items-center space-x-4">
            <span>⚠️ 免责声明：本工具仅供学习研究使用，不构成投资建议</span>
          </div>
        </div>
      </footer>
    </div>
  )
}

export default App
