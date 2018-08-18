package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

const Version = "0.1.0"

func usage() {
	fmt.Printf("Usage: png-to-jpeg\n\n")
	fmt.Print("Flags:\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	input := flag.String("input", "", "Input image file")
	quality := flag.Int("quality", 100, "JPEG output quality")
	version := flag.Bool("version", false, "Print the version and exit")
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Printf("png-to-jpeg v%s\n", Version)
		os.Exit(0)
	}
	if *input == "" {
		fmt.Println("Missing input file")
		os.Exit(1)
	}
	fmt.Printf("==> Input image %q\n", *input)
	ext := filepath.Ext(*input)
	switch ext {
	case ".png", ".PNG":
		fmt.Println("==> Reading image data...")
		img, err := ImageRead(*input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("==> Converting image to JPEG...")
		outFilepath := strings.TrimSuffix(*input, ext) + ".jpeg"
		err = FormatPNG(outFilepath, img, *quality)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Println("File format not supported")
		os.Exit(1)
	}
}

func ImageRead(ImageFile string) (image.Image, error) {
	var img image.Image
	var err error
	file, err := os.Open(ImageFile)
	if err != nil {
		return img, err
	}
	img, err = png.Decode(file)
	if err != nil {
		return img, err
	}
	file.Close()
	return img, nil
}

func FormatPNG(src string, img image.Image, quality int) error {
	out, err := os.Create(src)
	if err != nil {
		return err
	}
	opt := jpeg.Options{quality}
	return jpeg.Encode(out, img, &opt)
}
