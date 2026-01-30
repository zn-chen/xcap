//go:build darwin

package darwin

import (
	"image/png"
	"os"
	"testing"
)

func TestGetAllMonitors(t *testing.T) {
	monitors, err := GetAllMonitors()
	if err != nil {
		t.Fatalf("GetAllMonitors failed: %v", err)
	}

	if len(monitors) == 0 {
		t.Fatal("Expected at least one monitor")
	}

	for i, m := range monitors {
		t.Logf("Monitor %d: ID=%d, Name=%s, Position=(%d,%d), Size=%dx%d",
			i, m.ID, m.Name, m.X, m.Y, m.Width, m.Height)
	}
}

func TestCaptureMonitor(t *testing.T) {
	monitors, err := GetAllMonitors()
	if err != nil {
		t.Fatalf("GetAllMonitors failed: %v", err)
	}

	if len(monitors) == 0 {
		t.Fatal("Expected at least one monitor")
	}

	// Capture first monitor
	result, err := CaptureMonitor(monitors[0].ID)
	if err != nil {
		t.Fatalf("CaptureMonitor failed: %v", err)
	}

	t.Logf("Captured: %dx%d, bytes_per_row=%d, data_len=%d",
		result.Width, result.Height, result.BytesPerRow, len(result.Data))

	// Convert to image
	img := CaptureResultToImage(result)
	if img == nil {
		t.Fatal("Failed to convert capture result to image")
	}

	// Save to file for manual inspection
	f, err := os.Create("/tmp/xcap_monitor_test.png")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	t.Logf("Saved screenshot to /tmp/xcap_monitor_test.png")
}

func TestAllMonitors(t *testing.T) {
	monitors, err := AllMonitors()
	if err != nil {
		t.Fatalf("AllMonitors failed: %v", err)
	}

	if len(monitors) == 0 {
		t.Fatal("Expected at least one monitor")
	}

	for i, m := range monitors {
		t.Logf("Monitor %d: ID=%d, Name=%s, Position=(%d,%d), Size=%dx%d",
			i, m.ID(), m.Name(), m.X(), m.Y(), m.Width(), m.Height())
	}
}
