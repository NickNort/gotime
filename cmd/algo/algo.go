package main

import (
	"fmt"
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
	printBoolMatrix(pbmOpts{label: "visitedMatrix:", m: visitedMatrix})
	print("sizeY: ", sizeY, "\n")

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

func runTests() {
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

func main() {
	// runTests()
	
}
