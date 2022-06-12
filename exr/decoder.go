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

// DecodeConfig returns the color model and dimensions of an EXR image without
// decoding the entire image.
//
// This function supports all version 2 EXR images.
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

// Decode reads an EXR image from in and returns it as an image.Image.
// The type of the Image is RGBAImage.
//
// Only a limited set of EXR image types are supported at the moment.
// The main restrictions are as follows, though others apply as well:
//
// 	- They have to be single-part scan line images.
// 	- They have to use no compression or zip compression.
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

	dataWindow := header.DataWindow
	if dataWindow.Width() <= 0 || dataWindow.Height() <= 0 {
		return nil, fmt.Errorf("invalid data window size (%d x %d)", dataWindow.Width(), dataWindow.Height())
	}

	displayWindow := header.DisplayWindow
	if !dataWindow.Contains(displayWindow) {
		return nil, fmt.Errorf("invalid display window: not contained by data window")
	}

	lineOrder := header.LineOrder
	if lineOrder != exr.LineOrderIncreasingY {
		return nil, fmt.Errorf("unsupported line order %q", lineOrder)
	}

	var (
		decompressor exr.Decompressor
	)

	compression := header.Compression
	switch compression {
	case exr.CompressionNone:
		decompressor = exr.NewNopDecompressor()
	case exr.CompressionZIP:
		decompressor = exr.NewZipDecompressor()
	default:
		return nil, fmt.Errorf("unsupported compression %q", compression)
	}

	var (
		dataChannelR = exr.NewNopPixelData(0.0)
		dataChannelG = exr.NewNopPixelData(0.0)
		dataChannelB = exr.NewNopPixelData(0.0)
		dataChannelA = exr.NewNopPixelData(1.0)
	)

	dataChannels := make([]exr.PixelData, len(header.Channels))
	for i, channel := range header.Channels {
		switch channel.PixelType {
		case exr.PixelTypeUint:
			dataChannels[i] = exr.NewUint32PixelData(dataWindow, channel.XSampling, channel.YSampling)
		case exr.PixelTypeHalf:
			dataChannels[i] = exr.NewFloat16PixelData(dataWindow, channel.XSampling, channel.YSampling)
		case exr.PixelTypeFloat:
			dataChannels[i] = exr.NewFloat32PixelData(dataWindow, channel.XSampling, channel.YSampling)
		default:
			return nil, fmt.Errorf("unsupported channel pixel type %q", channel.PixelType)
		}
		switch channel.Name {
		case "R":
			dataChannelR = dataChannels[i]
		case "G":
			dataChannelG = dataChannels[i]
		case "B":
			dataChannelB = dataChannels[i]
		case "A":
			dataChannelA = dataChannels[i]
		}
	}

	chunkCount := exr.ChunkCount(dataWindow, compression)

	if err := exr.ReadOffsets(in, chunkCount); err != nil {
		return nil, fmt.Errorf("error reading offsets: %w", err)
	}

	for i := 0; i < chunkCount; i++ {
		if err := exr.ReadScanLineBlock(in, dataWindow, compression, decompressor, dataChannels); err != nil {
			return nil, fmt.Errorf("error reading scan line block: %w", err)
		}
	}

	return &RGBAImage{
		rect: image.Rect(
			int(displayWindow.XMin), int(displayWindow.YMin),
			int(displayWindow.XMax+1), int(displayWindow.YMax+1),
		),
		channelR: dataChannelR,
		channelG: dataChannelG,
		channelB: dataChannelB,
		channelA: dataChannelA,
	}, nil
}
