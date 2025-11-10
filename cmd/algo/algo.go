package main

import (
	"fmt"
	"strconv"
	"os"
)

type Coord struct {
	X, Y int
}

func makeVisitedMatrix(bitmap [][]bool) [][]bool {
	var sizeX, sizeY int

	sizeY = len(bitmap)
	if sizeY > 0 {
		sizeX = len(bitmap[0])
	}

	visitedMatrix := make([][]bool, sizeY)
	for i := range visitedMatrix {
		visitedMatrix[i] = make([]bool, sizeX)
	}

	return visitedMatrix
}

// arguments for printBoolmatrix() so that `label` is optional
type pbmOpts struct {
	label  string
	m      [][]bool
	coords [][]Coord
}

func printBoolMatrix(opts pbmOpts) {
	if opts.label != "" {
		println(opts.label)
	}

	m := opts.m

	coordMap := make(map[Coord]bool)
	for _, group := range opts.coords {
		for _, coord := range group {
			coordMap[coord] = true
		}
	}

	for j, row := range m {
		for i, val := range row {
			if coordMap[Coord{X: i, Y: j}] {
				if val {
					print("[1] ")
				} else {
					print("[0] ")
				}
			} else {
				if val {
					print(" 1  ")
				} else {
					print(" 0  ")
				}
			}
		}
		println()
	}
	println()
}

// arguments for printIntmatrix() so that `label` is optional
type pimOpts struct {
	label string
	m     [][]int
}

func printIntMatrix(opts pimOpts) {
	if opts.label != "" {
		println(opts.label)
	}

	m := opts.m
	for _, row := range m {
		for i := 0; i < len(row); i += 2 {
			if i+1 < len(row) {
				print("(", row[i], ",", row[i+1], ") ")
			}
		}
		println()
	}
}

type pcmOpts struct {
	label  string
	coords [][]Coord
}

func printCoordMatrix(opts pcmOpts) {
	if opts.label != "" {
		println(opts.label)
	}

	coords := opts.coords
	for i, group := range coords {
		fmt.Printf("Group %d: ", i)
		for _, coord := range group {
			fmt.Printf("(%d,%d) ", coord.Y, coord.X)
		}
		println()
	}
	println()
}

func algo(bitmap [][]bool) ([][]Coord, error) {
	sizeY := len(bitmap)
	if sizeY <= 0 {
		return nil, fmt.Errorf("bitmap cannot be empty")
	}

	sizeX := len(bitmap[0])
	if sizeX <= 0 {
		return nil, fmt.Errorf("bitmap rows cannot be empty")
	}

	visitedMatrix := makeVisitedMatrix(bitmap)

	printBoolMatrix(pbmOpts{label: "bitmap:", m: bitmap})
	// printBoolMatrix(pbmOpts{label: "visitedMatrix:", m: visitedMatrix})
	// print("sizeY: ", sizeY, "\n")

	res := new([][]Coord)
	var currentGroup *[]Coord

	for y, r := range bitmap {
		for x, c := range r {
			if c == true && visitedMatrix[y][x] == false {
				currentGroup = new([]Coord)
				DFS(bitmap, visitedMatrix, x, y, currentGroup)
				*res = append(*res, *currentGroup)
			}
		}
	}

	// printIntMatrix(pimOpts{label: "res:", m: *res})

	return *res, nil
}

func DFS(bitmap [][]bool, visitedMatrix [][]bool, x int, y int, currentGroup *[]Coord) bool {
	if y >= len(bitmap) || y < 0 || x >= len(bitmap[y]) || x < 0 || visitedMatrix[y][x] == true || bitmap[y][x] == false {
		return false
	}
	visitedMatrix[y][x] = true
	*currentGroup = append(*currentGroup, Coord{X: x, Y: y})

	DFS(bitmap, visitedMatrix, x-1, y, currentGroup)
	DFS(bitmap, visitedMatrix, x+1, y, currentGroup)
	DFS(bitmap, visitedMatrix, x, y-1, currentGroup)
	DFS(bitmap, visitedMatrix, x, y+1, currentGroup)

	return true
}

