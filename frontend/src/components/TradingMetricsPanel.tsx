import React from 'react'
import { StockDetail } from '../types'
import { TrendingUp, Clock, BarChart2 } from 'lucide-react'

interface TradingMetricsPanelProps {
  stockDetail: StockDetail | undefined
}

const TradingMetricsPanel: React.FC<TradingMetricsPanelProps> = ({ stockDetail }) => {
  if (!stockDetail || !stockDetail.stockData) {
    return (
      <div className="p-4 bg-gray-800 rounded-lg shadow-lg">
        <h3 className="text-lg font-semibold text-gray-300 mb-2 flex items-center">
          <BarChart2 className="w-5 h-5 mr-2 text-blue-400" />
          实时交易指标
        </h3>
        <p className="text-gray-500">暂无实时数据</p>
      </div>
    )
  }

  const quote = stockDetail.stockData

  const formatNumber = (num: number | undefined, unit: string = '', precision: number = 2): string => {
    if (num === undefined || isNaN(num)) return '--'
    return `${num.toFixed(precision)}${unit}`
  }

  const dataPoints = [
    { label: '量比', value: formatNumber(quote.volumeRatio, 'x'), icon: <BarChart2 className="w-4 h-4 text-blue-400" /> },
    { label: '委比', value: formatNumber(quote.warrantRatio, '%'), icon: <BarChart2 className="w-4 h-4 text-blue-400" /> },
    { label: '换手率', value: formatNumber(quote.turnover, '%'), icon: <Clock className="w-4 h-4 text-yellow-400" /> },
    { label: '振幅', value: formatNumber(quote.amplitude * 100, '%'), icon: <TrendingUp className="w-4 h-4 text-green-400" /> },
    { label: '成交额', value: formatNumber(quote.amount / 100000000, '亿', 2), icon: <TrendingUp className="w-4 h-4 text-green-400" /> },
    { label: '总市值', value: formatNumber(quote.totalMV / 100000000, '亿', 2), icon: <BarChart2 className="w-4 h-4 text-blue-400" /> },
  ]

  return (
    <div className="p-4 bg-gray-800 rounded-lg shadow-lg">
      <h3 className="text-lg font-semibold text-gray-300 mb-4 flex items-center">
        <BarChart2 className="w-5 h-5 mr-2 text-blue-400" />
        实时交易指标
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
    </div>
  )
}

export default TradingMetricsPanel
