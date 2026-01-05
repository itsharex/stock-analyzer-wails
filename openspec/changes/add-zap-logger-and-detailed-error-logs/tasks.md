## 1. Implementation
- [ ] 1.1 新增公共日志模块（建议 `internal/logger`）：封装 zap，提供 `Debug/Info/Warn/Error` 方法
- [ ] 1.2 支持 stdout + 文件双写，并通过 `lumberjack` 实现日志滚动；支持 `LOG_LEVEL` 等基础配置
- [ ] 1.3 在 `main.go` 初始化 logger，并在退出前 `Sync()`；替换 `log.Fatal`
- [ ] 1.4 在 `app.go`、`services/*` 的错误返回点补充结构化日志（字段含 module/op/err + 关键业务参数 + 耗时）
- [ ] 1.5 补充文档：README/部署文档中说明日志位置、级别与滚动策略

## 2. Validation
- [ ] 2.1 运行 `openspec validate add-zap-logger-and-detailed-error-logs --strict` 并修复问题

