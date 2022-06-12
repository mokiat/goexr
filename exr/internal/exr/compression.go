package exr

import (
	"fmt"
	"io"
)

func ReadCompression(in io.Reader, target *Compression) error {
	return Read(in, target)
}

const (
	CompressionNone  Compression = 0
	CompressionRLE   Compression = 1
	CompressionZIPS  Compression = 2
	CompressionZIP   Compression = 3
	CompressionPIZ   Compression = 4
	CompressionPXR24 Compression = 5
	CompressionB44   Compression = 6
	CompressionB44A  Compression = 7
)

type Compression uint8

func (c Compression) LineCount() int {
	switch c {
	case CompressionNone:
		return 1
	case CompressionRLE:
		return 1
	case CompressionZIPS:
		return 1
	case CompressionZIP:
		return 16
	case CompressionPIZ:
		return 32
	case CompressionPXR24:
		return 16
	case CompressionB44:
		return 32
	case CompressionB44A:
		return 32
	default:
		panic(fmt.Errorf("unknown compression type %d", c))
	}
}

func (c Compression) String() string {
	switch c {
	case CompressionNone:
		return "NONE"
	case CompressionRLE:
		return "RLE"
	case CompressionZIPS:
		return "ZIPS"
	case CompressionZIP:
		return "ZIP"
	case CompressionPIZ:
		return "PIZ"
	case CompressionPXR24:
		return "PXR24"
	case CompressionB44:
		return "B44"
	case CompressionB44A:
		return "B44A"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", c)
	}
}
