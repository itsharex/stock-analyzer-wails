import { useState, useEffect } from 'react'
import { useWailsAPI } from '../hooks/useWailsAPI'
import type { AppConfig } from '../types'

interface SettingsProps {
  onConfigSaved?: () => void
}

function Settings({ onConfigSaved }: SettingsProps) {
  const [config, setConfig] = useState<AppConfig>({
    apiKey: '',
    baseUrl: '',
    model: '',
    models: []
  })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState({ type: '', text: '' })
  
  const { getConfig, saveConfig } = useWailsAPI()

  useEffect(() => {
    loadConfig()
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
            阿里百炼 API Key
          </label>
          <input
            type="password"
            value={config.apiKey}
            onChange={(e) => setConfig({ ...config, apiKey: e.target.value })}
            placeholder="sk-..."
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none transition"
          />
          <p className="mt-1 text-xs text-gray-500">
            从阿里云百炼控制台获取的 API 密钥
          </p>
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
            {config.models?.map((m) => (
              <option key={m} value={m}>
                {m}
              </option>
            ))}
            {!config.models?.includes(config.model) && config.model && (
              <option value={config.model}>{config.model} (自定义)</option>
            )}
          </select>
          <p className="mt-1 text-xs text-gray-500">
            选择要使用的通义千问模型版本
          </p>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Base URL
          </label>
          <input
            type="text"
            value={config.baseUrl}
            onChange={(e) => setConfig({ ...config, baseUrl: e.target.value })}
            placeholder="https://dashscope.aliyuncs.com/compatible-mode/v1"
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

      <div className="mt-8 p-4 bg-blue-50 rounded-lg border border-blue-100">
        <h3 className="text-sm font-semibold text-blue-800 mb-2">💡 提示</h3>
        <ul className="text-xs text-blue-700 space-y-1 list-disc pl-4">
          <li>配置将保存在本地 `config.yaml` 文件中。</li>
          <li>保存后 AI 服务将立即使用新配置重新初始化。</li>
          <li>推荐使用 `qwen-plus` 以获得最佳的分析性价比。</li>
        </ul>
      </div>
    </div>
  )
}

export default Settings
