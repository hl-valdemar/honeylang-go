package main

import (
	"fmt"
	"honey/lexer"
	"honey/source"
	"os"
)

func main() {
	source_paths := make([]string, 0, 10)
	for _, arg := range os.Args[1:] {
		source_paths = append(source_paths, arg)
	}

	if len(source_paths) == 0 {
		fmt.Fprintf(os.Stderr, "honey requires at least one source file to compile\n")
		os.Exit(1)
	}

	src, err := source.Load(source_paths[0])
	if err != nil {
		fmt.Println(err)
	}

	tokens, scanErrors := lexer.Scan(src)

	if scanErrors.Len() > 0 {
		for i := 0; i < scanErrors.Len(); i++ {
			line, col := src.LineCol(scanErrors.Starts[i])
			fmt.Fprintf(os.Stderr, "%s:%d:%d: error: %s\n",
				src.FullPath, line, col, scanErrors.Kinds[i].String())
		}
	}

	for i := 0; i < tokens.Len(); i++ {
		start := tokens.Starts[i]
		end := tokens.Ends[i]

		val := string(src.Contents[start:end])
		if val == "\n" {
			val = "\\n"
		}
		fmt.Printf("%-20s %s\n", tokens.Kinds[i].String(), val)
	}
}
