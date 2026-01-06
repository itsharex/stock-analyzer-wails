## ADDED Requirements

### Requirement: DeepAnalysisAndEntryAnalysisButtonsCoexist
在自选股详情页中，系统 MUST 同时提供“深度分析”和“建仓分析”两个显式按钮入口，且两者互不遮挡、不会相互替代导致入口消失。

#### Scenario: 两按钮同时展示
- **WHEN** 用户打开自选股详情页
- **THEN** 页面展示“深度分析”按钮
- **AND** 页面展示“建仓分析”按钮

#### Scenario: 深度分析触发技术面分析
- **WHEN** 用户点击“深度分析”
- **THEN** 系统调用 `AnalyzeTechnical` 并展示 loading 状态
- **AND** 分析完成后更新技术分析结果展示区域

#### Scenario: 深度分析在数据未就绪时可理解提示
- **WHEN** K 线数据未就绪且用户尝试触发“深度分析”
- **THEN** 系统禁用按钮或给出明确提示（例如“请先加载K线数据/切换到K线视图”）


