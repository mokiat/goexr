# goexr

Go library for parsing OpenEXR files.

Not all EXR files are supported at the moment. Make sure to check the
[Limitations](#Limitations) section.

## Usage

Add the library as a dependency to your project.

```sh
go get github.com/mokiat/goexr
```

Parse the image as with any other format. For example:

```go
package main

import (
	"fmt"
	"image"
	"os"

	_ "github.com/mokiat/goexr"
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
  _ "github.com/mokiat/goexr"
)
```

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

- Krita
- GIMP
