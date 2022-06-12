package exr

import (
	"fmt"
	"image"
	"io"

	"github.com/mokiat/goexr/exr/internal/exr"
)

func init() {
	image.RegisterFormat(exr.Extension, string(exr.MagicSequence[:]), Decode, DecodeConfig)
}

func DecodeConfig(in io.Reader) (image.Config, error) {
	var magic exr.Magic
	if err := exr.ReadMagic(in, &magic); err != nil {
		return image.Config{}, fmt.Errorf("error reading magic: %w", err)
	}
	if !magic.IsCorrect() {
		return image.Config{}, fmt.Errorf("incorrect magic sequence \"0x%x\"", magic)
	}

	var version exr.Version
	if err := exr.ReadVersion(in, &version); err != nil {
		return image.Config{}, fmt.Errorf("error reading version: %w", err)
	}
	if version.Number() != exr.SupportedVersion {
		return image.Config{}, fmt.Errorf("unsupported version %d", version.Number())
	}

	var header exr.Header
	if err := exr.ReadHeader(in, &header); err != nil {
		return image.Config{}, fmt.Errorf("error reading header: %w", err)
	}

	return image.Config{
		ColorModel: RGBAModel,
		Width:      int(header.DisplayWindow.Width()),
		Height:     int(header.DisplayWindow.Height()),
	}, nil
}

func Decode(in io.Reader) (image.Image, error) {
	var magic exr.Magic
	if err := exr.ReadMagic(in, &magic); err != nil {
		return nil, fmt.Errorf("error reading magic: %w", err)
	}
	if !magic.IsCorrect() {
		return nil, fmt.Errorf("incorrect magic sequence \"0x%x\"", magic)
	}

	var version exr.Version
	if err := exr.ReadVersion(in, &version); err != nil {
		return nil, fmt.Errorf("error reading version: %w", err)
	}
	if version.Number() != exr.SupportedVersion {
		return nil, fmt.Errorf("unsupported version %d", version.Number())
	}
	if version.HasFlag(exr.FlagSingleTile) {
		return nil, fmt.Errorf("tiled format not supported")
	}
	if version.HasFlag(exr.FlagNonImage) {
		return nil, fmt.Errorf("deep data not supported")
	}
	if version.HasFlag(exr.FlagMultipart) {
		return nil, fmt.Errorf("multipart not supported")
	}

	var header exr.Header
	if err := exr.ReadHeader(in, &header); err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	return nil, fmt.Errorf("TODO")
}
