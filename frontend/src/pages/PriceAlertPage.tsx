import React, { useState, useEffect } from 'react';
import { parseError } from '../utils/errorHandler';
import { useWailsAPI } from '../hooks/useWailsAPI';
import {
  Plus, Edit3, Trash2, Bell, History, Search,
  CheckCircle, XCircle, Clock, Target, Shield,
  Play, Pause, PlusCircle, MinusCircle
} from 'lucide-react';

interface PriceAlert {
  id: number;
  stockCode: string;
  stockName: string;
  alertType: string;
  conditions: string;
  isActive: boolean;
  sensitivity: number;
  cooldownHours: number;
  postTriggerAction: string;
  enableSound: boolean;
  enableDesktop: boolean;
  templateId: string;
  createdAt: string;
  updatedAt: string;
  lastTriggeredAt: string;
}

interface AlertTemplate {
  id: string;
  name: string;
  description: string;
  alertType: string;
  conditions: string;
  createdAt: string;
}

interface TriggerHistory {
  id: number;
  alertId: number;
  stockCode: string;
  stockName: string;
  alertType: string;
  triggerPrice: number;
  triggerMessage: string;
  triggeredAt: string;
}

const PriceAlertPage: React.FC = () => {
  const {
    getAllPriceAlerts,
    getPriceAlertTemplates,
    getPriceAlertHistory,
    createPriceAlert,
    updatePriceAlert,
    deletePriceAlert,
    togglePriceAlert,
    createPriceAlertFromTemplate,
    getStockData,
  } = useWailsAPI();

  const [alerts, setAlerts] = useState<PriceAlert[]>([]);
  const [templates, setTemplates] = useState<AlertTemplate[]>([]);
  const [histories, setHistories] = useState<TriggerHistory[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // 表单状态
  const [showModal, setShowModal] = useState<boolean>(false);
  const [editingAlert, setEditingAlert] = useState<PriceAlert | null>(null);
  const [selectedTemplate, setSelectedTemplate] = useState<string>('');
  const [formData, setFormData] = useState({
    stockCode: '',
    stockName: '',
    alertType: '',
    conditions: '',
    sensitivity: 0.001,
    cooldownHours: 1,
    postTriggerAction: 'continue',
    enableSound: true,
    enableDesktop: true,
  });

  // 结构化预警条件状态
  const [alertConditions, setAlertConditions] = useState<any>({
    logic: 'AND',
    conditions: [
      { field: '', operator: '', value: 0, reference: '' }
    ]
  });

  // 股票代码查询状态
  const [searchingStock, setSearchingStock] = useState(false);
  const [stockCodeError, setStockCodeError] = useState<string | null>(null);

  // 筛选状态
  const [filterCode, setFilterCode] = useState<string>('');
  const [filterActive, setFilterActive] = useState<boolean | null>(null);
  const [activeTab, setActiveTab] = useState<'alerts' | 'templates' | 'history'>('alerts');

  useEffect(() => {
    fetchData();

    // 监听价格预警触发事件
    const unsubscribe = (window as any).runtime?.EventsOn('price_alert_triggered', (data: any) => {
      console.log('Price alert triggered:', data);

      // 如果当前在历史标签页，刷新历史数据
      if (activeTab === 'history') {
        loadTriggerHistory();
      }

      // 显示浏览器通知（如果启用了桌面通知）
      if (data.enableDesktop && 'Notification' in window) {
        if (Notification.permission === 'granted') {
          new Notification(`${data.stockName} (${data.stockCode}) - 预警触发`, {
            body: `${data.message}，当前价格: ¥${data.triggerPrice.toFixed(2)}`,
            icon: '/icon.png',
            requireInteraction: true,
          });
        } else if (Notification.permission !== 'denied') {
          Notification.requestPermission().then(permission => {
            if (permission === 'granted') {
              new Notification(`${data.stockName} (${data.stockCode}) - 预警触发`, {
                body: `${data.message}，当前价格: ¥${data.triggerPrice.toFixed(2)}`,
                icon: '/icon.png',
                requireInteraction: true,
              });
            }
          });
        }
      }

      // 播放提示音（如果启用了声音）
      if (data.enableSound) {
        try {
          const audio = new Audio('/alert.mp3');
          audio.play().catch(err => console.error('Failed to play alert sound:', err));
        } catch (err) {
          console.error('Failed to play alert sound:', err);
        }
      }
    });

    return () => {
      if (unsubscribe) unsubscribe();
    };
  }, [activeTab]);

  const fetchData = async () => {
    setLoading(true);
    setError(null);
    try {
      const [alertsRes, templatesRes] = await Promise.all([
        getAllPriceAlerts(),
        getPriceAlertTemplates(),
      ]);

      if (alertsRes?.success) {
        setAlerts(alertsRes.alerts || []);
      }
      if (templatesRes?.success) {
        setTemplates(templatesRes.templates || []);
      }
    } catch (err) {
      const errorResult = parseError(err);
      setError(errorResult.message);
    } finally {
      setLoading(false);
    }
  };

  const loadTriggerHistory = async (stockCode?: string) => {
    try {
      const res = await getPriceAlertHistory(stockCode || '', 100);
      if (res?.success) {
        setHistories(res.histories || []);
      }
    } catch (err) {
      console.error('Failed to load trigger history:', err);
    }
  };

  // 股票代码自动查询
  const handleStockCodeBlur = async (code: string) => {
    if (!code || code.length !== 6) {
      return;
    }

    setSearchingStock(true);
    setStockCodeError(null);

    try {
      const stockData = await getStockData(code);
      if (stockData) {
        setFormData(prev => ({ ...prev, stockName: stockData.name }));
      } else {
        setStockCodeError('未找到该股票');
      }
    } catch (err) {
      console.error('Failed to fetch stock data:', err);
      setStockCodeError('查询股票信息失败');
    } finally {
      setSearchingStock(false);
    }
  };

  // 预警类型变化时重置条件
  const handleAlertTypeChange = (alertType: string) => {
    setFormData(prev => ({ ...prev, alertType }));

    // 根据预警类型初始化默认条件
    const defaultConditions = getDefaultConditions(alertType);
    setAlertConditions(defaultConditions);
  };

  // 获取默认预警条件
  const getDefaultConditions = (alertType: string) => {
    switch (alertType) {
      case 'price_change':
        return {
          logic: 'AND',
          conditions: [
            { field: 'price_change_percent', operator: '>', value: 5 }
          ]
        };
      case 'target_price':
        return {
          logic: 'AND',
          conditions: [
            { field: 'close_price', operator: '>=', value: 0 }
          ]
        };
      case 'stop_loss':
        return {
          logic: 'AND',
          conditions: [
            { field: 'close_price', operator: '<=', value: 0 }
          ]
        };
      case 'high_low':
        return {
          logic: 'AND',
          conditions: [
            { field: 'high_price', operator: '>', value: 0, reference: 'historical_high' }
          ]
        };
      case 'price_range':
        return {
          logic: 'AND',
          conditions: [
            { field: 'close_price', operator: '>=', value: 0 },
            { field: 'close_price', operator: '<=', value: 0 }
          ]
        };
      case 'ma_deviation':
        return {
          logic: 'AND',
          conditions: [
            { field: 'ma5', operator: '>', value: 0, reference: 'ma20' }
          ]
        };
      case 'combined':
        return {
          logic: 'AND',
          conditions: [
            { field: 'price_change_percent', operator: '>', value: 0 },
            { field: 'volume_ratio', operator: '>', value: 0 }
          ]
        };
      default:
        return {
          logic: 'AND',
          conditions: [
            { field: '', operator: '', value: 0 }
          ]
        };
    }
  };

  // 添加条件（组合预警）
  const addCondition = () => {
    setAlertConditions((prev: any) => ({
      ...prev,
      conditions: [...prev.conditions, { field: '', operator: '', value: 0 }]
    }));
  };

  // 删除条件
  const removeCondition = (index: number) => {
    if (alertConditions.conditions.length <= 1) {
      return;
    }
    setAlertConditions((prev: any) => ({
      ...prev,
      conditions: prev.conditions.filter((_: any, i: number) => i !== index)
    }));
  };

  // 更新条件
  const updateCondition = (index: number, field: string, value: any) => {
    setAlertConditions((prev: any) => ({
      ...prev,
      conditions: prev.conditions.map((cond: any, i: number) =>
        i === index ? { ...cond, [field]: value } : cond
      )
    }));
  };

  // 将结构化条件转换为JSON字符串
  const serializeConditions = (): string => {
    return JSON.stringify(alertConditions);
  };

  // 渲染预警条件表单
  const renderAlertConditions = () => {
    switch (formData.alertType) {
      case 'price_change':
        return <PriceChangeForm alertConditions={alertConditions} updateCondition={updateCondition} />;
      case 'target_price':
        return <TargetPriceForm alertConditions={alertConditions} updateCondition={updateCondition} />;
      case 'stop_loss':
        return <StopLossForm alertConditions={alertConditions} updateCondition={updateCondition} />;
      case 'high_low':
        return <HighLowForm alertConditions={alertConditions} updateCondition={updateCondition} />;
      case 'price_range':
        return <PriceRangeForm alertConditions={alertConditions} updateCondition={updateCondition} />;
      case 'ma_deviation':
        return <MADeviationForm alertConditions={alertConditions} updateCondition={updateCondition} />;
      case 'combined':
        return (
          <CombinedForm
            alertConditions={alertConditions}
            updateCondition={updateCondition}
            addCondition={addCondition}
            removeCondition={removeCondition}
            setLogic={(logic: string) => setAlertConditions((prev: any) => ({ ...prev, logic }))}
          />
        );
      default:
        return <div className="text-sm text-gray-500">请选择预警类型</div>;
    }
  };

  // 涨跌幅预警表单
  const PriceChangeForm = ({ alertConditions, updateCondition }: any) => {
    const condition = alertConditions.conditions[0] || {};
    return (
      <div className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">涨跌幅类型</label>
            <select
              value={condition.operator}
              onChange={(e) => updateCondition(0, 'operator', e.target.value)}
              className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            >
              <option value=">">涨幅超过</option>
              <option value="<">跌幅超过</option>
            </select>
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">百分比 (%)</label>
            <input
              type="number"
              step="0.1"
              min="0"
              max="100"
              value={condition.value || ''}
              onChange={(e) => updateCondition(0, 'value', parseFloat(e.target.value))}
              className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
              placeholder="例如: 5"
            />
          </div>
        </div>
        <p className="text-xs text-gray-500">
          {condition.operator === '>' ? '当涨幅达到设定值时触发预警' : '当跌幅达到设定值时触发预警'}
        </p>
      </div>
    );
  };

  // 目标价预警表单
  const TargetPriceForm = ({ alertConditions, updateCondition }: any) => {
    const condition = alertConditions.conditions[0] || {};
    return (
      <div className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">操作符</label>
            <select
              value={condition.operator}
              onChange={(e) => updateCondition(0, 'operator', e.target.value)}
              className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            >
              <option value=">=">价格达到或高于</option>
              <option value="<=">价格达到或低于</option>
            </select>
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">目标价格 (元)</label>
            <input
              type="number"
              step="0.01"
              min="0"
              value={condition.value || ''}
              onChange={(e) => updateCondition(0, 'value', parseFloat(e.target.value))}
              className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
              placeholder="例如: 100.00"
            />
          </div>
        </div>
        <p className="text-xs text-gray-500">
          当股价达到设定的目标价格时触发预警
        </p>
      </div>
    );
  };

  // 止损价预警表单
  const StopLossForm = ({ alertConditions, updateCondition }: any) => {
    const condition = alertConditions.conditions[0] || {};
    return (
      <div className="space-y-4">
        <div>
          <label className="block text-xs font-medium text-gray-600 mb-1">止损价格 (元)</label>
          <input
            type="number"
            step="0.01"
            min="0"
            value={condition.value || ''}
            onChange={(e) => updateCondition(0, 'value', parseFloat(e.target.value))}
            className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            placeholder="例如: 90.00"
          />
        </div>
        <p className="text-xs text-gray-500">
          当股价跌破设定的止损价格时触发预警，用于风险控制
        </p>
      </div>
    );
  };

  // 突破高低点预警表单
  const HighLowForm = ({ alertConditions, updateCondition }: any) => {
    const condition = alertConditions.conditions[0] || {};
    return (
      <div className="space-y-4">
        <div>
          <label className="block text-xs font-medium text-gray-600 mb-1">突破类型</label>
          <select
            value={condition.reference}
            onChange={(e) => updateCondition(0, 'reference', e.target.value)}
            className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
          >
            <option value="historical_high">突破历史新高</option>
            <option value="historical_low">跌破历史新低</option>
          </select>
        </div>
        <p className="text-xs text-gray-500">
          当股价突破近期的历史最高价或最低价时触发预警
        </p>
      </div>
    );
  };

  // 价格区间预警表单
  const PriceRangeForm = ({ alertConditions, updateCondition }: any) => {
    const condition1 = alertConditions.conditions[0] || {};
    const condition2 = alertConditions.conditions[1] || {};
    return (
      <div className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">价格下限 (元)</label>
            <input
              type="number"
              step="0.01"
              min="0"
              value={condition1.value || ''}
              onChange={(e) => updateCondition(0, 'value', parseFloat(e.target.value))}
              className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
              placeholder="例如: 90.00"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">价格上限 (元)</label>
            <input
              type="number"
              step="0.01"
              min="0"
              value={condition2.value || ''}
              onChange={(e) => updateCondition(1, 'value', parseFloat(e.target.value))}
              className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
              placeholder="例如: 110.00"
            />
          </div>
        </div>
        <p className="text-xs text-gray-500">
          当股价进入或超出设定的价格区间时触发预警
        </p>
      </div>
    );
  };

  // 均线偏离预警表单
  const MADeviationForm = ({ alertConditions, updateCondition }: any) => {
    const condition = alertConditions.conditions[0] || {};
    return (
      <div className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">均线类型</label>
            <select
              value={condition.reference}
              onChange={(e) => updateCondition(0, 'reference', e.target.value)}
              className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            >
              <option value="ma20">MA5 上穿 MA20 (金叉)</option>
              <option value="ma20">MA5 下穿 MA20 (死叉)</option>
              <option value="ma10">MA5 上穿 MA10</option>
              <option value="ma10">MA5 下穿 MA10</option>
            </select>
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">偏离幅度 (%)</label>
            <input
              type="number"
              step="0.1"
              min="0"
              max="10"
              value={condition.value || ''}
              onChange={(e) => updateCondition(0, 'value', parseFloat(e.target.value))}
              className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
              placeholder="例如: 0.5"
            />
          </div>
        </div>
        <p className="text-xs text-gray-500">
          当股价均线发生交叉或偏离设定幅度时触发预警
        </p>
      </div>
    );
  };

  // 组合预警表单
  const CombinedForm = ({ alertConditions, updateCondition, addCondition, removeCondition, setLogic }: any) => {
    const fields = [
      { value: 'price_change_percent', label: '涨跌幅' },
      { value: 'close_price', label: '收盘价' },
      { value: 'high_price', label: '最高价' },
      { value: 'low_price', label: '最低价' },
      { value: 'volume_ratio', label: '量比' },
      { value: 'ma5', label: 'MA5' },
      { value: 'ma10', label: 'MA10' },
      { value: 'ma20', label: 'MA20' },
    ];

    return (
      <div className="space-y-4">
        {/* 逻辑关系选择 */}
        <div>
          <label className="block text-xs font-medium text-gray-600 mb-1">逻辑关系</label>
          <select
            value={alertConditions.logic}
            onChange={(e) => setLogic(e.target.value)}
            className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
          >
            <option value="AND">满足所有条件 (AND)</option>
            <option value="OR">满足任一条件 (OR)</option>
          </select>
        </div>

        {/* 条件列表 */}
        <div className="space-y-3">
          {alertConditions.conditions.map((condition: any, index: number) => (
            <div key={index} className="flex items-start space-x-2 bg-white p-3 rounded-lg border border-gray-200">
              <div className="flex-1 grid grid-cols-3 gap-2">
                <select
                  value={condition.field}
                  onChange={(e) => updateCondition(index, 'field', e.target.value)}
                  className="px-2 py-1.5 bg-gray-50 border border-gray-200 rounded text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                >
                  <option value="">选择字段</option>
                  {fields.map(field => (
                    <option key={field.value} value={field.value}>{field.label}</option>
                  ))}
                </select>

                <select
                  value={condition.operator}
                  onChange={(e) => updateCondition(index, 'operator', e.target.value)}
                  className="px-2 py-1.5 bg-gray-50 border border-gray-200 rounded text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                >
                  <option value=">">&gt;</option>
                  <option value=">=">&gt;=</option>
                  <option value="<">&lt;</option>
                  <option value="<=">&lt;=</option>
                  <option value="==">==</option>
                  <option value="!=">!=</option>
                </select>

                <input
                  type="number"
                  step="0.01"
                  value={condition.value || ''}
                  onChange={(e) => updateCondition(index, 'value', parseFloat(e.target.value))}
                  className="px-2 py-1.5 bg-gray-50 border border-gray-200 rounded text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  placeholder="值"
                />
              </div>

              <button
                onClick={() => removeCondition(index)}
                disabled={alertConditions.conditions.length <= 1}
                className="mt-1 p-1.5 text-red-500 hover:bg-red-50 rounded disabled:opacity-30 disabled:cursor-not-allowed"
              >
                <MinusCircle className="w-4 h-4" />
              </button>
            </div>
          ))}
        </div>

        <button
          onClick={addCondition}
          className="w-full px-4 py-2 border-2 border-dashed border-blue-300 text-blue-600 rounded-lg hover:bg-blue-50 transition-colors flex items-center justify-center space-x-2 text-sm"
        >
          <PlusCircle className="w-4 h-4" />
          <span>添加条件</span>
        </button>

        <p className="text-xs text-gray-500">
          {alertConditions.logic === 'AND' ? '需要所有条件同时满足才触发预警' : '任意一个条件满足即触发预警'}
        </p>
      </div>
    );
  };

  const handleCreateAlert = () => {
    console.log('handleCreateAlert called');
    setEditingAlert(null);
    setSelectedTemplate('');
    setFormData({
      stockCode: '',
      stockName: '',
      alertType: '',
      conditions: '',
      sensitivity: 0.001,
      cooldownHours: 1,
      postTriggerAction: 'continue',
      enableSound: true,
      enableDesktop: true,
    });
    setAlertConditions({
      logic: 'AND',
      conditions: [{ field: '', operator: '', value: 0 }]
    });
    setStockCodeError(null);
    console.log('Setting showModal to true');
    setShowModal(true);
  };

  const handleEditAlert = (alert: PriceAlert) => {
    setEditingAlert(alert);
    setSelectedTemplate(alert.templateId || '');
    setFormData({
      stockCode: alert.stockCode,
      stockName: alert.stockName,
      alertType: alert.alertType,
      conditions: alert.conditions,
      sensitivity: alert.sensitivity,
      cooldownHours: alert.cooldownHours,
      postTriggerAction: alert.postTriggerAction,
      enableSound: alert.enableSound,
      enableDesktop: alert.enableDesktop,
    });

    // 解析JSON条件
    try {
      const parsed = JSON.parse(alert.conditions);
      setAlertConditions(parsed);
    } catch (err) {
      console.error('Failed to parse conditions:', err);
      setAlertConditions({
        logic: 'AND',
        conditions: [{ field: '', operator: '', value: 0 }]
      });
    }

    setStockCodeError(null);
    setShowModal(true);
  };

  const handleDeleteAlert = async (id: number) => {
    if (!confirm('确定要删除这个预警吗？')) return;

    try {
      const res = await deletePriceAlert(id);
      if (res?.success) {
        await fetchData();
      } else {
        setError(res?.message || '删除失败');
      }
    } catch (err) {
      const errorResult = parseError(err);
      setError(errorResult.message);
    }
  };

  const handleToggleAlert = async (id: number, isActive: boolean) => {
    try {
      const res = await togglePriceAlert(id, isActive);
      if (res?.success) {
        await fetchData();
      } else {
        setError(res?.message || '切换状态失败');
      }
    } catch (err) {
      const errorResult = parseError(err);
      setError(errorResult.message);
    }
  };

  const handleSaveAlert = async () => {
    if (!formData.stockCode || !formData.alertType) {
      setError('请填写完整信息');
      return;
    }

    // 序列化预警条件
    const conditionsJSON = serializeConditions();

    try {
      let res;
      if (editingAlert) {
        // 更新
        const updateData = {
          ...formData,
          id: editingAlert.id,
          isActive: editingAlert.isActive,
          conditions: conditionsJSON
        };
        res = await updatePriceAlert(JSON.stringify(updateData));
      } else {
        // 从模板创建或直接创建
        if (selectedTemplate) {
          res = await createPriceAlertFromTemplate(
            selectedTemplate,
            formData.stockCode,
            formData.stockName,
            JSON.stringify({
              sensitivity: formData.sensitivity,
              cooldownHours: formData.cooldownHours,
              postTriggerAction: formData.postTriggerAction,
              enableSound: formData.enableSound,
              enableDesktop: formData.enableDesktop,
            })
          );
        } else {
          const createData = {
            ...formData,
            conditions: conditionsJSON
          };
          res = await createPriceAlert(JSON.stringify(createData));
        }
      }

      if (res?.success) {
        setShowModal(false);
        await fetchData();
      } else {
        setError(res?.message || '保存失败');
      }
    } catch (err) {
      const errorResult = parseError(err);
      setError(errorResult.message);
    }
  };

  const handleTemplateSelect = (templateId: string) => {
    setSelectedTemplate(templateId);
    const template = templates.find(t => t.id === templateId);
    if (template) {
      setFormData({
        ...formData,
        alertType: template.alertType,
        conditions: template.conditions,
      });
    }
  };

  const getAlertTypeLabel = (type: string) => {
    const labels: Record<string, string> = {
      'price_change': '涨跌幅预警',
      'target_price': '目标价预警',
      'stop_loss': '止损价预警',
      'high_low': '突破高低点',
      'price_range': '价格区间预警',
      'ma_deviation': '均线偏离预警',
      'combined': '组合预警',
    };
    return labels[type] || type;
  };

  const getPostActionLabel = (action: string) => {
    const labels: Record<string, string> = {
      'continue': '继续监控',
      'disable': '触发后禁用',
      'once': '仅触发一次',
    };
    return labels[action] || action;
  };

  const filteredAlerts = alerts.filter(alert => {
    if (filterCode && !alert.stockCode.toLowerCase().includes(filterCode.toLowerCase())) {
      return false;
    }
    if (filterActive !== null && alert.isActive !== filterActive) {
      return false;
    }
    return true;
  });

  if (loading) {
    return (
      <div className="min-h-screen bg-slate-50 flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-blue-100 border-t-blue-600 rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-600">加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-50 p-6">
      {/* 错误提示 */}
      {error && (
        <div className="mb-6 bg-red-50 border border-red-200 rounded-xl p-4 flex items-center">
          <XCircle className="w-5 h-5 text-red-600 mr-3 flex-shrink-0" />
          <p className="text-red-700 text-sm">{error}</p>
          <button onClick={() => setError(null)} className="ml-auto text-red-400 hover:text-red-600">
            <XCircle className="w-4 h-4" />
          </button>
        </div>
      )}

      {/* 标题栏 */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-800 mb-2 flex items-center">
          <Bell className="w-8 h-8 mr-3 text-blue-600" />
          价格预警中心
        </h1>
        <p className="text-gray-600">
          设置股票价格预警，实时监控市场变化，及时捕捉交易机会
        </p>
      </div>

      {/* 标签页切换 */}
      <div className="bg-white rounded-2xl shadow-sm border border-gray-200 mb-6">
        <div className="border-b border-gray-100">
          <div className="flex space-x-8 px-6">
            <button
              onClick={() => { setActiveTab('alerts'); }}
              className={`py-4 px-1 border-b-2 font-medium transition-colors ${
                activeTab === 'alerts'
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              <div className="flex items-center space-x-2">
                <Bell className="w-4 h-4" />
                <span>预警列表</span>
                <span className="bg-blue-100 text-blue-700 text-xs px-2 py-0.5 rounded-full">
                  {alerts.length}
                </span>
              </div>
            </button>
            <button
              onClick={() => { setActiveTab('templates'); }}
              className={`py-4 px-1 border-b-2 font-medium transition-colors ${
                activeTab === 'templates'
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              <div className="flex items-center space-x-2">
                <Shield className="w-4 h-4" />
                <span>预警模板</span>
                <span className="bg-green-100 text-green-700 text-xs px-2 py-0.5 rounded-full">
                  {templates.length}
                </span>
              </div>
            </button>
            <button
              onClick={() => { setActiveTab('history'); loadTriggerHistory(); }}
              className={`py-4 px-1 border-b-2 font-medium transition-colors ${
                activeTab === 'history'
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              <div className="flex items-center space-x-2">
                <History className="w-4 h-4" />
                <span>触发历史</span>
              </div>
            </button>
          </div>
        </div>

        {/* 工具栏 */}
        <div className="p-6 flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <div className="relative">
              <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
              <input
                type="text"
                placeholder="搜索股票代码..."
                value={filterCode}
                onChange={(e) => setFilterCode(e.target.value)}
                className="pl-10 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 w-64"
              />
            </div>
            <select
              value={filterActive === null ? 'all' : filterActive ? 'active' : 'inactive'}
              onChange={(e) => setFilterActive(e.target.value === 'all' ? null : e.target.value === 'active')}
              className="px-4 py-2 bg-gray-50 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            >
              <option value="all">全部状态</option>
              <option value="active">已启用</option>
              <option value="inactive">已禁用</option>
            </select>
          </div>
          <button
            onClick={handleCreateAlert}
            className="flex items-center space-x-2 bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-xl transition-colors shadow-lg shadow-blue-600/20"
          >
            <Plus className="w-4 h-4" />
            <span className="font-medium">创建预警</span>
          </button>
        </div>
      </div>

      {/* 内容区域 */}
      {activeTab === 'alerts' && (
        <div className="space-y-4">
          {filteredAlerts.length === 0 ? (
            <div className="bg-white rounded-2xl shadow-sm border border-gray-200 p-12 text-center">
              <Bell className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <h3 className="text-xl font-semibold text-gray-800 mb-2">暂无预警</h3>
              <p className="text-gray-500 mb-6">点击"创建预警"按钮，为您的股票设置价格预警</p>
            </div>
          ) : (
            filteredAlerts.map((alert) => (
              <div key={alert.id} className="bg-white rounded-2xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center space-x-3 mb-3">
                      <div className={`p-2 rounded-xl ${alert.isActive ? 'bg-green-50 text-green-600' : 'bg-gray-50 text-gray-400'}`}>
                        {alert.isActive ? <Play className="w-5 h-5" /> : <Pause className="w-5 h-5" />}
                      </div>
                      <div>
                        <h3 className="text-lg font-bold text-gray-800">
                          {alert.stockName}
                          <span className="ml-2 text-sm font-normal text-gray-500 font-mono">
                            ({alert.stockCode})
                          </span>
                        </h3>
                        <span className="text-xs font-medium text-blue-600 bg-blue-50 px-2 py-0.5 rounded">
                          {getAlertTypeLabel(alert.alertType)}
                        </span>
                      </div>
                    </div>

                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4 text-sm">
                      <div className="flex items-center text-gray-600">
                        <Target className="w-4 h-4 mr-2 text-blue-500" />
                        <span>触发后行为: {getPostActionLabel(alert.postTriggerAction)}</span>
                      </div>
                      <div className="flex items-center text-gray-600">
                        <Clock className="w-4 h-4 mr-2 text-blue-500" />
                        <span>冷却时间: {alert.cooldownHours}小时</span>
                      </div>
                      <div className="flex items-center text-gray-600">
                        <Shield className="w-4 h-4 mr-2 text-blue-500" />
                        <span>灵敏度: {(alert.sensitivity * 100).toFixed(2)}%</span>
                      </div>
                      <div className="flex items-center text-gray-600">
                        <Bell className="w-4 h-4 mr-2 text-blue-500" />
                        <span>声音: {alert.enableSound ? '开' : '关'} | 桌面: {alert.enableDesktop ? '开' : '关'}</span>
                      </div>
                    </div>

                    {alert.lastTriggeredAt && (
                      <div className="text-xs text-gray-500 flex items-center">
                        <History className="w-3 h-3 mr-1" />
                        最后触发: {new Date(alert.lastTriggeredAt).toLocaleString('zh-CN')}
                      </div>
                    )}
                  </div>

                  <div className="flex items-center space-x-2 ml-4">
                    <button
                      onClick={() => handleToggleAlert(alert.id, !alert.isActive)}
                      className={`p-2 rounded-xl transition-colors ${
                        alert.isActive
                          ? 'bg-orange-50 text-orange-600 hover:bg-orange-100'
                          : 'bg-green-50 text-green-600 hover:bg-green-100'
                      }`}
                      title={alert.isActive ? '禁用' : '启用'}
                    >
                      {alert.isActive ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
                    </button>
                    <button
                      onClick={() => handleEditAlert(alert)}
                      className="p-2 rounded-xl bg-blue-50 text-blue-600 hover:bg-blue-100 transition-colors"
                      title="编辑"
                    >
                      <Edit3 className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDeleteAlert(alert.id)}
                      className="p-2 rounded-xl bg-red-50 text-red-600 hover:bg-red-100 transition-colors"
                      title="删除"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      )}

      {activeTab === 'templates' && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {templates.map((template) => (
            <div
              key={template.id}
              className="bg-white rounded-2xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow cursor-pointer"
              onClick={() => {
                setSelectedTemplate(template.id);
                handleCreateAlert();
                handleTemplateSelect(template.id);
              }}
            >
              <div className="flex items-center space-x-3 mb-4">
                <div className="p-3 bg-blue-50 rounded-xl">
                  <Shield className="w-6 h-6 text-blue-600" />
                </div>
                <div>
                  <h3 className="font-bold text-gray-800">{template.name}</h3>
                  <span className="text-xs text-blue-600">{getAlertTypeLabel(template.alertType)}</span>
                </div>
              </div>
              <p className="text-sm text-gray-600">{template.description}</p>
            </div>
          ))}
        </div>
      )}

      {activeTab === 'history' && (
        <div className="space-y-4">
          {histories.length === 0 ? (
            <div className="bg-white rounded-2xl shadow-sm border border-gray-200 p-12 text-center">
              <History className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <h3 className="text-xl font-semibold text-gray-800 mb-2">暂无触发记录</h3>
              <p className="text-gray-500">预警触发后记录会显示在这里</p>
            </div>
          ) : (
            histories.map((history) => (
              <div key={history.id} className="bg-white rounded-2xl shadow-sm border border-gray-200 p-6">
                <div className="flex items-start justify-between">
                  <div className="flex items-start space-x-4">
                    <div className="p-2 bg-green-50 rounded-xl">
                      <CheckCircle className="w-5 h-5 text-green-600" />
                    </div>
                    <div>
                      <h3 className="font-bold text-gray-800">
                        {history.stockName}
                        <span className="ml-2 text-sm font-normal text-gray-500 font-mono">
                          ({history.stockCode})
                        </span>
                      </h3>
                      <p className="text-sm text-gray-600 mt-1">{history.triggerMessage}</p>
                      <div className="flex items-center space-x-4 mt-2 text-xs text-gray-500">
                        <span className="font-mono">触发价格: ¥{history.triggerPrice.toFixed(2)}</span>
                        <span>{new Date(history.triggeredAt).toLocaleString('zh-CN')}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      )}

      {/* 创建/编辑预警弹窗 */}
      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
            <div className="p-6 border-b border-gray-200">
              <h2 className="text-2xl font-bold text-gray-800">
                {editingAlert ? '编辑预警' : '创建预警'}
              </h2>
            </div>

            <div className="p-6 space-y-6">
              {/* 预警模板选择（仅创建时显示） */}
              {!editingAlert && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    选择预警模板（可选）
                  </label>
                  <select
                    value={selectedTemplate}
                    onChange={(e) => {
                      setSelectedTemplate(e.target.value);
                      if (e.target.value) {
                        handleTemplateSelect(e.target.value);
                      }
                    }}
                    className="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  >
                    <option value="">自定义预警</option>
                    {templates.map((template) => (
                      <option key={template.id} value={template.id}>
                        {template.name} - {template.description}
                      </option>
                    ))}
                  </select>
                </div>
              )}

              {/* 股票信息 */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    股票代码 *
                  </label>
                  <div className="relative">
                    <input
                      type="text"
                      value={formData.stockCode}
                      onChange={(e) => setFormData({ ...formData, stockCode: e.target.value })}
                      onBlur={(e) => handleStockCodeBlur(e.target.value)}
                      disabled={searchingStock}
                      className={`w-full px-4 py-2 bg-gray-50 border ${stockCodeError ? 'border-red-300' : 'border-gray-200'} rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 pr-10`}
                      placeholder="例如: 600519"
                    />
                    {searchingStock && (
                      <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
                      </div>
                    )}
                  </div>
                  {stockCodeError && (
                    <p className="text-xs text-red-500 mt-1">{stockCodeError}</p>
                  )}
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    股票名称
                  </label>
                  <input
                    type="text"
                    value={formData.stockName}
                    onChange={(e) => setFormData({ ...formData, stockName: e.target.value })}
                    className="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                    placeholder="自动查询或手动输入"
                  />
                </div>
              </div>

              {/* 预警类型 */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  预警类型 *
                </label>
                <select
                  value={formData.alertType}
                  onChange={(e) => handleAlertTypeChange(e.target.value)}
                  className="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                >
                  <option value="">请选择预警类型</option>
                  <option value="price_change">涨跌幅预警</option>
                  <option value="target_price">目标价预警</option>
                  <option value="stop_loss">止损价预警</option>
                  <option value="high_low">突破高低点</option>
                  <option value="price_range">价格区间预警</option>
                  <option value="ma_deviation">均线偏离预警</option>
                  <option value="combined">组合预警</option>
                </select>
              </div>

              {/* 预警条件（表单化） */}
              {formData.alertType && (
                <div className="bg-gray-50 rounded-xl p-4 border border-gray-200">
                  <h4 className="text-sm font-semibold text-gray-700 mb-3">预警条件</h4>
                  {renderAlertConditions()}
                </div>
              )}

              {/* 配置项 */}
              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    灵敏度
                  </label>
                  <input
                    type="number"
                    step="0.001"
                    min="0"
                    max="0.1"
                    value={formData.sensitivity}
                    onChange={(e) => setFormData({ ...formData, sensitivity: parseFloat(e.target.value) })}
                    className="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    冷却时间（小时）
                  </label>
                  <input
                    type="number"
                    min="0"
                    max="24"
                    value={formData.cooldownHours}
                    onChange={(e) => setFormData({ ...formData, cooldownHours: parseInt(e.target.value) })}
                    className="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    触发后行为
                  </label>
                  <select
                    value={formData.postTriggerAction}
                    onChange={(e) => setFormData({ ...formData, postTriggerAction: e.target.value })}
                    className="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  >
                    <option value="continue">继续监控</option>
                    <option value="disable">触发后禁用</option>
                    <option value="once">仅触发一次</option>
                  </select>
                </div>
              </div>

              {/* 通知设置 */}
              <div className="flex items-center space-x-6">
                <label className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    checked={formData.enableSound}
                    onChange={(e) => setFormData({ ...formData, enableSound: e.target.checked })}
                    className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <span className="text-sm text-gray-700">启用声音提醒</span>
                </label>
                <label className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    checked={formData.enableDesktop}
                    onChange={(e) => setFormData({ ...formData, enableDesktop: e.target.checked })}
                    className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <span className="text-sm text-gray-700">启用桌面通知</span>
                </label>
              </div>
            </div>

            <div className="p-6 border-t border-gray-200 flex justify-end space-x-3">
              <button
                onClick={() => setShowModal(false)}
                className="px-6 py-2 border border-gray-300 text-gray-700 rounded-xl hover:bg-gray-50 transition-colors"
              >
                取消
              </button>
              <button
                onClick={handleSaveAlert}
                className="px-6 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition-colors"
              >
                {editingAlert ? '更新' : '创建'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PriceAlertPage;
