package ref_test

import (
	"testing"

	"github.com/JoachimTislov/refviz/content/ref"
	"github.com/JoachimTislov/refviz/internal/path"
	"github.com/JoachimTislov/refviz/internal/types"
	"github.com/JoachimTislov/refviz/test/helpers"
	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
	mainPath := path.QuickfeedRootMain()

	wantRefs := make(map[string]*types.Ref)
	key := mainPath + ":main"
	wantRefs[key] = helpers.CheckDomainRefInQFMain

	helpers.QFCheckDomain.Refs = make(map[string]*types.Ref)

	if err := ref.Get(mainPath, helpers.QFCheckDomain)(); err != nil {
		t.Fatalf("ref Get() failed: %v", err)
	}

	if diff := cmp.Diff(helpers.QFCheckDomain.Refs, wantRefs); diff != "" {
		t.Errorf("ref Get() mismatch (-got +want):\n%s", diff)
	}
}
