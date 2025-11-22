package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/phasecurve/zhuji/internal/compiler"
)

func main() {
	outputFile := flag.String("o", "", "output file (default: input.x86.s)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: zhuji [-o output] <input.s>")
		os.Exit(1)
	}

	inputFile := args[0]
	input, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", inputFile, err)
		os.Exit(1)
	}

	result := compiler.Compile(string(input))

	outPath := *outputFile
	if outPath == "" {
		outPath = strings.TrimSuffix(inputFile, ".s") + ".x86.s"
	}

	err = os.WriteFile(outPath, []byte(result), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing %s: %v\n", outPath, err)
		os.Exit(1)
	}

	fmt.Printf("wrote %s\n", outPath)
}
