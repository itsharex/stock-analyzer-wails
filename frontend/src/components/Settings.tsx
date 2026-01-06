import { useState, useEffect } from 'react'
import { useWailsAPI } from '../hooks/useWailsAPI'
import type { AppConfig } from '../types'

interface SettingsProps {
  onConfigSaved?: () => void
}

interface StrategyConfig {
  trailingStopActivation: number
  trailingStopCallback: number
}

function Settings({ onConfigSaved }: SettingsProps) {
  const [config, setConfig] = useState<AppConfig>({
    provider: 'Qwen',
    apiKey: '',
    baseUrl: '',
    model: '',
    providerModels: {}
  })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState({ type: '', text: '' })
  const [strategyConfig, setStrategyConfig] = useState<StrategyConfig>({
    trailingStopActivation: 0.05,
    trailingStopCallback: 0.03
  })
  // const [strategyLoading, setStrategyLoading] = useState(true)
  
  const { getConfig, saveConfig } = useWailsAPI()

  useEffect(() => {
    loadConfig()
    loadStrategyConfig()
  }, [])

  const loadConfig = async () => {
    try {
      const data = await getConfig()
      setConfig(data)
    } catch (err) {
      setMessage({ type: 'error', text: '加载配置失败' })
    } finally {
      setLoading(false)
    }
  }

  const loadStrategyConfig = async () => {
    try {
      const result = await (window as any).runtime?.Call?.('app.GetGlobalStrategyConfig')
      if (result) {
        setStrategyConfig(result)
      }
    } catch (err) {
      console.error('加载策略配置失败:', err)
    } finally {
      // setStrategyLoading(false)
    }
  }

  const handleSaveStrategy = async () => {
    setSaving(true)
    setMessage({ type: '', text: '' })
    try {
      await (window as any).runtime?.Call?.('app.UpdateGlobalStrategyConfig', strategyConfig)
      setMessage({ type: 'success', text: '交易策略配置已保存' })
    } catch (err: any) {
      setMessage({ text: `保存失败: ${err.message || err}`, type: 'error' })
    } finally {
      setSaving(false)
    }
  }

  const handleProviderChange = (newProvider: string) => {
    const models = config.providerModels[newProvider] || []
    setConfig({
      ...config,
      provider: newProvider,
      model: models[0] || '',
      // 自动填充一些常见的 BaseURL
      baseUrl: getBaseURLForProvider(newProvider)
    })
  }

  const getBaseURLForProvider = (provider: string) => {
    switch (provider) {
      case 'Qwen':
      case 'DashScope': return 'https://dashscope.aliyuncs.com/compatible-mode/v1'
      case 'DeepSeek': return 'https://api.deepseek.com'
      case 'OpenAI': return 'https://api.openai.com/v1'
      case 'Claude': return 'https://api.anthropic.com/v1'
      case 'ARK': return 'https://ark.cn-beijing.volces.com/api/v3'
      default: return config.baseUrl
    }
  }

  const handleSave = async () => {
    setSaving(true)
    setMessage({ type: '', text: '' })
    try {
      await saveConfig(config)
      setMessage({ type: 'success', text: '配置已保存并生效' })
      if (onConfigSaved) {
        onConfigSaved()
      }
    } catch (err: any) {
      setMessage({ text: `保存失败: ${err.message || err}`, type: 'error' })
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto bg-white rounded-xl shadow-lg p-8">
      <h2 className="text-2xl font-bold text-gray-800 mb-6 flex items-center">
        <span className="mr-2">⚙️</span> 系统设置
      </h2>

      {message.text && (
        <div className={`mb-6 p-4 rounded-lg flex items-center ${
          message.type === 'success' ? 'bg-green-50 text-green-700 border border-green-200' : 'bg-red-50 text-red-700 border border-red-200'
        }`}>
          <span className="mr-2">{message.type === 'success' ? '✅' : '❌'}</span>
          {message.text}
        </div>
      )}

      <div className="space-y-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            AI 供应商 (Provider)
          </label>
          <select
            value={config.provider}
            onChange={(e) => handleProviderChange(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none transition bg-white"
          >
            {Object.keys(config.providerModels).map((p) => (
              <option key={p} value={p}>{p}</option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            API Key
          </label>
          <input
            type="password"
            value={config.apiKey}
            onChange={(e) => setConfig({ ...config, apiKey: e.target.value })}
            placeholder="请输入 API Key"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none transition"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            模型选择 (Model)
          </label>
          <select
            value={config.model}
            onChange={(e) => setConfig({ ...config, model: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none transition bg-white"
          >
            {(config.providerModels[config.provider] || []).map((m) => (
              <option key={m} value={m}>{m}</option>
            ))}
            {!config.providerModels[config.provider]?.includes(config.model) && config.model && (
              <option value={config.model}>{config.model} (自定义)</option>
            )}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Base URL
          </label>
          <input
            type="text"
            value={config.baseUrl}
            onChange={(e) => setConfig({ ...config, baseUrl: e.target.value })}
            placeholder="https://api.example.com/v1"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none transition"
          />
        </div>

        <div className="pt-4">
          <button
            onClick={handleSave}
            disabled={saving}
            className="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-lg transition shadow-md disabled:opacity-50"
          >
            {saving ? '正在保存...' : '保存配置'}
          </button>
        </div>
      </div>

      <div className="mt-8 border-t pt-8">
        <h2 className="text-2xl font-bold text-gray-800 mb-6 flex items-center">
          <span className="mr-2">📈</span> 交易策略默认配置
        </h2>

        <div className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              移动止损启动阈值 (%)
              <span className="text-xs text-gray-500 ml-2">当盈利达到此比例时启动移动止损</span>
            </label>
            <div className="flex items-center gap-4">
              <input
                type="range"
                min="0"
                max="0.20"
                step="0.01"
                value={strategyConfig.trailingStopActivation}
                onChange={(e) => setStrategyConfig({ ...strategyConfig, trailingStopActivation: parseFloat(e.target.value) })}
                className="flex-1 h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
              />
              <input
                type="number"
                min="0"
                max="0.20"
                step="0.01"
                value={strategyConfig.trailingStopActivation}
                onChange={(e) => setStrategyConfig({ ...strategyConfig, trailingStopActivation: parseFloat(e.target.value) })}
                className="w-20 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
              />
              <span className="text-sm text-gray-600 w-12 text-right">{(strategyConfig.trailingStopActivation * 100).toFixed(1)}%</span>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              移动止损回撤比例 (%)
              <span className="text-xs text-gray-500 ml-2">价格回撤此比例时触发止盈</span>
            </label>
            <div className="flex items-center gap-4">
              <input
                type="range"
                min="0"
                max="0.20"
                step="0.01"
                value={strategyConfig.trailingStopCallback}
                onChange={(e) => setStrategyConfig({ ...strategyConfig, trailingStopCallback: parseFloat(e.target.value) })}
                className="flex-1 h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
              />
              <input
                type="number"
                min="0"
                max="0.20"
                step="0.01"
                value={strategyConfig.trailingStopCallback}
                onChange={(e) => setStrategyConfig({ ...strategyConfig, trailingStopCallback: parseFloat(e.target.value) })}
                className="w-20 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
              />
              <span className="text-sm text-gray-600 w-12 text-right">{(strategyConfig.trailingStopCallback * 100).toFixed(1)}%</span>
            </div>
          </div>

          <div className="pt-4">
            <button
              onClick={handleSaveStrategy}
              disabled={saving}
              className="w-full bg-green-600 hover:bg-green-700 text-white font-bold py-3 px-6 rounded-lg transition shadow-md disabled:opacity-50"
            >
              {saving ? '正在保存...' : '保存交易策略配置'}
            </button>
          </div>
        </div>

        <div className="mt-6 p-4 bg-green-50 rounded-lg border border-green-100">
          <h3 className="text-sm font-semibold text-green-800 mb-2">💡 提示</h3>
          <ul className="text-xs text-green-700 space-y-1 list-disc pl-4">
            <li>这些参数将作为建仓时的默认值，用户仍可在建仓时手动调整。</li>
            <li>启动阈值：建议 3% - 10%，表示盈利多少后启动移动止损。</li>
            <li>回撤比例：建议 2% - 5%，表示从最高点回撤多少后止盈。</li>
          </ul>
        </div>
      </div>

      <div className="mt-8 p-4 bg-blue-50 rounded-lg border border-blue-100">
        <h3 className="text-sm font-semibold text-blue-800 mb-2">💡 提示</h3>
        <ul className="text-xs text-blue-700 space-y-1 list-disc pl-4">
          <li>支持 OpenAI 兼容协议的所有供应商。</li>
          <li>切换供应商后，Base URL 会尝试自动填充默认值。</li>
          <li>配置将保存在本地 `config.yaml` 文件中。</li>
        </ul>
      </div>
    </div>
  )
}

export default Settings
