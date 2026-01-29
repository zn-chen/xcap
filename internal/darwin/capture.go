//go:build darwin

package darwin

import "image"

// BGRAToRGBA converts BGRA pixel data to RGBA format
// It also handles row alignment (bytes_per_row may be larger than width * 4)
func BGRAToRGBA(data []byte, width, height, bytesPerRow uint32) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	srcOffset := 0
	dstOffset := 0
	rowBytes := int(width * 4)

	for y := uint32(0); y < height; y++ {
		// Copy one row, handling alignment
		for x := 0; x < rowBytes; x += 4 {
			// BGRA -> RGBA: swap B and R
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

// CaptureResultToImage converts a CaptureResult to an image.RGBA
func CaptureResultToImage(result *CaptureResult) *image.RGBA {
	return BGRAToRGBA(result.Data, result.Width, result.Height, result.BytesPerRow)
}
