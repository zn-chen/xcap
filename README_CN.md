# xcap

中文 | [English](README.md)

跨平台**原生窗口截图**库，使用 Go 语言实现，参考 [xcap](https://github.com/nashaofu/xcap)。

## 为什么选择 xcap？

大多数 Go 截图库只支持**区域截图** —— 从屏幕的矩形区域抓取像素。这种方式有根本性的局限：

| 方式 | 区域截图 | xcap（原生窗口截图）|
|------|:-------:|:------------------:|
| 截取被遮挡的窗口 | ❌ | ✅ |
| 截取后台窗口 | ❌ | ✅ |
| 截取最小化窗口 | ❌ | ✅ (Windows) |
| 获取窗口元数据（标题、应用名、PID）| ❌ | ✅ |
| 独立窗口隔离 | ❌ | ✅ |

**xcap 使用操作系统原生 API**（macOS 的 CoreGraphics，Windows 的 GDI/Win32）将窗口作为独立实体进行截图，而非简单的屏幕区域。这使得：

- **截取任意窗口** —— 无论是否可见或被遮挡
- **枚举所有窗口** —— 获取标题、应用名、进程ID、位置、尺寸等元数据
- **多显示器支持** —— 独立截取每个显示器
- **高 DPI 感知** —— 正确处理 Retina/HiDPI 缩放

## 功能特性

| 功能 | macOS | Windows | 说明 |
|------|:-----:|:-------:|------|
| 显示器截图 | ✅ | ✅ | 按显示器独立截取，支持高 DPI |
| **窗口截图** | ✅ | ✅ | **独立于可见性和遮挡状态** |
| 多显示器支持 | ✅ | ✅ | |
| Monitor.IsPrimary | ✅ | ✅ | 是否主显示器 |
| Monitor.ScaleFactor | ✅ | ✅ | Retina/HiDPI 缩放比例 |
| Monitor.Rotation | ✅ | ✅ | 屏幕旋转角度 |
| Monitor.Frequency | ✅ | ✅ | 刷新率（Hz）|
| Window.IsFocused | ✅ | ✅ | 是否获得焦点 |
| Window.IsMinimized | ❌ | ✅ | macOS 返回 `ErrNotSupported` |
| Window.IsMaximized | ❌ | ✅ | macOS 返回 `ErrNotSupported` |
| 排除当前进程窗口 | ✅ | ✅ | 过滤自身窗口 |
| 区域截图 | ❌ | ❌ | 请使用其他库 |

## 安装

```bash
go get github.com/zn-chen/xcap
```

## 快速开始

### 窗口截图（核心功能）

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
    // 获取所有窗口 - 每个都是独立可截取的实体
    windows, err := xcap.AllWindows()
    if err != nil {
        log.Fatal(err)
    }

    for _, w := range windows {
        // 丰富的窗口元数据
        fmt.Printf("[%s] %s (PID: %d, %dx%d at %d,%d)\n",
            w.AppName(), w.Title(), w.PID(),
            w.Width(), w.Height(), w.X(), w.Y())

        // 跳过小窗口（系统 UI 元素）
        if w.Width() < 200 || w.Height() < 200 {
            continue
        }

        // 截取窗口 - 即使被遮挡或在后台也能正常工作！
        img, err := w.CaptureImage()
        if err != nil {
            continue
        }

        filename := fmt.Sprintf("%s_%d.png", xcap.SanitizeFilename(w.AppName()), w.ID())
        f, _ := os.Create(filename)
        png.Encode(f, img)
        f.Close()
    }
}
```

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

## CLI 工具

```bash
# 构建
make build

# 截取所有显示器和窗口到 ./output/
./bin/xcap

# 只截取显示器
./bin/xcap --disable_windows

# 只截取窗口
./bin/xcap --disable_monitor
```

## API 参考

### Window 接口

```go
type Window interface {
    // 标识
    ID() uint32              // 操作系统窗口句柄（HWND/CGWindowID）
    PID() uint32             // 所属进程 ID
    AppName() string         // 应用程序名称
    Title() string           // 窗口标题

    // 几何信息
    X() int                  // X 坐标
    Y() int                  // Y 坐标
    Z() int                  // Z 顺序（值越大越靠前）
    Width() uint32           // 宽度（像素）
    Height() uint32          // 高度（像素）

    // 状态（不支持时返回 ErrNotSupported）
    IsMinimized() (bool, error)
    IsMaximized() (bool, error)
    IsFocused() (bool, error)

    // 截图
    CurrentMonitor() (Monitor, error)
    CaptureImage() (*image.RGBA, error)  // 截取窗口内容
}
```

### Monitor 接口

```go
type Monitor interface {
    ID() uint32              // 显示器 ID
    Name() string            // 显示器名称
    X() int                  // 虚拟屏幕中的 X 坐标
    Y() int                  // 虚拟屏幕中的 Y 坐标
    Width() uint32           // 宽度（像素）
    Height() uint32          // 高度（像素）
    Rotation() float64       // 旋转角度
    ScaleFactor() float64    // DPI 缩放（Retina 为 2.0）
    Frequency() float64      // 刷新率（Hz）
    IsPrimary() bool         // 是否主显示器
    IsBuiltin() bool         // 是否内置显示器
    CaptureImage() (*image.RGBA, error)
}
```

### 函数

```go
func AllMonitors() ([]Monitor, error)
func AllWindows() ([]Window, error)
func AllWindowsWithOptions(excludeCurrentProcess bool) ([]Window, error)
func SanitizeFilename(name string) string
```

## 工作原理

与简单读取屏幕坐标像素的区域截图不同，xcap 使用**操作系统级别的窗口合成 API**：

| 平台 | API | 能力 |
|------|-----|------|
| macOS | `CGWindowListCreateImage` | 直接捕获窗口的离屏缓冲区 |
| Windows | `PrintWindow` / `BitBlt` | 从合成器捕获窗口内容 |

这意味着每个窗口都作为独立实体被截取，拥有自己的位图，与屏幕上的可见状态无关。

## 项目结构

```
xcap/
├── cmd/xcap/           # 命令行工具
├── pkg/xcap/           # 公共 API（跨平台接口）
├── internal/
│   ├── darwin/         # macOS: CoreGraphics + AppKit (CGO)
│   └── windows/        # Windows: GDI + Win32 (CGO)
├── examples/           # 使用示例
└── docs/               # 实现文档
```

## 平台要求

### macOS

- macOS 10.15 (Catalina) 或更高版本
- **需要屏幕录制权限**
  - 系统设置 > 隐私与安全性 > 屏幕录制
  - 将您的应用程序添加到允许列表
- Xcode Command Line Tools: `xcode-select --install`

### Windows

- Windows 8.1 或更高版本
- MinGW-w64（用于 CGO）：通过 [MSYS2](https://www.msys2.org/) 或 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) 安装
- 无需额外权限

## 文档

- [macOS 实现原理](docs/macos-implementation.md) - CoreGraphics API 详解
- [Windows 实现原理](docs/windows-implementation.md) - GDI/Win32 API 详解
- [架构设计](docs/architecture.md) - 设计概述

## 贡献

欢迎贡献代码！请参阅 [CONTRIBUTING.md](CONTRIBUTING.md) 了解贡献指南。

## 许可证

Apache-2.0
