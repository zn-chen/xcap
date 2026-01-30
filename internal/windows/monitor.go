//go:build windows

package windows

import (
	"errors"
	"image"
	"syscall"
	"unsafe"
)

// ErrNotSupported 在功能未实现时返回
var ErrNotSupported = errors.New("not supported")

// ErrCaptureFailed 在截图失败时返回
var ErrCaptureFailed = errors.New("capture failed")

// MonitorInfo 表示从系统获取的显示器信息
type MonitorInfo struct {
	Handle  HMONITOR
	Name    string
	X       int32
	Y       int32
	Width   uint32
	Height  uint32
	Primary bool
}

// Monitor 表示 Windows 上的显示器
type Monitor struct {
	info MonitorInfo
}

// NewMonitor 从 MonitorInfo 创建新的 Monitor
func NewMonitor(info MonitorInfo) *Monitor {
	return &Monitor{info: info}
}

// 用于枚举显示器的回调数据
type enumMonitorData struct {
	monitors []MonitorInfo
}

// monitorEnumCallback 是 EnumDisplayMonitors 的回调函数
func monitorEnumCallback(hMonitor HMONITOR, hdcMonitor HDC, lprcMonitor *RECT, dwData uintptr) uintptr {
	data := (*enumMonitorData)(unsafe.Pointer(dwData))

	var mi MONITORINFOEXW
	mi.CbSize = uint32(unsafe.Sizeof(mi))

	if GetMonitorInfoW(hMonitor, &mi) {
		info := MonitorInfo{
			Handle:  hMonitor,
			Name:    UTF16ToString(mi.SzDevice[:]),
			X:       mi.RcMonitor.Left,
			Y:       mi.RcMonitor.Top,
			Width:   uint32(mi.RcMonitor.Right - mi.RcMonitor.Left),
			Height:  uint32(mi.RcMonitor.Bottom - mi.RcMonitor.Top),
			Primary: mi.DwFlags&1 != 0, // MONITORINFOF_PRIMARY = 1
		}
		data.monitors = append(data.monitors, info)
	}

	return 1 // 继续枚举
}

// GetAllMonitors 获取所有显示器信息
func GetAllMonitors() ([]MonitorInfo, error) {
	data := &enumMonitorData{
		monitors: make([]MonitorInfo, 0),
	}

	callback := syscall.NewCallback(monitorEnumCallback)
	EnumDisplayMonitors(0, nil, callback, uintptr(unsafe.Pointer(data)))

	if len(data.monitors) == 0 {
		return nil, errors.New("no monitors found")
	}

	return data.monitors, nil
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

// Y 返回显示器左上角的 y 坐标
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
	return 0 // TODO: 在完整版本中通过 EnumDisplaySettings 实现
}

// ScaleFactor 返回 DPI 缩放因子
func (m *Monitor) ScaleFactor() float32 {
	var dpiX, dpiY uint32
	hr := GetDpiForMonitor(m.info.Handle, MDT_EFFECTIVE_DPI, &dpiX, &dpiY)
	if hr == 0 && dpiX > 0 {
		return float32(dpiX) / 96.0
	}
	return 1.0
}

// Frequency 返回刷新率
func (m *Monitor) Frequency() float32 {
	return 0 // TODO: 在完整版本中通过 EnumDisplaySettings 实现
}

// IsPrimary 返回是否为主显示器
func (m *Monitor) IsPrimary() bool {
	return m.info.Primary
}

// IsBuiltin 返回是否为内置显示器
func (m *Monitor) IsBuiltin() bool {
	return false // TODO: 在完整版本中实现
}

// CaptureImage 截取整个显示器，返回 RGBA 图像
func (m *Monitor) CaptureImage() (*image.RGBA, error) {
	return CaptureMonitor(m.info)
}

// CaptureRegion 截取显示器的指定区域
func (m *Monitor) CaptureRegion(x, y, width, height uint32) (*image.RGBA, error) {
	return nil, ErrNotSupported // TODO: 在完整版本中实现
}

// CaptureMonitor 截取指定显示器
func CaptureMonitor(info MonitorInfo) (*image.RGBA, error) {
	width := int32(info.Width)
	height := int32(info.Height)

	if width <= 0 || height <= 0 {
		return nil, ErrCaptureFailed
	}

	// 获取桌面窗口的 DC
	hwndDesktop := GetDesktopWindow()
	hdcScreen := GetDC(hwndDesktop)
	if hdcScreen == 0 {
		return nil, ErrCaptureFailed
	}
	defer ReleaseDC(hwndDesktop, hdcScreen)

	// 创建兼容 DC
	hdcMem := CreateCompatibleDC(hdcScreen)
	if hdcMem == 0 {
		return nil, ErrCaptureFailed
	}
	defer DeleteDC(hdcMem)

	// 创建兼容位图
	hBitmap := CreateCompatibleBitmap(hdcScreen, width, height)
	if hBitmap == 0 {
		return nil, ErrCaptureFailed
	}
	defer DeleteObject(HGDIOBJ(hBitmap))

	// 选择位图到内存 DC
	oldBitmap := SelectObject(hdcMem, HGDIOBJ(hBitmap))
	defer SelectObject(hdcMem, oldBitmap)

	// 使用 BitBlt 复制屏幕内容
	if !BitBlt(hdcMem, 0, 0, width, height, hdcScreen, info.X, info.Y, SRCCOPY) {
		return nil, ErrCaptureFailed
	}

	// 准备位图信息
	bi := BITMAPINFO{
		BmiHeader: BITMAPINFOHEADER{
			BiSize:        uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
			BiWidth:       width,
			BiHeight:      -height, // 负值表示自上而下的 DIB
			BiPlanes:      1,
			BiBitCount:    32,
			BiCompression: BI_RGB,
		},
	}

	// 分配像素数据缓冲区
	dataSize := int(width) * int(height) * 4
	pixelData := make([]byte, dataSize)

	// 获取位图数据
	ret := GetDIBits(hdcMem, hBitmap, 0, uint32(height), unsafe.Pointer(&pixelData[0]), &bi, DIB_RGB_COLORS)
	if ret == 0 {
		return nil, ErrCaptureFailed
	}

	// 将 BGRA 转换为 RGBA
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for i := 0; i < dataSize; i += 4 {
		img.Pix[i+0] = pixelData[i+2] // R <- B
		img.Pix[i+1] = pixelData[i+1] // G <- G
		img.Pix[i+2] = pixelData[i+0] // B <- R
		img.Pix[i+3] = pixelData[i+3] // A <- A
	}

	return img, nil
}
