# 部署与分发指南

本文档介绍如何构建、打包和分发A股股票分析AI-Agent应用。

## 构建生产版本

### 基本构建

```bash
# 构建当前平台的可执行文件
wails build

# 构建后的文件位于 build/bin/ 目录
```

### 平台特定构建

#### Windows

```bash
# 64位Windows
wails build -platform windows/amd64

# 输出文件：build/bin/stock-analyzer.exe
```

**Windows安装包制作**：
可以使用以下工具创建安装程序：
- NSIS (Nullsoft Scriptable Install System)
- Inno Setup
- WiX Toolset

#### macOS

```bash
# Universal Binary (支持Intel和Apple Silicon)
wails build -platform darwin/universal

# 仅Intel Mac
wails build -platform darwin/amd64

# 仅Apple Silicon Mac
wails build -platform darwin/arm64

# 输出文件：build/bin/stock-analyzer.app
```

**macOS应用签名和公证**：

1. 代码签名：
```bash
codesign --deep --force --verify --verbose --sign "Developer ID Application: Your Name" build/bin/stock-analyzer.app
```

2. 创建DMG安装包：
```bash
# 使用create-dmg工具
npm install -g create-dmg
create-dmg build/bin/stock-analyzer.app
```

3. 公证（Notarization）：
```bash
xcrun notarytool submit stock-analyzer.dmg --apple-id your@email.com --team-id TEAMID --password app-specific-password
```

#### Linux

```bash
# 64位Linux
wails build -platform linux/amd64

# ARM64 Linux
wails build -platform linux/arm64

# 输出文件：build/bin/stock-analyzer
```

**Linux打包选项**：
- AppImage
- Snap
- Flatpak
- DEB包
- RPM包

### 构建优化选项

```bash
# 压缩构建（减小文件大小）
wails build -clean -upx

# 跳过前端构建（如已手动构建）
wails build -skipbindings

# 详细输出
wails build -v 2

# 自定义输出目录
wails build -o custom-output-dir
```

## 应用图标

### 准备图标文件

应用图标应放置在以下位置：

```
build/
├── appicon.png          # 通用图标（512x512或1024x1024）
├── windows/
│   └── icon.ico         # Windows图标
├── darwin/
│   └── icon.icns        # macOS图标
└── linux/
    └── icon.png         # Linux图标
```

### 图标格式要求

- **Windows (.ico)**：包含多个尺寸（16x16, 32x32, 48x48, 256x256）
- **macOS (.icns)**：包含多个尺寸（16x16 到 1024x1024）
- **Linux (.png)**：推荐512x512或更高分辨率

### 图标生成工具

```bash
# 使用在线工具
# https://www.icoconverter.com/
# https://iconverticons.com/

# 或使用命令行工具
# ImageMagick (跨平台)
convert icon.png -resize 256x256 icon.ico

# iconutil (macOS)
iconutil -c icns icon.iconset
```

## 应用元数据配置

编辑 `wails.json` 文件：

```json
{
  "name": "stock-analyzer-wails",
  "outputfilename": "stock-analyzer",
  "author": {
    "name": "Your Name",
    "email": "your@email.com"
  },
  "info": {
    "companyName": "Your Company",
    "productName": "A股股票分析AI-Agent",
    "productVersion": "1.0.0",
    "copyright": "Copyright © 2026 Your Company",
    "comments": "基于AI的A股股票分析工具"
  }
}
```

## 依赖打包

### Windows

Wails会自动嵌入WebView2 Loader，但用户系统需要安装WebView2运行时。

**WebView2分发选项**：
1. 在线安装器（推荐）：用户首次运行时自动下载
2. 离线安装器：将WebView2安装包打包到安装程序中
3. Fixed Version：将特定版本的WebView2打包到应用中（增加约100MB）

### macOS

无需额外依赖，系统自带WebKit。

### Linux

需要确保用户系统安装了以下依赖：
- GTK3
- WebKit2GTK

可以在安装脚本中添加依赖检查：
```bash
#!/bin/bash
if ! dpkg -l | grep -q libwebkit2gtk-4.0; then
    echo "正在安装依赖..."
    sudo apt-get install -y libwebkit2gtk-4.0-37
fi
```

## 自动更新

### 实现自动更新功能

可以集成以下更新方案：

1. **Electron-updater风格**（推荐）
   - 使用GitHub Releases托管更新
   - 实现版本检查和下载逻辑

2. **自定义更新服务器**
   - 搭建更新API服务
   - 提供版本检查和下载接口

