package main

import (
	"fmt"
	"os"

	"github.com/vtereso/paper_pete/pkg/parser"
	"github.com/vtereso/paper_pete/pkg/path"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not get CWD", err)
		os.Exit(1)
	}
	fmt.Println("Grabbing Go packages used within CWD:", cwd)

	// Extension point for other file format parsing
	// Each parser should accept a distinct file type where the first parser that is satisified will "consume" the file
	// where it will be not passed to the next parser.
	fileParserChain := []parser.FileParser{
		parser.GoParser{},
	}
	packages, err := path.WalkDirectory(cwd, fileParserChain)
	if err != nil {
		fmt.Println("Error walking CWD", err)
		os.Exit(1)
	}
	for _, p := range packages {
		fmt.Println(p)
	}
}
