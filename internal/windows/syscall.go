//go:build windows

package windows

import (
	"syscall"
	"unsafe"
)

// init 确保所有 DLL 被正确加载
// 这解决了某些情况下 LazyDLL 延迟加载导致的问题
func init() {
	// 预加载关键 DLL
	user32.Load()
	gdi32.Load()
	kernel32.Load()
	dwmapi.Load()
}

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	shcore   = syscall.NewLazyDLL("shcore.dll")
	dwmapi   = syscall.NewLazyDLL("dwmapi.dll")
	psapi    = syscall.NewLazyDLL("psapi.dll")

	// user32.dll
	procEnumDisplayMonitors    = user32.NewProc("EnumDisplayMonitors")
	procGetMonitorInfoW        = user32.NewProc("GetMonitorInfoW")
	procEnumWindows            = user32.NewProc("EnumWindows")
	procGetWindowTextW         = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW   = user32.NewProc("GetWindowTextLengthW")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible        = user32.NewProc("IsWindowVisible")
	procIsIconic               = user32.NewProc("IsIconic")
	procIsZoomed               = user32.NewProc("IsZoomed")
	procGetForegroundWindow    = user32.NewProc("GetForegroundWindow")
	procGetWindowRect          = user32.NewProc("GetWindowRect")
	procGetClientRect          = user32.NewProc("GetClientRect")
	procGetWindowDC            = user32.NewProc("GetWindowDC")
	procGetDC                  = user32.NewProc("GetDC")
	procReleaseDC              = user32.NewProc("ReleaseDC")
	procGetDesktopWindow       = user32.NewProc("GetDesktopWindow")
	procGetWindowLongPtrW      = user32.NewProc("GetWindowLongPtrW")
	procClientToScreen         = user32.NewProc("ClientToScreen")
	procPrintWindow            = user32.NewProc("PrintWindow")
	procGetClassNameW          = user32.NewProc("GetClassNameW")

	// gdi32.dll
	procCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procSelectObject           = gdi32.NewProc("SelectObject")
	procDeleteObject           = gdi32.NewProc("DeleteObject")
	procDeleteDC               = gdi32.NewProc("DeleteDC")
	procBitBlt                 = gdi32.NewProc("BitBlt")
	procGetDIBits              = gdi32.NewProc("GetDIBits")
	procCreateDIBSection       = gdi32.NewProc("CreateDIBSection")

	// kernel32.dll
	procOpenProcess            = kernel32.NewProc("OpenProcess")
	procCloseHandle            = kernel32.NewProc("CloseHandle")
	procGetCurrentProcessId    = kernel32.NewProc("GetCurrentProcessId")

	// shcore.dll
	procGetDpiForMonitor = shcore.NewProc("GetDpiForMonitor")

	// dwmapi.dll
	procDwmGetWindowAttribute    = dwmapi.NewProc("DwmGetWindowAttribute")
	procDwmIsCompositionEnabled  = dwmapi.NewProc("DwmIsCompositionEnabled")

	// psapi.dll
	procGetModuleBaseNameW = psapi.NewProc("GetModuleBaseNameW")

	// ntdll.dll (for OS version)
	ntdll                   = syscall.NewLazyDLL("ntdll.dll")
	procRtlGetNtVersionNumbers = ntdll.NewProc("RtlGetNtVersionNumbers")
)

// Windows 常量定义
const (
	SRCCOPY     = 0x00CC0020
	BI_RGB      = 0
	DIB_RGB_COLORS = 0

	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_READ           = 0x0010

	GWL_EXSTYLE = -20
	GWL_STYLE   = -16

	WS_EX_TOOLWINDOW = 0x00000080
	WS_EX_NOREDIRECTIONBITMAP = 0x00200000

	WS_VISIBLE = 0x10000000

	// PrintWindow 标志
	PW_CLIENTONLY        = 1
	PW_RENDERFULLCONTENT = 2 // Windows 8.1+
	PW_DEFAULT           = 0

	DWMWA_CLOAKED              = 14
	DWMWA_EXTENDED_FRAME_BOUNDS = 9

	MDT_EFFECTIVE_DPI = 0
)

// RECT 结构体
type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

// POINT 结构体
type POINT struct {
	X int32
	Y int32
}

// MONITORINFOEXW 结构体
type MONITORINFOEXW struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
	SzDevice  [32]uint16
}

// BITMAPINFOHEADER 结构体
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

// BITMAPINFO 结构体
type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors [1]uint32
}

// HWND 和 HDC 类型别名
type HWND uintptr
type HDC uintptr
type HBITMAP uintptr
type HGDIOBJ uintptr
type HMONITOR uintptr
type HANDLE uintptr

// EnumDisplayMonitors 枚举所有显示器
func EnumDisplayMonitors(hdc HDC, lprcClip *RECT, lpfnEnum uintptr, dwData uintptr) bool {
	ret, _, _ := procEnumDisplayMonitors.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lprcClip)),
		lpfnEnum,
		dwData,
	)
	return ret != 0
}

