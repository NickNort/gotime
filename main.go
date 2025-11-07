package main

import (
	"image"
	"image/png"
	"os"

	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

func renderQR(bitmap [][]bool, moduleSize int) *image.RGBA {
	size := len(bitmap) * moduleSize

	dc := gg.NewContext(size, size)

	// Background
	dc.SetRGB(0.973, 0.949, 0.925) // #f8f2ec
	dc.Clear()

	// Foreground color
	dc.SetRGB(0.333, 0.125, 0.282) // #552048

	// Iterate through each row
	for y := 0; y < len(bitmap); y++ {
		x := 0
		for x < len(bitmap[y]) {
			if bitmap[y][x] {
				// Found start of a dark run
				startX := x

				// Count consecutive dark modules
				for x < len(bitmap[y]) && bitmap[y][x] {
					x++
				}
				endX := x

				// Draw one rounded rectangle for the entire run
				width := float64((endX - startX) * moduleSize)
				height := float64(moduleSize)
				xPos := float64(startX * moduleSize)
				yPos := float64(y * moduleSize)
				radius := float64(moduleSize) / 4.0 // Adjust for desired roundness

				dc.DrawRoundedRectangle(xPos, yPos, width, height, radius)
				dc.Fill()
			} else {
				x++
			}
		}
	}

	return dc.Image().(*image.RGBA)
}

func main() {
	// Generate QR code
	q, err := qrcode.New("https://nickn.dev", qrcode.Highest)
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
