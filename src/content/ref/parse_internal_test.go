package ref

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/types"
	"github.com/JoachimTislov/RefViz/test/helpers"
	"github.com/google/go-cmp/cmp"
)

func TestFindParent(t *testing.T) {
	mainPath := path.QuickfeedRootMain()

	tests := []struct {
		name         string
		referencePos string
		want         *types.Symbol
	}{
		{
			name:         "checkDomain",
			referencePos: "88",
			want:         helpers.QFMain,
		},
		{
			name:         "NewQuickFeedService",
			referencePos: "132",
			want:         helpers.QFMain,
		},
		{
			name:         "Serve",
			referencePos: "159",
			want:         helpers.QFMain,
		},
		{
			name:         "Domain in Main",
			referencePos: "75",
			want:         helpers.QFMain,
		},
		{
			name:         "Domain in CheckDomain",
			referencePos: "166",
			want:         helpers.QFCheckDomain,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			symbol, foundParent, err := findParent(mainPath, tt.referencePos)
			if err != nil {
				t.Fatalf("findParent() failed: %v", err)
			}
			if !foundParent {
				t.Errorf("findParent() failed: parent not found")
			}
			if diff := cmp.Diff(symbol, tt.want); diff != "" {
				t.Errorf("findParent() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParseRefs(t *testing.T) {
	mainPath := path.QuickfeedRootMain()
	output := fmt.Sprintf("%s:%s", mainPath, "88:51-62\n")
	helpers.QFCheckDomain.Refs = make(map[string]*types.Ref)
	want := map[string]*types.Ref{
		mainPath + ":main": helpers.CheckDomainRefInQFMain,
	}

	relPath, err := filepath.Rel(path.Project(), mainPath)
	if err != nil {
		t.Fatalf("filepath.Rel() failed: %v", err)
	}

	if err := parseRefs(output, helpers.QFCheckDomain, relPath); err != nil {
		t.Fatalf("parseRefs() failed: %v", err)
	}

	if diff := cmp.Diff(helpers.QFCheckDomain.Refs, want); diff != "" {
		t.Errorf("parseRefs() mismatch (-got +want):\n%s", diff)
	}

	symbol := cache.GetSymbol(relPath, "main")
	if symbol == nil {
		t.Fatalf("cache.GetSymbol() failed: symbol not found")
	}

	childCheckDomain := &types.ChildSymbol{
		Key:      filepath.Join("quickfeed", "main.go"),
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
