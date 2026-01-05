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
      message: '文件权限不足，无法保存自选股',
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
      message: '自选股数据文件已损坏，已自动修复',
      suggestion: '数据文件已重置，请重新添加自选股',
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
      message: '自选股功能暂时不可用',
      suggestion: '请重启应用或联系技术支持',
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
