import React, { useState, useEffect } from 'react';
import { X, Save } from 'lucide-react';
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
  parameters: {
    name: string;
    label: string;
    type: string;
    minValue?: number;
    maxValue?: number;
    defaultValue: any;
    options?: string[];
  }[];
}

interface StrategyEditorProps {
  isOpen: boolean;
  strategy: StrategyConfig | null;
  strategyTypes: StrategyTypeDefinition[];
  onClose: (saved: boolean) => void;
}

const StrategyEditor: React.FC<StrategyEditorProps> = ({
  isOpen,
  strategy,
  strategyTypes,
  onClose,
}) => {
  const { CreateStrategy, UpdateStrategy } = useWailsAPI();

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [strategyType, setStrategyType] = useState('simple_ma');
  const [parameters, setParameters] = useState<Record<string, any>>({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 初始化表单
  useEffect(() => {
    if (strategy) {
      setName(strategy.name);
      setDescription(strategy.description);
      setStrategyType(strategy.strategyType);
      setParameters({ ...strategy.parameters });
    } else {
      // 重置表单为默认值
      setName('');
      setDescription('');
      setStrategyType('simple_ma');
      // 使用第一个策略类型的默认参数
      if (strategyTypes.length > 0) {
        const defaults: Record<string, any> = {};
        strategyTypes[0].parameters.forEach(param => {
          defaults[param.name] = param.defaultValue;
        });
        setParameters(defaults);
      }
    }
    setError(null);
  }, [strategy, strategyTypes, isOpen]);

  // 当策略类型改变时，更新参数
  useEffect(() => {
    const selectedType = strategyTypes.find(t => t.type === strategyType);
    if (selectedType && !strategy) {
      const defaults: Record<string, any> = {};
      selectedType.parameters.forEach(param => {
        defaults[param.name] = param.defaultValue;
      });
      setParameters(defaults);
    }
  }, [strategyType, strategyTypes, strategy]);

  const handleParameterChange = (paramName: string, value: any) => {
    setParameters(prev => ({
      ...prev,
      [paramName]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    // 验证必填字段
    if (!name.trim()) {
      setError('策略名称不能为空');
      setLoading(false);
      return;
    }

    // 验证参数
    const selectedType = strategyTypes.find(t => t.type === strategyType);
    if (selectedType) {
      for (const param of selectedType.parameters) {
        const value = parameters[param.name];
        if (param.type === 'number') {
          const numValue = Number(value);
          if (param.minValue !== undefined && numValue < param.minValue) {
            setError(`${param.label} 不能小于 ${param.minValue}`);
            setLoading(false);
            return;
          }
          if (param.maxValue !== undefined && numValue > param.maxValue) {
            setError(`${param.label} 不能大于 ${param.maxValue}`);
            setLoading(false);
            return;
          }
        }
      }
    }

    try {
      if (strategy) {
        // 更新策略
        await UpdateStrategy(strategy.id, name, description, strategyType, parameters);
      } else {
        // 创建新策略
        await CreateStrategy(name, description, strategyType, parameters);
      }
      onClose(true);
    } catch (err: any) {
      setError(err.message || '保存策略失败');
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  const selectedType = strategyTypes.find(t => t.type === strategyType);

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-gray-800 rounded-lg shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-hidden">
        {/* Modal Header */}
        <div className="flex justify-between items-center p-6 border-b border-gray-700">
          <h2 className="text-2xl font-bold">
            {strategy ? '编辑策略' : '新建策略'}
          </h2>
          <button
            onClick={() => onClose(false)}
            className="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-md transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Modal Body */}
        <form onSubmit={handleSubmit} className="p-6 overflow-y-auto max-h-[calc(90vh-200px)]">
          {error && (
            <div className="mb-4 p-4 bg-red-900/50 border border-red-700 rounded-md text-red-300">
              {error}
            </div>
          )}

          {/* 基本信息 */}
          <div className="space-y-4 mb-6">
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                策略名称 <span className="text-red-400">*</span>
              </label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                placeholder="输入策略名称"
                disabled={loading}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">描述</label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={3}
                className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                placeholder="输入策略描述（可选）"
                disabled={loading}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                策略类型
              </label>
              <select
                value={strategyType}
                onChange={(e) => setStrategyType(e.target.value)}
                className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                disabled={loading}
              >
                {strategyTypes.map((st) => (
                  <option key={st.type} value={st.type}>
                    {st.name}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* 参数配置 */}
          {selectedType && (
            <div className="space-y-4">
              <h3 className="text-lg font-semibold text-gray-300">参数配置</h3>
              {selectedType.parameters.map((param) => (
                <div key={param.name}>
                  <label className="block text-sm font-medium text-gray-300 mb-2">
                    {param.label}
                  </label>
                  {param.type === 'number' ? (
                    <input
                      type="number"
                      value={parameters[param.name] ?? ''}
                      onChange={(e) => handleParameterChange(param.name, Number(e.target.value))}
                      min={param.minValue}
                      max={param.maxValue}
                      className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                      disabled={loading}
                    />
                  ) : param.type === 'select' ? (
                    <select
                      value={parameters[param.name] ?? ''}
                      onChange={(e) => handleParameterChange(param.name, e.target.value)}
                      className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                      disabled={loading}
                    >
                      {param.options?.map((option) => (
                        <option key={option} value={option}>
                          {option}
                        </option>
                      ))}
                    </select>
                  ) : (
                    <input
                      type="text"
                      value={parameters[param.name] ?? ''}
                      onChange={(e) => handleParameterChange(param.name, e.target.value)}
                      className="w-full px-4 py-2 rounded-md bg-gray-700 border border-gray-600 text-gray-100 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                      disabled={loading}
                    />
                  )}
                </div>
              ))}
            </div>
          )}
        </form>

        {/* Modal Footer */}
        <div className="flex justify-end gap-3 p-6 border-t border-gray-700">
          <button
            type="button"
            onClick={() => onClose(false)}
            className="px-6 py-2 bg-gray-700 hover:bg-gray-600 text-white font-medium rounded-md transition-colors"
            disabled={loading}
          >
            取消
          </button>
          <button
            type="submit"
            onClick={handleSubmit}
            className="flex items-center gap-2 px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md transition-colors disabled:opacity-50"
            disabled={loading}
          >
            <Save className="w-4 h-4" />
            {loading ? '保存中...' : '保存'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default StrategyEditor;
