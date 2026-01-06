## 1. Implementation
- [ ] 1.1 在 `frontend/src/components/WatchlistDetail.tsx` 的“AI 智能分析师”区域新增“深度分析”按钮（与“建仓分析”并列）
- [ ] 1.2 “深度分析”按钮绑定 `handleAnalyze(role)` 并展示 loading 状态（复用 `analysisLoading`）
- [ ] 1.3 统一按钮可用性与提示：
  - [ ] 1.3.1 深度分析：K线数据未就绪时禁用并提示
  - [ ] 1.3.2 建仓分析：分时/资金流向未就绪时提示保持现状
- [ ] 1.4 视觉与布局：两按钮在窄宽度下不挤压（必要时改为上下两行或自适应栅格）

## 2. Validation
- [ ] 2.1 交互验收：
  - [ ] 2.1.1 点击“深度分析”可触发一次 `AnalyzeTechnical`
  - [ ] 2.1.2 点击“建仓分析”可触发一次 `AnalyzeEntryStrategy`
  - [ ] 2.1.3 两按钮同时存在且互不覆盖/不消失
  - [ ] 2.1.4 “深度分析”在无K线数据时给出可理解提示
- [ ] 2.2 运行 `openspec validate fix-missing-deep-analysis-button --strict` 并修复问题


