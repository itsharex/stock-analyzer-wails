## 1. Implementation
- [ ] 1.1 修复 `app.go`：`AnalyzeTechnical(code, period, role)` 调用 AI 层时传入正确的 `role`
- [ ] 1.2 校验并补齐 `services/ai_service.go` 的缓存维度：确保缓存 key 使用 `stock_code + role + period`
- [ ] 1.3 为深度分析链路补齐结构化日志字段：`stock_code`、`period`、`role`、`cache_hit`、`duration_ms`
- [ ] 1.4 （可选）前端增强：切换角色后展示对应角色的报告（可做简单内存缓存 map），避免快速切换时互相覆盖

## 2. Tests
- [ ] 2.1 后端单测：同一股票同一周期下，不同 `role` 生成路径/缓存 key 不相同（可通过 mock ChatModel + 断言缓存命中）

## 3. Validation
- [ ] 3.1 运行 `openspec validate fix-role-specific-technical-analysis-report --strict` 并修复问题


