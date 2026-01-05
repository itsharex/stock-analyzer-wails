## ADDED Requirements

### Requirement: 公共日志模块（zap）
系统 SHALL 提供公共日志模块封装 zap logger，并提供 `Debug/Info/Warn/Error` 等级的日志能力。

#### Scenario: 应用启动初始化日志模块成功
- **GIVEN** 应用启动
- **WHEN** 日志模块初始化
- **THEN** 产生日志输出且包含 `module` 与 `op` 字段

### Requirement: 日志输出（stdout + 文件）与文件滚动
系统 SHALL 同时将日志写入 stdout 与日志文件，并使用 `lumberjack` 提供文件滚动能力。

#### Scenario: 日志双写与滚动启用
- **GIVEN** 日志模块启用文件输出与滚动
- **WHEN** 产生日志
- **THEN** stdout 与日志文件均可观测到同一条日志记录

### Requirement: 错误路径补全（结构化字段）
系统 MUST 在关键错误路径记录足够详细的结构化日志，至少包含：`module`、`op`、`err`，并按场景追加关键业务字段（如 `stock_code`、`url`、`http_status`、`duration_ms` 等）。

#### Scenario: HTTP 请求失败
- **GIVEN** 后端发起外部 HTTP 请求
- **WHEN** 请求失败（如超时/连接失败/非预期响应）
- **THEN** 记录 `Error` 日志并包含 `url`、`method`、`timeout_ms`、`duration_ms`、`err`

#### Scenario: 解析失败
- **GIVEN** 后端收到外部响应并尝试解析
- **WHEN** 解析失败（如 JSON 反序列化错误）
- **THEN** 记录 `Error` 日志并包含 `body_size`、`duration_ms`、`err`

