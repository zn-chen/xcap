# macOS 窗口截图实现原理

本文档详细介绍 xcap 在 macOS 平台上实现屏幕和窗口截图的技术原理。

## 核心框架

macOS 截图主要依赖以下系统框架：

- **Core Graphics (Quartz)** - 底层图形框架，提供屏幕和窗口捕获 API
- **AppKit** - 提供窗口管理和应用程序信息
- **Core Foundation** - 提供基础数据类型和集合操作

## 关键 API

### 1. 屏幕/显示器相关

#### CGGetActiveDisplayList
```c
CGError CGGetActiveDisplayList(
    uint32_t maxDisplays,
    CGDirectDisplayID *activeDisplays,
    uint32_t *displayCount
);
```
获取当前活动的所有显示器列表。

#### CGDisplayBounds
```c
CGRect CGDisplayBounds(CGDirectDisplayID display);
```
获取指定显示器的边界矩形（位置和尺寸）。

#### CGDisplayCopyDisplayMode
获取显示器的显示模式，用于计算缩放因子（scale factor）。

### 2. 窗口相关

#### CGWindowListCopyWindowInfo
```c
CFArrayRef CGWindowListCopyWindowInfo(
    CGWindowListOption option,
    CGWindowID relativeToWindow
);
```
获取窗口列表信息，返回按 Z 顺序排列的窗口数组（从顶层到最底层）。

**常用选项：**
- `kCGWindowListOptionOnScreenOnly` - 仅获取屏幕上可见的窗口
- `kCGWindowListExcludeDesktopElements` - 排除桌面元素

**窗口信息字典包含：**
- `kCGWindowNumber` - 窗口 ID
- `kCGWindowOwnerPID` - 拥有者进程 ID
- `kCGWindowOwnerName` - 应用程序名称
- `kCGWindowName` - 窗口标题
- `kCGWindowBounds` - 窗口边界矩形
- `kCGWindowSharingState` - 窗口共享状态
- `kCGWindowIsOnscreen` - 是否在屏幕上

### 3. 截图核心 API

#### CGWindowListCreateImage
```c
CGImageRef CGWindowListCreateImage(
    CGRect screenBounds,
    CGWindowListOption listOption,
    CGWindowID windowID,
    CGWindowImageOption imageOption
);
```

这是实现截图的核心函数。

**参数说明：**
- `screenBounds` - 截图区域（CGRect）
- `listOption` - 窗口列表选项
  - `kCGWindowListOptionAll` - 所有窗口（用于屏幕截图）
  - `kCGWindowListOptionIncludingWindow` - 包含指定窗口（用于窗口截图）
- `windowID` - 目标窗口 ID（屏幕截图时为 0）
- `imageOption` - 图像选项（通常使用 `kCGWindowImageDefault`）

## 实现流程

### 屏幕截图流程

```
1. CGGetActiveDisplayList() 获取所有活动显示器
2. 选择目标显示器
3. CGDisplayBounds() 获取显示器边界
4. CGWindowListCreateImage(bounds, kCGWindowListOptionAll, 0, default)
5. 从 CGImage 提取像素数据
6. BGRA → RGBA 颜色转换
7. 返回 RgbaImage
```

### 窗口截图流程

```
1. CGWindowListCopyWindowInfo() 获取窗口列表
2. 遍历查找目标窗口 ID
3. 获取窗口边界 (kCGWindowBounds)
4. CGWindowListCreateImage(bounds, kCGWindowListOptionIncludingWindow, windowID, default)
5. 从 CGImage 提取像素数据
6. BGRA → RGBA 颜色转换
7. 返回 RgbaImage
```

### 图像数据处理

```rust
// 从 CGImage 获取数据
let data_provider = CGImage::data_provider(cg_image);
let data = CGDataProvider::data(data_provider).to_vec();

// 处理行对齐（macOS 可能在每行末尾有额外字节）
let bytes_per_row = CGImage::bytes_per_row(cg_image);
for row in data.chunks_exact(bytes_per_row) {
    buffer.extend_from_slice(&row[..width * 4]);
}

// BGRA → RGBA 转换
for bgra in buffer.chunks_exact_mut(4) {
    bgra.swap(0, 2);  // 交换 B 和 R
}
```

## 权限要求

### Screen Recording 权限

macOS 10.15+ 需要 Screen Recording 权限才能捕获屏幕和窗口内容。

```rust
// 检查权限
if !CGPreflightScreenCaptureAccess() {
    // 提示用户在 系统设置 > 隐私与安全 > 屏幕录制 中授权
}
```

## 窗口过滤规则

xcap 会过滤掉以下窗口：

1. **StatusIndicator** - 系统状态指示器窗口
2. **窗口共享状态为 0** - 不可共享的窗口
3. **不在屏幕上的窗口** - 最小化的窗口

## 显示器信息获取

### 获取显示器名称
通过 NSScreen API 获取友好名称：

```rust
let screens = NSScreen::screens();
for screen in screens {
    let device_description = screen.deviceDescription();
    let screen_id = device_description["NSScreenNumber"];
    let name = screen.localizedName();
}
```

### 显示器属性

| 属性 | API |
|------|-----|
| ID | CGDirectDisplayID |
| 位置 | CGDisplayBounds().origin |
| 尺寸 | CGDisplayBounds().size |
| 旋转角度 | CGDisplayRotation() |
| 缩放因子 | pixel_width / logical_width |
| 刷新率 | CGDisplayMode::refresh_rate() |
| 是否主显示器 | CGDisplayIsMain() |
| 是否内置显示器 | CGDisplayIsBuiltin() |

## Go 语言实现建议

在 Go 中实现 macOS 截图，可以使用以下方式：

### 方案一：CGO + Objective-C

```go
// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework CoreGraphics -framework AppKit -framework CoreFoundation
// #include <CoreGraphics/CoreGraphics.h>
import "C"
```

### 方案二：使用现有 Go 绑定库

- `github.com/progrium/darwinkit` - macOS 系统框架的 Go 绑定
- `github.com/kbinani/screenshot` - 简单的跨平台截图库（可参考实现）

### 方案三：调用系统命令

```go
// 使用 screencapture 命令（简单但功能有限）
exec.Command("screencapture", "-x", "-t", "png", filename)
```

## 注意事项

1. **行对齐** - macOS 返回的图像数据每行可能有额外的填充字节
2. **颜色格式** - 原始数据是 BGRA 格式，需要转换为 RGBA
3. **Retina 显示器** - 需要考虑缩放因子，实际像素尺寸 = 逻辑尺寸 × 缩放因子
4. **权限检查** - 10.15+ 必须先检查并请求 Screen Recording 权限
5. **焦点窗口判断** - 使用 NSWorkspace.activeApplication 获取当前焦点应用
