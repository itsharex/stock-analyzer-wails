# 快速开始指南

本指南将帮助您在5分钟内启动A股股票分析AI-Agent应用。

## 第一步：安装前置环境

### Windows用户

1. **安装Go语言**
   - 访问 https://go.dev/dl/
   - 下载并安装 `go1.22.windows-amd64.msi`
   - 安装完成后，打开命令提示符，运行 `go version` 验证

2. **安装Node.js**
   - 访问 https://nodejs.org/
   - 下载并安装LTS版本
   - 安装完成后，运行 `node -v` 验证

3. **安装WebView2**
   - Windows 10/11通常已预装
   - 如未安装，访问：https://developer.microsoft.com/microsoft-edge/webview2/

### macOS用户

```bash
# 安装Homebrew（如已安装可跳过）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安装Go
brew install go

# 安装Node.js
brew install node

# 安装Xcode命令行工具
xcode-select --install
```

### Linux用户（Ubuntu/Debian）

```bash
# 安装Go
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安装GTK和WebKit依赖
sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.0-dev
```

## 第二步：安装Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

验证安装：
```bash
wails version
```

## 第三步：配置项目

1. **进入项目目录**
```bash
cd stock-analyzer-wails
```

2. **配置环境变量**

创建 `.env` 文件（复制 `.env.example`）：
```bash
cp .env.example .env
```

编辑 `.env` 文件，填入您的OpenAI API密钥：
```
OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

**如何获取OpenAI API密钥？**
- 访问 https://platform.openai.com/
- 注册/登录账号
- 进入 API Keys 页面创建新密钥

3. **安装依赖**

```bash
# 安装Go依赖
go mod download

# 安装前端依赖
cd frontend
npm install
cd ..
```

## 第四步：运行应用

### 开发模式（推荐）

```bash
wails dev
```

应用将自动打开，支持热重载。修改代码后会自动刷新。

### 构建生产版本

```bash
wails build
```

构建完成后，可执行文件位于：
- Windows: `build/bin/stock-analyzer.exe`
- macOS: `build/bin/stock-analyzer.app`
- Linux: `build/bin/stock-analyzer`

## 第五步：使用应用

1. 在输入框中输入股票代码（如：600519）
2. 点击"查询数据"查看实时行情
3. 点击"AI分析"获取专业分析报告

## 常见问题

### Q: 提示"go: command not found"
A: Go环境变量未配置，需要将Go的bin目录添加到PATH

### Q: 提示"wails: command not found"
A: Wails CLI未安装或未添加到PATH，运行：
```bash
export PATH=$PATH:$HOME/go/bin
```

### Q: 前端依赖安装失败
A: 尝试使用国内镜像：
```bash
npm config set registry https://registry.npmmirror.com
npm install
```

### Q: OpenAI API调用失败
A: 
1. 检查API密钥是否正确
2. 确认账户有余额
3. 如在国内，可能需要配置代理或使用国内OpenAI服务

### Q: 股票数据获取失败
A: 
1. 检查网络连接
2. 确认股票代码格式正确（6位数字）
3. 东方财富API可能有访问限制

## 技术支持

如遇到其他问题，请查看完整文档 `README.md` 或提交Issue。

---

**祝您使用愉快！** 🚀
