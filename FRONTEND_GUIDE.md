# 前端开发指南

## 项目结构

```
frontend/
├── src/
│   ├── components/              # React 组件
│   │   ├── StockSearch.tsx      # 股票搜索组件
│   │   ├── StockInfo.tsx        # 股票信息展示组件
│   │   └── AnalysisReport.tsx   # 分析报告展示组件
│   ├── hooks/                   # 自定义 React hooks
│   │   └── useWailsAPI.ts       # Wails API 调用 hooks
│   ├── __tests__/               # 测试文件
│   │   └── integration.test.ts  # 集成测试
│   ├── types.ts                 # TypeScript 类型定义
│   ├── global.d.ts              # 全局类型声明
│   ├── App.tsx                  # 主应用组件
│   ├── main.tsx                 # 应用入口
│   └── index.css                # 全局样式
├── package.json                 # 依赖管理
├── vite.config.ts               # Vite 配置
├── tsconfig.json                # TypeScript 配置
└── tailwind.config.js           # TailwindCSS 配置
```

## 开发环境设置

### 1. 安装依赖

```bash
cd frontend
npm install
```

### 2. 启动开发服务器

在项目根目录运行：

```bash
wails dev
```

这将启动 Wails 开发服务器，并自动打开桌面应用窗口。前端支持热重载。

### 3. 构建生产版本

```bash
wails build
```

## 前端架构

### 组件层次结构

```
App
├── StockSearch          # 搜索和输入组件
├── StockInfo            # 股票信息展示
└── AnalysisReport       # 分析报告展示
```

### 状态管理

使用 React 的 `useState` hook 进行本地状态管理：

- `stockData`: 当前选中的股票数据
- `analysisReport`: AI 生成的分析报告
- `loading`: 加载状态
- `error`: 错误信息

### Wails API 集成

通过 `useWailsAPI` hook 与后端通信：

```typescript
const { getStockData, analyzeStock, quickAnalyze } = useWailsAPI()

// 获取股票数据
const data = await getStockData('600519')

// 进行 AI 分析
const report = await analyzeStock('600519')

// 快速分析
const analysis = await quickAnalyze('600519')
```

## 类型定义

### StockData

股票实时行情数据：

```typescript
interface StockData {
  code: string              // 股票代码
  name: string              // 股票名称
  price: number             // 最新价
  change: number            // 涨跌额
  changeRate: number        // 涨跌幅 (%)
  volume: number            // 成交量 (手)
  amount: number            // 成交额
  high: number              // 最高价
  low: number               // 最低价
  open: number              // 今开
  preClose: number          // 昨收
  amplitude: number         // 振幅 (%)
  turnover: number          // 换手率 (%)
  pe: number                // 市盈率
  pb: number                // 市净率
  totalMV: number           // 总市值
  circMV: number            // 流通市值
}
```

### AnalysisReport

AI 分析报告：

```typescript
interface AnalysisReport {
  stockCode: string         // 股票代码
  stockName: string         // 股票名称
  summary: string           // 分析摘要
  fundamentals: string      // 基本面分析
  technical: string         // 技术面分析
  recommendation: string    // 投资建议
  riskLevel: string         // 风险等级
  targetPrice: string       // 目标价位
  generatedAt: string       // 生成时间
}
```

## 样式系统

项目使用 **TailwindCSS** 进行样式管理。

### 颜色方案

- **上涨**：红色 (`text-red-500`, `bg-red-100`)
- **下跌**：绿色 (`text-green-500`, `bg-green-100`)
- **主色**：蓝色 (`bg-blue-500`)
- **强调色**：紫色/靛蓝 (`from-purple-500 to-indigo-600`)

### 响应式设计

使用 TailwindCSS 的响应式前缀：

```tsx
<div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
  {/* 移动端：1列，桌面端：3列 */}
</div>
```

## 常见问题

### Q: 如何调试前端代码？

A: 在 Wails 开发模式下，可以使用浏览器开发者工具：
1. 右键点击应用窗口，选择"检查"
2. 在 DevTools 中调试 JavaScript 代码

### Q: 如何修改 API 调用？

A: 修改 `frontend/src/hooks/useWailsAPI.ts` 文件中的相应方法。

### Q: 如何添加新的组件？

A: 
1. 在 `frontend/src/components/` 目录创建新文件
2. 定义组件并导出
3. 在 `App.tsx` 中导入并使用

### Q: 如何处理错误？

A: 错误信息通过 `onError` 回调传递给父组件，最终显示在错误提示框中。

## 性能优化建议

1. **使用 React.memo** 避免不必要的重渲染
2. **使用 useCallback** 缓存回调函数
3. **使用 useMemo** 缓存计算结果
4. **懒加载组件** 使用 React.lazy

## 测试

### 运行测试

```bash
npm test
```

### 编写测试

在 `frontend/src/__tests__/` 目录创建测试文件，使用 Vitest 框架。

## 部署

前端代码会随 Wails 应用一起打包。运行以下命令生成可执行文件：

```bash
wails build -platform windows/amd64  # Windows
wails build -platform darwin/universal # macOS
wails build -platform linux/amd64     # Linux
```

生成的可执行文件位于 `build/bin/` 目录。

## 相关资源

- [Wails 官方文档](https://wails.io/)
- [React 文档](https://react.dev/)
- [TailwindCSS 文档](https://tailwindcss.com/)
- [TypeScript 文档](https://www.typescriptlang.org/)
