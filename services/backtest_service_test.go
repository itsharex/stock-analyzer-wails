package services

import (
	"testing"
)

// MockStockService 用于测试
// 我们不需要继承 StockService，只需要实现 BacktestService 依赖的接口
// 但 BacktestService 依赖 *StockService 具体类型，这导致难以 Mock。
// 这是一个设计上的改进点：BacktestService 应该依赖一个 Interface。
// 鉴于目前架构，我们无法轻易 Mock StockService 的方法，除非使用 GoMock 等工具或重构。
// 为了简单起见，我们只测试指标计算逻辑，这已经是目前能做的最好的单元测试。
// 如果要测试 RunBacktest，我们需要构造一个真实的 StockService，但这涉及数据库等。

// 不过，我们可以测试辅助函数，这已经覆盖了核心算法。
// 对于回测逻辑本身，由于它依赖 GetKLineData，如果 StockService 是具体结构体，确实难测。
// 幸运的是，StockService 的 GetKLineData 是一个方法。
// 我们可以临时创建一个 TestBacktestService 方法，把获取 K 线数据的逻辑解耦？
// 或者，我们可以信任 runBacktest 的逻辑，因为它是通用的。

// 让我们只测试指标计算，这对于本次任务（添加RSI）来说是最关键的。

func TestCalculateSMA(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	sma := calculateSMA(data, 3)
	// sma[0]=0, sma[1]=0, sma[2]=(1+2+3)/3=2, sma[3]=(2+3+4)/3=3, sma[4]=(3+4+5)/3=4
	expected := []float64{0, 0, 2, 3, 4}
	for i, v := range sma {
		if v != expected[i] {
			t.Errorf("Index %d: expected %f, got %f", i, expected[i], v)
		}
	}
}

func TestCalculateEMA(t *testing.T) {
	data := []float64{10, 11, 12, 13}
	// period = 2, alpha = 2/(2+1) = 0.666
	// ema[0] = 10
	// ema[1] = (11-10)*0.666 + 10 = 10.666
	// ema[2] = (12-10.666)*0.666 + 10.666 = 0.888 + 10.666 = 11.555
	
	ema := calculateEMA(data, 2)
	if ema[0] != 10 {
		t.Errorf("Expected 10, got %f", ema[0])
	}
	if ema[1] < 10.6 || ema[1] > 10.7 {
		t.Errorf("Expected ~10.66, got %f", ema[1])
	}
}

func TestCalculateRSI(t *testing.T) {
	// 简单的 RSI 测试
	// 构造一段连续上涨的数据，RSI 应该很高
	data := []float64{10, 11, 12, 13, 14, 15}
	rsi := calculateRSI(data, 5)
	// 只有最后一个点有值 (period=5, index 5)
	// gain avg = 1, loss avg = 0 -> RS = inf -> RSI = 100
	if rsi[5] != 100 {
		t.Errorf("Expected RSI 100 for uptrend, got %f", rsi[5])
	}
	
	// 构造一段震荡数据
	data2 := []float64{10, 9, 10, 9, 10, 9, 10}
	// period=2
	rsi2 := calculateRSI(data2, 2)
	// i=1: 10->9 (-1), avgG=0, avgL=1. rsi=0
	// i=2: 9->10 (+1), avgG=0.5, avgL=0.5. rs=1, rsi=50
	if rsi2[2] != 50 {
		t.Errorf("Expected RSI 50 for oscillating, got %f", rsi2[2])
	}
}
