package path

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vtereso/paper_pete/pkg/parser"
)

func TestWalkDirectory(t *testing.T) {
	// This could be improved as well since it needs to be modified any time a package is added within this repository
	expectedPackages := map[string]struct{}{
		"testing":                      struct{}{},
		"github.com/google/go-cmp/cmp": struct{}{},
		"strings":                      struct{}{},
		"path/filepath":                struct{}{},
		"fmt":                          struct{}{},
		"os":                           struct{}{},
		"os/exec":                      struct{}{},
		"github.com/vtereso/paper_pete/pkg/parser": struct{}{},
		"github.com/vtereso/paper_pete/pkg/path":   struct{}{},
		"io/ioutil":                                struct{}{},
		"regexp":                                   struct{}{},
	}
	fileParserChain := []parser.FileParser{
		parser.GoParser{},
	}
	topLevelGitCommand := exec.Command("git", "rev-parse", "--show-toplevel")
	bytes, err := topLevelGitCommand.Output()
	if err != nil {
		t.Error(`Failed to get top level directory using 'git rev-parse --show-toplevel' command`)
	}
	packageTopLevelDirectory := strings.TrimSpace(string(bytes))
	actualPackages, err := WalkDirectory(packageTopLevelDirectory, fileParserChain)
	if err != nil {
		t.Fatal("Error walking CWD", err)
	}
	if diff := cmp.Diff(expectedPackages, actualPackages); diff != "" {
		t.Fatalf("WalkDirectory() Diff: -want +got: %s", diff)
	}
	// if !reflect.DeepEqual(expectedPackages, actualPackages) {
	// 	t.Fatalf("WalkDirectory() did not return expected packages \nExpected: %v\nActual: %v\n", expectedPackages, actualPackages)
	// }

}
