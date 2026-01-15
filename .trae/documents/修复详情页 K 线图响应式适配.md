## 问题定位
- 图表仅监听 window.resize，父容器宽度变化（如侧栏折叠/布局变更）不会触发，导致不自适配。
- 组件高度固定为传入 props，高度与实际容器不一致，缩放时出现空白或裁剪。
- 初次挂载在“隐藏/零宽容器”场景（视图切换、Tab 初次呈现）会初始化为错误尺寸。

## 解决方案
- 使用 ResizeObserver 监听图表容器的宽高变化，实时调用 chart.applyOptions({ width, height })。
- 容器样式统一为 w-full h-full，父容器提供明确高度（flex 布局下以 min-h 或固定区域高度）；KLineChart 支持 height="auto"，初始化与重算时读取 clientHeight。
- 首次宽度为 0 或显示状态变化时延迟一次 resize（requestAnimationFrame）并在尺寸稳定后 fitContent。
- 移除全局 window.resize 监听，统一用 observer；保留数据更新时的价格/时间轴不强制 fit，避免打断用户交互。

## 实施步骤
1) 新增 hooks/useResizeObserver.ts：对传入的容器元素注册 ResizeObserver，节流（rAF）后回调尺寸。
2) 更新 KLineChart.tsx：
   - 初始化后注册 observer，尺寸变化时 chart.applyOptions({ width, height })；若 width 或 height 变化幅度大于阈值则调用 chart.timeScale().fitContent()。
   - 支持 height="auto"：优先读取容器 clientHeight，否则回退到 props height。
   - 解决“初次 0 宽”问题：当首次测得 width>0 时触发一次 resize+fitContent。
   - 清理 observer。
3) 优化 WatchlistDetail.tsx 容器：
   - 主图容器设置 class：w-full h-full min-h-[400px]，确保有稳定高度来源；量价分析与资金流模块在下方分区，不影响主图自适配。
4) 验证用例：
   - 浏览器窗口缩放 + 侧栏折叠/展开 + 切换分时/K线/周期，图表宽高即时适配；首屏与视图切换不出现 0 宽。
   - 数据更新过程中不跳回起始位置；缩放与拖拽交互保持。

## 风险与回退
- ResizeObserver 在 Wails/Electron 环境可用；若异常则降级为 window.resize 监听。
- 大尺寸频繁重排可能影响性能，使用 rAF 节流与阈值判断减少不必要 fitContent。

## 验收标准
- 页面放大缩小时，K 线图宽高自适应，无空白或溢出。
- 侧栏宽度变化立即响应；切换视图或周期时不出现初始化尺寸错误。
- 手动缩放和拖拽体验不被打断，性能平稳。