package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/anthropic-research/xcap/pkg/xcap"
)

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"|", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"?", "_",
		"*", "_",
	)
	result := replacer.Replace(name)
	// Limit length
	if len(result) > 50 {
		result = result[:50]
	}
	return result
}

func main() {
	// Get all windows
	windows, err := xcap.AllWindows()
	if err != nil {
		log.Fatalf("Failed to get windows: %v", err)
	}

	fmt.Printf("Found %d window(s)\n", len(windows))

	// Create output directory
	if err := os.MkdirAll("output", 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Capture each window
	captured := 0
	for i, w := range windows {
		// Skip tiny windows (likely system UI elements)
		if w.Width() < 100 || w.Height() < 100 {
			continue
		}

		fmt.Printf("Window %d: [%s] %s (%dx%d at %d,%d)\n",
			i, w.AppName(), w.Title(), w.Width(), w.Height(), w.X(), w.Y())

		// Capture
		img, err := w.CaptureImage()
		if err != nil {
			log.Printf("  Failed to capture: %v", err)
			continue
		}

		// Save to file
		title := w.Title()
		if title == "" {
			title = "untitled"
		}
		filename := fmt.Sprintf("output/window_%d_%s_%s.png",
			i, sanitizeFilename(w.AppName()), sanitizeFilename(title))

		f, err := os.Create(filename)
		if err != nil {
			log.Printf("  Failed to create file: %v", err)
			continue
		}

		if err := png.Encode(f, img); err != nil {
			f.Close()
			log.Printf("  Failed to encode PNG: %v", err)
			continue
		}
		f.Close()

		fmt.Printf("  Saved: %s (%dx%d)\n", filename, img.Bounds().Dx(), img.Bounds().Dy())
		captured++
	}

	fmt.Printf("Done! Captured %d window(s)\n", captured)
}
