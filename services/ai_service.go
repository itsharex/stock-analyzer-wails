package services

import (
	"context"
	"fmt"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"

	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

type AIService struct {
	chatModel model.ChatModel
}

func NewAIService(cfg DashscopeResolvedConfig) (*AIService, error) {
	ctx := context.Background()
	start := time.Now()

	// 初始化 Eino Qwen 适配器
	chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		BaseURL:     cfg.BaseURL,
		APIKey:      cfg.APIKey,
		Model:       cfg.Model,
		Temperature: of(float32(0.7)),
		MaxTokens:   of(2048),
	})
	if err != nil {
		logger.Error("初始化 Eino ChatModel 失败",
			zap.String("module", "services.ai"),
			zap.String("op", "NewAIService.qwen.NewChatModel"),
			zap.String("model", cfg.Model),
			zap.String("base_url", cfg.BaseURL),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("初始化 Eino ChatModel 失败: %w", err)
	}

	logger.Info("AI 服务初始化成功",
		zap.String("module", "services.ai"),
		zap.String("op", "NewAIService"),
		zap.String("model", cfg.Model),
		zap.String("base_url", cfg.BaseURL),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)

	return &AIService{
		chatModel: chatModel,
	}, nil
}

func (s *AIService) AnalyzeStock(stock *models.StockData) (*models.AnalysisReport, error) {
	start := time.Now()
	if s.chatModel == nil {
		logger.Error("AI服务未正确初始化",
			zap.String("module", "services.ai"),
			zap.String("op", "AnalyzeStock"),
			zap.String("stock_code", stock.Code),
			zap.String("stock_name", stock.Name),
		)
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
		logger.Error("AI分析请求失败",
			zap.String("module", "services.ai"),
			zap.String("op", "AnalyzeStock.Generate"),
			zap.String("stock_code", stock.Code),
			zap.String("stock_name", stock.Name),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("AI分析请求失败: %w", err)
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

	logger.Info("AI分析完成",
		zap.String("module", "services.ai"),
		zap.String("op", "AnalyzeStock"),
		zap.String("stock_code", stock.Code),
		zap.String("stock_name", stock.Name),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return report, nil
}

// extractSection 从分析文本中提取指定章节
func (s *AIService) extractSection(text, startMarker, endMarker string) string {
	return extractSectionImpl(text, startMarker, endMarker)
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

	// endMarker 为空则截取到末尾
	if endMarker == "" {
		return text[startIdx:]
	}

	// 查找 startIdx 之后的 endMarker
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

	// endMarker 若仅出现在 startMarker 之前，则认为顺序不合法
	for i := 0; i+len(endMarker) <= startIdx; i++ {
		if text[i:i+len(endMarker)] == endMarker {
			return ""
		}
	}

	// endMarker 不存在：截取到末尾
	return text[startIdx:]
}

func (s *AIService) QuickAnalyze(stock *models.StockData) (string, error) {
	start := time.Now()
	if s.chatModel == nil {
		logger.Error("AI服务未正确初始化",
			zap.String("module", "services.ai"),
			zap.String("op", "QuickAnalyze"),
			zap.String("stock_code", stock.Code),
			zap.String("stock_name", stock.Name),
		)
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
		logger.Error("快速分析请求失败",
			zap.String("module", "services.ai"),
			zap.String("op", "QuickAnalyze.Generate"),
			zap.String("stock_code", stock.Code),
			zap.String("stock_name", stock.Name),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return "", err
	}
	logger.Info("快速分析完成",
		zap.String("module", "services.ai"),
		zap.String("op", "QuickAnalyze"),
		zap.String("stock_code", stock.Code),
		zap.String("stock_name", stock.Name),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return resp.Content, nil
}

func of[T any](t T) *T {
	return &t
}
