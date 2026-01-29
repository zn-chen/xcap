package owl

import "image"

// Monitor represents a display/monitor
type Monitor interface {
	// ID returns the unique identifier of the monitor
	ID() uint32

	// Name returns the friendly name of the monitor
	Name() string

	// X returns the x coordinate of the monitor's top-left corner
	X() int

	// Y returns the y coordinate of the monitor's top-left corner
	Y() int

	// Width returns the width of the monitor in pixels
	Width() uint32

	// Height returns the height of the monitor in pixels
	Height() uint32

	// Rotation returns the rotation angle in degrees (0, 90, 180, 270)
	Rotation() float32

	// ScaleFactor returns the DPI scale factor (e.g., 2.0 for Retina/HiDPI)
	ScaleFactor() float32

	// Frequency returns the refresh rate in Hz
	Frequency() float32

	// IsPrimary returns true if this is the primary monitor
	IsPrimary() bool

	// IsBuiltin returns true if this is a built-in display (laptop screen)
	IsBuiltin() bool

	// CaptureImage captures the entire monitor and returns an RGBA image
	CaptureImage() (*image.RGBA, error)

	// CaptureRegion captures a specific region of the monitor
	CaptureRegion(x, y, width, height uint32) (*image.RGBA, error)
}
