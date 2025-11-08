package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var cliPath string

func TestMain(m *testing.M) {
	// Resolve repo root (dir of this test, go up to cmd, then up to root)
	_, testFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFile)
	repoRoot := filepath.Join(testDir, "..")

	// Create temp directory for the binary
	tmpDir, err := os.MkdirTemp("", "climain-test-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(tmpDir)

	// Build the CLI
	binaryName := "climain"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	cliPath = filepath.Join(tmpDir, binaryName)

	buildCmd := exec.Command("go", "build", "-o", cliPath, "./cmd/climain.go")
	buildCmd.Dir = repoRoot
	if err := buildCmd.Run(); err != nil {
		panic("failed to build CLI: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Exit(code)
}

func runCLI(t *testing.T, workDir string, args ...string) (int, string, string) {
	t.Helper()
	cmd := exec.Command(cliPath, args...)
	cmd.Dir = workDir
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	runErr := cmd.Run()
	code := 0
	if runErr != nil {
		if ee, ok := runErr.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			t.Fatalf("failed to run: %v", runErr)
		}
	}
	return code, out.String(), errb.String()
}

func TestNoArgsPrintsUsageAndExitsNonZero(t *testing.T) {
	workDir, err := os.MkdirTemp("", "climain-test-work-*")
	if err != nil {
		t.Fatalf("failed to create work dir: %v", err)
	}
	defer os.RemoveAll(workDir)

	exitCode, stdout, stderr := runCLI(t, workDir)

	if exitCode == 0 {
		t.Errorf("expected non-zero exit code, got %d", exitCode)
	}

	if !strings.Contains(stderr, "Usage:") {
		t.Errorf("expected usage message in stderr, got: %q", stderr)
	}

	// Verify qr.svg was not created
	qrPath := filepath.Join(workDir, "qr.svg")
	if _, err := os.Stat(qrPath); err == nil {
		t.Errorf("qr.svg should not have been created, but it exists")
	}

	if stdout != "" {
		t.Errorf("expected empty stdout, got: %q", stdout)
	}
}

func TestValidRunCreatesQRSVGWithRequiredColors(t *testing.T) {
	workDir, err := os.MkdirTemp("", "climain-test-work-*")
	if err != nil {
		t.Fatalf("failed to create work dir: %v", err)
	}
	defer os.RemoveAll(workDir)

	exitCode, stdout, stderr := runCLI(t, workDir, "hello")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d. stderr: %q", exitCode, stderr)
	}

	if stderr != "" {
		t.Errorf("expected empty stderr, got: %q", stderr)
	}

	// Verify qr.svg was created
	qrPath := filepath.Join(workDir, "qr.svg")
	svgBytes, err := os.ReadFile(qrPath)
	if err != nil {
		t.Fatalf("failed to read qr.svg: %v", err)
	}

	svgContent := string(svgBytes)

	// Check for SVG structure
	if !strings.Contains(svgContent, "<svg") {
		t.Errorf("qr.svg does not contain <svg tag")
	}

	// Check for required colors
	if !strings.Contains(svgContent, "#552048") {
		t.Errorf("qr.svg does not contain foreground color #552048")
	}

	if !strings.Contains(svgContent, "#f8f2ec") {
		t.Errorf("qr.svg does not contain background color #f8f2ec")
	}

	if stdout != "" {
		t.Errorf("expected empty stdout, got: %q", stdout)
	}
}

func TestFlagsAfterPositionalArg(t *testing.T) {
	workDir, err := os.MkdirTemp("", "climain-test-work-*")
	if err != nil {
		t.Fatalf("failed to create work dir: %v", err)
	}
	defer os.RemoveAll(workDir)

	exitCode, stdout, stderr := runCLI(t, workDir, "hello", "-c", "circle", "-f", "rounded", "-m", "diamond")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d. stderr: %q", exitCode, stderr)
	}

	if stderr != "" {
		t.Errorf("expected empty stderr, got: %q", stderr)
	}

	// Verify qr.svg was created
	qrPath := filepath.Join(workDir, "qr.svg")
	if _, err := os.Stat(qrPath); err != nil {
		t.Errorf("qr.svg should have been created, but it doesn't exist: %v", err)
	}

	if stdout != "" {
		t.Errorf("expected empty stdout, got: %q", stdout)
	}
}

