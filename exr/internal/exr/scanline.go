package exr

import (
	"bytes"
	"fmt"
	"io"
)

func ReadScanLineBlock(in io.Reader, dataWindow Box2i, compression Compression, decompressor Decompressor, dataChannels []PixelData) error {
	var yCoordinate int32
	if err := Read(in, &yCoordinate); err != nil {
		return fmt.Errorf("error reading block y coordinate: %w", err)
	}

	var dataSize int32
	if err := Read(in, &dataSize); err != nil {
		return fmt.Errorf("error reading block data size: %w", err)
	}

	buffer := &bytes.Buffer{}
	if _, err := io.CopyN(buffer, in, int64(dataSize)); err != nil {
		return fmt.Errorf("error reading block data: %w", err)
	}

	blockHeight := int32(compression.LineCount())
	if dataWindow.YMax-yCoordinate+1 < blockHeight {
		blockHeight = dataWindow.YMax - yCoordinate + 1
	}

	if compression == CompressionZIP {
		uncompressedSize := int32(0)
		for _, dataChannel := range dataChannels {
			uncompressedSize += dataChannel.LineSize()
		}
		uncompressedSize *= blockHeight

		if uncompressedSize > dataSize {
			var err error
			buffer, err = decompressor.Decompress(buffer)
			if err != nil {
				return fmt.Errorf("error decompressing block data: %w", err)
			}
		}
	}

	for y := yCoordinate; y < yCoordinate+blockHeight; y++ {
		for _, dataChannel := range dataChannels {
			if err := dataChannel.ReadLine(buffer, y); err != nil {
				return fmt.Errorf("error reading scan line: %w", err)
			}
		}
	}
	return nil
}
