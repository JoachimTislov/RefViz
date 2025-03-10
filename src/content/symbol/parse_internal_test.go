package symbol

import (
	"testing"

	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/types"
	"github.com/JoachimTislov/RefViz/test/helpers"
	"github.com/google/go-cmp/cmp"
)

func TestParseSymbols(t *testing.T) {
	// Taken from quickfeed/main.go, expected output of gopls symbols command: "gopls symbols *path/to/main.go*"
	output := "init Function 28:6-28:10\nmain Function 46:6-46:10\ncheckDomain Function 165:6-165:17\n"
	s := make(map[string]*types.Symbol)

	quickfeed := path.Quickfeed()
	main := quickfeed + "/main.go"

	parseSymbols(output, main, &s)

	if diff := cmp.Diff(s, helpers.MainFileSymbols); diff != "" {
		t.Errorf("parseSymbols() mismatch (-got +want):\n%s", diff)
	}
}