func TestLongShortFlagVariantsWithEqualsAndSpace(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{"long flag with equals", []string{"test", "-finder-center=circle"}},
		{"short flag with equals", []string{"test", "-c=circle"}},
		{"long flag with space", []string{"test", "-finder-center", "circle"}},
		{"short flag with space", []string{"test", "-c", "circle"}},
		{"finder-frame long equals", []string{"test", "-finder-frame=rounded"}},
		{"finder-frame short equals", []string{"test", "-f=rounded"}},
		{"finder-frame long space", []string{"test", "-finder-frame", "rounded"}},
		{"finder-frame short space", []string{"test", "-f", "rounded"}},
		{"module-shape long equals", []string{"test", "-module-shape=diamond"}},
		{"module-shape short equals", []string{"test", "-m=diamond"}},
		{"module-shape long space", []string{"test", "-module-shape", "diamond"}},
		{"module-shape short space", []string{"test", "-m", "diamond"}},
		{"all flags long equals", []string{"test", "-finder-center=circle", "-finder-frame=rounded", "-module-shape=diamond"}},
		{"all flags short equals", []string{"test", "-c=circle", "-f=rounded", "-m=diamond"}},
		{"all flags long space", []string{"test", "-finder-center", "circle", "-finder-frame", "rounded", "-module-shape", "diamond"}},
		{"all flags short space", []string{"test", "-c", "circle", "-f", "rounded", "-m", "diamond"}},
		{"mixed formats", []string{"test", "-c=circle", "-finder-frame", "rounded", "-m=diamond"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			workDir, err := os.MkdirTemp("", "climain-test-work-*")
			if err != nil {
				t.Fatalf("failed to create work dir: %v", err)
			}
			defer os.RemoveAll(workDir)

			exitCode, stdout, stderr := runCLI(t, workDir, tc.args...)

			if exitCode != 0 {
				t.Errorf("expected exit code 0, got %d. stderr: %q", exitCode, stderr)
			}

			if stderr != "" {
				t.Errorf("expected empty stderr, got: %q", stderr)
			}

			// Verify qr.svg was created
			qrPath := filepath.Join(workDir, "qr.svg")
			if _, err := os.Stat(qrPath); err != nil {
				t.Errorf("qr.svg should have been created, but it doesn't exist: %v", err)
			}

			if stdout != "" {
				t.Errorf("expected empty stdout, got: %q", stdout)
			}
		})
	}
}

func TestInvalidFlagValue(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{"invalid finder-center", []string{"test", "-finder-center", "triangle"}, "finder-center"},
		{"invalid finder-center short", []string{"test", "-c", "triangle"}, "finder-center"},
		{"invalid finder-frame", []string{"test", "-finder-frame", "triangle"}, "finder-frame"},
		{"invalid finder-frame short", []string{"test", "-f", "triangle"}, "finder-frame"},
		{"invalid module-shape", []string{"test", "-module-shape", "triangle"}, "module-shape"},
		{"invalid module-shape short", []string{"test", "-m", "triangle"}, "module-shape"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			workDir, err := os.MkdirTemp("", "climain-test-work-*")
			if err != nil {
				t.Fatalf("failed to create work dir: %v", err)
			}
			defer os.RemoveAll(workDir)

			exitCode, stdout, stderr := runCLI(t, workDir, tc.args...)

			if exitCode == 0 {
				t.Errorf("expected non-zero exit code, got %d", exitCode)
			}

			// Check that error mentions the flag or allowed values
			if !strings.Contains(stderr, tc.expected) && !strings.Contains(stderr, "circle") && !strings.Contains(stderr, "square") && !strings.Contains(stderr, "diamond") {
				t.Errorf("expected error message mentioning flag or allowed values, got stderr: %q", stderr)
			}

			// Verify qr.svg was not created
			qrPath := filepath.Join(workDir, "qr.svg")
			if _, err := os.Stat(qrPath); err == nil {
				t.Errorf("qr.svg should not have been created, but it exists")
			}

			if stdout != "" {
				t.Errorf("expected empty stdout, got: %q", stdout)
			}
		})
	}
}

func TestMissingFlagValue(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{"missing finder-center at end", []string{"test", "-finder-center"}, "finder-center"},
		{"missing finder-center short at end", []string{"test", "-c"}, "finder-center"},
		{"missing finder-center before flag", []string{"test", "-c", "-f", "rounded"}, "finder-center"},
		{"missing finder-frame at end", []string{"test", "-finder-frame"}, "finder-frame"},
		{"missing finder-frame short at end", []string{"test", "-f"}, "finder-frame"},
		{"missing finder-frame before flag", []string{"test", "-f", "-c", "circle"}, "finder-frame"},
		{"missing module-shape at end", []string{"test", "-module-shape"}, "module-shape"},
		{"missing module-shape short at end", []string{"test", "-m"}, "module-shape"},
		{"missing module-shape before flag", []string{"test", "-m", "-c", "circle"}, "module-shape"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			workDir, err := os.MkdirTemp("", "climain-test-work-*")
			if err != nil {
				t.Fatalf("failed to create work dir: %v", err)
			}
			defer os.RemoveAll(workDir)

			exitCode, stdout, stderr := runCLI(t, workDir, tc.args...)

			if exitCode == 0 {
				t.Errorf("expected non-zero exit code, got %d", exitCode)
			}

			// Check that error mentions the flag
			if !strings.Contains(stderr, tc.expected) {
				t.Errorf("expected error message mentioning %q, got stderr: %q", tc.expected, stderr)
			}

			// Verify qr.svg was not created
			qrPath := filepath.Join(workDir, "qr.svg")
			if _, err := os.Stat(qrPath); err == nil {
				t.Errorf("qr.svg should not have been created, but it exists")
			}

			if stdout != "" {
				t.Errorf("expected empty stdout, got: %q", stdout)
			}
		})
	}
}
