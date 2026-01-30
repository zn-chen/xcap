# xcap

[中文](README_CN.md) | English

A cross-platform screen and window capture library for Go, inspired by [xcap](https://github.com/nashaofu/xcap).

## Features

| Feature | macOS | Windows | Notes |
|---------|:-----:|:-------:|-------|
| Monitor capture | ✅ | ✅ | High-DPI aware |
| Window capture | ✅ | ✅ | Supports background windows |
| Multi-monitor support | ✅ | ✅ | |
| Monitor.IsPrimary | ✅ | ✅ | |
| Monitor.ScaleFactor | ✅ | ✅ | Returns display scaling (e.g., 2.0 for Retina) |
| Monitor.Rotation | ✅ | ✅ | |
| Monitor.Frequency | ✅ | ✅ | Refresh rate in Hz |
| Window.IsFocused | ✅ | ✅ | |
| Window.IsMinimized | ❌ | ✅ | macOS returns `ErrNotSupported` |
| Window.IsMaximized | ❌ | ✅ | macOS returns `ErrNotSupported` |
| Exclude current process | ✅ | ✅ | Filter out self windows |
| Region capture | ❌ | ❌ | Planned |

## Installation

```bash
go get github.com/zn-chen/xcap
```

## Quick Start

### Monitor Capture

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
        fmt.Printf("Monitor %d: %s (%dx%d @ %.0fHz, scale=%.1fx, primary=%v)\n",
            i, m.Name(), m.Width(), m.Height(), m.Frequency(), m.ScaleFactor(), m.IsPrimary())

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
    "image/png"
    "log"
    "os"

    "github.com/zn-chen/xcap/pkg/xcap"
)

func main() {
    // Use AllWindowsWithOptions(true) to exclude current process windows
    windows, err := xcap.AllWindows()
    if err != nil {
        log.Fatal(err)
    }

    for i, w := range windows {
        // Skip minimized windows (check error for unsupported platforms)
        if minimized, err := w.IsMinimized(); err == nil && minimized {
            continue
        }

        fmt.Printf("Window %d: [%s] %s (%dx%d)\n",
            i, w.AppName(), w.Title(), w.Width(), w.Height())

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

## API Reference

### Monitor Interface

```go
type Monitor interface {
    ID() uint32              // Unique identifier
    Name() string            // Display name
    X() int                  // X position
    Y() int                  // Y position
    Width() uint32           // Width in pixels
    Height() uint32          // Height in pixels
    Rotation() float64       // Rotation in degrees
    ScaleFactor() float64    // DPI scaling factor
    Frequency() float64      // Refresh rate in Hz
    IsPrimary() bool         // Is primary display
    IsBuiltin() bool         // Is built-in display (laptop)
    CaptureImage() (*image.RGBA, error)
}
```

### Window Interface

```go
type Window interface {
    ID() uint32              // Unique identifier (HWND on Windows, CGWindowID on macOS)
    PID() uint32             // Process ID
    AppName() string         // Application name
    Title() string           // Window title
    X() int                  // X position
    Y() int                  // Y position
    Z() int                  // Z-order (higher = front)
    Width() uint32           // Width in pixels
    Height() uint32          // Height in pixels
    IsMinimized() (bool, error)  // Returns ErrNotSupported if unavailable
    IsMaximized() (bool, error)  // Returns ErrNotSupported if unavailable
    IsFocused() (bool, error)    // Has input focus
    CurrentMonitor() (Monitor, error)
    CaptureImage() (*image.RGBA, error)
}
```

### Functions

```go
// Get all monitors
func AllMonitors() ([]Monitor, error)

// Get all visible windows
func AllWindows() ([]Window, error)

// Get all visible windows with options
// excludeCurrentProcess: filter out windows from current process
func AllWindowsWithOptions(excludeCurrentProcess bool) ([]Window, error)

// Sanitize filename for cross-platform compatibility
func SanitizeFilename(name string) string
```

## Project Structure

```
xcap/
├── cmd/xcap/           # CLI tool
├── pkg/xcap/           # Public API (platform-agnostic interfaces)
├── internal/
│   ├── darwin/         # macOS implementation (CGO + Objective-C)
│   └── windows/        # Windows implementation (CGO + C)
├── examples/           # Usage examples
└── docs/               # Documentation
```

## Platform Requirements

### macOS

- macOS 10.15 (Catalina) or later
- Screen Recording permission required
  - System Settings > Privacy & Security > Screen Recording
  - Add your application to the allowed list
- Xcode Command Line Tools (for CGO compilation)
  ```bash
  xcode-select --install
  ```

### Windows

- Windows 8.1 or later
- MinGW-w64 (for CGO compilation)
  - Install via [MSYS2](https://www.msys2.org/) or [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
- No additional permissions required

## Documentation

- [macOS Implementation](docs/macos-implementation.md) - CoreGraphics API details
- [Windows Implementation](docs/windows-implementation.md) - GDI/Win32 API details
- [Architecture](docs/architecture.md) - Design overview

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Apache-2.0
