import { HealthCheckResult } from '../types'
import { ShieldCheck, AlertTriangle, XCircle, Info, Activity, TrendingUp, DollarSign, PieChart } from 'lucide-react'

interface StockHealthPanelProps {
  data: HealthCheckResult
}

function StockHealthPanel({ data }: StockHealthPanelProps) {
  const getStatusIcon = (status: string) => {
    switch (status) {
      case '正常': return <ShieldCheck className="w-4 h-4 text-emerald-500" />
      case '警告': return <AlertTriangle className="w-4 h-4 text-amber-500" />
      case '异常': return <XCircle className="w-4 h-4 text-red-500" />
      default: return <Info className="w-4 h-4 text-slate-400" />
    }
  }

  const getCategoryIcon = (category: string) => {
    switch (category) {
      case '财务': return <DollarSign className="w-4 h-4" />
      case '资金': return <PieChart className="w-4 h-4" />
      case '技术': return <TrendingUp className="w-4 h-4" />
      default: return <Activity className="w-4 h-4" />
    }
  }

  const getScoreColor = (score: number) => {
    if (score >= 85) return 'text-emerald-500'
    if (score >= 60) return 'text-amber-500'
    return 'text-red-500'
  }



  return (
    <div className="bg-white rounded-2xl border border-slate-200 overflow-hidden shadow-sm">
      {/* 顶部评分区 */}
      <div className="p-6 border-b border-slate-100 bg-slate-50/50 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <div className={`w-16 h-16 rounded-full border-4 flex items-center justify-center bg-white shadow-inner ${
            data.score >= 85 ? 'border-emerald-500' : data.score >= 60 ? 'border-amber-500' : 'border-red-500'
          }`}>
            <span className={`text-2xl font-black ${getScoreColor(data.score)}`}>{data.score}</span>
          </div>
          <div>
            <h3 className="text-lg font-bold text-slate-900">AI 深度体检报告</h3>
            <div className="flex items-center space-x-2 mt-1">
              <span className={`px-2 py-0.5 rounded text-xs font-bold ${
                data.riskLevel === '低' ? 'bg-emerald-100 text-emerald-700' : 
                data.riskLevel === '中' ? 'bg-amber-100 text-amber-700' : 'bg-red-100 text-red-700'
              }`}>
                风险等级：{data.riskLevel}
              </span>
              <span className="text-slate-400 text-xs">更新于 {data.updatedAt}</span>
            </div>
          </div>
        </div>
        <div className="text-right">
          <div className={`text-sm font-bold ${getScoreColor(data.score)}`}>{data.status}状态</div>
        </div>
      </div>

      {/* AI 总结 */}
      <div className="p-4 bg-blue-50/50 border-b border-blue-100">
        <p className="text-sm text-blue-800 leading-relaxed">
          <span className="font-bold mr-2">AI 诊断：</span>
          {data.summary}
        </p>
      </div>

      {/* 体检子项列表 */}
      <div className="divide-y divide-slate-100">
        {data.items.map((item, index) => (
          <div key={index} className="p-4 hover:bg-slate-50 transition-colors">
            <div className="flex items-start justify-between mb-1">
              <div className="flex items-center space-x-2">
                <div className="p-1.5 bg-slate-100 rounded-lg text-slate-500">
                  {getCategoryIcon(item.category)}
                </div>
                <span className="text-sm font-bold text-slate-700">{item.name}</span>
                <span className="text-xs text-slate-400">({item.category})</span>
              </div>
              <div className="flex items-center space-x-1.5">
                <span className="text-sm font-mono font-medium text-slate-600">{item.value}</span>
                {getStatusIcon(item.status)}
              </div>
            </div>
            <p className="text-xs text-slate-500 ml-10 leading-relaxed">
              {item.description}
            </p>
          </div>
        ))}
      </div>

      {/* 底部操作建议 */}
      <div className="p-4 bg-slate-50 flex items-center justify-between">
        <div className="flex items-center space-x-2 text-xs text-slate-400">
          <Info className="w-3.5 h-3.5" />
          <span>本报告由 AI 自动生成，仅供参考，不构成投资建议。</span>
        </div>
      </div>
    </div>
  )
}

export default StockHealthPanel
