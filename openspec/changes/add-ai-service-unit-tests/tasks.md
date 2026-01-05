## 1. Implementation
- [ ] 1.1 为 `extractSection` 增加最小可测试性重构（不改变业务语义）
- [ ] 1.2 新增 `services/ai_service_test.go`，为 `extractSection` 编写表驱动单测覆盖正常与边界场景

## 2. Validation
- [ ] 2.1 运行 `openspec validate add-ai-service-unit-tests --strict`
- [ ] 2.2 运行 `go test ./...`

