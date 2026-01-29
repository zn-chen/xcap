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

// Error codes matching bridge.h
const (
	errOK            = C.OWL_OK
	errNoMonitors    = C.OWL_ERR_NO_MONITORS
	errNoWindows     = C.OWL_ERR_NO_WINDOWS
	errCaptureFailed = C.OWL_ERR_CAPTURE_FAILED
	errAllocFailed   = C.OWL_ERR_ALLOC_FAILED
	errNotFound      = C.OWL_ERR_NOT_FOUND
)

// MonitorInfo represents monitor information from C layer
type MonitorInfo struct {
	ID     uint32
	Name   string
	X      int32
	Y      int32
	Width  uint32
	Height uint32
}

// WindowInfo represents window information from C layer
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

// CaptureResult represents raw capture data from C layer
type CaptureResult struct {
	Data        []byte
	Width       uint32
	Height      uint32
	BytesPerRow uint32
}

// GetAllMonitors returns information about all active monitors
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

	// Convert C array to Go slice
	count := int(cCount)
	monitors := make([]MonitorInfo, count)

	// Create a Go slice backed by the C array
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

// GetAllWindows returns information about all visible windows
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

	// Convert C array to Go slice
	count := int(cCount)
	windows := make([]WindowInfo, count)

	// Create a Go slice backed by the C array
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

// CaptureMonitor captures the specified monitor and returns raw BGRA data
func CaptureMonitor(displayID uint32) (*CaptureResult, error) {
	var cResult C.OwlCaptureResult

	result := C.owl_capture_monitor(C.uint32_t(displayID), &cResult)
	if result != errOK {
		return nil, fmt.Errorf("failed to capture monitor: error code %d", result)
	}
	defer C.owl_free_capture_result(&cResult)

	// Copy data to Go slice
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

// CaptureWindow captures the specified window and returns raw BGRA data
func CaptureWindow(windowID uint32) (*CaptureResult, error) {
	var cResult C.OwlCaptureResult

	result := C.owl_capture_window(C.uint32_t(windowID), &cResult)
	if result != errOK {
		return nil, fmt.Errorf("failed to capture window: error code %d", result)
	}
	defer C.owl_free_capture_result(&cResult)

	// Copy data to Go slice
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
