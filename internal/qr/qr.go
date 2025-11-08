package qr

import (
	"bytes"
	"fmt"

	svg "github.com/ajstarks/svgo"
	"github.com/skip2/go-qrcode"
)

const (
	foregroundColor = "#552048"
	backgroundColor = "#f8f2ec"
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

// Options contains configuration for QR code generation
type Options struct {
	FinderCenter    string // "circle" | "square" | "diamond"
	FinderFrame     string // "square" | "rounded" | "circle" | "diamond"
	ModuleShape     string // "square" | "rounded" | "circle" | "diamond"
	ModuleSize      int    // pixels per module; 0 => default 10
	BackgroundColor string // hex color code (e.g., "#f8f2ec"); empty => default
	ForegroundColor string // hex color code (e.g., "#552048"); empty => default
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

func renderQR(bitmap [][]bool, moduleSize int, canvas *svg.SVG, corners CornerBounds, finderCenterStyle string, finderFrameStyle string, moduleShape string, bgColor string, fgColor string) {
	size := len(bitmap) * moduleSize

	// Start SVG
	canvas.Start(size, size)

	// Background
	canvas.Rect(0, 0, size, size, "fill:"+bgColor)

	// First pass: render all non-corner modules normally
	switch moduleShape {
	case "circle", "diamond":
		// Render each module individually for circle and diamond shapes
		for y := 0; y < len(bitmap); y++ {
			for x := 0; x < len(bitmap[y]); x++ {
				if bitmap[y][x] && !isInCorner(x, y, corners) {
					xPos := x * moduleSize
					yPos := y * moduleSize
					centerX := xPos + moduleSize/2
					centerY := yPos + moduleSize/2
					// Make circles and diamonds slightly larger (25% increase)
					radius := int(float64(moduleSize) / 2 * 1.25)

					if moduleShape == "circle" {
						canvas.Circle(centerX, centerY, radius, "fill:"+fgColor)
					} else { // diamond
						// Make diamonds larger than circles (45% increase)
						halfSize := int(float64(moduleSize) / 2 * 1.45)
						path := fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Z",
							centerX, centerY-halfSize, // Top
							centerX+halfSize, centerY, // Right
							centerX, centerY+halfSize, // Bottom
							centerX-halfSize, centerY) // Left
						canvas.Path(path, "fill:"+fgColor)
					}
				}
			}
		}
	default:
		// Render grouped modules for rounded and square shapes
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

					// Draw one shape for the entire run
					width := (endX - startX) * moduleSize
					height := moduleSize
					xPos := startX * moduleSize
					yPos := y * moduleSize

					if moduleShape == "rounded" {
						radius := moduleSize / 4
						canvas.Roundrect(xPos, yPos, width, height, radius, radius, "fill:"+fgColor)
					} else { // square
						canvas.Rect(xPos, yPos, width, height, "fill:"+fgColor)
					}
				} else {
					x++
				}
			}
		}
	}

	// Second pass: render corners specially
	cornerList := []CornerRect{corners.TopLeft, corners.TopRight, corners.BottomLeft}
	for _, corner := range cornerList {
		cornerX := corner.X * moduleSize
		cornerY := corner.Y * moduleSize

		// Render outer frame based on finderFrameStyle
		// The outer frame consists of the border (positions 0 and 6 on all edges)
		frameSize := 7 * moduleSize
		frameThickness := moduleSize
		switch finderFrameStyle {
		case "rounded":
			// Render as a rounded rectangle frame
			// Draw outer rounded rectangle, then subtract inner area
			radius := moduleSize / 2
			// Outer rounded rectangle
			canvas.Roundrect(cornerX, cornerY, frameSize, frameSize, radius, radius, "fill:"+fgColor)
			// Inner rounded rectangle (subtract by drawing background color)
			innerSize := 5 * moduleSize // Inner size is 5x5 (positions 1-5)
			innerX := cornerX + moduleSize
			innerY := cornerY + moduleSize
			canvas.Roundrect(innerX, innerY, innerSize, innerSize, radius, radius, "fill:"+bgColor)
		case "circle":
			// Render as a circular frame
			centerCX := cornerX + frameSize/2
			centerCY := cornerY + frameSize/2
			outerRadius := frameSize / 2
			innerRadius := (5 * moduleSize) / 2 // Inner radius for 5x5 area
			// Draw outer circle
			canvas.Circle(centerCX, centerCY, outerRadius, "fill:"+fgColor)
			// Subtract inner circle
			canvas.Circle(centerCX, centerCY, innerRadius, "fill:"+bgColor)
		case "diamond":
			// Render as a diamond-shaped frame
			centerCX := cornerX + frameSize/2
			centerCY := cornerY + frameSize/2
			// Outer diamond path (covers 7x7 area)
			outerPath := fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Z",
				centerCX, cornerY, // Top
				cornerX+frameSize, centerCY, // Right
				centerCX, cornerY+frameSize, // Bottom
				cornerX, centerCY) // Left
			canvas.Path(outerPath, "fill:"+fgColor)
			// Inner diamond path (subtract 5x5 area by drawing background color)
			innerPath := fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Z",
				centerCX, cornerY+moduleSize, // Top
				cornerX+6*moduleSize, centerCY, // Right
				centerCX, cornerY+6*moduleSize, // Bottom
				cornerX+moduleSize, centerCY) // Left
			canvas.Path(innerPath, "fill:"+bgColor)
		default:
			// Render as a square frame (default)
			// Top bar: full width (row 0, positions 0-6)
			canvas.Rect(cornerX, cornerY, frameSize, frameThickness, "fill:"+fgColor)
			// Bottom bar: full width (row 6, positions 0-6)
			canvas.Rect(cornerX, cornerY+6*moduleSize, frameSize, frameThickness, "fill:"+fgColor)
			// Left edge: full height (position 0, rows 0-6)
			canvas.Rect(cornerX, cornerY, frameThickness, frameSize, "fill:"+fgColor)
			// Right edge: full height (position 6, rows 0-6)
			canvas.Rect(cornerX+6*moduleSize, cornerY, frameThickness, frameSize, "fill:"+fgColor)
		}

		// Render inner center as either a circle, square, or diamond (positions 2-4, 3x3 block)
		centerX := cornerX + 2*moduleSize
		centerY := cornerY + 2*moduleSize
		switch finderCenterStyle {
		case "circle":
			// Center of the 3x3 block
			centerCX := centerX + (3*moduleSize)/2
			centerCY := centerY + (3*moduleSize)/2
			// Radius is half the width/height of the 3x3 block
			radius := (3 * moduleSize) / 2
			canvas.Circle(centerCX, centerCY, radius, "fill:"+fgColor)
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
			canvas.Path(path, "fill:"+fgColor)
		default:
			// Render as a square (default)
			canvas.Rect(centerX, centerY, 3*moduleSize, 3*moduleSize, "fill:"+fgColor)
		}
	}

	// End SVG
	canvas.End()
}

