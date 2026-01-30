# xcap

[中文](README_CN.md) | English

A cross-platform **native window capture** library for Go, inspired by [xcap](https://github.com/nashaofu/xcap).

## Why xcap?

Most Go screenshot libraries only support **region capture** - they grab pixels from a rectangular area of the screen. This approach has fundamental limitations:

| Approach | Region Capture | xcap (Native Window Capture) |
|----------|:-------------:|:----------------------------:|
| Capture overlapped windows | ❌ | ✅ |
| Capture background windows | ❌ | ✅ |
| Capture minimized windows | ❌ | ✅ (Windows) |
| Per-window metadata (title, app, PID) | ❌ | ✅ |
| Individual window isolation | ❌ | ✅ |

**xcap uses native OS APIs** (CoreGraphics on macOS, GDI/Win32 on Windows) to capture windows as independent entities, not just screen regions. This enables:

- **Capture any window** regardless of visibility or overlap
- **Enumerate all windows** with metadata (title, app name, process ID, position, size)
- **Capture specific monitors** in multi-display setups
- **High-DPI aware** capture with proper scaling

## Features

| Feature | macOS | Windows | Notes |
|---------|:-----:|:-------:|-------|
| Monitor capture | ✅ | ✅ | Per-display, high-DPI aware |
| **Window capture** | ✅ | ✅ | **Independent of visibility/overlap** |
| Multi-monitor support | ✅ | ✅ | |
| Monitor.IsPrimary | ✅ | ✅ | |
| Monitor.ScaleFactor | ✅ | ✅ | Retina/HiDPI scaling factor |
| Monitor.Rotation | ✅ | ✅ | |
| Monitor.Frequency | ✅ | ✅ | Refresh rate in Hz |
| Window.IsFocused | ✅ | ✅ | |
| Window.IsMinimized | ❌ | ✅ | macOS returns `ErrNotSupported` |
| Window.IsMaximized | ❌ | ✅ | macOS returns `ErrNotSupported` |
| Exclude current process | ✅ | ✅ | Filter out self windows |
| Region capture | ❌ | ❌ | Use other libraries for this |

## Installation

```bash
go get github.com/zn-chen/xcap
```

## Quick Start

### Window Capture (The Key Feature)

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
    // Get all windows - each is an independent capturable entity
    windows, err := xcap.AllWindows()
    if err != nil {
        log.Fatal(err)
    }

    for _, w := range windows {
        // Rich metadata for each window
        fmt.Printf("[%s] %s (PID: %d, %dx%d at %d,%d)\n",
            w.AppName(), w.Title(), w.PID(),
            w.Width(), w.Height(), w.X(), w.Y())

        // Skip small windows (UI elements)
        if w.Width() < 200 || w.Height() < 200 {
            continue
        }

        // Capture window - works even if overlapped or in background!
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

## CLI Tool

```bash
# Build
make build

# Capture all monitors and windows to ./output/
./bin/xcap

# Capture monitors only
./bin/xcap --disable_windows

# Capture windows only
./bin/xcap --disable_monitor
```

## API Reference

### Window Interface

```go
type Window interface {
    // Identity
    ID() uint32              // OS window handle (HWND/CGWindowID)
    PID() uint32             // Owner process ID
    AppName() string         // Application name
    Title() string           // Window title

    // Geometry
    X() int                  // X position
    Y() int                  // Y position
    Z() int                  // Z-order (higher = front)
    Width() uint32           // Width in pixels
    Height() uint32          // Height in pixels

    // State (returns ErrNotSupported if unavailable)
    IsMinimized() (bool, error)
    IsMaximized() (bool, error)
    IsFocused() (bool, error)

    // Capture
    CurrentMonitor() (Monitor, error)
    CaptureImage() (*image.RGBA, error)  // Capture window content
}
```

### Monitor Interface

```go
type Monitor interface {
    ID() uint32              // Display ID
    Name() string            // Display name
    X() int                  // X position in virtual screen
    Y() int                  // Y position in virtual screen
    Width() uint32           // Width in pixels
    Height() uint32          // Height in pixels
    Rotation() float64       // Rotation in degrees
    ScaleFactor() float64    // DPI scaling (2.0 for Retina)
    Frequency() float64      // Refresh rate in Hz
    IsPrimary() bool         // Is primary display
    IsBuiltin() bool         // Is built-in display
    CaptureImage() (*image.RGBA, error)
}
```

### Functions

```go
func AllMonitors() ([]Monitor, error)
func AllWindows() ([]Window, error)
func AllWindowsWithOptions(excludeCurrentProcess bool) ([]Window, error)
func SanitizeFilename(name string) string
```

## How It Works

Unlike region-based capture that simply reads pixels from screen coordinates, xcap uses **OS-level window compositing APIs**:

| Platform | API | Capability |
|----------|-----|------------|
| macOS | `CGWindowListCreateImage` | Captures window's off-screen buffer directly |
| Windows | `PrintWindow` / `BitBlt` | Captures window content from compositor |

This means each window is captured as an isolated entity with its own bitmap, independent of what's visible on screen.

## Project Structure

```
xcap/
├── cmd/xcap/           # CLI tool
├── pkg/xcap/           # Public API (cross-platform interfaces)
├── internal/
│   ├── darwin/         # macOS: CoreGraphics + AppKit via CGO
│   └── windows/        # Windows: GDI + Win32 via CGO
├── examples/           # Usage examples
└── docs/               # Implementation documentation
```

## Platform Requirements

### macOS

- macOS 10.15 (Catalina) or later
- **Screen Recording permission required**
  - System Settings > Privacy & Security > Screen Recording
  - Add your application to the allowed list
- Xcode Command Line Tools: `xcode-select --install`

### Windows

- Windows 8.1 or later
- MinGW-w64 for CGO: Install via [MSYS2](https://www.msys2.org/) or [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
- No additional permissions required

## Documentation

- [macOS Implementation](docs/macos-implementation.md) - CoreGraphics API internals
- [Windows Implementation](docs/windows-implementation.md) - GDI/Win32 API internals
- [Architecture](docs/architecture.md) - Design overview

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Apache-2.0
