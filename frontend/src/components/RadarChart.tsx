import {
  Radar,
  RadarChart as RechartsRadarChart,
  PolarGrid,
  PolarAngleAxis,
  ResponsiveContainer,
  Tooltip
} from 'recharts';
import { RadarData } from '../types';

interface RadarChartProps {
  data: RadarData;
}

const dimensionMap: Record<string, string> = {
  technical: '技术面',
  fundamental: '基本面',
  capital: '资金面',
  valuation: '估值面',
  risk: '风险面'
};

export default function RadarChart({ data }: RadarChartProps) {
  const chartData = Object.entries(data.scores).map(([key, value]) => ({
    subject: dimensionMap[key] || key,
    A: value,
    fullMark: 100,
    reason: data.reasons[key] || ''
  }));

  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const item = payload[0].payload;
      return (
        <div className="bg-slate-800 border border-slate-700 p-3 rounded-xl shadow-2xl max-w-[200px]">
          <p className="text-blue-400 font-bold text-sm mb-1">{item.subject}: {item.A}分</p>
          <p className="text-slate-300 text-xs leading-relaxed">{item.reason}</p>
        </div>
      );
    }
    return null;
  };

  return (
    <div className="w-full h-[280px] flex items-center justify-center">
      <ResponsiveContainer width="100%" height="100%">
        <RechartsRadarChart cx="50%" cy="50%" outerRadius="80%" data={chartData}>
          <PolarGrid stroke="#334155" />
          <PolarAngleAxis 
            dataKey="subject" 
            tick={{ fill: '#94a3b8', fontSize: 12, fontWeight: 600 }}
          />
          <Tooltip content={<CustomTooltip />} />
          <Radar
            name="股票评分"
            dataKey="A"
            stroke="#3b82f6"
            strokeWidth={2}
            fill="#3b82f6"
            fillOpacity={0.4}
          />
        </RechartsRadarChart>
      </ResponsiveContainer>
    </div>
  );
}
