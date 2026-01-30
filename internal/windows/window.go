//go:build windows

package windows

import (
	"image"
	"syscall"
	"unsafe"
)

// WindowInfo 表示从系统获取的窗口信息
type WindowInfo struct {
	Handle  HWND
	PID     uint32
	AppName string
	Title   string
	X       int32
	Y       int32
	Width   uint32
	Height  uint32
}

// Window 表示 Windows 上的应用程序窗口
type Window struct {
	info WindowInfo
}

// NewWindow 从 WindowInfo 创建新的 Window
func NewWindow(info WindowInfo) *Window {
	return &Window{info: info}
}

// 用于枚举窗口的回调数据
type enumWindowData struct {
	windows    []WindowInfo
	currentPid uint32
}

// isWindowCloaked 检查窗口是否被隐藏（DWM cloaked）
func isWindowCloaked(hwnd HWND) bool {
	var cloaked int32
	hr := DwmGetWindowAttribute(hwnd, DWMWA_CLOAKED, unsafe.Pointer(&cloaked), uint32(unsafe.Sizeof(cloaked)))
	return hr == 0 && cloaked != 0
}

// getWindowExtendedFrameBounds 获取 DWM 扩展边框（更精确的窗口边界）
func getWindowExtendedFrameBounds(hwnd HWND) (RECT, bool) {
	var rect RECT
	hr := DwmGetWindowAttribute(hwnd, DWMWA_EXTENDED_FRAME_BOUNDS, unsafe.Pointer(&rect), uint32(unsafe.Sizeof(rect)))
	return rect, hr == 0
}

// windowEnumCallback 是 EnumWindows 的回调函数
func windowEnumCallback(hwnd HWND, lParam uintptr) uintptr {
	data := (*enumWindowData)(unsafe.Pointer(lParam))

	// 检查窗口是否可见
	if !IsWindowVisible(hwnd) {
		return 1
	}

	// 检查窗口是否被 cloaked
	if isWindowCloaked(hwnd) {
		return 1
	}

	// 获取进程 ID
	var pid uint32
	GetWindowThreadProcessId(hwnd, &pid)

	// 跳过当前进程的窗口（避免死锁）
	if pid == data.currentPid {
		return 1
	}

	// 获取窗口样式
	exStyle := GetWindowLongPtrW(hwnd, GWL_EXSTYLE)

	// 获取窗口类名
	className := make([]uint16, 256)
	GetClassNameW(hwnd, &className[0], 256)
	classNameStr := UTF16ToString(className)

	// 过滤工具窗口（但保留特定系统窗口如任务栏）
	if exStyle&WS_EX_TOOLWINDOW != 0 {
		// 允许 Shell_TrayWnd（任务栏）和 Shell_SecondaryTrayWnd
		if classNameStr != "Shell_TrayWnd" && classNameStr != "Shell_SecondaryTrayWnd" {
			return 1
		}
	}

	// 过滤不重定向位图的窗口
	if exStyle&WS_EX_NOREDIRECTIONBITMAP != 0 {
		return 1
	}

	// 过滤系统窗口
	if classNameStr == "Progman" || classNameStr == "Button" || classNameStr == "Windows.UI.Core.CoreWindow" {
		return 1
	}

	// 获取窗口矩形
	var rect RECT
	if !GetWindowRect(hwnd, &rect) {
		return 1
	}

	// 尝试获取 DWM 扩展边框（更精确）
	if extRect, ok := getWindowExtendedFrameBounds(hwnd); ok {
		rect = extRect
	}

	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top

	// 过滤空矩形
	if width <= 0 || height <= 0 {
		return 1
	}

	// 获取窗口标题
	titleLen := GetWindowTextLengthW(hwnd)
	var title string
	if titleLen > 0 {
		titleBuf := make([]uint16, titleLen+1)
		GetWindowTextW(hwnd, &titleBuf[0], titleLen+1)
		title = UTF16ToString(titleBuf)
	}

	// 获取进程名称
	appName := getProcessName(pid)

	info := WindowInfo{
		Handle:  hwnd,
		PID:     pid,
		AppName: appName,
		Title:   title,
		X:       rect.Left,
		Y:       rect.Top,
		Width:   uint32(width),
		Height:  uint32(height),
	}

	data.windows = append(data.windows, info)
	return 1 // 继续枚举
}

