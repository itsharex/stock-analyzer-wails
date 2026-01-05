## ADDED Requirements

### Requirement: 自选股添加错误处理
系统在添加股票到自选股时 SHALL 提供完善的错误处理和用户反馈。

#### Scenario: JSON 解析错误处理
- **WHEN** 自选股文件存在但包含无效的 JSON 数据
- **THEN** 系统 SHALL 记录详细的错误日志并返回明确的错误信息
- **AND** 系统 SHALL 尝试备份损坏的文件并创建新的空自选股列表

#### Scenario: 文件权限错误处理
- **WHEN** 用户没有足够权限写入自选股文件
- **THEN** 系统 SHALL 记录权限错误的详细信息
- **AND** 前端 SHALL 显示友好的错误提示，建议用户检查文件权限

#### Scenario: 磁盘空间不足错误处理
- **WHEN** 磁盘空间不足导致文件写入失败
- **THEN** 系统 SHALL 记录磁盘空间错误
- **AND** 前端 SHALL 提示用户清理磁盘空间

### Requirement: 自选股操作日志记录
系统 SHALL 为所有自选股操作记录结构化日志。

#### Scenario: 添加操作日志
- **WHEN** 用户尝试添加股票到自选股
- **THEN** 系统 SHALL 记录包含以下字段的日志：
  - `module`: "services.watchlist"
  - `op`: "add_to_watchlist"
  - `stock_code`: 股票代码
  - `file_path`: 自选股文件路径
  - `duration_ms`: 操作耗时
  - `success`: 操作是否成功

#### Scenario: 错误操作日志
- **WHEN** 自选股操作失败
- **THEN** 系统 SHALL 记录错误日志包含：
  - 上述基本字段
  - `err`: 错误详情
  - `file_exists`: 文件是否存在
  - `file_size`: 文件大小（如果存在）

### Requirement: 前端错误提示优化
前端 SHALL 为自选股操作提供详细且用户友好的错误提示。

#### Scenario: 具体错误信息显示
- **WHEN** 后端返回具体的错误信息
- **THEN** 前端 SHALL 解析错误类型并显示相应的中文提示
- **AND** 对于权限错误，SHALL 提供解决建议

#### Scenario: 网络错误处理
- **WHEN** 前端调用后端 API 时发生网络错误
- **THEN** 前端 SHALL 显示 "网络连接异常，请检查网络设置" 提示
- **AND** 提供重试选项

### Requirement: 应用初始化错误处理
应用启动时 SHALL 正确处理自选股服务初始化错误。

#### Scenario: 自选股服务初始化失败
- **WHEN** `NewWatchlistService()` 返回错误
- **THEN** 应用 SHALL 记录初始化错误日志
- **AND** 应用 SHALL 继续启动但禁用自选股功能
- **AND** 前端 SHALL 显示自选股功能不可用的提示

#### Scenario: 自选股目录创建失败
- **WHEN** 无法创建应用数据目录
- **THEN** 系统 SHALL 记录目录创建失败的详细错误
- **AND** 尝试使用临时目录作为备选方案
