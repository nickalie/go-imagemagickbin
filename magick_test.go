package imagemagickbin

import (
	"net/http"
	"os"
	"io"
	"github.com/stretchr/testify/assert"
	"testing"
	"image/jpeg"
	"golang.org/x/image/webp"
	"fmt"
	"image/png"
)

func init() {
	downloadFile("https://upload.wikimedia.org/wikipedia/commons/e/e3/Avola-Syracuse-Sicilia-Italy_-_Creative_Commons_by_gnuckx_%283858115914%29.jpg", "source.jpg")
	downloadFile("https://upload.wikimedia.org/wikipedia/commons/d/d1/Snail_in_Forest_on_Lichtscheid_2.webp", "source.webp")
}

func downloadFile(url, target string) {
	_, err := os.Stat(target)

	if err != nil {
		resp, err := http.Get(url)

		if err != nil {
			fmt.Printf("Error while downloading test image: %v\n", err)
			panic(err)
		}

		defer resp.Body.Close()

		f, err := os.Create(target)

		if err != nil {
			panic(err)
		}

		defer f.Close()

		_, err = io.Copy(f, resp.Body)

		if err != nil {
			panic(err)
		}
	}
}

func TestEncodeImage(t *testing.T) {
	c := NewMagick()
	f, err := os.Open("source.jpg")
	assert.Nil(t, err)
	img, err := jpeg.Decode(f)
	assert.Nil(t, err)
	c.InputImage(img)
	c.OutputFile("target.webp")
	img, err = c.Run()
	assert.Nil(t, err)
	assert.Nil(t, img)
	validateWebp(t)
}

func TestEncodeReader(t *testing.T) {
	c := NewMagick()
	f, err := os.Open("source.jpg")
	assert.Nil(t, err)
	c.Input(f)
	c.OutputFile("target.webp")
	img, err := c.Run()
	assert.Nil(t, err)
	assert.Nil(t, img)
	validateWebp(t)
}

func TestEncodeFile(t *testing.T) {
	c := NewMagick()
	c.InputFile("source.jpg")
	c.OutputFile("target.webp")
	img, err := c.Run()
	assert.Nil(t, err)
	assert.Nil(t, img)
	validateWebp(t)
}

func TestEncodeWriter(t *testing.T) {
	f, err := os.Create("target.webp")
	assert.Nil(t, err)
	defer f.Close()

	c := NewMagick()
	c.InputFile("source.jpg")
	c.OutputFormat("webp")
	c.Output(f)
	img, err := c.Run()
	assert.Nil(t, err)
	assert.Nil(t, img)
	f.Close()
	validateWebp(t)
}

func TestVersion(t *testing.T) {
	c := NewMagick()
	_, err := c.Version()
	assert.Nil(t, err)
}

func TestGetTrimInfo(t *testing.T) {
	c := NewMagick()
	c.InputFile("source.jpg")
	rect, err := c.GetTrimInfo(25, 0)
	assert.Nil(t, err)
	assert.NotNil(t, rect)
}

func validateWebp(t *testing.T) {
	defer os.Remove("target.webp")
	fSource, err := os.Open("source.jpg")
	assert.Nil(t, err)
	imgSource, err := jpeg.Decode(fSource)
	assert.Nil(t, err)
	fTarget, err := os.Open("target.webp")
	assert.Nil(t, err)
	defer fTarget.Close()
	imgTarget, err := webp.Decode(fTarget)
	assert.Nil(t, err)
	assert.Equal(t, imgSource.Bounds(), imgTarget.Bounds())
}

func TestDecodeReader(t *testing.T) {
	c := NewMagick()
	f, err := os.Open("source.webp")
	assert.Nil(t, err)
	defer f.Close()
	c.Input(f)
	c.OutputFile("target.png")
	img, err := c.Run()
	assert.Nil(t, err)
	assert.Nil(t, img)
	validatePng(t)
}

func TestDecodeFile(t *testing.T) {
	c := NewMagick()
	c.InputFile("source.webp")
	c.OutputFile("target.png")
	img, err := c.Run()
	assert.Nil(t, err)
	assert.Nil(t, img)
	validatePng(t)
}

func TestDecodeImage(t *testing.T) {
	c := NewMagick()
	f, err := os.Open("source.webp")
	assert.Nil(t, err)
	defer f.Close()
	imgSource, err := webp.Decode(f)
	assert.Nil(t, err)
	f.Seek(0, 0)
	c.Input(f)
	imgTarget, err := c.Run()
	assert.Nil(t, err)
	assert.NotNil(t, imgTarget)
	assert.Equal(t, imgSource.Bounds(), imgTarget.Bounds())
}

func TestDecodeWriter(t *testing.T) {
	f, err := os.Create("target.png")
	assert.Nil(t, err)
	defer f.Close()
	c := NewMagick()
	c.InputFile("source.webp")
	c.Output(f)
	img, err := c.Run()
	assert.Nil(t, err)
	assert.Nil(t, img)
	f.Close()
	validatePng(t)
}

func validatePng(t *testing.T) {
	defer os.Remove("target.png")
	fSource, err := os.Open("source.webp")
	assert.Nil(t, err)
	imgSource, err := webp.Decode(fSource)
	assert.Nil(t, err)
	fTarget, err := os.Open("target.png")
	assert.Nil(t, err)
	defer fTarget.Close()
	imgTarget, err := png.Decode(fTarget)
	assert.Nil(t, err)
	assert.Equal(t, imgSource.Bounds(), imgTarget.Bounds())
}

