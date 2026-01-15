import React, { useEffect, useRef, useState } from 'react';
import { createChart, ColorType, IChartApi, SeriesMarker, Time } from 'lightweight-charts';
import { useWailsAPI } from '../hooks/useWailsAPI';
import { StrategySignal, KLineData, MoneyFlowData } from '../types';
import { useResizeObserver } from '../hooks/useResizeObserver';

interface EnhancedKLineChartProps {
  stockCode: string;
  signal: StrategySignal | null;
  onSignalClick?: (signal: StrategySignal) => void;
}

const EnhancedKLineChart: React.FC<EnhancedKLineChartProps> = React.memo(({ stockCode, onSignalClick }) => {
  const { getKLineData, getMoneyFlowData, GetSignalsByStockCode } = useWailsAPI();
  const chartContainerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<IChartApi | null>(null);
  const [klines, setKlines] = useState<KLineData[]>([]);
  const [moneyFlows, setMoneyFlows] = useState<MoneyFlowData[]>([]);
  const [historySignals, setHistorySignals] = useState<StrategySignal[]>([]);
  const [loading, setLoading] = useState(false);

  // Fetch Data
  useEffect(() => {
    if (!stockCode) return;
    
    const fetchData = async () => {
      setLoading(true);
      try {
        const [kData, mfData, sigData] = await Promise.all([
          getKLineData(stockCode, 100, 'daily'),
          getMoneyFlowData(stockCode),
          GetSignalsByStockCode(stockCode)
        ]);
        setKlines(kData || []);
        setMoneyFlows(mfData.data || []);
        setHistorySignals(sigData || []);
      } catch (error) {
        console.error("Failed to fetch chart data", error);
      } finally {
        setLoading(false);
      }
    };
    
    fetchData();
  }, [stockCode, getKLineData, getMoneyFlowData, GetSignalsByStockCode]);

  // Render Chart
  useEffect(() => {
    if (!chartContainerRef.current || klines.length === 0) return;

    const chart = createChart(chartContainerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: 'transparent' },
        textColor: '#9ca3af', // gray-400
      },
      grid: {
        vertLines: { color: '#1f2937' }, // gray-800
        horzLines: { color: '#1f2937' },
      },
      width: chartContainerRef.current.clientWidth,
      height: chartContainerRef.current.clientHeight || 500,
      timeScale: {
        borderColor: '#374151',
        timeVisible: true,
      },
      rightPriceScale: {
        borderColor: '#374151',
      },
    });

    // 1. Candlestick Series
    const candlestickSeries = chart.addCandlestickSeries({
      upColor: '#ef4444',
      downColor: '#22c55e',
      borderVisible: false,
      wickUpColor: '#ef4444',
      wickDownColor: '#22c55e',
    });

    const candleData = klines.map(d => ({
      time: d.time,
      open: d.open,
      high: d.high,
      low: d.low,
      close: d.close,
    }));
    candlestickSeries.setData(candleData);

    // 2. Markers (History Signals)
    if (historySignals.length > 0) {
      const markers: SeriesMarker<Time>[] = [];
      const klineTimes = new Set(klines.map(k => k.time));

      historySignals.forEach(sig => {
        // Only add marker if the date exists in current K-line data
        if (klineTimes.has(sig.tradeDate)) {
          markers.push({
            time: sig.tradeDate,
            position: 'belowBar',
            color: '#ef4444',
            shape: 'arrowUp',
            text: 'B',
            size: 2,
            // @ts-ignore - custom property for click handling
            signalId: sig.id
          });
        }
      });
      candlestickSeries.setMarkers(markers);
    }

    // 3. Money Flow Histogram (Sub-chart)
    const mfSeries = chart.addHistogramSeries({
      color: '#94a3b8',
      priceFormat: { type: 'volume' },
      priceScaleId: 'moneyFlow',
      title: '主力净流入(万)',
    });

    chart.priceScale('moneyFlow').applyOptions({
      scaleMargins: { top: 0.8, bottom: 0 },
    });

    // Align money flow data with klines
    // Need to map klines time to money flows
    const mfMap = new Map(moneyFlows.map(m => [m.time, m.mainNet]));
    
    const mfChartData = klines.map(k => {
      const net = mfMap.get(k.time) || 0;
      return {
        time: k.time,
        value: net / 10000, // Convert to Wan
        color: net > 0 ? '#ef444480' : '#22c55e80', // Red for inflow, Green for outflow
      };
    });
    
    mfSeries.setData(mfChartData);

    chart.timeScale().fitContent();

    // Click Event Handler
    chart.subscribeClick((param) => {
      if (!param || !param.seriesData) return;
      // param.seriesData is a Map<ISeriesApi<SeriesType>, SeriesData<SeriesType>>
      // We need to check if the click was near a marker on the candlestick series
      
      // Lightweight charts doesn't have a direct "click marker" event.
      // But we can check the time of the click and see if there's a signal on that date.
      if (param.time && historySignals.length > 0) {
        const clickedTime = param.time as string;
        const clickedSignal = historySignals.find(s => s.tradeDate === clickedTime);
        if (clickedSignal && onSignalClick) {
          onSignalClick(clickedSignal);
        }
      }
    });

    chartRef.current = chart;

    return () => {
      chart.remove();
    };
  }, [klines, moneyFlows, historySignals, onSignalClick]);

  // Resize Handler
  const onSize = (w: number, h: number) => {
    if (chartRef.current) {
      chartRef.current.applyOptions({ width: w, height: h });
      chartRef.current.timeScale().fitContent();
    }
  };
  
  useResizeObserver(chartContainerRef, onSize);

  return (
    <div className="flex-1 flex flex-col min-w-0 bg-[#0D1117] relative">
      {loading && (
        <div className="absolute inset-0 flex items-center justify-center z-10 bg-[#0D1117]/50 backdrop-blur-sm">
          <div className="text-blue-400">Loading Chart...</div>
        </div>
      )}
      <div className="flex-1 w-full p-4" ref={chartContainerRef}></div>
    </div>
  );
});

export default EnhancedKLineChart;
