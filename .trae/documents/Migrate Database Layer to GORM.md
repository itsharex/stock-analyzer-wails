I will help you migrate the entire database layer from raw `database/sql` to **GORM** with SQLite.

### Technical Plan

1.  **Dependency Management**
    *   Remove `modernc.org/sqlite` (direct usage).
    *   Add `gorm.io/gorm`.
    *   Add `github.com/glebarez/sqlite` (Pure Go SQLite driver, compatible with GORM and Wails cross-compilation).

2.  **Model Layer (Define GORM Schemas)**
    *   I will update `models/` or create `models/entities/` to define structs with GORM tags (`gorm:"primaryKey"`, `gorm:"index"`, etc.) for all existing tables:
        *   `Watchlist`
        *   `Alert`, `AlertHistory`
        *   `Position`
        *   `Config`, `SyncHistory`
        *   `Stock` (Market data)
        *   `PriceThresholdAlert`, `PriceAlertTemplate`, `PriceAlertTriggerHistory`
        *   `StockMoneyFlowHist`, `StockStrategySignal`
    *   Define a `KLine` struct for dynamic table usage.

3.  **Service Layer Refactor (`services/db_service.go`)**
    *   Change `DBService` to hold `*gorm.DB`.
    *   Replace `initTables` raw SQL with `db.AutoMigrate(...)`.
    *   Refactor K-Line operations to use `db.Table("kline_" + code)` for dynamic table queries and insertions.
    *   Implement efficient batch insertion for K-Lines using GORM's `CreateInBatches`.

4.  **Repository Layer Refactor (`repositories/*.go`)**
    *   Update all repositories (`WatchlistRepository`, `AlertRepository`, etc.) to use `*gorm.DB`.
    *   Replace `db.Query/Exec` with GORM's fluent API (`db.Find`, `db.Create`, `db.Where`, `db.Delete`).

5.  **Application Wiring (`app.go`)**
    *   Update dependency injection to pass the GORM instance.

### Execution Steps
1.  **Install Dependencies**: Run `go get` commands.
2.  **Create Entities**: Define the GORM models.
3.  **Migrate DB Service**: Switch `db_service.go` to GORM.
4.  **Migrate Repositories**: Rewrite each repository file one by one.
5.  **Verification**: Ensure the application compiles and runs (I will try to run a build or check syntax).

### Note
I will use `github.com/glebarez/sqlite` instead of the standard CGO driver to ensure your Wails application remains easy to compile on Windows without CGO requirements.
