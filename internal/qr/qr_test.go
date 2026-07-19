package qr

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/skip2/go-qrcode"
)

func TestFindCornersAcrossSmallQRVersions(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantVersion int
		wantSize    int
		wantCorners CornerBounds
	}{
		{
			name:        "version one",
			content:     "x.com",
			wantVersion: 1,
			wantSize:    29,
			wantCorners: CornerBounds{
				TopLeft:    CornerRect{X: 4, Y: 4, Width: 7, Height: 7},
				TopRight:   CornerRect{X: 18, Y: 4, Width: 7, Height: 7},
				BottomLeft: CornerRect{X: 4, Y: 18, Width: 7, Height: 7},
			},
		},
		{
			name:        "version two",
			content:     "https://x.co",
			wantVersion: 2,
			wantSize:    33,
			wantCorners: CornerBounds{
				TopLeft:    CornerRect{X: 4, Y: 4, Width: 7, Height: 7},
				TopRight:   CornerRect{X: 22, Y: 4, Width: 7, Height: 7},
				BottomLeft: CornerRect{X: 4, Y: 22, Width: 7, Height: 7},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, err := qrcode.New(test.content, qrcode.Highest)
			if err != nil {
				t.Fatalf("qrcode.New() error = %v", err)
			}
			if code.VersionNumber != test.wantVersion {
				t.Fatalf("version = %d; want %d", code.VersionNumber, test.wantVersion)
			}

			bitmap := code.Bitmap()
			if len(bitmap) != test.wantSize {
				t.Fatalf("bitmap size = %d; want %d", len(bitmap), test.wantSize)
			}

			corners, err := findCorners(bitmap)
			if err != nil {
				t.Fatalf("findCorners() error = %v", err)
			}
			if corners != test.wantCorners {
				t.Fatalf("findCorners() = %+v; want %+v", corners, test.wantCorners)
			}
		})
	}
}

func TestFindCornersAcrossAllQRVersions(t *testing.T) {
	for version := 1; version <= 40; version++ {
		t.Run(fmt.Sprintf("version %d", version), func(t *testing.T) {
			code, err := qrcode.NewWithForcedVersion("x.com", version, qrcode.Highest)
			if err != nil {
				t.Fatalf("qrcode.NewWithForcedVersion() error = %v", err)
			}

			bitmap := code.Bitmap()
			wantSize := 21 + (version-1)*4 + 8
			if len(bitmap) != wantSize {
				t.Fatalf("bitmap size = %d; want %d", len(bitmap), wantSize)
			}

			farFinderOffset := wantSize - 4 - 7
			wantCorners := CornerBounds{
				TopLeft:    CornerRect{X: 4, Y: 4, Width: 7, Height: 7},
				TopRight:   CornerRect{X: farFinderOffset, Y: 4, Width: 7, Height: 7},
				BottomLeft: CornerRect{X: 4, Y: farFinderOffset, Width: 7, Height: 7},
			}

			corners, err := findCorners(bitmap)
			if err != nil {
				t.Fatalf("findCorners() error = %v", err)
			}
			if corners != wantCorners {
				t.Fatalf("findCorners() = %+v; want %+v", corners, wantCorners)
			}
		})
	}
}

func TestGenerateSVGVersionOneDoesNotDrawFinderAtOrigin(t *testing.T) {
	svgBytes, err := GenerateSVG("x.com", Options{
		FinderCenter: "square",
		FinderFrame:  "square",
		ModuleShape:  "square",
	})
	if err != nil {
		t.Fatalf("GenerateSVG() error = %v", err)
	}

	bogusFinderAtOrigin := []byte(`<rect x="0" y="0" width="70" height="10" style="fill:#552048" />`)
	if bytes.Contains(svgBytes, bogusFinderAtOrigin) {
		t.Fatal("GenerateSVG() drew a finder pattern over the quiet zone at (0, 0)")
	}
}
