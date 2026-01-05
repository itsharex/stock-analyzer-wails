package services

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"stock-analyzer-wails/internal/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const AppName = "StockAnalyzer"

// Provider 类型定义
type Provider string

const (
	ProviderQwen      Provider = "Qwen"
	ProviderDashScope Provider = "DashScope"
	ProviderDeepSeek  Provider = "DeepSeek"
	ProviderOpenAI    Provider = "OpenAI"
	ProviderClaude    Provider = "Claude"
	ProviderGemini    Provider = "Gemini"
	ProviderARK       Provider = "ARK"
	ProviderQianfan   Provider = "Qianfan"
)

// 供应商及其默认模型
var ProviderModels = map[Provider][]string{
	ProviderQwen:      {"qwen-plus", "qwen-max", "qwen-turbo", "qwen-long"},
	ProviderDashScope: {"qwen-plus", "qwen-max", "qwen-turbo", "qwen-long"},
	ProviderDeepSeek:  {"deepseek-chat", "deepseek-reasoner"},
	ProviderOpenAI:    {"gpt-4o", "gpt-4o-mini", "gpt-4-turbo"},
	ProviderClaude:    {"claude-3-5-sonnet-20240620", "claude-3-opus-20240229"},
	ProviderGemini:    {"gemini-1.5-pro", "gemini-1.5-flash"},
	ProviderARK:       {"doubao-pro-4k", "doubao-lite-4k"},
	ProviderQianfan:   {"ernie-4.0-8k", "ernie-3.5-8k"},
}

type AIConfigYAML struct {
	Provider Provider `yaml:"provider"`
	APIKey   string   `yaml:"api_key"`
	BaseURL  string   `yaml:"base_url"`
	Model    string   `yaml:"model"`
}

type appYAML struct {
	AI AIConfigYAML `yaml:"ai"`
}

type AIResolvedConfig struct {
	Provider       Provider            `json:"provider"`
	APIKey         string              `json:"apiKey"`
	BaseURL        string              `json:"baseUrl"`
	Model          string              `json:"model"`
	ProviderModels map[Provider][]string `json:"providerModels"`
}

// GetAppDataDir 获取跨平台的应用数据目录
func GetAppDataDir() string {
	var dir string
	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("LOCALAPPDATA")
		if dir == "" {
			dir = os.Getenv("APPDATA")
		}
	case "darwin":
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, "Library", "Application Support")
	default: // Linux and others
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, ".local", "share")
	}

	appDir := filepath.Join(dir, AppName)
	// 确保目录存在
	os.MkdirAll(appDir, 0755)
	return appDir
}

func LoadAIConfig() (AIResolvedConfig, error) {
	start := time.Now()
	var cfg appYAML
	
	// 默认值
	cfg.AI.Provider = ProviderQwen
	cfg.AI.Model = "qwen-plus"
	cfg.AI.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"

	path := filepath.Join(GetAppDataDir(), "config.yaml")
	
	// 尝试从新位置读取
	raw, err := os.ReadFile(path)
	if err != nil {
		// 如果新位置没有，尝试从当前目录读取（兼容旧版本）
		if oldRaw, oldErr := os.ReadFile("config.yaml"); oldErr == nil {
			raw = oldRaw
			// 迁移到新位置
			os.WriteFile(path, oldRaw, 0644)
		}
	}

	if len(raw) > 0 {
		yaml.Unmarshal(raw, &cfg)
	}

	logger.Info("AI 配置加载完成",
		zap.String("module", "services.config"),
		zap.String("path", path),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)

	return AIResolvedConfig{
		Provider:       cfg.AI.Provider,
		APIKey:         cfg.AI.APIKey,
		BaseURL:        cfg.AI.BaseURL,
		Model:          cfg.AI.Model,
		ProviderModels: ProviderModels,
	}, nil
}

func SaveAIConfig(config AIResolvedConfig) error {
	path := filepath.Join(GetAppDataDir(), "config.yaml")

	cfg := appYAML{
		AI: AIConfigYAML{
			Provider: config.Provider,
			APIKey:   config.APIKey,
			BaseURL:  config.BaseURL,
			Model:    config.Model,
		},
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
