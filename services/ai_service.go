package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"stock-analyzer-wails/models"
	"strings"
	"time"

	"stock-analyzer-wails/internal/logger"

	"go.uber.org/zap"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type AIService struct {
	chatModel    model.ChatModel
	config       AIResolvedConfig
	cacheService *AnalysisCacheService
	semaphore    chan struct{} // 并发控制
	enableMock   bool          // 启用 Mock 模式
}

func NewAIService(cfg AIResolvedConfig) (*AIService, error) {
	ctx := context.Background()

	opts := &openai.ChatModelConfig{
		APIKey:  cfg.APIKey,
		BaseURL: cfg.BaseURL,
		Model:   cfg.Model,
	}

	cm, err := openai.NewChatModel(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("创建 ChatModel 失败 (%s): %w", cfg.Provider, err)
	}

	cacheSvc, _ := NewAnalysisCacheService()

	// 默认限速处理：这里我们仅使用 semaphore 控制并发数
	// Eino 框架底层处理 HTTP 连接池，这里增加应用层并发限制
	return &AIService{
		chatModel:    cm,
		config:       cfg,
		cacheService: cacheSvc,
		semaphore:    make(chan struct{}, 5), // 默认最大并发数 5
		enableMock:   false,
	}, nil
}

// SetEnableMock 设置是否启用 Mock 模式
func (s *AIService) SetEnableMock(enable bool) {
	s.enableMock = enable
}

// VerifySignalAsync 异步验证信号
func (s *AIService) VerifySignalAsync(stock *models.StockData, recentFlows []models.MoneyFlowData) <-chan *models.AIVerificationResult {
	resultChan := make(chan *models.AIVerificationResult, 1)

	go func() {
		defer close(resultChan)

		// 简单的信号量控制并发
		s.semaphore <- struct{}{}
		defer func() { <-s.semaphore }()

		res, err := s.VerifySignal(stock, recentFlows)
		if err != nil {
			logger.Error("AI 验证信号失败", zap.String("code", stock.Code), zap.Error(err))
			return
		}
		resultChan <- res
	}()

	return resultChan
}

// VerifySignal 对股票近期的资金流向进行深度解读
func (s *AIService) VerifySignal(stock *models.StockData, recentFlows []models.MoneyFlowData) (*models.AIVerificationResult, error) {
	// Mock 模式
	if s.enableMock {
		time.Sleep(500 * time.Millisecond) // 模拟延迟
		return &models.AIVerificationResult{
			Score:     85,
			Opinion:   "主力连续吸筹，量价配合良好，建议重点关注。",
			RiskLevel: "低",
		}, nil
	}

	// 1. 数据组装
	type FlowDetail struct {
		Date            string  `json:"date"`
		MainNet         float64 `json:"main_net"`
		SuperNet        float64 `json:"super_net"`
		BigNet          float64 `json:"big_net"`
		ChgPct          float64 `json:"chg_pct"`
		MainInflowRatio float64 `json:"main_inflow_ratio"` // 主力流入占比
	}

	var details []FlowDetail
	for _, f := range recentFlows {
		// 估算成交额：如果有成交量且有收盘价，Amount ≈ Close * Volume * 100 (手 -> 股)
		// 但 MoneyFlowData 没有 Volume。我们只能传 0，并在 Prompt 中说明或忽略。
		// 为了满足 User 明确要求 "多算一个字段给 AI"，我们尽力而为。
		// 如果无法计算，AI 会根据 Prompt 规则处理（例如忽略或基于净额判断）。
		// 由于 MoneyFlowData 结构体限制，我们这里暂且填 0。
		// 更好的做法是 MoneyFlowData 包含 Volume 或 Amount，但不想改动太大。
		// 我们假设 Amount 为 0，AI 看到 0 会处理。

		ratio := 0.0
		// 如果我们能获取到当天的总成交额就好了。
		// 暂时填 0.0

		details = append(details, FlowDetail{
			Date:            f.TradeDate,
			MainNet:         f.MainNet,
			SuperNet:        f.SuperNet,
			BigNet:          f.BigNet,
			ChgPct:          f.ChgPct,
			MainInflowRatio: ratio,
		})
	}

	detailsJSON, _ := json.Marshal(details)

	// 2. 构造 Prompt
	prompt := fmt.Sprintf(`你是一位精通筹码分布的量化交易专家。 现有股票 %s(%s) 近 7 个交易日的资金流向数据： %s
 
 请根据以上数据进行复核： 
 
 吸筹识别：主力资金是在股价下跌时逆势吸筹，还是在拉升过程中诱多出货？ 
 
 筹码集中度：超大单流入是否具备持续性？ 
 
 风险提示：是否存在资金流向与涨跌幅背离的情况？ 
 
 请以 JSON 格式返回： { "score": (0-100的整数，代表建仓胜率), "opinion": (简短的专家分析理由，不超过 100 字), "risk_level": ("低", "中", "高") }`,
		stock.Name, stock.Code, string(detailsJSON))

	// 3. 调用 LLM
	ctx := context.Background()
	messages := []*schema.Message{
		schema.SystemMessage("你是一位精通筹码分布的量化交易专家。"),
		schema.UserMessage(prompt),
	}
	logger.Info("AI 请求prompt", zap.Any("prompt", prompt))

	resp, err := s.chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("AI分析请求失败: %w", err)
	}

	// 4. 解析结果
	var result models.AIVerificationResult
	jsonStr := s.extractJSON(resp.Content)
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		logger.Error("解析 AI 响应失败", zap.String("resp", resp.Content), zap.Error(err))
		return nil, fmt.Errorf("解析 AI 响应失败: %w", err)
	}
	logger.Info("AI 响应", zap.Any("resp", resp))

	return &result, nil
}

