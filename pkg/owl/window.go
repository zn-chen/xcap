package owl

import "image"

// Window represents an application window
type Window interface {
	// ID returns the unique identifier of the window (HWND on Windows, CGWindowID on macOS)
	ID() uint32

	// PID returns the process ID of the window's owner
	PID() uint32

	// AppName returns the name of the application that owns the window
	AppName() string

	// Title returns the window title
	Title() string

	// X returns the x coordinate of the window's top-left corner
	X() int

	// Y returns the y coordinate of the window's top-left corner
	Y() int

	// Z returns the z-order of the window (higher values are on top)
	Z() int

	// Width returns the width of the window in pixels
	Width() uint32

	// Height returns the height of the window in pixels
	Height() uint32

	// IsMinimized returns true if the window is minimized
	IsMinimized() bool

	// IsMaximized returns true if the window is maximized
	IsMaximized() bool

	// IsFocused returns true if the window has input focus
	IsFocused() bool

	// CurrentMonitor returns the monitor that contains most of the window
	CurrentMonitor() (Monitor, error)

	// CaptureImage captures the window content and returns an RGBA image
	CaptureImage() (*image.RGBA, error)
}
