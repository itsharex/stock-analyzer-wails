package services

import (
	"fmt"
	"stock-analyzer-wails/repositories"
	"testing"
)

func TestFetchAndAlignHistory(t *testing.T) {
	// 1. 设置测试环境
	dbService, err := NewDBService()
	if err != nil {
		t.Fatalf("NewDBService failed: %v", err)
	}
	defer dbService.Close()
	stockMarketService := NewStockMarketService(dbService)

	moneyFlowRepository := repositories.NewMoneyFlowRepository(dbService.db)

	s := NewSyncService(dbService, stockMarketService, moneyFlowRepository)
	secid := "002202"
	klines, err := s.FetchHistoryFlowDataV2(secid, 120)
	if err != nil {
		t.Fatalf("FetchAndAlignHistory failed: %v", err)
	}
	if len(klines) == 0 {
		t.Fatalf("expected non-empty klines, got empty")
	}
	// 打印 2026-01-16 的数据验证对齐是否严谨
	targetDate := "2026-01-16"
	if d, ok := klines[targetDate]; ok {
		fmt.Printf("日期: %s\n", d.TradeDate)
		fmt.Printf("收盘价: %.2f\n", d.ClosePrice)
		fmt.Printf("主力净额: %.2f\n", d.MainNet)
		fmt.Printf("成交金额: %.2f\n", d.Amount)
		fmt.Printf("主力强度 (MainRate): %.2f%%\n", d.MainRate)
	} else {
		fmt.Println("未找到指定日期数据，请确认 limit 长度或当天是否停牌")
	}
	RunDecisionSignal(GetSortedData(klines))
}

func TestFetchAllDayTicks(t *testing.T) {
	// 1. 设置测试环境
	ticks, err := FetchAllDayTicks("600686")
	if err != nil {
		t.Fatalf("FetchAllDayTicks failed: %v", err)
		return
	}
	fmt.Println(ticks)

}
