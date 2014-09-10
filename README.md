convert
=======

An application to convert and resize images written in Go.

``` 
convert creates a copy of the given image file, changing its format as indicated
by the file's extension and optionally downscaling it. Only supports gif, 
jpg (jpeg), or png input and output.

Usage:

  convert [options] <image_filename> <out_filename>

Options:

  -h: Constrain height of output image to this many pixels.
  -s: Constrain size (height and width) of output image to this many pixels.
  -w: Constrain width of output image to this many pixels.

Example: (convert from png to jpg and constrain to max 1024 on either side)

  convert -s 1024 foo.png foo.jpg

```