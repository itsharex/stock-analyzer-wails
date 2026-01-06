import { useState, useEffect, useCallback } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { StockData, KLineData, TechnicalAnalysisResult, IntradayData, MoneyFlowResponse, HealthCheckResult, EntryStrategyResult, TrailingStopConfig } from '../types'
import { parseError } from '../utils/errorHandler'
import { useWailsAPI } from '../hooks/useWailsAPI'
import KLineChart from './KLineChart'
import IntradayChart from './IntradayChart'
import MoneyFlowChart from './MoneyFlowChart'
import SignalTicker from './SignalTicker'
import StockHealthPanel from './StockHealthPanel'
import RadarChart from './RadarChart'
import TradePlanCard from './TradePlanCard'
import EntryStrategyPanel from './EntryStrategyPanel'
import { 
  Activity, 
  Clock, 
  ChevronDown,
  BrainCircuit,
  Loader2,
  PencilRuler,
  ShieldCheck,
  Zap,
  BarChart3,
  Anchor,
  Sword,
  Cpu,
  TrendingUp,
  AlertTriangle,
  Info,
  Wallet,
  Target
} from 'lucide-react'
import { GlossaryPanel, GlossaryTooltip } from './GlossaryTooltip'
import { STOCK_GLOSSARY } from '../utils/glossary'

interface WatchlistDetailProps {
  stock: StockData
}

