## 1. Implementation
- [ ] 1.1 添加 YAML 配置读取模块：支持读取 `config.yaml`（可执行文件目录优先，其次工作目录），并回退环境变量
- [ ] 1.2 调整 AI 服务初始化：改为接收配置入参，移除硬编码 key
- [ ] 1.3 调整应用启动流程：在 `startup()` 时加载配置并初始化 `AIService`；初始化失败时给出明确错误
- [ ] 1.4 添加 `config.yaml.example`，并在 `.gitignore` 忽略 `config.yaml`
- [ ] 1.5 更新 `README.md`：新增 YAML 配置说明与回退规则
- [ ] 1.6 运行 `openspec validate update-dashscope-config-to-yaml --strict` 并修复问题

