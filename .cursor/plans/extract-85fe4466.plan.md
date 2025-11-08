<!-- 85fe4466-084c-4fd0-875c-8d9722d775b4 7936b369-5525-49e0-93cd-343078fc5fad -->
# Extract shared QR generation to internal package

### Overview

- Create `internal/qr` package with shared QR generation (finder detection, rendering, validation, defaults, colors).
- Update `main.go` (HTTP API) and `cmd/climain.go` (CLI) to import `gotime/internal/qr` and call a single exported function.

### New package: `internal/qr/qr.go`

- Implement the shared logic now duplicated in both files:
  - `CornerRect`, `CornerBounds`, `isFinderPattern`, `findCorners`, `isInCorner`, `renderQR`.
  - Exported API with defaults/validation and fixed brand colors (#552048 foreground, #f8f2ec background):
```go
package qr

type Options struct {
    FinderCenter string // "circle" | "square" | "diamond"
    FinderFrame  string // "square" | "rounded" | "circle" | "diamond"
    ModuleShape  string // "square" | "rounded" | "circle" | "diamond"
    ModuleSize   int    // pixels per module; 0 => default 10
}

func GenerateSVG(content string, opts Options) ([]byte, error)
```

- Behavior:
  - Validate options; apply defaults if empty: FinderCenter="square", FinderFrame="square", ModuleShape="rounded", ModuleSize=10.
  - Generate QR (`go-qrcode`), detect corners, render via `svgo` into an SVG buffer with the mandated colors.

### Refactor `main.go` (HTTP API)

- Remove duplicated structs/functions and the local `generateQRCodeSVG`.
- Import `gotime/internal/qr`.
- In `handleQR`, call `qr.GenerateSVG(req.Content, qr.Options{FinderCenter: req.FinderCenter, FinderFrame: req.FinderFrame, ModuleShape: req.ModuleShape})` and return bytes with `image/svg+xml`.
- Keep `QRRequest`/`ErrorResponse`, handlers, and server boot logic unchanged.

### Refactor `cmd/climain.go` (CLI)

- Remove duplicated structs/functions and direct `qrcode`/`svgo` usage.
- Import `gotime/internal/qr`.
- After flag parsing, build `qr.Options` from flags and call `qr.GenerateSVG(qrContent, opts)`.
- Write the returned bytes to `qr.svg` (preserve file name); on error, print to stderr and exit non-zero.
- Rely on library validation (drop CLI-side duplication of allowed values checks).

### Notes

- Colors remain hardcoded in the library to satisfy workspace rule (#552048 foreground, #f8f2ec background).
- No public API exposure needed; `internal/qr` keeps it encapsulated inside the module.
- No behavior change expected beyond deduplication and centralized validation/defaults.

### To-dos

- [ ] Create internal/qr/qr.go with shared types, validation, GenerateSVG
- [ ] Update main.go to use gotime/internal/qr.GenerateSVG
- [ ] Update cmd/climain.go to use gotime/internal/qr.GenerateSVG and write qr.svg
- [ ] Build and smoke test CLI and HTTP API paths