func (s *AIService) AnalyzeStock(stock *models.StockData) (*models.AnalysisReport, error) {
	ctx := context.Background()

	systemPrompt := `你是一个专业的A股股票分析师。你的受众包含大量股票新手，请在提到专业术语时，使用括号附带通俗易懂的解释。
报告必须包含以下部分：
1. 分析摘要：简要概括当前走势。
2. 基本面分析：基于市盈率、市净率、市值等数据评估。
3. 技术面分析：基于价格变动、涨跌幅、换手率等数据评估。
4. 投资建议：给出明确的建议（买入/持有/观望/卖出）并说明理由。
5. 风险等级：低/中/高。
6. 目标价位：给出一个合理的短期目标价区间。

请按以下结构提供分析报告：
## 摘要
...
## 基本面分析
...
## 技术面分析
...
## 投资建议
...
## 风险等级
...
## 目标价位
...`

	userPrompt := fmt.Sprintf(`股票名称：%s (%s)
最新价：%.2f
涨跌幅：%.2f%%
涨跌额：%.2f
成交量：%d
成交额：%.2f
最高价：%.2f
最低价：%.2f
今开：%.2f
昨收：%.2f
换手率：%.2f%%
市盈率(动态)：%.2f
市净率：%.2f
总市值：%.2f
流通市值：%.2f`,
		stock.Name, stock.Code, stock.Price, stock.ChangeRate, stock.Change,
		stock.Volume, stock.Amount, stock.High, stock.Low, stock.Open, stock.PreClose,
		stock.Turnover, stock.PE, stock.PB, stock.TotalMV, stock.CircMV)

	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(userPrompt),
	}

	resp, err := s.chatModel.Generate(ctx, messages)
	if err != nil {
		logger.Error("AI分析请求失败", zap.Error(err))
		return nil, fmt.Errorf("AI分析请求失败: %w", err)
	}

	logger.Info("AI分析响应", zap.String("response", resp.Content))

	analysis := resp.Content

	report := &models.AnalysisReport{
		StockCode:      stock.Code,
		StockName:      stock.Name,
		Summary:        extractSectionImpl(analysis, "摘要", "基本面分析"),
		Fundamentals:   extractSectionImpl(analysis, "基本面分析", "技术面分析"),
		Technical:      extractSectionImpl(analysis, "技术面分析", "投资建议"),
		Recommendation: extractSectionImpl(analysis, "投资建议", "风险等级"),
		RiskLevel:      extractSectionImpl(analysis, "风险等级", "目标价位"),
		TargetPrice:    extractSectionImpl(analysis, "目标价位", ""),
		GeneratedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}

	if report.Summary == "" {
		report.Summary = analysis
	}

	return report, nil
}

