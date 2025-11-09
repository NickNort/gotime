package main

import (
	"fmt"
)

type Path struct {
	points []Coord
	path   string
}

func SquarePath(x int, y int, size int) Path {
	path := new(Path)
	path.points = append(path.points, Coord{X: x, Y: y})                           // top left
	path.points = append(path.points, Coord{X: x + size, Y: y})              // top right
	path.points = append(path.points, Coord{X: x + size, Y: y - size}) // bottom right
	path.points = append(path.points, Coord{X: x, Y: y - size})              // bottom left

	// start at top left
	path.path = fmt.Sprintf("M %d,%d", path.points[0].X, path.points[0].Y)
	// add a line for each subsequent point
	for i := 1; i < len(path.points); i++ {
		path.path += fmt.Sprintf(" L %d,%d", path.points[i].X, path.points[i].Y)
	}
	// connects back to start
	path.path += " Z"

	return *path
}
