package content_test

import (
	"path/filepath"
	"testing"

	c "github.com/JoachimTislov/RefViz/content"
	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/internal/path"
)

func TestGet(t *testing.T) {
	negative := false
	content := filepath.Join(path.Quickfeed(), "main.go")

	if err := c.Get(content, negative, &negative)(); err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	cacheEntry, err := cache.Check(content)
	if err != nil {
		t.Fatalf("Check() failed: %v", err)
	}

	main := "main"
	init := "init"
	checkDomain := "checkDomain"

	wantSymbols := []string{main, init, checkDomain}
	wantAmount := len(wantSymbols)
	symbolsAmount := len(cacheEntry.Symbols)
	if symbolsAmount != wantAmount {
		t.Errorf("Got: %v, want: %v", symbolsAmount, wantAmount)
	}
	if cacheEntry.Name != "main.go" {
		t.Errorf("Got: %v, want: %v", cacheEntry.Name, "main.go")
	}
	for _, symbol := range wantSymbols {
		if s, ok := cacheEntry.Symbols[symbol]; !ok {
			t.Errorf("%v is missing", s)
		}
		childSymbols := cacheEntry.Symbols[symbol].ChildSymbols
		refs := cacheEntry.Symbols[symbol].Refs
		switch symbol {
		case main:
			if childSymbols == nil {
				t.Errorf("main is missing checkDomain as a child symbol")
			} else if len(childSymbols) != 1 {
				t.Errorf("main should only have one child symbol; %v", childSymbols)
			}
		case init:
			if refs != nil {
				t.Errorf("init should not have any references")
			}
		case checkDomain:
			if refs == nil {
				t.Errorf("checkDomain is missing references")
			}
		}
	}
}
