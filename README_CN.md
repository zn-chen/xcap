# xcap

中文 | [English](README.md)

跨平台屏幕和窗口截图 Go 语言库，参考 [xcap](https://github.com/nashaofu/xcap) 实现。

## 功能特性

- 跨平台支持：macOS 和 Windows
- 屏幕截图：支持多显示器
- 窗口截图：支持捕获单个应用窗口
- 区域截图：支持截取屏幕指定区域

## 安装

```bash
go get github.com/zn-chen/xcap
```

## 快速开始

### 屏幕截图

```go
package main

import (
    "fmt"
    "log"
    "image/png"
    "os"

    "github.com/zn-chen/xcap/pkg/xcap"
)

func main() {
    monitors, err := xcap.AllMonitors()
    if err != nil {
        log.Fatal(err)
    }

    for i, m := range monitors {
        img, err := m.CaptureImage()
        if err != nil {
            log.Printf("Failed to capture monitor %d: %v", i, err)
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
    "log"
    "image/png"
    "os"

    "github.com/zn-chen/xcap/pkg/xcap"
)

func main() {
    windows, err := xcap.AllWindows()
    if err != nil {
        log.Fatal(err)
    }

    for i, w := range windows {
        if w.IsMinimized() {
            continue
        }

        log.Printf("Window: %s (%s)", w.Title(), w.AppName())

        img, err := w.CaptureImage()
        if err != nil {
            log.Printf("Failed to capture window: %v", err)
            continue
        }

        f, _ := os.Create(fmt.Sprintf("window-%d.png", i))
        png.Encode(f, img)
        f.Close()
    }
}
```

## 平台要求

### macOS

- macOS 10.15+
- 需要 Screen Recording 权限（系统设置 > 隐私与安全 > 屏幕录制）

### Windows

- Windows 8.1+
- 无额外权限要求

## 文档

- [macOS 实现原理](docs/macos-implementation.md)
- [Windows 实现原理](docs/windows-implementation.md)
- [架构设计](docs/architecture.md)

## 许可证

Apache-2.0
