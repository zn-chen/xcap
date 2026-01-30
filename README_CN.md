# xcap

中文 | [English](README.md)

跨平台屏幕和窗口截图 Go 语言库，参考 [xcap](https://github.com/nashaofu/xcap) 实现。

## 功能特性

| 功能 | macOS | Windows |
|------|-------|---------|
| 显示器截图 | ✅ | ✅ |
| 窗口截图 | ✅ | ✅ |
| 多显示器支持 | ✅ | ✅ |
| 显示器 IsPrimary | ✅ | ✅ |
| 显示器 ScaleFactor | ✅ | ✅ |
| 窗口 IsFocused | ✅ | ✅ |
| 窗口 IsMinimized | ❌ | ✅ |
| 窗口 IsMaximized | ❌ | ✅ |
| 排除当前进程窗口 | ✅ | ✅ |
| 区域截图 | ❌ | ❌ |

## 安装

```bash
go get github.com/zn-chen/xcap
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
- 需要 Xcode Command Line Tools（用于 CGO）

### Windows

- Windows 8.1+
- 需要 MinGW-w64（用于 CGO）
- 无额外权限要求

## 文档

- [macOS 实现原理](docs/macos-implementation.md)
- [Windows 实现原理](docs/windows-implementation.md)
- [架构设计](docs/architecture.md)

## 许可证

Apache-2.0
