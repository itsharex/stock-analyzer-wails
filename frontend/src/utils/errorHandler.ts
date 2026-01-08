// 错误类型枚举
export enum ErrorType {
  NETWORK = 'network',
  PERMISSION = 'permission',
  DISK_SPACE = 'disk_space',
  FILE_CORRUPTED = 'file_corrupted',
  SERVICE_UNAVAILABLE = 'service_unavailable',
  VALIDATION = 'validation',
  UNKNOWN = 'unknown'
}

// 错误处理结果接口
export interface ErrorHandlingResult {
  type: ErrorType
  message: string
  suggestion?: string
  canRetry: boolean
}

// 解析错误信息并返回友好的中文提示
export function parseError(error: any): ErrorHandlingResult {
  const errorMessage = error?.message || error?.toString() || '未知错误'

  // 价格预警模块：保留更具体的报错，不要被“初始化/不可用”规则泛化
  if (
    errorMessage.includes('价格预警模块未初始化') ||
    errorMessage.includes('PriceAlertController')
  ) {
    return {
      type: ErrorType.SERVICE_UNAVAILABLE,
      message: errorMessage,
      suggestion: '请确认后端数据库初始化成功后再重试（如仍失败，请查看后端启动日志中的 SQLite 初始化错误）',
      canRetry: false
    }
  }

  // 建仓分析（ENTRY_*）错误码优先解析，避免被下面的通用规则误判
  const codeMatch = errorMessage.match(/code=([A-Z0-9_]+)/)
  if (codeMatch) {
    const code = codeMatch[1]
    if (code.startsWith('ENTRY_')) {
      const traceMatch = errorMessage.match(/traceId=([a-zA-Z0-9]+)/)
      const traceId = traceMatch ? traceMatch[1] : ''
      const withTrace = (msg: string) => traceId ? `${msg}（traceId=${traceId}）` : msg

      switch (code) {
        case 'ENTRY_AI_NOT_READY':
          return {
            type: ErrorType.SERVICE_UNAVAILABLE,
            message: withTrace('AI 服务未就绪，无法进行建仓分析'),
            suggestion: '请先在设置中配置 API Key/模型，然后重试',
            canRetry: false
          }
        case 'ENTRY_INPUT_INVALID':
          return {
            type: ErrorType.VALIDATION,
            message: withTrace('输入参数无效（股票代码为空或格式不正确）'),
            suggestion: '请检查股票代码后重试',
            canRetry: false
          }
        case 'ENTRY_KLINE_INSUFFICIENT':
          return {
            type: ErrorType.SERVICE_UNAVAILABLE,
            message: withTrace('K 线数据不足，暂无法生成建仓方案'),
            suggestion: '可稍后再试，或更换股票/检查数据源是否正常',
            canRetry: true
          }
        case 'ENTRY_AI_TIMEOUT':
          return {
            type: ErrorType.NETWORK,
            message: withTrace('AI 分析超时'),
            suggestion: '请稍后重试，或检查网络/AI 配置是否可用',
            canRetry: true
          }
        case 'ENTRY_AI_INVALID_JSON':
          return {
            type: ErrorType.SERVICE_UNAVAILABLE,
            message: withTrace('AI 返回内容无法解析（格式异常）'),
            suggestion: '请稍后重试；如频繁出现可切换模型或查看日志定位',
            canRetry: true
          }
        case 'ENTRY_PANIC':
          return {
            type: ErrorType.UNKNOWN,
            message: withTrace('后端发生异常，建仓分析中断'),
            suggestion: '请在日志中搜索 traceId 定位原因，然后重试',
            canRetry: true
          }
        default:
          return {
            type: ErrorType.UNKNOWN,
            message: withTrace(`建仓分析失败（${code}）`),
            suggestion: '请稍后重试；如持续失败请查看日志',
            canRetry: true
          }
      }
    }
  }
  
  // 网络相关错误
  if (errorMessage.includes('fetch') || 
      errorMessage.includes('network') || 
      errorMessage.includes('连接') ||
      errorMessage.includes('timeout')) {
    return {
      type: ErrorType.NETWORK,
      message: '网络连接异常，请检查网络设置',
      suggestion: '请检查网络连接后重试，或稍后再试',
      canRetry: true
    }
  }

  // 文件权限错误
  if (errorMessage.includes('权限') || 
      errorMessage.includes('permission') ||
      errorMessage.includes('access denied')) {
    return {
      type: ErrorType.PERMISSION,
      message: '文件权限不足，无法保存数据',
      suggestion: '请检查应用是否有足够的文件访问权限，或以管理员身份运行',
      canRetry: false
    }
  }

  // 磁盘空间不足
  if (errorMessage.includes('磁盘') || 
      errorMessage.includes('空间') ||
      errorMessage.includes('disk') ||
      errorMessage.includes('space')) {
    return {
      type: ErrorType.DISK_SPACE,
      message: '磁盘空间不足，无法保存数据',
      suggestion: '请清理磁盘空间后重试',
      canRetry: true
    }
  }

  // 文件损坏
  if (errorMessage.includes('解析') || 
      errorMessage.includes('JSON') ||
      errorMessage.includes('格式') ||
      errorMessage.includes('corrupted')) {
    return {
      type: ErrorType.FILE_CORRUPTED,
      message: '数据格式异常或解析失败',
      suggestion: '请检查输入内容或稍后重试；如持续出现请查看日志定位',
      canRetry: true
    }
  }

  // 服务不可用
  if (errorMessage.includes('服务') || 
      errorMessage.includes('初始化') ||
      errorMessage.includes('不可用') ||
      errorMessage.includes('service') ||
      errorMessage.includes('unavailable')) {
    return {
      type: ErrorType.SERVICE_UNAVAILABLE,
      message: '功能暂时不可用',
      suggestion: '请重启应用；如持续失败请查看后端日志/数据库初始化状态',
      canRetry: false
    }
  }

  // 验证错误
  if (errorMessage.includes('代码') || 
      errorMessage.includes('输入') ||
      errorMessage.includes('validation') ||
      errorMessage.includes('invalid')) {
    return {
      type: ErrorType.VALIDATION,
      message: errorMessage,
      suggestion: '请检查输入的股票代码是否正确',
      canRetry: false
    }
  }

  // 默认未知错误
  return {
    type: ErrorType.UNKNOWN,
    message: `操作失败: ${errorMessage}`,
    suggestion: '请稍后重试，如问题持续存在请联系技术支持',
    canRetry: true
  }
}

// 格式化错误消息用于显示
export function formatErrorMessage(result: ErrorHandlingResult): string {
  let message = result.message
  
  if (result.suggestion) {
    message += `\n\n💡 建议: ${result.suggestion}`
  }
  
  if (result.canRetry) {
    message += '\n\n🔄 您可以重试此操作'
  }
  
  return message
}

// 简化的错误处理函数，直接返回格式化的消息
export function handleError(error: any): string {
  const result = parseError(error)
  return formatErrorMessage(result)
}
