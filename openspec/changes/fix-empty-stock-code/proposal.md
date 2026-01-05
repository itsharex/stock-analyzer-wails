# Change: 修复自选股添加时股票代码为空的 Bug

## Why
在实施自选股添加失败修复时，发现了一个新的 bug：添加到自选股的股票数据中 `Code` 字段为空字符串。这会导致：

1. **自选股显示异常**：股票代码显示为空，用户无法识别股票
2. **功能逻辑错误**：依赖股票代码的功能（如删除、查询）可能失效
3. **数据完整性问题**：存储的数据不完整，影响数据质量

## What Changes
- 修复 `services/stock_service.go` 中 `GetStockByCode` 方法的股票代码获取逻辑
- 添加调试日志来追踪 API 返回数据的完整性
- 实现备用逻辑：当 API 返回的代码为空时，使用输入的代码参数
- 增强数据验证，确保返回的股票数据包含必要字段

## Impact
- Affected specs: `stock-data-retrieval`（新建）
- Affected code:
  - `services/stock_service.go`: 股票代码获取和数据转换逻辑
  - 测试验证：确认修复后新添加的股票有正确的代码字段
