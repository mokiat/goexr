package exr

import (
	"encoding/binary"
	"fmt"
	"io"
)

func ReadBox2i(in io.Reader, target *Box2i) error {
	if err := binary.Read(in, binary.LittleEndian, &target.XMin); err != nil {
		return fmt.Errorf("error reading min x: %w", err)
	}
	if err := binary.Read(in, binary.LittleEndian, &target.YMin); err != nil {
		return fmt.Errorf("error reading min y: %w", err)
	}
	if err := binary.Read(in, binary.LittleEndian, &target.XMax); err != nil {
		return fmt.Errorf("error reading max x: %w", err)
	}
	if err := binary.Read(in, binary.LittleEndian, &target.YMax); err != nil {
		return fmt.Errorf("error reading max y: %w", err)
	}
	return nil
}

type Box2i struct {
	XMin int32
	YMin int32
	XMax int32
	YMax int32
}

func (b Box2i) Width() int32 {
	return b.XMax - b.XMin + 1
}

func (b Box2i) Height() int32 {
	return b.YMax - b.YMin + 1
}

func (b Box2i) Contains(other Box2i) bool {
	return other.XMin >= b.XMin &&
		other.XMax <= b.XMax &&
		other.YMin >= b.YMin &&
		other.YMax <= b.YMax
}
