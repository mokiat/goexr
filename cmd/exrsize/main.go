package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/mokiat/goexr/exr"
)

func init() {
	flag.Usage = func() {
		cl := flag.CommandLine
		fmt.Fprintf(cl.Output(), "Usage:\texrsize <source.exr>\n")
		cl.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	source := flag.Arg(0)

	if err := runApp(source); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(2)
	}
}

func runApp(source string) error {
	cfg, err := openEXRConfig(source)
	if err != nil {
		return err
	}
	fmt.Printf("%d x %d\n", cfg.Width, cfg.Height)
	return nil
}

func openEXRConfig(location string) (image.Config, error) {
	file, err := os.Open(location)
	if err != nil {
		return image.Config{}, fmt.Errorf("error opening file %q: %w", location, err)
	}
	defer file.Close()

	cfg, err := exr.DecodeConfig(file)
	if err != nil {
		return image.Config{}, fmt.Errorf("error decoding exr image config: %w", err)
	}
	return cfg, nil
}
