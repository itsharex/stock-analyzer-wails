package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	Level string

	// FilePath 为空时，默认写入：<exeDir>/logs/app.log
	FilePath string

	MaxSizeMB    int
	MaxBackups   int
	MaxAgeDays   int
	Compress     bool
	Development  bool
	EnableStdout bool
	EnableFile   bool
}

var (
	mu sync.RWMutex
	l  = zap.NewNop()
)

// InitFromEnv 初始化全局 logger。
//
// 支持环境变量：
// - LOG_LEVEL: debug|info|warn|error
// - LOG_FILE_PATH: 自定义日志文件路径（默认 <exeDir>/logs/app.log）
// - LOG_MAX_SIZE_MB / LOG_MAX_BACKUPS / LOG_MAX_AGE_DAYS / LOG_COMPRESS
func InitFromEnv() error {
	opts := Options{
		Level:        getenvDefault("LOG_LEVEL", "info"),
		FilePath:     strings.TrimSpace(os.Getenv("LOG_FILE_PATH")),
		MaxSizeMB:    getenvIntDefault("LOG_MAX_SIZE_MB", 50),
		MaxBackups:   getenvIntDefault("LOG_MAX_BACKUPS", 10),
		MaxAgeDays:   getenvIntDefault("LOG_MAX_AGE_DAYS", 14),
		Compress:     getenvBoolDefault("LOG_COMPRESS", true),
		Development:  strings.EqualFold(getenvDefault("APP_ENV", "prod"), "dev"),
		EnableStdout: true,
		EnableFile:   true,
	}
	return Init(opts)
}

func Init(opts Options) error {
	level, err := parseLevel(opts.Level)
	if err != nil {
		return err
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	enc := zapcore.NewJSONEncoder(encCfg)

	cores := make([]zapcore.Core, 0, 2)

	if opts.EnableStdout {
		cores = append(cores, zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), level))
	}

	if opts.EnableFile {
		path, err := resolveLogFilePath(opts.FilePath)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("创建日志目录失败: %s: %w", filepath.Dir(path), err)
		}

		lj := &lumberjack.Logger{
			Filename:   path,
			MaxSize:    opts.MaxSizeMB,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAgeDays,
			Compress:   opts.Compress,
		}
		cores = append(cores, zapcore.NewCore(enc, zapcore.AddSync(lj), level))
	}

	core := zapcore.NewTee(cores...)

	newLogger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	mu.Lock()
	l = newLogger
	mu.Unlock()

	return nil
}

func L() *zap.Logger {
	mu.RLock()
	defer mu.RUnlock()
	return l
}

func Sync() {
	_ = L().Sync()
}

func Debug(msg string, fields ...zap.Field) { L().Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { L().Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { L().Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { L().Error(msg, fields...) }

func parseLevel(s string) (zapcore.Level, error) {
	var lvl zapcore.Level
	if err := lvl.Set(strings.ToLower(strings.TrimSpace(s))); err != nil {
		return zapcore.InfoLevel, fmt.Errorf("无效 LOG_LEVEL=%q（应为 debug|info|warn|error）", s)
	}
	return lvl, nil
}

func resolveLogFilePath(p string) (string, error) {
	if strings.TrimSpace(p) != "" {
		return p, nil
	}
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	return filepath.Join(filepath.Dir(exe), "logs", "app.log"), nil
}

func getenvDefault(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func getenvIntDefault(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func getenvBoolDefault(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return def
	}
}


