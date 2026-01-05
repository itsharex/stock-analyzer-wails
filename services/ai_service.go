package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"stock-analyzer-wails/models"

	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type AIService struct {
	chatModel model.ChatModel
}

func NewAIService() *AIService {
	ctx := context.Background()
	
	// 从环境变量获取配置
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	modelName := os.Getenv("DASHSCOPE_MODEL")
	if modelName == "" {
		modelName = "qwen-plus" // 默认使用 qwen-plus
	}
	baseURL := os.Getenv("DASHSCOPE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}

	// 初始化 Eino Qwen 适配器
	chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		BaseURL:     baseURL,
		APIKey:      apiKey,
		Model:       modelName,
		Temperature: of(float32(0.7)),
		MaxTokens:   of(2048),
	})
	if err != nil {
		fmt.Printf("初始化 Eino ChatModel 失败: %v\n", err)
		return &AIService{}
	}

	return &AIService{
		chatModel: chatModel,
	}
}

func (s *AIService) AnalyzeStock(stock *models.StockData) (*models.AnalysisReport, error) {
	if s.chatModel == nil {
		return nil, fmt.Errorf("AI服务未正确初始化，请检查 API Key 配置")
	}

	ctx := context.Background()

	// 构建提示词
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

	// 使用 Eino 生成分析
	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(userPrompt),
	}

	resp, err := s.chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("AI分析请求失败: %v", err)
	}

	// 解析 AI 返回的内容
	analysis := resp.Content
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

	if report.Summary == "" {
		report.Summary = analysis
	}

	return report, nil
}

// extractSection 从分析文本中提取指定章节
func (s *AIService) extractSection(text, startMarker, endMarker string) string {
	startIdx := -1
	endIdx := len(text)
	
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
	
	if endMarker != "" && startIdx != -1 {
		for i := startIdx; i < len(text); i++ {
			if i+len(endMarker) <= len(text) && text[i:i+len(endMarker)] == endMarker {
				endIdx = i
				break
			}
		}
	}
	
	if startIdx != -1 && startIdx < endIdx {
		return text[startIdx:endIdx]
	}
	
	return ""
}

func (s *AIService) QuickAnalyze(stock *models.StockData) (string, error) {
	if s.chatModel == nil {
		return "", fmt.Errorf("AI服务未正确初始化")
	}
	ctx := context.Background()
	prompt := fmt.Sprintf(`请用一段话（100字以内）快速分析股票 %s（%s）的投资价值。
当前价格：%.2f元，涨跌幅：%.2f%%，市盈率：%.2f，市净率：%.2f。`,
		stock.Name, stock.Code, stock.Price, stock.ChangeRate, stock.PE, stock.PB)
	
	resp, err := s.chatModel.Generate(ctx, []*schema.Message{
		schema.UserMessage(prompt),
	})
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

func of[T any](t T) *T {
	return &t
}