func runCCATests() {
	// Initialize a 2D bool array for testing
	bitmap1 := [][]bool{
		{true, false, true},
		{false, true, false},
		{true, true, false},
		{true, true, true},
	}

	res1, err := algo(bitmap1)
	if err != nil {
		fmt.Printf("Error processing bitmap1: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap1 with groups:", m: bitmap1, coords: res1})
	printCoordMatrix(pcmOpts{label: "res1:", coords: res1})

	bitmap2 := [][]bool{
		{true, false, true, true},
		{false, true, false, true},
		{true, true, false, true},
	}

	res2, err := algo(bitmap2)
	if err != nil {
		fmt.Printf("Error processing bitmap2: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap2 with groups:", m: bitmap2, coords: res2})
	printCoordMatrix(pcmOpts{label: "res2:", coords: res2})

	bitmap3 := [][]bool{
		{true, true, true, true},
		{true, true, true, true},
		{true, true, true, true},
	}

	res3, err := algo(bitmap3)
	if err != nil {
		fmt.Printf("Error processing bitmap3: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap3 with groups:", m: bitmap3, coords: res3})
	printCoordMatrix(pcmOpts{label: "res3:", coords: res3})

	bitmap4 := [][]bool{
		{true, false, false, false},
		{false, true, false, false},
		{false, false, true, false},
		{false, false, false, true},
	}

	res4, err := algo(bitmap4)
	if err != nil {
		fmt.Printf("Error processing bitmap4: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap4 with groups:", m: bitmap4, coords: res4})
	printCoordMatrix(pcmOpts{label: "res4:", coords: res4})

	bitmap5 := [][]bool{
		{false, false, false, false},
		{false, true, true, false},
		{false, true, false, false},
		{false, true, true, false},
		{false, false, false, false},
	}

	res5, err := algo(bitmap5)
	if err != nil {
		fmt.Printf("Error processing bitmap5: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap5 with groups:", m: bitmap5, coords: res5})
	printCoordMatrix(pcmOpts{label: "res5:", coords: res5})

	bitmap6 := [][]bool{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	}

	res6, err := algo(bitmap6)
	if err != nil {
		fmt.Printf("Error processing bitmap6: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap6 with groups:", m: bitmap6, coords: res6})
	printCoordMatrix(pcmOpts{label: "res6:", coords: res6})

	bitmap7 := [][]bool{
		{true},
	}

	res7, err := algo(bitmap7)
	if err != nil {
		fmt.Printf("Error processing bitmap7: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap7 with groups:", m: bitmap7, coords: res7})
	printCoordMatrix(pcmOpts{label: "res7:", coords: res7})

	bitmap8 := [][]bool{
		{true, false, true, false},
		{false, true, false, true},
		{true, false, true, false},
		{false, true, false, true},
	}

	res8, err := algo(bitmap8)
	if err != nil {
		fmt.Printf("Error processing bitmap8: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap8 with groups:", m: bitmap8, coords: res8})
	printCoordMatrix(pcmOpts{label: "res8:", coords: res8})

	bitmap9 := [][]bool{
		{false, false, true, false, false},
		{false, false, true, false, false},
		{true, true, true, true, true},
		{false, false, true, false, false},
		{false, false, true, false, false},
	}

	res9, err := algo(bitmap9)
	if err != nil {
		fmt.Printf("Error processing bitmap9: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap9 with groups:", m: bitmap9, coords: res9})
	printCoordMatrix(pcmOpts{label: "res9:", coords: res9})

	bitmap10 := [][]bool{
		{true, true, true, true, true},
		{true, false, false, false, true},
		{true, false, false, false, true},
		{true, false, false, false, true},
		{true, true, true, true, true},
	}

	res10, err := algo(bitmap10)
	if err != nil {
		fmt.Printf("Error processing bitmap10: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap10 with groups:", m: bitmap10, coords: res10})
	printCoordMatrix(pcmOpts{label: "res10:", coords: res10})

	bitmap11 := [][]bool{}

	res11, err := algo(bitmap11)
	if err != nil {
		fmt.Printf("Error processing bitmap11: %v\n", err)
	} else {
		printBoolMatrix(pbmOpts{label: "bitmap11 with groups:", m: bitmap11, coords: res11})
		printCoordMatrix(pcmOpts{label: "res11:", coords: res11})
	}

	bitmap12 := [][]bool{
		{true, false, true},
		{},
		{false, true, false},
	}

	res12, err := algo(bitmap12)
	if err != nil {
		fmt.Printf("Error processing bitmap12: %v\n", err)
	} else {
		printBoolMatrix(pbmOpts{label: "bitmap12 with groups:", m: bitmap12, coords: res12})
		printCoordMatrix(pcmOpts{label: "res12:", coords: res12})
	}

	bitmap13 := [][]bool{
		{true, false, true, false},
		{false, true},
		{true, false, true, false},
	}

	res13, err := algo(bitmap13)
	if err != nil {
		fmt.Printf("Error processing bitmap13: %v\n", err)
	} else {
		printBoolMatrix(pbmOpts{label: "bitmap13 with groups:", m: bitmap13, coords: res13})
		printCoordMatrix(pcmOpts{label: "res13:", coords: res13})
	}
}

func testSquarePath() {
	fmt.Println("Testing SquarePath...")

	// Test with x=10, y=20, size=5
	path := SquarePath(10, 20, 5)

	// Verify we have 4 points
	if len(path.points) != 4 {
		fmt.Printf("FAIL: Expected 4 points, got %d\n", len(path.points))
		return
	}

	// Verify points are correct: top-left (10,20), top-right (15,20), bottom-right (15,15), bottom-left (10,15)
	expectedPoints := []Coord{
		{X: 10, Y: 20}, // top left
		{X: 15, Y: 20}, // top right
		{X: 15, Y: 15}, // bottom right
		{X: 10, Y: 15}, // bottom left
	}

	for i, expected := range expectedPoints {
		if path.points[i].X != expected.X || path.points[i].Y != expected.Y {
			fmt.Printf("FAIL: Point %d expected (%d,%d), got (%d,%d)\n",
				i, expected.X, expected.Y, path.points[i].X, path.points[i].Y)
			return
		}
	}

	// Verify path string
	expectedPath := "M 10,20 L 15,20 L 15,15 L 10,15 Z"
	if path.path != expectedPath {
		fmt.Printf("FAIL: Expected path '%s', got '%s'\n", expectedPath, path.path)
		return
	}

	fmt.Println("PASS: SquarePath test passed!")
	fmt.Printf("  Points: %v\n", path.points)
	fmt.Printf("  Path: %s\n", path.path)
}

func testSquarePathSVG1() {
	fmt.Println("Testing SquarePath with SVG generation...")

	// Create multiple squares at different positions
	squares := []struct {
		x, y, size int
		fill       string
	}{
		{50, 200, 30, "#ff6b6b"},
		{100, 200, 30, "#4ecdc4"},
		{150, 200, 30, "#45b7d1"},
		{200, 200, 30, "#f9ca24"},
		{250, 200, 30, "#6c5ce7"},
		{100, 150, 50, "#a29bfe"},
		{175, 100, 40, "#fd79a8"},
		{50, 100, 25, "#00b894"},
		{300, 150, 35, "#e17055"},
	}

	// Start building SVG content
	svgContent := `<?xml version="1.0"?>
<svg width="400" height="300" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="400" height="300" style="fill:#f8f9fa" />
`

	// Generate paths for each square
	for i, sq := range squares {
		path := SquarePath(sq.x, sq.y, sq.size)
		svgContent += fmt.Sprintf(`<path d="%s" style="fill:%s;stroke:#2d3436;stroke-width:2" />`+"\n", path.path, sq.fill)
		fmt.Printf("  Square %d: %s\n", i+1, path.path)
	}

	svgContent += `</svg>`

	// Write to file
	filename := "test_squarepath.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  Generated %d squares\n", len(squares))
}

func testSquarePathSingleSVG() {
	fmt.Println("Testing SquarePath with single centered square...")

	// SVG dimensions
	svgWidth := 400
	svgHeight := 300
	squareSize := 150

	// Calculate center position
	centerX := svgWidth / 2
	centerY := svgHeight / 2

	// Calculate top-left position (SquarePath draws upward, so y is the top)
	topLeftX := centerX - squareSize/2
	topLeftY := centerY + squareSize/2

	// Generate the square path
	path := SquarePath(topLeftX, topLeftY, squareSize)

	// Build SVG content
	svgContent := fmt.Sprintf(`<?xml version="1.0"?>
<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="%d" height="%d" style="fill:#ffffff" />
<path d="%s" style="fill:#4a90e2;stroke:#2d3436;stroke-width:3" />
</svg>`, svgWidth, svgHeight, svgWidth, svgHeight, path.path)

	// Write to file
	filename := "test_squarepath_single.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  Square path: %s\n", path.path)
	fmt.Printf("  Position: (%d, %d), Size: %d\n", topLeftX, topLeftY, squareSize)
}

