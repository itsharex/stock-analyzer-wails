import React from 'react'
import { FinancialSummary } from '../types'
import { DollarSign, TrendingUp, Scale } from 'lucide-react'

interface FinancialPanelProps {
  financialSummary: FinancialSummary | undefined
}

const FinancialPanel: React.FC<FinancialPanelProps> = ({ financialSummary }) => {
  if (!financialSummary) {
    return (
      <div className="p-4 bg-gray-800 rounded-lg shadow-lg">
        <h3 className="text-lg font-semibold text-gray-300 mb-2 flex items-center">
          <DollarSign className="w-5 h-5 mr-2 text-green-400" />
          财务摘要
        </h3>
        <p className="text-gray-500">暂无财务数据</p>
      </div>
    )
  }

  const formatNumber = (num: number | undefined, unit: string = '', precision: number = 2): string => {
    if (num === undefined || isNaN(num)) return '--'
    // 假设 marketCap 是以“亿”为单位的，其他是百分比或倍数
    if (unit === '亿') {
      return `${num.toFixed(precision)}${unit}`
    }
    return `${num.toFixed(precision)}${unit}`
  }

  const dataPoints = [
    
    { label: '净资产收益率 (ROE)', value: formatNumber(financialSummary.roe, '%'), icon: <TrendingUp className="w-4 h-4 text-green-400" /> },
    { label: '净利润增长率', value: formatNumber(financialSummary.net_profit_growth_rate, '%'), icon: <TrendingUp className="w-4 h-4 text-green-400" /> },
    { label: '毛利率', value: formatNumber(financialSummary.gross_profit_margin, '%'), icon: <TrendingUp className="w-4 h-4 text-green-400" /> },
    { label: '股息率', value: formatNumber(financialSummary.dividend_yield, '%'), icon: <DollarSign className="w-4 h-4 text-yellow-400" /> },
    { label: '总市值', value: formatNumber(financialSummary.total_market_value, '亿', 2), icon: <Scale className="w-4 h-4 text-blue-400" /> },
    { label: '流通市值', value: formatNumber(financialSummary.circulating_market_value, '亿', 2), icon: <Scale className="w-4 h-4 text-blue-400" /> },
  ]

  return (
    <div className="p-4 bg-gray-800 rounded-lg shadow-lg">
      <h3 className="text-lg font-semibold text-gray-300 mb-4 flex items-center">
        <DollarSign className="w-5 h-5 mr-2 text-green-400" />
        财务摘要
      </h3>
      <div className="grid grid-cols-2 gap-4">
        {dataPoints.map((item, index) => (
          <div key={index} className="flex justify-between items-center p-2 bg-gray-700 rounded-md">
            <span className="text-sm text-gray-400 flex items-center">
              {item.icon}
              <span className="ml-2">{item.label}</span>
            </span>
            <span className="text-base font-medium text-white">{item.value}</span>
          </div>
        ))}
      </div>
      <p className="text-xs text-gray-500 mt-4">数据来源: 东方财富 (Placeholder)</p>
    </div>
  )
}

export default FinancialPanel