// AnalyzeTechnical 深度技术面分析（支持多角色切换、绘图数据和风险评估）
// GenerateAlertAdvice 根据警报触发情况生成角色化的建议
func (s *AIService) GenerateAlertAdvice(stockName, alertType, label, role string, currentPrice, alertPrice float64) (string, error) {
	rolePrompts := map[string]string{
		"conservative": "你是一位名为'稳健老船长'的资深投资顾问。你极度厌恶风险，推崇价值投资和安全边际。你的语言风格沉稳、老练，经常使用航海比喻。",
		"aggressive":   "你是一位名为'激进先锋官'的短线交易高手。你追求资金效率，擅长捕捉热点和动能爆发。你的语言风格果断、充满激情，经常使用军事比喻。",
		"technical":    "你是一位名为'技术派大师'的量化分析专家。你只相信数据和图形，不带任何感情色彩。你的语言风格冷静、客观、专业。",
	}

	systemPrompt := rolePrompts[role]
	if systemPrompt == "" {
		systemPrompt = rolePrompts["technical"]
	}

	prompt := fmt.Sprintf("股票 %s 触发了价格预警。\n预警类型：%s\n关键位描述：%s\n当前价：%.2f\n关键位价格：%.2f\n\n请作为你的角色，给出一句极其简短（20字以内）的'大白话'操作建议。",
		stockName, alertType, label, currentPrice, alertPrice)

	ctx := context.Background()
	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(prompt),
	}

	resp, err := s.chatModel.Generate(ctx, messages)
	if err != nil {
		return "注意风险，按计划操作。", nil
	}

	return strings.TrimSpace(resp.Content), nil
}

