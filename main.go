package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/skip2/go-qrcode"
)

func renderQR(bitmap [][]bool, moduleSize int) *image.RGBA {
	size := len(bitmap) * moduleSize
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Custom colors
	bgColor := color.RGBA{R: 248, G: 242, B: 236, A: 255} // #f8f2ec
	fgColor := color.RGBA{R: 85, G: 32, B: 72, A: 255}    // #552048

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw modules
	for y := 0; y < len(bitmap); y++ {
		for x := 0; x < len(bitmap[y]); x++ {
			if bitmap[y][x] {
				// Draw a square module
				rect := image.Rect(
					x*moduleSize, y*moduleSize,
					(x+1)*moduleSize, (y+1)*moduleSize,
				)
				draw.Draw(img, rect, &image.Uniform{fgColor}, image.Point{}, draw.Src)
			}
		}
	}

	return img
}

func main() {
	// Generate QR code
	q, err := qrcode.New("https://nickn.dev", qrcode.High)
	if err != nil {
		panic(err)
	}

	// Access the bitmap - this is the data matrix
	// It's a 2D boolean array where true = dark module, false = light module
	bitmap := q.Bitmap()

	// Render QR code as image
	moduleSize := 10 // Size of each module in pixels
	img := renderQR(bitmap, moduleSize)

	// Save it
	f, err := os.Create("qr.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
