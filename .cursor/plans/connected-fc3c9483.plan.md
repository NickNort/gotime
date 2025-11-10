<!-- fc3c9483-4a69-4dd6-8a5e-902d27fbeddc 5830bc8f-556f-462e-b16e-9be78d82d2ea -->
# Connected Squares Outline SVG Path

### What we’ll add

- A boundary-following algorithm that converts connected grid cells into a single closed SVG path (sharp or filleted corners).
- Parsing helper for input like "(0,2) (1,2) ...".
- Minimal test helpers to emit demo SVGs (commented by default), consistent with existing test style.

### Files to change

- `cmd/algo/algo.go` (new functions and optional test harness)

### Reuse existing types

We will reuse these existing types so the new API plugs in seamlessly:

```8:10:cmd/algo/algo.go
type Coord struct {
	X, Y int
}
```
```7:10:cmd/algo/subpixel.go
type Path struct {
	points []Coord
	path   string
}
```

### New API (signatures)

- `ParseGridCoords(input string) ([]Coord, error)`
- `ConnectedSquaresOutlinePath(input string, size int, fillet int) (Path, error)`
  - Internally:
    - `buildBoundaryEdges(cells []Coord, size int) (map[[4]int]bool, map[Coord]bool)`
    - `orderBoundaryVertices(edges map[[4]int]bool) ([]Coord, error)`
    - `compressColinear(points []Coord) []Coord`
    - `generateSharpPath(points []Coord) string`
    - `generateFilletedPath(points []Coord, radius int) string`

### Algorithm details

- **Parse**: Split by whitespace; accept tokens like `(x,y)`; robust to extra commas/spaces.
- **Cell set**: Use `map[Coord]bool` for O(1) neighbor checks.
- **Boundary edges**: For each occupied cell, add its 4 edges; cancel shared edges with neighboring cells by storing a canonical edge key.
  - Canonical edge key: `[4]int{x1,y1,x2,y2}` with `(x1,y1) <= (x2,y2)` lexicographically.
- **Order vertices**: Build adjacency `map[Coord][]Coord` from remaining edges, select start as leftmost–topmost (min X, then max Y), then walk edges clockwise keeping inside on right.
  - Direction set: East, South, West, North (cw), with priority right > straight > left > back.
- **Simplify**: Collapse colinear runs to only true corners.
- **Path generation**:
  - Sharp: `M x0,y0 L ... Z`.
  - Filleted: clamp `fillet <= size/2`, and per-corner clamp to segment lengths; for each corner use `L c1 Q cx,cy c2` where c1/c2 are radius-offset points on incident edges.

### Essential snippet (edge building idea)

```go
// For each occupied grid cell (cx, cy):
// top-left pixel is (cx*size, (cy+1)*size) to match existing upward-drawing Y
// Define 4 edges; toggle them in a map to cancel internal edges
func addEdge(edges map[[4]int]bool, a, b Coord) {
  key := canonical(a, b) // [4]int{x1,y1,x2,y2}
  if edges[key] { delete(edges, key) } else { edges[key] = true }
}
```

### Test/demo (commented in main)

- `testConnectedSquaresOutlineSVG()` that:
  - Uses input from your example.
  - Emits `test_connected_squares_outline.svg` (sharp) and `test_connected_squares_outline_filleted.svg` (fillet = size/5).

### Error handling and constraints

- Empty input → error.
- Nonpositive `size` → error.
- `fillet` auto-clamped to `<= size/2` and per-corner segment lengths.
- Single-cell and straight-line cases covered.

### To-dos

- [ ] Implement ParseGridCoords for "(x,y)" tokens
- [ ] Build boundary edges, cancel shared edges
- [ ] Order boundary vertices clockwise using right-hand rule
- [ ] Compress colinear vertices to true corners
- [ ] Generate sharp SVG path from vertices
- [ ] Generate filleted path with Q commands and clamped radius
- [ ] Wire up ConnectedSquaresOutlinePath(input,size,fillet) returning Path
- [ ] Add demo function to write sample SVGs (commented in main)