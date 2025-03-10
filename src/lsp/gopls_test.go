package lsp_test

import (
	"testing"

	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/lsp"
)

func TestRunGopls(t *testing.T) {
	project := path.Project()

	// Symbols
	// This is the expected output of the test
	output := "init Function 28:6-28:10\nmain Function 46:6-46:10\ncheckDomain Function 165:6-165:17\n"

	out, err := lsp.RunGopls(project, "symbols", "/home/yoa/RefViz/sample-code/quickfeed/main.go")
	if err != nil {
		t.Fatalf("RunGopls() failed: %v", err)
	}

	if string(out) != output {
		t.Errorf("Got: %v, want: %v", string(out), output)
	}

	// References

	output = "/home/yoa/RefViz/sample-code/quickfeed/main.go:88:51-62\n"
	out, err = lsp.RunGopls(project, "references", "/home/yoa/RefViz/sample-code/quickfeed/main.go:165:6-10")
	if err != nil {
		t.Fatalf("RunGopls() failed: %v", err)
	}

	if string(out) != output {
		t.Errorf("Got: %v, want: %v", string(out), output)
	}
}
