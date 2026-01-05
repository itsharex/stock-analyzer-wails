import { useState, useEffect } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { StockData, KLineData, TechnicalAnalysisResult } from '../types'
import { useWailsAPI } from '../hooks/useWailsAPI'
import KLineChart from './KLineChart'
import RadarChart from './RadarChart'
import TradePlanCard from './TradePlanCard'
import { 
  Activity, 
  Clock, 
  ChevronDown,
  BrainCircuit,
  Loader2,
  LineChart as LineChartIcon,
  PencilRuler,
  ShieldCheck,
  Zap,
  BarChart3
} from 'lucide-react'
import { GlossaryPanel, GlossaryTooltip } from './GlossaryTooltip'
import { STOCK_GLOSSARY } from '../utils/glossary'

interface WatchlistDetailProps {
  stock: StockData
}

function WatchlistDetail({ stock }: WatchlistDetailProps) {
  const { getKLineData, analyzeTechnical } = useWailsAPI()
  const [klineData, setKlineData] = useState<KLineData[]>([])
  const [period, setPeriod] = useState<string>('daily')
  const [loading, setLoading] = useState(false)
  const [analysisLoading, setAnalysisLoading] = useState(false)
  const [analysisResult, setAnalysisResult] = useState<TechnicalAnalysisResult | null>(null)
  
  // 指标显示控制
  const [showMACD, setShowMACD] = useState(false)
  const [showKDJ, setShowKDJ] = useState(false)
  const [showRSI, setShowRSI] = useState(false)
  const [showAIDrawings, setShowAIDrawings] = useState(true)

  useEffect(() => {
    loadKLineData()
  }, [stock.code, period])

  const loadKLineData = async () => {
    setLoading(true)
    try {
      const data = await getKLineData(stock.code, 100, period)
      setKlineData(data)
    } catch (error) {
      console.error('加载K线数据失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAnalyze = async () => {
    setAnalysisLoading(true)
    try {
      const result = await analyzeTechnical(stock.code, period)
      setAnalysisResult(result)
      setShowAIDrawings(true)
    } catch (error) {
      console.error('技术分析失败:', error)
    } finally {
      setAnalysisLoading(false)
    }
  }

  const getActionColor = (advice: string) => {
    switch (advice) {
      case '买入': case '增持': return 'bg-red-500 text-white'
      case '卖出': case '减持': return 'bg-green-500 text-white'
      default: return 'bg-slate-500 text-white'
    }
  }



  return (
    <div className="flex flex-col h-full bg-slate-50 overflow-hidden">
      {/* 顶部行情概览 */}
      <div className="bg-white p-4 border-b border-slate-200 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <div>
            <h2 className="text-2xl font-bold text-slate-900">{stock.name}</h2>
            <p className="text-sm text-slate-500">{stock.code}</p>
          </div>
          <div className="h-10 w-px bg-slate-200 mx-2" />
          <div>
            <div className={`text-2xl font-mono font-bold ${stock.change >= 0 ? 'text-red-500' : 'text-green-500'}`}>
              {stock.price.toFixed(2)}
            </div>
            <div className={`text-sm font-medium ${stock.change >= 0 ? 'text-red-500' : 'text-green-500'}`}>
              {stock.change >= 0 ? '+' : ''}{stock.change.toFixed(2)} ({stock.changeRate.toFixed(2)}%)
            </div>
          </div>
        </div>

        <div className="grid grid-cols-4 gap-8">
          <div className="text-center">
            <p className="text-xs text-slate-400 uppercase tracking-wider">最高</p>
            <p className="text-sm font-semibold text-slate-700">{stock.high.toFixed(2)}</p>
          </div>
          <div className="text-center">
            <p className="text-xs text-slate-400 uppercase tracking-wider">最低</p>
            <p className="text-sm font-semibold text-slate-700">{stock.low.toFixed(2)}</p>
          </div>
          <div className="text-center">
            <p className="text-xs text-slate-400 uppercase tracking-wider">成交量</p>
            <p className="text-sm font-semibold text-slate-700">{(stock.volume / 10000).toFixed(2)}万</p>
          </div>
          <div className="text-center">
            <p className="text-xs text-slate-400 uppercase tracking-wider">换手率</p>
            <p className="text-sm font-semibold text-slate-700">{stock.turnover.toFixed(2)}%</p>
          </div>
        </div>
      </div>

      {/* 主体内容区 */}
      <div className="flex-1 flex overflow-hidden">
        {/* 左侧图表区 */}
        <div className="flex-1 flex flex-col p-2 lg:p-4 overflow-y-auto">
          <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-4 flex-1 flex flex-col">
            {/* 图表控制栏 */}
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center space-x-2">
                <div className="relative">
                  <select 
                    value={period}
                    onChange={(e) => setPeriod(e.target.value)}
                    className="appearance-none bg-slate-100 border-none rounded-lg px-4 py-2 pr-10 text-sm font-medium text-slate-700 focus:ring-2 focus:ring-blue-500 cursor-pointer"
                  >
                    <option value="daily">日线</option>
                    <option value="weekly">周线</option>
                    <option value="monthly">月线</option>
                  </select>
                  <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                </div>
                
                <div className="h-6 w-px bg-slate-200 mx-2" />
                
                <div className="flex bg-slate-100 rounded-lg p-1">
                  <button 
                    onClick={() => setShowMACD(!showMACD)}
                    className={`px-3 py-1 text-xs font-bold rounded-md transition-all ${showMACD ? 'bg-white text-blue-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
                  >
                    MACD
                  </button>
                  <button 
                    onClick={() => setShowKDJ(!showKDJ)}
                    className={`px-3 py-1 text-xs font-bold rounded-md transition-all ${showKDJ ? 'bg-white text-purple-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
                  >
                    KDJ
                  </button>
                  <button 
                    onClick={() => setShowRSI(!showRSI)}
                    className={`px-3 py-1 text-xs font-bold rounded-md transition-all ${showRSI ? 'bg-white text-cyan-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
                  >
                    RSI
                  </button>
                </div>

                {analysisResult && (
                  <button 
                    onClick={() => setShowAIDrawings(!showAIDrawings)}
                    className={`flex items-center space-x-1 px-3 py-1 text-xs font-bold rounded-md transition-all ${showAIDrawings ? 'bg-blue-50 text-blue-600 border border-blue-200' : 'bg-slate-100 text-slate-500'}`}
                  >
                    <PencilRuler className="w-3 h-3" />
                    <span>AI 绘图</span>
                  </button>
                )}
              </div>

              <div className="flex items-center text-slate-400 text-xs space-x-4">
                <span className="flex items-center"><Clock className="w-3 h-3 mr-1" /> 实时更新</span>
                <span className="flex items-center"><Activity className="w-3 h-3 mr-1" /> 东方财富数据源</span>
              </div>
            </div>

            {/* K线图容器 */}
            <div className="flex-1 relative min-h-[500px] bg-slate-50 rounded-lg border border-slate-100 overflow-hidden">
              {loading ? (
                <div className="absolute inset-0 flex items-center justify-center bg-white/80 z-10">
                  <Loader2 className="w-8 h-8 text-blue-500 animate-spin" />
                </div>
              ) : (
                <KLineChart 
                  data={klineData} 
                  drawings={showAIDrawings ? analysisResult?.drawings : []}
                  showMACD={showMACD} 
                  showKDJ={showKDJ} 
                  showRSI={showRSI} 
                />
              )}
            </div>
          </div>
        </div>

        {/* 右侧技术分析师面板 */}
        <div className="w-[480px] bg-slate-50 border-l border-slate-200 flex flex-col shadow-[-10px_0_30px_rgba(0,0,0,0.03)] z-10">
          <div className="p-5 border-b border-slate-200 flex items-center justify-between bg-white/50 backdrop-blur-sm">
            <div className="flex items-center space-x-2.5 text-blue-600">
              <div className="p-1.5 bg-blue-50 rounded-lg">
                <BrainCircuit className="w-5 h-5" />
              </div>
              <h3 className="font-bold text-slate-800 tracking-tight">技术分析师</h3>
            </div>
            <button 
              onClick={handleAnalyze}
              disabled={analysisLoading || klineData.length === 0}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-300 text-white text-xs font-bold rounded-xl transition-all shadow-sm hover:shadow-md flex items-center space-x-1.5"
            >
              {analysisLoading ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <LineChartIcon className="w-3.5 h-3.5" />}
              <span>{analysisResult ? '重新分析' : '开始深度分析'}</span>
            </button>
          </div>

          <div className="flex-1 overflow-y-auto p-6 custom-scrollbar bg-gradient-to-b from-white/30 to-transparent">
            {analysisLoading ? (
              <div className="flex flex-col items-center justify-center h-64 space-y-5 text-slate-400">
                <div className="relative">
                  <div className="absolute inset-0 bg-blue-500/10 rounded-full blur-xl animate-pulse" />
                  <BrainCircuit className="w-14 h-14 text-blue-500/30 relative z-10" />
                  <Loader2 className="absolute inset-0 w-14 h-14 text-blue-500 animate-spin relative z-20" />
                </div>
                <p className="text-sm font-medium animate-pulse">正在识别形态并评估风险...</p>
              </div>
            ) : analysisResult ? (
              <div className="space-y-8">
                {/* 多维度评分雷达图 */}
                {analysisResult.radarData && (
                  <div className="bg-white border border-slate-200 rounded-2xl p-4 shadow-sm">
                    <div className="flex items-center space-x-2 mb-4">
                      <div className="p-1 bg-blue-50 rounded-lg">
                        <BarChart3 className="w-4 h-4 text-blue-600" />
                      </div>
                      <h4 className="text-sm font-bold text-slate-800">多维度投资评分</h4>
                    </div>
                    <RadarChart data={analysisResult.radarData} />
                  </div>
                )}

                {/* 风险与建议看板 */}
                <div className="grid grid-cols-2 gap-4">
                  <div className="bg-white border border-slate-200 rounded-2xl p-4 flex flex-col items-center justify-center shadow-sm hover:shadow-md transition-shadow">
                    <p className="text-[10px] text-slate-400 uppercase tracking-widest font-bold mb-2">操盘建议</p>
                    <div className={`px-5 py-1.5 rounded-full text-sm font-black shadow-sm ${getActionColor(analysisResult.actionAdvice)}`}>
                      {analysisResult.actionAdvice}
                    </div>
                  </div>
                  <div className="bg-white border border-slate-200 rounded-2xl p-4 flex flex-col items-center justify-center shadow-sm hover:shadow-md transition-shadow">
                    <p className="text-[10px] text-slate-400 uppercase tracking-widest font-bold mb-2">风险得分</p>
                    <div className="text-2xl font-black text-slate-800">
                      {analysisResult.riskScore}
                    </div>
                  </div>
                </div>

                {/* 智能交易计划 */}
                {analysisResult.tradePlan && (
                  <TradePlanCard plan={analysisResult.tradePlan} currentPrice={stock.price} />
                )}

                {/* 核心结论 */}
                <div className="bg-blue-600/5 border border-blue-600/10 rounded-2xl p-5 mb-6 relative overflow-hidden group">
                  <div className="absolute top-0 left-0 w-1 h-full bg-blue-500" />
                  <p className="text-blue-700 text-sm font-bold flex items-center mb-1">
                    <Zap className="w-4 h-4 mr-2 fill-blue-500" />
                    核心结论
                  </p>
                  <p className="text-slate-600 text-sm leading-relaxed">
                    {analysisResult.analysis.split('\n')[0].replace(/^[#\s*]+/, '')}
                  </p>
                </div>

                {/* Markdown 渲染内容 */}
                <div className="prose prose-slate prose-sm max-w-none prose-headings:text-slate-800 prose-p:text-slate-600 prose-strong:text-blue-600">
                  <ReactMarkdown 
                    remarkPlugins={[remarkGfm]}
                    components={{
                      strong: ({ children }) => {
                        const content = String(children);
                        if (STOCK_GLOSSARY[content]) {
                          return <GlossaryTooltip term={content}>{children}</GlossaryTooltip>;
                        }
                        return <strong className="font-bold text-blue-600">{children}</strong>;
                      }
                    }}
                  >
                    {analysisResult.analysis}
                  </ReactMarkdown>
                </div>

                {/* 术语百科面板 */}
                <GlossaryPanel text={analysisResult.analysis} />
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center h-80 text-center space-y-5">
                <div className="w-20 h-20 bg-white rounded-3xl shadow-sm border border-slate-100 flex items-center justify-center group hover:scale-105 transition-transform">
                  <ShieldCheck className="w-10 h-10 text-slate-300 group-hover:text-blue-400 transition-colors" />
                </div>
                <div>
                  <p className="text-slate-500 font-bold text-lg">暂无风险评估</p>
                  <p className="text-sm text-slate-400 mt-2 leading-relaxed">
                    点击上方按钮，获取 AI 深度<br/>
                    <span className="text-blue-500 font-medium">风险评估与操盘建议</span>
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default WatchlistDetail
