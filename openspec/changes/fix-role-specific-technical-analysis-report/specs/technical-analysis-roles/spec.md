## ADDED Requirements

### Requirement: RoleSpecificTechnicalAnalysisReport
系统 MUST 支持深度分析的多角色输出，并保证不同角色对应不同的分析报告内容（至少风格/侧重点/结论之一有差异），避免角色选择失效。

#### Scenario: 同一股票不同角色返回不同报告
- **GIVEN** 同一只股票、同一分析周期（period）
- **WHEN** 用户分别选择 `technical` / `conservative` / `aggressive` 触发深度分析
- **THEN** 系统返回的分析报告内容 SHOULD 体现角色差异（例如策略倾向、仓位建议、风险措辞等）

### Requirement: RoleAwareTechnicalAnalysisCache
系统 MUST 按 `stock_code + role + period` 区分深度分析缓存，防止不同角色/周期互相复用导致展示错误。

#### Scenario: 角色缓存隔离
- **GIVEN** 同一股票同一周期下已生成 `technical` 的分析报告并写入缓存
- **WHEN** 用户切换到 `aggressive` 角色触发深度分析
- **THEN** 系统不得直接返回 `technical` 的缓存结果

#### Scenario: 周期缓存隔离
- **GIVEN** 同一股票同一角色下已生成 `daily` 的分析报告并写入缓存
- **WHEN** 用户切换到 `week` 周期触发深度分析
- **THEN** 系统不得直接返回 `daily` 的缓存结果


