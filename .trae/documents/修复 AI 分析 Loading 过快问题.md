# 修复方案：实现强制刷新机制

问题分析：点击“重新分析”时 Loading 消失太快，是因为后端直接返回了缓存的历史数据，没有进行真正的 AI 分析。

为了解决这个问题，我们需要实现“强制刷新”功能，让按钮点击时忽略缓存，重新调用 AI 生成数据。

## 1. 后端修改 (`services/ai_service.go`)
修改 `AnalyzeTechnical` 方法，增加对 `role` 参数的特殊处理：
- 检查 `role` 参数是否包含 `:force` 后缀。
- 如果包含（如 `technical:force`），则**跳过缓存读取步骤**，强制执行 AI 生成。
- 生成完成后，使用原始角色名（如 `technical`）更新缓存，确保下次页面加载时能获取到最新数据。

## 2. 前端修改 (`frontend/src/components/WatchlistDetail.tsx`)
修改 `handleAnalyze` 方法：
- 在调用 `analyzeTechnical` 时，传入第三个参数 `'technical:force'`。
- 这样点击按钮时会触发后端的强制刷新逻辑，前端会等待 AI 分析完成（约 10-20 秒），期间 Loading 动画会一直保持。

## 预期效果
- **页面加载时**：依然使用默认逻辑（优先读缓存），快速展示历史数据。
- **点击按钮时**：强制重新分析，Loading 动画持续直到 AI 返回最新结果，确保用户看到的是实时分析。
