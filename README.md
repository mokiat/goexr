# goexr

Go library for parsing OpenEXR files.

Not all EXR files are supported at the moment. Make sure to check the
[Limitations](#limitations) section.

## Usage

Add the library as a dependency to your project.

```sh
go get github.com/mokiat/goexr/exr@latest
```

Parse the image as with any other format. For example:

```go
package main

import (
	"fmt"
	"image"
	"os"

	_ "github.com/mokiat/goexr/exr"
)

func main() {
	file, err := os.Open("example.exr")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bounds [%s]: %v\n", format, img.Bounds())
}
```

Make sure that you have added an anonymous import to the library which ensures
that there is a decoder registered with the `image` package.

```go
import (
  _ "github.com/mokiat/goexr/exr"
)
```

For more information check the Go documentation of the `exr` package.

## Limitations

The library supports a subset of EXR files.

Supported EXR versions:

- `2`

Supported image types:

- `single part scanline`

Supported compression modes:

- `NO_COMPRESSION`
- `ZIP_COMPRESSION`

Supported channels:

- `R`
- `G`
- `B`
- `A`

Supported channel formats:

- `HALF`
- `FLOAT`

At the time of writing, EXR images produced by the following programs appear
to be supported:

- [Krita](https://krita.org/)
- [GIMP](https://www.gimp.org/)

## Tools

Aside from the `exr` package, this library provides two tools that can be used
to check whether this project works on specific EXR files.

### exrsize

This tool prints the pixel size of an EXR image.

```sh
exrsize <img.exr>
```

It can be install as follows:

```sh
go install github.com/mokiat/goexr/cmd/exrsize@latest
```

It should work with all version 2 EXR images, even those that cannot be
decoded by the `exr` package.

### exrtopng

This tool converts an `exr` image into a `png` one by using basic tone mapping
and gamma correction to convert linear colors into sRGB space.

```sh
exrtopng <src.exr> <dst.png>
```

It can be install as follows:

```sh
go install github.com/mokiat/goexr/cmd/exrtopng@latest
```

This tool is subject to the above [Limitations](#limitations).

## Resources

For the development of this decoded the following resources were useful:

- [Online Documentation](https://openexr.readthedocs.io/en/latest/TechnicalIntroduction.html)

		Very easy to navigate and contains information on the core concepts.

- [PDF Documentation](https://www.openexr.com/documentation/openexrfilelayout.pdf)

		The Online Documentation appeared to have errors in some places and the PDF variant was useful in such cases.

- [Source Code](https://github.com/AcademySoftwareFoundation/openexr/tree/main/src/lib/OpenEXR)

		Some aspects of OpenEXR are not fully documented and in such cases the only place where answers can be found is the reference implementation. In fact, even the documentation refers to it at times.
