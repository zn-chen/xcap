# xcap

[中文](README_CN.md) | English

A cross-platform screen and window capture library for Go, inspired by [xcap](https://github.com/nashaofu/xcap).

## Features

| Feature | macOS | Windows |
|---------|-------|---------|
| Monitor capture | ✅ | ✅ |
| Window capture | ✅ | ✅ |
| Multi-monitor support | ✅ | ✅ |
| Monitor IsPrimary | ✅ | ✅ |
| Monitor ScaleFactor | ✅ | ✅ |
| Window IsFocused | ✅ | ✅ |
| Window IsMinimized | ❌ | ✅ |
| Window IsMaximized | ❌ | ✅ |
| Exclude current process | ✅ | ✅ |
| Region capture | ❌ | ❌ |

## Installation

```bash
go get github.com/zn-chen/xcap
```

## CLI Tool

```bash
# Build
make build

# Run (captures all monitors and windows to ./output/)
./bin/xcap

# Capture monitors only
./bin/xcap --disable_windows

# Capture windows only
./bin/xcap --disable_monitor
```

## Quick Start

### Screen Capture

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

### Window Capture

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

## Platform Requirements

### macOS

- macOS 10.15+
- Screen Recording permission required (System Settings > Privacy & Security > Screen Recording)
- Requires Xcode Command Line Tools (for CGO)

### Windows

- Windows 8.1+
- Requires MinGW-w64 (for CGO)
- No additional permissions required

## Documentation

- [macOS Implementation](docs/macos-implementation.md)
- [Windows Implementation](docs/windows-implementation.md)
- [Architecture](docs/architecture.md)

## License

Apache-2.0
