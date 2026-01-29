// Package main demonstrates the basic usage of owl-go library.
//
// This example shows how to:
// - Enumerate all monitors and their properties
// - Enumerate all windows and their properties
// - Capture screenshots of monitors and windows
//
// Run with: go run ./examples/basic
package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthropic-research/owl-go/pkg/owl"
)

func main() {
	fmt.Println("=== owl-go Basic Example ===")
	fmt.Println()

	// Create output directory
	outputDir := "owl_output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Demo 1: Monitor enumeration and capture
	demoMonitors(outputDir)

	fmt.Println()

	// Demo 2: Window enumeration and capture
	demoWindows(outputDir)

	fmt.Println()
	fmt.Printf("All screenshots saved to: %s/\n", outputDir)
}

func demoMonitors(outputDir string) {
	fmt.Println("--- Monitors ---")

	monitors, err := owl.AllMonitors()
	if err != nil {
		log.Printf("Failed to get monitors: %v", err)
		return
	}

	fmt.Printf("Found %d monitor(s)\n\n", len(monitors))

	for i, m := range monitors {
		// Display monitor properties
		fmt.Printf("Monitor #%d:\n", i)
		fmt.Printf("  ID:          %d\n", m.ID())
		fmt.Printf("  Name:        %s\n", m.Name())
		fmt.Printf("  Position:    (%d, %d)\n", m.X(), m.Y())
		fmt.Printf("  Size:        %d x %d\n", m.Width(), m.Height())
		fmt.Printf("  Rotation:    %.0f°\n", m.Rotation())
		fmt.Printf("  Scale:       %.1fx\n", m.ScaleFactor())
		fmt.Printf("  Frequency:   %.0f Hz\n", m.Frequency())
		fmt.Printf("  Primary:     %v\n", m.IsPrimary())
		fmt.Printf("  Built-in:    %v\n", m.IsBuiltin())

		// Capture screenshot
		img, err := m.CaptureImage()
		if err != nil {
			fmt.Printf("  Capture:     FAILED (%v)\n", err)
			continue
		}

		// Save to file
		filename := filepath.Join(outputDir, fmt.Sprintf("monitor_%d_%s.png", i, sanitize(m.Name())))
		if err := savePNG(filename, img); err != nil {
			fmt.Printf("  Capture:     FAILED to save (%v)\n", err)
			continue
		}

		fmt.Printf("  Capture:     %dx%d -> %s\n", img.Bounds().Dx(), img.Bounds().Dy(), filename)
		fmt.Println()
	}
}

func demoWindows(outputDir string) {
	fmt.Println("--- Windows ---")

	windows, err := owl.AllWindows()
	if err != nil {
		log.Printf("Failed to get windows: %v", err)
		return
	}

	fmt.Printf("Found %d window(s)\n\n", len(windows))

	// Only capture first 5 windows with reasonable size
	captured := 0
	maxCaptures := 5

	for i, w := range windows {
		// Skip tiny windows (likely system UI)
		if w.Width() < 200 || w.Height() < 200 {
			continue
		}

		// Display window properties
		fmt.Printf("Window #%d:\n", i)
		fmt.Printf("  ID:          %d\n", w.ID())
		fmt.Printf("  PID:         %d\n", w.PID())
		fmt.Printf("  App:         %s\n", w.AppName())
		fmt.Printf("  Title:       %s\n", truncate(w.Title(), 50))
		fmt.Printf("  Position:    (%d, %d)\n", w.X(), w.Y())
		fmt.Printf("  Size:        %d x %d\n", w.Width(), w.Height())
		fmt.Printf("  Z-Order:     %d\n", w.Z())
		fmt.Printf("  Minimized:   %v\n", w.IsMinimized())
		fmt.Printf("  Maximized:   %v\n", w.IsMaximized())
		fmt.Printf("  Focused:     %v\n", w.IsFocused())

		// Capture screenshot
		img, err := w.CaptureImage()
		if err != nil {
			fmt.Printf("  Capture:     FAILED (%v)\n", err)
			fmt.Println()
			continue
		}

		// Save to file
		title := w.Title()
		if title == "" {
			title = "untitled"
		}
		filename := filepath.Join(outputDir, fmt.Sprintf("window_%d_%s_%s.png",
			i, sanitize(w.AppName()), sanitize(title)))

		if err := savePNG(filename, img); err != nil {
			fmt.Printf("  Capture:     FAILED to save (%v)\n", err)
			fmt.Println()
			continue
		}

		fmt.Printf("  Capture:     %dx%d -> %s\n", img.Bounds().Dx(), img.Bounds().Dy(), filepath.Base(filename))
		fmt.Println()

		captured++
		if captured >= maxCaptures {
			fmt.Printf("(Showing first %d windows with size >= 200x200)\n", maxCaptures)
			break
		}
	}
}

// savePNG saves an image to a PNG file
func savePNG(filename string, img image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// sanitize removes invalid filename characters
func sanitize(name string) string {
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "|", "_",
		"\"", "_", "<", "_", ">", "_", "?", "_", "*", "_",
		"\n", "_", "\r", "_", "\t", "_",
	)
	result := replacer.Replace(name)
	if len(result) > 40 {
		result = result[:40]
	}
	return result
}

// truncate shortens a string with ellipsis
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
