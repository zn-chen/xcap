# Windows 窗口截图实现原理

本文档详细介绍 xcap 在 Windows 平台上实现屏幕和窗口截图的技术原理。

## 核心 API 框架

Windows 截图主要依赖以下 Win32 API：

- **GDI (Graphics Device Interface)** - 传统图形设备接口
- **DWM (Desktop Window Manager)** - 桌面窗口管理器
- **User32** - 用户界面相关 API
- **Shcore** - 高 DPI 支持

## 关键 API

### 1. 显示器枚举

#### EnumDisplayMonitors
```c
BOOL EnumDisplayMonitors(
    HDC hdc,
    LPCRECT lprcClip,
    MONITORENUMPROC lpfnEnum,
    LPARAM dwData
);
```
枚举所有显示器，通过回调函数收集 HMONITOR 句柄。

#### GetMonitorInfoW
```c
BOOL GetMonitorInfoW(
    HMONITOR hMonitor,
    LPMONITORINFO lpmi
);
```
获取显示器信息，包括工作区域和设备名称。

#### EnumDisplaySettingsW
获取显示器的显示模式（分辨率、刷新率、旋转角度等）。

### 2. 窗口枚举

#### EnumWindows
```c
BOOL EnumWindows(
    WNDENUMPROC lpEnumFunc,
    LPARAM lParam
);
```
枚举所有顶层窗口，按 Z 顺序遍历（从最顶层开始）。

#### 窗口状态检查函数

| 函数 | 用途 |
|------|------|
| IsWindow() | 检查窗口句柄是否有效 |
| IsWindowVisible() | 检查窗口是否可见 |
| IsIconic() | 检查窗口是否最小化 |
| IsZoomed() | 检查窗口是否最大化 |
| GetForegroundWindow() | 获取当前焦点窗口 |

### 3. DWM API

#### DwmGetWindowAttribute
```c
HRESULT DwmGetWindowAttribute(
    HWND hwnd,
    DWORD dwAttribute,
    PVOID pvAttribute,
    DWORD cbAttribute
);
```

**常用属性：**
- `DWMWA_CLOAKED` - 检查窗口是否被隐藏（UWP 应用相关）
- `DWMWA_EXTENDED_FRAME_BOUNDS` - 获取窗口实际边界（不含阴影）

### 4. 截图 API

#### BitBlt - 位块传输
```c
BOOL BitBlt(
    HDC hdc,        // 目标 DC
    int x, int y,   // 目标位置
    int cx, int cy, // 尺寸
    HDC hdcSrc,     // 源 DC
    int x1, int y1, // 源位置
    DWORD rop       // 光栅操作（SRCCOPY）
);
```
从源 DC 复制位图到目标 DC，用于屏幕截图。

#### PrintWindow - 窗口打印
```c
BOOL PrintWindow(
    HWND hwnd,
    HDC hdcBlt,
    UINT nFlags
);
```

**标志位：**
- `PW_CLIENTONLY (1)` - 仅客户区
- `PW_RENDERFULLCONTENT (2)` - 完整渲染内容（Windows 8+）

PrintWindow 可以捕获不在屏幕上可见的窗口内容（如被遮挡的窗口）。

## 实现流程

### 屏幕截图流程

```
1. GetDesktopWindow() 获取桌面窗口句柄
2. GetWindowDC() 获取桌面窗口的设备上下文 (DC)
3. CreateCompatibleDC() 创建内存 DC
4. CreateCompatibleBitmap() 创建兼容位图
5. SelectObject() 将位图选入内存 DC
6. BitBlt() 将屏幕内容复制到内存 DC
7. GetDIBits() 读取位图数据到缓冲区
8. BGRA → RGBA 颜色转换
9. 释放资源 (DeleteDC, DeleteObject, ReleaseDC)
```

### 窗口截图流程

```
1. 获取窗口信息 (GetWindowInfo)
2. 处理 DPI 缩放
3. GetWindowDC() 获取窗口 DC
4. CreateCompatibleDC() 创建内存 DC
5. CreateCompatibleBitmap() 创建兼容位图
6. 尝试多种捕获方法：
   a. PrintWindow(PW_RENDERFULLCONTENT) - Windows 8+ 首选
   b. PrintWindow(0) - DWM 启用时
   c. PrintWindow(PW_CLIENTONLY | PW_RENDERFULLCONTENT)
   d. BitBlt() - 后备方案
7. GetDIBits() 读取位图数据
8. 裁剪到客户区
9. BGRA → RGBA 颜色转换
```

## 窗口过滤规则

xcap 使用以下规则过滤无效窗口：

```rust
fn is_valid_window(hwnd: HWND) -> bool {
    // 1. 必须是有效且可见的窗口
    if !IsWindow(hwnd) || !IsWindowVisible(hwnd) {
        return false;
    }

    // 2. 过滤特殊窗口类
    // - "Progman" - Program Manager
    // - "Button" - 开始按钮 (Vista/7)

    // 3. 过滤 WS_EX_TOOLWINDOW 样式的窗口
    // 例外：Shell_TrayWnd (任务栏)

    // 4. 排除当前进程的窗口（避免死锁）
    if GetWindowPid(hwnd) == GetCurrentProcessId() {
        return false;
    }

    // 5. 排除 Cloaked 窗口（隐藏的 UWP 应用）
    if is_window_cloaked(hwnd) {
        return false;
    }

    // 6. 排除空矩形窗口
    if IsRectEmpty(&rect) {
        return false;
    }

    true
}
```

