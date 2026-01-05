package services

import (
	"context"
	"fmt"
	"stock-analyzer-wails/models"
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
	
	// 大多数供应商（DeepSeek, DashScope, OpenAI, Claude via Proxy）都支持 OpenAI 兼容协议
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
