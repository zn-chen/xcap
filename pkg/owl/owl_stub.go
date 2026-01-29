//go:build !darwin && !windows

package owl

// AllMonitors returns all available monitors on the system
func AllMonitors() ([]Monitor, error) {
	return nil, ErrNotSupported
}

// AllWindows returns all visible windows on the system
func AllWindows() ([]Window, error) {
	return nil, ErrNotSupported
}