func testDiamondPath() {
	fmt.Println("Testing DiamondPath...")

	// Test with x=10, y=20, size=10
	path := DiamondPath(10, 20, 10)

	// Verify we have 4 points
	if len(path.points) != 4 {
		fmt.Printf("FAIL: Expected 4 points, got %d\n", len(path.points))
		return
	}

	// Verify points are correct: top (15,20), right (20,15), bottom (15,10), left (10,15)
	expectedPoints := []Coord{
		{X: 15, Y: 20}, // top
		{X: 20, Y: 15}, // right
		{X: 15, Y: 10}, // bottom
		{X: 10, Y: 15}, // left
	}

	for i, expected := range expectedPoints {
		if path.points[i].X != expected.X || path.points[i].Y != expected.Y {
			fmt.Printf("FAIL: Point %d expected (%d,%d), got (%d,%d)\n",
				i, expected.X, expected.Y, path.points[i].X, path.points[i].Y)
			return
		}
	}

	// Verify path string
	expectedPath := "M 15,20 L 20,15 L 15,10 L 10,15 Z"
	if path.path != expectedPath {
		fmt.Printf("FAIL: Expected path '%s', got '%s'\n", expectedPath, path.path)
		return
	}

	fmt.Println("PASS: DiamondPath test passed!")
	fmt.Printf("  Points: %v\n", path.points)
	fmt.Printf("  Path: %s\n", path.path)
}

func testDiamondPathSVG1() {
	fmt.Println("Testing DiamondPath with SVG generation...")

	// Create multiple diamonds at different positions
	diamonds := []struct {
		x, y, size int
		fill       string
	}{
		{50, 200, 30, "#ff6b6b"},
		{100, 200, 30, "#4ecdc4"},
		{150, 200, 30, "#45b7d1"},
		{200, 200, 30, "#f9ca24"},
		{250, 200, 30, "#6c5ce7"},
		{100, 150, 50, "#a29bfe"},
		{175, 100, 40, "#fd79a8"},
		{50, 100, 25, "#00b894"},
		{300, 150, 35, "#e17055"},
	}

	// Start building SVG content
	svgContent := `<?xml version="1.0"?>
<svg width="400" height="300" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="400" height="300" style="fill:#f8f9fa" />
`

	// Generate paths for each diamond
	for i, d := range diamonds {
		path := DiamondPath(d.x, d.y, d.size)
		svgContent += fmt.Sprintf(`<path d="%s" style="fill:%s;stroke:#2d3436;stroke-width:2" />`+"\n", path.path, d.fill)
		fmt.Printf("  Diamond %d: %s\n", i+1, path.path)
	}

	svgContent += `</svg>`

	// Write to file
	filename := "test_diamondpath.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  Generated %d diamonds\n", len(diamonds))
}

func testDiamondPathSingleSVG() {
	fmt.Println("Testing DiamondPath with single centered diamond...")

	// SVG dimensions
	svgWidth := 400
	svgHeight := 300
	diamondSize := 150

	// Calculate center position
	centerX := svgWidth / 2
	centerY := svgHeight / 2

	// Calculate top position (DiamondPath draws upward, so y is the top)
	topX := centerX - diamondSize/2
	topY := centerY + diamondSize/2

	// Generate the diamond path
	path := DiamondPath(topX, topY, diamondSize)

	// Build SVG content
	svgContent := fmt.Sprintf(`<?xml version="1.0"?>
<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="%d" height="%d" style="fill:#ffffff" />
<path d="%s" style="fill:#4a90e2;stroke:#2d3436;stroke-width:3" />
</svg>`, svgWidth, svgHeight, svgWidth, svgHeight, path.path)

	// Write to file
	filename := "test_diamondpath_single.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  Diamond path: %s\n", path.path)
	fmt.Printf("  Position: (%d, %d), Size: %d\n", topX, topY, diamondSize)
}

func testTrianglePath() {
	fmt.Println("Testing TrianglePath...")

	// Test with x=10, y=20, size=10
	path := TrianglePath(10, 20, 10)

	// Verify we have 3 points
	if len(path.points) != 3 {
		fmt.Printf("FAIL: Expected 3 points, got %d\n", len(path.points))
		return
	}

	// Verify points are correct: top left (10,20), top right (20,20), bottom (15,10)
	expectedPoints := []Coord{
		{X: 10, Y: 20}, // top left (base)
		{X: 20, Y: 20}, // top right (base)
		{X: 15, Y: 10}, // bottom (apex)
	}

	for i, expected := range expectedPoints {
		if path.points[i].X != expected.X || path.points[i].Y != expected.Y {
			fmt.Printf("FAIL: Point %d expected (%d,%d), got (%d,%d)\n",
				i, expected.X, expected.Y, path.points[i].X, path.points[i].Y)
			return
		}
	}

	// Verify path string
	expectedPath := "M 10,20 L 20,20 L 15,10 Z"
	if path.path != expectedPath {
		fmt.Printf("FAIL: Expected path '%s', got '%s'\n", expectedPath, path.path)
		return
	}

	fmt.Println("PASS: TrianglePath test passed!")
	fmt.Printf("  Points: %v\n", path.points)
	fmt.Printf("  Path: %s\n", path.path)
}

func testTrianglePathSVG1() {
	fmt.Println("Testing TrianglePath with SVG generation...")

	// Create multiple triangles at different positions
	triangles := []struct {
		x, y, size int
		fill       string
	}{
		{50, 200, 30, "#ff6b6b"},
		{100, 200, 30, "#4ecdc4"},
		{150, 200, 30, "#45b7d1"},
		{200, 200, 30, "#f9ca24"},
		{250, 200, 30, "#6c5ce7"},
		{100, 150, 50, "#a29bfe"},
		{175, 100, 40, "#fd79a8"},
		{50, 100, 25, "#00b894"},
		{300, 150, 35, "#e17055"},
	}

	// Start building SVG content
	svgContent := `<?xml version="1.0"?>
<svg width="400" height="300" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="400" height="300" style="fill:#f8f9fa" />
`

	// Generate paths for each triangle
	for i, t := range triangles {
		path := TrianglePath(t.x, t.y, t.size)
		svgContent += fmt.Sprintf(`<path d="%s" style="fill:%s;stroke:#2d3436;stroke-width:2" />`+"\n", path.path, t.fill)
		fmt.Printf("  Triangle %d: %s\n", i+1, path.path)
	}

	svgContent += `</svg>`

	// Write to file
	filename := "test_trianglepath.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  Generated %d triangles\n", len(triangles))
}