// GenerateSVG generates a QR code and returns the SVG as bytes
func GenerateSVG(content string, opts Options) ([]byte, error) {
	// Apply defaults
	if opts.FinderCenter == "" {
		opts.FinderCenter = "square"
	}
	if opts.FinderFrame == "" {
		opts.FinderFrame = "square"
	}
	if opts.ModuleShape == "" {
		opts.ModuleShape = "rounded"
	}
	if opts.ModuleSize == 0 {
		opts.ModuleSize = 10
	}
	if opts.BackgroundColor == "" {
		opts.BackgroundColor = backgroundColor
	}
	if opts.ForegroundColor == "" {
		opts.ForegroundColor = foregroundColor
	}

	// Validate parameters
	if opts.FinderCenter != "circle" && opts.FinderCenter != "square" && opts.FinderCenter != "diamond" {
		return nil, fmt.Errorf("finder-center must be either 'circle', 'square', or 'diamond' (got: %q)", opts.FinderCenter)
	}
	if opts.FinderFrame != "square" && opts.FinderFrame != "rounded" && opts.FinderFrame != "circle" && opts.FinderFrame != "diamond" {
		return nil, fmt.Errorf("finder-frame must be either 'square', 'rounded', 'circle', or 'diamond' (got: %q)", opts.FinderFrame)
	}
	if opts.ModuleShape != "square" && opts.ModuleShape != "rounded" && opts.ModuleShape != "circle" && opts.ModuleShape != "diamond" {
		return nil, fmt.Errorf("module-shape must be either 'square', 'rounded', 'circle', or 'diamond' (got: %q)", opts.ModuleShape)
	}

	// Generate QR code
	q, err := qrcode.New(content, qrcode.Highest)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Access the bitmap - this is the data matrix
	// It's a 2D boolean array where true = dark module, false = light module
	bitmap := q.Bitmap()

	// Find corner finder patterns
	corners := findCorners(bitmap)

	// Create buffer to capture SVG output
	var buf bytes.Buffer

	// Create SVG canvas
	canvas := svg.New(&buf)

	// Render QR code as SVG
	renderQR(bitmap, opts.ModuleSize, canvas, corners, opts.FinderCenter, opts.FinderFrame, opts.ModuleShape, opts.BackgroundColor, opts.ForegroundColor)

	return buf.Bytes(), nil
}
