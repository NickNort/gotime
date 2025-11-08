package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type TestResult struct {
	Name   string
	Passed bool
	Error  string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	baseURL := flag.String("base", "", "Base URL for the API (default: http://localhost:8080)")
	flag.Parse()

	if *baseURL == "" {
		*baseURL = os.Getenv("BASE_URL")
		if *baseURL == "" {
			*baseURL = "http://localhost:8080"
		}
	}

	// Remove trailing slash if present
	*baseURL = strings.TrimSuffix(*baseURL, "/")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var results []TestResult

	fmt.Println("Testing gotime API endpoints...")
	fmt.Println("Base URL:", *baseURL)
	fmt.Println()

	// Test GET /health
	results = append(results, testHealthGET(client, *baseURL))

	// Test POST /health (should fail)
	results = append(results, testHealthPOST(client, *baseURL))

	// Test POST /qr success (default colors)
	results = append(results, testQRSuccess(client, *baseURL))

	// Test POST /qr with custom colors
	results = append(results, testQRCustomColors(client, *baseURL))

	// Test POST /qr with partial custom colors
	results = append(results, testQRPartialColors(client, *baseURL))

	// Test POST /qr missing content
	results = append(results, testQRMissingContent(client, *baseURL))

	// Test POST /qr invalid JSON
	results = append(results, testQRInvalidJSON(client, *baseURL))

	// Test POST /qr invalid params
	results = append(results, testQRInvalidParams(client, *baseURL))

	// Print results
	fmt.Println()
	fmt.Println("Results:")
	fmt.Println(strings.Repeat("-", 60))
	passed := 0
	for _, result := range results {
		status := "FAIL"
		if result.Passed {
			status = "PASS"
			passed++
		}
		fmt.Printf("[%s] %s\n", status, result.Name)
		if !result.Passed && result.Error != "" {
			fmt.Printf("      Error: %s\n", result.Error)
		}
	}
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("PASS %d / %d\n", passed, len(results))

	if passed < len(results) {
		os.Exit(1)
	}
}

func testHealthGET(client *http.Client, baseURL string) TestResult {
	name := "GET /health"
	resp, err := client.Get(baseURL + "/health")
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected 200, got %d", resp.StatusCode)}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected Content-Type application/json, got %s", contentType)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to read body: %v", err)}
	}

	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to parse JSON: %v", err)}
	}

	if data["status"] != "ok" {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected status 'ok', got %q", data["status"])}
	}

	return TestResult{Name: name, Passed: true}
}

func testHealthPOST(client *http.Client, baseURL string) TestResult {
	name := "POST /health (method not allowed)"
	resp, err := client.Post(baseURL+"/health", "application/json", strings.NewReader("{}"))
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected 405, got %d", resp.StatusCode)}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected Content-Type application/json, got %s", contentType)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to read body: %v", err)}
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to parse JSON: %v", err)}
	}

	if !strings.Contains(errorResp.Error, "Method not allowed") || !strings.Contains(errorResp.Error, "GET") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected error message about method not allowed, got %q", errorResp.Error)}
	}

	return TestResult{Name: name, Passed: true}
}

func testQRSuccess(client *http.Client, baseURL string) TestResult {
	name := "POST /qr (success with default colors)"
	reqBody := `{"content":"https://example.com"}`
	resp, err := client.Post(baseURL+"/qr", "application/json", strings.NewReader(reqBody))
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected 200, got %d", resp.StatusCode)}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/svg+xml") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected Content-Type image/svg+xml, got %s", contentType)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to read body: %v", err)}
	}

	bodyStr := string(body)
	if !strings.Contains(bodyStr, "<svg") {
		return TestResult{Name: name, Passed: false, Error: "Response body does not contain <svg"}
	}

	// Check for default colors (workspace rule)
	if !strings.Contains(bodyStr, "#552048") {
		return TestResult{Name: name, Passed: false, Error: "SVG does not contain default foreground color #552048"}
	}

	if !strings.Contains(bodyStr, "#f8f2ec") {
		return TestResult{Name: name, Passed: false, Error: "SVG does not contain default background color #f8f2ec"}
	}

	return TestResult{Name: name, Passed: true}
}

