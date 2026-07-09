// SPDX-FileCopyrightText: 2026 Renan Fernandes
//
// SPDX-License-Identifier: EUPL-1.2

package escpos

import (
	"image"
	"image/color"
)

// PrintImageRaster prints an image using GS v 0 raster mode.
func (e *Escpos) PrintImageRaster(img image.Image) (int, error) {
	b := imageToRaster(img)
	return e.WriteRaw(b)
}

// PrintImageBitImage prints an image using ESC * 24-dot bit image mode.
// This mode is usually more compatible with POS-80 clone printers than GS v 0.
func (e *Escpos) PrintImageBitImage(img image.Image) (int, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return 0, nil
	}

	var out []byte
	for y := 0; y < height; y += 24 {
		out = append(out, esc, '*', 33, byte(width&0xff), byte((width>>8)&0xff))
		for x := 0; x < width; x++ {
			for block := 0; block < 3; block++ {
				var b byte
				for bit := 0; bit < 8; bit++ {
					yy := y + block*8 + bit
					if yy < height && isBlack(img.At(bounds.Min.X+x, bounds.Min.Y+yy)) {
						b |= 1 << uint(7-bit)
					}
				}
				out = append(out, b)
			}
		}
		out = append(out, '\n')
	}
	return e.WriteRaw(out)
}

func imageToRaster(img image.Image) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	widthBytes := (width + 7) / 8

	data := make([]byte, widthBytes*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if isBlack(img.At(bounds.Min.X+x, bounds.Min.Y+y)) {
				data[y*widthBytes+x/8] |= 1 << uint(7-(x%8))
			}
		}
	}

	out := []byte{gs, 'v', '0', 0, byte(widthBytes & 0xff), byte((widthBytes >> 8) & 0xff), byte(height & 0xff), byte((height >> 8) & 0xff)}
	return append(out, data...)
}

func isBlack(c color.Color) bool {
	r, g, b, a := c.RGBA()
	if a == 0 {
		return false
	}
	// 16-bit luminance, with alpha blended over white.
	alpha := a >> 8
	r8 := (r>>8)*alpha/255 + 255*(255-alpha)/255
	g8 := (g>>8)*alpha/255 + 255*(255-alpha)/255
	b8 := (b>>8)*alpha/255 + 255*(255-alpha)/255
	lum := (299*r8 + 587*g8 + 114*b8) / 1000
	return lum < 128
}
