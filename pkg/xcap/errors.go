package xcap

import "errors"

// 通用错误定义
var (
	// ErrNoMonitor 在找不到显示器时返回
	ErrNoMonitor = errors.New("xcap: no monitor found")

	// ErrNoWindow 在找不到窗口时返回
	ErrNoWindow = errors.New("xcap: no window found")

	// ErrCaptureFailed 在屏幕/窗口截图失败时返回
	ErrCaptureFailed = errors.New("xcap: capture failed")

	// ErrPermissionDenied 在应用缺少必要权限时返回
	// 在 macOS 上，通常表示未授予 Screen Recording 权限
	ErrPermissionDenied = errors.New("xcap: permission denied")

	// ErrWindowMinimized 在尝试截取最小化窗口时返回
	ErrWindowMinimized = errors.New("xcap: window is minimized")

	// ErrInvalidRegion 在截图区域无效时返回
	ErrInvalidRegion = errors.New("xcap: invalid capture region")

	// ErrNotSupported 在当前平台不支持该功能时返回
	ErrNotSupported = errors.New("xcap: not supported on this platform")
)
