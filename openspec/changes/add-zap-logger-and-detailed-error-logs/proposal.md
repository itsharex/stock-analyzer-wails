# Change: 后端统一 zap 日志与错误路径补全

## Why
当前后端主要通过 `log.Fatal` 与 `fmt.Printf` 输出错误，缺乏统一格式与结构化字段，排障时难以定位“哪个步骤失败、输入是什么、外部依赖返回了什么、耗时多少”。需要引入统一日志能力（zap），并在关键错误路径补充足够详细的日志。

## What Changes
- 引入 `zap` 并封装公共日志模块（提供 `Debug/Info/Warn/Error`）。
- 日志输出同时写入 stdout 与文件，并支持文件滚动（`lumberjack`）。
- 在后端关键错误路径补全日志（尤其是外部调用、解析失败、初始化失败、关键参数校验失败）。

## Impact
- Affected specs: `logging`
- Affected code (initial scan):
  - `main.go`: `log.Fatal("启动应用失败:", err)`
  - `app.go`: 启动初始化错误使用 `fmt.Printf`；各 API 方法返回错误但缺少上下文日志
  - `services/stock_service.go`: HTTP 请求构建/发送/读取/解析失败与未找到股票等错误返回点
  - `services/ai_service.go`: AI 初始化失败、Generate 失败、chatModel 未初始化错误返回点
  - `services/config.go`: 配置文件查找/读取/解析失败与关键字段缺失

## Logging Design Notes
- 默认使用 JSON encoder（便于检索）；如需更适合开发调试的 console encoder，可后续通过环境变量/构建参数扩展。
- 文件输出默认路径建议：可执行文件目录下 `logs/app.log`（自动创建目录）。
- 统一字段建议：
  - `module`: 模块名（app/services/ai/services/stock/...）
  - `op`: 操作名/函数名
  - `err`: 错误对象（zap.Error）
  - `duration_ms`: 耗时（如有）
  - 与业务相关的关键字段：如 `stock_code`、`keyword`、`page_num`、`page_size`、`url`、`http_status`、`body_size`

