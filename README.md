My name is [`paper_pete`](https://www.youtube.com/watch?v=TwkKnv1H0fU), but you can call me PP!

Go has great first party tooling (the `go` CLI), however I normally only use a subset of the commands.
`go list` is one such neglected command and this repository is a simple recreation of some of its functionality.
With the sourcecode for Go being opensource, the "optimal" implementation would be to copy https://github.com/golang/go/blob/master/src/cmd/go/internal/list/list.go, but for the sake of doing it learning, I'll recreate it in a more simplisitic manner.
In order to list all the imported packages, this could be accomplished by the parsing the specified path and recursing for all directory elements.
In this initial implementation, I will implicitly assume `./...`, which is the go syntax for the current path and all nested directories.
This could also be extended to use the go syntax for pathing at a later point.
The go list command used the following struct for packages:
```go
type Package struct {
        Dir           string   // directory containing package sources
        ImportPath    string   // import path of package in dir
        ImportComment string   // path in import comment on package statement
        Name          string   // package name
        Doc           string   // package documentation string
        Target        string   // install path
        Shlib         string   // the shared library that contains this package (only set when -linkshared)
        Goroot        bool     // is this package in the Go root?
        Standard      bool     // is this package part of the standard Go library?
        Stale         bool     // would 'go install' do anything for this package?
        StaleReason   string   // explanation for Stale==true
        Root          string   // Go root or Go path dir containing this package
        ConflictDir   string   // this directory shadows Dir in $GOPATH
        BinaryOnly    bool     // binary-only package: cannot be recompiled from sources
        ForTest       string   // package is only for use in named test
        Export        string   // file containing export data (when using -export)
        Module        *Module  // info about package's containing module, if any (can be nil)
        Match         []string // command-line patterns matching this package
        DepOnly       bool     // package is only a dependency, not explicitly listed

        // Source files
        GoFiles         []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
        CgoFiles        []string // .go source files that import "C"
        CompiledGoFiles []string // .go files presented to compiler (when using -compiled)
        IgnoredGoFiles  []string // .go source files ignored due to build constraints
        CFiles          []string // .c source files
        CXXFiles        []string // .cc, .cxx and .cpp source files
        MFiles          []string // .m source files
        HFiles          []string // .h, .hh, .hpp and .hxx source files
        FFiles          []string // .f, .F, .for and .f90 Fortran source files
        SFiles          []string // .s source files
        SwigFiles       []string // .swig files
        SwigCXXFiles    []string // .swigcxx files
        SysoFiles       []string // .syso object files to add to archive
        TestGoFiles     []string // _test.go files in package
        XTestGoFiles    []string // _test.go files outside package

        // Cgo directives
        CgoCFLAGS    []string // cgo: flags for C compiler
        CgoCPPFLAGS  []string // cgo: flags for C preprocessor
        CgoCXXFLAGS  []string // cgo: flags for C++ compiler
        CgoFFLAGS    []string // cgo: flags for Fortran compiler
        CgoLDFLAGS   []string // cgo: flags for linker
        CgoPkgConfig []string // cgo: pkg-config names

        // Dependency information
        Imports      []string          // import paths used by this package
        ImportMap    map[string]string // map from source import to ImportPath (identity entries omitted)
        Deps         []string          // all (recursively) imported dependencies
        TestImports  []string          // imports from TestGoFiles
        XTestImports []string          // imports from XTestGoFiles

        // Error information
        Incomplete bool            // this package or a dependency has an error
        Error      *PackageError   // error loading package
        DepsErrors []*PackageError // errors loading dependencies
    }
```
In this way, the scope of this program could extend to files that do not have a `.go` file format, but I will constrain the function to just go files.
In terms of obtaining the packages, a top down parse of the file could be done in order to find the import packages.
Imports can be specified as follows:
```go
import "fmt"
import "math"
```
,but standard practice (as specified by [gotour](https://tour.golang.org/basics/2)) would be to use to use specify imports as follows
```go
import (
	"fmt"
	"math"
)
```
This is more technically outlined [here](https://golang.org/ref/spec#Import_declarations).
Since most any project would use the second syntax, I will initially target my implemention to handle this pattern.

A similar, but different approach would be to consult the Gopkg.lock/Gopkg.toml files for repositories with vendor directories (or the go.sum/go.mod files for repositories using go modules) in order to grab the repositories that are utilized by the project.
This would be signficantly faster to parse since ultimately it is the source of truth for the dependencies already, but it does provide the level of package granularity specified in a go import statement.

From a high level design, I will create FileParser interface that takes an [`os.File`](https://golang.org/src/os/types.go?s=369:411#L6) and returns the packages.
Using an interface for this _could_ be overkill (this is a toy project presenting pathways to things that likely will not get updated), but at the same point is also provides an extensibility point for the different file types as mentioned above/previous.
Also, go is somewhat unique with its package imports such that it would take more thought to make sense of the utility of this functionality across different programming languages, but I will create the extension point abstraction nonetheless.
For the sake of simplicity, the packages returned will be a slice of strings; for the sake of having a (starting point/place to interate), this seems sufficient.

The following command `go list -f '{{.Imports}}'` allows for a template to be passed, which allows the output of these packages to be customized; this too is an optimization/feature/improvement that I will forgo.
At a later point, this code could be updated to return a struct that would make sense with this sort of templating.
The output could also be stored per file, sorted, etc., which I will also forgo as it could be added at another point.

As for the specifics, multiple files may use the same packages and in order to prevent duplicates, a map makes most sense.
In this scenario, it is not a concern to pull members out, but rather check for the existance of packages already being used (hashset behavior).
For this reason, a further optimization would be to use a [bloom filter](https://llimllib.github.io/bloomfilter-tutorial/).
In terms of grabbing the packages per files, this could also be stored within a slice; depending on your use-case a hash can be considerred expensive (the same reason bloom filters could be preferable).

As another potential "optimization" would be to parse each file in a go routine.
Each routine could send its results over a channel (multiple writers, one reader) that would "merge" the result whether they be slices or maps. Since goroutines have _some_ overhead to get created, that may outweighs the cost to parse the files; I would also need to look into the parsing strategy to determine if there were ways to exit the files early (rather the read the entire thing).
Although in practice all imports are at the top, I don't know if it can be guaranteed this would be the case for arbitrary packages.
Benchmarks would be necessary to determine if this is truly an optimization; this implementation for `go list` doesn't do so, which presents further doubt.

I have tested this by using `go run main.go` within the project directory. In order to run this against different directories.
In order to test this against other repositories, run `go install`, get into the working directory/path that you want to execute this program against and then execute `paper_pete`.

From the time I started development, I recognized that testing was going to be very different from traditional unit testing because having to deal with files doesn't lend itself to using testing tables as traditionally. Since I am just about at time, I will make a test for `WalkDirectory` and ensure that the expected `map[string]struct{}` matches the actual values when run against this repository.

As a further improvement, I would further tune the regex as it does not support packages that are commented out within the import blob.