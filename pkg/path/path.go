package path

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vtereso/paper_pete/pkg/parser"
)

// WalkDirectory walks the directory path specified and returns all of the packages used within. The packages are determined by the
// fileparsers.
func WalkDirectory(rootpath string, fileParserChain []parser.FileParser) (map[string]struct{}, error) {
	packages := map[string]struct{}{}
	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Could not open file", file)
				return err
			}
			for _, fileParser := range fileParserChain {
				if fileParser.Accepts(file) {
					filePackages, err := fileParser.GetPackages(file)
					if err != nil {
						fmt.Printf("Error getting packages for file %s: %s\n", file.Name(), err)
						return err
					}
					for _, p := range filePackages {
						packages[p] = struct{}{}
					}
				}
			}
		}
		return nil
	})
	// Returns named parameters
	return packages, err
}
