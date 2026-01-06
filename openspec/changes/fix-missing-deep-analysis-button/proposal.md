# Change: 修复建仓分析上线后“深度分析”按钮缺失

## Why
近期在自选股详情页（`WatchlistDetail`）新增了“建仓分析”能力后，用户反馈**“深度分析按钮不见了”**，导致无法显式触发深度技术分析流程，影响核心使用路径与排障效率。

目前 UI 在右侧“AI 智能分析师”区域仅提供“建仓分析”按钮，而深度技术分析的触发入口被弱化为角色切换按钮（且并非所有用户能直观发现/理解）。

## What Changes
- 在 `WatchlistDetail` 右侧“AI 智能分析师”区域恢复**显式的“深度分析”按钮**，与“建仓分析”并存，互不遮挡。
- 明确交互语义：
  - “深度分析”用于触发 `AnalyzeTechnical`（深度技术分析）。
  - “建仓分析”用于触发 `AnalyzeEntryStrategy`（建仓方案）。
- 补齐按钮可用性规则与错误提示（例如数据尚未加载、正在分析中等）。

## Impact
- Affected specs: `analysis-ui-controls`（新建）
- Affected code（预估）:
  - `frontend/src/components/WatchlistDetail.tsx`: 恢复“深度分析”按钮与布局调整
  - `frontend/src/utils/errorHandler.ts`: 若需要，为“深度分析”补充更友好错误提示（可复用现有逻辑）

## Open Questions
- “深度分析”按钮是否仅在 `chartType === 'kline'` 时显示？
  - 建议：**始终显示**，但在 K 线数据未就绪时禁用并提示“请先加载K线数据/切到K线视图”。


