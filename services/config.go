package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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

// AIConfigYAML 保持不变，用于兼容旧的 YAML 配置
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

// GlobalStrategyConfig 全局策略配置
type GlobalStrategyConfig struct {
	TrailingStopActivation float64 `json:"trailingStopActivation"` // 移动止损启动阈值 (0.05 = 5%)
	TrailingStopCallback   float64 `json:"trailingStopCallback"`   // 移动止损回撤比例 (0.03 = 3%)
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

// ConfigService 负责全局配置的读取和写入
type ConfigService struct {
	db *sql.DB
}

// NewConfigService 构造函数
func NewConfigService(dbSvc *DBService) *ConfigService {
	return &ConfigService{db: dbSvc.GetDB()}
}

// getConfigValue 从数据库中读取配置值
func (s *ConfigService) getConfigValue(key string) (string, error) {
	var value string
	err := s.db.QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // 配置项不存在
		}
		return "", fmt.Errorf("查询配置项 %s 失败: %w", key, err)
	}
	return value, nil
}

// setConfigValue 向数据库中写入配置值
func (s *ConfigService) setConfigValue(key string, value string) error {
	query := `
		INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)
	`
	_, err := s.db.Exec(query, key, value)
	if err != nil {
		return fmt.Errorf("保存配置项 %s 失败: %w", key, err)
	}
	return nil
}

// LoadAIConfig 从 YAML 文件加载 AI 配置 (保持兼容性)
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

	// 兼容性修复：规范化 DashScope BaseURL
	if normalized, changed := normalizeDashscopeBaseURL(cfg.AI.BaseURL); changed {
		logger.Warn("检测到 DashScope BaseURL 需要修正，已自动规范化",
			zap.String("module", "services.config"),
			zap.String("op", "normalize_base_url"),
			zap.String("before", cfg.AI.BaseURL),
			zap.String("after", normalized),
		)
		cfg.AI.BaseURL = normalized
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

// SaveAIConfig 保存 AI 配置到 YAML 文件 (保持兼容性)
func SaveAIConfig(config AIResolvedConfig) error {
	path := filepath.Join(GetAppDataDir(), "config.yaml")

	cfg := appYAML{
		AI: AIConfigYAML{
			Provider: config.Provider,
			APIKey:   config.APIKey,
			BaseURL:  func() string { s, _ := normalizeDashscopeBaseURL(config.BaseURL); return s }(),
			Model:    config.Model,
		},
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetGlobalStrategyConfig 从 SQLite 获取全局策略配置
func (s *ConfigService) GetGlobalStrategyConfig() (GlobalStrategyConfig, error) {
	var config GlobalStrategyConfig
	
	// 默认值 (与 db_service.go 中保持一致)
	config.TrailingStopActivation = 0.05
	config.TrailingStopCallback = 0.03

	// 读取启动阈值
	activationStr, err := s.getConfigValue("trailing_stop_default_activation")
	if err != nil {
		return config, err
	}
	if activationStr != "" {
		if val, err := strconv.ParseFloat(activationStr, 64); err == nil {
			config.TrailingStopActivation = val
		}
	}

	// 读取回撤比例
	callbackStr, err := s.getConfigValue("trailing_stop_default_callback")
	if err != nil {
		return config, err
	}
	if callbackStr != "" {
		if val, err := strconv.ParseFloat(callbackStr, 64); err == nil {
			config.TrailingStopCallback = val
		}
	}

	return config, nil
}

// UpdateGlobalStrategyConfig 更新全局策略配置到 SQLite
func (s *ConfigService) UpdateGlobalStrategyConfig(config GlobalStrategyConfig) error {
	if err := s.setConfigValue("trailing_stop_default_activation", fmt.Sprintf("%f", config.TrailingStopActivation)); err != nil {
		return err
	}
	if err := s.setConfigValue("trailing_stop_default_callback", fmt.Sprintf("%f", config.TrailingStopCallback)); err != nil {
		return err
	}
	return nil
}

func normalizeDashscopeBaseURL(in string) (string, bool) {
	orig := in
	s := strings.TrimSpace(in)
	if s == "" {
		return s, strings.TrimSpace(orig) != ""
	}

	// 移除末尾 /
	s = strings.TrimRight(s, "/")

	// 修复常见 typo
	s = strings.ReplaceAll(s, "/compatible-moe/", "/compatible-mode/")
	s = strings.ReplaceAll(s, "/compatible-moe", "/compatible-mode")
	s = strings.ReplaceAll(s, "/dv1", "/v1")

	// 仅 compatible-mode 未带版本时，补全 /v1
	if strings.HasSuffix(s, "/compatible-mode") {
		s = s + "/v1"
	}

	return s, s != strings.TrimSpace(orig)
}
