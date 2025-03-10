package content_test

import (
	"path/filepath"
	"testing"

	c "github.com/JoachimTislov/RefViz/content"
	"github.com/JoachimTislov/RefViz/internal/path"
)

func TestScan(t *testing.T) {
	content := filepath.Join(path.Quickfeed(), "main.go")

	if err := c.Scan(&content, false, false); err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}
}
