package exr

import "io"

const (
	SupportedVersion = 2
)

func ReadVersion(in io.Reader, target *Version) error {
	return Read(in, target)
}

type Version int32

func (v Version) Number() int {
	return int(v & 0xFF)
}

func (v Version) HasFlag(flag Flag) bool {
	return int32(v)&int32(flag) == int32(flag)
}

type Flag int32

const (
	FlagSingleTile Flag = 1 << 8  // one at 9-th bit in version
	FlagLongName   Flag = 1 << 9  // one at 10-th bit in version
	FlagNonImage   Flag = 1 << 10 // one at 11-th bit in version
	FlagMultipart  Flag = 1 << 11 // one at 12-th bit in version
)