func testTrianglePathSingleSVG() {
	fmt.Println("Testing TrianglePath with single centered triangle...")

	// SVG dimensions
	svgWidth := 400
	svgHeight := 300
	triangleSize := 150

	// Calculate center position
	centerX := svgWidth / 2
	centerY := svgHeight / 2

	// Calculate top position (TrianglePath draws downward, so y is the top of the base)
	topX := centerX - triangleSize/2
	topY := centerY + triangleSize/2

	// Generate the triangle path
	path := TrianglePath(topX, topY, triangleSize)

	// Build SVG content
	svgContent := fmt.Sprintf(`<?xml version="1.0"?>
<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="%d" height="%d" style="fill:#ffffff" />
<path d="%s" style="fill:#4a90e2;stroke:#2d3436;stroke-width:3" />
</svg>`, svgWidth, svgHeight, svgWidth, svgHeight, path.path)

	// Write to file
	filename := "test_trianglepath_single.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  Triangle path: %s\n", path.path)
	fmt.Printf("  Position: (%d, %d), Size: %d\n", topX, topY, triangleSize)
}

func testInverseTrianglePath() {
	fmt.Println("Testing InverseTrianglePath...")

	// Test with x=10, y=20, size=10
	path := InverseTrianglePath(10, 20, 10)

	// Verify we have 3 points
	if len(path.points) != 3 {
		fmt.Printf("FAIL: Expected 3 points, got %d\n", len(path.points))
		return
	}

	// Verify points are correct: top (15,20), bottom right (20,10), bottom left (10,10)
	expectedPoints := []Coord{
		{X: 15, Y: 20}, // top (apex)
		{X: 20, Y: 10}, // bottom right
		{X: 10, Y: 10}, // bottom left
	}

	for i, expected := range expectedPoints {
		if path.points[i].X != expected.X || path.points[i].Y != expected.Y {
			fmt.Printf("FAIL: Point %d expected (%d,%d), got (%d,%d)\n",
				i, expected.X, expected.Y, path.points[i].X, path.points[i].Y)
			return
		}
	}

	// Verify path string
	expectedPath := "M 15,20 L 20,10 L 10,10 Z"
	if path.path != expectedPath {
		fmt.Printf("FAIL: Expected path '%s', got '%s'\n", expectedPath, path.path)
		return
	}

	fmt.Println("PASS: InverseTrianglePath test passed!")
	fmt.Printf("  Points: %v\n", path.points)
	fmt.Printf("  Path: %s\n", path.path)
}

func testInverseTrianglePathSVG1() {
	fmt.Println("Testing InverseTrianglePath with SVG generation...")

	// Create multiple inverse triangles at different positions
	inverseTriangles := []struct {
		x, y, size int
		fill       string
	}{
		{50, 200, 30, "#ff6b6b"},
		{100, 200, 30, "#4ecdc4"},
		{150, 200, 30, "#45b7d1"},
		{200, 200, 30, "#f9ca24"},
		{250, 200, 30, "#6c5ce7"},
		{100, 150, 50, "#a29bfe"},
		{175, 100, 40, "#fd79a8"},
		{50, 100, 25, "#00b894"},
		{300, 150, 35, "#e17055"},
	}

	// Start building SVG content
	svgContent := `<?xml version="1.0"?>
<svg width="400" height="300" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="400" height="300" style="fill:#f8f9fa" />
`

	// Generate paths for each inverse triangle
	for i, t := range inverseTriangles {
		path := InverseTrianglePath(t.x, t.y, t.size)
		svgContent += fmt.Sprintf(`<path d="%s" style="fill:%s;stroke:#2d3436;stroke-width:2" />`+"\n", path.path, t.fill)
		fmt.Printf("  InverseTriangle %d: %s\n", i+1, path.path)
	}

	svgContent += `</svg>`

	// Write to file
	filename := "test_inversetrianglepath.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  Generated %d inverse triangles\n", len(inverseTriangles))
}

func testInverseTrianglePathSingleSVG() {
	fmt.Println("Testing InverseTrianglePath with single centered inverse triangle...")

	// SVG dimensions
	svgWidth := 400
	svgHeight := 300
	triangleSize := 150

	// Calculate center position
	centerX := svgWidth / 2
	centerY := svgHeight / 2

	// Calculate top position (InverseTrianglePath draws upward, so y is the top)
	topX := centerX - triangleSize/2
	topY := centerY + triangleSize/2

	// Generate the inverse triangle path
	path := InverseTrianglePath(topX, topY, triangleSize)

	// Build SVG content
	svgContent := fmt.Sprintf(`<?xml version="1.0"?>
<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="%d" height="%d" style="fill:#ffffff" />
<path d="%s" style="fill:#4a90e2;stroke:#2d3436;stroke-width:3" />
</svg>`, svgWidth, svgHeight, svgWidth, svgHeight, path.path)

	// Write to file
	filename := "test_inversetrianglepath_single.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  InverseTriangle path: %s\n", path.path)
	fmt.Printf("  Position: (%d, %d), Size: %d\n", topX, topY, triangleSize)
}

func testHalfCircleLeftHalfSquarePath() {
	fmt.Println("Testing HalfCircleLeftHalfSquarePath...")

	// Test with x=10, y=20, size=10
	path := HalfCircleLeftHalfSquarePath(10, 20, 10)

	// Verify we have 4 points
	if len(path.points) != 4 {
		fmt.Printf("FAIL: Expected 4 points, got %d\n", len(path.points))
		return
	}

	// Verify points are correct: top left (10,20), bottom left (10,10), bottom right (20,10), top right (20,20)
	expectedPoints := []Coord{
		{X: 10, Y: 20}, // top left
		{X: 10, Y: 10}, // bottom left
		{X: 20, Y: 10}, // bottom right
		{X: 20, Y: 20}, // top right
	}

	for i, expected := range expectedPoints {
		if path.points[i].X != expected.X || path.points[i].Y != expected.Y {
			fmt.Printf("FAIL: Point %d expected (%d,%d), got (%d,%d)\n",
				i, expected.X, expected.Y, path.points[i].X, path.points[i].Y)
			return
		}
	}

	// Verify path string contains arc command
	expectedPathStart := "M 10,20 A"
	if len(path.path) < len(expectedPathStart) || path.path[:len(expectedPathStart)] != expectedPathStart {
		fmt.Printf("FAIL: Expected path to start with '%s', got '%s'\n", expectedPathStart, path.path)
		return
	}

	// Verify path contains arc and lines
	if len(path.path) == 0 {
		fmt.Printf("FAIL: Path string is empty\n")
		return
	}

	fmt.Println("PASS: HalfCircleLeftHalfSquarePath test passed!")
	fmt.Printf("  Points: %v\n", path.points)
	fmt.Printf("  Path: %s\n", path.path)
}

