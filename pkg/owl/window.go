package owl

import "image"

// Window 表示一个应用程序窗口
type Window interface {
	// ID 返回窗口的唯一标识符（Windows 上为 HWND，macOS 上为 CGWindowID）
	ID() uint32

	// PID 返回窗口所属进程的 ID
	PID() uint32

	// AppName 返回拥有该窗口的应用程序名称
	AppName() string

	// Title 返回窗口标题
	Title() string

	// X 返回窗口左上角的 x 坐标
	X() int

	// Y 返回窗口左上角的 y 坐标
	Y() int

	// Z 返回窗口的 Z 顺序（值越大越靠前）
	Z() int

	// Width 返回窗口的宽度（像素）
	Width() uint32

	// Height 返回窗口的高度（像素）
	Height() uint32

	// IsMinimized 返回窗口是否最小化
	IsMinimized() bool

	// IsMaximized 返回窗口是否最大化
	IsMaximized() bool

	// IsFocused 返回窗口是否拥有输入焦点
	IsFocused() bool

	// CurrentMonitor 返回窗口所在的显示器
	CurrentMonitor() (Monitor, error)

	// CaptureImage 截取窗口内容，返回 RGBA 图像
	CaptureImage() (*image.RGBA, error)
}
