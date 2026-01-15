# Implementation Plan: AI Verification for Quantitative Signals

I will implement the AI verification process for strategy signals as requested. This involves designing the prompt, implementing the `VerifySignal` method in `AIService`, integrating it into the strategy execution flow, and updating the database schema.

## 1. Database Schema Update
**File:** `services/db_service.go`

- Update `stock_strategy_signals` table schema to include `ai_score` and `ai_reason` fields.
- Add migration logic to add these columns if the table already exists.

## 2. AI Service Enhancement
**File:** `services/ai_service.go`

- **New Method:** `VerifySignal(stock *models.StockData, recentFlows []models.MoneyFlowData) (*models.AIVerificationResult, error)`
- **Prompt Design:**
  - Construct the prompt as specified, including stock info and recent 7 days money flow data.
  - Calculate "Main Inflow Ratio" (MainNet / Amount) and add to the data sent to AI.
  - Request JSON output with `score`, `opinion`, and `risk_level`.
- **Response Parsing:**
  - Use `extractJSON` helper to parse the AI response.
- **Concurrency Control:**
  - Implement a `VerifySignalAsync` method that uses a buffered channel (semaphore) to limit concurrent API calls (e.g., 5 concurrent requests).

## 3. Strategy Service Integration
**File:** `services/strategy_service.go`

- Update `CalculateBuildSignals` (or the caller in `app.go`) to trigger AI verification.
- **Note:** The user requested "Immediate async trigger after signal generation". Since `CalculateBuildSignals` returns a signal, it's better to handle the async call in `app.go` or a dedicated orchestration method to avoid blocking the strategy loop. However, to keep logic encapsulated, I can add an `onSignal` callback or handle it in `RunStrategyScan`.
- **Decision:** I will modify `app.go`'s `RunStrategyScan` to trigger the AI verification asynchronously after a signal is found.

## 4. App Logic Update
**File:** `app.go`

- In `RunStrategyScan`:
  - When a signal is generated (`CalculateBuildSignals` returns non-nil):
    - Launch a goroutine to call `aiService.VerifySignalAsync`.
    - On success: Update the database record with AI score and reason.
    - Emit a Wails event `signal_verified` to notify the frontend.

## 5. Mock Mechanism
**File:** `services/ai_service.go` / `config`

- Add a `EnableMock` flag in `AIResolvedConfig`.
- If enabled, `VerifySignal` returns dummy data without calling the LLM.

## 6. Verification
- Create a test case in `services/ai_service_test.go` to verify prompt construction and response parsing.
- (Manual) Run the strategy scan and check logs/DB for AI results.
