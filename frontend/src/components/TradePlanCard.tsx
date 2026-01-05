import { TradePlan } from '../types';
import { Target, ShieldAlert, TrendingUp, PieChart, Info } from 'lucide-react';

interface TradePlanCardProps {
  plan: TradePlan;
  currentPrice: number;
}

export default function TradePlanCard({ plan, currentPrice }: TradePlanCardProps) {
  // 计算止损和止盈的百分比
  const stopLossPercent = ((plan.stopLoss - currentPrice) / currentPrice * 100).toFixed(2);
  const takeProfitPercent = ((plan.takeProfit - currentPrice) / currentPrice * 100).toFixed(2);

  return (
    <div className="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-sm hover:shadow-md transition-all">
      {/* 头部：建议仓位 */}
      <div className="bg-slate-50 px-4 py-3 border-b border-slate-200 flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <PieChart className="w-4 h-4 text-blue-600" />
          <span className="text-sm font-bold text-slate-800">智能仓位建议</span>
        </div>
        <div className="flex items-center space-x-1">
          <span className="text-lg font-black text-blue-600">{plan.suggestedPosition}</span>
          <span className="text-[10px] text-slate-400 font-bold uppercase tracking-tighter">仓位</span>
        </div>
      </div>

      <div className="p-4 space-y-4">
        {/* 价格阶梯 */}
        <div className="space-y-3">
          {/* 止盈 */}
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <div className="p-1 bg-green-50 rounded">
                <TrendingUp className="w-3.5 h-3.5 text-green-600" />
              </div>
              <span className="text-xs font-medium text-slate-600">建议止盈</span>
            </div>
            <div className="text-right">
              <span className="text-sm font-bold text-green-600">{plan.takeProfit.toFixed(2)}</span>
              <span className="text-[10px] text-green-500 ml-1">+{takeProfitPercent}%</span>
            </div>
          </div>

          {/* 现价指示器 */}
          <div className="relative h-1.5 bg-slate-100 rounded-full overflow-hidden">
            <div 
              className="absolute top-0 bottom-0 bg-blue-500 rounded-full"
              style={{ 
                left: '30%', 
                right: '30%' 
              }}
            />
          </div>

          {/* 止损 */}
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <div className="p-1 bg-red-50 rounded">
                <ShieldAlert className="w-3.5 h-3.5 text-red-600" />
              </div>
              <span className="text-xs font-medium text-slate-600">建议止损</span>
            </div>
            <div className="text-right">
              <span className="text-sm font-bold text-red-600">{plan.stopLoss.toFixed(2)}</span>
              <span className="text-[10px] text-red-500 ml-1">{stopLossPercent}%</span>
            </div>
          </div>
        </div>

        {/* 盈亏比与策略 */}
        <div className="pt-3 border-t border-slate-100 space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Target className="w-3.5 h-3.5 text-slate-400" />
              <span className="text-xs font-medium text-slate-500">预期盈亏比</span>
            </div>
            <span className={`text-xs font-bold ${plan.riskRewardRatio >= 2 ? 'text-blue-600' : 'text-slate-600'}`}>
              {plan.riskRewardRatio.toFixed(1)} : 1
            </span>
          </div>

          <div className="bg-blue-50/50 rounded-xl p-3 border border-blue-100/50">
            <div className="flex items-start space-x-2">
              <Info className="w-3.5 h-3.5 text-blue-500 mt-0.5 shrink-0" />
              <p className="text-xs text-blue-800 leading-relaxed font-medium">
                {plan.strategy}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
