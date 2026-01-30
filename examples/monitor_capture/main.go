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
	return replacer.Replace(name)
}

func main() {
	// Get all monitors
	monitors, err := xcap.AllMonitors()
	if err != nil {
		log.Fatalf("Failed to get monitors: %v", err)
	}

	fmt.Printf("Found %d monitor(s)\n", len(monitors))

	// Create output directory
	if err := os.MkdirAll("output", 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Capture each monitor
	for i, m := range monitors {
		fmt.Printf("Monitor %d: %s (%dx%d at %d,%d)\n",
			i, m.Name(), m.Width(), m.Height(), m.X(), m.Y())

		// Capture
		img, err := m.CaptureImage()
		if err != nil {
			log.Printf("Failed to capture monitor %d: %v", i, err)
			continue
		}

		// Save to file
		filename := fmt.Sprintf("output/monitor_%d_%s.png", i, sanitizeFilename(m.Name()))
		f, err := os.Create(filename)
		if err != nil {
			log.Printf("Failed to create file: %v", err)
			continue
		}

		if err := png.Encode(f, img); err != nil {
			f.Close()
			log.Printf("Failed to encode PNG: %v", err)
			continue
		}
		f.Close()

		fmt.Printf("  Saved: %s (%dx%d)\n", filename, img.Bounds().Dx(), img.Bounds().Dy())
	}

	fmt.Println("Done!")
}
