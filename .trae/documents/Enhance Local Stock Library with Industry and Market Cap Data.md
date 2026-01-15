I will implement the following changes to add the new fields (Industry, Region, Board, Market Cap) to your local stock library:

### 1. Database Schema Update (`services/db_service.go`)
-   Modify `initTables` to update the `stocks` table schema.
-   Add columns: `industry` (TEXT), `region` (TEXT), `board` (TEXT), `total_mv` (REAL), `circ_mv` (REAL).
-   Add auto-migration logic: Execute `ALTER TABLE stocks ADD COLUMN ...` for each new field if it doesn't exist, ensuring existing databases are updated without data loss.

### 2. Backend Service Update (`services/stock_market_service.go`)
-   Update `StockMarketData` struct to include the new fields: `Industry`, `Region`, `Board`, `TotalMV`, `CircMV`.
-   Update `SyncAllStocks` function:
    -   Request additional fields from EastMoney API: `f100` (Industry), `f102` (Region), `f103` (Board), `f20` (Total Market Cap), `f21` (Circulating Market Cap).
    -   Update `parseStockItem` to extract and map these fields.
    -   Update the `upsertSQL` to save these fields into the database.
-   Update `GetStocksList` to query and return these new fields.

### 3. Frontend Update
-   **Types (`frontend/src/types.ts`)**: Update `StockMarketData` interface to include the new properties.
-   **UI (`frontend/src/pages/StockListPage.tsx`)**:
    -   Add new columns to the stock list table: "Industry", "Region/Board", "Total Market Cap", "Circulating Market Cap".
    -   Add formatting for Market Cap (e.g., displaying in "Billions" for better readability).

This plan ensures that your local stock library is enriched with the requested data, and the system automatically handles the database schema update upon restart.