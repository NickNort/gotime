package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"gotime/internal/qr"
)

// QRRequest represents the JSON request body for QR code generation
type QRRequest struct {
	Content         string `json:"content"`
	FinderCenter    string `json:"finder_center,omitempty"`
	FinderFrame     string `json:"finder_frame,omitempty"`
	ModuleShape     string `json:"module_shape,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`
	ForegroundColor string `json:"foreground_color,omitempty"`
}

// ErrorResponse represents an error response in JSON format
type ErrorResponse struct {
	Error string `json:"error"`
}

// handleQR handles POST requests to /qr endpoint
func handleQR(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Method not allowed. Use POST",
		})
		return
	}

	// Parse JSON request body
	var req QRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: fmt.Sprintf("Invalid JSON: %v", err),
		})
		return
	}

	// Validate required field
	if req.Content == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "content field is required",
		})
		return
	}

	// Generate QR code SVG
	svgBytes, err := qr.GenerateSVG(req.Content, qr.Options{
		FinderCenter:    req.FinderCenter,
		FinderFrame:     req.FinderFrame,
		ModuleShape:     req.ModuleShape,
		BackgroundColor: req.BackgroundColor,
		ForegroundColor: req.ForegroundColor,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Return SVG with proper content type
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	w.Write(svgBytes)
}

// handleHealth handles GET requests to /health endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Method not allowed. Use GET",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func main() {
	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Register HTTP handlers
	http.HandleFunc("/qr", handleQR)
	http.HandleFunc("/health", handleHealth)

	// Start HTTP server
	log.Printf("Starting HTTP server on port %s", port)
	log.Printf("POST /qr - Generate QR code")
	log.Printf("GET /health - Health check")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
