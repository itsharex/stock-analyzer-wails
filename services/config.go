package services

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/repositories"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const AppName = "StockAnalyzer"

var (
	appDataDirOnce  sync.Once
	cachedAppDataDir string
)

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

// MigrateAIConfigFromYAML 负责将旧的 config.yaml 迁移到 SQLite
func (s *ConfigService) MigrateAIConfigFromYAML() error {
	path := filepath.Join(GetAppDataDir(), "config.yaml")

	// 检查文件是否存在
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil // 文件不存在，无需迁移
	}
	if err != nil {
		return fmt.Errorf("检查 config.yaml 状态失败: %w", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取 config.yaml 失败: %w", err)
	}

	var cfg appYAML
	if len(raw) > 0 {
		if err := yaml.Unmarshal(raw, &cfg); err != nil {
			return fmt.Errorf("解析 config.yaml 失败: %w", err)
		}
	}

	// 转换为新的配置结构
	newConfig := AIResolvedConfig{
		Provider: cfg.AI.Provider,
		APIKey:   cfg.AI.APIKey,
		BaseURL:  cfg.AI.BaseURL,
		Model:    cfg.AI.Model,
	}

	// 保存到 SQLite
	if err := s.SaveAIConfig(newConfig); err != nil {
		return fmt.Errorf("保存 AI 配置到 SQLite 失败: %w", err)
	}

	// 迁移成功，重命名旧文件
	backupPath := path + ".bak." + time.Now().Format("20060102150405")
	if err := os.Rename(path, backupPath); err != nil {
		logger.Error("重命名旧 config.yaml 文件失败", zap.Error(err))
		// 即使重命名失败，也认为迁移成功，只是下次启动会再次尝试迁移
	} else {
		logger.Info("成功将 config.yaml 迁移到 SQLite", zap.String("backup_path", backupPath))
	}

	return nil
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
	appDataDirOnce.Do(func() {
		var (
			baseDir string
			source  string
		)

		switch runtime.GOOS {
		case "windows":
			baseDir = strings.TrimSpace(os.Getenv("LOCALAPPDATA"))
			source = "LOCALAPPDATA"
			if baseDir == "" {
				baseDir = strings.TrimSpace(os.Getenv("APPDATA"))
				source = "APPDATA"
			}
		case "darwin":
			home, err := os.UserHomeDir()
			if err != nil {
				logger.Warn("获取用户目录失败，将回退到当前目录",
					zap.String("module", "services.config"),
					zap.String("op", "GetAppDataDir"),
					zap.String("os", runtime.GOOS),
					zap.Error(err),
				)
				home = "."
			}
			baseDir = filepath.Join(home, "Library", "Application Support")
			source = "UserHomeDir"
		default: // Linux and others
			home, err := os.UserHomeDir()
			if err != nil {
				logger.Warn("获取用户目录失败，将回退到当前目录",
					zap.String("module", "services.config"),
					zap.String("op", "GetAppDataDir"),
					zap.String("os", runtime.GOOS),
					zap.Error(err),
				)
				home = "."
			}
			baseDir = filepath.Join(home, ".local", "share")
			source = "UserHomeDir"
		}

		if strings.TrimSpace(baseDir) == "" {
			// 极端情况下环境变量为空，保证至少有个可写路径
			baseDir = "."
			source = "fallback-dot"
			logger.Warn("应用数据目录基路径为空，已回退到当前目录（可能导致权限问题）",
				zap.String("module", "services.config"),
				zap.String("op", "GetAppDataDir"),
				zap.String("os", runtime.GOOS),
			)
		}

		appDir := filepath.Join(baseDir, AppName)
		if err := os.MkdirAll(appDir, 0o755); err != nil {
			logger.Error("创建应用数据目录失败（SQLite 初始化可能失败）",
				zap.String("module", "services.config"),
				zap.String("op", "GetAppDataDir"),
				zap.String("os", runtime.GOOS),
				zap.String("source", source),
				zap.String("baseDir", baseDir),
				zap.String("appDir", appDir),
				zap.Error(err),
			)
		} else {
			logger.Info("应用数据目录已解析",
				zap.String("module", "services.config"),
				zap.String("op", "GetAppDataDir"),
				zap.String("os", runtime.GOOS),
				zap.String("source", source),
				zap.String("baseDir", baseDir),
				zap.String("appDir", appDir),
			)
		}

		cachedAppDataDir = appDir
	})

	return cachedAppDataDir
}

