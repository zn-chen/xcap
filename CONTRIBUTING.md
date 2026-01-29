# 贡献指南

## 代码规范

### 注释语言

- **代码注释使用中文**，专业术语保持英文
- **日志输出使用英文**，便于国际化和日志分析

示例：

```go
// GetAllMonitors 获取所有活动的显示器
// 返回 Monitor 列表，如果没有显示器则返回 error
func GetAllMonitors() ([]Monitor, error) {
    monitors, err := enumerate()
    if err != nil {
        log.Printf("Failed to enumerate monitors: %v", err)  // 日志用英文
        return nil, err
    }
    return monitors, nil
}
```

### 命名规范

- 公共函数、类型使用 PascalCase
- 私有函数、变量使用 camelCase
- 常量使用 PascalCase 或全大写下划线分隔

### 错误处理

- 使用 `pkg/owl/errors.go` 中定义的标准错误
- 错误消息使用英文

## 目录结构

```
owl-go/
├── pkg/owl/          # 公共 API（用户使用）
├── internal/darwin/  # macOS 内部实现
├── internal/windows/ # Windows 内部实现
├── examples/         # 示例代码
└── docs/             # 文档
```

## 提交规范

提交消息格式：
```
<type>: <description>

[optional body]
```

类型：
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具相关
