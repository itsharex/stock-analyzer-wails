package services

import (
	"fmt"
	"math"
	"stock-analyzer-wails/models"
)

// BacktestService 提供回测功能
type BacktestService struct {
	stockService *StockService
}

// NewBacktestService 创建新的 BacktestService
func NewBacktestService(stockService *StockService) *BacktestService {
	return &BacktestService{
		stockService: stockService,
	}
}

// SignalGenerator 是产生买卖信号的函数类型
// i: 当前K线的索引
// dates: 日期列表
// closes: 收盘价列表
// 返回: "BUY", "SELL" 或 "" (无操作)
type SignalGenerator func(i int, dates []string, closes []float64) string

// runBacktest 执行通用回测逻辑
func (s *BacktestService) runBacktest(
	code string,
	strategyName string,
	initialCapital float64,
	startDate string,
	endDate string,
	limit int,
	signalGen SignalGenerator,
) (*models.BacktestResult, error) {

	// 1. 获取数据
	klines, err := s.stockService.GetKLineData(code, limit, "daily")
	if err != nil {
		return nil, fmt.Errorf("获取K线失败: %w", err)
	}
	if len(klines) == 0 {
		return nil, fmt.Errorf("无K线数据")
	}

	// 2. 过滤区间 & 准备数据数组
	var dates []string
	var closes []float64
	// 为了确保指标计算准确，我们通常需要比 startDate 更早的数据
	// 这里简单处理：保留所有获取到的数据用于计算指标，但在回测循环中根据 startDate 过滤交易
	// 更好的做法是：获取足够多的历史数据计算指标，然后截取回测区间
	
	// 注意：为了简化逻辑和保持与旧版一致，我们这里先全部加载到 slice 中
	// 真正的日期过滤在交易循环中判断，或者先过滤数据再计算指标（但这会影响指标初期的准确性）
	// 旧版 app.go 是先过滤数据再计算指标，这其实会导致区间开始时的指标不准确（因为没有前置数据）。
	// 改进：我们使用全量数据计算指标，然后在回测循环中只在指定区间内交易。
	
	for _, k := range klines {
		dates = append(dates, k.Time)
		closes = append(closes, k.Close)
	}
	
	// 3. 执行回测循环
	cash := initialCapital
	units := 0.0
	inPosition := false
	entryPrice := 0.0
	trades := make([]models.TradeRecord, 0)
	equityCurve := make([]float64, 0)
	equityDates := make([]string, 0)

	// 确定回测的起始索引
	// 我们遍历所有数据，但只在时间范围内记录净值和交易
	for i := 0; i < len(closes); i++ {
		date := dates[i]
		price := closes[i]

		// 检查是否在回测区间内
		inRange := (startDate == "" || date >= startDate) && (endDate == "" || date <= endDate)
		
		if !inRange {
			// 如果还没到开始时间，保持初始状态
			// 如果已经过了结束时间，可以提前结束（但为了画图完整性，也可以继续算净值但不交易）
			if endDate != "" && date > endDate {
				break
			}
			continue
		}

		// 获取信号
		signal := signalGen(i, dates, closes)

		if signal == "BUY" && !inPosition {
			// 全仓买入
			if price > 0 {
				units = cash / price
				cash = 0
				inPosition = true
				entryPrice = price
				trades = append(trades, models.TradeRecord{
					Time: date, Type: "BUY", Price: price, Volume: 0, Amount: units * price,
				})
			}
		} else if signal == "SELL" && inPosition {
			// 全部卖出
			cash = cash + units*price
			profit := (price - entryPrice) * units
			trades = append(trades, models.TradeRecord{
				Time: date, Type: "SELL", Price: price, Volume: 0, Amount: units * price, Profit: profit,
			})
			units = 0
			inPosition = false
		}

		// 记录每日净值
		equity := cash + units*price
		equityCurve = append(equityCurve, equity)
		equityDates = append(equityDates, date)
	}

	// 如果最后仍持仓，按最后一天（回测区间内的最后一天）价格平仓
	if inPosition && len(equityDates) > 0 {
		lastDate := equityDates[len(equityDates)-1]
		lastPrice := 0.0
		// 找到 lastDate 对应的 price
		for i := len(closes) - 1; i >= 0; i-- {
			if dates[i] == lastDate {
				lastPrice = closes[i]
				break
			}
		}
		
		if lastPrice > 0 {
			cash = cash + units*lastPrice
			profit := (lastPrice - entryPrice) * units
			trades = append(trades, models.TradeRecord{
				Time: lastDate, Type: "SELL", Price: lastPrice, Volume: 0, Amount: units * lastPrice, Profit: profit,
			})
			units = 0
			inPosition = false
		}
	}
	
	// 如果区间内没有数据，equityCurve 可能为空
	if len(equityCurve) == 0 {
		return nil, fmt.Errorf("指定日期范围内没有有效交易数据")
	}

	final := cash
	// 注意：如果最后强制平仓了，final 已经包含了变现后的价值
	// 如果没平仓（逻辑上我们在上面强制平仓了），final 只是现金。
	// 这里 final 应该是最后时刻的总资产。
	// 上面的强制平仓逻辑已经把 units 变成了 cash，所以 final = cash 是对的。

	// 计算统计指标
	ret := final/initialCapital - 1
	
	// 年化收益
	annualized := 0.0
	days := len(equityCurve)
	if days > 0 {
		annualized = math.Pow(final/initialCapital, 252.0/float64(days)) - 1
	}

	// 最大回撤
	peak := equityCurve[0]
	maxDD := 0.0
	for _, v := range equityCurve {
		if v > peak {
			peak = v
		}
		dd := 0.0
		if peak > 0 {
			dd = (peak - v) / peak
		}
		if dd > maxDD {
			maxDD = dd
		}
	}

	// 胜率
	wins := 0
	finishedTrades := 0
	for _, t := range trades {
		if t.Type == "SELL" {
			finishedTrades++
			if t.Profit > 0 {
				wins++
			}
		}
	}
	winRate := 0.0
	if finishedTrades > 0 {
		winRate = float64(wins) / float64(finishedTrades)
	}

	return &models.BacktestResult{
		StrategyName:     strategyName,
		StockCode:        code,
		StartDate:        equityDates[0],
		EndDate:          equityDates[len(equityDates)-1],
		InitialCapital:   initialCapital,
		FinalCapital:     final,
		TotalReturn:      ret,
		AnnualizedReturn: annualized,
		MaxDrawdown:      maxDD,
		WinRate:          winRate,
		TradeCount:       finishedTrades,
		Trades:           trades,
		EquityCurve:      equityCurve,
		EquityDates:      equityDates,
	}, nil
}

