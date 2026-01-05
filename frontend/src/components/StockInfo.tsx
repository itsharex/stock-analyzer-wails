import type { StockData } from '../types'

interface StockInfoProps {
  stockData: StockData
}

function StockInfo({ stockData }: StockInfoProps) {
  const isRise = stockData.changeRate >= 0
  const changeRateText = `${isRise ? '+' : ''}${stockData.changeRate.toFixed(2)}%`
  const changeText = `${isRise ? '+' : ''}${stockData.change.toFixed(2)}`

  const formatNumber = (num: number, decimals: number = 2): string => {
    return num.toLocaleString('zh-CN', {
      minimumFractionDigits: decimals,
      maximumFractionDigits: decimals,
    })
  }

  const formatLargeNumber = (num: number): string => {
    if (num >= 100000000) {
      return `${(num / 100000000).toFixed(2)}亿`
    } else if (num >= 10000) {
      return `${(num / 10000).toFixed(2)}万`
    }
    return num.toFixed(2)
  }

  return (
    <div className="bg-white rounded-lg shadow-lg p-6 overflow-y-auto max-h-[calc(100vh-400px)]">
      <div className="mb-4">
        <div className="flex items-baseline justify-between mb-2">
          <h2 className="text-2xl font-bold text-gray-800">{stockData.name}</h2>
          <span className="text-sm text-gray-500">{stockData.code}</span>
        </div>
        
        <div className="flex items-baseline space-x-3">
          <span className={`text-3xl font-bold ${isRise ? 'text-red-500' : 'text-green-500'}`}>
            ¥{formatNumber(stockData.price)}
          </span>
          <span className={`text-lg font-medium ${isRise ? 'text-red-500' : 'text-green-500'}`}>
            {changeText}
          </span>
          <span className={`text-lg font-medium px-2 py-1 rounded ${isRise ? 'bg-red-100 text-red-600' : 'bg-green-100 text-green-600'}`}>
            {changeRateText}
          </span>
        </div>
      </div>

      <div className="space-y-3">
        <div className="grid grid-cols-2 gap-3">
          <div className="bg-gray-50 rounded-lg p-3">
            <div className="text-xs text-gray-500 mb-1">今开</div>
            <div className="text-base font-semibold text-gray-800">
              ¥{formatNumber(stockData.open)}
            </div>
          </div>
          
          <div className="bg-gray-50 rounded-lg p-3">
            <div className="text-xs text-gray-500 mb-1">昨收</div>
            <div className="text-base font-semibold text-gray-800">
              ¥{formatNumber(stockData.preClose)}
            </div>
          </div>
          
          <div className="bg-gray-50 rounded-lg p-3">
            <div className="text-xs text-gray-500 mb-1">最高</div>
            <div className="text-base font-semibold text-red-500">
              ¥{formatNumber(stockData.high)}
            </div>
          </div>
          
          <div className="bg-gray-50 rounded-lg p-3">
            <div className="text-xs text-gray-500 mb-1">最低</div>
            <div className="text-base font-semibold text-green-500">
              ¥{formatNumber(stockData.low)}
            </div>
          </div>
        </div>

        <div className="border-t border-gray-200 pt-3">
          <div className="grid grid-cols-2 gap-y-2 text-sm">
            <div className="flex justify-between px-2">
              <span className="text-gray-600">成交量</span>
              <span className="font-medium text-gray-800">
                {formatLargeNumber(stockData.volume)}手
              </span>
            </div>
            
            <div className="flex justify-between px-2">
              <span className="text-gray-600">成交额</span>
              <span className="font-medium text-gray-800">
                {formatLargeNumber(stockData.amount)}
              </span>
            </div>
            
            <div className="flex justify-between px-2">
              <span className="text-gray-600">振幅</span>
              <span className="font-medium text-gray-800">
                {stockData.amplitude.toFixed(2)}%
              </span>
            </div>
            
            <div className="flex justify-between px-2">
              <span className="text-gray-600">换手率</span>
              <span className="font-medium text-gray-800">
                {stockData.turnover.toFixed(2)}%
              </span>
            </div>
            
            <div className="flex justify-between px-2">
              <span className="text-gray-600">市盈率</span>
              <span className="font-medium text-gray-800">
                {stockData.pe.toFixed(2)}
              </span>
            </div>
            
            <div className="flex justify-between px-2">
              <span className="text-gray-600">市净率</span>
              <span className="font-medium text-gray-800">
                {stockData.pb.toFixed(2)}
              </span>
            </div>
            
            <div className="flex justify-between px-2">
              <span className="text-gray-600">总市值</span>
              <span className="font-medium text-gray-800">
                {formatLargeNumber(stockData.totalMV)}
              </span>
            </div>
            
            <div className="flex justify-between px-2">
              <span className="text-gray-600">流通市值</span>
              <span className="font-medium text-gray-800">
                {formatLargeNumber(stockData.circMV)}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default StockInfo
