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

  // 兼容两种结构：
  // - Go 后端当前返回：buy5/sell5
  // - 旧前端类型：buy/sell + volume/amount
  const sells = Array.isArray((orderBook as any).sell5)
    ? (orderBook as any).sell5
    : Array.isArray((orderBook as any).sell)
      ? (orderBook as any).sell
      : []

  const buys = Array.isArray((orderBook as any).buy5)
    ? (orderBook as any).buy5
    : Array.isArray((orderBook as any).buy)
      ? (orderBook as any).buy
      : []

  const totalVolume = typeof (orderBook as any).volume === 'number' ? (orderBook as any).volume : null
  const totalAmount = typeof (orderBook as any).amount === 'number' ? (orderBook as any).amount : null

  return (
    <div className="bg-white rounded-lg shadow-md p-4 col-span-1">
      <h3 className="text-lg font-semibold text-slate-800 mb-2 flex items-center">
        <Wallet className="w-5 h-5 mr-2 text-blue-500" /> 实时盘口
      </h3>
      <div className="grid grid-cols-2 gap-2 text-xs mb-3">
        <div className="flex justify-between">
          <span className="text-slate-500">总成交量:</span>
          <span className="font-medium text-slate-700">{totalVolume != null ? formatVolume(totalVolume) : '--'}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-slate-500">总成交额:</span>
          <span className="font-medium text-slate-700">{totalAmount != null ? formatAmount(totalAmount) : '--'}</span>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-1 text-xs font-mono">
        {/* 卖盘 Sell */}
        {sells.slice().reverse().map((item: any, index: number) => (
          <React.Fragment key={`sell-${index}`}>
            <div className="text-slate-500">卖{5 - index}</div>
            <div className="text-red-500 text-right">{item.price.toFixed(2)}</div>
            <div className="text-slate-600 text-right">{formatVolume(item.volume)}</div>
          </React.Fragment>
        ))}

        {/* 盘口分隔 */}
        <div className="col-span-3 h-px bg-slate-200 my-1" />

        {/* 买盘 Buy */}
        {buys.map((item: any, index: number) => (
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
