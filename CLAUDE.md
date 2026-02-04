# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

```bash
# Build CLI tool
make build

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/darwin/...
go test ./internal/windows/...

# Run a single test
go test -run TestCaptureMonitor ./internal/darwin/

# Clean build artifacts
make clean
```

## Architecture Overview

xcap is a cross-platform screen capture library using CGO to bridge Go with native APIs.

### Layer Structure

```
pkg/xcap/           → Public API (interfaces + platform wrappers)
    ├── monitor.go, window.go    → Interface definitions
    ├── xcap_darwin.go           → macOS wrapper (wraps internal/darwin)
    ├── xcap_windows.go          → Windows wrapper (wraps internal/windows)
    └── xcap_stub.go             → Unsupported platforms

internal/darwin/    → macOS implementation
    ├── bridge.h, bridge.m       → C/Objective-C layer (CoreGraphics, AppKit)
    ├── bridge.go                → CGO bindings
    ├── monitor.go, window.go    → Go types implementing interfaces
    └── capture.go               → Image conversion (BGRA → RGBA)

internal/windows/   → Windows implementation
    ├── bridge.h, bridge.c       → C layer (GDI, Win32 API)
    ├── bridge.go                → CGO bindings
    ├── monitor.go, window.go    → Go types implementing interfaces
    └── capture.go               → Image capture logic
```

### CGO Bridge Pattern

Each platform uses a three-file bridge:
- `bridge.h` - C struct definitions and function declarations
- `bridge.m` (macOS) / `bridge.c` (Windows) - Native implementation
- `bridge.go` - CGO bindings that convert C types to Go types

The bridge layer handles memory management: Go code calls C functions that allocate data, copies it to Go-managed memory, then calls free functions.

### Platform Wrappers

`pkg/xcap/xcap_darwin.go` and `xcap_windows.go` contain `windowWrapper` and (implicitly for monitors) adapters that wrap platform-specific implementations to satisfy the public interfaces in `pkg/xcap/`.

## Code Conventions

- **Comments in Chinese**, technical terms in English
- **Logs in English** for internationalization
- Use errors from `pkg/xcap/errors.go` (e.g., `ErrNotSupported`, `ErrCaptureFailed`)
- Methods returning platform-specific results use `(T, error)` where unsupported platforms return `ErrNotSupported`

## Platform Notes

### macOS
- Requires Screen Recording permission
- CGO flags: `-framework CoreGraphics -framework AppKit -framework CoreFoundation`
- `IsMinimized`/`IsMaximized` return `ErrNotSupported` (require Accessibility API)
- To suppress `duplicate libraries` warning when using `go build` directly:
  ```bash
  CGO_LDFLAGS="-Wl,-no_warn_duplicate_libraries" go build ./...
  ```

### Windows
- Uses GDI for capture, Win32 for window enumeration
- CGO requires MinGW-w64
- All window state methods fully supported
