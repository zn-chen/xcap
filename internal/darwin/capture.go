//go:build darwin

package darwin

import "image"

// BGRAToRGBA 将 BGRA 像素数据转换为 RGBA 格式
// 同时处理行对齐问题（bytes_per_row 可能大于 width * 4）
func BGRAToRGBA(data []byte, width, height, bytesPerRow uint32) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	srcOffset := 0
	dstOffset := 0
	rowBytes := int(width * 4)

	for y := uint32(0); y < height; y++ {
		// 逐行复制，处理对齐
		for x := 0; x < rowBytes; x += 4 {
			// BGRA -> RGBA: 交换 B 和 R
			img.Pix[dstOffset+x+0] = data[srcOffset+x+2] // R <- B
			img.Pix[dstOffset+x+1] = data[srcOffset+x+1] // G <- G
			img.Pix[dstOffset+x+2] = data[srcOffset+x+0] // B <- R
			img.Pix[dstOffset+x+3] = data[srcOffset+x+3] // A <- A
		}
		srcOffset += int(bytesPerRow)
		dstOffset += rowBytes
	}

	return img
}

// CaptureResultToImage 将 CaptureResult 转换为 image.RGBA
func CaptureResultToImage(result *CaptureResult) *image.RGBA {
	return BGRAToRGBA(result.Data, result.Width, result.Height, result.BytesPerRow)
}
