package services

import (
	"context"
	"fmt"
	"os"
	"stock-analyzer-wails/models"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// AIService AI分析服务
type AIService struct {
	client *openai.Client
}

// NewAIService 创建AI服务实例
func NewAIService() *AIService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// 如果环境变量未设置，使用默认值（用户需要自行配置）
		apiKey = "your-api-key-here"
	}
	
	config := openai.DefaultConfig(apiKey)
	
	// 支持自定义API基础URL（用于代理或第三方服务）
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	
	return &AIService{
		client: openai.NewClientWithConfig(config),
	}
}

// AnalyzeStock 分析股票并生成报告
func (s *AIService) AnalyzeStock(stock *models.StockData) (*models.AnalysisReport, error) {
	// 构建分析提示词
	prompt := s.buildAnalysisPrompt(stock)
	
	// 调用OpenAI API
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini, // 使用GPT-4o-mini，性价比高
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你是一位专业的股票分析师，擅长A股市场分析。请基于提供的股票数据，进行专业、客观的分析，并给出合理的投资建议。",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   2000,
		},
	)
	
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("AI未返回分析结果")
	}
	
	// 解析AI响应
	analysis := resp.Choices[0].Message.Content
	
	// 构建分析报告
	report := &models.AnalysisReport{
		StockCode:      stock.Code,
		StockName:      stock.Name,
		Summary:        s.extractSection(analysis, "摘要", "基本面分析"),
		Fundamentals:   s.extractSection(analysis, "基本面分析", "技术面分析"),
		Technical:      s.extractSection(analysis, "技术面分析", "投资建议"),
		Recommendation: s.extractSection(analysis, "投资建议", "风险等级"),
		RiskLevel:      s.extractSection(analysis, "风险等级", "目标价位"),
		TargetPrice:    s.extractSection(analysis, "目标价位", ""),
		GeneratedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}
	
	// 如果提取失败，使用完整分析作为摘要
	if report.Summary == "" {
		report.Summary = analysis
	}
	
	return report, nil
}

// buildAnalysisPrompt 构建分析提示词
func (s *AIService) buildAnalysisPrompt(stock *models.StockData) string {
	return fmt.Sprintf(`请对以下A股股票进行专业分析：

股票信息：
- 股票代码：%s
- 股票名称：%s
- 最新价：%.2f元
- 涨跌幅：%.2f%%
- 涨跌额：%.2f元
- 成交量：%d手
- 成交额：%.2f万元
- 最高价：%.2f元
- 最低价：%.2f元
- 今开：%.2f元
- 昨收：%.2f元
- 振幅：%.2f%%
- 换手率：%.2f%%
- 市盈率：%.2f
- 市净率：%.2f
- 总市值：%.2f亿元
- 流通市值：%.2f亿元

请按以下结构提供分析报告：

## 摘要
（简要概述该股票的当前状况和主要特点）

## 基本面分析
（分析市盈率、市净率、市值等基本面指标，评估公司估值水平）

## 技术面分析
（分析价格走势、成交量、换手率等技术指标，评估短期走势）

## 投资建议
（基于以上分析，给出明确的投资建议：买入/持有/观望/卖出，并说明理由）

## 风险等级
（评估投资风险：低风险/中等风险/高风险）

## 目标价位
（给出合理的目标价位区间）

注意：
1. 分析要客观、专业，基于数据而非主观臆断
2. 投资建议要谨慎，充分提示风险
3. 使用简洁明了的语言，避免过度专业术语
4. 每个部分控制在100-200字左右`,
		stock.Code,
		stock.Name,
		stock.Price,
		stock.ChangeRate,
		stock.Change,
		stock.Volume,
		stock.Amount/10000,
		stock.High,
		stock.Low,
		stock.Open,
		stock.PreClose,
		stock.Amplitude,
		stock.Turnover,
		stock.PE,
		stock.PB,
		stock.TotalMV/100000000,
		stock.CircMV/100000000,
	)
}

// extractSection 从分析文本中提取指定章节
func (s *AIService) extractSection(text, startMarker, endMarker string) string {
	// 简单的文本提取逻辑
	// 实际应用中可以使用更复杂的解析方法
	startIdx := -1
	endIdx := len(text)
	
	// 查找起始标记
	if startMarker != "" {
		for i := 0; i < len(text); i++ {
			if i+len(startMarker) <= len(text) && text[i:i+len(startMarker)] == startMarker {
				startIdx = i + len(startMarker)
				break
			}
		}
	} else {
		startIdx = 0
	}
	
	// 查找结束标记
	if endMarker != "" && startIdx != -1 {
		for i := startIdx; i < len(text); i++ {
			if i+len(endMarker) <= len(text) && text[i:i+len(endMarker)] == endMarker {
				endIdx = i
				break
			}
		}
	}
	
	// 提取文本
	if startIdx != -1 && startIdx < endIdx {
		section := text[startIdx:endIdx]
		// 清理文本（去除多余的空白字符）
		section = trimMultipleSpaces(section)
		return section
	}
	
	return ""
}

// trimMultipleSpaces 清理多余的空白字符
func trimMultipleSpaces(s string) string {
	// 简单实现，实际可以使用正则表达式
	result := ""
	prevSpace := false
	
	for _, c := range s {
		if c == ' ' || c == '\n' || c == '\r' || c == '\t' {
			if !prevSpace {
				result += " "
				prevSpace = true
			}
		} else {
			result += string(c)
			prevSpace = false
		}
	}
	
	return result
}

// QuickAnalyze 快速分析（使用更简单的提示词，响应更快）
func (s *AIService) QuickAnalyze(stock *models.StockData) (string, error) {
	prompt := fmt.Sprintf(`请用一段话（100字以内）快速分析股票 %s（%s）的投资价值。
当前价格：%.2f元，涨跌幅：%.2f%%，市盈率：%.2f，市净率：%.2f。`,
		stock.Name, stock.Code, stock.Price, stock.ChangeRate, stock.PE, stock.PB)
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你是一位专业的股票分析师。",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   200,
		},
	)
	
	if err != nil {
		return "", fmt.Errorf("快速分析失败: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI未返回分析结果")
	}
	
	return resp.Choices[0].Message.Content, nil
}
