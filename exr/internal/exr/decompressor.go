package exr

import (
	"bytes"
	"compress/zlib"
	"io"
)

type Decompressor interface {
	Decompress(src *bytes.Buffer) (*bytes.Buffer, error)
}

func NewNopDecompressor() Decompressor {
	return &nopDecompressor{}
}

type nopDecompressor struct{}

func (d *nopDecompressor) Decompress(src *bytes.Buffer) (*bytes.Buffer, error) {
	return src, nil
}

func NewZipDecompressor() Decompressor {
	return &zipDecompressor{}
}

type zipDecompressor struct{}

func (d *zipDecompressor) Decompress(src *bytes.Buffer) (*bytes.Buffer, error) {
	zlibIn, err := zlib.NewReader(src)
	if err != nil {
		return nil, err
	}
	out := &bytes.Buffer{}
	if _, err := io.Copy(out, zlibIn); err != nil {
		return nil, err
	}
	if err := zlibIn.Close(); err != nil {
		return nil, err
	}

	data := out.Bytes()

	// reconstruct scalar
	for i := 1; i < len(data); i++ {
		v := int(data[i-1]) + int(data[i]) - 128
		data[i] = byte(v)
	}

	// interleave scalar
	result := make([]byte, len(data))
	i1 := 0
	i2 := (len(data) + 1) / 2
	j := 0
	for j < len(result) {
		result[j] = data[i1]
		j++
		i1++

		if j >= len(result) {
			break
		}
		result[j] = data[i2]
		j++
		i2++
	}

	return bytes.NewBuffer(result), nil
}