// BacktestSimpleMA 双均线策略
func (s *BacktestService) BacktestSimpleMA(code string, shortPeriod int, longPeriod int, initialCapital float64, startDate string, endDate string) (*models.BacktestResult, error) {
	if shortPeriod <= 0 || longPeriod <= 0 || shortPeriod >= longPeriod {
		return nil, fmt.Errorf("参数错误: shortPeriod 必须 > 0 且 < longPeriod")
	}

	// 预先计算指标所需的闭包
	var shortMA, longMA []float64

	return s.runBacktest(code, fmt.Sprintf("SMA(%d,%d)", shortPeriod, longPeriod), initialCapital, startDate, endDate, 5000, 
		func(i int, dates []string, closes []float64) string {
			// 懒加载计算指标 (只计算一次)
			if shortMA == nil {
				shortMA = calculateSMA(closes, shortPeriod)
				longMA = calculateSMA(closes, longPeriod)
			}

			// 信号逻辑
			if i > 0 && shortMA[i-1] > 0 && longMA[i-1] > 0 && shortMA[i] > 0 && longMA[i] > 0 {
				prevCrossUp := shortMA[i-1] <= longMA[i-1] && shortMA[i] > longMA[i]
				prevCrossDown := shortMA[i-1] >= longMA[i-1] && shortMA[i] < longMA[i]

				if prevCrossUp {
					return "BUY"
				} else if prevCrossDown {
					return "SELL"
				}
			}
			return ""
		},
	)
}

