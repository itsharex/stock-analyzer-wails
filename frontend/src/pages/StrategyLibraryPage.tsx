import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, Play, TrendingUp, TrendingDown, Activity } from 'lucide-react';
import StrategyEditor from '../components/StrategyEditor';
import { useWailsAPI } from '../hooks/useWailsAPI';

interface StrategyConfig {
  id: number;
  name: string;
  description: string;
  strategyType: string;
  parameters: Record<string, any>;
  lastBacktestResult?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

interface StrategyTypeDefinition {
  type: string;
  name: string;
  parameters: any[];
}

const StrategyLibraryPage: React.FC = () => {
  const [strategies, setStrategies] = useState<StrategyConfig[]>([]);
  const [strategyTypes, setStrategyTypes] = useState<StrategyTypeDefinition[]>([]);
  const [loading, setLoading] = useState(false);
  const [showEditor, setShowEditor] = useState(false);
  const [editingStrategy, setEditingStrategy] = useState<StrategyConfig | null>(null);
  const [error, setError] = useState<string | null>(null);

  const { GetAllStrategies, GetStrategyTypes, DeleteStrategy } = useWailsAPI();

  const loadStrategies = async () => {
    setLoading(true);
    setError(null);
    try {
      console.log('开始加载策略列表...');
      const result = await GetAllStrategies();
      console.log('获取到的策略列表结果:', result);
      if (Array.isArray(result)) {
        setStrategies(result);
        console.log('策略列表已加载，共', result.length, '个策略');
      } else {
        console.warn('返回的结果不是数组:', result);
        setError('返回的数据格式错误');
      }
    } catch (err: any) {
      console.error('加载策略列表失败:', err);
      console.error('错误对象:', err);
      console.error('错误消息:', err.message);
      setError(err.message || '加载策略列表失败');
    } finally {
      setLoading(false);
    }
  };

  const loadStrategyTypes = async () => {
    try {
      const result = await GetStrategyTypes();
      if (Array.isArray(result)) {
        setStrategyTypes(result);
      }
    } catch (err: any) {
      console.error('加载策略类型失败:', err);
    }
  };

  useEffect(() => {
    loadStrategies();
    loadStrategyTypes();
  }, []);

  const handleAdd = () => {
    setEditingStrategy(null);
    setShowEditor(true);
  };

  const handleEdit = (strategy: StrategyConfig) => {
    setEditingStrategy(strategy);
    setShowEditor(true);
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('确定要删除这个策略吗？')) {
      return;
    }
    try {
      await DeleteStrategy(id);
      await loadStrategies();
    } catch (err: any) {
      setError(err.message || '删除策略失败');
    }
  };

  const handleEditorClose = (saved: boolean) => {
    setShowEditor(false);
    setEditingStrategy(null);
    if (saved) {
      loadStrategies();
    }
  };

  const getStrategyTypeName = (type: string) => {
    const st = strategyTypes.find(t => t.type === type);
    return st ? st.name : type;
  };

  const formatPercentage = (value: number) => `${value.toFixed(2)}%`;

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* 页面标题 */}
        <div className="mb-8 flex justify-between items-center">
          <div>
            <h1 className="text-4xl font-bold mb-2">策略库</h1>
            <p className="text-gray-400">管理和保存您的交易策略配置</p>
          </div>
          <button
            onClick={handleAdd}
            className="flex items-center gap-2 px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md transition-colors"
          >
            <Plus className="w-5 h-5" />
            新建策略
          </button>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-900/50 border border-red-700 rounded-md text-red-300">
            {error}
          </div>
        )}

        {loading ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-400">加载中...</p>
          </div>
        ) : strategies.length === 0 ? (
          <div className="text-center py-12 bg-gray-800 rounded-lg">
            <Activity className="w-16 h-16 mx-auto text-gray-600 mb-4" />
            <p className="text-gray-400 text-lg">还没有策略</p>
            <p className="text-gray-500 text-sm mt-2">点击"新建策略"开始创建您的第一个策略</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {strategies.map((strategy) => (
              <div
                key={strategy.id}
                className="bg-gray-800 rounded-lg p-6 shadow-lg hover:shadow-xl transition-shadow"
              >
                {/* 策略头部 */}
                <div className="flex justify-between items-start mb-4">
                  <div className="flex-1">
                    <h3 className="text-xl font-bold mb-1">{strategy.name}</h3>
                    <p className="text-sm text-blue-400">{getStrategyTypeName(strategy.strategyType)}</p>
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleEdit(strategy)}
                      className="p-2 text-gray-400 hover:text-blue-400 hover:bg-gray-700 rounded-md transition-colors"
                      title="编辑"
                    >
                      <Edit className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDelete(strategy.id)}
                      className="p-2 text-gray-400 hover:text-red-400 hover:bg-gray-700 rounded-md transition-colors"
                      title="删除"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>

                {/* 策略描述 */}
                {strategy.description && (
                  <p className="text-sm text-gray-400 mb-4 line-clamp-2">{strategy.description}</p>
                )}

                {/* 策略参数摘要 */}
                <div className="mb-4 p-3 bg-gray-700 rounded-md">
                  <p className="text-xs text-gray-400 mb-2">参数配置:</p>
                  <div className="text-xs space-y-1">
                    {Object.entries(strategy.parameters).slice(0, 3).map(([key, value]) => (
                      <div key={key} className="flex justify-between">
                        <span className="text-gray-500">{key}:</span>
                        <span className="text-gray-300">{String(value)}</span>
                      </div>
                    ))}
                    {Object.keys(strategy.parameters).length > 3 && (
                      <p className="text-gray-500">等 {Object.keys(strategy.parameters).length} 项...</p>
                    )}
                  </div>
                </div>

                {/* 最后回测结果 */}
                {strategy.lastBacktestResult && (
                  <div className="mb-4 p-3 bg-gray-700 rounded-md">
                    <p className="text-xs text-gray-400 mb-2">最后回测结果:</p>
                    <div className="flex justify-between items-center">
                      <div className="flex items-center gap-2">
                        {(strategy.lastBacktestResult.totalReturn || 0) >= 0 ? (
                          <TrendingUp className="w-4 h-4 text-green-400" />
                        ) : (
                          <TrendingDown className="w-4 h-4 text-red-400" />
                        )}
                        <span className="text-sm font-semibold">
                          {formatPercentage(strategy.lastBacktestResult.totalReturn || 0)}
                        </span>
                      </div>
                      <span className="text-xs text-gray-500">
                        胜率: {formatPercentage(strategy.lastBacktestResult.winRate || 0)}
                      </span>
                    </div>
                  </div>
                )}

                {/* 底部信息 */}
                <div className="flex justify-between items-center text-xs text-gray-500">
                  <span>更新: {new Date(strategy.updatedAt).toLocaleDateString()}</span>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* 策略编辑器 Modal */}
        {showEditor && (
          <StrategyEditor
            isOpen={showEditor}
            strategy={editingStrategy}
            strategyTypes={strategyTypes}
            onClose={handleEditorClose}
          />
        )}
      </div>
    </div>
  );
};

export default StrategyLibraryPage;
