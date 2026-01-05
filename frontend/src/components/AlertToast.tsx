import React, { useEffect, useState } from 'react';
import { X, TrendingUp, TrendingDown, Anchor, Sword, Cpu } from 'lucide-react';

interface AlertData {
  stockCode: string;
  stockName: string;
  message: string;
  advice: string;
  type: 'support' | 'resistance';
  price: number;
  role: 'conservative' | 'aggressive' | 'technical';
}

export const AlertToast: React.FC = () => {
  const [alerts, setAlerts] = useState<(AlertData & { id: number })[]>([]);

  useEffect(() => {
    // 监听来自 Wails 后端的 price_alert 事件
    const unsubscribe = (window as any).runtime?.EventsOn('price_alert', (data: AlertData) => {
      const id = Date.now();
      setAlerts(prev => [...prev, { ...data, id }]);
      
      // 10秒后自动移除
      setTimeout(() => {
        setAlerts(prev => prev.filter(a => a.id !== id));
      }, 10000);
    });

    return () => {
      if (unsubscribe) unsubscribe();
    };
  }, []);

  const removeAlert = (id: number) => {
    setAlerts(prev => prev.filter(a => a.id !== id));
  };

  const getRoleIcon = (role: string) => {
    switch (role) {
      case 'conservative': return <Anchor className="w-4 h-4 text-blue-600" />;
      case 'aggressive': return <Sword className="w-4 h-4 text-red-600" />;
      default: return <Cpu className="w-4 h-4 text-slate-600" />;
    }
  };

  const getRoleName = (role: string) => {
    switch (role) {
      case 'conservative': return '稳健老船长';
      case 'aggressive': return '激进先锋官';
      default: return '技术派大师';
    }
  };

  if (alerts.length === 0) return null;

  return (
    <div className="fixed bottom-6 right-6 z-[9999] flex flex-col gap-4 max-w-md w-full">
      {alerts.map((alert) => (
        <div 
          key={alert.id}
          className="bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden animate-in slide-in-from-right duration-300"
        >
          <div className={`h-1 w-full ${alert.type === 'resistance' ? 'bg-red-500' : 'bg-green-500'}`} />
          <div className="p-4">
            <div className="flex justify-between items-start mb-2">
              <div className="flex items-center space-x-2">
                <div className={`p-1.5 rounded-lg ${alert.type === 'resistance' ? 'bg-red-50' : 'bg-green-50'}`}>
                  {alert.type === 'resistance' ? 
                    <TrendingUp className={`w-4 h-4 ${alert.type === 'resistance' ? 'text-red-600' : 'text-green-600'}`} /> : 
                    <TrendingDown className="w-4 h-4 text-green-600" />
                  }
                </div>
                <span className="font-bold text-slate-800">{alert.stockName} ({alert.stockCode})</span>
              </div>
              <button 
                onClick={() => removeAlert(alert.id)}
                className="text-slate-400 hover:text-slate-600 transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            
            <p className="text-sm text-slate-600 mb-3 leading-relaxed">
              {alert.message}
            </p>
            
            <div className="bg-slate-50 rounded-xl p-3 border border-slate-100">
              <div className="flex items-center space-x-2 mb-1">
                {getRoleIcon(alert.role)}
                <span className="text-xs font-bold text-slate-500">{getRoleName(alert.role)} 的建议：</span>
              </div>
              <p className="text-sm font-medium text-blue-700 italic">
                "{alert.advice}"
              </p>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};
