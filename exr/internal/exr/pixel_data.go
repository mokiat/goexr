package exr

import (
	"fmt"
	"io"

	"github.com/x448/float16"
)

type PixelData interface {
	ReadLine(in io.Reader, y int32) error
	Float32(x, y int) float32
}

func NewNopPixelData(value float32) PixelData {
	return &nopPixelData{
		value: value,
	}
}

type nopPixelData struct {
	value float32
}

func (d *nopPixelData) ReadLine(in io.Reader, y int32) error {
	return fmt.Errorf("cannot read into nop pixel data")
}

func (d *nopPixelData) Float32(x, y int) float32 {
	return d.value
}

func NewFloat16PixelData(window Box2i, xSampling, ySampling int32) PixelData {
	width := window.Width() / xSampling
	height := window.Height() / ySampling
	return &float16PixelData{
		window:    window,
		xSampling: xSampling,
		ySampling: ySampling,
		pixels:    make([]float16.Float16, width*height),
	}
}

type float16PixelData struct {
	window    Box2i
	xSampling int32
	ySampling int32
	pixels    []float16.Float16
}

func (d *float16PixelData) ReadLine(in io.Reader, y int32) error {
	width := d.window.Width() / d.xSampling
	y = (y - d.window.YMin) / d.ySampling
	offset := y * width
	if err := Read(in, d.pixels[offset:offset+width:offset+width]); err != nil {
		return fmt.Errorf("error reading pixel slice: %w", err)
	}
	return nil
}

func (d *float16PixelData) Float32(x, y int) float32 {
	offX := (int32(x) - d.window.XMin) / d.xSampling
	offY := (int32(y) - d.window.YMin) / d.ySampling
	width := d.window.Width() / d.xSampling

	value := d.pixels[offX+width*offY]
	if value.IsInf(0) {
		value = float16.Frombits(uint16(0x7bff)) // max value
	}
	if value.IsNaN() {
		value = float16.Frombits(uint16(0x0000)) // min value
	}
	return value.Float32()
}
