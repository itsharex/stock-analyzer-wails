import React from 'react'
import { OrderBook } from '../types'
import { Wallet } from 'lucide-react'

interface OrderBookPanelProps {
  orderBook: OrderBook | undefined
}

const OrderBookPanel: React.FC<OrderBookPanelProps> = ({ orderBook }) => {
  if (!orderBook) {
    return (
      <div className="bg-white rounded-lg shadow-md p-4 col-span-1">
        <h3 className="text-lg font-semibold text-slate-800 mb-2 flex items-center">
          <Wallet className="w-5 h-5 mr-2 text-blue-500" /> 实时盘口
        </h3>
        <p className="text-sm text-slate-500">数据加载中...</p>
      </div>
    )
  }

  const formatVolume = (volume: number) => (volume / 100).toFixed(0) + '手'
  const formatAmount = (amount: number) => (amount / 10000).toFixed(2) + '万'

  return (
    <div className="bg-white rounded-lg shadow-md p-4 col-span-1">
      <h3 className="text-lg font-semibold text-slate-800 mb-2 flex items-center">
        <Wallet className="w-5 h-5 mr-2 text-blue-500" /> 实时盘口
      </h3>
      <div className="grid grid-cols-2 gap-2 text-xs mb-3">
        <div className="flex justify-between">
          <span className="text-slate-500">总成交量:</span>
          <span className="font-medium text-slate-700">{formatVolume(orderBook.volume)}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-slate-500">总成交额:</span>
          <span className="font-medium text-slate-700">{formatAmount(orderBook.amount)}</span>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-1 text-xs font-mono">
        {/* 卖盘 Sell */}
        {orderBook.sell.slice().reverse().map((item, index) => (
          <React.Fragment key={`sell-${index}`}>
            <div className="text-slate-500">卖{5 - index}</div>
            <div className="text-red-500 text-right">{item.price.toFixed(2)}</div>
            <div className="text-slate-600 text-right">{formatVolume(item.volume)}</div>
          </React.Fragment>
        ))}

        {/* 盘口分隔 */}
        <div className="col-span-3 h-px bg-slate-200 my-1" />

        {/* 买盘 Buy */}
        {orderBook.buy.map((item, index) => (
          <React.Fragment key={`buy-${index}`}>
            <div className="text-slate-500">买{index + 1}</div>
            <div className="text-green-500 text-right">{item.price.toFixed(2)}</div>
            <div className="text-slate-600 text-right">{formatVolume(item.volume)}</div>
          </React.Fragment>
        ))}
      </div>
    </div>
  )
}

export default OrderBookPanel
