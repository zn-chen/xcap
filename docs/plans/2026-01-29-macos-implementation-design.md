# macOS 实现设计文档

**日期**: 2026-01-29
**状态**: 已确认

## 概述

实现 owl-go 的 macOS 平台支持，包括显示器枚举、窗口枚举和截图功能。

## 设计决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| CGO 绑定方式 | Objective-C 混合 | 需要 AppKit API（如 NSScreen）获取完整信息 |
| 实现范围 | 最小可用版本 | 先跑通核心流程，再逐步添加功能 |
| 代码组织 | 分离方式 (.go + .h + .m) | Objective-C 代码量较大，分离更清晰 |
| 公共 API 入口 | 条件编译 (build tags) | Go 标准做法，编译产物干净 |

## 实现范围（最小可用版本）

### Monitor 功能
- [x] 枚举所有显示器 (`AllMonitors`)
- [x] 获取 ID
- [x] 获取名称
- [x] 获取位置 (X, Y)
- [x] 获取尺寸 (Width, Height)
- [x] 屏幕截图 (`CaptureImage`)

### Window 功能
- [x] 枚举所有窗口 (`AllWindows`)
- [x] 获取 ID
- [x] 获取 PID
- [x] 获取应用名称
- [x] 获取窗口标题
- [x] 获取位置 (X, Y)
- [x] 获取尺寸 (Width, Height)
- [x] 窗口截图 (`CaptureImage`)

### 暂不实现
- 旋转角度、缩放因子、刷新率
- 是否主显示器、是否内置显示器
- Z 顺序、最小化/最大化/焦点状态
- 区域截图 (`CaptureRegion`)
- `CurrentMonitor()`
- 权限检查 API

## 文件结构

```
owl-go/
├── pkg/owl/
│   ├── monitor.go              # Monitor 接口定义
│   ├── window.go               # Window 接口定义
│   ├── errors.go               # 错误定义
│   ├── owl_darwin.go           # //go:build darwin
│   │                           # 实现 AllMonitors(), AllWindows()
│   └── owl_stub.go             # //go:build !darwin && !windows
│                               # 返回 ErrNotSupported
│
├── internal/darwin/
│   ├── bridge.h                # C 头文件
│   │                           # - 结构体定义 (MonitorInfo, WindowInfo)
│   │                           # - 函数声明
│   │
│   ├── bridge.m                # Objective-C 实现
│   │                           # - 调用 Core Graphics API
│   │                           # - 调用 AppKit API
│   │
│   ├── bridge.go               # CGO 桥接层
│   │                           # - #cgo CFLAGS/LDFLAGS
│   │                           # - Go 调用 C 函数的封装
│   │
│   ├── monitor.go              # darwinMonitor 结构体
│   │                           # - 实现 owl.Monitor 接口
│   │
│   ├── window.go               # darwinWindow 结构体
│   │                           # - 实现 owl.Window 接口
│   │
│   └── capture.go              # 图像处理
│                               # - BGRA → RGBA 转换
│                               # - CGImage → image.RGBA
```

## 核心数据流

### 显示器枚举流程

```
AllMonitors()
    │
    ▼
bridge.go: GetAllMonitors()
    │
    ▼
bridge.m: owl_get_all_monitors()
    │
    ├── CGGetActiveDisplayList()     获取显示器 ID 列表
    ├── CGDisplayBounds()            获取位置和尺寸
    └── NSScreen.localizedName       获取友好名称
    │
    ▼
返回 []MonitorInfo 结构体
    │
    ▼
monitor.go: 转换为 []owl.Monitor
```

### 窗口枚举流程

```
AllWindows()
    │
    ▼
bridge.go: GetAllWindows()
    │
    ▼
bridge.m: owl_get_all_windows()
    │
    ├── CGWindowListCopyWindowInfo()  获取窗口列表
    ├── 遍历 CFArray
    │   ├── kCGWindowNumber          窗口 ID
    │   ├── kCGWindowOwnerPID        进程 ID
    │   ├── kCGWindowOwnerName       应用名称
    │   ├── kCGWindowName            窗口标题
    │   └── kCGWindowBounds          位置和尺寸
    └── 过滤无效窗口
    │
    ▼
返回 []WindowInfo 结构体
    │
    ▼
window.go: 转换为 []owl.Window
```

### 截图流程

