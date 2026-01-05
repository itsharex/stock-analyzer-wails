import { useState, useEffect } from 'react'
import { Loader2, CheckCircle2, XCircle, Play, X } from 'lucide-react'

interface BatchAnalyzeModalProps {
  isOpen: boolean
  onClose: () => void
  stocks: { code: string; name: string }[]
  onStart: (codes: string[]) => Promise<void>
}

interface ProgressEvent {
  code: string
  name: string
  completed: number
  total: number
  percent: number
}

function BatchAnalyzeModal({ isOpen, onClose, stocks, onStart }: BatchAnalyzeModalProps) {
  const [status, setStatus] = useState<'idle' | 'running' | 'completed' | 'error'>('idle')
  const [progress, setProgress] = useState<ProgressEvent | null>(null)
  const [currentStock, setCurrentStock] = useState<string>('')

  useEffect(() => {
    if (!isOpen) {
      setStatus('idle')
      setProgress(null)
      setCurrentStock('')
      return
    }

    // 监听来自后端的进度事件
    // @ts-ignore
    const unbind = window.runtime.EventsOn('batch_analyze_progress', (data: ProgressEvent) => {
      setProgress(data)
      setCurrentStock(data.name)
      if (data.completed === data.total) {
        setStatus('completed')
      }
    })

    return () => {
      if (unbind) unbind()
    }
  }, [isOpen])

  const handleStart = async () => {
    setStatus('running')
    try {
      await onStart(stocks.map(s => s.code))
    } catch (error) {
      console.error('批量分析失败:', error)
      setStatus('error')
    }
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md overflow-hidden">
        <div className="p-6 border-b border-slate-100 flex items-center justify-between">
          <h3 className="text-lg font-bold text-slate-900">批量 AI 深度分析</h3>
          <button onClick={onClose} className="p-1 hover:bg-slate-100 rounded-full transition-colors">
            <X className="w-5 h-5 text-slate-400" />
          </button>
        </div>

        <div className="p-8">
          {status === 'idle' && (
            <div className="text-center">
              <div className="w-16 h-16 bg-blue-50 rounded-full flex items-center justify-center mx-auto mb-4">
                <Play className="w-8 h-8 text-blue-500 fill-current" />
              </div>
              <p className="text-slate-600 mb-6">
                准备分析 <span className="font-bold text-blue-600">{stocks.length}</span> 只自选股。
                分析结果将自动缓存，提升后续查看速度。
              </p>
              <button
                onClick={handleStart}
                className="w-full py-3 bg-blue-600 text-white rounded-xl font-bold hover:bg-blue-700 transition-all shadow-lg shadow-blue-200"
              >
                开始批量分析
              </button>
            </div>
          )}

          {status === 'running' && (
            <div className="space-y-6">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium text-slate-500">正在分析: {currentStock}</span>
                <span className="text-sm font-bold text-blue-600">{progress?.percent.toFixed(0)}%</span>
              </div>
              <div className="w-full h-3 bg-slate-100 rounded-full overflow-hidden">
                <div 
                  className="h-full bg-blue-500 transition-all duration-500 ease-out"
                  style={{ width: `${progress?.percent || 0}%` }}
                />
              </div>
              <div className="flex items-center justify-center space-x-2 text-slate-400">
                <Loader2 className="w-4 h-4 animate-spin" />
                <span className="text-sm">正在处理第 {progress?.completed || 0}/{stocks.length} 只股票...</span>
              </div>
            </div>
          )}

          {status === 'completed' && (
            <div className="text-center">
              <div className="w-16 h-16 bg-emerald-50 rounded-full flex items-center justify-center mx-auto mb-4">
                <CheckCircle2 className="w-8 h-8 text-emerald-500" />
              </div>
              <h4 className="text-xl font-bold text-slate-900 mb-2">分析完成！</h4>
              <p className="text-slate-500 mb-6">所有自选股已完成深度分析并存入本地缓存。</p>
              <button
                onClick={onClose}
                className="w-full py-3 bg-slate-900 text-white rounded-xl font-bold hover:bg-slate-800 transition-all"
              >
                返回列表
              </button>
            </div>
          )}

          {status === 'error' && (
            <div className="text-center">
              <div className="w-16 h-16 bg-red-50 rounded-full flex items-center justify-center mx-auto mb-4">
                <XCircle className="w-8 h-8 text-red-500" />
              </div>
              <h4 className="text-xl font-bold text-slate-900 mb-2">分析中断</h4>
              <p className="text-slate-500 mb-6">处理过程中遇到错误，请检查网络或 AI 配置。</p>
              <button
                onClick={handleStart}
                className="w-full py-3 bg-red-600 text-white rounded-xl font-bold hover:bg-red-700 transition-all"
              >
                重试
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default BatchAnalyzeModal
