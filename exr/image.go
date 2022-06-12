package exr

import (
	"image"
	"image/color"

	"github.com/mokiat/goexr/exr/internal/exr"
)

// RGBAImage represents an EXR image that consists of R, G, B, and A components.
//
// Even if the original image that is loaded does not contain all of the
// components, default ones will be assigned.
type RGBAImage struct {
	rect     image.Rectangle
	channelR exr.PixelData
	channelG exr.PixelData
	channelB exr.PixelData
	channelA exr.PixelData
}

// ColorModel returns the RGBAImage's color model.
func (i *RGBAImage) ColorModel() color.Model {
	return RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (i *RGBAImage) Bounds() image.Rectangle {
	return i.rect
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
//
// The returned color is of type RGBAColor which can be used to acquire the
// linear (float) components of the color.
func (i *RGBAImage) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(i.rect)) {
		return RGBAColor{}
	}
	return RGBAColor{
		R: i.channelR.Float32(x, y),
		G: i.channelG.Float32(x, y),
		B: i.channelB.Float32(x, y),
		A: i.channelA.Float32(x, y),
	}
}
