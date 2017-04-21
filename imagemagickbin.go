package imagemagickbin

import (
	"bytes"
	"github.com/nickalie/go-binwrapper"
	"runtime"
	"image"
	"image/png"
	"io"
)

var skipDownload bool
var dest string = "vendor/imagemagick"

//Detects platforms without prebuilt binaries (alpine and arm).
//For this platforms imagemagick tools should be built manually.
func init() {
	if runtime.GOARCH == "arm" || runtime.GOOS == "linux" {
		SkipDownload()
	}
}

//Skips binary download.
func SkipDownload() {
	skipDownload = true
	dest = ""
}

//Sets directory to download imagemagick binaries or where to look for them if SkipDownload is used.
func Dest(value string) {
	dest = value
}

func createBinWrapper() *binwrapper.BinWrapper {
	base := "https://www.imagemagick.org/download/binaries/"

	b := binwrapper.NewBinWrapper().AutoExe()

	if !skipDownload {
		b.Src(
			binwrapper.NewSrc().
				Url(base + "ImageMagick-x86_64-apple-darwin16.4.0.tar.gz").
				Os("darwin").ExecPath("bin/convert")).
			Src(
			binwrapper.NewSrc().
				Url(base + "ImageMagick-7.0.5-4-portable-Q16-x64.zip").
				Os("win32").
				Arch("x64")).
			Src(
			binwrapper.NewSrc().
				Url(base + "ImageMagick-7.0.5-4-portable-Q16-x86.zip").
				Os("win32").
				Arch("x86"))
	}

	return b.ExecPath("convert").Strip(1).Dest(dest)
}

func createReaderFromImage(img image.Image) (io.Reader, error) {
	enc := &png.Encoder{
		CompressionLevel: png.NoCompression,
	}

	var buffer bytes.Buffer
	err := enc.Encode(&buffer, img)

	if err != nil {
		return nil, err
	}

	return &buffer, nil
}