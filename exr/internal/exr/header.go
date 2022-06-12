package exr

import (
	"bytes"
	"fmt"
	"io"
)

func ReadHeader(in io.Reader, target *Header) error {
	for {
		var attributeName AttributeName
		if err := ReadAttributeName(in, &attributeName); err != nil {
			return fmt.Errorf("error reading attribute name: %w", err)
		}
		if attributeName == "" {
			return nil
		}

		var attributeType AttributeType
		if err := ReadAttributeType(in, &attributeType); err != nil {
			return fmt.Errorf("error reading attribute type: %w", err)
		}

		var attributeSize int32
		if err := Read(in, &attributeSize); err != nil {
			return fmt.Errorf("error reading attribute size: %w", err)
		}

		attributeValue := make([]byte, attributeSize)
		if err := Read(in, &attributeValue); err != nil {
			return fmt.Errorf("error reading attribute value: %w", err)
		}

		switch attributeName {
		case AttributeNameChannels:
			if attributeType != AttributeTypeChannelList {
				return fmt.Errorf("incorrect channels attribute type %q", attributeType)
			}
			if err := ReadChannelList(bytes.NewReader(attributeValue), &target.Channels); err != nil {
				return fmt.Errorf("error reading channels: %w", err)
			}

		case AttributeNameCompression:
			if attributeType != AttributeTypeCompression {
				return fmt.Errorf("incorrect compression attribute type %q", attributeType)
			}
			if err := ReadCompression(bytes.NewReader(attributeValue), &target.Compression); err != nil {
				return fmt.Errorf("error reading compression: %w", err)
			}

		case AttributeNameDataWindow:
			if attributeType != AttributeTypeBox2i {
				return fmt.Errorf("incorrect data window attribute type %q", attributeType)
			}
			if err := ReadBox2i(bytes.NewReader(attributeValue), &target.DataWindow); err != nil {
				return fmt.Errorf("error reading data window: %w", err)
			}

		case AttributeNameDisplayWindow:
			if attributeType != AttributeTypeBox2i {
				return fmt.Errorf("incorrect display window attribute type %q", attributeType)
			}
			if err := ReadBox2i(bytes.NewReader(attributeValue), &target.DisplayWindow); err != nil {
				return fmt.Errorf("error reading display window: %w", err)
			}

		case AttributeNameLineOrder:
			if attributeType != AttributeTypeLineOrder {
				return fmt.Errorf("incorrect line order attribute type %q", attributeType)
			}
			if err := ReadLineOrder(bytes.NewReader(attributeValue), &target.LineOrder); err != nil {
				return fmt.Errorf("error reading line order: %w", err)
			}

		default:
			// Skip unknown / unnecessary attributes
		}
	}
}

type Header struct {
	Channels      ChannelList
	Compression   Compression
	DataWindow    Box2i
	DisplayWindow Box2i
	LineOrder     LineOrder
}
