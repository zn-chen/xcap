package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zn-chen/xcap/pkg/xcap"
)

var (
	version        = "dev"
	disableMonitor bool
	disableWindows bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "xcap",
		Short:   "跨平台屏幕截图工具",
		Version: version,
		Run:     run,
	}

	rootCmd.Flags().BoolVar(&disableMonitor, "disable_monitor", false, "禁用显示器截图")
	rootCmd.Flags().BoolVar(&disableWindows, "disable_windows", false, "禁用窗口截图")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	outputDir := "output"

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "创建输出目录失败: %v\n", err)
		os.Exit(1)
	}

	var totalCaptured int

	if !disableMonitor {
		count := captureMonitors(outputDir)
		totalCaptured += count
	}

	if !disableWindows {
		count := captureWindows(outputDir)
		totalCaptured += count
	}

	fmt.Printf("\n完成! 共截取 %d 张图片，保存到 %s/\n", totalCaptured, outputDir)
}

func captureMonitors(outputDir string) int {
	monitors, err := xcap.AllMonitors()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取显示器列表失败: %v\n", err)
		return 0
	}

	fmt.Printf("找到 %d 个显示器\n", len(monitors))

	captured := 0
	for i, m := range monitors {
		img, err := m.CaptureImage()
		if err != nil {
			fmt.Printf("  显示器 %d: 截图失败 - %v\n", i+1, err)
			continue
		}

		filename := fmt.Sprintf("monitor_%d_%s.png", i+1, xcap.SanitizeFilename(m.Name()))
		path := filepath.Join(outputDir, filename)

		if err := saveImage(path, img); err != nil {
			fmt.Printf("  显示器 %d: 保存失败 - %v\n", i+1, err)
			continue
		}

		fmt.Printf("  显示器 %d: %s -> %s\n", i+1, m.Name(), filename)
		captured++
	}

	return captured
}

func captureWindows(outputDir string) int {
	windows, err := xcap.AllWindows()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取窗口列表失败: %v\n", err)
		return 0
	}

	fmt.Printf("找到 %d 个窗口\n", len(windows))

	captured := 0
	for i, w := range windows {
		minimized, _ := w.IsMinimized()
		if minimized || w.Width() < 50 || w.Height() < 50 {
			continue
		}

		img, err := w.CaptureImage()
		if err != nil {
			continue
		}

		title := w.Title()
		if len(title) > 30 {
			title = title[:30]
		}

		filename := fmt.Sprintf("window_%d_%s_%s.png", i+1, xcap.SanitizeFilename(w.AppName()), xcap.SanitizeFilename(title))
		path := filepath.Join(outputDir, filename)

		if err := saveImage(path, img); err != nil {
			continue
		}

		// 获取焦点状态
		focusedMark := ""
		if focused, err := w.IsFocused(); err == nil && focused {
			focusedMark = " [焦点]"
		}

		fmt.Printf("  窗口 %d: [%s] %s%s -> %s\n", i+1, w.AppName(), w.Title(), focusedMark, filename)
		captured++
	}

	return captured
}

func saveImage(path string, img *image.RGBA) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
