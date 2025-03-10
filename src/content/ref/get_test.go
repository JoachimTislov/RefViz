package ref_test

import (
	"testing"

	"github.com/JoachimTislov/RefViz/content/ref"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/types"
	"github.com/JoachimTislov/RefViz/test/helpers"
	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
	mainPath := path.QuickfeedRootMain()

	wantRefs := make(map[string]*types.Ref)
	key := mainPath + ":main"
	wantRefs[key] = helpers.CheckDomainRefInQFMain

	qfCheckDomain := helpers.QFCheckDomain
	qfCheckDomain.Refs = make(map[string]*types.Ref)

	if err := ref.Get(mainPath, qfCheckDomain)(); err != nil {
		t.Fatalf("ref Get() failed: %v", err)
	}

	if diff := cmp.Diff(qfCheckDomain.Refs, wantRefs); diff != "" {
		t.Errorf("ref Get() mismatch (-got +want):\n%s", diff)
	}
}
