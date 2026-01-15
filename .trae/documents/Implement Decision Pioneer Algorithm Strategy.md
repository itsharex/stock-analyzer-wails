# Implementation Plan: Decision Pioneer Algorithm (Quantitative Strategy)

I will implement the "Decision Pioneer Algorithm" to detect "Build Position" signals based on money flow history and technical indicators.

## 1. Database Schema Update
**File:** `services/db_service.go`

- Add a new table `stock_strategy_signals` to `initTables`.
- **Table Schema:**
  - `id` (INTEGER PRIMARY KEY AUTOINCREMENT)
  - `code` (TEXT NOT NULL)
  - `trade_date` (TEXT NOT NULL)
  - `signal_type` (TEXT NOT NULL) -- 'B' for Buy
  - `strategy_name` (TEXT NOT NULL) -- '决策先锋'
  - `score` (REAL) -- Optional score
  - `details` (TEXT) -- JSON or text description of why it triggered
  - `created_at` (DATETIME)
  - **Unique Constraint:** `(code, trade_date, strategy_name)` to ensure idempotency.

## 2. Implement Strategy Logic
**File:** `services/strategy_service.go`

- Add a method `CalculateBuildSignals(code string) (*models.StrategySignal, error)`.
- **Logic Steps:**
  1.  **Fetch Data:** `SELECT * FROM stock_money_flow_hist WHERE code = ? ORDER BY trade_date DESC LIMIT 20`.
  2.  **Pre-check:** Ensure at least 20 records exist (for MA20).
  3.  **Dimension A: Funds (Last 5 days: indices 0 to 4)**
      - `Count(MainNet > 0) >= 3`
      - `Sum(MainNet) > 0`
      - `Abs(MainNet[0]) > 1.5 * Average(Abs(MainNet[0..4]))`
  4.  **Dimension B: Technical (MA20)**
      - Calculate `MA20` using `ClosePrice` of indices 0 to 19.
      - `ClosePrice[0] >= MA20`
      - `0 <= (ClosePrice[0] - MA20) / MA20 <= 0.03`
  5.  **Dimension C: Momentum**
      - `0.5 <= ChgPct[0] <= 5.0` (Assuming `chg_pct` is stored as percentage, e.g., 1.5)
  6.  **Persistence:**
      - If all conditions met, insert into `stock_strategy_signals`.

## 3. Wails Interface
**File:** `app.go`

- Add `RunStrategyScan(codes []string) []map[string]interface{}`.
- Loop through codes, call `CalculateBuildSignals`.
- Return list of triggered signals.

## 4. Verification
- Create a test in `services/strategy_service_test.go` to mock data and verify logic.