// getProcessName 获取进程名称
func getProcessName(pid uint32) string {
	hProcess := OpenProcess(PROCESS_QUERY_INFORMATION|PROCESS_VM_READ, false, pid)
	if hProcess == 0 {
		return ""
	}
	defer CloseHandle(hProcess)

	nameBuf := make([]uint16, 260)
	ret := GetModuleBaseNameW(hProcess, 0, &nameBuf[0], 260)
	if ret == 0 {
		return ""
	}

	return UTF16ToString(nameBuf)
}

// dummyEnumProc 是一个空的枚举回调，用于初始化 syscall.NewCallback
func dummyEnumProc(hMonitor HMONITOR, hdcMonitor HDC, lprcMonitor *RECT, dwData uintptr) uintptr {
	return 1
}

// initCallbackSystem 通过调用 EnumDisplayMonitors 初始化回调系统
// 这解决了直接调用 EnumWindows 时返回空结果的问题
var callbackInitialized = false

func ensureCallbackInitialized() {
	if callbackInitialized {
		return
	}
	callback := syscall.NewCallback(dummyEnumProc)
	EnumDisplayMonitors(0, nil, callback, 0)
	callbackInitialized = true
}

// GetAllWindows 获取所有可见窗口信息
func GetAllWindows() ([]WindowInfo, error) {
	// 确保回调系统已初始化
	ensureCallbackInitialized()

	data := &enumWindowData{
		windows:    make([]WindowInfo, 0),
		currentPid: GetCurrentProcessId(),
	}

	callback := syscall.NewCallback(windowEnumCallback)
	EnumWindows(callback, uintptr(unsafe.Pointer(data)))

	return data.windows, nil
}

// AllWindows 返回所有可见的窗口
func AllWindows() ([]*Window, error) {
	infos, err := GetAllWindows()
	if err != nil {
		return nil, err
	}

	windows := make([]*Window, len(infos))
	for i, info := range infos {
		windows[i] = NewWindow(info)
	}

	return windows, nil
}

// ID 返回窗口的唯一标识符
func (w *Window) ID() uint32 {
	return uint32(w.info.Handle)
}

// PID 返回窗口所属进程的 ID
func (w *Window) PID() uint32 {
	return w.info.PID
}

// AppName 返回拥有该窗口的应用程序名称
func (w *Window) AppName() string {
	return w.info.AppName
}

// Title 返回窗口标题
func (w *Window) Title() string {
	return w.info.Title
}

// X 返回窗口左上角的 x 坐标
func (w *Window) X() int {
	return int(w.info.X)
}

// Y 返回窗口左上角的 y 坐标
func (w *Window) Y() int {
	return int(w.info.Y)
}

// Z 返回窗口的 Z 顺序
func (w *Window) Z() int {
	return 0 // TODO: 在完整版本中通过 GetWindow 实现
}

// Width 返回窗口的宽度（像素）
func (w *Window) Width() uint32 {
	return w.info.Width
}

// Height 返回窗口的高度（像素）
func (w *Window) Height() uint32 {
	return w.info.Height
}

// IsMinimized 返回窗口是否最小化
func (w *Window) IsMinimized() bool {
	return IsIconic(w.info.Handle)
}

// IsMaximized 返回窗口是否最大化
func (w *Window) IsMaximized() bool {
	return IsZoomed(w.info.Handle)
}

// IsFocused 返回窗口是否拥有输入焦点
func (w *Window) IsFocused() bool {
	return GetForegroundWindow() == w.info.Handle
}

// CurrentMonitor 返回窗口所在的显示器
func (w *Window) CurrentMonitor() (*Monitor, error) {
	return nil, ErrNotSupported // TODO: 在完整版本中通过 MonitorFromWindow 实现
}

