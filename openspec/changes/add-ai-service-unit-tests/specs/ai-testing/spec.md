## ADDED Requirements

### Requirement: AIService extractSection 单元测试
系统 SHALL 为 `services/ai_service.go` 中用于章节提取的 `extractSection` 提供单元测试，且测试无需网络、无需环境变量即可运行。

#### Scenario: 提取章节成功
- **GIVEN** 一段包含 startMarker 与 endMarker 的文本
- **WHEN** 调用章节提取逻辑
- **THEN** 返回 startMarker 之后、endMarker 之前的文本片段

#### Scenario: 章节提取失败
- **GIVEN** 文本缺失 startMarker 或 marker 顺序不合法
- **WHEN** 调用章节提取逻辑
- **THEN** 返回空字符串