// GetMonitorInfoW 获取显示器信息
func GetMonitorInfoW(hMonitor HMONITOR, lpmi *MONITORINFOEXW) bool {
	ret, _, _ := procGetMonitorInfoW.Call(
		uintptr(hMonitor),
		uintptr(unsafe.Pointer(lpmi)),
	)
	return ret != 0
}

// EnumWindows 枚举所有顶层窗口
func EnumWindows(lpEnumFunc uintptr, lParam uintptr) bool {
	ret, _, _ := procEnumWindows.Call(lpEnumFunc, lParam)
	return ret != 0
}

// GetWindowTextW 获取窗口标题
func GetWindowTextW(hwnd HWND, lpString *uint16, nMaxCount int32) int32 {
	ret, _, _ := procGetWindowTextW.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(lpString)),
		uintptr(nMaxCount),
	)
	return int32(ret)
}

// GetWindowTextLengthW 获取窗口标题长度
func GetWindowTextLengthW(hwnd HWND) int32 {
	ret, _, _ := procGetWindowTextLengthW.Call(uintptr(hwnd))
	return int32(ret)
}

// GetWindowThreadProcessId 获取窗口的进程 ID
func GetWindowThreadProcessId(hwnd HWND, lpdwProcessId *uint32) uint32 {
	ret, _, _ := procGetWindowThreadProcessId.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(lpdwProcessId)),
	)
	return uint32(ret)
}

// IsWindowVisible 检查窗口是否可见
func IsWindowVisible(hwnd HWND) bool {
	ret, _, _ := procIsWindowVisible.Call(uintptr(hwnd))
	return ret != 0
}

// IsIconic 检查窗口是否最小化
func IsIconic(hwnd HWND) bool {
	ret, _, _ := procIsIconic.Call(uintptr(hwnd))
	return ret != 0
}

// IsZoomed 检查窗口是否最大化
func IsZoomed(hwnd HWND) bool {
	ret, _, _ := procIsZoomed.Call(uintptr(hwnd))
	return ret != 0
}

// GetForegroundWindow 获取当前焦点窗口
func GetForegroundWindow() HWND {
	ret, _, _ := procGetForegroundWindow.Call()
	return HWND(ret)
}

// GetWindowRect 获取窗口矩形
func GetWindowRect(hwnd HWND, lpRect *RECT) bool {
	ret, _, _ := procGetWindowRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(lpRect)),
	)
	return ret != 0
}

// GetClientRect 获取窗口客户区矩形
func GetClientRect(hwnd HWND, lpRect *RECT) bool {
	ret, _, _ := procGetClientRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(lpRect)),
	)
	return ret != 0
}

// GetWindowDC 获取窗口设备上下文
func GetWindowDC(hwnd HWND) HDC {
	ret, _, _ := procGetWindowDC.Call(uintptr(hwnd))
	return HDC(ret)
}

// GetDC 获取设备上下文
func GetDC(hwnd HWND) HDC {
	ret, _, _ := procGetDC.Call(uintptr(hwnd))
	return HDC(ret)
}

// ReleaseDC 释放设备上下文
func ReleaseDC(hwnd HWND, hdc HDC) int32 {
	ret, _, _ := procReleaseDC.Call(uintptr(hwnd), uintptr(hdc))
	return int32(ret)
}

// GetDesktopWindow 获取桌面窗口句柄
func GetDesktopWindow() HWND {
	ret, _, _ := procGetDesktopWindow.Call()
	return HWND(ret)
}

// GetWindowLongPtrW 获取窗口属性
func GetWindowLongPtrW(hwnd HWND, nIndex int32) uintptr {
	ret, _, _ := procGetWindowLongPtrW.Call(uintptr(hwnd), uintptr(nIndex))
	return ret
}

// GetClassNameW 获取窗口类名
func GetClassNameW(hwnd HWND, lpClassName *uint16, nMaxCount int32) int32 {
	ret, _, _ := procGetClassNameW.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(lpClassName)),
		uintptr(nMaxCount),
	)
	return int32(ret)
}

// ClientToScreen 将客户区坐标转换为屏幕坐标
func ClientToScreen(hwnd HWND, lpPoint *POINT) bool {
	ret, _, _ := procClientToScreen.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(lpPoint)),
	)
	return ret != 0
}

// PrintWindow 打印窗口内容到 DC
func PrintWindow(hwnd HWND, hdcBlt HDC, nFlags uint32) bool {
	ret, _, _ := procPrintWindow.Call(
		uintptr(hwnd),
		uintptr(hdcBlt),
		uintptr(nFlags),
	)
	return ret != 0
}

// CreateCompatibleDC 创建兼容的设备上下文
func CreateCompatibleDC(hdc HDC) HDC {
	ret, _, _ := procCreateCompatibleDC.Call(uintptr(hdc))
	return HDC(ret)
}

