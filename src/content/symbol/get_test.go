package symbol_test

import (
	"testing"

	"github.com/JoachimTislov/RefViz/content/symbol"
	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/types"
	"github.com/JoachimTislov/RefViz/test/helpers"
	"github.com/google/go-cmp/cmp"
)

func TestGetOne(t *testing.T) {
	main := path.QuickfeedRootMain()

	symbol, err := symbol.GetOne(main, "main", false)
	if err != nil {
		t.Fatalf("GetOne() failed: %v", err)
	}

	if diff := cmp.Diff(symbol, &types.Symbol{
		Name:     "main",
		Kind:     "Function",
		Path:     main + ":46:6-46:10",
		FilePath: main,
		Position: types.Position{Line: "46", CharRange: "6-10"},
	}); diff != "" {
		t.Errorf("GetOne() mismatch (-got +want):\n%s", diff)
	}
}

func TestGetMany(t *testing.T) {
	main := path.QuickfeedRootMain()
	want := cache.NewCacheEntry("main.go", 0, helpers.MainFileSymbols)

	cacheEntry, _, err := symbol.GetMany(main, false)
	if err != nil {
		t.Errorf("GetMany() failed: %v", err)
	}
	if cacheEntry.ModTime == 0 {
		t.Errorf("ModTime is 0, expected non-zero value; cacheEntry: %v", cacheEntry)
	}
	want.ModTime = cacheEntry.ModTime // Allows for successful diff comparison

	if diff := cmp.Diff(cacheEntry, want); diff != "" {
		t.Errorf("GetMany() mismatch (-got +want):\n%s", diff)
	}
}
