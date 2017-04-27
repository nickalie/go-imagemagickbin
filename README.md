# ImageMagick wrapper for Golang

[![](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/nickalie/go-imagemagickbin)

## Install

```go get -u github.com/nickalie/go-imagemagickbin```

## Example of usage

```go
package main

import (
	"github.com/nickalie/go-imagemagickbin"
	"image"
	"image/color"
)

func main() {
	const width, height = 256, 256

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8((x + y) & 255),
				G: uint8((x + y) << 1 & 255),
				B: uint8((x + y) << 2 & 255),
				A: 255,
			})
		}
	}

	// Convert generated image to webp, jpg, png, bmp
	m := imagemagickbin.NewMagick().InputImage(img)
	m.OutputFile("image.webp").Run()
	m.OutputFile("image.jpg").Run()
	m.OutputFile("image.png").Run()
	m.OutputFile("image.bmp").Run()
}

```

Library uses official Windows and macOS imagemagick distributions.
To make it work on other platforms make sure imagemagick is installed and *convert* available through command line.

Refer to [docs](https://godoc.org/github.com/nickalie/go-imagemagickbin) for more information.
