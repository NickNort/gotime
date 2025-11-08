<!-- e9e3de2f-00ff-4217-8673-509b8a55c3e8 ee5ab20d-f210-4679-ab39-63a79233c047 -->
# Convert CLI to HTTP API

## Overview

Transform the command-line QR code generator into an HTTP API server that accepts generation requests via HTTP and returns SVG QR codes.

## Implementation Plan

### 1. Replace CLI with HTTP Server (`main.go`)

- Remove `flag` package usage and command-line argument parsing
- Add HTTP server using Go's standard `net/http` package
- Create HTTP handler function that:
- Accepts POST requests with JSON body containing:
- `content` (required): QR code content string
- `finder_center` (optional): corner center style (default: "square")
- `finder_frame` (optional): finder frame style (default: "square")
- `module_shape` (optional): module shape (default: "rounded")
- Validates input parameters
- Generates QR code using existing logic
- Returns SVG directly in response body with `Content-Type: image/svg+xml`
- Add GET endpoint for health check (optional but useful)
- Make server port configurable via environment variable (default: 8080)

### 2. Extract Core Logic

- Keep all existing QR generation functions (`isFinderPattern`, `findCorners`, `isInCorner`, `renderQR`) unchanged
- Extract QR generation logic from `main()` into a reusable function that:
- Takes parameters as function arguments instead of flags
- Returns SVG as `[]byte` instead of writing to file
- Can be called from HTTP handler

### 3. Request/Response Handling

- Define JSON request struct for API input
- Add proper error handling with HTTP status codes (400 for bad requests, 500 for server errors)
- Return JSON error responses for validation failures
- Return SVG directly for successful requests

### 4. Dependencies

- No new dependencies required - use standard library `net/http` and `encoding/json`
- Keep existing dependencies (`github.com/ajstarks/svgo`, `github.com/skip2/go-qrcode`)

## Files to Modify

- `main.go`: Refactor `main()` function, add HTTP handlers, extract QR generation logic

## API Design

- **POST /qr** - Generate QR code
- Request body: `{"content": "text", "finder_center": "square", "finder_frame": "square", "module_shape": "rounded"}`
- Response: SVG image (200 OK) or JSON error (400/500)
- **GET /health** - Health check endpoint (optional)

## Notes

- Maintains all existing QR code styling functionality
- Preserves color scheme (#552048 foreground, #f8f2ec background) per workspace rules
- SVG output format unchanged