package services

import (
	"os"
	"path/filepath"
	"time"

	"stock-analyzer-wails/internal/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Provider 类型定义
type Provider string

const (
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
	Provider       Provider              `json:"provider"`
	APIKey         string                `json:"apiKey"`
	BaseURL        string                `json:"baseUrl"`
	Model          string                `json:"model"`
	ProviderModels map[Provider][]string `json:"providerModels"`
}

func LoadAIConfig() (AIResolvedConfig, error) {
	start := time.Now()
	var cfg appYAML

	// 默认值
	cfg.AI.Provider = ProviderDashScope
	cfg.AI.Model = "qwen-plus"
	cfg.AI.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"

	if path, ok, err := findConfigYAMLPath(); err == nil && ok {
		raw, err := os.ReadFile(path)
		if err == nil {
			yaml.Unmarshal(raw, &cfg)
		}
	}

	logger.Info("AI 配置加载完成",
		zap.String("module", "services.config"),
		zap.String("provider", string(cfg.AI.Provider)),
		zap.String("model", cfg.AI.Model),
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
	path, ok, err := findConfigYAMLPath()
	if err != nil {
		return err
	}
	if !ok {
		cwd, _ := os.Getwd()
		path = filepath.Join(cwd, "config.yaml")
	}

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

func findConfigYAMLPath() (path string, ok bool, err error) {
	paths := []string{}
	if exe, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(exe), "config.yaml"))
	}
	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, "config.yaml"))
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, true, nil
		}
	}
	return "", false, nil
}