```
CaptureImage()
    │
    ▼
bridge.go: CaptureMonitor(id) 或 CaptureWindow(id)
    │
    ▼
bridge.m: owl_capture_monitor() 或 owl_capture_window()
    │
    ├── CGWindowListCreateImage()    核心截图函数
    ├── CGImageGetDataProvider()     获取数据提供者
    ├── CGDataProviderCopyData()     复制像素数据
    └── 返回 BGRA 数据 + 宽高
    │
    ▼
capture.go: BGRAToRGBA()
    │
    ├── 处理行对齐 (bytes_per_row)
    └── 交换 B 和 R 通道
    │
    ▼
返回 *image.RGBA
```

## C 层接口设计

### bridge.h

```c
// 显示器信息
typedef struct {
    uint32_t id;
    char name[256];
    int32_t x;
    int32_t y;
    uint32_t width;
    uint32_t height;
} OwlMonitorInfo;

// 窗口信息
typedef struct {
    uint32_t id;
    uint32_t pid;
    char app_name[256];
    char title[256];
    int32_t x;
    int32_t y;
    uint32_t width;
    uint32_t height;
} OwlWindowInfo;

// 截图结果
typedef struct {
    uint8_t *data;
    uint32_t width;
    uint32_t height;
    uint32_t bytes_per_row;
} OwlCaptureResult;

// API 函数
int owl_get_all_monitors(OwlMonitorInfo **monitors, int *count);
int owl_get_all_windows(OwlWindowInfo **windows, int *count);
int owl_capture_monitor(uint32_t display_id, OwlCaptureResult *result);
int owl_capture_window(uint32_t window_id, OwlCaptureResult *result);
void owl_free_monitors(OwlMonitorInfo *monitors);
void owl_free_windows(OwlWindowInfo *windows);
void owl_free_capture_result(OwlCaptureResult *result);
```

## 窗口过滤规则

参考 xcap 实现，过滤以下窗口：

1. `kCGWindowSharingState == 0` - 不可共享的窗口
2. `kCGWindowName == "StatusIndicator" && kCGWindowOwnerName == "Window Server"` - 系统状态指示器
3. 不在屏幕上的窗口（通过 `kCGWindowListOptionOnScreenOnly` 选项）

## 错误处理

| 错误场景 | 返回错误 |
|----------|----------|
| 无显示器 | `ErrNoMonitor` |
| 无窗口 | `ErrNoWindow` |
| CGWindowListCreateImage 返回 nil | `ErrCaptureFailed` |
| 内存分配失败 | `ErrCaptureFailed` |

## 实现计划

### Phase 1: CGO 基础设施
- [ ] 创建 `internal/darwin/bridge.h` - 定义结构体和函数声明
- [ ] 创建 `internal/darwin/bridge.m` - 桩函数实现
- [ ] 创建 `internal/darwin/bridge.go` - CGO 指令和桥接
- [ ] 验收: `go build ./internal/darwin` 编译通过

### Phase 2: Monitor 实现
- [ ] 实现 `owl_get_all_monitors()` - 显示器枚举
- [ ] 创建 `internal/darwin/monitor.go` - darwinMonitor 结构体
- [ ] 实现 `owl_capture_monitor()` - 屏幕截图
- [ ] 创建 `internal/darwin/capture.go` - BGRA→RGBA 转换
- [ ] 验收: 能枚举显示器并截图保存为 PNG

### Phase 3: Window 实现
- [ ] 实现 `owl_get_all_windows()` - 窗口枚举
- [ ] 创建 `internal/darwin/window.go` - darwinWindow 结构体
- [ ] 实现 `owl_capture_window()` - 窗口截图
- [ ] 验收: 能枚举窗口并截图保存为 PNG

### Phase 4: 公共 API 集成
- [ ] 创建 `pkg/owl/owl_darwin.go` - AllMonitors(), AllWindows()
- [ ] 创建 `pkg/owl/owl_stub.go` - 非 darwin 平台返回错误
- [ ] 创建 `examples/monitor_capture/main.go` - 示例代码
- [ ] 创建 `examples/window_capture/main.go` - 示例代码
- [ ] 验收: 用户可以 import "owl-go/pkg/owl" 使用

## 依赖

### 系统要求
- macOS 10.15+ (Catalina)
- Xcode Command Line Tools

### 链接框架
```
-framework CoreGraphics
-framework AppKit
-framework CoreFoundation
```

## 参考

- [xcap macOS 实现](https://github.com/nashaofu/xcap/tree/main/src/macos)
- [Core Graphics 文档](https://developer.apple.com/documentation/coregraphics)
- [CGWindowListCreateImage](https://developer.apple.com/documentation/coregraphics/1454
