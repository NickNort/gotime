package algo

import (
	"fmt"
)

// Coord represents a coordinate in a 2D bitmap
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

// Algo finds all connected groups of true values in a bitmap using DFS
// Returns a slice of coordinate groups, where each group contains coordinates
// of connected true values
func Algo(bitmap [][]bool) ([][]Coord, error) {
	sizeY := len(bitmap)
	if sizeY <= 0 {
		return nil, fmt.Errorf("bitmap cannot be empty")
	}

	sizeX := len(bitmap[0])
	if sizeX <= 0 {
		return nil, fmt.Errorf("bitmap rows cannot be empty")
	}

	visitedMatrix := makeVisitedMatrix(bitmap)

	res := new([][]Coord)
	var currentGroup *[]Coord

	for y, r := range bitmap {
		for x, c := range r {
			if c == true && visitedMatrix[y][x] == false {
				currentGroup = new([]Coord)
				dfs(bitmap, visitedMatrix, x, y, currentGroup)
				*res = append(*res, *currentGroup)
			}
		}
	}

	return *res, nil
}

func dfs(bitmap [][]bool, visitedMatrix [][]bool, x int, y int, currentGroup *[]Coord) bool {
	if y >= len(bitmap) || y < 0 || x >= len(bitmap[y]) || x < 0 || visitedMatrix[y][x] == true || bitmap[y][x] == false {
		return false
	}
	visitedMatrix[y][x] = true
	*currentGroup = append(*currentGroup, Coord{X: x, Y: y})

	dfs(bitmap, visitedMatrix, x-1, y, currentGroup)
	dfs(bitmap, visitedMatrix, x+1, y, currentGroup)
	dfs(bitmap, visitedMatrix, x, y-1, currentGroup)
	dfs(bitmap, visitedMatrix, x, y+1, currentGroup)

	return true
}
