import React, { useState, useEffect } from 'react';
import { useWailsAPI } from '../hooks/useWailsAPI';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import { Radar, Loader2, CheckCircle2, AlertCircle } from 'lucide-react';

const ScanButton: React.FC = () => {
  const { StartMassScan } = useWailsAPI();
  const [scanning, setScanning] = useState(false);
  const [progress, setProgress] = useState({ current: 0, total: 0, found: 0 });
  const [status, setStatus] = useState<'idle' | 'scanning' | 'completed' | 'error'>('idle');
  const [errorMessage, setErrorMessage] = useState('');

  useEffect(() => {
    // 监听扫描开始
    const onStart = (data: any) => {
      console.log('Scan started:', data);
      setScanning(true);
      setStatus('scanning');
      setProgress({ current: 0, total: data.total || 0, found: 0 });
    };

    // 监听扫描进度
    const onProgress = (data: any) => {
      // console.log('Scan progress:', data);
      setProgress(prev => ({
        ...prev,
        current: data.current,
        total: data.total,
        found: data.found
      }));
    };

    // 监听扫描完成
    const onComplete = (data: any) => {
      console.log('Scan complete:', data);
      setScanning(false);
      setStatus('completed');
      setProgress(prev => ({ ...prev, found: data.found }));
      
      // 3秒后重置为空闲状态
      setTimeout(() => setStatus('idle'), 3000);
    };

    // 监听扫描错误
    const onError = (msg: string) => {
      console.error('Scan error:', msg);
      setScanning(false);
      setStatus('error');
      setErrorMessage(msg);
      
      // 3秒后重置为空闲状态
      setTimeout(() => setStatus('idle'), 3000);
    };

    EventsOn('scan_start', onStart);
    EventsOn('scan_progress', onProgress);
    EventsOn('scan_complete', onComplete);
    EventsOn('scan_error', onError);

    return () => {
      EventsOff('scan_start');
      EventsOff('scan_progress');
      EventsOff('scan_complete');
      EventsOff('scan_error');
    };
  }, []);

  const handleScan = async () => {
    if (scanning) return;
    try {
      await StartMassScan();
    } catch (e) {
      console.error('Failed to start scan:', e);
      setStatus('error');
    }
  };

  if (status === 'scanning') {
    const percent = progress.total > 0 ? Math.round((progress.current / progress.total) * 100) : 0;
    return (
      <div className="flex flex-col items-end gap-1">
        <button disabled className="flex items-center gap-2 bg-blue-600/50 text-white/80 px-4 py-2 rounded-lg cursor-not-allowed border border-blue-500/30">
          <Loader2 className="w-4 h-4 animate-spin" />
          <span className="font-mono text-sm">{percent}%</span>
        </button>
        <span className="text-[10px] text-gray-400 font-mono">
          正在扫描 {progress.current}/{progress.total}...
        </span>
      </div>
    );
  }

  if (status === 'completed') {
    return (
      <button className="flex items-center gap-2 bg-green-600 text-white px-4 py-2 rounded-lg shadow-lg shadow-green-900/20 animate-in fade-in zoom-in duration-300">
        <CheckCircle2 className="w-4 h-4" />
        <span>已完成 (发现 {progress.found})</span>
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
        <span>扫描失败</span>
      </button>
    );
  }

  return (
    <button 
      onClick={handleScan}
      className="flex items-center gap-2 bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-lg transition-all shadow-lg shadow-blue-900/20 hover:shadow-blue-600/30 active:scale-95"
    >
      <Radar className="w-4 h-4" />
      <span>全市场扫描</span>
    </button>
  );
};

export default ScanButton;
