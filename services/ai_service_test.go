package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"stock-analyzer-wails/models"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// MockChatModel 模拟 Eino 的 ChatModel
type MockChatModel struct {
	model.ChatModel
	mockResponse string
	mockError    error
}

func (m *MockChatModel) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return &schema.Message{
		Role:    schema.Assistant,
		Content: m.mockResponse,
	}, nil
}

func TestExtractSectionImpl(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		text        string
		startMarker string
		endMarker   string
		want        string
	}{
		{
			name:        "basic_between_markers",
			text:        "摘要AAA基本面分析BBB技术面分析CCC",
			startMarker: "摘要",
			endMarker:   "基本面分析",
			want:        "AAA",
		},
		{
			name:        "start_missing_returns_empty",
			text:        "摘要AAA基本面分析BBB",
			startMarker: "不存在",
			endMarker:   "基本面分析",
			want:        "",
		},
		{
			name:        "end_missing_returns_to_end",
			text:        "摘要AAA基本面分析BBB",
			startMarker: "摘要",
			endMarker:   "不存在",
			want:        "AAA基本面分析BBB",
		},
		{
			name:        "start_empty_from_zero",
			text:        "AAA基本面分析BBB",
			startMarker: "",
			endMarker:   "基本面分析",
			want:        "AAA",
		},
		{
			name:        "end_empty_to_end",
			text:        "摘要AAA基本面分析BBB",
			startMarker: "摘要",
			endMarker:   "",
			want:        "AAA基本面分析BBB",
		},
		{
			name:        "end_before_start_returns_empty",
			text:        "ENDxxxSTARTyyy",
			startMarker: "START",
			endMarker:   "END",
			want:        "",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := extractSectionImpl(tc.text, tc.startMarker, tc.endMarker)
			if got != tc.want {
				t.Fatalf("extractSectionImpl() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestAnalyzeStock_Mock(t *testing.T) {
	t.Parallel()
	// 1. 准备测试数据
	mockStock := &models.StockData{
		Code:       "600519",
		Name:       "贵州茅台",
		Price:      1700.00,
		ChangeRate: 1.5,
	}

	// 2. 初始化带 Mock 的服务
	mockModel := &MockChatModel{
		mockResponse: `## 摘要
贵州茅台今日表现强劲。
## 基本面分析
略
## 技术面分析
略
## 投资建议
略
## 风险等级
低风险。
## 目标价位
略`,
	}
	service := &AIService{
		chatModel: mockModel,
	}

	// 3. 执行测试
	report, err := service.AnalyzeStock(mockStock)

	// 4. 验证结果
	if err != nil {
		t.Fatalf("AnalyzeStock 失败: %v", err)
	}

	if report.StockCode != "600519" {
		t.Errorf("期望代码 600519, 实际 %s", report.StockCode)
	}

	if report.Summary != "\n贵州茅台今日表现强劲。\n" {
		t.Errorf("摘要提取错误, 实际: %q", report.Summary)
	}

	if report.RiskLevel != "\n低风险。\n" {
		t.Errorf("风险等级提取错误, 实际: %q", report.RiskLevel)
	}
}

func TestAnalyzeEntryStrategy_KlinesInsufficient(t *testing.T) {
	t.Parallel()

	service := &AIService{
		chatModel: &MockChatModel{mockResponse: `{}`},
	}

	_, err := service.AnalyzeEntryStrategy(
		&models.StockData{Code: "600519", Name: "贵州茅台", Price: 1700, ChangeRate: 1.5, Turnover: 1.2},
		[]*models.KLineData{
			{Time: "2026-01-01", Close: 100},
		},
		&models.MoneyFlowResponse{TodayMain: 1000000, Status: "平稳运行"},
		&models.HealthCheckResult{Score: 80, RiskLevel: "中"},
	)

	if err == nil {
		t.Fatalf("期望返回错误，但 err=nil")
	}
	if !strings.Contains(err.Error(), "code=ENTRY_KLINE_INSUFFICIENT") {
		t.Fatalf("期望包含 ENTRY_KLINE_INSUFFICIENT，实际: %v", err)
	}
}

func TestAnalyzeEntryStrategy_InvalidJSON(t *testing.T) {
	t.Parallel()

	service := &AIService{
		chatModel: &MockChatModel{mockResponse: `not a json`},
	}

	_, err := service.AnalyzeEntryStrategy(
		&models.StockData{Code: "600519", Name: "贵州茅台", Price: 1700, ChangeRate: 1.5, Turnover: 1.2},
		[]*models.KLineData{
			{Time: "2026-01-01", Close: 100},
			{Time: "2026-01-02", Close: 101},
		},
		&models.MoneyFlowResponse{TodayMain: 1000000, Status: "平稳运行"},
		&models.HealthCheckResult{Score: 80, RiskLevel: "中"},
	)

	if err == nil {
		t.Fatalf("期望返回错误，但 err=nil")
	}
	if !strings.Contains(err.Error(), "code=ENTRY_AI_INVALID_JSON") {
		t.Fatalf("期望包含 ENTRY_AI_INVALID_JSON，实际: %v", err)
	}
}

func TestAnalyzeEntryStrategy_EntryPriceRangeArray(t *testing.T) {
	t.Parallel()

	service := &AIService{
		chatModel: &MockChatModel{mockResponse: `{
  "recommendation": "分批建仓",
  "entryPriceRange": [21.50, 22.20],
  "initialPosition": "20%",
  "stopLossPrice": 20.10,
  "takeProfitPrice": 25.00,
  "coreReasons": [{"type":"technical","description":"x","threshold":"y"}],
  "riskRewardRatio": 1.78,
  "actionPlan": "test"
}`},
	}

	got, err := service.AnalyzeEntryStrategy(
		&models.StockData{Code: "002202", Name: "金风科技", Price: 22.01, ChangeRate: 1.5, Turnover: 1.2},
		[]*models.KLineData{
			{Time: "2026-01-01", Close: 100},
			{Time: "2026-01-02", Close: 101},
		},
		&models.MoneyFlowResponse{TodayMain: 0, Status: "平稳运行"},
		&models.HealthCheckResult{Score: 100, RiskLevel: "低"},
	)
	if err != nil {
		t.Fatalf("期望成功，但失败: %v", err)
	}
	if got.EntryPriceRange == "" {
		t.Fatalf("期望 entryPriceRange 非空")
	}
	if !strings.Contains(got.EntryPriceRange, "21.50") {
		t.Fatalf("期望 entryPriceRange 被格式化包含 21.50，实际: %q", got.EntryPriceRange)
	}
}

func TestAnalyzeEntryStrategy_Timeout(t *testing.T) {
	t.Parallel()

	service := &AIService{
		chatModel: &MockChatModel{mockError: context.DeadlineExceeded},
	}

	_, err := service.AnalyzeEntryStrategy(
		&models.StockData{Code: "600519", Name: "贵州茅台", Price: 1700, ChangeRate: 1.5, Turnover: 1.2},
		[]*models.KLineData{
			{Time: "2026-01-01", Close: 100},
			{Time: "2026-01-02", Close: 101},
		},
		&models.MoneyFlowResponse{TodayMain: 1000000, Status: "平稳运行"},
		&models.HealthCheckResult{Score: 80, RiskLevel: "中"},
	)

	if err == nil {
		t.Fatalf("期望返回错误，但 err=nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) && !strings.Contains(err.Error(), "code=ENTRY_AI_TIMEOUT") {
		t.Fatalf("期望包含 ENTRY_AI_TIMEOUT（或可识别为超时），实际: %v", err)
	}
}


func TestAnalyzeStock(t *testing.T) {
	cfg, err := LoadAIConfig()
	if err != nil {
		t.Skipf("跳过需要真实配置的测试：LoadAIConfig 失败: %v", err)
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		t.Skip("跳过需要真实 APIKey 的测试：APIKey 为空")
	}
	aiService, err := NewAIService(cfg)
	if err != nil {
		t.Skipf("跳过需要真实配置的测试：NewAIService 失败: %v", err)
	}
	report, err := aiService.AnalyzeStock(&models.StockData{
		Code: "600519",
		Name: "贵州茅台",
		Price: 1700.00,
		ChangeRate: 1.5,
	})
	if err != nil {
		t.Fatalf("AnalyzeStock 失败: %v", err)
	}
	if report.StockCode != "600519" {
		t.Errorf("期望代码 600519, 实际 %s", report.StockCode)
	}
	if strings.TrimSpace(report.Summary) == "" {
		t.Errorf("期望摘要非空, 实际: %q", report.Summary)
	}
}