# AI 股票投资助手功能路线与技术方案

## 核心功能版图
- 实时行情与数据源：稳健 SSE、指数/ETF/行业板块覆盖、数据质量监控
- 持仓与资产管理：多账户、交易记录、盈亏分析、绩效指标（年化收益、最大回撤、Sharpe）
- 风险控制：仓位约束、止盈止损、行业/单票暴露、波动率与 VaR、回撤预警
- 预警与自动化：价格/涨跌幅/量能/形态/公告模板化预警，冷却与去重，系统通知/邮件/企业微信
- 策略与回测：指标库（MA/MACD/RSI/KDJ/布林带）、事件驱动回测、滑点/手续费模型、参数优化、绩效报告
- AI 分析与对话：自然语言筛选、新闻/公告情绪、研报摘要、解释型多维分析
- 研究数据接入：新闻 RSS、巨潮公告、财务与行业数据融合
- 可视化增强：多图层 K 线/分时、AI 绘图标注、信号叠加、交互筛选
- 性能与可靠性：缓存与离线同步、日志分级与采样、并发与资源治理、健康检查
- 安全与合规：隐私、免责声明、频控与熔断、敏感操作权限
- 可扩展性：策略插件化、脚本化、REST/gRPC 接口

## 技术实施要点与代码对齐
- 后端服务
  - SSE 与状态管理：参考 [stock_service.go](file:///c:/Users/o_wanyuqin/GolandProjects/stock-analyzer-wails/services/stock_service.go)
  - Wails 绑定入口：[app.go](file:///c:/Users/o_wanyuqin/GolandProjects/stock-analyzer-wails/app.go)
  - 持久化：DBService 与各 Repository（扩展持仓/交易结构）
- 前端
  - API 封装与事件订阅：[useWailsAPI.ts](file:///c:/Users/o_wanyuqin/GolandProjects/stock-analyzer-wails/frontend/src/hooks/useWailsAPI.ts)
  - 自选与详情 SSE 管理：[Watchlist.tsx](file:///c:/Users/o_wanyuqin/GolandProjects/stock-analyzer-wails/frontend/src/components/Watchlist.tsx)、[WatchlistDetail.tsx](file:///c:/Users/o_wanyuqin/GolandProjects/stock-analyzer-wails/frontend/src/components/WatchlistDetail.tsx)
  - 图表组件：K 线/分时扩展 AI 标注层
 - 现有控制器与服务参考
   - 策略控制器：[app.go:StrategyController](file:///c:/Users/o_wanyuqin/GolandProjects/stock-analyzer-wails/app.go#L148-L158)
   - 价格预警控制器转发接口：[app.go:PriceAlert*](file:///c:/Users/o_wanyuqin/GolandProjects/stock-analyzer-wails/app.go#L517-L591)
   - 持仓监控与移动止损：[app.go:startPositionMonitor](file:///c:/Users/o_wanyuqin/GolandProjects/stock-analyzer-wails/app.go#L231-L244)、[app.go:updateTrailingStop](file:///c:/Users/o_wanyuqin/GolandProjects/stock-analyzer-wails/app.go#L313-L357)
   - 价格预警引擎：[app.go:startAlertMonitor](file:///c:/Users/o_wanyuqin/GolandProjects/stock-analyzer-wails/app.go#L359-L373)、[app.go:checkAlerts](file:///c:/Users/o_wanyuqin/GolandProjects/stock-analyzer-wails/app.go#L374-L435)

## 分阶段落地计划
### 阶段1：实时行情与稳定性
- 每只股票连接状态与失败计数
- 重连与节流策略完善、健康事件
- 列表/详情统一清理并停止 SSE
 - 交付物：前后端状态事件、连接健康面板、日志降噪
### 阶段2：预警与通知
- 预警模板（价格/涨跌幅/量能/形态/公告）
- 冷却与去重、交易时段策略
- 通知渠道适配（系统/邮件/企业微信）
 - 交付物：预警中心、模板管理、触发历史与通知联动
### 阶段3：持仓与风险
- 扩展持仓结构（多账户、成本、税费、现金）与交易记录
- 绩效与风险度量（最大回撤、波动、VaR）
- 风控规则引擎（仓位上限、止损止盈、行业暴露）
 - 交付物：持仓仪表盘、风险面板、规则配置
### 阶段4：策略与回测
- 指标库与信号生成
- 事件驱动回测、滑点/手续费、复权
- 参数优化与绩效报告
 - 交付物：策略管理页、回测配置与报告
### 阶段5：AI 分析与对话
- NLP 查询到筛选执行
- 新闻/公告情绪与研报摘要
- 解释型报告与联动操作
 - 交付物：AI 助手面板、解释报告生成与操作联动
### 阶段6：数据接入与可视化
- RSS/公告接入与合并
- 图表 AI 标注层与信号叠加
 - 交付物：研究信息流整合、图形化洞见

## 数据与性能治理
- 缓存：热点行情内存缓存 + SQLite 落地
- 日志：分级与采样、错误节流、诊断事件
- 并发：按 code 限流与连接复用、统一 Stop 清理
- 可观测：连接数、失败率、重试次数指标与健康页
 - 性能基线：SSE 连接延迟、事件处理耗时、前端刷新 FPS
 - 容量规划：并发连接上限、数据库写入速率、缓存命中率
 - 故障回退：数据源不可用时的退化策略与提示

## 安全与合规
- 明示免责声明与风险提示
- 请求频控与熔断、接口时段控制
- 权限边界与敏感操作确认
 - 合规要点：不提供交易建议与收益承诺，标注数据来源与时效限制
 - 运维合规：访问频次控制、异常重试策略、IP 封禁回避与告警

## 交互体验要点
- 自然语言操作：筛选、对比、解释、设置预警
- 列表/详情一致的实时刷新与清理
- 清晰的连接与健康状态提示
 - 视图：自选行情概览、个股详情（K线/分时/资金流）、预警中心、持仓仪表盘、策略与回测、AI 助手面板
 - 交互：信号点击高亮、图表区间选择联动、AI 标注显隐切换、预警模板快速应用
