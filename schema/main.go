package main

import (
	"6enten/garlicphone/schema/gowriter"
	"6enten/garlicphone/schema/parser"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	_ "embed"
)

//go:embed js.lua
var jsTemplater string

func main() {
	langFlag := flag.String("lang", "go", "The output language for the generated file (e.g., 'go', 'ts', 'js').")
	inputFlag := flag.String("i", "", "Path to the input schema file. (Required)")
	outputFlag := flag.String("o", "", "Path to the output directory. (Required)")
	namespaceFlag := flag.String("n", "", "Namespace for the generated file. (Optional: defaults to the input file name)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "This tool generates code from a schema file.")
		fmt.Fprintln(os.Stderr, "The 'flag' package automatically provides -h and --help flags.")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *inputFlag == "" {
		fmt.Fprintln(os.Stderr, "Error: Input file path is required. Use the -i flag.")
		flag.Usage()
		os.Exit(1)
	}

	if *outputFlag == "" {
		fmt.Fprintln(os.Stderr, "Error: Output directory path is required. Use the -o flag.")
		flag.Usage()
		os.Exit(1)
	}

	var namespace string
	if *namespaceFlag != "" {
		namespace = *namespaceFlag
	} else {
		base := filepath.Base(*inputFlag)
		ext := filepath.Ext(base)
		namespace = strings.TrimSuffix(base, ext)
	}

	schema, err := parser.GenerateSchema(*inputFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input file '%s': %v\n", *inputFlag, err)
		os.Exit(1)
	}

	var generatedCode string
	var fileExtension string

	switch *langFlag {
	case "go":
		generatedCode, err = gowriter.Print(schema, namespace)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating Go code: %v\n", err)
			os.Exit(1)
		}
		fileExtension = ".go"
	case "ts", "js":
		generatedCode, err = parser.RunLuaCodegen(schema, jsTemplater)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating JavaScript/TypeScript code: %v\n", err)
			os.Exit(1)
		}
		fileExtension = ".js"
	default:
		fmt.Fprintf(os.Stderr, "Error: Language '%s' is not supported.\n", *langFlag)
		os.Exit(1)
	}

	// 6. Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputFlag, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory '%s': %v\n", *outputFlag, err)
		os.Exit(1)
	}

	// 7. Construct the full output path and write the file
	outputFilename := namespace + fileExtension
	outputPath := filepath.Join(*outputFlag, outputFilename)

	err = os.WriteFile(outputPath, []byte(generatedCode), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to output file '%s': %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated file: %s\n", outputPath)
}
