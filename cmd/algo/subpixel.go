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

func HalfCircleLeftHalfSquarePath(x int, y int, size int) Path {
	path := new(Path)
	halfSize := size / 2
	radius := halfSize

	// Define key points
	topLeft := Coord{X: x, Y: y}                   // top left (start of semicircle)
	bottomLeft := Coord{X: x, Y: y - size}         // bottom left (end of semicircle)
	bottomRight := Coord{X: x + size, Y: y - size} // bottom right
	topRight := Coord{X: x + size, Y: y}           // top right

	path.points = append(path.points, topLeft)
	path.points = append(path.points, bottomLeft)
	path.points = append(path.points, bottomRight)
	path.points = append(path.points, topRight)

	// Start at top left
	path.path = fmt.Sprintf("M %d,%d", topLeft.X, topLeft.Y)
	// Arc to bottom left (semicircle on the left side, convex/inward)
	// A rx ry x-axis-rotation large-arc-flag sweep-flag x y
	// large-arc-flag=0 (semicircle), sweep-flag=1 (clockwise, curves inward/convex)
	path.path += fmt.Sprintf(" A %d,%d 0 0 1 %d,%d", radius, radius, bottomLeft.X, bottomLeft.Y)
	// Line to bottom right
	path.path += fmt.Sprintf(" L %d,%d", bottomRight.X, bottomRight.Y)
	// Line to top right
	path.path += fmt.Sprintf(" L %d,%d", topRight.X, topRight.Y)
	// Close path back to start
	path.path += " Z"

	return *path
}

func HalfCircleTopHalfSquarePath(x int, y int, size int) Path {
	path := new(Path)
	halfSize := size / 2
	radius := halfSize

	// Define key points
	topLeft := Coord{X: x, Y: y}                   // top left (start of semicircle)
	topRight := Coord{X: x + size, Y: y}           // top right (end of semicircle)
	bottomRight := Coord{X: x + size, Y: y - size} // bottom right
	bottomLeft := Coord{X: x, Y: y - size}         // bottom left

	path.points = append(path.points, topLeft)
	path.points = append(path.points, topRight)
	path.points = append(path.points, bottomRight)
	path.points = append(path.points, bottomLeft)

	// Start at top left
	path.path = fmt.Sprintf("M %d,%d", topLeft.X, topLeft.Y)
	// Arc to top right (semicircle on the top side, convex/inward)
	// A rx ry x-axis-rotation large-arc-flag sweep-flag x y
	// large-arc-flag=0 (semicircle), sweep-flag=0 (counterclockwise, curves inward/convex)
	path.path += fmt.Sprintf(" A %d,%d 0 0 0 %d,%d", radius, radius, topRight.X, topRight.Y)
	// Line to bottom right
	path.path += fmt.Sprintf(" L %d,%d", bottomRight.X, bottomRight.Y)
	// Line to bottom left
	path.path += fmt.Sprintf(" L %d,%d", bottomLeft.X, bottomLeft.Y)
	// Close path back to start
	path.path += " Z"

	return *path
}

func HalfCircleRightHalfSquarePath(x int, y int, size int) Path {
	path := new(Path)
	halfSize := size / 2
	radius := halfSize

	// Define key points
	topRight := Coord{X: x + size, Y: y}           // top right (start of semicircle)
	bottomRight := Coord{X: x + size, Y: y - size} // bottom right (end of semicircle)
	bottomLeft := Coord{X: x, Y: y - size}         // bottom left
	topLeft := Coord{X: x, Y: y}                   // top left

	path.points = append(path.points, topRight)
	path.points = append(path.points, bottomRight)
	path.points = append(path.points, bottomLeft)
	path.points = append(path.points, topLeft)

	// Start at top right
	path.path = fmt.Sprintf("M %d,%d", topRight.X, topRight.Y)
	// Arc to bottom right (semicircle on the right side, convex/inward)
	// A rx ry x-axis-rotation large-arc-flag sweep-flag x y
	// large-arc-flag=0 (semicircle), sweep-flag=0 (counterclockwise, curves inward/convex)
	path.path += fmt.Sprintf(" A %d,%d 0 0 0 %d,%d", radius, radius, bottomRight.X, bottomRight.Y)
	// Line to bottom left
	path.path += fmt.Sprintf(" L %d,%d", bottomLeft.X, bottomLeft.Y)
	// Line to top left
	path.path += fmt.Sprintf(" L %d,%d", topLeft.X, topLeft.Y)
	// Close path back to start
	path.path += " Z"

	return *path
}

func HalfCircleBottomHalfSquarePath(x int, y int, size int) Path {
	path := new(Path)
	halfSize := size / 2
	radius := halfSize

	// Define key points
	bottomLeft := Coord{X: x, Y: y - size}         // bottom left (start of semicircle)
	bottomRight := Coord{X: x + size, Y: y - size} // bottom right (end of semicircle)
	topRight := Coord{X: x + size, Y: y}           // top right
	topLeft := Coord{X: x, Y: y}                   // top left

	path.points = append(path.points, bottomLeft)
	path.points = append(path.points, bottomRight)
	path.points = append(path.points, topRight)
	path.points = append(path.points, topLeft)

	// Start at bottom left
	path.path = fmt.Sprintf("M %d,%d", bottomLeft.X, bottomLeft.Y)
	// Arc to bottom right (semicircle on the bottom side, convex/inward)
	// A rx ry x-axis-rotation large-arc-flag sweep-flag x y
	// large-arc-flag=0 (semicircle), sweep-flag=1 (clockwise, curves inward/convex)
	path.path += fmt.Sprintf(" A %d,%d 0 0 1 %d,%d", radius, radius, bottomRight.X, bottomRight.Y)
	// Line to top right
	path.path += fmt.Sprintf(" L %d,%d", topRight.X, topRight.Y)
	// Line to top left
	path.path += fmt.Sprintf(" L %d,%d", topLeft.X, topLeft.Y)
	// Close path back to start
	path.path += " Z"

	return *path
}
