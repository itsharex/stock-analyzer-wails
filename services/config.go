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
	APIKey  string
	Model   string
	BaseURL string
}

func LoadDashscopeConfig() (DashscopeResolvedConfig, error) {
	const (
		defaultModel   = "qwen-plus-2025-07-28"
		defaultBaseURL = "https://dashscope.aliyuncs.com/compatible-moe/dv1"
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
			logger.Error("读取 config.yaml 失败",
				zap.String("module", "services.config"),
				zap.String("op", "LoadDashscopeConfig.os.ReadFile"),
				zap.String("config_path", path),
				zap.Int64("duration_ms", time.Since(start).Milliseconds()),
				zap.Error(err),
			)
			return DashscopeResolvedConfig{}, fmt.Errorf("读取配置文件失败: %s: %w", path, err)
		}
		if err := yaml.Unmarshal(raw, &cfg); err != nil {
			logger.Error("解析 config.yaml 失败",
				zap.String("module", "services.config"),
				zap.String("op", "LoadDashscopeConfig.yaml.Unmarshal"),
				zap.String("config_path", path),
				zap.Int("body_size", len(raw)),
				zap.Int64("duration_ms", time.Since(start).Milliseconds()),
				zap.Error(err),
			)
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

	if apiKey == "" {
		logger.Error("DashScope API Key 缺失",
			zap.String("module", "services.config"),
			zap.String("op", "LoadDashscopeConfig.validate"),
			zap.Bool("has_yaml_api_key", strings.TrimSpace(cfg.Dashscope.APIKey) != ""),
			zap.Bool("has_env_api_key", strings.TrimSpace(os.Getenv("DASHSCOPE_API_KEY")) != ""),
		)
		return DashscopeResolvedConfig{}, errors.New("DashScope API Key 缺失：请在 config.yaml 的 dashscope.api_key 或环境变量 DASHSCOPE_API_KEY 中配置")
	}

	logger.Info("DashScope 配置加载完成",
		zap.String("module", "services.config"),
		zap.String("op", "LoadDashscopeConfig"),
		zap.String("model", model),
		zap.String("base_url", baseURL),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return DashscopeResolvedConfig{
		APIKey:  apiKey,
		Model:   model,
		BaseURL: baseURL,
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func findConfigYAMLPath() (path string, ok bool, err error) {
	paths := make([]string, 0, 2)

	// 1) 可执行文件目录
	exe, err := os.Executable()
	if err != nil {
		return "", false, fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	paths = append(paths, filepath.Join(filepath.Dir(exe), "config.yaml"))

	// 2) 当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		return "", false, fmt.Errorf("获取当前工作目录失败: %w", err)
	}
	paths = append(paths, filepath.Join(cwd, "config.yaml"))

	for _, p := range paths {
		_, statErr := os.Stat(p)
		if statErr == nil {
			return p, true, nil
		}
		if errors.Is(statErr, os.ErrNotExist) {
			continue
		}
		return "", false, fmt.Errorf("检查配置文件失败: %s: %w", p, statErr)
	}

	return "", false, nil
}