// CaptureImage 截取窗口内容，返回 RGBA 图像
func (w *Window) CaptureImage() (*image.RGBA, error) {
	return CaptureWindow(w.info)
}

// CaptureWindow 截取指定窗口
// 与 xcap Rust 实现保持一致的捕获策略：
// 1. Windows 8+ 使用 PrintWindow(flag=2)
// 2. DWM 合成启用时使用 PrintWindow(flag=0)
// 3. 使用 PrintWindow(flag=4)
// 4. 最后使用 BitBlt 作为回退
func CaptureWindow(info WindowInfo) (*image.RGBA, error) {
	hwnd := info.Handle

	// 获取窗口矩形
	var rect RECT
	if !GetWindowRect(hwnd, &rect) {
		return nil, ErrCaptureFailed
	}

	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top

	if width <= 0 || height <= 0 {
		return nil, ErrCaptureFailed
	}

	// 获取窗口 DC
	hdcWindow := GetWindowDC(hwnd)
	if hdcWindow == 0 {
		return nil, ErrCaptureFailed
	}
	defer ReleaseDC(hwnd, hdcWindow)

	// 创建兼容 DC
	hdcMem := CreateCompatibleDC(hdcWindow)
	if hdcMem == 0 {
		return nil, ErrCaptureFailed
	}
	defer DeleteDC(hdcMem)

	// 创建兼容位图
	hBitmap := CreateCompatibleBitmap(hdcWindow, width, height)
	if hBitmap == 0 {
		return nil, ErrCaptureFailed
	}
	defer DeleteObject(HGDIOBJ(hBitmap))

	// 选择位图到内存 DC
	oldBitmap := SelectObject(hdcMem, HGDIOBJ(hBitmap))
	defer SelectObject(hdcMem, oldBitmap)

	// 与 xcap Rust 一致的回退策略
	captured := false

	// 1. Windows 8+ 使用 PrintWindow(flag=2) - PW_RENDERFULLCONTENT
	if GetOSMajorVersion() >= 8 {
		if PrintWindow(hwnd, hdcMem, PW_RENDERFULLCONTENT) {
			captured = true
		}
	}

	// 2. DWM 合成启用时使用 PrintWindow(flag=0) - PW_DEFAULT
	if !captured {
		if dwmEnabled, err := DwmIsCompositionEnabled(); err == nil && dwmEnabled {
			if PrintWindow(hwnd, hdcMem, PW_DEFAULT) {
				captured = true
			}
		}
	}

	// 3. 使用 PrintWindow(flag=4)
	if !captured {
		if PrintWindow(hwnd, hdcMem, 4) {
			captured = true
		}
	}

	// 4. 最后使用 BitBlt 作为回退（使用窗口 DC，与 xcap Rust 一致）
	if !captured {
		if BitBlt(hdcMem, 0, 0, width, height, hdcWindow, 0, 0, SRCCOPY) {
			captured = true
		}
	}

	if !captured {
		return nil, ErrCaptureFailed
	}

	// 准备位图信息
	bi := BITMAPINFO{
		BmiHeader: BITMAPINFOHEADER{
			BiSize:        uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
			BiWidth:       width,
			BiHeight:      -height, // 负值表示自上而下的 DIB
			BiPlanes:      1,
			BiBitCount:    32,
			BiCompression: BI_RGB,
		},
	}

	// 分配像素数据缓冲区
	dataSize := int(width) * int(height) * 4
	pixelData := make([]byte, dataSize)

	// 获取位图数据
	ret := GetDIBits(hdcMem, hBitmap, 0, uint32(height), unsafe.Pointer(&pixelData[0]), &bi, DIB_RGB_COLORS)
	if ret == 0 {
		return nil, ErrCaptureFailed
	}

	// 将 BGRA 转换为 RGBA
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for i := 0; i < dataSize; i += 4 {
		img.Pix[i+0] = pixelData[i+2] // R <- B
		img.Pix[i+1] = pixelData[i+1] // G <- G
		img.Pix[i+2] = pixelData[i+0] // B <- R
		img.Pix[i+3] = pixelData[i+3] // A <- A
	}

	return img, nil
}
