package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"gotime/internal/qr"
)

func main() {
	// Parse command line flags
	finderCenter := flag.String("finder-center", "square", "Finder center style: 'circle', 'square', or 'diamond'")
	finderFrame := flag.String("finder-frame", "square", "Finder frame style: 'square', 'rounded', 'circle', or 'diamond'")
	moduleShape := flag.String("module-shape", "rounded", "Module shape: 'square', 'rounded', 'circle', or 'diamond'")
	moduleSize := flag.Int("module-size", 0, "Module size in pixels (0 = default 10)")

	// Shorthand flags (aliases)
	flag.StringVar(finderCenter, "c", "square", "Shorthand for -finder-center")
	flag.StringVar(finderFrame, "f", "square", "Shorthand for -finder-frame")
	flag.StringVar(moduleShape, "m", "rounded", "Shorthand for -module-shape")
	flag.IntVar(moduleSize, "s", 0, "Shorthand for -module-size")

	// Parse flags (this stops at first positional arg)
	flag.Parse()

	// Manually check for flags that come after positional args
	// Go's flag package stops at first positional arg, so flags after are included in flag.Args()
	parsedArgs := flag.Args() // Args that flag.Parse() considered positional (includes flags after first pos arg)

	var qrContent string
	var positionalArgs []string

	// Parse remaining arguments to handle flags and positional args in any order
	skipNext := false
	for i := 0; i < len(parsedArgs); i++ {
		if skipNext {
			skipNext = false
			continue
		}

		arg := parsedArgs[i]

		// Check if this is a flag
		if len(arg) > 0 && arg[0] == '-' {
			// Handle -flag=value format
			if len(arg) >= 16 && arg[:16] == "-finder-center=" {
				value := arg[16:]
				if value == "" {
					fmt.Fprintf(os.Stderr, "Error: finder-center cannot be empty. Must be 'circle', 'square', or 'diamond'\n")
					os.Exit(1)
				}
				*finderCenter = value
				continue
			} else if len(arg) >= 3 && arg[:3] == "-c=" {
				value := arg[3:]
				if value == "" {
					fmt.Fprintf(os.Stderr, "Error: finder-center cannot be empty. Must be 'circle', 'square', or 'diamond'\n")
					os.Exit(1)
				}
				*finderCenter = value
				continue
			} else if len(arg) >= 14 && arg[:14] == "-finder-frame=" {
				value := arg[14:]
				if value == "" {
					fmt.Fprintf(os.Stderr, "Error: finder-frame cannot be empty. Must be 'square', 'rounded', 'circle', or 'diamond'\n")
					os.Exit(1)
				}
				*finderFrame = value
				continue
			} else if len(arg) >= 3 && arg[:3] == "-f=" {
				value := arg[3:]
				if value == "" {
					fmt.Fprintf(os.Stderr, "Error: finder-frame cannot be empty. Must be 'square', 'rounded', 'circle', or 'diamond'\n")
					os.Exit(1)
				}
				*finderFrame = value
				continue
			} else if len(arg) >= 14 && arg[:14] == "-module-shape=" {
				value := arg[14:]
				if value == "" {
					fmt.Fprintf(os.Stderr, "Error: module-shape cannot be empty. Must be 'square', 'rounded', 'circle', or 'diamond'\n")
					os.Exit(1)
				}
				*moduleShape = value
				continue
			} else if len(arg) >= 3 && arg[:3] == "-m=" {
				value := arg[3:]
				if value == "" {
					fmt.Fprintf(os.Stderr, "Error: module-shape cannot be empty. Must be 'square', 'rounded', 'circle', or 'diamond'\n")
					os.Exit(1)
				}
				*moduleShape = value
				continue
			} else if len(arg) >= 13 && arg[:13] == "-module-size=" {
				value := arg[13:]
				size, err := strconv.Atoi(value)
				if err != nil || size < 0 {
					fmt.Fprintf(os.Stderr, "Error: module-size must be a non-negative integer\n")
					os.Exit(1)
				}
				*moduleSize = size
				continue
			} else if len(arg) >= 3 && arg[:3] == "-s=" {
				value := arg[3:]
				size, err := strconv.Atoi(value)
				if err != nil || size < 0 {
					fmt.Fprintf(os.Stderr, "Error: module-size must be a non-negative integer\n")
					os.Exit(1)
				}
				*moduleSize = size
				continue
			}

			// Handle -flag value format
			if arg == "-finder-center" || arg == "-c" {
				if i+1 >= len(parsedArgs) {
					fmt.Fprintf(os.Stderr, "Error: finder-center requires a value. Must be 'circle', 'square', or 'diamond'\n")
					os.Exit(1)
				}
				value := parsedArgs[i+1]
				if value == "" || (len(value) > 0 && value[0] == '-') {
					fmt.Fprintf(os.Stderr, "Error: finder-center cannot be empty. Must be 'circle', 'square', or 'diamond'\n")
					os.Exit(1)
				}
				*finderCenter = value
				skipNext = true
				continue
			} else if arg == "-finder-frame" || arg == "-f" {
				if i+1 >= len(parsedArgs) {
					fmt.Fprintf(os.Stderr, "Error: finder-frame requires a value. Must be 'square', 'rounded', 'circle', or 'diamond'\n")
					os.Exit(1)
				}
				value := parsedArgs[i+1]
				if value == "" || (len(value) > 0 && value[0] == '-') {
					fmt.Fprintf(os.Stderr, "Error: finder-frame cannot be empty. Must be 'square', 'rounded', 'circle', or 'diamond'\n")
					os.Exit(1)
				}
				*finderFrame = value
				skipNext = true
				continue
			} else if arg == "-module-shape" || arg == "-m" {
				if i+1 >= len(parsedArgs) {
					fmt.Fprintf(os.Stderr, "Error: module-shape requires a value. Must be 'square', 'rounded', 'circle', or 'diamond'\n")
					os.Exit(1)
				}
				value := parsedArgs[i+1]
				if value == "" || (len(value) > 0 && value[0] == '-') {
					fmt.Fprintf(os.Stderr, "Error: module-shape cannot be empty. Must be 'square', 'rounded', 'circle', or 'diamond'\n")
					os.Exit(1)
				}
				*moduleShape = value
				skipNext = true
				continue
			} else if arg == "-module-size" || arg == "-s" {
				if i+1 >= len(parsedArgs) {
					fmt.Fprintf(os.Stderr, "Error: module-size requires a value (non-negative integer)\n")
					os.Exit(1)
				}
				value := parsedArgs[i+1]
				if value == "" || (len(value) > 0 && value[0] == '-') {
					fmt.Fprintf(os.Stderr, "Error: module-size must be a non-negative integer\n")
					os.Exit(1)
				}
				size, err := strconv.Atoi(value)
				if err != nil || size < 0 {
					fmt.Fprintf(os.Stderr, "Error: module-size must be a non-negative integer\n")
					os.Exit(1)
				}
				*moduleSize = size
				skipNext = true
				continue
			}
		}

		// This is a positional argument
		positionalArgs = append(positionalArgs, arg)
	}

	// Get QR code content from positional argument
	if len(positionalArgs) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <qr-content> [-finder-center=<style>] [-finder-frame=<style>] [-module-shape=<style>] [-module-size=<pixels>]\n", os.Args[0])
		os.Exit(1)
	}
	qrContent = positionalArgs[0]

	// Build options from flags
	opts := qr.Options{
		FinderCenter: *finderCenter,
		FinderFrame:  *finderFrame,
		ModuleShape:  *moduleShape,
		ModuleSize:   *moduleSize,
	}

	// Generate QR code SVG
	svgBytes, err := qr.GenerateSVG(qrContent, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Write to qr.svg
	if err := os.WriteFile("qr.svg", svgBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing qr.svg: %v\n", err)
		os.Exit(1)
	}
}
