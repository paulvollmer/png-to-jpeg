package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Version of the commandline tool
const Version = "0.2.0"

var quality int
var recursive bool

func usage() {
	fmt.Printf("Usage: png-to-jpeg\n\n")
	fmt.Print("Flags:\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	flag.IntVar(&quality, "q", 100, "JPEG output quality")
	flag.BoolVar(&recursive, "r", false, "if the input is a directory, run the process recursive")
	version := flag.Bool("v", false, "Print the version and exit")
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Printf("png-to-jpeg v%s\n", Version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Missing input file/folder")
		os.Exit(1)
	}

	err := Process(args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// Process an input path
func Process(input string) error {
	inputType, err := os.Stat(input)
	if err != nil {
		return err
	}
	if inputType.IsDir() {
		files, err := ioutil.ReadDir(input)
		if err != nil {
			return err
		}
		for _, file := range files {
			err := Process(path.Join(input, file.Name()))
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}
		}
	} else {
		return ProcessImage(input)
	}
	return nil
}

// ProcessImage read and convert the given input PNG to a JPEG
func ProcessImage(input string) error {
	fmt.Printf("==> %q ", input)
	ext := filepath.Ext(input)
	switch ext {
	case ".png", ".PNG":
		img, err := ImageRead(input)
		if err != nil {
			return err
		}
		fmt.Printf("converted to JPEG\n")
		outFilepath := strings.TrimSuffix(input, ext) + ".jpeg"
		err = FormatPNG(outFilepath, img, quality)
		if err != nil {
			return err
		}
		break
	case ".jpg", ".JPG", ".jpeg", ".JPEG":
		fmt.Printf("skip JPEG file\n")
		break
	default:
		return errors.New("File format not supported")
	}
	return nil
}

// ImageRead the given input filepath and return an Image interface
func ImageRead(input string) (image.Image, error) {
	var img image.Image
	var err error
	file, err := os.Open(input)
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

// FormatPNG format the given image and save it to the source directory
func FormatPNG(src string, img image.Image, quality int) error {
	out, err := os.Create(src)
	if err != nil {
		return err
	}
	opt := jpeg.Options{quality}
	return jpeg.Encode(out, img, &opt)
}
