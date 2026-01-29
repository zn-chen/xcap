//go:build darwin

package darwin

import (
	"errors"
	"image"
)

// ErrNotSupported is returned when a feature is not implemented
var ErrNotSupported = errors.New("not supported")

// Monitor represents a display/monitor on macOS
type Monitor struct {
	info MonitorInfo
}

// NewMonitor creates a new Monitor from MonitorInfo
func NewMonitor(info MonitorInfo) *Monitor {
	return &Monitor{info: info}
}

// AllMonitors returns all available monitors
func AllMonitors() ([]*Monitor, error) {
	infos, err := GetAllMonitors()
	if err != nil {
		return nil, err
	}

	monitors := make([]*Monitor, len(infos))
	for i, info := range infos {
		monitors[i] = NewMonitor(info)
	}

	return monitors, nil
}

// ID returns the unique identifier of the monitor
func (m *Monitor) ID() uint32 {
	return m.info.ID
}

// Name returns the friendly name of the monitor
func (m *Monitor) Name() string {
	return m.info.Name
}

// X returns the x coordinate of the monitor's top-left corner
func (m *Monitor) X() int {
	return int(m.info.X)
}

// Y returns the y coordinate of the monitor's top-left corner
func (m *Monitor) Y() int {
	return int(m.info.Y)
}

// Width returns the width of the monitor in pixels
func (m *Monitor) Width() uint32 {
	return m.info.Width
}

// Height returns the height of the monitor in pixels
func (m *Monitor) Height() uint32 {
	return m.info.Height
}

// Rotation returns the rotation angle in degrees (not implemented in minimal version)
func (m *Monitor) Rotation() float32 {
	return 0 // TODO: implement in full version
}

// ScaleFactor returns the DPI scale factor (not implemented in minimal version)
func (m *Monitor) ScaleFactor() float32 {
	return 1.0 // TODO: implement in full version
}

// Frequency returns the refresh rate in Hz (not implemented in minimal version)
func (m *Monitor) Frequency() float32 {
	return 0 // TODO: implement in full version
}

// IsPrimary returns true if this is the primary monitor (not implemented in minimal version)
func (m *Monitor) IsPrimary() bool {
	return false // TODO: implement in full version
}

// IsBuiltin returns true if this is a built-in display (not implemented in minimal version)
func (m *Monitor) IsBuiltin() bool {
	return false // TODO: implement in full version
}

// CaptureImage captures the entire monitor and returns an RGBA image
func (m *Monitor) CaptureImage() (*image.RGBA, error) {
	result, err := CaptureMonitor(m.info.ID)
	if err != nil {
		return nil, err
	}

	return CaptureResultToImage(result), nil
}

// CaptureRegion captures a specific region of the monitor (not implemented in minimal version)
func (m *Monitor) CaptureRegion(x, y, width, height uint32) (*image.RGBA, error) {
	return nil, ErrNotSupported // TODO: implement in full version
}
