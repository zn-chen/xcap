//go:build windows

package windows

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
		t.Logf("Monitor %d: Handle=%d, Name=%s, Position=(%d,%d), Size=%dx%d, Primary=%v",
			i, m.Handle, m.Name, m.X, m.Y, m.Width, m.Height, m.Primary)
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
	img, err := CaptureMonitor(monitors[0])
	if err != nil {
		t.Fatalf("CaptureMonitor failed: %v", err)
	}

	t.Logf("Captured: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())

	// Save to file for manual inspection
	f, err := os.Create("xcap_monitor_test.png")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	t.Logf("Saved screenshot to %s", f.Name())
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
		t.Logf("Monitor %d: ID=%d, Name=%s, Position=(%d,%d), Size=%dx%d, Primary=%v, Scale=%.2f",
			i, m.ID(), m.Name(), m.X(), m.Y(), m.Width(), m.Height(), m.IsPrimary(), m.ScaleFactor())
	}
}
