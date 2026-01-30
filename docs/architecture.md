# xcap 架构设计

基于 xcap 的 Go 语言跨平台屏幕截图库。

## 项目概述

xcap 是一个 Go 语言实现的跨平台屏幕和窗口截图库，参考 Rust 库 [xcap](https://github.com/nashaofu/xcap) 的设计和实现。

## 功能特性

### 核心功能

| 功能 | macOS | Windows | Linux (未来) |
|------|-------|---------|--------------|
| 屏幕截图 | 计划 | 计划 | 计划 |
| 窗口截图 | 计划 | 计划 | 计划 |
| 区域截图 | 计划 | 计划 | 计划 |
| 显示器枚举 | 计划 | 计划 | 计划 |
| 窗口枚举 | 计划 | 计划 | 计划 |

## 项目结构

```
xcap/
├── doc/                          # 文档
│   ├── macos-implementation.md   # macOS 实现原理
│   ├── windows-implementation.md # Windows 实现原理
│   └── architecture.md           # 架构设计
├── pkg/
│   └── xcap/                     # 主包
│       ├── monitor.go            # 显示器接口
│       ├── window.go             # 窗口接口
│       ├── capture.go            # 截图通用逻辑
│       └── errors.go             # 错误定义
├── internal/
│   ├── darwin/                   # macOS 实现
│   │   ├── monitor.go
│   │   ├── window.go
│   │   ├── capture.go
│   │   └── cgo.go                # CGO 绑定
│   └── windows/                  # Windows 实现
│       ├── monitor.go
│       ├── window.go
│       ├── capture.go
│       └── syscall.go            # Windows API 调用
├── examples/                     # 示例代码
│   ├── monitor_capture/
│   ├── window_capture/
│   └── list_windows/
├── go.mod
├── go.sum
└── README.md
```

## 公共接口设计

### Monitor 接口

```go
package xcap

import "image"

// Monitor 表示一个显示器
type Monitor interface {
    // 基本信息
    ID() uint32
    Name() string

    // 位置和尺寸
    X() int
    Y() int
    Width() uint32
    Height() uint32

    // 属性
    Rotation() float32
    ScaleFactor() float32
    Frequency() float32
    IsPrimary() bool
    IsBuiltin() bool

    // 截图
    CaptureImage() (*image.RGBA, error)
    CaptureRegion(x, y, width, height uint32) (*image.RGBA, error)
}

// 获取所有显示器
func AllMonitors() ([]Monitor, error)

// 根据坐标获取显示器
func MonitorFromPoint(x, y int) (Monitor, error)
```

### Window 接口

```go
package xcap

import "image"

// Window 表示一个窗口
type Window interface {
    // 基本信息
    ID() uint32
    PID() uint32
    AppName() string
    Title() string

    // 位置和尺寸
    X() int
    Y() int
    Z() int  // Z 顺序
    Width() uint32
    Height() uint32

    // 状态
    IsMinimized() bool
    IsMaximized() bool
    IsFocused() bool

    // 关联
    CurrentMonitor() (Monitor, error)

    // 截图
    CaptureImage() (*image.RGBA, error)
}

// 获取所有窗口
func AllWindows() ([]Window, error)

// 获取当前焦点窗口
func FocusedWindow() (Window, error)
```

## 平台实现策略

### macOS 实现

使用 CGO 调用 Core Graphics 和 AppKit 框架：

```go
// internal/darwin/cgo.go

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework AppKit -framework CoreFoundation

#include <CoreGraphics/CoreGraphics.h>
#include <AppKit/AppKit.h>

// C 函数声明
CGImageRef captureScreen(CGRect bounds);
CGImageRef captureWindow(uint32_t windowID);
*/
import "C"
```

### Windows 实现

使用 syscall 调用 Win32 API：

```go
// internal/windows/syscall.go

import "golang.org/x/sys/windows"

var (
    user32 = windows.NewLazySystemDLL("user32.dll")
    gdi32  = windows.NewLazySystemDLL("gdi32.dll")
    dwmapi = windows.NewLazySystemDLL("dwmapi.dll")
    shcore = windows.NewLazySystemDLL("shcore.dll")
)
```

## 图像处理

统一使用 Go 标准库 `image` 包：

```go
import "image"

// 所有截图返回 *image.RGBA
func CaptureImage() (*image.RGBA, error) {
    // 1. 调用平台 API 获取原始数据
    // 2. BGRA → RGBA 转换
    // 3. 创建 image.RGBA
    return img, nil
}
```

## 错误处理

```go
package xcap

import "errors"

var (
    ErrNoMonitor       = errors.New("no monitor found")
    ErrNoWindow        = errors.New("no window found")
    ErrCaptureFailed   = errors.New("capture failed")
    ErrPermissionDenied = errors.New("permission denied")
    ErrWindowMinimized = errors.New("window is minimized")
)

// XcapError 包装平台特定错误
type XcapError struct {
    Op      string  // 操作名称
    Err     error   // ��层错误
}

func (e *XcapError) Error() string {
    return fmt.Sprintf("xcap: %s: %v", e.Op, e.Err)
}

func (e *XcapError) Unwrap() error {
    return e.Err
}
```

## 依赖

- Go 1.21+
- macOS: Xcode Command Line Tools (for CGO)
- Windows: 无额外依赖

## 开发计划

### Phase 1: 基础框架
- [x] 项目初始化
- [x] 文档整理
- [ ] 定义公共接口

### Phase 2: macOS 实现
- [ ] Monitor 枚举和截图
- [ ] Window 枚举和截图
- [ ] 权限检查

### Phase 3: Windows 实现
- [ ] Monitor 枚举和截图
- [ ] Window 枚举和截图
- [ ] DPI 处理

### Phase 4: 完善
- [ ] 示例代码
- [ ] 单元测试
- [ ] 性能优化
- [ ] 文档完善
