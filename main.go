package main

import (
	"flag"
	"fmt"
	"os"

	"symbol-art/internal/analyzer"
	"symbol-art/internal/banner"
	"symbol-art/internal/printer"
)

func main() {
	analyze := flag.Bool("analyze", false, "print an AI-style analysis of the input")
	suggest := flag.Bool("suggest", false, "print AI-style suggestions for the input")
	bannerName := flag.String("banner", "standard", "banner font to use: standard, shadow, or thinkertoy")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: go run . [--analyze] [--suggest] [--banner=standard|shadow|thinkertoy] STRING")
		os.Exit(1)
	}
	input := args[0]

	if !banner.Validate([]rune(input)) {
		fmt.Fprintln(os.Stderr, "Error: input contains unsupported characters")
		os.Exit(1)
	}

	b, err := banner.Load(*bannerName + ".txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	lines := printer.Render(input, b)
	printer.Print(lines)

	if *analyze || *suggest {
		height, width := printer.Dimensions(lines)
		if *analyze {
			fmt.Println()
			fmt.Println(analyzer.FormatAnalysis(input, height, width))
		}
		if *suggest {
			fmt.Println()
			fmt.Println(analyzer.FormatSuggestions(analyzer.Suggest(input, height, width)))
		}
	}
}
