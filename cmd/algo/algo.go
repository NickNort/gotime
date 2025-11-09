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

func main() {
	// Initialize a 2D bool array for testing
	bitmap1 := [][]bool{
		{true, false, true},
		{false, true, false},
		{true, true, false},
		{true, true, true},
	}

	bitmap2 := [][]bool{
		{true, false, true, true},
		{false, true, false, true},
		{true, true, false, true},
	}

	res1, err := algo(bitmap1)
	if err != nil {
		fmt.Printf("Error processing bitmap1: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap1 with groups:", m: bitmap1, coords: res1})
	printCoordMatrix(pcmOpts{label: "res1:", coords: res1})

	res2, err := algo(bitmap2)
	if err != nil {
		fmt.Printf("Error processing bitmap2: %v\n", err)
	}
	printBoolMatrix(pbmOpts{label: "bitmap2 with groups:", m: bitmap2, coords: res2})
	printCoordMatrix(pcmOpts{label: "res2:", coords: res2})
}
