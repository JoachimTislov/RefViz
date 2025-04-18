package ref

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/internal/lsp"
	p "github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/types"
)

const (
	references = "references"
)

func Get(path string, symbol *types.Symbol) func() error {
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

		if err := parseRefs(string(output), symbol, relPath); err != nil {
			return fmt.Errorf("error parsing references: %s, err: %v", pathToSymbol, err)
		}

		return nil
	}
}