// ConfigService 负责全局配置的读取和写入
type ConfigService struct {
	repo repositories.ConfigRepository
}

// NewConfigService 构造函数
func NewConfigService(repo repositories.ConfigRepository) *ConfigService {
	return &ConfigService{repo: repo}
}

// getConfigValue 从数据库中读取配置值
func (s *ConfigService) getConfigValue(key string) (string, error) {
	return s.repo.GetConfigValue(key)
}

// setConfigValue 向数据库中写入配置值
func (s *ConfigService) setConfigValue(key string, value string) error {
	return s.repo.SetConfigValue(key, value)
}

// LoadAIConfig 从 SQLite 加载 AI 配置
func (s *ConfigService) LoadAIConfig() (AIResolvedConfig, error) {
	// 确保在加载前执行迁移
	if err := s.MigrateAIConfigFromYAML(); err != nil {
		logger.Error("执行 YAML 配置迁移失败", zap.Error(err))
		// 迁移失败不影响后续加载，继续执行
	}
	start := time.Now()

	// 默认值
	cfg := AIResolvedConfig{
		Provider:       ProviderQwen,
		Model:          "qwen-plus",
		BaseURL:        "https://dashscope.aliyuncs.com/compatible-mode/v1",
		ProviderModels: ProviderModels,
	}

	// 从数据库读取配置
	providerStr, _ := s.getConfigValue("ai_provider")
	apiKey, _ := s.getConfigValue("ai_api_key")
	baseURL, _ := s.getConfigValue("ai_base_url")
	model, _ := s.getConfigValue("ai_model")

	if providerStr != "" {
		cfg.Provider = Provider(providerStr)
	}
	if apiKey != "" {
		cfg.APIKey = apiKey
	}
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	if model != "" {
		cfg.Model = model
	}

	// 兼容性修复：规范化 BaseURL
	if normalized, changed := normalizeDashscopeBaseURL(cfg.BaseURL); changed {
		logger.Warn("检测到 BaseURL 需要修正，已自动规范化",
			zap.String("module", "services.config"),
			zap.String("op", "normalize_base_url"),
			zap.String("before", cfg.BaseURL),
			zap.String("after", normalized),
		)
		cfg.BaseURL = normalized
	}

	logger.Info("AI 配置加载完成",
		zap.String("module", "services.config"),
		zap.String("provider", string(cfg.Provider)),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)

	return cfg, nil
}

// SaveAIConfig 保存 AI 配置到 SQLite
func (s *ConfigService) SaveAIConfig(config AIResolvedConfig) error {
	if err := s.setConfigValue("ai_provider", string(config.Provider)); err != nil {
		return err
	}
	if err := s.setConfigValue("ai_api_key", config.APIKey); err != nil {
		return err
	}

	// 规范化 BaseURL 后保存
	normalizedBaseURL, _ := normalizeDashscopeBaseURL(config.BaseURL)
	if err := s.setConfigValue("ai_base_url", normalizedBaseURL); err != nil {
		return err
	}
	if err := s.setConfigValue("ai_model", config.Model); err != nil {
		return err
	}
	return nil
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

// LoadAIConfig is a backward-compatible helper used by older code/tests.
// It builds a minimal ConfigService using the SQLite-backed repository and returns the resolved config.
func LoadAIConfig() (AIResolvedConfig, error) {
	dbSvc, err := NewDBService()
	if err != nil {
		return AIResolvedConfig{}, err
	}
	defer dbSvc.Close()

	repo := repositories.NewSQLiteConfigRepository(dbSvc.GetDB())
	svc := NewConfigService(repo)
	return svc.LoadAIConfig()
}
