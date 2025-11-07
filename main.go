package main

import (
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/skip2/go-qrcode"
)

// CornerRect represents a rectangular region of a corner finder pattern
type CornerRect struct {
	X      int // Top-left x coordinate in modules
	Y      int // Top-left y coordinate in modules
	Width  int // Width in modules (always 7 for finder patterns)
	Height int // Height in modules (always 7 for finder patterns)
}

// CornerBounds contains the three corner finder pattern rectangles
type CornerBounds struct {
	TopLeft    CornerRect
	TopRight   CornerRect
	BottomLeft CornerRect
}

// isFinderPattern checks if a 7x7 region starting at (x, y) matches a QR code finder pattern
// Finder pattern structure:
// - Outer border (all edges): dark (true)
// - Inner ring (positions 1-5 excluding center): light (false)
// - Center (positions 2-4): dark (true)
func isFinderPattern(bitmap [][]bool, x, y int) bool {
	size := len(bitmap)
	if x+7 > size || y+7 > size {
		return false
	}

	// Check outer border - all must be dark
	for i := 0; i < 7; i++ {
		// Top and bottom edges
		if !bitmap[y][x+i] || !bitmap[y+6][x+i] {
			return false
		}
		// Left and right edges
		if !bitmap[y+i][x] || !bitmap[y+i][x+6] {
			return false
		}
	}

	// Check inner ring - must be light
	// Row 1 and row 5: all columns 1-5 must be light
	for i := 1; i < 6; i++ {
		if bitmap[y+1][x+i] || bitmap[y+5][x+i] {
			return false
		}
	}
	// Column 1 and column 5: rows 2-4 must be light (center rows)
	for i := 2; i < 5; i++ {
		if bitmap[y+i][x+1] || bitmap[y+i][x+5] {
			return false
		}
	}

	// Check center 3x3 - must be dark
	for i := 2; i < 5; i++ {
		for j := 2; j < 5; j++ {
			if !bitmap[y+i][x+j] {
				return false
			}
		}
	}

	return true
}

// findCorners identifies the three corner finder patterns in a QR code bitmap
func findCorners(bitmap [][]bool) CornerBounds {
	size := len(bitmap)
	var bounds CornerBounds

	// Search region for top corners (vertical limit)
	topSearchLimit := size / 3 // Search top third for top corners
	if topSearchLimit < 7 {
		topSearchLimit = 7
	}
	if topSearchLimit > size-7 {
		topSearchLimit = size - 7
	}

	// Search region for left/right corners (horizontal limit)
	sideSearchLimit := size / 3 // Search outer third for side corners
	if sideSearchLimit < 7 {
		sideSearchLimit = 7
	}
	if sideSearchLimit > size-7 {
		sideSearchLimit = size - 7
	}

	// Find top-left corner (search from top-left)
	for y := 0; y < topSearchLimit; y++ {
		for x := 0; x < sideSearchLimit; x++ {
			if isFinderPattern(bitmap, x, y) {
				bounds.TopLeft = CornerRect{X: x, Y: y, Width: 7, Height: 7}
				goto foundTopLeft
			}
		}
	}
foundTopLeft:

	// Find top-right corner (search from top-right, right to left)
	maxX := size - 7
	minX := size - sideSearchLimit
	if minX < 0 {
		minX = 0
	}
	for y := 0; y < topSearchLimit; y++ {
		// Search from right edge going left
		for x := maxX; x >= minX; x-- {
			if isFinderPattern(bitmap, x, y) {
				bounds.TopRight = CornerRect{X: x, Y: y, Width: 7, Height: 7}
				goto foundTopRight
			}
		}
	}
foundTopRight:

	// Find bottom-left corner (search from bottom-left, bottom to top)
	maxY := size - 7
	minY := size - topSearchLimit
	if minY < 0 {
		minY = 0
	}
	// Search from bottom going up
	for y := maxY; y >= minY; y-- {
		for x := 0; x < sideSearchLimit; x++ {
			if isFinderPattern(bitmap, x, y) {
				bounds.BottomLeft = CornerRect{X: x, Y: y, Width: 7, Height: 7}
				goto foundBottomLeft
			}
		}
	}
foundBottomLeft:

	return bounds
}

// isInCorner checks if a module at (x, y) is within any of the corner rectangles
func isInCorner(x, y int, corners CornerBounds) bool {
	// Check top-left corner
	if x >= corners.TopLeft.X && x < corners.TopLeft.X+corners.TopLeft.Width &&
		y >= corners.TopLeft.Y && y < corners.TopLeft.Y+corners.TopLeft.Height {
		return true
	}
	// Check top-right corner
	if x >= corners.TopRight.X && x < corners.TopRight.X+corners.TopRight.Width &&
		y >= corners.TopRight.Y && y < corners.TopRight.Y+corners.TopRight.Height {
		return true
	}
	// Check bottom-left corner
	if x >= corners.BottomLeft.X && x < corners.BottomLeft.X+corners.BottomLeft.Width &&
		y >= corners.BottomLeft.Y && y < corners.BottomLeft.Y+corners.BottomLeft.Height {
		return true
	}
	return false
}

func renderQR(bitmap [][]bool, moduleSize int, canvas *svg.SVG, corners CornerBounds) {
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
				inCorner := isInCorner(x, y, corners)

				// Count consecutive dark modules, checking if any are in corners
				for x < len(bitmap[y]) && bitmap[y][x] {
					if isInCorner(x, y, corners) {
						inCorner = true
					}
					x++
				}
				endX := x

				// If any part of the run is in a corner, draw individual modules
				// Otherwise, draw as a single rounded rectangle
				if inCorner {
					// Draw each module individually as regular rectangles
					for i := startX; i < endX; i++ {
						xPos := i * moduleSize
						yPos := y * moduleSize
						if isInCorner(i, y, corners) {
							// Regular rectangle for corner modules
							canvas.Rect(xPos, yPos, moduleSize, moduleSize, "fill:#552048")
						} else {
							// Rounded rectangle for non-corner modules in the same run
							radius := moduleSize / 4
							canvas.Roundrect(xPos, yPos, moduleSize, moduleSize, radius, radius, "fill:#552048")
						}
					}
				} else {
					// Draw one rounded rectangle for the entire run
					width := (endX - startX) * moduleSize
					height := moduleSize
					xPos := startX * moduleSize
					yPos := y * moduleSize
					radius := moduleSize / 4 // Adjust for desired roundness

					canvas.Roundrect(xPos, yPos, width, height, radius, radius, "fill:#552048")
				}
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
	q, err := qrcode.New("https://x.com/ItsNickNorton/status/1986520090666569982", qrcode.Highest)
	if err != nil {
		panic(err)
	}

	// Access the bitmap - this is the data matrix
	// It's a 2D boolean array where true = dark module, false = light module
	bitmap := q.Bitmap()

	// Find corner finder patterns
	corners := findCorners(bitmap)

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
	renderQR(bitmap, moduleSize, canvas, corners)
}
