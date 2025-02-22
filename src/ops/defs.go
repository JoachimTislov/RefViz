package ops

import (
	"github.com/JoachimTislov/RefViz/types"
)

const (
	references     = "references"
	symbols        = "symbols"
	method         = "Method"
	function       = "Function"
	tempFolder     = "/refViz"
	refVizRootPath = "refVizProjectRoot"
	scanAll        = "Scan all content"
	exit           = "Cancel"
	scanSelected   = "Scan selected content"
	yes            = "y"
	// customPath is used to adjust the root path of the project
	// This only for development, TODO: remove later
	customPath = "" // /sample-code/quickfeed
)

var (
	config = types.NewConfig()
	cache  = types.NewCache()
)