// BacktestMACD MACD策略
func (s *BacktestService) BacktestMACD(code string, fastPeriod int, slowPeriod int, signalPeriod int, initialCapital float64, startDate string, endDate string) (*models.BacktestResult, error) {
	if fastPeriod <= 0 || slowPeriod <= 0 || fastPeriod >= slowPeriod {
		return nil, fmt.Errorf("参数错误: fastPeriod 必须 > 0 且 < slowPeriod")
	}

	var dif, dea []float64

	return s.runBacktest(code, fmt.Sprintf("MACD(%d,%d,%d)", fastPeriod, slowPeriod, signalPeriod), initialCapital, startDate, endDate, 5000,
		func(i int, dates []string, closes []float64) string {
			if dif == nil {
				dif, dea, _ = calculateMACD(closes, fastPeriod, slowPeriod, signalPeriod)
			}

			if i > 0 && dif[i-1] != 0 && dea[i-1] != 0 && dif[i] != 0 && dea[i] != 0 {
				crossUp := dif[i-1] <= dea[i-1] && dif[i] > dea[i]
				crossDown := dif[i-1] >= dea[i-1] && dif[i] < dea[i]

				if crossUp {
					return "BUY"
				} else if crossDown {
					return "SELL"
				}
			}
			return ""
		},
	)
}

// BacktestRSI RSI策略
func (s *BacktestService) BacktestRSI(code string, period int, buyThreshold float64, sellThreshold float64, initialCapital float64, startDate string, endDate string) (*models.BacktestResult, error) {
	if period <= 0 {
		return nil, fmt.Errorf("参数错误: period 必须 > 0")
	}
	if buyThreshold >= sellThreshold {
		return nil, fmt.Errorf("参数错误: buyThreshold 必须 < sellThreshold")
	}

	var rsi []float64

	return s.runBacktest(code, fmt.Sprintf("RSI(%d,%.0f,%.0f)", period, buyThreshold, sellThreshold), initialCapital, startDate, endDate, 5000,
		func(i int, dates []string, closes []float64) string {
			if rsi == nil {
				rsi = calculateRSI(closes, period)
			}

			if i > 0 && rsi[i] > 0 {
				// RSI 策略：
				// 低于阈值 (超卖) -> 买入
				// 高于阈值 (超买) -> 卖出
				
				// 简单的阈值突破策略
				// 可以优化为：从下方上穿 buyThreshold 买入，从上方下穿 sellThreshold 卖出
				// 或者：低于 buyThreshold 买入，高于 sellThreshold 卖出
				
				// 这里使用：
				// 买入：RSI < buyThreshold
				// 卖出：RSI > sellThreshold
				
				// 为了避免在超卖区域反复买入（虽然全仓模式下只能买一次），我们只在首次满足条件时触发
				// 但由于 runBacktest 内部会检查 !inPosition，所以这里只要返回 BUY 即可。
				
				// 改进策略逻辑：
				// 买入信号：RSI 跌破买入阈值 (寻找反弹机会，或者更稳健的是 RSI 从下向上突破买入阈值)
				// 卖出信号：RSI 突破卖出阈值
				
				// 采用反转逻辑：
				// 买入：前一天 RSI < Buy，今天 RSI >= Buy (金叉 Buy Line) -- 这是一个常见的稳健策略
				// 或者简单的：只要 RSI < Buy 就买。
				
				// 让我们使用最直观的：
				// 当 RSI < buyThreshold 时买入
				// 当 RSI > sellThreshold 时卖出
				
				// 增加一个条件：前一天不在区间内，今天进入区间（或者反之），防止信号闪烁？
				// 不，对于全仓模型，只要给出 BUY，如果已有仓位会忽略。
				
				// 策略 A: 
				// if rsi[i] < buyThreshold { return "BUY" }
				// if rsi[i] > sellThreshold { return "SELL" }
				
				// 策略 B (经典反转):
				// 买入: RSI 上穿 BuyThreshold (脱离超卖区) -> rsi[i-1] < buy && rsi[i] >= buy
				// 卖出: RSI 下穿 SellThreshold (脱离超买区) -> rsi[i-1] > sell && rsi[i] <= sell
				
				// 考虑到用户通常理解的“超卖买入”可能是指“进入超卖区就买”或者“离开超卖区才买”。
				// 为了捕捉底部，往往是“掉进坑里”或者“爬出坑”时。
				// 这里我们实现“掉进坑里就买” (RSI < Buy)，这意味着可能会抄在半山腰，但能保证买到。
				// 卖出同理。
				
				if rsi[i] < buyThreshold {
					return "BUY"
				} else if rsi[i] > sellThreshold {
					return "SELL"
				}
			}
			return ""
		},
	)
}


