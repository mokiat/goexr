package exr

import (
	"image/color"
	"math"
)

const (
	gammaFactor = 1.0 / 2.2
)

var (
	// RGBAModel returns the color.Model for RGBAColor colors.
	RGBAModel color.Model = color.ModelFunc(rgbaModel)
)

// RGBAColor represents a linear EXR color that implements the color.Color
// interface and is composed of R, G, B, and A components.
type RGBAColor struct {

	// R holds the amount of red in this color.
	R float32

	// G holds the amount of green in this color.
	G float32

	// B holds the amount of blue in this color.
	B float32

	// A holds the amount of alpha in this color.
	A float32
}

// RGBA returns the alpha-premultiplied red, green, blue and alpha values
// for the color. Each value ranges within [0, 0xffff], but is represented
// by a uint32 so that multiplying by a blend factor up to 0xffff will not
// overflow.
//
// An alpha-premultiplied color component c has been scaled by alpha (a),
// so has valid values 0 <= c <= a.
//
// Reinhard tone mapping and gamma correction are performed to convert the
// color into sRGB space.
func (c RGBAColor) RGBA() (r, g, b, a uint32) {
	// tone mapping
	floatR := float64(c.R / (c.R + 1.0))
	floatG := float64(c.G / (c.G + 1.0))
	floatB := float64(c.B / (c.B + 1.0))

	// gamma correction
	floatR = math.Pow(floatR, gammaFactor)
	floatG = math.Pow(floatG, gammaFactor)
	floatB = math.Pow(floatB, gammaFactor)

	// alpha pre-multiplication
	floatR *= float64(c.A)
	floatG *= float64(c.A)
	floatB *= float64(c.A)

	// uint32 conversion
	r = uint32(floatR*0xFFFF) & 0xFFFF
	g = uint32(floatG*0xFFFF) & 0xFFFF
	b = uint32(floatB*0xFFFF) & 0xFFFF
	a = uint32(c.A*0xFFFF) & 0xFFFF
	return
}

func rgbaModel(c color.Color) color.Color {
	if _, ok := c.(RGBAColor); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return RGBAColor{
		R: float32(r) / float32(0xFFFF),
		G: float32(g) / float32(0xFFFF),
		B: float32(b) / float32(0xFFFF),
		A: float32(a) / float32(0xFFFF),
	}
}
