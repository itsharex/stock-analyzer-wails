import { useState, useEffect } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { StockData, KLineData, TechnicalAnalysisResult, IntradayData, MoneyFlowResponse, HealthCheckResult, EntryStrategyResult, StockDetail } from '../types'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
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
import OrderBookPanel from './OrderBookPanel'
import FinancialPanel from './FinancialPanel'
import IndustryPanel from './IndustryPanel'
import VolumePriceAnalysis from './VolumePriceAnalysis'
import { useSmartSignalsModal } from '../hooks/useSmartSignalsModal'
import { 
  Activity, 
  BrainCircuit,
  Loader2,
  ShieldCheck,
  BarChart3,
  Cpu,
  Target,
  Wallet
} from 'lucide-react'
import { GlossaryTooltip } from './GlossaryTooltip'

interface WatchlistDetailProps {
  stock: StockData
}

export default function WatchlistDetail({ stock }: WatchlistDetailProps) {
  const { 
    getStockDetail, 
    getKLineData, 
    getIntradayData, 
    getMoneyFlowData, 
    analyzeTechnical, 
    getStockHealthCheck,
    analyzeEntryStrategy,
    streamIntradayData,
    stopIntradayStream
  } = useWailsAPI()

  const [stockDetail, setStockDetail] = useState<StockDetail | null>(null)
  const [klines, setKlines] = useState<KLineData[]>([])
  const [intraday, setIntraday] = useState<IntradayData[]>([])
  const [moneyFlow, setMoneyFlow] = useState<MoneyFlowResponse | null>(null)
  const [analysis, setAnalysis] = useState<TechnicalAnalysisResult | null>(null)
  const [healthCheck, setHealthCheck] = useState<HealthCheckResult | null>(null)
  const [entryStrategy, setEntryStrategy] = useState<EntryStrategyResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [analyzing, setAnalyzing] = useState(false)
  const [generatingStrategy, setGeneratingStrategy] = useState(false)
  
  const { 
    SmartSignalsModal: ModalComponent, 
    show: showSmartSignals, 
    klines: sKlines, 
    drawings: sDrawings, 
    moneyFlow: sMoneyFlow, 
    hide: sHide, 
    onAddAlert, 
    onCreatePosition 
  } = useSmartSignalsModal({
    code: stock.code,
    name: stock.name,
    price: stock.price
  })

  const handleAnalyze = async () => {
    if (analyzing) return
    setAnalyzing(true)
    try {
      // Use :force suffix to bypass cache and force re-analysis
      const res = await analyzeTechnical(stock.code, 'daily', 'technical:force')
      setAnalysis(res)
    } catch (e) {
      console.error(parseError(e))
    } finally {
      setAnalyzing(false)
    }
  }

  const handleEntryStrategy = async () => {
    if (generatingStrategy) return
    setGeneratingStrategy(true)
    try {
      const res = await analyzeEntryStrategy(stock.code)
      setEntryStrategy(res)
    } catch (e) {
      console.error(parseError(e))
    } finally {
      setGeneratingStrategy(false)
    }
  }

  useEffect(() => {
    let mounted = true
    const fetchData = async () => {
      setLoading(true)
      try {
        streamIntradayData(stock.code)

        const [detail, k, intra, mf, health] = await Promise.all([
          getStockDetail(stock.code),
          getKLineData(stock.code, 200, 'daily'),
          getIntradayData(stock.code),
          getMoneyFlowData(stock.code),
          getStockHealthCheck(stock.code)
        ])
        
        if (mounted) {
          setStockDetail(detail)
          // Sort klines ascending
          const sortedKlines = (k || []).sort((a, b) => a.time.localeCompare(b.time))
          setKlines(sortedKlines)
          setIntraday(intra?.data || [])
          setMoneyFlow(mf)
          setHealthCheck(health)
          
          // Try to fetch cached technical analysis
          analyzeTechnical(stock.code, 'daily').then(res => {
             if (mounted) setAnalysis(res)
          }).catch(console.error)
        }
      } catch (e) {
        console.error(parseError(e))
      } finally {
        if (mounted) setLoading(false)
      }
    }

    fetchData()

    const handleIntradayUpdate = (data: any) => {
       if (mounted && data && Array.isArray(data)) {
          // Parse and update logic would go here
       }
    }
    
    EventsOn("intradayDataUpdate:" + stock.code, handleIntradayUpdate)
    
    return () => {
      mounted = false
      EventsOff("intradayDataUpdate:" + stock.code)
      stopIntradayStream(stock.code)
    }
  }, [stock.code, analyzeEntryStrategy, analyzeTechnical, getIntradayData, getKLineData, getMoneyFlowData, getStockDetail, getStockHealthCheck, stopIntradayStream, streamIntradayData])

  if (loading && !stockDetail) {
    return (
      <div className="flex items-center justify-center h-96">
        <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
      </div>
    )
  }

  return (
    <div className="space-y-6 pb-20">
      {/* Signal Ticker */}
      <SignalTicker data={moneyFlow?.data || []} />

      {/* Top Cards Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {/* Intraday Chart */}
        <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-200">
           <div className="flex items-center justify-between mb-4">
             <h3 className="font-bold text-slate-700 flex items-center gap-2">
               <Activity className="w-5 h-5 text-blue-500" />
               分时走势
             </h3>
             <span className="text-xs font-mono text-slate-400">{stock.code}</span>
           </div>
           <div className="h-64">
             <IntradayChart data={intraday} preClose={stock.preClose} />
           </div>
        </div>

        {/* Money Flow */}
        <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-200">
           <div className="flex items-center justify-between mb-4">
             <h3 className="font-bold text-slate-700 flex items-center gap-2">
               <Wallet className="w-5 h-5 text-orange-500" />
               资金流向
             </h3>
             <GlossaryTooltip term="资金流向">
               <span className="text-xs text-slate-400 cursor-help">?</span>
             </GlossaryTooltip>
           </div>
           <div className="h-64">
             {moneyFlow && <MoneyFlowChart data={moneyFlow.data} />}
           </div>
        </div>

        {/* Stock Health */}
        <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-200">
           <div className="flex items-center justify-between mb-4">
             <h3 className="font-bold text-slate-700 flex items-center gap-2">
               <ShieldCheck className="w-5 h-5 text-green-500" />
               健康体检
             </h3>
           </div>
           {healthCheck && <StockHealthPanel data={healthCheck} />}
        </div>
      </div>

      {/* Main KLine Chart */}
      <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-200">
        <div className="flex items-center justify-between mb-4">
          <h3 className="font-bold text-slate-700 flex items-center gap-2">
            <BarChart3 className="w-5 h-5 text-indigo-500" />
            技术分析
          </h3>
          <div className="flex gap-2">
            <button 
              onClick={showSmartSignals}
              className="px-3 py-1.5 bg-indigo-50 text-indigo-600 rounded-lg text-sm font-medium hover:bg-indigo-100 flex items-center gap-2"
            >
              <BrainCircuit className="w-4 h-4" />
              智能信号
            </button>
          </div>
        </div>
        <div className="h-96">
           <KLineChart 
             data={klines} 
             drawings={analysis?.drawings || []}
             height={384}
           />
        </div>
        <div className="mt-4">
            <VolumePriceAnalysis stock={stock} />
        </div>
      </div>

      {/* Analysis Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Radar & Trade Plan */}
        <div className="space-y-6">
          <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-200">
             <h3 className="font-bold text-slate-700 mb-4 flex items-center gap-2">
               <Target className="w-5 h-5 text-purple-500" />
               综合评分
             </h3>
             <div className="h-64">
               {analysis?.radarData && <RadarChart data={analysis.radarData} />}
             </div>
          </div>
          
          {analysis?.tradePlan && (
            <TradePlanCard plan={analysis.tradePlan} currentPrice={stock.price} />
          )}
        </div>

        {/* AI Analysis Text */}
        <div className="lg:col-span-2 bg-white p-6 rounded-xl shadow-sm border border-slate-200">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-bold text-slate-700 flex items-center gap-2">
              <Cpu className="w-5 h-5 text-blue-600" />
              AI 深度分析
            </h3>
            <button 
              onClick={handleAnalyze}
              disabled={analyzing}
              className="px-3 py-1.5 bg-blue-50 text-blue-600 rounded-lg text-sm font-medium hover:bg-blue-100 disabled:opacity-50 transition-colors flex items-center gap-2"
            >
              {analyzing ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  正在分析...
                </>
              ) : (
                analysis ? '重新分析' : '开始分析'
              )}
            </button>
          </div>
          <div className="prose prose-slate max-w-none">
            {analysis ? (
              <ReactMarkdown remarkPlugins={[remarkGfm]}>
                {analysis.analysis}
              </ReactMarkdown>
            ) : analyzing ? (
              <div className="flex flex-col items-center justify-center py-12 text-slate-400">
                <Loader2 className="w-8 h-8 animate-spin mb-2 text-blue-500" />
                <p>AI 正在深度分析市场数据...</p>
              </div>
            ) : (
              <div className="text-slate-400 text-center py-12">
                点击“开始分析”获取 AI 深度研报
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Bottom Panels */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {entryStrategy ? (
          <EntryStrategyPanel strategy={entryStrategy} />
        ) : (
          <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200 flex flex-col items-center justify-center min-h-[300px]">
             <Target className="w-12 h-12 text-slate-300 mb-4" />
             <h3 className="font-bold text-slate-700 mb-2">AI 智能建仓方案</h3>
             <p className="text-sm text-slate-500 mb-6 text-center">
               基于技术面、资金面和基本面的多维度分析，<br/>生成科学的建仓与止损止盈策略。
             </p>
             <button 
               onClick={handleEntryStrategy}
               disabled={generatingStrategy}
               className="px-6 py-2 bg-blue-600 text-white rounded-lg font-bold text-sm hover:bg-blue-700 transition-all shadow-md shadow-blue-200 disabled:opacity-50 flex items-center gap-2"
             >
               {generatingStrategy ? (
                 <>
                   <Loader2 className="w-4 h-4 animate-spin" />
                   正在生成方案...
                 </>
               ) : (
                 '生成建仓方案'
               )}
             </button>
          </div>
        )}
        
        <div className="space-y-6">
          <OrderBookPanel orderBook={stockDetail?.orderBook} />
          <FinancialPanel financialSummary={stockDetail?.financial} />
          <IndustryPanel industryInfo={stockDetail?.industry} />
        </div>
      </div>

      {/* Modals */}
      {ModalComponent && (
        <ModalComponent 
          stock={stock}
          klines={sKlines}
          drawings={sDrawings}
          moneyFlow={sMoneyFlow}
          onClose={sHide}
          onAddAlert={onAddAlert}
          onCreatePosition={onCreatePosition}
        />
      )}
    </div>
  )
}
