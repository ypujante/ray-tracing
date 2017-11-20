package main

import "math"

// Color defines the basic Red/Green/Blue as raw float64 values
type Color struct {
	R, G, B float64
}

// PixelValue converts a raw Color into a pixel value (0-255) packed into a uint32
func (color *Color) PixelValue() uint32 {
	r := uint32(math.Min(255.0, color.R*255.99))
	g := uint32(math.Min(255.0, color.G*255.99))
	b := uint32(math.Min(255.0, color.B*255.99))

	return ((r & 0xFF) << 16) | ((g & 0xFF) << 8) | (b & 0xFF)
}