func testHalfCircleLeftHalfSquarePathSVG1() {
	fmt.Println("Testing HalfCircleLeftHalfSquarePath with SVG generation...")

	// Create multiple half-circle-half-square shapes at different positions
	shapes := []struct {
		x, y, size int
		fill       string
	}{
		{50, 200, 30, "#ff6b6b"},
		{100, 200, 30, "#4ecdc4"},
		{150, 200, 30, "#45b7d1"},
		{200, 200, 30, "#f9ca24"},
		{250, 200, 30, "#6c5ce7"},
		{100, 150, 50, "#a29bfe"},
		{175, 100, 40, "#fd79a8"},
		{50, 100, 25, "#00b894"},
		{300, 150, 35, "#e17055"},
	}

	// Start building SVG content
	svgContent := `<?xml version="1.0"?>
<svg width="400" height="300" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="400" height="300" style="fill:#f8f9fa" />
`

	// Generate paths for each shape
	for i, s := range shapes {
		path := HalfCircleLeftHalfSquarePath(s.x, s.y, s.size)
		svgContent += fmt.Sprintf(`<path d="%s" style="fill:%s;stroke:#2d3436;stroke-width:2" />`+"\n", path.path, s.fill)
		fmt.Printf("  HalfCircleLeftHalfSquare %d: %s\n", i+1, path.path)
	}

	svgContent += `</svg>`

	// Write to file
	filename := "test_halfcirclelefthalfsquarepath.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  Generated %d half-circle-half-square shapes\n", len(shapes))
}

func testHalfCircleLeftHalfSquarePathSingleSVG() {
	fmt.Println("Testing HalfCircleLeftHalfSquarePath with single centered shape...")

	// SVG dimensions
	svgWidth := 400
	svgHeight := 300
	shapeSize := 150

	// Calculate center position
	centerX := svgWidth / 2
	centerY := svgHeight / 2

	// Calculate top-left position (HalfCircleLeftHalfSquarePath draws upward, so y is the top)
	topLeftX := centerX - shapeSize/2
	topLeftY := centerY + shapeSize/2

	// Generate the half-circle-half-square path
	path := HalfCircleLeftHalfSquarePath(topLeftX, topLeftY, shapeSize)

	// Build SVG content
	svgContent := fmt.Sprintf(`<?xml version="1.0"?>
<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="%d" height="%d" style="fill:#ffffff" />
<path d="%s" style="fill:#4a90e2;stroke:#2d3436;stroke-width:3" />
</svg>`, svgWidth, svgHeight, svgWidth, svgHeight, path.path)

	// Write to file
	filename := "test_halfcirclelefthalfsquarepath_single.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  HalfCircleLeftHalfSquare path: %s\n", path.path)
	fmt.Printf("  Position: (%d, %d), Size: %d\n", topLeftX, topLeftY, shapeSize)
}

func testHalfCircleTopHalfSquarePathSVG1() {
	fmt.Println("Testing HalfCircleTopHalfSquarePath with SVG generation...")

	// Create multiple half-circle-top-half-square shapes at different positions
	shapes := []struct {
		x, y, size int
		fill       string
	}{
		{50, 200, 30, "#ff6b6b"},
		{100, 200, 30, "#4ecdc4"},
		{150, 200, 30, "#45b7d1"},
		{200, 200, 30, "#f9ca24"},
		{250, 200, 30, "#6c5ce7"},
		{100, 150, 50, "#a29bfe"},
		{175, 100, 40, "#fd79a8"},
		{50, 100, 25, "#00b894"},
		{300, 150, 35, "#e17055"},
	}

	// Start building SVG content
	svgContent := `<?xml version="1.0"?>
<svg width="400" height="300" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="400" height="300" style="fill:#f8f9fa" />
`

	// Generate paths for each shape
	for i, s := range shapes {
		path := HalfCircleTopHalfSquarePath(s.x, s.y, s.size)
		svgContent += fmt.Sprintf(`<path d="%s" style="fill:%s;stroke:#2d3436;stroke-width:2" />`+"\n", path.path, s.fill)
		fmt.Printf("  HalfCircleTopHalfSquare %d: %s\n", i+1, path.path)
	}

	svgContent += `</svg>`

	// Write to file
	filename := "test_halfcircletophalfsquarepath.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  Generated %d half-circle-top-half-square shapes\n", len(shapes))
}

func testHalfCircleTopHalfSquarePathSingleSVG() {
	fmt.Println("Testing HalfCircleTopHalfSquarePath with single centered shape...")

	// SVG dimensions
	svgWidth := 400
	svgHeight := 300
	shapeSize := 150

	// Calculate center position
	centerX := svgWidth / 2
	centerY := svgHeight / 2

	// Calculate top-left position (HalfCircleTopHalfSquarePath draws upward, so y is the top)
	topLeftX := centerX - shapeSize/2
	topLeftY := centerY + shapeSize/2

	// Generate the half-circle-top-half-square path
	path := HalfCircleTopHalfSquarePath(topLeftX, topLeftY, shapeSize)

	// Build SVG content
	svgContent := fmt.Sprintf(`<?xml version="1.0"?>
<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="%d" height="%d" style="fill:#ffffff" />
<path d="%s" style="fill:#4a90e2;stroke:#2d3436;stroke-width:3" />
</svg>`, svgWidth, svgHeight, svgWidth, svgHeight, path.path)

	// Write to file
	filename := "test_halfcircletophalfsquarepath_single.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  HalfCircleTopHalfSquare path: %s\n", path.path)
	fmt.Printf("  Position: (%d, %d), Size: %d\n", topLeftX, topLeftY, shapeSize)
}

func testHalfCircleRightHalfSquarePathSVG1() {
	fmt.Println("Testing HalfCircleRightHalfSquarePath with SVG generation...")

	// Create multiple half-circle-right-half-square shapes at different positions
	shapes := []struct {
		x, y, size int
		fill       string
	}{
		{50, 200, 30, "#ff6b6b"},
		{100, 200, 30, "#4ecdc4"},
		{150, 200, 30, "#45b7d1"},
		{200, 200, 30, "#f9ca24"},
		{250, 200, 30, "#6c5ce7"},
		{100, 150, 50, "#a29bfe"},
		{175, 100, 40, "#fd79a8"},
		{50, 100, 25, "#00b894"},
		{300, 150, 35, "#e17055"},
	}

	// Start building SVG content
	svgContent := `<?xml version="1.0"?>
<svg width="400" height="300" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="400" height="300" style="fill:#f8f9fa" />
`

	// Generate paths for each shape
	for i, s := range shapes {
		path := HalfCircleRightHalfSquarePath(s.x, s.y, s.size)
		svgContent += fmt.Sprintf(`<path d="%s" style="fill:%s;stroke:#2d3436;stroke-width:2" />`+"\n", path.path, s.fill)
		fmt.Printf("  HalfCircleRightHalfSquare %d: %s\n", i+1, path.path)
	}

	svgContent += `</svg>`

	// Write to file
	filename := "test_halfcirclerighthalfsquarepath.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  Generated %d half-circle-right-half-square shapes\n", len(shapes))
}

