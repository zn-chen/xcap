//go:build windows

package windows

/*
#cgo CFLAGS: -DUNICODE -D_UNICODE
#cgo LDFLAGS: -luser32 -lgdi32 -ldwmapi -lshcore -lpsapi

#include "bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"image"
	"syscall"
	"unsafe"
)

// 错误码，与 bridge.h 中的定义对应
const (
	errOK            = C.XCAP_OK
	errNoMonitors    = C.XCAP_ERR_NO_MONITORS
	errNoWindows     = C.XCAP_ERR_NO_WINDOWS
	errCaptureFailed = C.XCAP_ERR_CAPTURE_FAILED
	errAllocFailed   = C.XCAP_ERR_ALLOC_FAILED
	errNotFound      = C.XCAP_ERR_NOT_FOUND
)

// ErrNotSupported 在功能未实现时返回
var ErrNotSupported = errors.New("not supported")

// ErrCaptureFailed 在截图失败时返回
var ErrCaptureFailed = errors.New("capture failed")

// ErrNoMonitors 在没有找到显示器时返回
var ErrNoMonitors = errors.New("no monitors found")

// ErrNoWindows 在没有找到窗口时返回
var ErrNoWindows = errors.New("no windows found")

// HMONITOR 类型别名
type HMONITOR uintptr

// HWND 类型别名
type HWND uintptr

// MonitorInfo 表示从 C 层获取的显示器信息
type MonitorInfo struct {
	Handle  HMONITOR
	Name    string
	X       int32
	Y       int32
	Width   uint32
	Height  uint32
	Primary bool
}

// WindowInfo 表示从 C 层获取的窗口信息
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

// utf16ToString 将 UTF-16 数组转换为 Go 字符串
func utf16ToString(s []uint16) string {
	for i, v := range s {
		if v == 0 {
			s = s[:i]
			break
		}
	}
	return syscall.UTF16ToString(s)
}

// GetAllMonitors 返回所有活动显示器的信息
func GetAllMonitors() ([]MonitorInfo, error) {
	var cMonitors *C.XcapMonitorInfo
	var cCount C.int

	result := C.xcap_get_all_monitors(&cMonitors, &cCount)
	if result != errOK {
		if result == errNoMonitors {
			return nil, ErrNoMonitors
		}
		return nil, errors.New("failed to get monitors")
	}
	defer C.xcap_free_monitors(cMonitors)

	if cCount == 0 {
		return nil, nil
	}

	// 将 C 数组转换为 Go slice
	count := int(cCount)
	monitors := make([]MonitorInfo, count)

	// 创建由 C 数组支持的 Go slice
	cSlice := unsafe.Slice(cMonitors, count)

	for i := 0; i < count; i++ {
		// 将 UTF-16 name 转换为 Go 字符串
		nameSlice := make([]uint16, 32)
		for j := 0; j < 32; j++ {
			nameSlice[j] = uint16(cSlice[i].name[j])
		}

		monitors[i] = MonitorInfo{
			Handle:  HMONITOR(cSlice[i].handle),
			Name:    utf16ToString(nameSlice),
			X:       int32(cSlice[i].x),
			Y:       int32(cSlice[i].y),
			Width:   uint32(cSlice[i].width),
			Height:  uint32(cSlice[i].height),
			Primary: bool(cSlice[i].is_primary),
		}
	}

	return monitors, nil
}

// GetAllWindows 获取所有可见窗口信息（包括当前进程的窗口）
func GetAllWindows() ([]WindowInfo, error) {
	return GetAllWindowsWithOptions(false)
}

// GetAllWindowsWithOptions 获取所有可见窗口信息
// excludeCurrentProcess: 是否排除当前进程的窗口
func GetAllWindowsWithOptions(excludeCurrentProcess bool) ([]WindowInfo, error) {
	var cWindows *C.XcapWindowInfo
	var cCount C.int

	result := C.xcap_get_all_windows(&cWindows, &cCount, C.bool(excludeCurrentProcess))
	if result != errOK {
		if result == errNoWindows {
			return nil, ErrNoWindows
		}
		return nil, errors.New("failed to get windows")
	}
	defer C.xcap_free_windows(cWindows)

	if cCount == 0 {
		return nil, nil
	}

	// 将 C 数组转换为 Go slice
	count := int(cCount)
	windows := make([]WindowInfo, count)

	// 创建由 C 数组支持的 Go slice
	cSlice := unsafe.Slice(cWindows, count)

	for i := 0; i < count; i++ {
		// 将 UTF-16 app_name 转换为 Go 字符串
		appNameSlice := make([]uint16, 260)
		for j := 0; j < 260; j++ {
			appNameSlice[j] = uint16(cSlice[i].app_name[j])
		}

		// 将 UTF-16 title 转换为 Go 字符串
		titleSlice := make([]uint16, 256)
		for j := 0; j < 256; j++ {
			titleSlice[j] = uint16(cSlice[i].title[j])
		}

		windows[i] = WindowInfo{
			Handle:  HWND(cSlice[i].handle),
			PID:     uint32(cSlice[i].pid),
			AppName: utf16ToString(appNameSlice),
			Title:   utf16ToString(titleSlice),
			X:       int32(cSlice[i].x),
			Y:       int32(cSlice[i].y),
			Width:   uint32(cSlice[i].width),
			Height:  uint32(cSlice[i].height),
		}
	}

	return windows, nil
}

// CaptureMonitor 截取指定显示器，返回 RGBA 图像
func CaptureMonitor(info MonitorInfo) (*image.RGBA, error) {
	var cResult C.XcapCaptureResult

	result := C.xcap_capture_monitor(
		C.uintptr_t(info.Handle),
		C.int32_t(info.X),
		C.int32_t(info.Y),
		C.uint32_t(info.Width),
		C.uint32_t(info.Height),
		&cResult,
	)
	if result != errOK {
		return nil, ErrCaptureFailed
	}
	defer C.xcap_free_capture_result(&cResult)

	return convertBGRAToRGBA(&cResult), nil
}

// CaptureWindow 截取指定窗口，返回 RGBA 图像
func CaptureWindow(info WindowInfo) (*image.RGBA, error) {
	var cResult C.XcapCaptureResult

	result := C.xcap_capture_window(C.uintptr_t(info.Handle), &cResult)
	if result != errOK {
		return nil, ErrCaptureFailed
	}
	defer C.xcap_free_capture_result(&cResult)

	return convertBGRAToRGBA(&cResult), nil
}

// convertBGRAToRGBA 将 BGRA 像素数据转换为 RGBA 格式
func convertBGRAToRGBA(cResult *C.XcapCaptureResult) *image.RGBA {
	width := int(cResult.width)
	height := int(cResult.height)
	dataLen := int(cResult.data_length)

	// 复制数据到 Go slice
	pixelData := C.GoBytes(unsafe.Pointer(cResult.data), C.int(dataLen))

	// 创建 RGBA 图像
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// BGRA -> RGBA 转换
	for i := 0; i < dataLen; i += 4 {
		img.Pix[i+0] = pixelData[i+2] // R <- B
		img.Pix[i+1] = pixelData[i+1] // G <- G
		img.Pix[i+2] = pixelData[i+0] // B <- R
		img.Pix[i+3] = pixelData[i+3] // A <- A
	}

	return img
}

// IsWindowMinimized 检查窗口是否最小化
func IsWindowMinimized(handle HWND) bool {
	return bool(C.xcap_is_window_minimized(C.uintptr_t(handle)))
}

// IsWindowMaximized 检查窗口是否最大化
func IsWindowMaximized(handle HWND) bool {
	return bool(C.xcap_is_window_maximized(C.uintptr_t(handle)))
}

// IsWindowFocused 检查窗口是否拥有焦点
func IsWindowFocused(handle HWND) bool {
	return bool(C.xcap_is_window_focused(C.uintptr_t(handle)))
}

// GetMonitorDPI 获取显示器 DPI
func GetMonitorDPI(handle HMONITOR) (uint32, uint32) {
	var dpiX, dpiY C.uint32_t
	C.xcap_get_monitor_dpi(C.uintptr_t(handle), &dpiX, &dpiY)
	return uint32(dpiX), uint32(dpiY)
}
