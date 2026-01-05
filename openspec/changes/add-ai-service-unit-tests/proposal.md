# Change: 为 AIService 的 extractSection 添加单元测试

## Why
`services/ai_service.go` 中的 `extractSection` 负责从大模型返回文本中提取章节内容，是报告解析的关键纯逻辑。补充单元测试可以锁定边界条件，避免后续改动导致解析回归。

## What Changes
- 为 `extractSection` 增加表驱动单元测试，文件为 `services/ai_service_test.go`。
- 对 `extractSection` 做最小可测试性重构（不改变对外调用方式），使边界行为更明确、便于测试。

## Impact
- Affected specs: `ai-testing`
- Affected code:
  - `services/ai_service.go`
  - `services/ai_service_test.go`

