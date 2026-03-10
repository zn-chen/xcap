# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

```bash
# Build CLI tool (outputs to bin/xcap, version from git tags)
make build

# Run all tests (requires display access — tests capture real screens/windows)
go test ./...

# Run tests for a specific package
go test ./internal/darwin/...
go test ./internal/windows/...

# Run a single test
go test -run TestCaptureMonitor ./internal/darwin/

# Clean build artifacts
make clean

# Run CLI: captures all monitors + windows to output/
./bin/xcap
./bin/xcap --disable_windows   # monitors only
./bin/xcap --disable_monitor   # windows only
```

Note: `make build` handles macOS `duplicate libraries` warning automatically. When using `go build` directly on macOS, prefix with:
```bash
CGO_LDFLAGS="-Wl,-no_warn_duplicate_libraries" go build ./...
```

## Architecture Overview

xcap is a cross-platform screen capture library (Go 1.22+) using CGO to bridge Go with native OS compositing APIs. It captures windows as isolated entities via their off-screen buffers, not by reading screen pixels.

### Data Flow

```
User code → pkg/xcap (interfaces) → internal/{darwin,windows}
              ↓                         ↓
         windowWrapper              bridge.go (CGO bindings)
         adapters                       ↓
                                    bridge.h + bridge.m/.c (native C/ObjC)
                                        ↓
                                    CaptureResult (raw BGRA)
                                        ↓
                                    capture.go: BGRAToRGBA → *image.RGBA
```

Public entry points: `xcap.AllMonitors()`, `xcap.AllWindows()`, `xcap.AllWindowsWithOptions(excludeCurrentProcess bool)`

### CGO Bridge Pattern

Each platform uses a three-file bridge:
- `bridge.h` - C struct definitions, function declarations, and error codes (`XCAP_OK`, `XCAP_ERR_*`)
- `bridge.m` (macOS) / `bridge.c` (Windows) - Native implementation
- `bridge.go` - CGO bindings that convert C types to Go types

Memory management: C functions allocate data → Go copies to managed memory via `unsafe.Slice` → defers C free functions (`xcap_free_monitors`, `xcap_free_windows`, `xcap_free_capture_result`).

### Platform Wrappers

`pkg/xcap/xcap_darwin.go` and `xcap_windows.go` contain `windowWrapper` adapters that wrap platform-specific implementations to satisfy the public `Monitor`/`Window` interfaces. `xcap_stub.go` returns `ErrNotSupported` on unsupported platforms.

## Code Conventions

- **Comments in Chinese**, technical terms in English
- **Logs in English** for internationalization
- Use sentinel errors from `pkg/xcap/errors.go`: `ErrNoMonitor`, `ErrNoWindow`, `ErrCaptureFailed`, `ErrPermissionDenied`, `ErrWindowMinimized`, `ErrInvalidRegion`, `ErrNotSupported`
- Methods returning platform-specific results use `(T, error)` where unsupported platforms return `ErrNotSupported`
- Version injected via ldflags: `-X main.version=$(VERSION)` (git describe or "dev")

## Platform Notes

### macOS
- Minimum deployment target: macOS 10.15 (set in CGO CFLAGS)
- Requires Screen Recording permission (System Settings > Privacy & Security)
- CGO flags: `-framework CoreGraphics -framework AppKit -framework CoreFoundation`
- `IsMinimized`/`IsMaximized` return `ErrNotSupported` (would require Accessibility API)
- Unimplemented: `Window.Z()`, `Window.CurrentMonitor()`, `Monitor.Rotation()`, `Monitor.Frequency()`, `Monitor.IsBuiltin()`, `CaptureRegion`

### Windows
- Uses GDI for capture, Win32 for window enumeration
- CGO flags: `-luser32 -lgdi32 -ldwmapi -lshcore -lpsapi`
- CGO requires MinGW-w64
- All window state methods (`IsMinimized`, `IsMaximized`, `IsFocused`) fully supported
- Unimplemented: `Window.Z()`, `Window.CurrentMonitor()`, `Monitor.Rotation()`, `Monitor.Frequency()`, `Monitor.IsBuiltin()`, `CaptureRegion`

## Testing Notes

- Tests require an active display — they enumerate real monitors/windows and perform actual captures
- Tests skip gracefully with `t.Skip()` when no suitable windows are found
- Capture tests write output files for manual inspection (e.g., `/tmp/xcap_window_test.png` on macOS)
- Only `internal/{darwin,windows}` have tests; `pkg/xcap` has no tests (thin wrapper layer)
