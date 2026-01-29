//go:build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations -mmacosx-version-min=10.15
#cgo LDFLAGS: -framework CoreGraphics -framework AppKit -framework CoreFoundation

#include "bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// 错误码，与 bridge.h 中的定义对应
const (
	errOK            = C.OWL_OK
	errNoMonitors    = C.OWL_ERR_NO_MONITORS
	errNoWindows     = C.OWL_ERR_NO_WINDOWS
	errCaptureFailed = C.OWL_ERR_CAPTURE_FAILED
	errAllocFailed   = C.OWL_ERR_ALLOC_FAILED
	errNotFound      = C.OWL_ERR_NOT_FOUND
)

// MonitorInfo 表示从 C 层获取的显示器信息
type MonitorInfo struct {
	ID     uint32
	Name   string
	X      int32
	Y      int32
	Width  uint32
	Height uint32
}

// WindowInfo 表示从 C 层获取的窗口信息
type WindowInfo struct {
	ID      uint32
	PID     uint32
	AppName string
	Title   string
	X       int32
	Y       int32
	Width   uint32
	Height  uint32
}

// CaptureResult 表示从 C 层获取的原始截图数据
type CaptureResult struct {
	Data        []byte
	Width       uint32
	Height      uint32
	BytesPerRow uint32
}

// GetAllMonitors 返回所有活动显示器的信息
func GetAllMonitors() ([]MonitorInfo, error) {
	var cMonitors *C.OwlMonitorInfo
	var cCount C.int

	result := C.owl_get_all_monitors(&cMonitors, &cCount)
	if result != errOK {
		return nil, fmt.Errorf("failed to get monitors: error code %d", result)
	}
	defer C.owl_free_monitors(cMonitors)

	if cCount == 0 {
		return nil, nil
	}

	// 将 C 数组转换为 Go slice
	count := int(cCount)
	monitors := make([]MonitorInfo, count)

	// 创建由 C 数组支持的 Go slice
	cSlice := unsafe.Slice(cMonitors, count)

	for i := 0; i < count; i++ {
		monitors[i] = MonitorInfo{
			ID:     uint32(cSlice[i].id),
			Name:   C.GoString(&cSlice[i].name[0]),
			X:      int32(cSlice[i].x),
			Y:      int32(cSlice[i].y),
			Width:  uint32(cSlice[i].width),
			Height: uint32(cSlice[i].height),
		}
	}

	return monitors, nil
}

// GetAllWindows 返回所有可见窗口的信息
func GetAllWindows() ([]WindowInfo, error) {
	var cWindows *C.OwlWindowInfo
	var cCount C.int

	result := C.owl_get_all_windows(&cWindows, &cCount)
	if result != errOK {
		return nil, fmt.Errorf("failed to get windows: error code %d", result)
	}
	defer C.owl_free_windows(cWindows)

	if cCount == 0 {
		return nil, nil
	}

	// 将 C 数组转换为 Go slice
	count := int(cCount)
	windows := make([]WindowInfo, count)

	// 创建由 C 数组支持的 Go slice
	cSlice := unsafe.Slice(cWindows, count)

	for i := 0; i < count; i++ {
		windows[i] = WindowInfo{
			ID:      uint32(cSlice[i].id),
			PID:     uint32(cSlice[i].pid),
			AppName: C.GoString(&cSlice[i].app_name[0]),
			Title:   C.GoString(&cSlice[i].title[0]),
			X:       int32(cSlice[i].x),
			Y:       int32(cSlice[i].y),
			Width:   uint32(cSlice[i].width),
			Height:  uint32(cSlice[i].height),
		}
	}

	return windows, nil
}

// CaptureMonitor 截取指定显示器，返回原始 BGRA 数据
func CaptureMonitor(displayID uint32) (*CaptureResult, error) {
	var cResult C.OwlCaptureResult

	result := C.owl_capture_monitor(C.uint32_t(displayID), &cResult)
	if result != errOK {
		return nil, fmt.Errorf("failed to capture monitor: error code %d", result)
	}
	defer C.owl_free_capture_result(&cResult)

	// 将数据复制到 Go slice
	dataLen := int(cResult.data_length)
	data := make([]byte, dataLen)
	copy(data, unsafe.Slice((*byte)(unsafe.Pointer(cResult.data)), dataLen))

	return &CaptureResult{
		Data:        data,
		Width:       uint32(cResult.width),
		Height:      uint32(cResult.height),
		BytesPerRow: uint32(cResult.bytes_per_row),
	}, nil
}

// CaptureWindow 截取指定窗口，返回原始 BGRA 数据
func CaptureWindow(windowID uint32) (*CaptureResult, error) {
	var cResult C.OwlCaptureResult

	result := C.owl_capture_window(C.uint32_t(windowID), &cResult)
	if result != errOK {
		return nil, fmt.Errorf("failed to capture window: error code %d", result)
	}
	defer C.owl_free_capture_result(&cResult)

	// 将数据复制到 Go slice
	dataLen := int(cResult.data_length)
	data := make([]byte, dataLen)
	copy(data, unsafe.Slice((*byte)(unsafe.Pointer(cResult.data)), dataLen))

	return &CaptureResult{
		Data:        data,
		Width:       uint32(cResult.width),
		Height:      uint32(cResult.height),
		BytesPerRow: uint32(cResult.bytes_per_row),
	}, nil
}
