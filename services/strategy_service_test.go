package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"stock-analyzer-wails/repositories"
	"testing"
)

func TestCalculateBuildSignals(t *testing.T) {
	// Setup temporary DB
	tempDir, err := os.MkdirTemp("", "strategy_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Init tables
	svc := &DBService{db: db, dbPath: dbPath}
	if err := svc.initTables(); err != nil {
		t.Fatal(err)
	}

	// Create repositories
	strategyRepo := repositories.NewStrategyRepository(db)
	moneyFlowRepo := repositories.NewMoneyFlowRepository(db)

	strategySvc := NewStrategyService(strategyRepo, moneyFlowRepo)

	// Mock Data for a "Buy" signal
	// T-0: MainNet > 0, Close > MA20, Deviation < 3%, ChgPct 0.5-5%
	// Last 5 days: 3+ days MainNet > 0, Sum > 0, T-0 > 1.5 * Avg

	code := "000001"

	// Prepare 20 days of data
	// Day 0 (Today): MainNet=100, Close=10.2, ChgPct=2.0 (MA20 will be approx 10.0)
	// Day 1-4: MainNet=20, 20, -10, 20. AvgAbs(0..4) = (100+20+20+10+20)/5 = 34. 100 > 1.5*34 (51) -> OK
	// Day 5-19: Fill with dummy data to make MA20 approx 10.0

	// Insert T-0 to T-4
	// T-0
	insertHist(t, db, code, "2023-10-20", 100, 10.2, 2.0)
	// T-1
	insertHist(t, db, code, "2023-10-19", 20, 10.0, 0.0)
	// T-2
	insertHist(t, db, code, "2023-10-18", 20, 10.0, 0.0)
	// T-3
	insertHist(t, db, code, "2023-10-17", -10, 10.0, 0.0)
	// T-4
	insertHist(t, db, code, "2023-10-16", 20, 10.0, 0.0)

	// Insert T-5 to T-19 (15 days)
	// Make Close=10.0 mostly
	for i := 0; i < 15; i++ {
		date := fmt.Sprintf("2023-10-%02d", 15-i)   // simplistic date gen
		insertHist(t, db, code, date, 0, 9.98, 0.0) // Close slightly less than 10 to keep MA20 low enough
	}

	// Run Strategy
	signal, err := strategySvc.CalculateBuildSignals(code)
	if err != nil {
		t.Fatalf("CalculateBuildSignals failed: %v", err)
	}

	if signal == nil {
		t.Fatal("Expected a signal, but got nil")
	}

	if signal.SignalType != "B" {
		t.Errorf("Expected signal type B, got %s", signal.SignalType)
	}

	t.Logf("Signal Triggered: %+v", signal)
}

func insertHist(t *testing.T, db *sql.DB, code, date string, mainNet, closePrice, chgPct float64) {
	_, err := db.Exec(`
		INSERT INTO stock_money_flow_hist (code, trade_date, main_net, close_price, chg_pct)
		VALUES (?, ?, ?, ?, ?)
	`, code, date, mainNet, closePrice, chgPct)
	if err != nil {
		t.Fatal(err)
	}
}
