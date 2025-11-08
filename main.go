package main

import (
	"flag"
	"fmt"
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

// getCornerRect returns the corner rectangle that contains (x, y), or nil if not in any corner
func getCornerRect(x, y int, corners CornerBounds) *CornerRect {
	// Check top-left corner
	if x >= corners.TopLeft.X && x < corners.TopLeft.X+corners.TopLeft.Width &&
		y >= corners.TopLeft.Y && y < corners.TopLeft.Y+corners.TopLeft.Height {
		return &corners.TopLeft
	}
	// Check top-right corner
	if x >= corners.TopRight.X && x < corners.TopRight.X+corners.TopRight.Width &&
		y >= corners.TopRight.Y && y < corners.TopRight.Y+corners.TopRight.Height {
		return &corners.TopRight
	}
	// Check bottom-left corner
	if x >= corners.BottomLeft.X && x < corners.BottomLeft.X+corners.BottomLeft.Width &&
		y >= corners.BottomLeft.Y && y < corners.BottomLeft.Y+corners.BottomLeft.Height {
		return &corners.BottomLeft
	}
	return nil
}

// isCornerOuterFrame checks if a module at (x, y) relative to corner rect is part of the outer frame
// Outer frame: positions 0 and 6 on all edges (the border)
func isCornerOuterFrame(x, y int, corner *CornerRect) bool {
	// Convert to relative coordinates within the 7x7 corner
	relX := x - corner.X
	relY := y - corner.Y

	// Outer frame: all positions where relX is 0 or 6, or relY is 0 or 6
	return relX == 0 || relX == 6 || relY == 0 || relY == 6
}

// isCornerInnerCenter checks if a module at (x, y) relative to corner rect is part of the inner center
// Inner center: positions 2-4 (the 3x3 center block)
func isCornerInnerCenter(x, y int, corner *CornerRect) bool {
	// Convert to relative coordinates within the 7x7 corner
	relX := x - corner.X
	relY := y - corner.Y

	// Inner center: positions 2-4 in both X and Y
	return relX >= 2 && relX <= 4 && relY >= 2 && relY <= 4
}

func renderQR(bitmap [][]bool, moduleSize int, canvas *svg.SVG, corners CornerBounds, cornerCenterStyle string) {
	size := len(bitmap) * moduleSize

	// Start SVG
	canvas.Start(size, size)

	// Background
	canvas.Rect(0, 0, size, size, "fill:#f8f2ec")

	// First pass: render all non-corner modules normally
	for y := 0; y < len(bitmap); y++ {
		x := 0
		for x < len(bitmap[y]) {
			if bitmap[y][x] && !isInCorner(x, y, corners) {
				// Found start of a dark run (not in corner)
				startX := x

				// Count consecutive dark modules (not in corners)
				for x < len(bitmap[y]) && bitmap[y][x] && !isInCorner(x, y, corners) {
					x++
				}
				endX := x

				// Draw one rounded rectangle for the entire run
				width := (endX - startX) * moduleSize
				height := moduleSize
				xPos := startX * moduleSize
				yPos := y * moduleSize
				radius := moduleSize / 4

				canvas.Roundrect(xPos, yPos, width, height, radius, radius, "fill:#552048")
			} else {
				x++
			}
		}
	}

	// Second pass: render corners specially
	cornerList := []CornerRect{corners.TopLeft, corners.TopRight, corners.BottomLeft}
	for _, corner := range cornerList {
		cornerX := corner.X * moduleSize
		cornerY := corner.Y * moduleSize

		// Render outer frame as a connected square frame (with gap for inner ring)
		// The outer frame consists of the border (positions 0 and 6 on all edges)
		// Top bar: full width (row 0, positions 0-6)
		canvas.Rect(cornerX, cornerY, 7*moduleSize, moduleSize, "fill:#552048")
		// Bottom bar: full width (row 6, positions 0-6)
		canvas.Rect(cornerX, cornerY+6*moduleSize, 7*moduleSize, moduleSize, "fill:#552048")
		// Left edge: full height (position 0, rows 0-6) - this creates the connected frame
		canvas.Rect(cornerX, cornerY, moduleSize, 7*moduleSize, "fill:#552048")
		// Right edge: full height (position 6, rows 0-6) - this creates the connected frame
		canvas.Rect(cornerX+6*moduleSize, cornerY, moduleSize, 7*moduleSize, "fill:#552048")

		// Render inner center as either a circle, square, or diamond (positions 2-4, 3x3 block)
		centerX := cornerX + 2*moduleSize
		centerY := cornerY + 2*moduleSize
		switch cornerCenterStyle {
		case "circle":
			// Center of the 3x3 block
			centerCX := centerX + (3*moduleSize)/2
			centerCY := centerY + (3*moduleSize)/2
			// Radius is half the width/height of the 3x3 block
			radius := (3 * moduleSize) / 2
			canvas.Circle(centerCX, centerCY, radius, "fill:#552048")
		case "diamond":
			// Render as a diamond (rotated square)
			// Center of the 3x3 block
			centerCX := centerX + (3*moduleSize)/2
			centerCY := centerY + (3*moduleSize)/2
			// Half-size for the diamond (diagonal distance from center to corner)
			halfSize := (3 * moduleSize) / 2
			// Diamond path: M (move to top), L (line to right), L (line to bottom), L (line to left), Z (close)
			path := fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Z",
				centerCX, centerCY-halfSize, // Move to top
				centerCX+halfSize, centerCY, // Line to right
				centerCX, centerCY+halfSize, // Line to bottom
				centerCX-halfSize, centerCY) // Line to left, close
			canvas.Path(path, "fill:#552048")
		default:
			// Render as a square (default)
			canvas.Rect(centerX, centerY, 3*moduleSize, 3*moduleSize, "fill:#552048")
		}
	}

	// End SVG
	canvas.End()
}

func main() {
	// Parse command line flags
	cornerCenter := flag.String("corner-center", "square", "Corner center style: 'circle', 'square', or 'diamond'")
	flag.Parse()

	// Get QR code content from positional argument
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <qr-content> [-corner-center=<style>]\n", os.Args[0])
		os.Exit(1)
	}
	qrContent := args[0]

	// Validate corner center style
	if *cornerCenter != "circle" && *cornerCenter != "square" && *cornerCenter != "diamond" {
		panic("corner-center must be either 'circle', 'square', or 'diamond'")
	}

	// Generate QR code
	q, err := qrcode.New(qrContent, qrcode.Highest)
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
	renderQR(bitmap, moduleSize, canvas, corners, *cornerCenter)
}
