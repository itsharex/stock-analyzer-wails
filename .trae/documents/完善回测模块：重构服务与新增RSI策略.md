我已对现有的回测模块进行了全面分析。目前的实现虽然功能可用，但代码逻辑主要堆积在 `app.go` 中，且策略扩展性有限。

为了完善回测模块，我建议进行一次**重构与功能增强**，具体计划如下：

### 1. 后端重构：独立回测服务 (Go)
目前回测逻辑与应用程序主入口耦合严重。我将创建一个独立的 `BacktestService` 来专门处理回测逻辑。
*   **新建 `services/backtest_service.go`**:
    *   **提取通用回测引擎**: 将资金管理、交易记录、收益计算等公共逻辑提取为通用方法，避免重复代码。
    *   **迁移现有策略**: 将 `SimpleMA` (双均线) 和 `MACD` 策略逻辑迁移至该服务中。
    *   **新增 RSI 策略**: 实现相对强弱指标 (RSI) 回测策略。逻辑为：RSI < 超卖阈值(如30) 买入，RSI > 超买阈值(如70) 卖出。

### 2. 后端集成 (`app.go`)
*   在 `App` 中注入新的 `BacktestService`。
*   更新 `BacktestSimpleMA` 和 `BacktestMACD` 方法，改为调用 Service 层接口。
*   **新增 API**: `BacktestRSI(code, period, buyThreshold, sellThreshold, initialCapital, startDate, endDate)`。

### 3. 前端功能增强 (React)
*   **更新 API Hook**: 在 `useWailsAPI.ts` 中添加 `BacktestRSI` 方法定义。
*   **升级回测面板 (`BacktestPanelEnhanced.tsx`)**:
    *   **增加策略选项**: 在下拉菜单中增加 "RSI 策略"。
    *   **动态参数界面**: 当选择 RSI 策略时，动态展示 "RSI周期"、"买入阈值"、"卖出阈值" 等参数输入框。
    *   **适配调用**: 根据用户选择的策略类型，自动调用对应的后端接口。

### 预期效果
完成上述工作后，回测模块将拥有更好的代码结构，更容易扩展新策略，并且用户将立即可用一个新的 RSI 回测策略。
