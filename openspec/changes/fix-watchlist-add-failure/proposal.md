# Change: 修复自选股添加失败的 Bug

## Why
当前自选股添加功能存在多个潜在的失败点，导致用户无法成功添加股票到自选股列表。主要问题包括：

1. **错误处理不完善**：`getAllInternal()` 方法中 JSON 解析错误被忽略（第109行），可能导致数据损坏
2. **缺乏详细的错误日志**：文件操作失败时没有记录具体的错误信息，难以排查问题
3. **前端错误提示不友好**：只显示简单的 "添加失败" 提示，用户无法了解具体原因
4. **文件权限和目录创建问题**：虽然 `GetAppDataDir()` 会创建目录，但可能存在权限问题
5. **初始化错误被忽略**：`NewApp()` 中忽略了 `NewWatchlistService()` 的错误返回

## What Changes
- **BREAKING**: 无破坏性变更
- 完善 `watchlist_service.go` 中的错误处理逻辑
- 添加详细的结构化日志记录（基于现有的 zap 日志系统）
- 改进前端错误提示，显示具体的错误信息
- 增强文件操作的健壮性和错误恢复能力
- 修复应用初始化时的错误处理

## Impact
- Affected specs: `watchlist-management`（新建）
- Affected code:
  - `services/watchlist_service.go`: 错误处理和日志记录
  - `app.go`: 初始化错误处理和自选股操作日志
  - `frontend/src/components/StockSearch.tsx`: 错误提示改进
  - `frontend/src/hooks/useWailsAPI.ts`: 错误处理优化
