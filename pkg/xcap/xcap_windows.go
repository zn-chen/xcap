//go:build windows

package xcap

import (
	"image"

	"github.com/zn-chen/xcap/internal/windows"
)

// windowWrapper 包装 windows.Window 以实现 xcap.Window 接口
type windowWrapper struct {
	w *windows.Window
}

func (w *windowWrapper) ID() uint32        { return w.w.ID() }
func (w *windowWrapper) PID() uint32       { return w.w.PID() }
func (w *windowWrapper) AppName() string   { return w.w.AppName() }
func (w *windowWrapper) Title() string     { return w.w.Title() }
func (w *windowWrapper) X() int            { return w.w.X() }
func (w *windowWrapper) Y() int            { return w.w.Y() }
func (w *windowWrapper) Z() int            { return w.w.Z() }
func (w *windowWrapper) Width() uint32     { return w.w.Width() }
func (w *windowWrapper) Height() uint32    { return w.w.Height() }
func (w *windowWrapper) IsMinimized() bool { return w.w.IsMinimized() }
func (w *windowWrapper) IsMaximized() bool { return w.w.IsMaximized() }
func (w *windowWrapper) IsFocused() bool   { return w.w.IsFocused() }

func (w *windowWrapper) CurrentMonitor() (Monitor, error) {
	m, err := w.w.CurrentMonitor()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (w *windowWrapper) CaptureImage() (*image.RGBA, error) {
	return w.w.CaptureImage()
}

// AllMonitors 返回系统上所有可用的显示器
func AllMonitors() ([]Monitor, error) {
	monitors, err := windows.AllMonitors()
	if err != nil {
		return nil, err
	}

	result := make([]Monitor, len(monitors))
	for i, m := range monitors {
		result[i] = m
	}

	return result, nil
}

// AllWindows 返回系统上所有可见的窗口（包括当前进程的窗口）
func AllWindows() ([]Window, error) {
	return AllWindowsWithOptions(false)
}

// AllWindowsWithOptions 返回系统上所有可见的窗口
// excludeCurrentProcess: 是否排除当前进程的窗口
func AllWindowsWithOptions(excludeCurrentProcess bool) ([]Window, error) {
	wins, err := windows.AllWindowsWithOptions(excludeCurrentProcess)
	if err != nil {
		return nil, err
	}

	result := make([]Window, len(wins))
	for i, w := range wins {
		result[i] = &windowWrapper{w: w}
	}

	return result, nil
}