// ============ 指标计算辅助函数 ============

func calculateSMA(data []float64, period int) []float64 {
	res := make([]float64, len(data))
	if period <= 0 {
		return res
	}
	var sum float64
	for i := 0; i < len(data); i++ {
		sum += data[i]
		if i >= period {
			sum -= data[i-period]
		}
		if i >= period-1 {
			res[i] = sum / float64(period)
		}
	}
	return res
}

func calculateEMA(data []float64, period int) []float64 {
	res := make([]float64, len(data))
	if period <= 0 || len(data) == 0 {
		return res
	}
	res[0] = data[0]
	multiplier := 2.0 / (float64(period) + 1.0)
	for i := 1; i < len(data); i++ {
		res[i] = (data[i]-res[i-1])*multiplier + res[i-1]
	}
	return res
}

func calculateMACD(data []float64, fastPeriod, slowPeriod, signalPeriod int) ([]float64, []float64, []float64) {
	fastEMA := calculateEMA(data, fastPeriod)
	slowEMA := calculateEMA(data, slowPeriod)
	
	dif := make([]float64, len(data))
	for i := 0; i < len(data); i++ {
		dif[i] = fastEMA[i] - slowEMA[i]
	}
	
	dea := calculateEMA(dif, signalPeriod)
	
	bar := make([]float64, len(data))
	for i := 0; i < len(data); i++ {
		bar[i] = (dif[i] - dea[i]) * 2
	}
	
	return dif, dea, bar
}

func calculateRSI(data []float64, period int) []float64 {
	res := make([]float64, len(data))
	if period <= 0 || len(data) <= period {
		return res
	}

	// RSI 计算需要计算每日涨跌幅
	gains := make([]float64, len(data))
	losses := make([]float64, len(data))

	for i := 1; i < len(data); i++ {
		change := data[i] - data[i-1]
		if change > 0 {
			gains[i] = change
		} else {
			losses[i] = -change
		}
	}

	// Wilder's Smoothing
	// 第一个平均值是简单的算术平均
	var avgGain, avgLoss float64
	for i := 1; i <= period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// 计算第一个 RSI
	if avgLoss == 0 {
		res[period] = 100
	} else {
		rs := avgGain / avgLoss
		res[period] = 100 - (100 / (1 + rs))
	}

	// 后续计算使用平滑公式
	for i := period + 1; i < len(data); i++ {
		avgGain = (avgGain*float64(period-1) + gains[i]) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + losses[i]) / float64(period)

		if avgLoss == 0 {
			res[i] = 100
		} else {
			rs := avgGain / avgLoss
			res[i] = 100 - (100 / (1 + rs))
		}
	}

	return res
}