示例更新检查代码（Go）：

```go
type UpdateInfo struct {
    Version     string `json:"version"`
    DownloadURL string `json:"download_url"`
    Changelog   string `json:"changelog"`
}

func CheckForUpdates() (*UpdateInfo, error) {
    resp, err := http.Get("https://your-server.com/api/updates/latest")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var info UpdateInfo
    if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
        return nil, err
    }
    
    return &info, nil
}
```

## 分发渠道

### 1. GitHub Releases

```bash
# 创建Release
gh release create v1.0.0 \
  build/bin/stock-analyzer-windows-amd64.exe \
  build/bin/stock-analyzer-darwin-universal.app.zip \
  build/bin/stock-analyzer-linux-amd64 \
  --title "v1.0.0 - 首次发布" \
  --notes "初始版本发布"
```

### 2. 应用商店

#### Microsoft Store (Windows)

1. 注册Microsoft开发者账号
2. 使用Windows Application Packaging Project打包
3. 提交到Microsoft Store

#### Mac App Store

1. 注册Apple Developer账号
2. 配置App ID和证书
3. 使用Xcode打包并提交

#### Snap Store (Linux)

```bash
# 创建snapcraft.yaml
snapcraft

# 发布到Snap Store
snapcraft upload stock-analyzer_1.0.0_amd64.snap
```

### 3. 自建下载站点

创建简单的下载页面：

```html
<!DOCTYPE html>
<html>
<head>
    <title>下载 A股股票分析AI-Agent</title>
</head>
<body>
    <h1>下载应用</h1>
    <ul>
        <li><a href="/downloads/stock-analyzer-windows.exe">Windows版本</a></li>
        <li><a href="/downloads/stock-analyzer-macos.dmg">macOS版本</a></li>
        <li><a href="/downloads/stock-analyzer-linux.AppImage">Linux版本</a></li>
    </ul>
</body>
</html>
```

## 版本管理

### 语义化版本

遵循 [Semantic Versioning](https://semver.org/)：

- **主版本号**：不兼容的API修改
- **次版本号**：向下兼容的功能性新增
- **修订号**：向下兼容的问题修正

示例：`1.2.3`
- 1：主版本
- 2：次版本
- 3：修订版本

### 版本发布流程

1. 更新版本号（`wails.json`）
2. 更新CHANGELOG.md
3. 提交代码并打标签
4. 构建所有平台版本
5. 创建GitHub Release
6. 发布到各分发渠道
7. 通知用户更新

## 许可证

在项目根目录添加 `LICENSE` 文件，常见选择：

- **MIT License**：宽松，允许商业使用
- **Apache 2.0**：宽松，提供专利授权
- **GPL v3**：Copyleft，要求开源

## 安全建议

### 1. 保护API密钥

- ❌ 不要将API密钥硬编码在代码中
- ✅ 使用环境变量或配置文件
- ✅ 提供用户自行配置API密钥的界面

### 2. 代码签名

- Windows：使用Code Signing证书
- macOS：使用Apple Developer证书
- Linux：使用GPG签名

### 3. 安全更新

- 使用HTTPS分发更新
- 验证更新包的签名
- 提供校验和（SHA256）

## 性能优化

### 构建优化

```bash
# 启用UPX压缩
wails build -upx

# 去除调试信息
wails build -ldflags "-s -w"
```

### 资源优化

- 压缩前端资源（JS、CSS、图片）
- 使用代码分割减小初始加载大小
- 优化图标和图片资源

## 监控和分析

### 错误追踪

集成错误追踪服务（如Sentry）：

```go
import "github.com/getsentry/sentry-go"

func init() {
    sentry.Init(sentry.ClientOptions{
        Dsn: "your-sentry-dsn",
    })
}
```

### 使用统计

可以集成匿名使用统计（需用户同意）：
- 应用启动次数
- 功能使用频率
- 错误发生率

## 部署检查清单

发布前确认：

- [ ] 所有测试通过
- [ ] 版本号已更新
- [ ] CHANGELOG已更新
- [ ] 所有平台构建成功
- [ ] 应用图标正确
- [ ] 代码已签名（如适用）
- [ ] README和文档完整
- [ ] LICENSE文件存在
- [ ] 环境变量配置说明清晰
- [ ] 已在目标平台测试运行

## 发布后

- 监控用户反馈
- 收集错误报告
- 准备热修复版本（如需要）
- 规划下一版本功能

---

**发布负责人**：_________  
**发布日期**：_________  
**发布版本**：v1.0.0
