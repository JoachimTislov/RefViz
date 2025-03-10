package ref

import (
	"path/filepath"
	"testing"

	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/types"
	"github.com/JoachimTislov/RefViz/test/helpers"
	"github.com/google/go-cmp/cmp"
)

func TestParseRefs(t *testing.T) {
	mainPath := path.QuickfeedRootMain()
	output := mainPath + ":88:51-62\n"
	helpers.QFCheckDomain.Refs = make(map[string]*types.Ref)
	want := map[string]*types.Ref{
		mainPath + ":main": helpers.CheckDomainRefInQFMain,
	}

	relPath, err := filepath.Rel(path.Project(), mainPath)
	if err != nil {
		t.Fatalf("filepath.Rel() failed: %v", err)
	}

	parseRefs(output, helpers.QFCheckDomain, relPath)

	if diff := cmp.Diff(helpers.QFCheckDomain.Refs, want); diff != "" {
		t.Errorf("parseRefs() mismatch (-got +want):\n%s", diff)
	}

	symbol := cache.GetSymbol(relPath, "main")

	if symbol == nil {
		t.Fatalf("cache.GetSymbol() failed: symbol not found")
	}

	childCheckDomain := &types.ChildSymbol{
		Key:      "quickfeed/main.go",
		FilePath: mainPath,
	}

	checkDomain := "checkDomain"
	child, ok := symbol.ChildSymbols[checkDomain]
	if !ok {
		t.Errorf("failed: child symbol: %s not found", checkDomain)
	}

	if diff := cmp.Diff(child, childCheckDomain); diff != "" {
		t.Errorf("cache.GetSymbol() mismatch (-got +want):\n%s", diff)
	}
}

func TestFindParent(t *testing.T) {
	mainPath := path.QuickfeedRootMain()

	symbol, foundParent, err := findParent(mainPath, "88") // 88 is the reference line number of checkDomain in main
	if err != nil {
		t.Fatalf("findParent() failed: %v", err)
	}

	if !foundParent {
		t.Errorf("findParent() failed: parent not found")
	}

	if diff := cmp.Diff(symbol, helpers.QFMain); diff != "" {
		t.Errorf("findParent() mismatch (-got +want):\n%s", diff)
	}
}
