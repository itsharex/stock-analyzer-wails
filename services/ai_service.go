package services

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"stock-analyzer-wails/models"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type AIService struct {
	chatModel model.ChatModel
	config    AIResolvedConfig
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

	return &AIService{
		chatModel: cm,
		config:    cfg,
	}, nil
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
		return nil, fmt.Errorf("AI分析请求失败: %w", err)
	}

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
func (s *AIService) AnalyzeTechnical(stock *models.StockData, klines []*models.KLineData, role string) (*models.TechnicalAnalysisResult, error) {
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
	if startIdx < 0 { startIdx = 0 }
	
	var maxPrice, minPrice float64
	for i := startIdx; i < len(klines); i++ {
		k := klines[i]
		if i == startIdx || k.High > maxPrice { maxPrice = k.High }
		if i == startIdx || k.Low < minPrice { minPrice = k.Low }
		
		if i > len(klines)-15 || i%3 == 0 {
			change := 0.0
			if k.Open != 0 {
				change = (k.Close - k.Open) / k.Open * 100
			}
			klineSummary = append(klineSummary, fmt.Sprintf("T-%d(%s): O:%.2f, C:%.2f, H:%.2f, L:%.2f, Vol:%d, Chg:%.2f%%", len(klines)-1-i, k.Time, k.Close, k.High, k.Low, k.Volume, change))
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

		prompt := fmt.Sprintf(`%s 请对股票 %s (%s) 进行深度多维度评估。
			%s
			你的受众包含大量股票新手，请在提到专业术语时，使用括号附带通俗易懂的解释。
		
		最近60个交易日数据(T-0为最新):
		%s
		
		当前指标: %s
		
		请输出五部分内容，**必须严格遵守以下标签格式，不要在标签内包含任何 Markdown 代码块标记（如 \`\`\`json）**：
		1. 【文字分析】：识别经典形态、量价配合、趋势阶段及操盘建议。
		2. 【风险评估】：请以纯 JSON 格式输出风险得分和操盘建议，放在 <RISK_JSON> 标签内。
		示例：<RISK_JSON>{"riskScore": 65, "actionAdvice": "观望"}</RISK_JSON>
		3. 【绘图数据】：请以纯 JSON 格式输出识别到的关键线段，放在 <DRAWING_JSON> 标签内。
		4. 【多维度评分】：请以纯 JSON 格式输出五个维度的评分（0-100）及理由，放在 <RADAR_JSON> 标签内。
		5. 【智能交易计划】：请以纯 JSON 格式输出具体的交易建议，放在 <TRADE_JSON> 标签内。
		包括：建议仓位(suggestedPosition, 如"30%%")、止损价(stopLoss)、止盈价(takeProfit)、盈亏比(riskRewardRatio)、操作策略(strategy)。
		
		**重要：即使你正在扮演特定角色，也请确保 JSON 标签内的内容是纯净的 JSON 字符串，以便程序解析。**`, selectedRole.System, stock.Name, stock.Code, selectedRole.Style, strings.Join(klineSummary, "\n"), indicatorInfo)
	
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

	// 提取绘图 JSON
	drawings := []models.TechnicalDrawing{}
	reDrawing := regexp.MustCompile(`(?s)<DRAWING_JSON>(.*?)</DRAWING_JSON>`)
	matchDrawing := reDrawing.FindStringSubmatch(content)
	if len(matchDrawing) > 1 {
		jsonStr := cleanJSON(matchDrawing[1])
		json.Unmarshal([]byte(jsonStr), &drawings)
	}

	// 提取风险 JSON
	riskData := struct {
		RiskScore    int    `json:"riskScore"`
		ActionAdvice string `json:"actionAdvice"`
	}{RiskScore: 50, ActionAdvice: "观望"}
	reRisk := regexp.MustCompile(`(?s)<RISK_JSON>(.*?)</RISK_JSON>`)
	matchRisk := reRisk.FindStringSubmatch(content)
	if len(matchRisk) > 1 {
		jsonStr := cleanJSON(matchRisk[1])
		json.Unmarshal([]byte(jsonStr), &riskData)
	}

	// 提取雷达图 JSON
	radarData := &models.RadarData{}
	reRadar := regexp.MustCompile(`(?s)<RADAR_JSON>(.*?)</RADAR_JSON>`)
	matchRadar := reRadar.FindStringSubmatch(content)
	if len(matchRadar) > 1 {
		jsonStr := cleanJSON(matchRadar[1])
		json.Unmarshal([]byte(jsonStr), radarData)
	}

	// 提取交易计划 JSON
	tradePlan := &models.TradePlan{}
	reTrade := regexp.MustCompile(`(?s)<TRADE_JSON>(.*?)</TRADE_JSON>`)
	matchTrade := reTrade.FindStringSubmatch(content)
	if len(matchTrade) > 1 {
		jsonStr := cleanJSON(matchTrade[1])
		json.Unmarshal([]byte(jsonStr), tradePlan)
	}

	// 移除 JSON 标签后的纯文字分析
	cleanAnalysis := reDrawing.ReplaceAllString(content, "")
	cleanAnalysis = reRisk.ReplaceAllString(cleanAnalysis, "")
	cleanAnalysis = reRadar.ReplaceAllString(cleanAnalysis, "")
	cleanAnalysis = reTrade.ReplaceAllString(cleanAnalysis, "")

	return &models.TechnicalAnalysisResult{
		Analysis:     cleanAnalysis,
		Drawings:     drawings,
		RiskScore:    riskData.RiskScore,
		ActionAdvice: riskData.ActionAdvice,
		RadarData:    radarData,
		TradePlan:    tradePlan,
	}, nil
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
			return text[startIdx:endRel]
		}
		return ""
	}

	return text[startIdx:]
}
