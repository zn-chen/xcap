package owl

import "errors"

// Common errors
var (
	// ErrNoMonitor is returned when no monitor is found
	ErrNoMonitor = errors.New("owl: no monitor found")

	// ErrNoWindow is returned when no window is found
	ErrNoWindow = errors.New("owl: no window found")

	// ErrCaptureFailed is returned when screen/window capture fails
	ErrCaptureFailed = errors.New("owl: capture failed")

	// ErrPermissionDenied is returned when the app lacks necessary permissions
	// On macOS, this typically means Screen Recording permission is not granted
	ErrPermissionDenied = errors.New("owl: permission denied")

	// ErrWindowMinimized is returned when trying to capture a minimized window
	ErrWindowMinimized = errors.New("owl: window is minimized")

	// ErrInvalidRegion is returned when the capture region is invalid
	ErrInvalidRegion = errors.New("owl: invalid capture region")

	// ErrNotSupported is returned when a feature is not supported on the current platform
	ErrNotSupported = errors.New("owl: not supported on this platform")
)
