package ref

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/JoachimTislov/RefViz/content/symbol"
	"github.com/JoachimTislov/RefViz/core/cache"
	p "github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/types"
	"github.com/JoachimTislov/RefViz/lsp"
)

const (
	references = "references"
)

func Get(path string, symbol *types.Symbol, refs *map[string]*types.Ref) func() error {
	return func() error {
		pathToSymbol := fmt.Sprintf("%s:%s", path, symbol.Position.String())
		relPath, err := filepath.Rel(p.Project(), path)
		if err != nil {
			return fmt.Errorf("error getting relative path: %s, err: %v", path, err)
		}

		log.Printf("\t\t Finding references for symbol: %s\n", symbol.Name)

		output, err := lsp.RunGopls(p.Project(), references, pathToSymbol)
		if err != nil {
			cache.LogError(fmt.Sprintf("gopls %s %s", references, pathToSymbol))
			symbol.ZeroRefs = true
			return nil
		}
		// if there are no references, add the symbol to the unused symbols list
		if string(output) == "" {
			symbol.ZeroRefs = true
			// Add to unused map in the cache
			cache.AddUnusedSymbol(relPath, symbol.Name, cache.NewOrphan(
				filepath.Base(filepath.Dir(path)),
				filepath.Base(path),
				pathToSymbol,
			))
		}

		if err := parseRefs(string(output), refs); err != nil {
			return fmt.Errorf("error parsing references: %s, err: %v", pathToSymbol, err)
		}

		return nil
	}
}

const (
	function = "Function"
	method   = "Method"
)

// getRelatedMethod finds the closest method above the reference
func getRelatedMethod(path string, refLinePos string) (*string, error) {
	c, _, err := symbol.GetMany(path, false)
	if err != nil {
		return nil, fmt.Errorf("error getting symbols: %s, err: %v", path, err)
	}
	if len(c.Symbols) == 0 {
		return nil, fmt.Errorf("zero symbols found in %s", path)
	}
	var parentSymbol *types.Symbol
	// loop through potential parent symbols
	for _, s := range c.Symbols {
		// skip if the symbol is not a function
		if s.Kind != function && s.Kind != method {
			continue
		}
		if parentSymbol == nil {
			parentSymbol = s
			continue
		}
		isFurtherDown := parentSymbol.Position.Line < s.Position.Line
		isAboveRef := s.Position.Line < refLinePos
		if isFurtherDown && isAboveRef {
			parentSymbol = s
		}
	}
	if parentSymbol == nil {
		for _, s := range c.Symbols {
			if s.Position.Line == refLinePos {
				return &s.Name, nil
			}
		}
		panic(fmt.Sprintf("Parent symbol is nil, path: %s, line: %s", path, refLinePos))
	}
	return &parentSymbol.Name, nil
}
