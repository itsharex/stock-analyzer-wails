I will implement the following changes to add the industry filtering functionality:

### 1. Backend (Go)
*   **`services/stock_market_service.go`**:
    *   Add a new method `GetIndustries` to fetch the list of industries from the EastMoney API (`https://push2.eastmoney.com/api/qt/clist/get...`).
    *   Update `GetStocksList` to accept an `industry` parameter (string).
    *   Modify the SQL query in `GetStocksList` to filter by `industry` if the parameter is provided (and not empty).
*   **`controllers/stock_market_controller.go`**:
    *   Add a wrapper method `GetIndustries` to expose the service method.
    *   Update `GetStocksList` wrapper to accept the `industry` parameter.
*   **`app.go`**:
    *   Expose `GetIndustries` via the `App` struct.
    *   Update `GetStocksList` signature in `App` struct to match the controller/service change.

### 2. Frontend (React)
*   **`frontend/src/hooks/useWailsAPI.ts`**:
    *   Add `getIndustries` function.
    *   Update `getStocksList` to accept the `industry` parameter.
*   **`frontend/src/pages/StockListPage.tsx`**:
    *   Add state for `industries` (list of industries) and `selectedIndustry` (current filter).
    *   Add a `useEffect` to fetch industries on component mount.
    *   Add a dropdown (select) in the UI to allow users to choose an industry.
    *   Update the `loadStocks` function to pass `selectedIndustry` to the backend.
    *   Update the `handleSearch` and pagination logic to persist the selected industry filter.

This plan ensures that users can filter the stock list by industry using the data from the provided API.