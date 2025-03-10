package helpers

import (
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/types"
)

// Taken Quickfeed's main.go

var mainPath = path.Quickfeed() + "/main.go"

var MainFileSymbols = map[string]*types.Symbol{
	"checkDomain": QFCheckDomain,
	"init": {
		Name: "init",
		Kind: "Function",
		Position: types.Position{
			Line:      "28",
			CharRange: "6-10",
		},
		Path:     mainPath + ":28:6-28:10",
		FilePath: mainPath,
	},
	"main": QFMain,
}

var QFCheckDomain = &types.Symbol{
	Name: "checkDomain",
	Kind: "Function",
	Position: types.Position{
		Line:      "165",
		CharRange: "6-17",
	},
	Path:     mainPath + ":165:6-165:17",
	FilePath: mainPath,
}

var CheckDomainRefInQFMain = &types.Ref{
	FilePath:   mainPath,
	Path:       mainPath + ":88:51-62",
	FolderName: "quickfeed",
	FileName:   "main.go",
	MethodName: "main",
}

var QFMain = &types.Symbol{
	Name: "main",
	Kind: "Function",
	Position: types.Position{
		Line:      "46",
		CharRange: "6-10",
	},
	Path:     mainPath + ":46:6-46:10",
	FilePath: mainPath,
}
