//go:build !darwin && !windows

package owl

// AllMonitors 返回系统上所有可用的显示器
func AllMonitors() ([]Monitor, error) {
	return nil, ErrNotSupported
}

// AllWindows 返回系统上所有可见的窗口
func AllWindows() ([]Window, error) {
	return nil, ErrNotSupported
}
