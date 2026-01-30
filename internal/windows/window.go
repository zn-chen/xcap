//go:build windows

package windows

import "image"

// Window 表示 Windows 上的应用程序窗口
type Window struct {
	info WindowInfo
}

// NewWindow 从 WindowInfo 创建新的 Window
func NewWindow(info WindowInfo) *Window {
	return &Window{info: info}
}

// AllWindows 返回所有可见的窗口（包括当前进程的窗口）
func AllWindows() ([]*Window, error) {
	return AllWindowsWithOptions(false)
}

// AllWindowsWithOptions 返回所有可见的窗口
// excludeCurrentProcess: 是否排除当前进程的窗口
func AllWindowsWithOptions(excludeCurrentProcess bool) ([]*Window, error) {
	infos, err := GetAllWindowsWithOptions(excludeCurrentProcess)
	if err != nil {
		return nil, err
	}

	windows := make([]*Window, len(infos))
	for i, info := range infos {
		windows[i] = NewWindow(info)
	}

	return windows, nil
}

// ID 返回窗口的唯一标识符
func (w *Window) ID() uint32 {
	return uint32(w.info.Handle)
}

// PID 返回窗口所属进程的 ID
func (w *Window) PID() uint32 {
	return w.info.PID
}

// AppName 返回拥有该窗口的应用程序名称
func (w *Window) AppName() string {
	return w.info.AppName
}

// Title 返回窗口标题
func (w *Window) Title() string {
	return w.info.Title
}

// X 返回窗口左上角的 x 坐标
func (w *Window) X() int {
	return int(w.info.X)
}

// Y 返回窗口左上角的 y 坐标
func (w *Window) Y() int {
	return int(w.info.Y)
}

// Z 返回窗口的 Z 顺序
func (w *Window) Z() int {
	return 0 // TODO: 通过 GetWindow 实现
}

// Width 返回窗口的宽度（像素）
func (w *Window) Width() uint32 {
	return w.info.Width
}

// Height 返回窗口的高度（像素）
func (w *Window) Height() uint32 {
	return w.info.Height
}

// IsMinimized 返回窗口是否最小化
func (w *Window) IsMinimized() (bool, error) {
	return IsWindowMinimized(w.info.Handle), nil
}

// IsMaximized 返回窗口是否最大化
func (w *Window) IsMaximized() (bool, error) {
	return IsWindowMaximized(w.info.Handle), nil
}

// IsFocused 返回窗口是否拥有输入焦点
func (w *Window) IsFocused() (bool, error) {
	return IsWindowFocused(w.info.Handle), nil
}

// CurrentMonitor 返回窗口所在的显示器
func (w *Window) CurrentMonitor() (*Monitor, error) {
	return nil, ErrNotSupported // TODO: 通过 MonitorFromWindow 实现
}

// CaptureImage 截取窗口内容，返回 RGBA 图像
func (w *Window) CaptureImage() (*image.RGBA, error) {
	return CaptureWindow(w.info)
}