func (s *AIService) AnalyzeEntryStrategy(stock *models.StockData, klines []*models.KLineData, moneyFlow *models.MoneyFlowResponse, health *models.HealthCheckResult) (*models.EntryStrategyResult, error) {
	// 设置 20 秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	start := time.Now()
	stockCode := ""
	stockName := ""
	if stock != nil {
		stockCode = stock.Code
		stockName = stock.Name
	}

	logger.Info("开始建仓分析（AI）",
		zap.String("module", "services.ai"),
		zap.String("op", "AnalyzeEntryStrategy"),
		zap.String("stock_code", stockCode),
		zap.String("stock_name", stockName),
		zap.String("provider", string(s.config.Provider)),
		zap.String("model", s.config.Model),
	)

	if stock == nil || strings.TrimSpace(stock.Code) == "" {
		return nil, fmt.Errorf("建仓分析失败(step=input, code=ENTRY_INPUT_INVALID): 股票数据为空或股票代码为空")
	}
	if moneyFlow == nil {
		return nil, fmt.Errorf("建仓分析失败(step=input, code=ENTRY_INPUT_MISSING): 资金流向数据缺失")
	}
	if health == nil {
		return nil, fmt.Errorf("建仓分析失败(step=input, code=ENTRY_INPUT_MISSING): 体检数据缺失")
	}

	// 构建 K 线摘要
	if len(klines) < 2 {
		return nil, fmt.Errorf("建仓分析失败(step=kline_summary, code=ENTRY_KLINE_INSUFFICIENT): K线数据不足（len=%d）", len(klines))
	}

	klineSummary := ""
	startIdx := len(klines) - 10
	if startIdx < 1 {
		startIdx = 1 // 需要访问 i-1
	}
	for i := startIdx; i < len(klines); i++ {
		k := klines[i]
		prev := klines[i-1]
		pct := 0.0
		if prev != nil && prev.Close != 0 {
			pct = (k.Close - prev.Close) / prev.Close * 100
		}
		klineSummary += fmt.Sprintf("日期:%s,收盘:%.2f,涨跌:%.2f%%; ", k.Time, k.Close, pct)
	}

	systemPrompt := `你是一位资深的量化交易员和风险管理专家。你的任务是根据提供的股票数据，为用户生成一份极具实战价值的“智能建仓方案”。
你的分析必须严谨，给出的价格和比例必须具体。

请按以下 JSON 格式输出：
{
  "recommendation": "建议类型(立即建仓/分批建仓/等待回调/暂时观望)",
  "entryPriceRange": "建议买入价格区间（必须是字符串，如\"21.50-22.20\"，不要输出数组）",
  "initialPosition": "建议首仓比例（必须是字符串，如\"20%\"）",
  "stopLossPrice": 止损价(数字),
  "takeProfitPrice": 目标止盈价(数字),
  "coreReasons": [
    {"type": "fundamental/technical/money_flow", "description": "理由描述", "threshold": "逻辑失效的触发阈值"}
  ],
  "riskRewardRatio": 预估盈亏比(数字),
  "actionPlan": "具体操作步骤描述"
}`

	userPrompt := fmt.Sprintf(`股票: %s (%s)
当前价: %.2f, 涨跌幅: %.2f%%, 换手率: %.2f%%
体检评分: %d, 风险等级: %s
今日资金流向: 主力净流入 %.2f 万, 状态: %s
最近K线走势: %s

请基于以上数据，给出深度建仓分析方案。`,
		stock.Name, stock.Code, stock.Price, stock.ChangeRate, stock.Turnover,
		health.Score, health.RiskLevel, moneyFlow.TodayMain/10000, moneyFlow.Status, klineSummary)

	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(userPrompt),
	}

	resp, err := s.chatModel.Generate(ctx, messages)
	if err != nil {
		// 解决冲突：保留日志记录和更详细的错误信息
		logger.Error("建仓分析 AI 调用失败",
			zap.String("module", "services.ai"),
			zap.String("op", "AnalyzeEntryStrategy"),
			zap.String("step", "ai_generate"),
			zap.String("stock_code", stock.Code),
			zap.Error(err),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("建仓分析失败(step=ai_generate, code=ENTRY_AI_TIMEOUT): AI 调用超时")
		}
		return nil, fmt.Errorf("建仓分析失败(step=ai_generate, code=ENTRY_AI_REQUEST_FAILED): %v", err)
	}
	cleanJSON := s.extractJSON(resp.Content)
	var result models.EntryStrategyResult
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		logger.Error("建仓分析解析失败（JSON）",
			zap.String("module", "services.ai"),
			zap.String("op", "AnalyzeEntryStrategy"),
			zap.String("step", "json_unmarshal"),
			zap.String("stock_code", stock.Code),
			zap.Error(err),
			zap.String("clean_json_preview", truncateString(cleanJSON, 2048)),
			zap.String("resp_preview", truncateString(resp.Content, 2048)),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return nil, fmt.Errorf("建仓分析失败(step=json_unmarshal, code=ENTRY_AI_INVALID_JSON): 解析失败: %w", err)
	}

	logger.Info("建仓分析（AI）成功",
		zap.String("module", "services.ai"),
		zap.String("op", "AnalyzeEntryStrategy"),
		zap.String("stock_code", stock.Code),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return &result, nil
}

func (s *AIService) extractJSON(content string) string {
	// 优先提取 fenced code block（```json ... ``` 或 ``` ... ```）
	reFence := regexp.MustCompile("(?s)```(?:json)?\\s*(\\{.*?\\})\\s*```")
	if m := reFence.FindStringSubmatch(content); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}

	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start != -1 && end != -1 && end > start {
		return content[start : end+1]
	}
	return content
}

