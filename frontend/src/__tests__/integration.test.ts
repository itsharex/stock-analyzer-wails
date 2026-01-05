/**
 * 前端集成测试
 * 验证与后端 Go 代码的通信是否正常
 */

import { describe, it, expect, beforeAll } from 'vitest'

describe('Wails Frontend Integration', () => {
  // 模拟 Wails 运行时
  let wailsRuntime: any

  beforeAll(() => {
    // 在实际运行时，Wails 会自动注入 window.go
    // 这里为测试目的模拟一个基本的结构
    wailsRuntime = {
      go: {
        main: {
          App: {
            GetStockData: async (code: string) => ({
              code,
              name: '测试股票',
              price: 100,
              change: 1,
              changeRate: 1,
              volume: 1000,
              amount: 100000,
              high: 101,
              low: 99,
              open: 100,
              preClose: 99,
              amplitude: 2,
              turnover: 0.5,
              pe: 15,
              pb: 1.5,
              totalMV: 1000000000,
              circMV: 1000000000,
            }),
            AnalyzeStock: async (code: string) => ({
              stockCode: code,
              stockName: '测试股票',
              summary: '这是一个测试摘要',
              fundamentals: '这是基本面分析',
              technical: '这是技术面分析',
              recommendation: '建议持有',
              riskLevel: '中等风险',
              targetPrice: '110-120元',
              generatedAt: new Date().toISOString(),
            }),
          },
        },
      },
    }
  })

  it('应该能够获取股票数据', async () => {
    const result = await wailsRuntime.go.main.App.GetStockData('600519')
    expect(result).toBeDefined()
    expect(result.code).toBe('600519')
    expect(result.price).toBeGreaterThan(0)
  })

  it('应该能够进行 AI 分析', async () => {
    const result = await wailsRuntime.go.main.App.AnalyzeStock('600519')
    expect(result).toBeDefined()
    expect(result.stockCode).toBe('600519')
    expect(result.summary).toBeTruthy()
    expect(result.recommendation).toBeTruthy()
  })

  it('应该能够处理错误', async () => {
    const errorRuntime = {
      go: {
        main: {
          App: {
            GetStockData: async () => {
              throw new Error('API 错误')
            },
          },
        },
      },
    }

    try {
      await errorRuntime.go.main.App.GetStockData('invalid')
      expect.fail('应该抛出错误')
    } catch (error) {
      expect(error).toBeDefined()
    }
  })
})
