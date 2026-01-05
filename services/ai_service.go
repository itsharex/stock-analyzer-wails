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
	
	systemPrompt := `你是一个专业的A股股票分析师。请根据提供的股票实时行情数据，给出一份简短、专业且具有深度的分析报告。
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

// AnalyzeTechnical 深度技术面分析（支持绘图数据输出）
func (s *AIService) AnalyzeTechnical(stock *models.StockData, klines []*models.KLineData) (*models.TechnicalAnalysisResult, error) {
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

	prompt := fmt.Sprintf(`你是一位顶级技术分析师。请对股票 %s (%s) 进行深度形态识别。

最近60个交易日数据(T-0为最新):
%s

当前指标: %s

请输出两部分内容：
1. 【文字分析】：识别经典形态（头肩、双底、三角形等）、量价配合、趋势阶段及操盘建议。
2. 【绘图数据】：请以 JSON 格式输出识别到的关键线段，放在 <DRAWING_JSON> 标签内。
JSON 格式示例：
[
  {"type": "support", "price": 15.5, "label": "强支撑位"},
  {"type": "resistance", "price": 18.2, "label": "近期压力位"},
  {"type": "trendline", "start": "2023-10-01", "startPrice": 14.2, "end": "2023-12-01", "endPrice": 17.5, "label": "上升趋势线"}
]

注意：时间格式必须与数据中的 YYYY-MM-DD 一致。`, stock.Name, stock.Code, strings.Join(klineSummary, "\n"), indicatorInfo)

	ctx := context.Background()
	messages := []*schema.Message{
		schema.SystemMessage("你是一个精通K线绘图的技术分析师。"),
		schema.UserMessage(prompt),
	}

	resp, err := s.chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, err
	}

	content := resp.Content
	
	// 提取 JSON
	drawings := []models.TechnicalDrawing{}
	re := regexp.MustCompile(`(?s)<DRAWING_JSON>(.*?)</DRAWING_JSON>`)
	match := re.FindStringSubmatch(content)
	if len(match) > 1 {
		jsonStr := strings.TrimSpace(match[1])
		json.Unmarshal([]byte(jsonStr), &drawings)
	}

	// 移除 JSON 标签后的纯文字分析
	cleanAnalysis := re.ReplaceAllString(content, "")

	return &models.TechnicalAnalysisResult{
		Analysis: cleanAnalysis,
		Drawings: drawings,
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
