//go:build windows

package windows

import "image"

// Monitor 表示 Windows 上的显示器
type Monitor struct {
	info MonitorInfo
}

// NewMonitor 从 MonitorInfo 创建新的 Monitor
func NewMonitor(info MonitorInfo) *Monitor {
	return &Monitor{info: info}
}

// AllMonitors 返回所有可用的显示器
func AllMonitors() ([]*Monitor, error) {
	infos, err := GetAllMonitors()
	if err != nil {
		return nil, err
	}

	monitors := make([]*Monitor, len(infos))
	for i, info := range infos {
		monitors[i] = NewMonitor(info)
	}

	return monitors, nil
}

// ID 返回显示器的唯一标识符
func (m *Monitor) ID() uint32 {
	return uint32(m.info.Handle)
}

// Name 返回显示器的友好名称
func (m *Monitor) Name() string {
	return m.info.Name
}

// X 返回显示器左上角的 x 坐标
func (m *Monitor) X() int {
	return int(m.info.X)
}

// Y 返回��示器左上角的 y 坐标
func (m *Monitor) Y() int {
	return int(m.info.Y)
}

// Width 返回显示器的宽度（像素）
func (m *Monitor) Width() uint32 {
	return m.info.Width
}

// Height 返回显示器的高度（像素）
func (m *Monitor) Height() uint32 {
	return m.info.Height
}

// Rotation 返回旋转角度
func (m *Monitor) Rotation() float32 {
	return 0 // TODO: 通过 EnumDisplaySettings 实现
}

// ScaleFactor 返回 DPI 缩放因子
func (m *Monitor) ScaleFactor() float32 {
	dpiX, _ := GetMonitorDPI(m.info.Handle)
	if dpiX > 0 {
		return float32(dpiX) / 96.0
	}
	return 1.0
}

// Frequency 返回刷新率
func (m *Monitor) Frequency() float32 {
	return 0 // TODO: 通过 EnumDisplaySettings 实现
}

// IsPrimary 返回是否为主显示器
func (m *Monitor) IsPrimary() bool {
	return m.info.Primary
}

// IsBuiltin 返回是否为内置显示器
func (m *Monitor) IsBuiltin() bool {
	return false // TODO: 实现
}

// CaptureImage 截取整个显示器，返回 RGBA 图像
func (m *Monitor) CaptureImage() (*image.RGBA, error) {
	return CaptureMonitor(m.info)
}

// CaptureRegion 截取显示器的指定区域
func (m *Monitor) CaptureRegion(x, y, width, height uint32) (*image.RGBA, error) {
	return nil, ErrNotSupported // TODO: 实现
}