// CreateCompatibleBitmap 创建兼容的位图
func CreateCompatibleBitmap(hdc HDC, cx, cy int32) HBITMAP {
	ret, _, _ := procCreateCompatibleBitmap.Call(
		uintptr(hdc),
		uintptr(cx),
		uintptr(cy),
	)
	return HBITMAP(ret)
}

// SelectObject 选择对象到 DC
func SelectObject(hdc HDC, h HGDIOBJ) HGDIOBJ {
	ret, _, _ := procSelectObject.Call(uintptr(hdc), uintptr(h))
	return HGDIOBJ(ret)
}

// DeleteObject 删除 GDI 对象
func DeleteObject(ho HGDIOBJ) bool {
	ret, _, _ := procDeleteObject.Call(uintptr(ho))
	return ret != 0
}

// DeleteDC 删除设备上下文
func DeleteDC(hdc HDC) bool {
	ret, _, _ := procDeleteDC.Call(uintptr(hdc))
	return ret != 0
}

// BitBlt 位块传输
func BitBlt(hdc HDC, x, y, cx, cy int32, hdcSrc HDC, x1, y1 int32, rop uint32) bool {
	ret, _, _ := procBitBlt.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(hdcSrc),
		uintptr(x1),
		uintptr(y1),
		uintptr(rop),
	)
	return ret != 0
}

// GetDIBits 获取设备无关位图数据
func GetDIBits(hdc HDC, hbm HBITMAP, start, cLines uint32, lpvBits unsafe.Pointer, lpbmi *BITMAPINFO, usage uint32) int32 {
	ret, _, _ := procGetDIBits.Call(
		uintptr(hdc),
		uintptr(hbm),
		uintptr(start),
		uintptr(cLines),
		uintptr(lpvBits),
		uintptr(unsafe.Pointer(lpbmi)),
		uintptr(usage),
	)
	return int32(ret)
}

// OpenProcess 打开进程
func OpenProcess(dwDesiredAccess uint32, bInheritHandle bool, dwProcessId uint32) HANDLE {
	var inherit uintptr
	if bInheritHandle {
		inherit = 1
	}
	ret, _, _ := procOpenProcess.Call(
		uintptr(dwDesiredAccess),
		inherit,
		uintptr(dwProcessId),
	)
	return HANDLE(ret)
}

// CloseHandle 关闭句柄
func CloseHandle(hObject HANDLE) bool {
	ret, _, _ := procCloseHandle.Call(uintptr(hObject))
	return ret != 0
}

// GetCurrentProcessId 获取当前进程 ID
func GetCurrentProcessId() uint32 {
	ret, _, _ := procGetCurrentProcessId.Call()
	return uint32(ret)
}

// GetModuleBaseNameW 获取模块基本名称
func GetModuleBaseNameW(hProcess HANDLE, hModule HANDLE, lpBaseName *uint16, nSize uint32) uint32 {
	ret, _, _ := procGetModuleBaseNameW.Call(
		uintptr(hProcess),
		uintptr(hModule),
		uintptr(unsafe.Pointer(lpBaseName)),
		uintptr(nSize),
	)
	return uint32(ret)
}

// DwmGetWindowAttribute 获取 DWM 窗口属性
func DwmGetWindowAttribute(hwnd HWND, dwAttribute uint32, pvAttribute unsafe.Pointer, cbAttribute uint32) int32 {
	ret, _, _ := procDwmGetWindowAttribute.Call(
		uintptr(hwnd),
		uintptr(dwAttribute),
		uintptr(pvAttribute),
		uintptr(cbAttribute),
	)
	return int32(ret)
}

// DwmIsCompositionEnabled 检查 DWM 合成是否启用
func DwmIsCompositionEnabled() (bool, error) {
	var enabled int32
	ret, _, _ := procDwmIsCompositionEnabled.Call(uintptr(unsafe.Pointer(&enabled)))
	if ret != 0 {
		return false, syscall.Errno(ret)
	}
	return enabled != 0, nil
}

// GetDpiForMonitor 获取显示器 DPI
func GetDpiForMonitor(hMonitor HMONITOR, dpiType uint32, dpiX, dpiY *uint32) int32 {
	ret, _, _ := procGetDpiForMonitor.Call(
		uintptr(hMonitor),
		uintptr(dpiType),
		uintptr(unsafe.Pointer(dpiX)),
		uintptr(unsafe.Pointer(dpiY)),
	)
	return int32(ret)
}

// GetOSMajorVersion 获取 Windows 主版本号
// Windows 8 = 6.2, Windows 8.1 = 6.3, Windows 10 = 10.0
func GetOSMajorVersion() uint32 {
	var major, minor, build uint32
	procRtlGetNtVersionNumbers.Call(
		uintptr(unsafe.Pointer(&major)),
		uintptr(unsafe.Pointer(&minor)),
		uintptr(unsafe.Pointer(&build)),
	)
	return major
}

// UTF16ToString 将 UTF16 转换为 Go 字符串
func UTF16ToString(s []uint16) string {
	for i, v := range s {
		if v == 0 {
			s = s[:i]
			break
		}
	}
	return syscall.UTF16ToString(s)
}
