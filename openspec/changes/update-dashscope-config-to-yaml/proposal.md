# Change: 启动时从 YAML 配置读取 DashScope Key/Model

## Why
当前 AI 配置依赖环境变量，且代码中存在对 `apiKey` 的硬编码覆盖，带来安全风险与部署维护成本。需要统一为“启动时加载配置文件”为主、环境变量为兼容回退的方式。

## What Changes
- 新增 `config.yaml`（不入库）配置读取机制，支持 `dashscope.api_key` 与 `dashscope.model`。
- 应用启动时加载配置，优先级：`config.yaml` → 环境变量（兼容现有 README 用法）。
- 当 `api_key` 最终仍缺失时，AI 服务初始化失败并给出明确错误信息。
- 移除代码中硬编码覆盖 `apiKey` 的逻辑。

## Impact
- Affected specs: `ai-configuration`
- Affected code:
  - `services/ai_service.go`
  - `app.go`
  - （新增）`services/config.go`
  - `README.md`
  - `.gitignore`
  - `config.yaml.example`

