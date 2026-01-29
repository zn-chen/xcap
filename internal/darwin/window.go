//go:build darwin

package darwin

import (
	"image"
)

// Window represents an application window on macOS
type Window struct {
	info WindowInfo
}

// NewWindow creates a new Window from WindowInfo
func NewWindow(info WindowInfo) *Window {
	return &Window{info: info}
}

// AllWindows returns all visible windows
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

// ID returns the unique identifier of the window
func (w *Window) ID() uint32 {
	return w.info.ID
}

// PID returns the process ID of the window's owner
func (w *Window) PID() uint32 {
	return w.info.PID
}

// AppName returns the name of the application that owns the window
func (w *Window) AppName() string {
	return w.info.AppName
}

// Title returns the window title
func (w *Window) Title() string {
	return w.info.Title
}

// X returns the x coordinate of the window's top-left corner
func (w *Window) X() int {
	return int(w.info.X)
}

// Y returns the y coordinate of the window's top-left corner
func (w *Window) Y() int {
	return int(w.info.Y)
}

// Z returns the z-order of the window (not implemented in minimal version)
func (w *Window) Z() int {
	return 0 // TODO: implement in full version
}

// Width returns the width of the window in pixels
func (w *Window) Width() uint32 {
	return w.info.Width
}

// Height returns the height of the window in pixels
func (w *Window) Height() uint32 {
	return w.info.Height
}

// IsMinimized returns true if the window is minimized (not implemented in minimal version)
func (w *Window) IsMinimized() bool {
	return false // TODO: implement in full version
}

// IsMaximized returns true if the window is maximized (not implemented in minimal version)
func (w *Window) IsMaximized() bool {
	return false // TODO: implement in full version
}

// IsFocused returns true if the window has input focus (not implemented in minimal version)
func (w *Window) IsFocused() bool {
	return false // TODO: implement in full version
}

// CurrentMonitor returns the monitor that contains most of the window (not implemented in minimal version)
func (w *Window) CurrentMonitor() (*Monitor, error) {
	return nil, ErrNotSupported // TODO: implement in full version
}

// CaptureImage captures the window content and returns an RGBA image
func (w *Window) CaptureImage() (*image.RGBA, error) {
	result, err := CaptureWindow(w.info.ID)
	if err != nil {
		return nil, err
	}

	return CaptureResultToImage(result), nil
}