func testHalfCircleRightHalfSquarePathSingleSVG() {
	fmt.Println("Testing HalfCircleRightHalfSquarePath with single centered shape...")

	// SVG dimensions
	svgWidth := 400
	svgHeight := 300
	shapeSize := 150

	// Calculate center position
	centerX := svgWidth / 2
	centerY := svgHeight / 2

	// Calculate top-left position (HalfCircleRightHalfSquarePath draws upward, so y is the top)
	topLeftX := centerX - shapeSize/2
	topLeftY := centerY + shapeSize/2

	// Generate the half-circle-right-half-square path
	path := HalfCircleRightHalfSquarePath(topLeftX, topLeftY, shapeSize)

	// Build SVG content
	svgContent := fmt.Sprintf(`<?xml version="1.0"?>
<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="%d" height="%d" style="fill:#ffffff" />
<path d="%s" style="fill:#4a90e2;stroke:#2d3436;stroke-width:3" />
</svg>`, svgWidth, svgHeight, svgWidth, svgHeight, path.path)

	// Write to file
	filename := "test_halfcirclerighthalfsquarepath_single.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  HalfCircleRightHalfSquare path: %s\n", path.path)
	fmt.Printf("  Position: (%d, %d), Size: %d\n", topLeftX, topLeftY, shapeSize)
}

func testHalfCircleBottomHalfSquarePathSVG1() {
	fmt.Println("Testing HalfCircleBottomHalfSquarePath with SVG generation...")

	// Create multiple half-circle-bottom-half-square shapes at different positions
	shapes := []struct {
		x, y, size int
		fill       string
	}{
		{50, 200, 30, "#ff6b6b"},
		{100, 200, 30, "#4ecdc4"},
		{150, 200, 30, "#45b7d1"},
		{200, 200, 30, "#f9ca24"},
		{250, 200, 30, "#6c5ce7"},
		{100, 150, 50, "#a29bfe"},
		{175, 100, 40, "#fd79a8"},
		{50, 100, 25, "#00b894"},
		{300, 150, 35, "#e17055"},
	}

	// Start building SVG content
	svgContent := `<?xml version="1.0"?>
<svg width="400" height="300" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="400" height="300" style="fill:#f8f9fa" />
`

	// Generate paths for each shape
	for i, s := range shapes {
		path := HalfCircleBottomHalfSquarePath(s.x, s.y, s.size)
		svgContent += fmt.Sprintf(`<path d="%s" style="fill:%s;stroke:#2d3436;stroke-width:2" />`+"\n", path.path, s.fill)
		fmt.Printf("  HalfCircleBottomHalfSquare %d: %s\n", i+1, path.path)
	}

	svgContent += `</svg>`

	// Write to file
	filename := "test_halfcirclebottomhalfsquarepath.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  Generated %d half-circle-bottom-half-square shapes\n", len(shapes))
}

func testHalfCircleBottomHalfSquarePathSingleSVG() {
	fmt.Println("Testing HalfCircleBottomHalfSquarePath with single centered shape...")

	// SVG dimensions
	svgWidth := 400
	svgHeight := 300
	shapeSize := 150

	// Calculate center position
	centerX := svgWidth / 2
	centerY := svgHeight / 2

	// Calculate top-left position (HalfCircleBottomHalfSquarePath draws upward, so y is the top)
	topLeftX := centerX - shapeSize/2
	topLeftY := centerY + shapeSize/2

	// Generate the half-circle-bottom-half-square path
	path := HalfCircleBottomHalfSquarePath(topLeftX, topLeftY, shapeSize)

	// Build SVG content
	svgContent := fmt.Sprintf(`<?xml version="1.0"?>
<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="%d" height="%d" style="fill:#ffffff" />
<path d="%s" style="fill:#4a90e2;stroke:#2d3436;stroke-width:3" />
</svg>`, svgWidth, svgHeight, svgWidth, svgHeight, path.path)

	// Write to file
	filename := "test_halfcirclebottomhalfsquarepath_single.svg"
	err := os.WriteFile(filename, []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("FAIL: Error writing SVG file: %v\n", err)
		return
	}

	fmt.Printf("PASS: SVG file '%s' created successfully!\n", filename)
	fmt.Printf("  HalfCircleBottomHalfSquare path: %s\n", path.path)
	fmt.Printf("  Position: (%d, %d), Size: %d\n", topLeftX, topLeftY, shapeSize)
}

// =========================
// Connected Squares Outline
// =========================

// ParseGridCoords parses input like "(0,2) (1,2)" into []Coord.
// Robust to extra spaces and optional signs.
func ParseGridCoords(input string) ([]Coord, error) {
	var coords []Coord
	i := 0
	n := len(input)

	parseInt := func() (int, bool) {
		// optional sign
		sign := 1
		if i < n && (input[i] == '+' || input[i] == '-') {
			if input[i] == '-' {
				sign = -1
			}
			i++
		}
		start := i
		for i < n && input[i] >= '0' && input[i] <= '9' {
			i++
		}
		if start == i {
			return 0, false
		}
		val, err := strconv.Atoi(input[start:i])
		if err != nil {
			return 0, false
		}
		return sign * val, true
	}

	skipSpaceComma := func() {
		for i < n && (input[i] == ' ' || input[i] == '\t' || input[i] == ',') {
			i++
		}
	}

	for i < n {
		if input[i] != '(' {
			i++
			continue
		}
		// consume '('
		i++
		skipSpaceComma()
		x, ok := parseInt()
		if !ok {
			return nil, fmt.Errorf("parse error: expected integer X near index %d", i)
		}
		skipSpaceComma()
		y, ok := parseInt()
		if !ok {
			return nil, fmt.Errorf("parse error: expected integer Y near index %d", i)
		}
		// move to ')'
		for i < n && input[i] != ')' {
			i++
		}
		if i >= n || input[i] != ')' {
			return nil, fmt.Errorf("parse error: missing closing ')' for coord (%d,%d)", x, y)
		}
		// consume ')'
		i++
		coords = append(coords, Coord{X: x, Y: y})
	}
	if len(coords) == 0 {
		return nil, fmt.Errorf("no coordinates found")
	}
	return coords, nil
}

