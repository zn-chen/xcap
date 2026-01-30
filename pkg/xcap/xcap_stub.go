//go:build !darwin && !windows

package xcap

// AllMonitors 返回系统上所有可用的显示器
func AllMonitors() ([]Monitor, error) {
	return nil, ErrNotSupported
}

// AllWindows 返回系统上所有可见的窗口（包括当前进程的窗口）
func AllWindows() ([]Window, error) {
	return nil, ErrNotSupported
}

// AllWindowsWithOptions 返回系统上所有可见的窗口
// excludeCurrentProcess: 是否排除当前进程的窗口
func AllWindowsWithOptions(excludeCurrentProcess bool) ([]Window, error) {
	return nil, ErrNotSupported
}