function WatchlistDetail({ stock }: WatchlistDetailProps) {
  const { getKLineData, analyzeTechnical, getIntradayData, getMoneyFlowData, getStockHealthCheck, analyzeEntryStrategy, addPosition } = useWailsAPI()
  const [klineData, setKlineData] = useState<KLineData[]>([])
  const [intradayData, setIntradayData] = useState<IntradayData[]>([])
  const [moneyFlowResponse, setMoneyFlowResponse] = useState<MoneyFlowResponse | null>(null)
  const [healthCheck, setHealthCheck] = useState<HealthCheckResult | null>(null)
  const [preClose, setPreClose] = useState<number>(0)
  const [chartType, setChartType] = useState<'intraday' | 'kline'>('intraday')
  const [period, setPeriod] = useState<string>('daily')
  const [loading, setLoading] = useState(false)
  const [analysisLoading, setAnalysisLoading] = useState(false)
  const [analysisResult, setAnalysisResult] = useState<TechnicalAnalysisResult | null>(null)
  const [entryStrategy, setEntryStrategy] = useState<EntryStrategyResult | null>(null)
  const [entryLoading, setEntryLoading] = useState(false)
  const [role, setRole] = useState('technical')

  const roles = [
    { id: 'technical', name: '技术派大师', icon: Cpu, color: 'text-blue-500', bg: 'bg-blue-50' },
    { id: 'conservative', name: '稳健老船长', icon: Anchor, color: 'text-emerald-600', bg: 'bg-emerald-50' },
    { id: 'aggressive', name: '激进先锋官', icon: Sword, color: 'text-rose-600', bg: 'bg-rose-50' },
  ]
  
  const [showMACD, setShowMACD] = useState(false)
  const [showKDJ, setShowKDJ] = useState(false)
  const [showRSI, setShowRSI] = useState(false)
  const [showAIDrawings, setShowAIDrawings] = useState(true)
  const [showMoneyFlow, setShowMoneyFlow] = useState(true)
  const [entryAnalysisStatus, setEntryAnalysisStatus] = useState<'idle' | 'checking' | 'analyzing' | 'error'>('idle')
  const [entryAnalysisError, setEntryAnalysisError] = useState<string>('')

  const loadKLineData = useCallback(async () => {
    if (chartType !== 'kline') return
    setLoading(true)
    try {
      const data = await getKLineData(stock.code, 100, period)
      setKlineData(data)
    } catch (error) {
      console.error('加载K线数据失败:', error)
    } finally {
      setLoading(false)
    }
  }, [stock.code, period, chartType, getKLineData])

  const loadIntradayData = useCallback(async () => {
    setLoading(true)
    try {
      const [intraResp, flowResp, healthResp] = await Promise.all([
        getIntradayData(stock.code),
        getMoneyFlowData(stock.code),
        getStockHealthCheck(stock.code)
      ])
      setIntradayData(intraResp.data)
      setPreClose(intraResp.preClose)
      setMoneyFlowResponse(flowResp)
      setHealthCheck(healthResp)
    } catch (error) {
      console.error('加载分时/资金流向数据失败:', error)
    } finally {
      setLoading(false)
    }
  }, [stock.code, getIntradayData, getMoneyFlowData, getStockHealthCheck])

  useEffect(() => {
    setAnalysisResult(null)
    setHealthCheck(null)
    setEntryStrategy(null)
    setEntryAnalysisStatus('idle')
    setEntryAnalysisError('')
  }, [stock.code])

  useEffect(() => {
    if (chartType === 'kline') {
      loadKLineData()
    } else {
      loadIntradayData()
    }
  }, [chartType, loadKLineData, loadIntradayData])

  useEffect(() => {
    if (chartType !== 'intraday') return
    const timer = setInterval(() => {
      loadIntradayData()
    }, 30000)
    return () => clearInterval(timer)
  }, [chartType, loadIntradayData])

  const handleAnalyze = async (selectedRole = role) => {
    setAnalysisLoading(true)
    try {
      const result = await analyzeTechnical(stock.code, period, selectedRole)
      setAnalysisResult(result)
      setShowAIDrawings(true)
    } catch (error) {
      console.error('技术分析失败:', error)
    } finally {
      setAnalysisLoading(false)
    }
  }

  const handleDeepAnalyzeClick = async () => {
    // 深度分析更适合在 K 线视图下展示（可看到 AI 绘图）
    if (chartType !== 'kline') {
      setChartType('kline')
    }
    return handleAnalyze(role)
  }

  const handleEntryAnalyze = async () => {
    // 数据完整性检查
    if (!intradayData || intradayData.length === 0) {
      setEntryAnalysisStatus('checking')
      setEntryAnalysisError('正在同步分时数据...')
      await loadIntradayData()
      setTimeout(() => {
        setEntryAnalysisStatus('idle')
        setEntryAnalysisError('')
      }, 1500)
      return
    }

    if (!moneyFlowResponse) {
      setEntryAnalysisStatus('checking')
      setEntryAnalysisError('正在同步资金流向数据...')
      await loadIntradayData()
      setTimeout(() => {
        setEntryAnalysisStatus('idle')
        setEntryAnalysisError('')
      }, 1500)
      return
    }

    setEntryLoading(true)
    setEntryAnalysisStatus('analyzing')
    setEntryAnalysisError('')
    
    try {
      const result = await analyzeEntryStrategy(stock.code)
      if (!result) {
        throw new Error('不幸的是，AI 引擎暂时无法给出建议，请稍后再试')
      }
      setEntryStrategy(result)
      setEntryAnalysisStatus('idle')
    } catch (error) {
      const parsed = parseError(error)
      const errorMsg = parsed.suggestion ? `${parsed.message}，建议：${parsed.suggestion}` : parsed.message
      console.error('建仓分析失败:', error)
      setEntryAnalysisStatus('error')
      setEntryAnalysisError(`分析失败：${errorMsg}`)
      setTimeout(() => {
        setEntryAnalysisStatus('idle')
        setEntryAnalysisError('')
      }, 3000)
    } finally {
      setEntryLoading(false)
    }
  }

  const handleConfirmEntry = async (config: TrailingStopConfig) => {
    if (!entryStrategy || !stock) return
    try {
      await addPosition({
        stockCode: stock.code,
        stockName: stock.name,
        entryPrice: stock.price,
        entryTime: new Date().toISOString(),
        strategy: entryStrategy,
        trailingConfig: config,
        currentStatus: 'holding',
        logicStatus: 'valid'
      })
      setEntryStrategy(null)
      alert('建仓成功，系统已开始实时监控逻辑！')
    } catch (error) {
      console.error('确认建仓失败:', error)
      alert('建仓失败，请重试')
    }
  }

  const getActionColor = (advice: string) => {
    switch (advice) {
      case '买入': case '增持': return 'bg-red-500 text-white'
      case '卖出': case '减持': return 'bg-green-500 text-white'
      default: return 'bg-slate-500 text-white'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case '主力建仓': return <TrendingUp className="w-4 h-4 text-red-500" />
      case '散户追高': return <AlertTriangle className="w-4 h-4 text-amber-500" />
      case '机构洗盘': return <Activity className="w-4 h-4 text-blue-500" />
      default: return <Info className="w-4 h-4 text-slate-500" />
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case '主力建仓': return 'bg-red-50 text-red-700 border-red-100'
      case '散户追高': return 'bg-amber-50 text-amber-700 border-amber-100'
      case '机构洗盘': return 'bg-blue-50 text-blue-700 border-blue-100'
      default: return 'bg-slate-50 text-slate-700 border-slate-100'
    }
  }

  return (
    <div className="flex flex-col h-full bg-slate-50 overflow-hidden">
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

      <div className="flex-1 flex overflow-hidden">
        <div className="flex-1 flex flex-col p-2 lg:p-4 overflow-y-auto">
          <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-4 flex-1 flex flex-col">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center space-x-2">
                <div className="flex bg-slate-100 rounded-lg p-1 mr-2">
                  <button 
                    onClick={() => setChartType('intraday')}
                    className={`px-4 py-1.5 text-sm font-bold rounded-md transition-all ${chartType === 'intraday' ? 'bg-white text-blue-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
                  >
                    分时
                  </button>
                  <button 
                    onClick={() => setChartType('kline')}
                    className={`px-4 py-1.5 text-sm font-bold rounded-md transition-all ${chartType === 'kline' ? 'bg-white text-blue-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
                  >
                    K线
                  </button>
                </div>

                {chartType === 'kline' ? (
                  <>
                    <div className="relative">
                      <select 
                        value={period}
                        onChange={(e) => setPeriod(e.target.value)}
                        className="appearance-none bg-slate-100 border-none rounded-lg px-4 py-2 pr-10 text-sm font-medium text-slate-700 focus:ring-2 focus:ring-blue-500 cursor-pointer"
                      >
                        <option value="daily">日线</option>
                        <option value="week">周线</option>
                        <option value="month">月线</option>
                      </select>
                      <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                    </div>
                    <div className="h-6 w-px bg-slate-200 mx-2" />
                    <div className="flex bg-slate-100 rounded-lg p-1">
                      <button onClick={() => setShowMACD(!showMACD)} className={`px-3 py-1 text-xs font-bold rounded-md transition-all ${showMACD ? 'bg-white text-blue-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}>MACD</button>
                      <button onClick={() => setShowKDJ(!showKDJ)} className={`px-3 py-1 text-xs font-bold rounded-md transition-all ${showKDJ ? 'bg-white text-purple-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}>KDJ</button>
                      <button onClick={() => setShowRSI(!showRSI)} className={`px-3 py-1 text-xs font-bold rounded-md transition-all ${showRSI ? 'bg-white text-cyan-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}>RSI</button>
                    </div>
                    {analysisResult && (
                      <button onClick={() => setShowAIDrawings(!showAIDrawings)} className={`flex items-center space-x-1 px-3 py-1 text-xs font-bold rounded-md transition-all ${showAIDrawings ? 'bg-blue-50 text-blue-600 border border-blue-200' : 'bg-slate-100 text-slate-500'}`}>
                        <PencilRuler className="w-3 h-3" />
                        <span>AI 绘图</span>
                      </button>
                    )}
                  </>
                ) : (
                  <div className="flex bg-slate-100 rounded-lg p-1">
                    <button onClick={() => setShowMoneyFlow(!showMoneyFlow)} className={`px-3 py-1 text-xs font-bold rounded-md transition-all ${showMoneyFlow ? 'bg-white text-blue-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}>资金流向</button>
                  </div>
                )}
              </div>
              <div className="flex items-center text-slate-400 text-xs space-x-4">
                <span className="flex items-center"><Clock className="w-3 h-3 mr-1" /> {chartType === 'intraday' ? '30秒自动刷新' : '实时更新'}</span>
                <span className="flex items-center"><Activity className="w-3 h-3 mr-1" /> 东方财富数据源</span>
              </div>
            </div>

            {chartType === 'intraday' && moneyFlowResponse && (
              <div className="mb-4 rounded-lg overflow-hidden border border-slate-800">
                <SignalTicker data={moneyFlowResponse.data} />
              </div>
            )}

            {healthCheck && (
              <div className="mb-6">
                <StockHealthPanel data={healthCheck} />
              </div>
            )}

            {chartType === 'intraday' && moneyFlowResponse && (
              <div className={`mb-4 p-3 rounded-xl border flex items-start space-x-3 transition-all ${getStatusColor(moneyFlowResponse.status)}`}>
                <div className="mt-0.5">{getStatusIcon(moneyFlowResponse.status)}</div>
                <div>
                  <div className="flex items-center space-x-2">
                    <span className="font-bold text-sm">智能资金识别：{moneyFlowResponse.status}</span>
                    <span className="text-[10px] px-1.5 py-0.5 rounded-md bg-white/50 font-medium">实时监控中</span>
                  </div>
                  <p className="text-xs mt-1 opacity-90 leading-relaxed">{moneyFlowResponse.description}</p>
                </div>
              </div>
            )}

            <div className="flex-1 flex flex-col space-y-2 min-h-[500px]">
              <div className="flex-[2] relative bg-slate-50 rounded-lg border border-slate-100 overflow-hidden">
                {loading && (
                  <div className="absolute inset-0 flex items-center justify-center bg-white/80 z-10">
                    <Loader2 className="w-8 h-8 text-blue-500 animate-spin" />
                  </div>
                )}
                {chartType === 'kline' ? (
                  <KLineChart data={klineData} drawings={showAIDrawings ? analysisResult?.drawings : []} showMACD={showMACD} showKDJ={showKDJ} showRSI={showRSI} />
                ) : (
                  <IntradayChart data={intradayData} moneyFlowData={moneyFlowResponse?.data} preClose={preClose} height={400} />
                )}
              </div>
              {chartType === 'intraday' && showMoneyFlow && moneyFlowResponse && (
                <div className="flex-1 relative bg-slate-50 rounded-lg border border-slate-100 overflow-hidden">
                  <MoneyFlowChart data={moneyFlowResponse.data} height={180} />
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="w-[480px] bg-slate-50 border-l border-slate-200 flex flex-col shadow-[-10px_0_30px_rgba(0,0,0,0.03)] z-10">
          <div className="p-5 border-b border-slate-200 flex flex-col space-y-4 bg-white/50 backdrop-blur-sm">
            {chartType === 'intraday' && moneyFlowResponse && (
              <div className="grid grid-cols-2 gap-3 mb-2">
                <div className="bg-white p-3 rounded-xl border border-slate-100 shadow-sm">
                  <div className="flex items-center space-x-2 text-slate-400 mb-1">
                    <Wallet className="w-3 h-3" />
                    <span className="text-[10px] font-bold uppercase tracking-wider">主力净流入</span>
                  </div>
                  <div className={`text-lg font-black ${moneyFlowResponse.todayMain >= 0 ? 'text-red-500' : 'text-green-500'}`}>
                    {moneyFlowResponse.todayMain >= 0 ? '+' : ''}{(moneyFlowResponse.todayMain / 10000).toFixed(2)}万
                  </div>
                </div>
                <div className="bg-white p-3 rounded-xl border border-slate-100 shadow-sm">
                  <div className="flex items-center space-x-2 text-slate-400 mb-1">
                    <Target className="w-3 h-3" />
                    <span className="text-[10px] font-bold uppercase tracking-wider">主力占比</span>
                  </div>
                  <div className="text-lg font-black text-slate-700">
                    {((moneyFlowResponse.todayMain / (Math.abs(moneyFlowResponse.todayMain) + Math.abs(moneyFlowResponse.todayRetail) || 1)) * 100).toFixed(1)}%
                  </div>
                </div>
              </div>
            )}

            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <div className="p-2 bg-blue-600 rounded-xl shadow-lg shadow-blue-200">
                  <BrainCircuit className="w-5 h-5 text-white" />
                </div>
                <div>
                  <h3 className="text-sm font-black text-slate-800">AI 智能分析师</h3>
                  <p className="text-[10px] text-slate-400 font-medium">DeepSeek-V3 引擎驱动</p>
                </div>
              </div>
              <div className="flex-1">
                <div className="grid grid-cols-2 gap-2">
                  <button
                    onClick={handleDeepAnalyzeClick}
                    disabled={analysisLoading || entryLoading}
                    className="w-full flex items-center justify-center space-x-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-bold transition-all shadow-lg shadow-blue-100 disabled:opacity-50"
                    title={chartType !== 'kline' ? '将自动切换到K线视图以展示AI绘图' : ''}
                  >
                    {analysisLoading ? (
                      <>
                        <Loader2 className="w-3 h-3 animate-spin" />
                        <span>深度分析中...</span>
                      </>
                    ) : (
                      <>
                        <Cpu className="w-3 h-3" />
                        <span>深度分析</span>
                      </>
                    )}
                  </button>
                  <button
                    onClick={handleEntryAnalyze}
                    disabled={entryLoading || analysisLoading}
                    className="w-full flex items-center justify-center space-x-2 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl text-xs font-bold transition-all shadow-lg shadow-indigo-100 disabled:opacity-50"
                  >
                    {entryLoading ? (
                      <>
                        <Loader2 className="w-3 h-3 animate-spin" />
                        <span>建仓分析中...</span>
                      </>
                    ) : (
                      <>
                        <Target className="w-3 h-3" />
                        <span>建仓分析</span>
                      </>
                    )}
                  </button>
                </div>
                {/* 状态提示文字 */}
                {entryAnalysisStatus !== 'idle' && (
                  <div className={`mt-2 text-xs text-center font-medium px-2 py-1 rounded-lg transition-all ${
                    entryAnalysisStatus === 'checking' ? 'text-blue-600 bg-blue-50' :
                    entryAnalysisStatus === 'analyzing' ? 'text-indigo-600 bg-indigo-50' :
                    'text-red-600 bg-red-50'
                  }`}>
                    {entryAnalysisError}
                  </div>
                )}
              </div>
            </div>

            <div className="flex p-1 bg-slate-100 rounded-2xl border border-slate-200">
              {roles.map((r) => (
                <button
                  key={r.id}
                  onClick={() => {
                    setRole(r.id);
                    handleAnalyze(r.id);
                  }}
                  disabled={analysisLoading}
                  className={`flex-1 flex items-center justify-center space-x-2 py-2 rounded-xl transition-all ${
                    role === r.id 
                      ? 'bg-white shadow-sm text-slate-800' 
                      : 'text-slate-500 hover:text-slate-700'
                  }`}
                >
                  <r.icon className={`w-3.5 h-3.5 ${role === r.id ? r.color : 'text-slate-400'}`} />
                  <span className={`text-[11px] font-bold ${role === r.id ? '' : 'opacity-70'}`}>{r.name}</span>
                </button>
              ))}
            </div>
          </div>

          <div className="flex-1 overflow-y-auto p-6 custom-scrollbar bg-gradient-to-b from-white/30 to-transparent">
            {entryStrategy && (
              <div className="mt-4 mb-8">
                <EntryStrategyPanel 
                  strategy={entryStrategy} 
                  onConfirm={handleConfirmEntry}
                />
              </div>
            )}

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
                {analysisResult.radarData && analysisResult.radarData.scores && Object.keys(analysisResult.radarData.scores).length > 0 && (
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

                {analysisResult.tradePlan && (
                  <TradePlanCard plan={analysisResult.tradePlan} currentPrice={stock.price} />
                )}

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
