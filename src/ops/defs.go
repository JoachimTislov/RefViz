package ops

import (
	"github.com/JoachimTislov/Project-Visualizer/types"
)

const (
	samplePath     = "../../sample-code"
	src            = "src"
	references     = "references"
	symbols        = "symbols"
	method         = "Method"
	function       = "Function"
	tempFolder     = "/refViz"
	refVizRootPath = "refVizProjectRoot"
	scanAll        = "Scan all content"
)

var (
	config = types.NewConfig()
	cache  = types.NewCache()
)
