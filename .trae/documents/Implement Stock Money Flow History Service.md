# Implementation Plan: Stock Money Flow History Service

I will create a new service to fetch and store historical money flow data for stocks, as per your requirements.

## 1. Database Schema Update
**File:** `services/db_service.go`

- Add a new table initialization step in `initTables` to create the `stock_money_flow_hist` table.
- **Table Schema:**
  - `code` (TEXT)
  - `trade_date` (TEXT)
  - `main_net` (REAL)
  - `super_net` (REAL)
  - `big_net` (REAL)
  - `mid_net` (REAL)
  - `small_net` (REAL)
  - `close_price` (REAL)
  - `chg_pct` (REAL)
  - **Primary Key:** `(code, trade_date)`
- **Index:** `idx_flow_date` on `(code, trade_date DESC)`

## 2. Create Money Flow Service
**File:** `services/money_flow_service.go` (New File)

- **Structs:**
  - `MoneyFlowService`: Holds `DBService` and `http.Client`.
  - `MoneyFlowData`: Internal struct to hold parsed data before DB insertion.
- **Key Methods:**
  - `NewMoneyFlowService(db *DBService)`: Constructor.
  - `FetchAndSaveHistory(code string) error`: Main logic.

### Logic for `FetchAndSaveHistory`
1.  **Generate `secid`**:
    - `6*` -> `1.code` (SH)
    - `0*`, `3*` -> `0.code` (SZ)
    - `4*`, `8*` -> `2.code` (BJ)
2.  **API Request**:
    - **URL:** `https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get`
    - **Query Params:**
      - `lmt=100`, `klt=101`
      - `fields1=f1,f2,f3,f7`
      - `fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65`
      - `secid={generated_secid}`
    - **Headers:** `Referer: https://quote.eastmoney.com/`
3.  **Data Parsing**:
    - Parse JSON response.
    - Iterate over `data.klines` (string array).
    - Split each string by comma `,`.
    - **Mapping:**
      - `[0]` -> `trade_date`
      - `[1]` -> `main_net`
      - `[5]` -> `super_net`
      - `[4]` -> `big_net`
      - `[3]` -> `mid_net`
      - `[2]` -> `small_net`
      - `[11]` -> `close_price`
      - `[12]` -> `chg_pct`
    - **Data Cleaning:** Convert `"-"` to `0.0`.
4.  **Database Storage**:
    - Begin Transaction (`db.Begin()`).
    - Use `INSERT OR REPLACE INTO ...` to ensure idempotency.
    - Commit Transaction (`tx.Commit()`).

## 3. Verification
- I will verify the implementation by:
  - Checking if the table is created correctly.
  - (Optional) If you wish, I can create a test case or a temporary main function to run a fetch for a specific stock (e.g., `002202` as in your example) and verify the data is saved.
