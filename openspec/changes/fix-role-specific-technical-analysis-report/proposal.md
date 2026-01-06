# Change: 深度分析不同角色返回不同分析报告

## Why
当前前端支持“稳健/技术/激进”三种深度分析角色，但用户发现**不同角色实际返回同一份分析报告**，导致角色差异化失效。

根因初步定位为后端参数传递错误：`App.AnalyzeTechnical(code, period, role)` 调用 AI 层时将 `period` 误作为 `role` 传入，导致 AI 层角色选择/缓存 key 不随用户选择变化。

## What Changes
- 修复后端 `AnalyzeTechnical` 参数传递：将 `role` 正确传入 `AIService.AnalyzeTechnical`。
- 明确行为规范：
  - 同一只股票同一周期下，不同角色 MUST 返回不同风格/内容侧重点的分析报告（至少提示语/结论/策略不同）。
  - 缓存 MUST 按 `stock_code + role + period` 区分（防止角色串档）。
- 补充日志字段，便于排查：在深度分析链路中记录 `stock_code`、`period`、`role`、`cache_hit`、`duration_ms`。

## Impact
- Affected specs: `technical-analysis-roles`（新建）
- Affected code:
  - `app.go`: `AnalyzeTechnical` 传参修正
  - `services/ai_service.go`: （如需）补齐日志与缓存 key 的 period 处理一致性
  - `frontend/src/components/WatchlistDetail.tsx`: （可选）切换角色时保留各自结果，避免 UI 覆盖（非必须）

## Non-Goals
- 不调整角色 prompt 的具体文案策略（除非需要更明显的差异化输出）。


