package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/mokiat/goexr/exr"
)

func init() {
	flag.Usage = func() {
		cl := flag.CommandLine
		fmt.Fprintf(cl.Output(), "Usage:\texrtopng <source.exr> <target.png>\n")
		cl.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	source, target := flag.Arg(0), flag.Arg(1)

	if err := runApp(source, target); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(2)
	}
}

func runApp(source, target string) error {
	img, err := openEXR(source)
	if err != nil {
		return err
	}
	return savePNG(target, img)
}

func openEXR(location string) (image.Image, error) {
	file, err := os.Open(location)
	if err != nil {
		return nil, fmt.Errorf("error opening file %q: %w", location, err)
	}
	defer file.Close()

	img, err := exr.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding exr image: %w", err)
	}
	return img, nil
}

func savePNG(location string, img image.Image) error {
	file, err := os.Create(location)
	if err != nil {
		return fmt.Errorf("error creating file %q: %w", location, err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("error encoding png image: %w", err)
	}
	return nil
}
