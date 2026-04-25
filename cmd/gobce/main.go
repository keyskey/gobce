package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/keyskey/gobce"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "analyze" {
		printUsage()
		os.Exit(2)
	}

	analyzeCmd := flag.NewFlagSet("analyze", flag.ExitOnError)
	coverProfile := analyzeCmd.String("coverprofile", "", "path to go coverprofile")
	outputFormat := analyzeCmd.String("format", "json", "output format (json)")
	outputPath := analyzeCmd.String("output", "", "optional path to write analysis result JSON")
	if err := analyzeCmd.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "parse flags: %v\n", err)
		os.Exit(2)
	}

	if *coverProfile == "" {
		fmt.Fprintln(os.Stderr, "--coverprofile is required")
		os.Exit(2)
	}
	if *outputFormat != "json" {
		fmt.Fprintf(os.Stderr, "unsupported format: %s\n", *outputFormat)
		os.Exit(2)
	}

	result, err := gobce.Analyze(gobce.AnalyzeInput{
		CoverProfilePath: *coverProfile,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "analyze failed: %v\n", err)
		os.Exit(1)
	}

	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode result: %v\n", err)
		os.Exit(1)
	}

	// Keep stdout output for terminal use and piping.
	if _, err := fmt.Fprintln(os.Stdout, string(encoded)); err != nil {
		fmt.Fprintf(os.Stderr, "write stdout: %v\n", err)
		os.Exit(1)
	}

	if *outputPath != "" {
		parent := filepath.Dir(*outputPath)
		if parent != "." {
			if err := os.MkdirAll(parent, 0o755); err != nil {
				fmt.Fprintf(os.Stderr, "create output directory: %v\n", err)
				os.Exit(1)
			}
		}
		if err := os.WriteFile(*outputPath, append(encoded, '\n'), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write output file: %v\n", err)
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  gobce analyze --coverprofile <path> --format json [--output <path>]")
}
