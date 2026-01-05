import { useState, useEffect } from 'react'
import StockSearch from './components/StockSearch'
import StockInfo from './components/StockInfo'
import AnalysisReport from './components/AnalysisReport'
import Settings from './components/Settings'
import Watchlist from './components/Watchlist'
import { useWailsAPI } from './hooks/useWailsAPI'
import type { StockData, AnalysisReport as AnalysisReportType, NavItem, AppConfig } from './types'

function App() {
  const [activeTab, setActiveTab] = useState<NavItem>('analysis')
  const [stockData, setStockData] = useState<StockData | null>(null)
  const [analysisReport, setAnalysisReport] = useState<AnalysisReportType | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string>('')
  const [currentConfig, setCurrentConfig] = useState<AppConfig | null>(null)
  const [watchlistRefresh, setWatchlistRefresh] = useState(0)

  const { getConfig, getStockData: fetchStockData } = useWailsAPI()

  useEffect(() => {
    fetchConfig()
  }, [])

  const fetchConfig = async () => {
    try {
      const config = await getConfig()
      setCurrentConfig(config)
    } catch (err) {
      console.error('Failed to fetch config:', err)
    }
  }

  const handleStockDataLoaded = (data: StockData) => {
    setStockData(data)
    setAnalysisReport(null)
    setError('')
  }

  const handleAnalysisComplete = (report: AnalysisReportType) => {
    setAnalysisReport(report)
    setError('')
  }

  const handleConfigSaved = () => {
    fetchConfig()
  }

  const handleWatchlistChanged = () => {
    setWatchlistRefresh(prev => prev + 1)
  }

  const handleSelectFromWatchlist = async (code: string) => {
    setLoading(true)
    setError('')
    try {
      const data = await fetchStockData(code)
      setStockData(data)
      setAnalysisReport(null)
    } catch (err: any) {
      setError(err.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex h-screen bg-gray-100 overflow-hidden">
      {/* 侧边菜单栏 */}
      <aside className="w-64 bg-slate-900 text-white flex flex-col shadow-xl z-20">
        <div className="p-6 flex items-center space-x-3 border-b border-slate-800">
          <div className="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center shadow-lg">
            <span className="text-2xl font-bold">股</span>
          </div>
          <div>
            <h1 className="text-lg font-bold tracking-tight">AI-Agent</h1>
            <p className="text-[10px] text-slate-400 uppercase tracking-widest">Stock Analyzer</p>
          </div>
        </div>

        <nav className="flex-1 p-4 space-y-2 mt-4">
          <button
            onClick={() => setActiveTab('analysis')}
            className={`w-full flex items-center space-x-3 px-4 py-3 rounded-xl transition-all duration-200 ${
              activeTab === 'analysis' 
                ? 'bg-blue-600 text-white shadow-lg shadow-blue-900/20' 
                : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'
            }`}
          >
            <span className="text-xl">📊</span>
            <span className="font-medium">股票分析</span>
          </button>

          <button
            onClick={() => setActiveTab('settings')}
            className={`w-full flex items-center space-x-3 px-4 py-3 rounded-xl transition-all duration-200 ${
              activeTab === 'settings' 
                ? 'bg-blue-600 text-white shadow-lg shadow-blue-900/20' 
                : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'
            }`}
          >
            <span className="text-xl">⚙️</span>
            <span className="font-medium">系统设置</span>
          </button>
        </nav>

        <div className="p-6 border-t border-slate-800">
          <div className="bg-slate-800/50 rounded-lg p-3">
            <p className="text-[10px] text-slate-500 mb-1">当前版本</p>
            <p className="text-xs font-mono text-slate-300">v1.2.0 (Watchlist)</p>
          </div>
        </div>
      </aside>

      {/* 主内容区域 */}
      <main className="flex-1 flex flex-col min-w-0 bg-slate-50 relative">
        {/* 顶部状态栏 */}
        <header className="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-8 z-10">
          <h2 className="text-lg font-semibold text-gray-800">
            {activeTab === 'analysis' ? '股票分析工作台' : '系统参数配置'}
          </h2>
          <div className="flex items-center space-x-6 text-sm">
            {currentConfig && (
              <div className="flex items-center bg-blue-50 text-blue-700 px-3 py-1.5 rounded-lg border border-blue-100">
                <span className="mr-2">🤖</span>
                <span className="font-medium mr-1">当前模型:</span>
                <span className="font-mono text-xs">{currentConfig.model}</span>
              </div>
            )}
            <span className="flex items-center text-gray-500">
              <span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
              API 状态: 正常
            </span>
          </div>
        </header>

        {/* 内容滚动区 */}
        <div className="flex-1 overflow-y-auto p-8">
          {activeTab === 'analysis' ? (
            <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-4 gap-8">
              {/* 左侧：搜索和自选股 */}
              <div className="lg:col-span-1 space-y-8">
                <StockSearch
                  onStockDataLoaded={handleStockDataLoaded}
                  onAnalysisComplete={handleAnalysisComplete}
                  onError={setError}
                  onLoadingChange={setLoading}
                  onWatchlistChanged={handleWatchlistChanged}
                />
                <Watchlist 
                  onSelectStock={handleSelectFromWatchlist} 
                  refreshTrigger={watchlistRefresh}
                />
              </div>

              {/* 中间/右侧：股票信息和分析报告 */}
              <div className="lg:col-span-3 space-y-8">
                {error && (
                  <div className="bg-red-50 border border-red-200 rounded-xl p-4 flex items-start">
                    <span className="mr-3 mt-0.5">⚠️</span>
                    <p className="text-red-700 text-sm">{error}</p>
                  </div>
                )}

                {stockData && <StockInfo stockData={stockData} />}

                {loading ? (
                  <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-12 flex flex-col items-center justify-center min-h-[400px]">
                    <div className="relative w-20 h-20 mb-6">
                      <div className="absolute inset-0 border-4 border-blue-100 rounded-full"></div>
                      <div className="absolute inset-0 border-4 border-blue-600 rounded-full border-t-transparent animate-spin"></div>
                    </div>
                    <h3 className="text-lg font-semibold text-gray-800 mb-2">AI 正在深度分析中</h3>
                    <p className="text-gray-500 text-sm">正在调用 {currentConfig?.model || 'Qwen'} 模型进行多维度评估...</p>
                  </div>
                ) : analysisReport ? (
                  <AnalysisReport report={analysisReport} />
                ) : !stockData ? (
                  <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-12 flex flex-col items-center justify-center text-center min-h-[400px]">
                    <div className="w-24 h-24 bg-slate-50 rounded-full flex items-center justify-center mb-6">
                      <span className="text-4xl">📈</span>
                    </div>
                    <h3 className="text-xl font-bold text-gray-800 mb-2">准备就绪</h3>
                    <p className="text-gray-500 max-w-sm text-sm leading-relaxed">
                      请输入 A 股代码或从自选股中选择，开始您的智能投资分析之旅。
                    </p>
                  </div>
                ) : null}
              </div>
            </div>
          ) : (
            <Settings onConfigSaved={handleConfigSaved} />
          )}
        </div>

        {/* 底部免责声明 */}
        <footer className="h-10 bg-white border-t border-gray-100 flex items-center justify-center px-8 text-[10px] text-gray-400 uppercase tracking-widest">
          ⚠️ Disclaimer: AI-generated content is for reference only and does not constitute investment advice.
        </footer>
      </main>
    </div>
  )
}

export default App