func testQRCustomColors(client *http.Client, baseURL string) TestResult {
	name := "POST /qr (custom colors)"
	reqBody := `{"content":"https://example.com","background_color":"#ffffff","foreground_color":"#000000"}`
	resp, err := client.Post(baseURL+"/qr", "application/json", strings.NewReader(reqBody))
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected 200, got %d", resp.StatusCode)}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/svg+xml") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected Content-Type image/svg+xml, got %s", contentType)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to read body: %v", err)}
	}

	bodyStr := string(body)
	if !strings.Contains(bodyStr, "<svg") {
		return TestResult{Name: name, Passed: false, Error: "Response body does not contain <svg"}
	}

	// Check for custom colors
	if !strings.Contains(bodyStr, "#000000") {
		return TestResult{Name: name, Passed: false, Error: "SVG does not contain custom foreground color #000000"}
	}

	if !strings.Contains(bodyStr, "#ffffff") {
		return TestResult{Name: name, Passed: false, Error: "SVG does not contain custom background color #ffffff"}
	}

	// Verify default colors are NOT present
	if strings.Contains(bodyStr, "#552048") {
		return TestResult{Name: name, Passed: false, Error: "SVG still contains default foreground color #552048 instead of custom color"}
	}

	if strings.Contains(bodyStr, "#f8f2ec") {
		return TestResult{Name: name, Passed: false, Error: "SVG still contains default background color #f8f2ec instead of custom color"}
	}

	return TestResult{Name: name, Passed: true}
}

func testQRPartialColors(client *http.Client, baseURL string) TestResult {
	name := "POST /qr (partial custom colors - only foreground)"
	reqBody := `{"content":"https://example.com","foreground_color":"#ff0000"}`
	resp, err := client.Post(baseURL+"/qr", "application/json", strings.NewReader(reqBody))
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected 200, got %d", resp.StatusCode)}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/svg+xml") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected Content-Type image/svg+xml, got %s", contentType)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to read body: %v", err)}
	}

	bodyStr := string(body)
	if !strings.Contains(bodyStr, "<svg") {
		return TestResult{Name: name, Passed: false, Error: "Response body does not contain <svg"}
	}

	// Check for custom foreground color
	if !strings.Contains(bodyStr, "#ff0000") {
		return TestResult{Name: name, Passed: false, Error: "SVG does not contain custom foreground color #ff0000"}
	}

	// Check for default background color (should still be used)
	if !strings.Contains(bodyStr, "#f8f2ec") {
		return TestResult{Name: name, Passed: false, Error: "SVG does not contain default background color #f8f2ec"}
	}

	return TestResult{Name: name, Passed: true}
}

func testQRMissingContent(client *http.Client, baseURL string) TestResult {
	name := "POST /qr (missing content field)"
	resp, err := client.Post(baseURL+"/qr", "application/json", strings.NewReader("{}"))
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected 400, got %d", resp.StatusCode)}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected Content-Type application/json, got %s", contentType)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to read body: %v", err)}
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to parse JSON: %v", err)}
	}

	if !strings.Contains(errorResp.Error, "content field is required") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected error about content field required, got %q", errorResp.Error)}
	}

	return TestResult{Name: name, Passed: true}
}

func testQRInvalidJSON(client *http.Client, baseURL string) TestResult {
	name := "POST /qr (invalid JSON)"
	resp, err := client.Post(baseURL+"/qr", "application/json", strings.NewReader("{"))
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected 400, got %d", resp.StatusCode)}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected Content-Type application/json, got %s", contentType)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to read body: %v", err)}
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to parse JSON: %v", err)}
	}

	if !strings.Contains(strings.ToLower(errorResp.Error), "invalid json") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected error containing 'Invalid JSON', got %q", errorResp.Error)}
	}

	return TestResult{Name: name, Passed: true}
}

func testQRInvalidParams(client *http.Client, baseURL string) TestResult {
	name := "POST /qr (invalid finder_center param)"
	reqBody := `{"content":"https://example.com","finder_center":"triangle"}`
	resp, err := client.Post(baseURL+"/qr", "application/json", strings.NewReader(reqBody))
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected 400, got %d", resp.StatusCode)}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected Content-Type application/json, got %s", contentType)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to read body: %v", err)}
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Failed to parse JSON: %v", err)}
	}

	// Check that error mentions allowed values (circle, square, diamond)
	errorLower := strings.ToLower(errorResp.Error)
	if !strings.Contains(errorLower, "finder-center") && !strings.Contains(errorLower, "finder_center") {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected error mentioning finder-center, got %q", errorResp.Error)}
	}

	// Check that it mentions at least one allowed value
	hasAllowedValue := strings.Contains(errorLower, "circle") ||
		strings.Contains(errorLower, "square") ||
		strings.Contains(errorLower, "diamond")
	if !hasAllowedValue {
		return TestResult{Name: name, Passed: false, Error: fmt.Sprintf("Expected error to mention allowed values, got %q", errorResp.Error)}
	}

	return TestResult{Name: name, Passed: true}
}

