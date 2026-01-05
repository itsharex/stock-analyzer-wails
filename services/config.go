package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"stock-analyzer-wails/internal/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type dashscopeYAML struct {
	APIKey  string `yaml:"api_key"`
	Model   string `yaml:"model"`
	BaseURL string `yaml:"base_url"`
}

type appYAML struct {
	Dashscope dashscopeYAML `yaml:"dashscope"`
}

type DashscopeResolvedConfig struct {
	APIKey  string `json:"apiKey"`
	Model   string `json:"model"`
	BaseURL string `json:"baseUrl"`
}

func LoadDashscopeConfig() (DashscopeResolvedConfig, error) {
	const (
		defaultModel = "qwen-plus-2025-07-28"
		defaultBaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	)

	start := time.Now()
	var cfg appYAML
	if path, ok, err := findConfigYAMLPath(); err != nil {
		logger.Error("查找 config.yaml 失败",
			zap.String("module", "services.config"),
			zap.String("op", "LoadDashscopeConfig.findConfigYAMLPath"),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Error(err),
		)
		return DashscopeResolvedConfig{}, err
	} else if ok {
		logger.Info("发现 config.yaml",
			zap.String("module", "services.config"),
			zap.String("op", "LoadDashscopeConfig"),
			zap.String("config_path", path),
		)
		raw, err := os.ReadFile(path)
		if err != nil {
			return DashscopeResolvedConfig{}, fmt.Errorf("读取配置文件失败: %s: %w", path, err)
		}
		if err := yaml.Unmarshal(raw, &cfg); err != nil {
			return DashscopeResolvedConfig{}, fmt.Errorf("解析配置文件失败: %s: %w", path, err)
		}
	}

	apiKey := firstNonEmpty(
		strings.TrimSpace(cfg.Dashscope.APIKey),
		strings.TrimSpace(os.Getenv("DASHSCOPE_API_KEY")),
	)
	model := firstNonEmpty(
		strings.TrimSpace(cfg.Dashscope.Model),
		strings.TrimSpace(os.Getenv("DASHSCOPE_MODEL")),
		defaultModel,
	)
	baseURL := firstNonEmpty(
		strings.TrimSpace(cfg.Dashscope.BaseURL),
		strings.TrimSpace(os.Getenv("DASHSCOPE_BASE_URL")),
		defaultBaseURL,
	)
	if normalized, changed := normalizeDashscopeBaseURL(baseURL); changed {
		baseURL = normalized
	}

	return DashscopeResolvedConfig{
		APIKey:  apiKey,
		Model:   model,
		BaseURL: baseURL,
	}, nil
}

// SaveDashscopeConfig 保存配置到 config.yaml
func SaveDashscopeConfig(config DashscopeResolvedConfig) error {
	path, ok, err := findConfigYAMLPath()
	if err != nil {
		return err
	}
	if !ok {
		// 如果不存在，默认在当前工作目录创建
		cwd, _ := os.Getwd()
		path = filepath.Join(cwd, "config.yaml")
	}

	cfg := appYAML{
		Dashscope: dashscopeYAML{
			APIKey:  config.APIKey,
			Model:   config.Model,
			BaseURL: config.BaseURL,
		},
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func normalizeDashscopeBaseURL(in string) (out string, changed bool) {
	s := strings.TrimSpace(in)
	if s == "" {
		return s, false
	}
	s = strings.TrimRight(s, "/")

	normalized := s
	normalized = strings.ReplaceAll(normalized, "compatible-moe", "compatible-mode")
	normalized = strings.ReplaceAll(normalized, "/dv1", "/v1")

	if strings.HasSuffix(normalized, "/compatible-mode") {
		normalized += "/v1"
	}

	return normalized, normalized != s
}

func findConfigYAMLPath() (path string, ok bool, err error) {
	paths := make([]string, 0, 2)
	exe, err := os.Executable()
	if err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(exe), "config.yaml"))
	}
	cwd, err := os.Getwd()
	if err == nil {
		paths = append(paths, filepath.Join(cwd, "config.yaml"))
	}

	for _, p := range paths {
		_, statErr := os.Stat(p)
		if statErr == nil {
			return p, true, nil
		}
	}
	return "", false, nil
}
