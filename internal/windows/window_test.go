//go:build windows

package windows

import (
	"image/png"
	"os"
	"testing"
)

func TestGetAllWindows(t *testing.T) {
	windows, err := GetAllWindows()
	if err != nil {
		t.Fatalf("GetAllWindows failed: %v", err)
	}

	t.Logf("Found %d windows", len(windows))

	for i, w := range windows {
		t.Logf("Window %d: Handle=%d, PID=%d, App=%s, Title=%s, Position=(%d,%d), Size=%dx%d",
			i, w.Handle, w.PID, w.AppName, w.Title, w.X, w.Y, w.Width, w.Height)
	}
}

func TestCaptureWindow(t *testing.T) {
	windows, err := GetAllWindows()
	if err != nil {
		t.Fatalf("GetAllWindows failed: %v", err)
	}

	if len(windows) == 0 {
		t.Skip("No windows found")
	}

	// Find a window with reasonable size
	var targetWindow *WindowInfo
	for i := range windows {
		if windows[i].Width > 100 && windows[i].Height > 100 {
			targetWindow = &windows[i]
			break
		}
	}

	if targetWindow == nil {
		t.Skip("No suitable window found")
	}

	t.Logf("Capturing window: Handle=%d, App=%s, Title=%s",
		targetWindow.Handle, targetWindow.AppName, targetWindow.Title)

	img, err := CaptureWindow(*targetWindow)
	if err != nil {
		t.Fatalf("CaptureWindow failed: %v", err)
	}

	t.Logf("Captured: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())

	// Save to file for manual inspection
	f, err := os.Create("xcap_window_test.png")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	t.Logf("Saved screenshot to %s", f.Name())
}

func TestAllWindows(t *testing.T) {
	windows, err := AllWindows()
	if err != nil {
		t.Fatalf("AllWindows failed: %v", err)
	}

	t.Logf("Found %d windows", len(windows))

	for i, w := range windows {
		t.Logf("Window %d: ID=%d, PID=%d, App=%s, Title=%s, Position=(%d,%d), Size=%dx%d, Minimized=%v, Maximized=%v, Focused=%v",
			i, w.ID(), w.PID(), w.AppName(), w.Title(), w.X(), w.Y(), w.Width(), w.Height(),
			w.IsMinimized(), w.IsMaximized(), w.IsFocused())
	}
}

func TestWindowCapture(t *testing.T) {
	wins, err := AllWindows()
	if err != nil {
		t.Fatalf("AllWindows failed: %v", err)
	}

	// Find first large enough window
	for _, w := range wins {
		if w.Width() > 200 && w.Height() > 200 && !w.IsMinimized() {
			t.Logf("Capturing: %s - %s", w.AppName(), w.Title())

			img, err := w.CaptureImage()
			if err != nil {
				t.Logf("Failed to capture: %v", err)
				continue
			}

			t.Logf("Successfully captured %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
			return
		}
	}

	t.Skip("No suitable window found for capture test")
}
