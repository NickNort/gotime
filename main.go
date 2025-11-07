package main

import (
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/skip2/go-qrcode"
)

func renderQR(bitmap [][]bool, moduleSize int, canvas *svg.SVG) {
	size := len(bitmap) * moduleSize

	// Start SVG
	canvas.Start(size, size)

	// Background
	canvas.Rect(0, 0, size, size, "fill:#f8f2ec")

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
				width := (endX - startX) * moduleSize
				height := moduleSize
				xPos := startX * moduleSize
				yPos := y * moduleSize
				radius := moduleSize / 4 // Adjust for desired roundness

				canvas.Roundrect(xPos, yPos, width, height, radius, radius, "fill:#552048")
			} else {
				x++
			}
		}
	}

	// End SVG
	canvas.End()
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

	// Create output file
	f, err := os.Create("qr.svg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create SVG canvas
	canvas := svg.New(f)

	// Render QR code as SVG
	moduleSize := 10 // Size of each module in pixels
	renderQR(bitmap, moduleSize, canvas)
}