## DPI 处理

Windows 10+ 支持每窗口 DPI，需要特殊处理：

### 获取进程 DPI 感知状态
```rust
// 使用 Shcore.dll 的 GetProcessDpiAwareness
fn get_process_is_dpi_awareness(process: HANDLE) -> bool {
    let get_dpi = GetProcAddress(shcore, "GetProcessDpiAwareness");
    let mut awareness = 0;
    get_dpi(process, &mut awareness);
    awareness != 0  // PROCESS_DPI_UNAWARE = 0
}
```

### 获取显示器 DPI
```rust
// 使用 Shcore.dll 的 GetDpiForMonitor
fn get_hi_dpi_scale_factor(h_monitor: HMONITOR) -> f32 {
    let get_dpi = GetProcAddress(shcore, "GetDpiForMonitor");
    let mut dpi_x = 0;
    get_dpi(h_monitor, MDT_EFFECTIVE_DPI, &mut dpi_x, &mut dpi_y);
    dpi_x as f32 / 96.0  // 96 DPI = 100% 缩放
}
```

### 缩放因子计算

```rust
let scale_factor = if !window_is_dpi_awareness || current_process_is_dpi_awareness {
    1.0  // 不需要缩放
} else {
    monitor.scale_factor()  // 需要按显示器缩放
};
```

## 颜色格式转换

Windows GDI 返回 BGRA 格式数据：

```rust
fn bgra_to_rgba(buffer: &mut [u8]) {
    for chunk in buffer.chunks_exact_mut(4) {
        chunk.swap(0, 2);  // 交换 B 和 R

        // Windows 7 及更早版本的 alpha 修复
        if chunk[3] == 0 && is_old_version {
            chunk[3] = 255;
        }
    }
}
```

## 获取应用名称

Windows 获取应用名称比较复杂：

```rust
fn get_app_name(pid: u32) -> String {
    // 1. 打开进程
    let handle = OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, pid);

    // 2. 获取模块文件名
    let filename = GetModuleFileNameExW(handle);

    // 3. 获取文件版本信息
    let info_size = GetFileVersionInfoSizeW(filename);
    let info = GetFileVersionInfoW(filename);

    // 4. 按优先级查找属性
    let keys = [
        "FileDescription",   // 首选
        "ProductName",
        "ProductShortName",
        "InternalName",
        "OriginalFilename",  // 后备
    ];

    // 5. 如果都失败，使用模块基本名称
    GetModuleBaseNameW(handle)
}
```

## 显示器配置信息

使用 DisplayConfig API 获取详细信息：

```rust
// 获取显示器友好名称
fn get_monitor_config(monitor_info: MONITORINFOEXW) -> DISPLAYCONFIG_TARGET_DEVICE_NAME {
    // 1. 获取显示配置缓冲区大小
    GetDisplayConfigBufferSizes(QDC_ONLY_ACTIVE_PATHS, &paths_count, &modes_count);

    // 2. 查询显示配置
    QueryDisplayConfig(QDC_ONLY_ACTIVE_PATHS, paths, modes);

    // 3. 遍历路径匹配源设备名称
    for path in paths {
        let source = DisplayConfigGetDeviceInfo(DISPLAYCONFIG_DEVICE_INFO_GET_SOURCE_NAME);
        if source.viewGdiDeviceName == monitor_info.szDevice {
            return DisplayConfigGetDeviceInfo(DISPLAYCONFIG_DEVICE_INFO_GET_TARGET_NAME);
        }
    }
}
```

## Go 语言实现建议

### 方案一：使用 syscall/windows

```go
import (
    "syscall"
    "unsafe"
    "golang.org/x/sys/windows"
)

var (
    user32 = windows.NewLazySystemDLL("user32.dll")
    gdi32  = windows.NewLazySystemDLL("gdi32.dll")

    procEnumWindows = user32.NewProc("EnumWindows")
    procBitBlt      = gdi32.NewProc("BitBlt")
)
```

### 方案二：使用现有 Go 库

- `github.com/kbinani/screenshot` - 跨平台截图库
- `github.com/vova616/screenshot` - Windows 截图
- `golang.org/x/sys/windows` - Windows 系统调用

### 核心结构体

```go
type RECT struct {
    Left, Top, Right, Bottom int32
}

type BITMAPINFOHEADER struct {
    BiSize          uint32
    BiWidth         int32
    BiHeight        int32
    BiPlanes        uint16
    BiBitCount      uint16
    BiCompression   uint32
    BiSizeImage     uint32
    BiXPelsPerMeter int32
    BiYPelsPerMeter int32
    BiClrUsed       uint32
    BiClrImportant  uint32
}

type MONITORINFO struct {
    CbSize    uint32
    RcMonitor RECT
    RcWork    RECT
    DwFlags   uint32
}
```

## 注意事项

1. **资源释放** - 必须正确释放 DC、位图等 GDI 对象，否则会内存泄漏
2. **DPI 感知** - Windows 10+ 需要处理每窗口 DPI
3. **Cloaked 窗口** - UWP 应用可能处于 Cloaked 状态
4. **PrintWindow 降级** - 某些窗口可能不支持 PrintWindow，需要降级到 BitBlt
5. **当前进程窗口** - 不要尝试截取当前进程的窗口（可能死锁）
6. **Windows 版本** - 需要检查 Windows 版本以使用正确的 API
   - Windows 8+ 使用 PrintWindow(PW_RENDERFULLCONTENT)
   - Windows 7 需要 alpha 通道修复
