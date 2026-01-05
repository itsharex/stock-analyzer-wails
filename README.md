# A股股票分析AI-Agent

基于Go + Wails框架的A股股票分析桌面应用程序，通过东方财富API获取实时股票数据，结合OpenAI进行智能分析，生成专业的投资建议报告。

## 功能特性

- **股票数据获取**：通过东方财富API实时获取A股股票行情数据
- **AI智能分析**：集成OpenAI GPT模型，对股票进行深度分析
- **专业报告生成**：自动生成包含基本面分析、技术面分析和投资建议的专业报告
- **现代化界面**：基于React的美观用户界面
- **跨平台支持**：支持Windows、macOS和Linux

## 技术栈

- **后端**：Go 1.22+
- **框架**：Wails v2
- **前端**：React + TypeScript + TailwindCSS
- **AI服务**：OpenAI API
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
# OpenAI API密钥
export OPENAI_API_KEY="your-openai-api-key-here"

# 可选：OpenAI API基础URL（如使用代理或第三方服务）
export OPENAI_BASE_URL="https://api.openai.com/v1"
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
│   └── ai_service.go      # AI分析服务
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

## API说明

### 东方财富API

应用使用东方财富网的公开API获取股票数据：

**接口地址**：
```
http://78.push2.eastmoney.com/api/qt/clist/get
```

**主要参数**：
- `pn`: 页码
- `pz`: 每页数量
- `fields`: 返回字段（股票代码、名称、价格、涨跌幅等）
- `fs`: 市场筛选（沪深A股）

### OpenAI API

使用OpenAI的GPT模型进行股票分析：

**模型**：GPT-4或GPT-3.5-turbo
**功能**：
- 基本面分析
- 技术指标解读
- 投资建议生成
- 风险评估

## 常见问题

### Q: 如何获取OpenAI API密钥？
A: 访问 https://platform.openai.com/ 注册账号并创建API密钥。

### Q: 应用启动失败怎么办？
A: 
1. 确认已安装所有前置依赖
2. 检查Go和Node.js版本是否符合要求
3. 确认环境变量已正确配置
4. 查看终端错误信息

### Q: 股票数据获取失败？
A: 
1. 检查网络连接
2. 确认输入的股票代码格式正确
3. 东方财富API可能有访问频率限制

### Q: AI分析响应慢？
A: 
1. OpenAI API响应时间取决于网络和服务器负载
2. 可以考虑使用国内的OpenAI API代理服务
3. 检查OPENAI_BASE_URL配置

## 开发指南

### 添加新功能

1. **后端功能**：在 `services/` 目录添加新的服务文件
2. **前端组件**：在 `frontend/src/components/` 添加React组件
3. **数据模型**：在 `models/` 目录定义新的数据结构

### 调试技巧

- 使用 `wails dev` 启动开发模式，支持热重载
- 在浏览器中调试：应用启动后访问 http://localhost:34115
- 查看Go日志：在代码中使用 `runtime.LogInfo()` 等方法

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！

## 联系方式

如有问题或建议，请通过GitHub Issues联系。

---

**免责声明**：本应用仅供学习和研究使用，不构成任何投资建议。投资有风险，入市需谨慎。
