# xcap

中文 | [English](README.md)

跨平台屏幕和窗口截图 Go 语言库，参考 [xcap](https://github.com/nashaofu/xcap) 实现。

## 功能特性

| 功能 | macOS | Windows | 说明 |
|------|:-----:|:-------:|------|
| 显示器截图 | ✅ | ✅ | 支持高 DPI |
| 窗口截图 | ✅ | ✅ | 支持后台窗口 |
| 多显示器支持 | ✅ | ✅ | |
| Monitor.IsPrimary | ✅ | ✅ | 是否主显示器 |
| Monitor.ScaleFactor | ✅ | ✅ | 返回缩放比例（如 Retina 为 2.0）|
| Monitor.Rotation | ✅ | ✅ | 屏幕旋转角度 |
| Monitor.Frequency | ✅ | ✅ | 刷新率（Hz）|
| Window.IsFocused | ✅ | ✅ | 是否获得焦点 |
| Window.IsMinimized | ❌ | ✅ | macOS 返回 `ErrNotSupported` |
| Window.IsMaximized | ❌ | ✅ | macOS 返回 `ErrNotSupported` |
| 排除当前进程窗口 | ✅ | ✅ | 过滤自身窗口 |
| 区域截图 | ❌ | ❌ | 计划中 |

## 安装

```bash
go get github.com/zn-chen/xcap
```

## 快速开始

### 显示器截图

```go
package main

import (
    "fmt"
    "image/png"
    "log"
    "os"

    "github.com/zn-chen/xcap/pkg/xcap"
)

func main() {
    monitors, err := xcap.AllMonitors()
    if err != nil {
        log.Fatal(err)
    }

    for i, m := range monitors {
        fmt.Printf("显示器 %d: %s (%dx%d @ %.0fHz, 缩放=%.1fx, 主显示器=%v)\n",
            i, m.Name(), m.Width(), m.Height(), m.Frequency(), m.ScaleFactor(), m.IsPrimary())

        img, err := m.CaptureImage()
        if err != nil {
            log.Printf("截取显示器 %d 失败: %v", i, err)
            continue
        }

        f, _ := os.Create(fmt.Sprintf("monitor-%d.png", i))
        png.Encode(f, img)
        f.Close()
    }
}
```

### 窗口截图

```go
package main

import (
    "fmt"
    "image/png"
    "log"
    "os"

    "github.com/zn-chen/xcap/pkg/xcap"
)

func main() {
    // 使用 AllWindowsWithOptions(true) 排除当前进程的窗口
    windows, err := xcap.AllWindows()
    if err != nil {
        log.Fatal(err)
    }

    for i, w := range windows {
        // 跳过最小化窗口（检查错误以处理不支持的平台）
        if minimized, err := w.IsMinimized(); err == nil && minimized {
            continue
        }

        fmt.Printf("窗口 %d: [%s] %s (%dx%d)\n",
            i, w.AppName(), w.Title(), w.Width(), w.Height())

        img, err := w.CaptureImage()
        if err != nil {
            log.Printf("截取窗口失败: %v", err)
            continue
        }

        f, _ := os.Create(fmt.Sprintf("window-%d.png", i))
        png.Encode(f, img)
        f.Close()
    }
}
```

## CLI 工具

```bash
# 构建
make build

# 运行（截取所有显示器和窗口到 ./output/）
./bin/xcap

# 只截取显示器
./bin/xcap --disable_windows

# 只截取窗口
./bin/xcap --disable_monitor
```

## API 参考

### Monitor 接口

```go
type Monitor interface {
    ID() uint32              // 唯一标识符
    Name() string            // 显示器名称
    X() int                  // X 坐标
    Y() int                  // Y 坐标
    Width() uint32           // 宽度（像素）
    Height() uint32          // 高度（像素）
    Rotation() float64       // 旋转角度
    ScaleFactor() float64    // DPI 缩放比例
    Frequency() float64      // 刷新率（Hz）
    IsPrimary() bool         // 是否主显示器
    IsBuiltin() bool         // 是否内置显示器（笔记本）
    CaptureImage() (*image.RGBA, error)
}
```

### Window 接口

```go
type Window interface {
    ID() uint32              // 唯一标识符（Windows 为 HWND，macOS 为 CGWindowID）
    PID() uint32             // 进程 ID
    AppName() string         // 应用程序名称
    Title() string           // 窗口标题
    X() int                  // X 坐标
    Y() int                  // Y 坐标
    Z() int                  // Z 顺序（值越大越靠前）
    Width() uint32           // 宽度（像素）
    Height() uint32          // 高度（像素）
    IsMinimized() (bool, error)  // 不支持时返回 ErrNotSupported
    IsMaximized() (bool, error)  // 不支持时返回 ErrNotSupported
    IsFocused() (bool, error)    // 是否拥有输入焦点
    CurrentMonitor() (Monitor, error)
    CaptureImage() (*image.RGBA, error)
}
```

### 函数

```go
// 获取所有显示器
func AllMonitors() ([]Monitor, error)

// 获取所有可见窗口
func AllWindows() ([]Window, error)

// 获取所有可见窗口（带选项）
// excludeCurrentProcess: 是否排除当前进程的窗口
func AllWindowsWithOptions(excludeCurrentProcess bool) ([]Window, error)

// 清理文件名中的非法字符
func SanitizeFilename(name string) string
```

## 项目结构

```
xcap/
├── cmd/xcap/           # 命令行工具
├── pkg/xcap/           # 公共 API（跨平台接口）
├── internal/
│   ├── darwin/         # macOS 实现（CGO + Objective-C）
│   └── windows/        # Windows 实现（CGO + C）
├── examples/           # 使用示例
└── docs/               # 文档
```

## 平台要求

### macOS

- macOS 10.15 (Catalina) 或更高版本
- 需要屏幕录制权限
  - 系统设置 > 隐私与安全性 > 屏幕录制
  - 将您的应用程序添加到允许列表
- 需要 Xcode Command Line Tools（用于 CGO 编译）
  ```bash
  xcode-select --install
  ```

### Windows

- Windows 8.1 或更高版本
- 需要 MinGW-w64（用于 CGO 编译）
  - 通过 [MSYS2](https://www.msys2.org/) 或 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) 安装
- 无需额外权限

## 文档

- [macOS 实现原理](docs/macos-implementation.md) - CoreGraphics API 详解
- [Windows 实现原理](docs/windows-implementation.md) - GDI/Win32 API 详解
- [架构设计](docs/architecture.md) - 设计概述

## 贡献

欢迎贡献代码！请参阅 [CONTRIBUTING.md](CONTRIBUTING.md) 了解贡献指南。

## 许可证

Apache-2.0
