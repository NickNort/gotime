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
	path.points = append(path.points, Coord{X: x, Y: y})               // top left
	path.points = append(path.points, Coord{X: x + size, Y: y})        // top right
	path.points = append(path.points, Coord{X: x + size, Y: y - size}) // bottom right
	path.points = append(path.points, Coord{X: x, Y: y - size})        // bottom left

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

func DiamondPath(x int, y int, size int) Path {
	path := new(Path)
	halfSize := size / 2
	path.points = append(path.points, Coord{X: x + halfSize, Y: y})        // top
	path.points = append(path.points, Coord{X: x + size, Y: y - halfSize}) // right
	path.points = append(path.points, Coord{X: x + halfSize, Y: y - size}) // bottom
	path.points = append(path.points, Coord{X: x, Y: y - halfSize})        // left

	// start at top point
	path.path = fmt.Sprintf("M %d,%d", path.points[0].X, path.points[0].Y)
	// add a line for each subsequent point
	for i := 1; i < len(path.points); i++ {
		path.path += fmt.Sprintf(" L %d,%d", path.points[i].X, path.points[i].Y)
	}
	// connects back to start
	path.path += " Z"

	return *path
}

func InverseTrianglePath(x int, y int, size int) Path {
	path := new(Path)
	halfSize := size / 2
	path.points = append(path.points, Coord{X: x + halfSize, Y: y})    // top (apex)
	path.points = append(path.points, Coord{X: x + size, Y: y - size}) // bottom right
	path.points = append(path.points, Coord{X: x, Y: y - size})        // bottom left

	// start at top point
	path.path = fmt.Sprintf("M %d,%d", path.points[0].X, path.points[0].Y)
	// add a line for each subsequent point
	for i := 1; i < len(path.points); i++ {
		path.path += fmt.Sprintf(" L %d,%d", path.points[i].X, path.points[i].Y)
	}
	// connects back to start
	path.path += " Z"

	return *path
}

func TrianglePath(x int, y int, size int) Path {
	path := new(Path)
	halfSize := size / 2
	path.points = append(path.points, Coord{X: x, Y: y})                   // top left (base)
	path.points = append(path.points, Coord{X: x + size, Y: y})            // top right (base)
	path.points = append(path.points, Coord{X: x + halfSize, Y: y - size}) // bottom (apex)

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
