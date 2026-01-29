// Package owl 提供跨平台的屏幕和窗口截图功能。
//
// owl-go 是一个参考 Rust 库 xcap 实现的 Go 语言屏幕截图库，
// 支持在 macOS 和 Windows 上截取单个窗口或整个显示器。
//
// 基本用法：
//
//	// 截取所有显示器
//	monitors, err := owl.AllMonitors()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, m := range monitors {
//	    img, err := m.CaptureImage()
//	    if err != nil {
//	        continue
//	    }
//	    // 使用 img...
//	}
//
//	// 截取所有窗口
//	windows, err := owl.AllWindows()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, w := range windows {
//	    if w.IsMinimized() {
//	        continue
//	    }
//	    img, err := w.CaptureImage()
//	    if err != nil {
//	        continue
//	    }
//	    // 使用 img...
//	}
package owl
