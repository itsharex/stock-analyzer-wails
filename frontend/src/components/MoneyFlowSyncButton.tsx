import React, { useState, useEffect } from 'react';
import { useWailsAPI } from '../hooks/useWailsAPI';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import { Database, Loader2, CheckCircle2, AlertCircle } from 'lucide-react';

const MoneyFlowSyncButton: React.FC = () => {
  const { StartFullMarketSync } = useWailsAPI();
  const [syncing, setSyncing] = useState(false);
  const [progress, setProgress] = useState({ current: 0, total: 0, currentStock: '', success: 0, failed: 0 });
  const [status, setStatus] = useState<'idle' | 'syncing' | 'completed' | 'error'>('idle');
  const [errorMessage, setErrorMessage] = useState('');

  useEffect(() => {
    const onProgress = (data: any) => {
      // console.log('Sync progress:', data);
      setSyncing(true);
      setStatus('syncing');
      setProgress({
        current: data.current,
        total: data.total,
        currentStock: data.currentStock,
        success: data.successCount,
        failed: data.failedCount
      });

      if (data.status === 'completed') {
        setSyncing(false);
        setStatus('completed');
        setTimeout(() => setStatus('idle'), 5000);
      } else if (data.status === 'error') {
        setSyncing(false);
        setStatus('error');
        setErrorMessage(data.currentStock || '同步失败'); // error message might be in currentStock field for simplicity in backend
        setTimeout(() => setStatus('idle'), 5000);
      }
    };

    EventsOn('sync_progress', onProgress);

    return () => {
      EventsOff('sync_progress');
    };
  }, []);

  const handleSync = async () => {
    if (syncing) return;
    try {
      await StartFullMarketSync();
    } catch (e) {
      console.error('Failed to start sync:', e);
      setStatus('error');
      setErrorMessage('启动同步失败');
    }
  };

  if (status === 'syncing') {
    const percent = progress.total > 0 ? Math.round((progress.current / progress.total) * 100) : 0;
    return (
      <div className="flex flex-col items-end gap-1">
        <button disabled className="flex items-center gap-2 bg-purple-600/50 text-white/80 px-4 py-2 rounded-lg cursor-not-allowed border border-purple-500/30">
          <Loader2 className="w-4 h-4 animate-spin" />
          <span className="font-mono text-sm">{percent}%</span>
        </button>
        <span className="text-[10px] text-gray-400 font-mono">
          {progress.current}/{progress.total} | {progress.currentStock}
        </span>
      </div>
    );
  }

  if (status === 'completed') {
    return (
      <button className="flex items-center gap-2 bg-green-600 text-white px-4 py-2 rounded-lg shadow-lg shadow-green-900/20 animate-in fade-in zoom-in duration-300">
        <CheckCircle2 className="w-4 h-4" />
        <span>同步完成</span>
      </button>
    );
  }

  if (status === 'error') {
    return (
      <button 
        className="flex items-center gap-2 bg-red-600 text-white px-4 py-2 rounded-lg shadow-lg shadow-red-900/20"
        title={errorMessage}
      >
        <AlertCircle className="w-4 h-4" />
        <span>同步失败</span>
      </button>
    );
  }

  return (
    <button 
      onClick={handleSync}
      className="flex items-center gap-2 bg-purple-600 hover:bg-purple-500 text-white px-4 py-2 rounded-lg transition-all shadow-lg shadow-purple-900/20 hover:shadow-purple-600/30 active:scale-95"
    >
      <Database className="w-4 h-4" />
      <span>资金同步</span>
    </button>
  );
};

export default MoneyFlowSyncButton;
