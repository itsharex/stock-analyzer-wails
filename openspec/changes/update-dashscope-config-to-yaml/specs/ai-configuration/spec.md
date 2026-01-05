## ADDED Requirements

### Requirement: DashScope 配置加载（YAML 优先，环境变量回退）
系统 SHALL 在应用启动时加载 DashScope 配置，加载优先级为：`config.yaml` → 环境变量。

#### Scenario: 配置文件存在且字段齐全
- **GIVEN** 可执行文件目录（或工作目录）存在 `config.yaml`
- **AND** `dashscope.api_key` 与 `dashscope.model` 均为非空字符串
- **WHEN** 应用启动
- **THEN** 系统使用 `config.yaml` 中的值初始化 AI 服务

#### Scenario: 配置文件缺失，回退环境变量成功
- **GIVEN** 未找到 `config.yaml`
- **AND** 环境变量 `DASHSCOPE_API_KEY` 为非空字符串
- **WHEN** 应用启动
- **THEN** 系统使用环境变量初始化 AI 服务

### Requirement: DashScope API Key 缺失时的错误提示
若 `dashscope.api_key` 与 `DASHSCOPE_API_KEY` 均缺失或为空，系统 MUST 明确提示用户需要配置 API Key，且 AI 服务不得被视为已初始化。

#### Scenario: 配置与环境变量均缺失
- **GIVEN** `config.yaml` 缺失或 `dashscope.api_key` 为空
- **AND** 环境变量 `DASHSCOPE_API_KEY` 为空或不存在
- **WHEN** 应用启动并尝试初始化 AI 服务
- **THEN** 初始化失败并返回/记录包含“API Key 缺失/未配置”的错误信息

### Requirement: DashScope Model 缺省规则
若 `dashscope.model` 与 `DASHSCOPE_MODEL` 均缺失或为空，系统 SHALL 使用内置默认模型名初始化 AI 服务。

#### Scenario: 未配置 model
- **GIVEN** `dashscope.model` 为空或缺失
- **AND** 环境变量 `DASHSCOPE_MODEL` 为空或不存在
- **WHEN** 应用启动并初始化 AI 服务
- **THEN** 系统使用默认模型名初始化 AI 服务

