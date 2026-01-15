# Implementation Plan: Refactor to Repository Pattern

The goal is to align the recently added `stock_money_flow_hist` and `stock_strategy_signals` related logic with the existing MVC/Repository pattern of the project.

## 1. Create `MoneyFlowRepository`
**File:** `repositories/money_flow_repository.go` (New)

- **Struct:** `MoneyFlowRepository` wrapping `*sql.DB`.
- **Methods:**
  - `SaveMoneyFlows(flows []models.MoneyFlowData) error`: Handles batch `INSERT OR REPLACE` into `stock_money_flow_hist`.
  - `GetMoneyFlowHistory(code string, limit int) ([]models.MoneyFlowData, error)`: Fetches recent history for strategy calculation.
  - `SaveStrategySignal(signal *models.StrategySignal) error`: Handles `INSERT` into `stock_strategy_signals`.

## 2. Refactor `MoneyFlowService`
**File:** `services/money_flow_service.go`

- Update struct to use `repositories.MoneyFlowRepository` instead of direct `dbService`.
- **Method `FetchAndSaveHistory`**:
  - Keep API fetching logic in Service.
  - Move DB saving logic to `MoneyFlowRepository.SaveMoneyFlows`.

## 3. Refactor `StrategyService`
**File:** `services/strategy_service.go`

- Update struct to use `repositories.MoneyFlowRepository` (for reading history and saving signals).
- **Method `CalculateBuildSignals`**:
  - Use `MoneyFlowRepository.GetMoneyFlowHistory` to get data.
  - Use `MoneyFlowRepository.SaveStrategySignal` to persist results.

## 4. Define Models
**File:** `models/money_flow.go` (New or update existing if appropriate)

- Move `MoneyFlowData` and `StrategySignal` structs to the `models` package to avoid circular dependencies and ensure clean architecture.

## 5. Dependency Injection Update
**File:** `app.go`

- Initialize `MoneyFlowRepository`.
- Inject `MoneyFlowRepository` into `MoneyFlowService` and `StrategyService`.

## 6. Cleanup
- Remove raw SQL queries from `services/money_flow_service.go` and `services/strategy_service.go`.
- Fix imports.

## 7. Verification
- Update tests (`services/money_flow_service_test.go`, `services/strategy_service_test.go`) to reflect the architectural changes (or ensure they still pass if they test public methods).
