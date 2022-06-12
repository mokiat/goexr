package exr

import (
	"fmt"
	"io"
)

func ChunkCount(dataWindow Box2i, compression Compression) int {
	lineCount := compression.LineCount()
	return (int(dataWindow.YMax-dataWindow.YMin) + lineCount) / lineCount
}

func ReadOffsets(in io.Reader, chunkCount int) error {
	var lastOffset uint64
	for i := 0; i < chunkCount; i++ {
		var offset uint64
		if err := Read(in, &offset); err != nil {
			return fmt.Errorf("error reading offset: %w", err)
		}
		if offset < lastOffset {
			return fmt.Errorf("non-incrementing chunk offsets")
		}
		lastOffset = offset
	}
	return nil
}
