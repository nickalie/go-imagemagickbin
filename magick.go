package imagemagickbin

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/nickalie/go-binwrapper"
	"image"
	"image/png"
	"io"
	"strconv"
	"strings"
)

// Magick is a wrapper for convert tool
type Magick struct {
	*binwrapper.BinWrapper
	inputFile    string
	inputImage   image.Image
	input        io.Reader
	outputFile   string
	output       io.Writer
	quality      int
	outputFormat string
	fuzz         uint
	trim         bool
}

// NewMagick creates new Magick instance.
func NewMagick() *Magick {
	return &Magick{
		BinWrapper:   createBinWrapper(),
		quality:      -1,
		outputFormat: "png",
	}
}

// InputFile sets image file to convert.
// Input or InputImage called before will be ignored.
func (c *Magick) InputFile(file string) *Magick {
	c.input = nil
	c.inputImage = nil
	c.inputFile = file
	return c
}

// Input sets reader to convert.
// InputFile or InputImage called before will be ignored.
func (c *Magick) Input(reader io.Reader) *Magick {
	c.inputFile = ""
	c.inputImage = nil
	c.input = reader
	return c
}

// InputImage sets image to convert.
// InputFile or Input called before will be ignored.
func (c *Magick) InputImage(img image.Image) *Magick {
	c.inputFile = ""
	c.input = nil
	c.inputImage = img
	return c
}

// OutputFile specify the name of the output image file.
// Output called before will be ignored.
func (c *Magick) OutputFile(file string) *Magick {
	c.output = nil
	c.outputFile = file
	return c
}

// Output specify writer to write image file content.
// OutputFile called before will be ignored.
func (c *Magick) Output(writer io.Writer) *Magick {
	c.outputFile = ""
	c.output = writer
	return c
}

// OutputFormat specifies output format of the image.
func (c *Magick) OutputFormat(format string) *Magick {
	c.outputFormat = format
	return c
}

// Fuzz uses a number of algorithms search for a target color.
// By default the color must be exact.
// Use this option to match colors that are close to the target color in RGB space.
// For example, if you want to automagically trim the edges of an image with Trim, but the image was scanned and the target background color may differ by a small amount.
// This option can account for these differences.
func (c *Magick) Fuzz(percents uint) *Magick {
	c.fuzz = percents
	return c
}

// Trim removes any edges that are exactly the same color as the corner pixels.
// Use Fuzz to make Trim remove edges that are nearly the same color as the corner pixels.
func (c *Magick) Trim(trim bool) *Magick {
	c.trim = trim
	return c
}

// GetTrimInfo returns trimmed rectangle without performin trim.
func (c *Magick) GetTrimInfo(fuzz uint, threshold int) (*image.Rectangle, error) {
	defer c.BinWrapper.Reset()
	c.setInput()
	err := c.Arg("-fuzz", fmt.Sprintf("%d%%", fuzz)).Arg("-trim").Arg("info:").Run()

	if err != nil {
		return nil, err
	}

	outputs := strings.Split(string(c.StdOut()), " ")
	var xStr string
	var yStr string
	var wStr string
	var hStr string
	var initialWStr string
	var initialHStr string

	for _, v := range outputs {
		if strings.Count(v, "+") == 2 && strings.Count(v, "x") == 1 {
			t := strings.Split(v, "+")
			xStr, yStr = t[1], t[2]
			t = strings.Split(t[0], "x")
			initialWStr, initialHStr = t[0], t[1]

		} else if strings.Count(v, "x") == 1 {
			t := strings.Split(v, "x")
			wStr, hStr = t[0], t[1]
		}
	}

	var x, y, width, height int

	x, err = strconv.Atoi(xStr)

	if err != nil {
		return nil, err
	}

	y, err = strconv.Atoi(yStr)

	if err != nil {
		return nil, err
	}

	width, err = strconv.Atoi(wStr)

	if err != nil {
		return nil, err
	}

	height, err = strconv.Atoi(hStr)

	if err != nil {
		return nil, err
	}

	initialWidth, err := strconv.Atoi(initialWStr)

	if err != nil {
		return nil, err
	}

	initialHeight, err := strconv.Atoi(initialHStr)

	if err != nil {
		return nil, err
	}

	if 100*(width*height)/(initialWidth*initialHeight) < threshold {
		return nil, errors.New("To much trimmed")
	}

	if x == 0 && y == 0 && width == initialWidth && height == initialHeight {
		return nil, errors.New("Nothing to trim")
	}

	result := image.Rect(x, y, x+width, y+height)
	return &result, nil
}

// Quality specifies quality to compress the image.
// quality is 1 (lowest image quality and highest compression) to 100 (best quality but least effective compression).
// The default is to use the estimated quality of your input image if it can be determined, otherwise 92.
func (c *Magick) Quality(quality uint) *Magick {
	if quality > 100 {
		quality = 100
	}

	c.quality = int(quality)
	return c
}

// Version returns convert version.
func (c *Magick) Version() (string, error) {
	err := c.BinWrapper.Run("-version")

	if err != nil {
		return "", err
	}

	return string(bytes.Split(c.StdOut(), []byte("\n"))[0]), nil
}

// Run starts convert with specified parameters.
func (c *Magick) Run() (image.Image, error) {
	defer c.BinWrapper.Reset()
	err := c.setInput()

	if err != nil {
		return nil, err
	}

	if c.quality > 0 {
		c.Arg("-quality", fmt.Sprintf("%d", c.quality))
	}

	if c.fuzz > 0 {
		c.Arg("-fuzz", fmt.Sprintf("%d%%", c.fuzz))
	}

	if c.trim {
		c.Arg("-trim").Arg("+repage")
	}

	c.setOutput()

	err = c.BinWrapper.Run()

	if err != nil {
		return nil, errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	if c.output == nil && c.outputFile == "" {
		return png.Decode(bytes.NewReader(c.BinWrapper.StdOut()))
	}

	return nil, nil
}

// Reset resets all parameters to default values
func (c *Magick) Reset() *Magick {
	c.quality = -1
	c.fuzz = 0
	c.outputFormat = "png"
	c.trim = false
	return c
}

func (c *Magick) setOutput() {
	if c.outputFile != "" {
		c.Arg(c.outputFile)
	} else if c.output != nil {
		c.Arg(c.outputFormat + ":-")
		c.SetStdOut(c.output)
	} else {
		c.Arg("png:-")
	}
}

func (c *Magick) setInput() error {
	if c.input != nil {
		c.Arg("-")
		c.StdIn(c.input)
	} else if c.inputImage != nil {
		r, err := createReaderFromImage(c.inputImage)

		if err != nil {
			return err
		}

		c.Arg("-")
		c.StdIn(r)
	} else if c.inputFile != "" {
		c.Arg(c.inputFile)
	} else {
		return errors.New("Undefined input")
	}

	return nil
}
