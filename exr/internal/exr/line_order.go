package exr

import (
	"fmt"
	"io"
)

func ReadLineOrder(in io.Reader, target *LineOrder) error {
	return Read(in, target)
}

const (
	LineOrderIncreasingY LineOrder = 0
	LineOrderDecreasingY LineOrder = 1
	LineOrderRandomY     LineOrder = 2
)

type LineOrder uint8

func (o LineOrder) String() string {
	switch o {
	case LineOrderIncreasingY:
		return "INCREASING_Y"
	case LineOrderDecreasingY:
		return "DECREASING_Y"
	case LineOrderRandomY:
		return "RANDOM_Y"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", o)
	}
}
