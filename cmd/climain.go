package main

import (
	"flag"
	"fmt"
	"os"

	"gotime/internal/qr"
)

func main() {
	// Parse command line flags
	cornerCenter := flag.String("finder-center", "square", "Corner center style: 'circle', 'square', or 'diamond'")
	finderFrame := flag.String("finder-frame", "square", "Finder frame style: 'square', 'rounded', 'circle', or 'diamond'")
	moduleShape := flag.String("module-shape", "rounded", "Module shape: 'square', 'rounded', 'circle', or 'diamond'")

	// Shorthand flags (aliases)
	flag.StringVar(cornerCenter, "c", "square", "Shorthand for -finder-center")
	flag.StringVar(finderFrame, "f", "square", "Shorthand for -finder-frame")
	flag.StringVar(moduleShape, "m", "rounded", "Shorthand for -module-shape")

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
				*cornerCenter = value
				continue
			} else if len(arg) >= 3 && arg[:3] == "-c=" {
				value := arg[3:]
				if value == "" {
					fmt.Fprintf(os.Stderr, "Error: finder-center cannot be empty. Must be 'circle', 'square', or 'diamond'\n")
					os.Exit(1)
				}
				*cornerCenter = value
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
				*cornerCenter = value
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
			}
		}

		// This is a positional argument
		positionalArgs = append(positionalArgs, arg)
	}

	// Get QR code content from positional argument
	if len(positionalArgs) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <qr-content> [-finder-center=<style>] [-finder-frame=<style>] [-module-shape=<style>]\n", os.Args[0])
		os.Exit(1)
	}
	qrContent = positionalArgs[0]

	// Build options from flags
	opts := qr.Options{
		FinderCenter: *cornerCenter,
		FinderFrame:  *finderFrame,
		ModuleShape:  *moduleShape,
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
