package owl

import "image"

// Monitor 表示一个显示器/监视器
type Monitor interface {
	// ID 返回显示器的唯一标识符
	ID() uint32

	// Name 返回显示器的友好名称
	Name() string

	// X 返回显示器左上角的 x 坐标
	X() int

	// Y 返回显示器左上角的 y 坐标
	Y() int

	// Width 返回显示器的宽度（像素）
	Width() uint32

	// Height 返回显示器的高度（像素）
	Height() uint32

	// Rotation 返回旋转角度（0, 90, 180, 270）
	Rotation() float32

	// ScaleFactor 返回 DPI 缩放因子（如 Retina 显示器为 2.0）
	ScaleFactor() float32

	// Frequency 返回刷新率（Hz）
	Frequency() float32

	// IsPrimary 返回是否为主显示器
	IsPrimary() bool

	// IsBuiltin 返回是否为内置显示器（如笔记本屏幕）
	IsBuiltin() bool

	// CaptureImage 截取整个显示器，返回 RGBA 图像
	CaptureImage() (*image.RGBA, error)

	// CaptureRegion 截取显示器的指定区域
	CaptureRegion(x, y, width, height uint32) (*image.RGBA, error)
}
