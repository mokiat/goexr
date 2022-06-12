package exr

import (
	"encoding/binary"
	"io"
)

var (
	order = binary.LittleEndian
)

func Read(in io.Reader, data any) error {
	return binary.Read(in, order, data)
}

func ReadNullTerminatedString(in io.Reader, target *string) error {
	var buffer []byte
	for {
		var char byte
		if err := Read(in, &char); err != nil {
			return err
		}
		if char == 0x00 {
			break
		}
		buffer = append(buffer, char)
	}
	*target = string(buffer)
	return nil
}
