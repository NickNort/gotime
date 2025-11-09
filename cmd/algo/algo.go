package main

import (
	"fmt"
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

func main() {
	// testing the Connected Component Analysis algorithm
	// runCCATests()

	// testing SquarePath
	// testSquarePath()

	// testing SquarePath with SVG generation
	// testSquarePathSVG1()

	// testing SquarePath with single centered square
	testSquarePathSingleSVG()

	testDiamondPathSVG1()
	testDiamondPathSingleSVG()

	testInverseTrianglePathSVG1()
	testInverseTrianglePathSingleSVG()

	testTrianglePathSVG1()
	testTrianglePathSingleSVG()
}
