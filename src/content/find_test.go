package content_test

import (
	"path/filepath"
	"testing"

	c "github.com/JoachimTislov/RefViz/content"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/google/go-cmp/cmp"
)

func TestFind(t *testing.T) {

	// Test absolute path to main.go
	content := path.QuickfeedRootMain()
	want := []string{content}

	paths, err := c.Find(&content, false)
	if err != nil {
		t.Fatalf("Find() failed: %v", err)
	}
	if !cmp.Equal(paths, want) {
		t.Errorf("Got: %v, want: %v", paths, want)
	}

	// Test empty content
	content = ""
	want = []string{path.Project()}
	paths, err = c.Find(&content, false)
	if err != nil {
		t.Fatalf("Find() failed: %v", err)
	}
	if !cmp.Equal(paths, want) {
		t.Errorf("Got: %v, want: %v", paths, want)
	}

	// Test name of file
	content = "save.go"
	want = []string{filepath.Join(path.Quickfeed(), "internal", "env", "save.go")}
	paths, err = c.Find(&content, false)
	if err != nil {
		t.Fatalf("Find() failed: %v", err)
	}
	if !cmp.Equal(paths, want) {
		t.Errorf("Got: %v, want: %v", paths, want)
	}

	// Test name of directory
	content = "env"
	want = []string{filepath.Join(path.Quickfeed(), "internal", "env")}
	paths, err = c.Find(&content, false)
	if err != nil {
		t.Errorf("Got: %v, want: %v", paths, want)
	}
	if !cmp.Equal(paths, want) {
		t.Errorf("Got: %v, want: %v", paths, want)
	}
}