// canonicalEdgeKey ensures consistent ordering of an edge's endpoints.
func canonicalEdgeKey(a, b Coord) [4]int {
	if a.X < b.X || (a.X == b.X && a.Y <= b.Y) {
		return [4]int{a.X, a.Y, b.X, b.Y}
	}
	return [4]int{b.X, b.Y, a.X, a.Y}
}

func addEdge(edges map[[4]int]bool, a, b Coord) {
	key := canonicalEdgeKey(a, b)
	if edges[key] {
		delete(edges, key) // toggle off when shared
	} else {
		edges[key] = true
	}
}

// buildBoundaryEdges builds outer boundary edges (in pixel coords), cancelling shared edges.
// Returns the edge set and a set of occupied cells for potential future use.
func buildBoundaryEdges(cells []Coord, size int) (map[[4]int]bool, map[Coord]bool) {
	edges := make(map[[4]int]bool)
	cellSet := make(map[Coord]bool, len(cells))
	for _, c := range cells {
		cellSet[c] = true
	}
	for _, c := range cells {
		x0 := c.X * size
		x1 := (c.X + 1) * size
		y0 := c.Y * size
		y1 := (c.Y + 1) * size

		tl := Coord{X: x0, Y: y1}
		tr := Coord{X: x1, Y: y1}
		br := Coord{X: x1, Y: y0}
		bl := Coord{X: x0, Y: y0}

		addEdge(edges, tl, tr) // top
		addEdge(edges, tr, br) // right
		addEdge(edges, br, bl) // bottom
		addEdge(edges, bl, tl) // left
	}
	return edges, cellSet
}

// normalize returns unit axis-aligned direction for delta between two vertices.
func normalize(from, to Coord) Coord {
	dx := to.X - from.X
	dy := to.Y - from.Y
	if dx > 0 {
		return Coord{X: 1, Y: 0} // East
	} else if dx < 0 {
		return Coord{X: -1, Y: 0} // West
	} else if dy > 0 {
		return Coord{X: 0, Y: 1} // North (up)
	}
	return Coord{X: 0, Y: -1} // South (down)
}

func rotateRight(d Coord) Coord {
	// cw order: E (1,0) -> S (0,-1) -> W (-1,0) -> N (0,1) -> E
	switch {
	case d.X == 1 && d.Y == 0:
		return Coord{X: 0, Y: -1}
	case d.X == 0 && d.Y == -1:
		return Coord{X: -1, Y: 0}
	case d.X == -1 && d.Y == 0:
		return Coord{X: 0, Y: 1}
	default: // 0,1
		return Coord{X: 1, Y: 0}
	}
}

func rotateLeft(d Coord) Coord {
	// ccw order: E -> N -> W -> S -> E
	switch {
	case d.X == 1 && d.Y == 0:
		return Coord{X: 0, Y: 1}
	case d.X == 0 && d.Y == 1:
		return Coord{X: -1, Y: 0}
	case d.X == -1 && d.Y == 0:
		return Coord{X: 0, Y: -1}
	default: // 0,-1
		return Coord{X: 1, Y: 0}
	}
}

func reverseDir(d Coord) Coord { return Coord{X: -d.X, Y: -d.Y} }

