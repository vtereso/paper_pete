package parser

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// goPackageRegex matches import blocks in the standard convention as follows:
// import(
// 	_ "a"
// 	. "b"
// 	"c"
// "d"
// "ee"
// )
var goPackageRegex = regexp.MustCompile(`[ \t]*import[ \t]*\(\n(?:[ \t]*\"([a-zA-Z0-9._]*)\"[ \t]*\n)*\)`)

// GoParser accepts files with a `.go` file format
type GoParser struct{}

// Accepts returns whether the file has a `.go` postfix/suffix
func (g GoParser) Accepts(file *os.File) bool {
	return strings.HasSuffix(file.Name(), ".go")
}

// GetPackages grabs all go import packages within `.go` files
func (g GoParser) GetPackages(file *os.File) (packages []string, err error) {
	b, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}
	matches := goPackageRegex.FindAllSubmatch(b, -1)
	for _, match := range matches {
		packageCaptureGroup := string(match[1])
		packages = append(packages, packageCaptureGroup)
	}
	return
}
