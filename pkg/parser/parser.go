package parser

import (
	"os"
)

// FileParser is an interface that parses a particular file format (e.g. `.go` files)
// and returns a slice of strings, which corresponds with the packages used within.
type FileParser interface {
	// Returns whether this file can be parsed
	Accepts(*os.File) bool
	// Reads a file and returns the import packages used within
	GetPackages(*os.File) ([]string, error)
}