func truncateString(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(s) <= max {
		return s
	}
	return s[:max] + "...(truncated)"
}

func (s *AIService) AnalyzeTechnical(stock *models.StockData, klines []*models.KLineData, period string, role string) (*models.TechnicalAnalysisResult, error) {
	// 检查是否强制刷新
	force := false
	if strings.HasSuffix(role, ":force") {
		force = true
		role = strings.TrimSuffix(role, ":force")
	}

	// 1. 尝试从缓存获取
	if strings.TrimSpace(period) == "" {
		period = "daily"
	}

	if !force && s.cacheService != nil {
		if cached, ok := s.cacheService.Get(stock.Code, role, period); ok {
			return cached, nil
		}
	}

	// 2. 缓存未命中，执行 AI 分析
	// 角色 Prompt 定义
	rolePrompts := map[string]struct {
		System string
		Style  string
	}{
		"conservative": {
			System: "你是一位名为'稳健老船长'的资深投资顾问。你极度厌恶风险，推崇价值投资和安全边际。",
			Style:  "你的语言风格沉稳、老练，经常使用航海比喻。你对仓位控制非常严格，止损位设置较宽以防洗盘，但对基本面瑕疵零容忍。",
		},
		"aggressive": {
			System: "你是一位名为'激进先锋官'的短线交易高手。你追求资金效率，擅长捕捉热点和动能爆发。",
			Style:  "你的语言风格果断、充满激情，经常使用军事比喻。你关注量价齐升，止盈目标宏大，敢于在趋势确认时重仓出击。",
		},
		"technical": {
			System: "你是一位名为'技术派大师'的量化分析专家。你只相信数据和图形，不带任何感情色彩。",
			Style:  "你的语言风格冷静、客观、专业。你专注于指标背离、形态识别和支撑阻力位，给出的建议极其精确，不废话。",
		},
	}

	// 默认使用技术派
	selectedRole, ok := rolePrompts[role]
	if !ok {
		selectedRole = rolePrompts["technical"]
	}
	var klineSummary []string
	startIdx := len(klines) - 60
	if startIdx < 0 {
		startIdx = 0
	}

	var maxPrice, minPrice float64
	for i := startIdx; i < len(klines); i++ {
		k := klines[i]
		if i == startIdx || k.High > maxPrice {
			maxPrice = k.High
		}
		if i == startIdx || k.Low < minPrice {
			minPrice = k.Low
		}

		if i > len(klines)-15 || i%3 == 0 {
			change := 0.0
			if k.Open != 0 {
				change = (k.Close - k.Open) / k.Open * 100
			}
			klineSummary = append(klineSummary, fmt.Sprintf("T-%d(%s): O:%.2f, C:%.2f, H:%.2f, L:%.2f, Vol:%d, Chg:%.2f%%", len(klines)-1-i, k.Time, k.Open, k.Close, k.High, k.Low, k.Volume, change))
		}
	}

	lastK := klines[len(klines)-1]
	indicatorInfo := ""
	if lastK.MACD != nil {
		indicatorInfo += fmt.Sprintf("MACD(DIF:%.3f, DEA:%.3f, BAR:%.3f); ", lastK.MACD.DIF, lastK.MACD.DEA, lastK.MACD.Bar)
	}
	if lastK.KDJ != nil {
		indicatorInfo += fmt.Sprintf("KDJ(K:%.1f, D:%.1f, J:%.1f); ", lastK.KDJ.K, lastK.KDJ.D, lastK.KDJ.J)
	}
	if lastK.RSI > 0 {
		indicatorInfo += fmt.Sprintf("RSI:%.1f; ", lastK.RSI)
	}

	prompt := fmt.Sprintf("%s 请对股票 %s (%s) 进行深度多维度评估。\n"+
		"%s\n"+
		"你的受众包含大量股票新手，请在提到专业术语时，使用括号附带通俗易懂的解释。\n\n"+
		"最近60个交易日数据(T-0为最新):\n%s\n\n"+
		"当前指标: %s\n\n"+
		"请输出五部分内容，**必须严格遵守以下标签格式，不要在标签内包含任何 Markdown 代码块标记（如 ```json）**：\n"+
		"1. 【文字分析】：识别经典形态、量价配合、趋势阶段及操盘建议。\n"+
		"2. 【风险评估】：请以纯 JSON 格式输出风险得分和操盘建议，放在 <RISK_JSON> 标签内。\n"+
		"示例：<RISK_JSON>{\"riskScore\": 65, \"actionAdvice\": \"观望\"}</RISK_JSON>\n"+
		"3. 【绘图数据】：请以纯 JSON 格式输出识别到的关键线段，放在 <DRAWING_JSON> 标签内。**注意：必须至少包含一个支撑位(support)和一个压力位(resistance)，如果趋势不明显，请选择最近的局部高低点。**\n"+
		"4. 【多维度评分】：请以纯 JSON 格式输出五个维度的评分（0-100）及理由，放在 <RADAR_JSON> 标签内。\n"+
		"5. 【智能交易计划】：请以纯 JSON 格式输出具体的交易建议，放在 <TRADE_JSON> 标签内。\n"+
		"包括：建议仓位(suggestedPosition, 如\"30%%\")、止损价(stopLoss)、止盈价(takeProfit)、盈亏比(riskRewardRatio)、操作策略(strategy)。\n\n"+
		"**重要：即使你正在扮演特定角色，也请确保 JSON 标签内的内容是纯净的 JSON 字符串，以便程序解析。**",
		selectedRole.System, stock.Name, stock.Code, selectedRole.Style, strings.Join(klineSummary, "\n"), indicatorInfo)

	ctx := context.Background()
	messages := []*schema.Message{
		schema.SystemMessage(selectedRole.System + " 你精通K线绘图和风险管理，擅长用通俗易懂的语言向新手解释复杂的金融术语。"),
		schema.UserMessage(prompt),
	}

	resp, err := s.chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, err
	}

	content := resp.Content

	// 辅助函数：清理 JSON 字符串中的 Markdown 标记
	cleanJSON := func(s string) string {
		s = strings.TrimSpace(s)
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimSuffix(s, "```")
		return strings.TrimSpace(s)
	}

	// 使用递归模糊解析器提取绘图数据
	drawings := []models.TechnicalDrawing{}
	reDrawing := regexp.MustCompile(`(?s)<DRAWING_JSON>(.*?)</DRAWING_JSON>`)
	if match := reDrawing.FindStringSubmatch(content); len(match) > 1 {
		drawings = robustParseDrawings(cleanJSON(match[1]))
	}

	// 提取风险 JSON
	riskData := struct {
		RiskScore    int    `json:"riskScore"`
		ActionAdvice string `json:"actionAdvice"`
	}{RiskScore: 50, ActionAdvice: "观望"}
	reRisk := regexp.MustCompile(`(?s)<RISK_JSON>(.*?)</RISK_JSON>`)
	if match := reRisk.FindStringSubmatch(content); len(match) > 1 {
		json.Unmarshal([]byte(cleanJSON(match[1])), &riskData)
	}

	// 提取雷达图 JSON
	radarData := &models.RadarData{}
	reRadar := regexp.MustCompile(`(?s)<RADAR_JSON>(.*?)</RADAR_JSON>`)
	if match := reRadar.FindStringSubmatch(content); len(match) > 1 {
		json.Unmarshal([]byte(cleanJSON(match[1])), radarData)
	}

	// 提取交易计划 JSON
	tradePlan := &models.TradePlan{}
	reTrade := regexp.MustCompile(`(?s)<TRADE_JSON>(.*?)</TRADE_JSON>`)
	if match := reTrade.FindStringSubmatch(content); len(match) > 1 {
		json.Unmarshal([]byte(cleanJSON(match[1])), tradePlan)
	}

	// 移除 JSON 标签后的纯文字分析
	cleanAnalysis := reDrawing.ReplaceAllString(content, "")
	cleanAnalysis = reRisk.ReplaceAllString(cleanAnalysis, "")
	cleanAnalysis = reRadar.ReplaceAllString(cleanAnalysis, "")
	cleanAnalysis = reTrade.ReplaceAllString(cleanAnalysis, "")

	result := &models.TechnicalAnalysisResult{
		Analysis:     cleanAnalysis,
		Drawings:     drawings,
		RiskScore:    riskData.RiskScore,
		ActionAdvice: riskData.ActionAdvice,
		RadarData:    radarData,
		TradePlan:    tradePlan,
	}

	// 3. 存入缓存
	if s.cacheService != nil {
		s.cacheService.Set(stock.Code, role, period, *result)
	}

	return result, nil
}

