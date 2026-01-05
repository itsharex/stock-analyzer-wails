package main

import (
	"embed"
	"os"

	"stock-analyzer-wails/internal/logger"

	"go.uber.org/zap"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := logger.InitFromEnv(); err != nil {
		// logger 尚未初始化完成，只能兜底输出到 stderr
		_, _ = os.Stderr.WriteString("初始化 logger 失败: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer logger.Sync()

	// 创建应用实例
	app := NewApp()

	// 创建应用配置
	err := wails.Run(&options.App{
		Title:  "A股股票分析AI-Agent",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		// Windows特定配置
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			Theme:                windows.SystemDefault,
		},
		// macOS特定配置
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameAqua,
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "A股股票分析AI-Agent",
				Message: "基于AI的专业股票分析工具\n\n版本: 1.0.0\n\n© 2026 Stock Analyzer",
			},
		},
	})

	if err != nil {
		logger.Error("启动应用失败",
			zap.String("module", "main"),
			zap.String("op", "wails.Run"),
			zap.Error(err),
		)
		os.Exit(1)
	}
}
