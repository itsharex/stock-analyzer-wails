# A股股票分析AI-Agent

基于Go + Wails框架的A股股票分析桌面应用程序，通过东方财富API获取实时股票数据，结合 **CloudWeGo Eino** 框架和 **阿里百炼 (DashScope)** 进行智能分析，生成专业的投资建议报告。

## 功能特性

- **股票数据获取**：通过东方财富API实时获取A股股票行情数据
- **AI智能分析**：集成 **Eino** 框架和 **阿里百炼 (DashScope)** 进行智能股票分析
- **专业报告生成**：自动生成包含基本面分析、技术面分析和投资建议的专业报告
- **现代化界面**：基于React的美观用户界面
- **跨平台支持**：支持Windows、macOS和Linux

## 技术栈

- **后端**：Go 1.22+
- **框架**：Wails v2
- **AI框架**：CloudWeGo Eino
- **前端**：React + TypeScript + TailwindCSS
- **AI服务**：阿里百炼 (DashScope) - 通义千问 (Qwen)
- **数据源**：东方财富网API

## 前置要求

### 1. 安装Go语言环境

**Windows**:
- 下载并安装：https://go.dev/dl/
- 推荐版本：Go 1.22或更高

**macOS**:
```bash
brew install go
```

**Linux**:
```bash
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 2. 安装Node.js

**Windows/macOS**:
- 下载并安装：https://nodejs.org/
- 推荐版本：Node.js 18+

**Linux**:
```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

### 3. 安装Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 4. 平台特定依赖

**Windows**:
- 需要安装WebView2运行时（Windows 10/11通常已预装）
- 下载地址：https://developer.microsoft.com/en-us/microsoft-edge/webview2/

**macOS**:
- 需要Xcode命令行工具
```bash
xcode-select --install
```

**Linux (Ubuntu/Debian)**:
```bash
sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev
```

## 安装步骤

### 1. 克隆或下载项目

将项目文件解压到本地目录，例如：`/path/to/stock-analyzer-wails`

### 2. 配置环境变量
创建 `.env` 文件（或设置系统环境变量）：

```bash
# 阿里百炼 API密钥
export DASHSCOPE_API_KEY="your-dashscope-api-key-here"

# 可选：模型名称（默认 qwen-plus）
export DASHSCOPE_MODEL="qwen-plus"

# 可选：API基础URL
export DASHSCOPE_BASE_URL="https://dashscope.aliyuncs.com/compatible-mode/v1"
```

### 3. 安装依赖

```bash
cd stock-analyzer-wails

# 安装Go依赖
go mod download

# 安装前端依赖
cd frontend
npm install
cd ..
```

## 运行应用

### 开发模式

```bash
wails dev
```

这将启动开发服务器，支持热重载。应用会自动打开桌面窗口。

### 构建生产版本

```bash
# 构建当前平台的可执行文件
wails build

# 构建后的文件位于 build/bin/ 目录
```

**跨平台构建**：
```bash
# Windows
wails build -platform windows/amd64

# macOS
wails build -platform darwin/universal

# Linux
wails build -platform linux/amd64
```

## 使用说明

1. **启动应用**：运行开发模式或双击构建后的可执行文件
2. **输入股票代码**：在输入框中输入A股代码（如：600519、000001）
3. **点击分析**：系统将自动获取股票数据并进行AI分析
4. **查看报告**：分析完成后显示详细的投资分析报告

## 项目结构

```
stock-analyzer-wails/
├── main.go                 # 应用程序入口
├── app.go                  # 应用程序主逻辑
├── go.mod                  # Go依赖管理
├── go.sum                  # Go依赖校验
├── wails.json             # Wails配置文件
├── services/              # 业务服务层
│   ├── stock_service.go   # 股票数据服务
│   └── ai_service.go      # AI分析服务 (使用 Eino 框架)
├── models/                # 数据模型
│   └── stock.go           # 股票数据模型
├── frontend/              # 前端代码
│   ├── src/
│   │   ├── App.tsx        # 主应用组件
│   │   ├── main.tsx       # 前端入口
│   │   └── components/    # React组件
│   ├── package.json       # 前端依赖
│   └── vite.config.ts     # Vite配置
├── build/                 # 构建输出目录
└── README.md             # 项目文档
```

## AI 框架：CloudWeGo Eino

本项目使用字节跳动开源的 **Eino** 框架进行 AI 逻辑编排。Eino 提供了高度组件化和可扩展的 AI 应用开发体验。

**核心组件**：
- `ChatModel`: 使用 `eino-ext/components/model/qwen` 适配器接入阿里百炼。
- `Schema`: 使用 Eino 标准消息结构进行提示词管理。

## API说明
### 东方财富API

应用使用东方财富网的公开API获取股票数据：

**接口地址**：
```
http://78.push2.eastmoney.com/api/qt/clist/get
```

### 阿里百炼 (DashScope) API

使用阿里百炼提供的通义千问模型进行股票分析：

**模型**：qwen-plus, qwen-max 等
**功能**：
- 基本面分析
- 技术指标解读
- 投资建议生成
- 风险评估

## 常见问题

### Q: 如何获取阿里百炼 API密钥？
A: 访问 [阿里云百炼控制台](https://bailian.console.aliyun.com/) 申请 API Key。

### Q: 应用启动失败怎么办？
A: 
1. 确认已安装所有前置依赖
2. 检查Go和Node.js版本是否符合要求
3. 确认环境变量已正确配置
4. 查看终端错误信息

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！

---

**免责声明**：本应用仅供学习和研究使用，不构成任何投资建议。投资有风险，入市需谨慎。