// orderBoundaryVertices walks the boundary edges clockwise keeping inside on the right.
func orderBoundaryVertices(edges map[[4]int]bool) ([]Coord, error) {
	if len(edges) == 0 {
		return nil, fmt.Errorf("no boundary edges")
	}
	// Build adjacency list
	adj := make(map[Coord][]Coord)
	for k := range edges {
		a := Coord{X: k[0], Y: k[1]}
		b := Coord{X: k[2], Y: k[3]}
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	// Find start: leftmost (min X), then topmost (max Y)
	var start Coord
	first := true
	for v := range adj {
		if first {
			start = v
			first = false
			continue
		}
		if v.X < start.X || (v.X == start.X && v.Y > start.Y) {
			start = v
		}
	}
	// Right-hand rule traversal
	cur := start
	dir := Coord{X: 0, Y: 1} // pretend we came from South -> facing North; prefer right=East first

	path := []Coord{start}
	used := make(map[[4]int]bool)
	maxSteps := len(edges)*2 + 10

	pickNext := func(cur, dir Coord) (Coord, Coord, bool) {
		candidates := []Coord{rotateRight(dir), dir, rotateLeft(dir), reverseDir(dir)}
		neighbors := adj[cur]
		// Prefer unused edges first
		for _, cd := range candidates {
			for _, nb := range neighbors {
				if normalize(cur, nb) == cd {
					key := canonicalEdgeKey(cur, nb)
					if !used[key] {
						return nb, cd, true
					}
				}
			}
		}
		// Fallback to any matching direction (even if used) to avoid dead-ends
		for _, cd := range candidates {
			for _, nb := range neighbors {
				if normalize(cur, nb) == cd {
					return nb, cd, true
				}
			}
		}
		return Coord{}, Coord{}, false
	}

	for steps := 0; steps < maxSteps; steps++ {
		next, ndir, ok := pickNext(cur, dir)
		if !ok {
			return nil, fmt.Errorf("failed to advance boundary walk")
		}
		used[canonicalEdgeKey(cur, next)] = true
		cur = next
		dir = ndir
		path = append(path, cur)
		if cur == start {
			break
		}
	}
	if cur != start {
		return nil, fmt.Errorf("boundary walk did not close")
	}
	// Drop the closing duplicate
	return path[:len(path)-1], nil
}

// compressColinear removes intermediate points on straight runs (axis-aligned).
func compressColinear(points []Coord) []Coord {
	n := len(points)
	if n <= 2 {
		return points
	}

	result := make([]Coord, 0, n)
	for i := 0; i < n; i++ {
		prev := points[(i-1+n)%n]
		cur := points[i]
		next := points[(i+1)%n]
		if (prev.X == cur.X && cur.X == next.X) || (prev.Y == cur.Y && cur.Y == next.Y) {
			continue // drop colinear middle
		}
		result = append(result, cur)
	}
	if len(result) < 3 {
		return points
	}
	return result
}

func generateSharpPath(points []Coord) string {
	if len(points) == 0 {
		return ""
	}
	path := fmt.Sprintf("M %d,%d", points[0].X, points[0].Y)
	for i := 1; i < len(points); i++ {
		path += fmt.Sprintf(" L %d,%d", points[i].X, points[i].Y)
	}
	path += " Z"
	return path
}

func manhattanLen(a, b Coord) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func moveToward(from, to Coord, dist int) Coord {
	// Assumes axis-aligned
	if to.X > from.X {
		return Coord{X: to.X - dist, Y: to.Y}
	}
	if to.X < from.X {
		return Coord{X: to.X + dist, Y: to.Y}
	}
	if to.Y > from.Y {
		return Coord{X: to.X, Y: to.Y - dist}
	}
	return Coord{X: to.X, Y: to.Y + dist}
}

// generateFilletedPath creates quarter-circle-like fillets at each corner using Q commands.
// Radius is expected pre-clamped to a reasonable global maximum (e.g., size/2).
func generateFilletedPath(points []Coord, radius int) string {
	n := len(points)
	if n == 0 {
		return ""
	}
	if radius <= 0 {
		return generateSharpPath(points)
	}

	c1 := make([]Coord, n)
	c2 := make([]Coord, n)
	for i := 0; i < n; i++ {
		prev := points[(i-1+n)%n]
		cur := points[i]
		next := points[(i+1)%n]
		// Clamp per-corner radius to segment lengths
		r := radius
		lenIn := manhattanLen(prev, cur)
		lenOut := manhattanLen(cur, next)
		if r > lenIn {
			r = lenIn
		}
		if r > lenOut {
			r = lenOut
		}
		if r < 0 {
			r = 0
		}
		// c1 is along incoming edge toward prev; c2 is along outgoing edge toward next
		c1[i] = moveToward(prev, cur, r)
		c2[i] = moveToward(next, cur, r)
	}

	// Start at first c2
	path := fmt.Sprintf("M %d,%d", c2[0].X, c2[0].Y)
	for i := 1; i < n; i++ {
		path += fmt.Sprintf(" L %d,%d", c1[i].X, c1[i].Y)
		path += fmt.Sprintf(" Q %d,%d %d,%d", points[i].X, points[i].Y, c2[i].X, c2[i].Y)
	}
	// Close back to 0
	path += fmt.Sprintf(" L %d,%d", c1[0].X, c1[0].Y)
	path += fmt.Sprintf(" Q %d,%d %d,%d", points[0].X, points[0].Y, c2[0].X, c2[0].Y)
	path += " Z"
	return path
}

// ConnectedSquaresOutlinePath converts connected grid cells into a single closed SVG Path.
// size: square size in pixels. fillet: radius in pixels (will be clamped to <= size/2).
func ConnectedSquaresOutlinePath(input string, size int, fillet int) (Path, error) {
	if size <= 0 {
		return Path{}, fmt.Errorf("size must be positive")
	}
	cells, err := ParseGridCoords(input)
	if err != nil {
		return Path{}, err
	}
	edges, _ := buildBoundaryEdges(cells, size)
	loop, err := orderBoundaryVertices(edges)
	if err != nil {
		return Path{}, err
	}
	loop = compressColinear(loop)

	radius := fillet
	if radius < 0 {
		radius = 0
	}
	if radius > size/2 {
		radius = size / 2
	}

	var pathStr string
	if radius > 0 {
		pathStr = generateFilletedPath(loop, radius)
	} else {
		pathStr = generateSharpPath(loop)
	}
	return Path{points: loop, path: pathStr}, nil
}

// Demo: Writes two SVG files, sharp and filleted, for a sample set of connected squares.
func testConnectedSquaresOutlineSVG() {
	fmt.Println("Testing ConnectedSquaresOutlinePath with SVG generation...")

	// Sample rectangle: 3x2 cells
	cellsStr := "(0,0) (1,0) (2,0) (0,1) (1,1) (2,1)"
	size := 40

	// Sharp
	shapeSharp, err := ConnectedSquaresOutlinePath(cellsStr, size, 0)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}
	// Filleted with radius = size/5
	shapeFillet, err := ConnectedSquaresOutlinePath(cellsStr, size, size/5)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}

	// Compute bounding box for SVG sizing
	minX, minY := shapeSharp.points[0].X, shapeSharp.points[0].Y
	maxX, maxY := minX, minY
	for _, p := range shapeSharp.points {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	margin := size
	width := maxX + margin
	height := maxY + margin

	makeSVG := func(d string) string {
		return fmt.Sprintf(`<?xml version="1.0"?>
<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
<rect x="0" y="0" width="%d" height="%d" style="fill:#ffffff" />
<path d="%s" style="fill:#4a90e2;stroke:#2d3436;stroke-width:2" />
</svg>`, width, height, width, height, d)
	}

	// Write files
	if err := os.WriteFile("test_connected_squares_outline.svg", []byte(makeSVG(shapeSharp.path)), 0644); err != nil {
		fmt.Printf("FAIL: writing sharp SVG: %v\n", err)
		return
	}
	if err := os.WriteFile("test_connected_squares_outline_filleted.svg", []byte(makeSVG(shapeFillet.path)), 0644); err != nil {
		fmt.Printf("FAIL: writing filleted SVG: %v\n", err)
		return
	}

	fmt.Println("PASS: ConnectedSquaresOutline SVGs created:")
	fmt.Println("  - test_connected_squares_outline.svg")
	fmt.Println("  - test_connected_squares_outline_filleted.svg")
}

func main() {
	// testing the Connected Component Analysis algorithm
	runCCATests()

	// testing SquarePath
	// testSquarePath()

	// testing SquarePath with SVG generation
	// testSquarePathSVG1()

	// testing SquarePath with single centered square
	// testSquarePathSingleSVG()

	// testDiamondPathSVG1()
	// testDiamondPathSingleSVG()

	// testInverseTrianglePathSVG1()
	// testInverseTrianglePathSingleSVG()

	// testTrianglePathSVG1()
	// testTrianglePathSingleSVG()

	// testHalfCircleLeftHalfSquarePathSVG1()
	// testHalfCircleLeftHalfSquarePathSingleSVG()

	// testHalfCircleTopHalfSquarePathSVG1()
	// testHalfCircleTopHalfSquarePathSingleSVG()

	// testHalfCircleRightHalfSquarePathSVG1()
	// testHalfCircleRightHalfSquarePathSingleSVG()

	// testHalfCircleBottomHalfSquarePathSVG1()
	// testHalfCircleBottomHalfSquarePathSingleSVG()

	// Connected Squares Outline SVG demo (sharp + filleted)
	// testConnectedSquaresOutlineSVG()
}