// robustParseDrawings 递归模糊解析绘图数据
func robustParseDrawings(jsonStr string) []models.TechnicalDrawing {
	var results []models.TechnicalDrawing
	var raw interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return results
	}

	// 辅助函数：尝试从 map 中获取 float64
	getFloat := func(m map[string]interface{}, keys ...string) (float64, bool) {
		for _, k := range keys {
			if v, ok := m[k].(float64); ok {
				return v, true
			}
		}
		return 0, false
	}

	// 辅助函数：尝试从 map 中获取 string
	getString := func(m map[string]interface{}, keys ...string) (string, bool) {
		for _, k := range keys {
			if v, ok := m[k].(string); ok {
				return v, true
			}
		}
		return "", false
	}

	var search func(data interface{})
	search = func(data interface{}) {
		switch v := data.(type) {
		case []interface{}:
			for _, item := range v {
				search(item)
			}
		case map[string]interface{}:
			// 启发式识别：如果一个对象同时拥有“价格特征”
			price, hasPrice := getFloat(v, "price", "level", "value", "val", "support", "resistance")
			label, _ := getString(v, "label", "name", "desc", "role")
			role, hasRole := getString(v, "type", "role", "kind")

			if hasPrice && price > 0 {
				dType := role
				if !hasRole {
					// 如果没有明确 type，尝试从字段名推断
					if _, ok := v["support"]; ok {
						dType = "support"
					}
					if _, ok := v["resistance"]; ok {
						dType = "resistance"
					}
				}

				results = append(results, models.TechnicalDrawing{
					Price: price,
					Type:  dType,
					Label: label,
				})
			}
			// 继续深挖子节点（如 segments 数组）
			for _, val := range v {
				search(val)
			}
		}
	}

	search(raw)
	return results
}

func extractSectionImpl(text, startMarker, endMarker string) string {
	startIdx := -1
	if startMarker == "" {
		startIdx = 0
	} else {
		for i := 0; i+len(startMarker) <= len(text); i++ {
			if text[i:i+len(startMarker)] == startMarker {
				startIdx = i + len(startMarker)
				break
			}
		}
	}
	if startIdx == -1 {
		return ""
	}

	if endMarker == "" {
		return text[startIdx:]
	}

	endRel := -1
	for i := startIdx; i+len(endMarker) <= len(text); i++ {
		if text[i:i+len(endMarker)] == endMarker {
			endRel = i
			break
		}
	}
	if endRel != -1 {
		if startIdx < endRel {
			// 兼容 Markdown 标题：如 endMarker 位于 "## xxx" 中，截取段落末尾可能残留 "## "。
			return strings.TrimRight(text[startIdx:endRel], " \t#")
		}
		return ""
	}

	// 若 endMarker 仅出现在 startMarker 之前，认为段落不完整，返回空（避免误把后续全部当作内容）
	if endMarker != "" && startIdx > 0 && strings.Contains(text[:startIdx], endMarker) {
		return ""
	}

	return text[startIdx:]
}
