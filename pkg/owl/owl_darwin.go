//go:build darwin

package owl

import (
	"image"

	"github.com/anthropic-research/owl-go/internal/darwin"
)

// windowWrapper wraps darwin.Window to implement owl.Window interface
type windowWrapper struct {
	w *darwin.Window
}

func (w *windowWrapper) ID() uint32           { return w.w.ID() }
func (w *windowWrapper) PID() uint32          { return w.w.PID() }
func (w *windowWrapper) AppName() string      { return w.w.AppName() }
func (w *windowWrapper) Title() string        { return w.w.Title() }
func (w *windowWrapper) X() int               { return w.w.X() }
func (w *windowWrapper) Y() int               { return w.w.Y() }
func (w *windowWrapper) Z() int               { return w.w.Z() }
func (w *windowWrapper) Width() uint32        { return w.w.Width() }
func (w *windowWrapper) Height() uint32       { return w.w.Height() }
func (w *windowWrapper) IsMinimized() bool    { return w.w.IsMinimized() }
func (w *windowWrapper) IsMaximized() bool    { return w.w.IsMaximized() }
func (w *windowWrapper) IsFocused() bool      { return w.w.IsFocused() }

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

// AllMonitors returns all available monitors on the system
func AllMonitors() ([]Monitor, error) {
	monitors, err := darwin.AllMonitors()
	if err != nil {
		return nil, err
	}

	result := make([]Monitor, len(monitors))
	for i, m := range monitors {
		result[i] = m
	}

	return result, nil
}

// AllWindows returns all visible windows on the system
func AllWindows() ([]Window, error) {
	windows, err := darwin.AllWindows()
	if err != nil {
		return nil, err
	}

	result := make([]Window, len(windows))
	for i, w := range windows {
		result[i] = &windowWrapper{w: w}
	}

	return result, nil
}
