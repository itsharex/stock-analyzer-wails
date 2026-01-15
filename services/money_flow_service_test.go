package services

import (
	"database/sql"
	"os"
	"path/filepath"
	"stock-analyzer-wails/repositories"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) (*DBService, *MoneyFlowService) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "stock_test")
	if err != nil {
		t.Fatal(err)
	}
	
	dbPath := filepath.Join(tempDir, "test.db")
	
	// Connect
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	
	svc := &DBService{db: db, dbPath: dbPath}
	
	// Init tables
	if err := svc.initTables(); err != nil {
		t.Fatal(err)
	}

	repo := repositories.NewMoneyFlowRepository(db)
	mfService := NewMoneyFlowService(repo)
	
	return svc, mfService
}

func cleanupTestDB(t *testing.T, svc *DBService) {
	svc.Close()
	os.RemoveAll(filepath.Dir(svc.dbPath))
}

func TestFetchAndSaveHistory(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=1 to run")
	}

	svc, mfService := setupTestDB(t)
	defer cleanupTestDB(t, svc)

	// Use a real stock code. 002202 (Goldwind)
	code := "002202"
	err := mfService.FetchAndSaveHistory(code)
	if err != nil {
		t.Fatalf("FetchAndSaveHistory failed: %v", err)
	}

	// Verify data in DB
	var count int
	err = svc.db.QueryRow("SELECT COUNT(*) FROM stock_money_flow_hist WHERE code = ?", code).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query count: %v", err)
	}
	
	if count == 0 {
		t.Error("Should have fetched some records, but got 0")
	} else {
		t.Logf("Fetched %d records for %s", count, code)
	}
    
    // Check fields of one record
    var mainNet, closePrice float64
    err = svc.db.QueryRow("SELECT main_net, close_price FROM stock_money_flow_hist WHERE code = ? LIMIT 1", code).Scan(&mainNet, &closePrice)
    if err != nil {
    	t.Fatalf("Failed to query row: %v", err)
    }
    
    if closePrice == 0 {
    	t.Error("Close price should not be 0")
    }
}
