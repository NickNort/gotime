package main

import (
	"image/color"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	// Generate QR code
	qrc, err := qrcode.New("https://nickn.dev")
	if err != nil {
		panic(err)
	}

	// Create writer with custom options
	w, err := standard.New("output.png",
		standard.WithBgColor(color.RGBA{R: 248, G: 242, B: 236, A: 255}), // Light beige background
		standard.WithFgColor(color.RGBA{R: 85, G: 32, B: 72, A: 255}),    // Dark purple foreground
		standard.WithCircleShape(),
		standard.WithQRWidth(128), // Size of QR code
	)
	if err != nil {
		panic(err)
	}

	// Save the QR code
	if err = qrc.Save(w); err != nil {
		panic(err)
	}
}
