package gouldian_test

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// every Go snippet of the user guide must be a copy of example_test.go,
// the examples are executed by go test.
func TestGuideSnippetsAreExamples(t *testing.T) {
	guide, err := os.ReadFile("doc/user-guide.md")
	if err != nil {
		t.Fatal(err)
	}
	examples, err := os.ReadFile("example_test.go")
	if err != nil {
		t.Fatal(err)
	}

	norm := func(s string) string {
		var out []string
		for _, line := range strings.Split(s, "\n") {
			if line = strings.TrimSpace(line); line != "" {
				out = append(out, line)
			}
		}
		return strings.Join(out, "\n")
	}

	source := norm(string(examples))
	blocks := regexp.MustCompile("(?s)```go\n(.*?)```").FindAllStringSubmatch(string(guide), -1)
	if len(blocks) == 0 {
		t.Fatal("the guide has no Go snippets")
	}

	for _, b := range blocks {
		if !strings.Contains(source, norm(b[1])) {
			t.Errorf("snippet is not a copy of example_test.go:\n%s", b[1])
		}
	}
}
