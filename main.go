package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

var (
	maxWidth  uint
	maxHeight uint
	max       uint
)

func init() {
	flag.UintVar(&maxWidth, "w", 0, "Constrain width of output image to this many pixels.")
	flag.UintVar(&maxHeight, "h", 0, "Constrain height of output image to this many pixels.")
	flag.UintVar(&max, "s", 0, "Constrain size (height and width) of output image to this many pixels.")
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		exit(nil)
	}

	filename := flag.Arg(0)
	if max > 0 && (maxWidth > 0 || maxHeight > 0) {
		exit(fmt.Errorf("-s cannot be combined with -h or -w."))
	}

	if max > 0 {
		maxWidth = max
		maxHeight = max
	}

	f, err := os.Open(filename)
	if err != nil {
		exit(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		exit(err)
	}

	if maxHeight > 0 || maxWidth > 0 {
		if maxHeight == 0 {
			maxHeight = uint(img.Bounds().Size().Y)
		}
		if maxWidth == 0 {
			maxWidth = uint(img.Bounds().Size().X)
		}
		img = resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)
	}
	out := flag.Arg(1)

	ext := filepath.Ext(out)

	switch ext {
	case ".gif", ".jpeg", ".jpg", ".png":
	default:
		err = fmt.Errorf("Unsupported output file format: %q", ext)
	}

	thumb, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		exit(err)
	}
	defer thumb.Close()

	switch ext {
	case ".gif":
		err = gif.Encode(thumb, img, nil)
	case ".jpeg", ".jpg":
		err = jpeg.Encode(thumb, img, nil)
	case ".png":
		err = png.Encode(thumb, img)
	default:
		panic(fmt.Errorf("Should be impossible: unsupported output file format: %q", ext))
	}
	if err != nil {
		exit(err)
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Fprintln(os.Stderr, `
convert creates a copy of the given image file, changing its format as indicated
by the file's extension and optionally downscaling it. Only supports gif, 
jpg (jpeg), or png input and output.
`)
	fmt.Fprintln(os.Stderr, "Usage:\n\n  convert [options] <image_filename> <out_filename>\n\nOptions:\n")

	flag.VisitAll(func(f *flag.Flag) {
		format := "  -%s: %s\n"
		fmt.Fprintf(os.Stderr, format, f.Name, f.Usage)
	})

	fmt.Fprintln(os.Stderr, "\nExample: (convert from png to jpg and constrain to max 1024 on either side)\n\n  convert foo.png foo.jpg -s 1024")
	os.Exit(1)
}
