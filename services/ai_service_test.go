package services

import (
	"context"
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
	// 1. 准备测试数据
	mockStock := &models.StockData{
		Code:       "600519",
		Name:       "贵州茅台",
		Price:      1700.00,
		ChangeRate: 1.5,
	}

	mockAIResponse := `## 摘要
贵州茅台今日表现强劲。
## 基本面分析
公司盈利能力极强。
## 技术面分析
均线多头排列。
## 投资建议
建议持有。
## 风险等级
低风险。
## 目标价位
1800-1900元。`

	// 2. 初始化带 Mock 的服务
	mockModel := &MockChatModel{
		mockResponse: mockAIResponse,
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
