package parser

import (
	"os"
	"strings"
)

// FileParser is an interface that parses a particular file format (e.g. `.go` files)
// and returns a slice of strings, which corresponds with the packages used within.
type FileParser interface {
	// Returns whether this file can be parsed
	Accepts(*os.File) bool
	// Reads a file and returns the import packages used within
	GetPackages(*os.File) ([]string, error)
}

// GoParser accepts files with a `.go` file format
type GoParser struct{}

// Accepts returns whether the file has a `.go` postfix/suffix
func (g GoParser) Accepts(file *os.File) bool {
	return strings.HasSuffix(file.Name(), ".go")
}

// GetPackages grabs all go import packages within `.go` files
// WIP
func (g GoParser) GetPackages(file *os.File) (packages []string, err error) {
	return
}
