// Package main 演示 xcap 库的基本用法。
//
// 本示例展示如何：
// - 枚举所有显示器及其属性
// - 枚举所有窗口及其属性
// - 截取显示器和窗口的屏幕截图
//
// 运行方式: go run ./examples/basic
package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zn-chen/xcap/pkg/xcap"
)

func main() {
	fmt.Println("=== xcap Basic Example ===")
	fmt.Println()

	// 创建输出目录
	outputDir := "xcap_output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// 演示 1: 显示器枚举和截图
	demoMonitors(outputDir)

	fmt.Println()

	// 演示 2: 窗口枚举和截图
	demoWindows(outputDir)

	fmt.Println()
	fmt.Printf("All screenshots saved to: %s/\n", outputDir)
}

// demoMonitors 演示显示器枚举和截图功能
func demoMonitors(outputDir string) {
	fmt.Println("--- Monitors ---")

	monitors, err := xcap.AllMonitors()
	if err != nil {
		log.Printf("Failed to get monitors: %v", err)
		return
	}

	fmt.Printf("Found %d monitor(s)\n\n", len(monitors))

	for i, m := range monitors {
		// 显示 Monitor 的各项属性
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

		// 截取屏幕
		img, err := m.CaptureImage()
		if err != nil {
			fmt.Printf("  Capture:     FAILED (%v)\n", err)
			continue
		}

		// 保存到文件
		filename := filepath.Join(outputDir, fmt.Sprintf("monitor_%d_%s.png", i, sanitize(m.Name())))
		if err := savePNG(filename, img); err != nil {
			fmt.Printf("  Capture:     FAILED to save (%v)\n", err)
			continue
		}

		fmt.Printf("  Capture:     %dx%d -> %s\n", img.Bounds().Dx(), img.Bounds().Dy(), filename)
		fmt.Println()
	}
}

// demoWindows 演示窗口枚举和截图功能
func demoWindows(outputDir string) {
	fmt.Println("--- Windows ---")

	windows, err := xcap.AllWindows()
	if err != nil {
		log.Printf("Failed to get windows: %v", err)
		return
	}

	fmt.Printf("Found %d window(s)\n\n", len(windows))

	// 只截取前 5 个尺寸合适的窗口
	captured := 0
	maxCaptures := 5

	for i, w := range windows {
		// 跳过太小的窗口（通常是系统 UI）
		if w.Width() < 200 || w.Height() < 200 {
			continue
		}

		// 显示 Window 的各项属性
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

		// 截取窗口
		img, err := w.CaptureImage()
		if err != nil {
			fmt.Printf("  Capture:     FAILED (%v)\n", err)
			fmt.Println()
			continue
		}

		// 保存到文件
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

// savePNG 将图像保存为 PNG 文件
func savePNG(filename string, img image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// sanitize 移除文件名中的非法字符
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

// truncate 截断字符串并添加省略号
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
