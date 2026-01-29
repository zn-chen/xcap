//go:build darwin

package darwin

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
		t.Logf("Window %d: ID=%d, PID=%d, App=%s, Title=%s, Position=(%d,%d), Size=%dx%d",
			i, w.ID, w.PID, w.AppName, w.Title, w.X, w.Y, w.Width, w.Height)
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

	t.Logf("Capturing window: ID=%d, App=%s, Title=%s",
		targetWindow.ID, targetWindow.AppName, targetWindow.Title)

	result, err := CaptureWindow(targetWindow.ID)
	if err != nil {
		t.Fatalf("CaptureWindow failed: %v", err)
	}

	t.Logf("Captured: %dx%d, bytes_per_row=%d, data_len=%d",
		result.Width, result.Height, result.BytesPerRow, len(result.Data))

	// Convert to image
	img := CaptureResultToImage(result)
	if img == nil {
		t.Fatal("Failed to convert capture result to image")
	}

	// Save to file for manual inspection
	f, err := os.Create("/tmp/owl_window_test.png")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	t.Logf("Saved screenshot to /tmp/owl_window_test.png")
}

func TestAllWindows(t *testing.T) {
	windows, err := AllWindows()
	if err != nil {
		t.Fatalf("AllWindows failed: %v", err)
	}

	t.Logf("Found %d windows", len(windows))

	for i, w := range windows {
		t.Logf("Window %d: ID=%d, PID=%d, App=%s, Title=%s, Position=(%d,%d), Size=%dx%d",
			i, w.ID(), w.PID(), w.AppName(), w.Title(), w.X(), w.Y(), w.Width(), w.Height())
	}
}
