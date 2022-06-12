package exr

import "io"

func ReadAttributeName(in io.Reader, target *AttributeName) error {
	return ReadNullTerminatedString(in, target)
}

const (
	AttributeNameChannels           AttributeName = "channels"
	AttributeNameCompression        AttributeName = "compression"
	AttributeNameDataWindow         AttributeName = "dataWindow"
	AttributeNameDisplayWindow      AttributeName = "displayWindow"
	AttributeNameLineOrder          AttributeName = "lineOrder"
	AttributeNamePixelAspectRatio   AttributeName = "pixelAspectRatio"
	AttributeNameScreenWindowCenter AttributeName = "screenWindowCenter"
	AttributeNameScreenWindowWidth  AttributeName = "screenWindowWidth"
)

type AttributeName string

func ReadAttributeType(in io.Reader, target *AttributeType) error {
	return ReadNullTerminatedString(in, target)
}

const (
	AttributeTypeChannelList AttributeType = "chlist"
	AttributeTypeCompression AttributeType = "compression"
	AttributeTypeBox2i       AttributeType = "box2i"
	AttributeTypeLineOrder   AttributeType = "lineOrder"
	AttributeTypeFloat       AttributeType = "float"
	AttributeTypeV2f         AttributeType = "v2f"
)

type AttributeType string
