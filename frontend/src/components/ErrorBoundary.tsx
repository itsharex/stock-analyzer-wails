import React from 'react'

/**
 * ErrorBoundary
 *
 * 目的：避免某个页面/组件在渲染期抛异常导致整个应用白屏。
 * 典型案例：syncStats.stock_list 为 null 时直接 join\(\) 报错。
 *
 * 行为：
 * - 捕获渲染/生命周期异常
 * - 控制台输出详细错误（开发排查）
 * - 在页面上显示友好的“应用出错”提示，并提供刷新按钮
 */
export class ErrorBoundary extends React.Component<
  React.PropsWithChildren<{}>,
  { hasError: boolean; message?: string }
> {
  state = { hasError: false as boolean, message: undefined as string | undefined }

  static getDerivedStateFromError(err: unknown) {
    const message = err instanceof Error ? err.message : String(err)
    return { hasError: true, message }
  }

  componentDidCatch(error: unknown, errorInfo: unknown) {
    // 这里尽量保留原始错误信息，方便定位。
    // Wails 桌面环境没有远程上报默认通道，先打到 console 即可。
    // 需要的话后续可以接入后端日志：window.go... 或 runtime.EventsEmit
    // eslint-disable-next-line no-console
    console.error('[ErrorBoundary] uncaught error', error)
    // eslint-disable-next-line no-console
    console.error('[ErrorBoundary] errorInfo', errorInfo)
  }

  private reload = () => {
    window.location.reload()
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen bg-gray-900 text-gray-100 flex items-center justify-center p-6">
          <div className="max-w-xl w-full bg-gray-800 border border-gray-700 rounded-lg p-6 shadow-lg">
            <h1 className="text-xl font-bold mb-2">页面发生错误</h1>
            <p className="text-gray-300 text-sm mb-4">
              应用捕获到一个未处理异常，为避免白屏已自动拦截。你可以尝试刷新页面；如果反复出现，请把控制台错误日志发给我。
            </p>
            {this.state.message && (
              <pre className="bg-gray-900 border border-gray-700 rounded p-3 text-xs text-gray-300 overflow-auto mb-4">
                {this.state.message}
              </pre>
            )}
            <div className="flex gap-3">
              <button
                onClick={this.reload}
                className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md transition-colors"
              >
                刷新
              </button>
            </div>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}

