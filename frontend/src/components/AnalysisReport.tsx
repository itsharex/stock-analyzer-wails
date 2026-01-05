import type { AnalysisReport } from '../types'

interface AnalysisReportProps {
  report: AnalysisReport
}

function AnalysisReport({ report }: AnalysisReportProps) {
  const getRiskLevelColor = (level: string): string => {
    if (level.includes('低')) return 'text-green-600 bg-green-50'
    if (level.includes('中')) return 'text-yellow-600 bg-yellow-50'
    if (level.includes('高')) return 'text-red-600 bg-red-50'
    return 'text-gray-600 bg-gray-50'
  }

  const getRecommendationIcon = (recommendation: string): string => {
    if (recommendation.includes('买入')) return '📈'
    if (recommendation.includes('持有')) return '🤝'
    if (recommendation.includes('观望')) return '👀'
    if (recommendation.includes('卖出')) return '📉'
    return '💡'
  }

  return (
    <div className="bg-white rounded-lg shadow-lg p-6 h-full overflow-y-auto">
      <div className="mb-6">
        <div className="flex items-center justify-between mb-2">
          <h2 className="text-2xl font-bold text-gray-800">AI分析报告</h2>
          <span className={`px-3 py-1 rounded-full text-sm font-medium ${getRiskLevelColor(report.riskLevel)}`}>
            {report.riskLevel}
          </span>
        </div>
        <div className="flex items-center justify-between text-sm text-gray-500">
          <div>
            <span className="font-medium">{report.stockName}</span>
            <span className="mx-2">|</span>
            <span>{report.stockCode}</span>
          </div>
          <div>{report.generatedAt}</div>
        </div>
      </div>

      <div className="space-y-6">
        {/* 摘要 */}
        {report.summary && (
          <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg p-4 border-l-4 border-blue-500">
            <h3 className="text-lg font-semibold text-gray-800 mb-2 flex items-center">
              <span className="mr-2">📊</span>
              分析摘要
            </h3>
            <p className="text-gray-700 leading-relaxed">{report.summary.trim()}</p>
          </div>
        )}

        {/* 基本面分析 */}
        {report.fundamentals && (
          <div className="border border-gray-200 rounded-lg p-4">
            <h3 className="text-lg font-semibold text-gray-800 mb-3 flex items-center">
              <span className="mr-2">📈</span>
              基本面分析
            </h3>
            <p className="text-gray-700 leading-relaxed whitespace-pre-wrap">{report.fundamentals.trim()}</p>
          </div>
        )}

        {/* 技术面分析 */}
        {report.technical && (
          <div className="border border-gray-200 rounded-lg p-4">
            <h3 className="text-lg font-semibold text-gray-800 mb-3 flex items-center">
              <span className="mr-2">📉</span>
              技术面分析
            </h3>
            <p className="text-gray-700 leading-relaxed whitespace-pre-wrap">{report.technical.trim()}</p>
          </div>
        )}

        {/* 投资建议 */}
        {report.recommendation && (
          <div className="bg-gradient-to-r from-purple-50 to-pink-50 rounded-lg p-4 border-l-4 border-purple-500">
            <h3 className="text-lg font-semibold text-gray-800 mb-3 flex items-center">
              <span className="mr-2">{getRecommendationIcon(report.recommendation)}</span>
              投资建议
            </h3>
            <p className="text-gray-700 leading-relaxed whitespace-pre-wrap">{report.recommendation.trim()}</p>
          </div>
        )}

        {/* 目标价位 */}
        {report.targetPrice && (
          <div className="bg-gradient-to-r from-green-50 to-emerald-50 rounded-lg p-4 border-l-4 border-green-500">
            <h3 className="text-lg font-semibold text-gray-800 mb-2 flex items-center">
              <span className="mr-2">🎯</span>
              目标价位
            </h3>
            <p className="text-gray-700 leading-relaxed">{report.targetPrice.trim()}</p>
          </div>
        )}
      </div>

      {/* 免责声明 */}
      <div className="mt-6 pt-6 border-t border-gray-200">
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
          <div className="flex items-start">
            <svg className="w-5 h-5 text-yellow-600 mr-2 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
            <div className="text-sm text-yellow-800">
              <p className="font-medium mb-1">免责声明</p>
              <p>本分析报告由AI生成，仅供参考，不构成任何投资建议。投资有风险，入市需谨慎。请根据自身情况谨慎决策，自行承担投资风险。</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default AnalysisReport
