// Package owl provides cross-platform screen and window capture functionality.
//
// owl-go is a Go implementation of screen capture inspired by the Rust library xcap.
// It supports capturing individual windows or entire monitors on macOS and Windows.
//
// Basic usage:
//
//	// Capture all monitors
//	monitors, err := owl.AllMonitors()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, m := range monitors {
//	    img, err := m.CaptureImage()
//	    if err != nil {
//	        continue
//	    }
//	    // Use img...
//	}
//
//	// Capture all windows
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
//	    // Use img...
//	}
package owl
