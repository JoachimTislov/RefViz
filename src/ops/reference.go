package ops

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/JoachimTislov/RefViz/lsp"
	"github.com/JoachimTislov/RefViz/types"
)

func getRefs(path string, symbol *types.Symbol, refs *map[string]*types.Ref, ch chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	pathToSymbol := fmt.Sprintf("%s:%s", path, symbol.Position.String())
	relPath, err := filepath.Rel(projectPath(), path)
	if err != nil {
		ch <- fmt.Errorf("error getting relative path: %s, err: %v", path, err)
	}

	log.Printf("\t\t Finding references for symbol: %s\n", symbol.Name)

	output, err := lsp.RunGopls(projectPath(), references, pathToSymbol)
	if err != nil {
		cache.LogError(fmt.Sprintf("gopls %s %s", references, pathToSymbol))
		symbol.ZeroRefs = true
		ch <- nil
	}
	// if there are no references, add the symbol to the unused symbols list
	if string(output) == "" {
		symbol.ZeroRefs = true
		// Add to unused map in the cache
		cache.AddUnusedSymbol(relPath, symbol.Name, types.NewUnusedSymbol(
			filepath.Base(filepath.Dir(path)),
			filepath.Base(path),
			pathToSymbol,
		))
	}

	if err := parseRefs(string(output), refs); err != nil {
		ch <- fmt.Errorf("error parsing references: %s, err: %v", pathToSymbol, err)
	}
}

func parseRefs(output string, refs *map[string]*types.Ref) error {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		args := strings.Split(line, ":")
		path := args[0]
		LinePos := args[1]
		refRelPath, err := filepath.Rel(projectPath(), path)
		if err != nil {
			fmt.Printf("error getting relative path: %s, err: %v", path, err)
		}
		fileName := filepath.Base(path)
		folderName := filepath.Base(filepath.Dir(path))

		parentSymbol := types.Symbol{}
		if err := getRelatedMethod(path, LinePos, &parentSymbol); err != nil {
			return fmt.Errorf("error getting related method: %s, err: %v", path, err)
		}
		key, err := filepath.Rel(projectPath(), path)
		if err != nil {
			return fmt.Errorf("error getting relative path: %s, err: %v", path, err)
		}
		(*refs)[key] = &types.Ref{
			Path:       refRelPath,
			FolderName: folderName,
			FileName:   fileName,
			MethodName: parentSymbol.Name,
		}
	}
	return nil
}

// getRelatedMethod finds the closest method above the reference
func getRelatedMethod(path string, refLinePos string, parentSymbol *types.Symbol) error {
	c, err := getSymbols(path, false)
	if err != nil {
		return fmt.Errorf("error getting symbols: %s, err: %v", path, err)
	}
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols found in %s", path)
	}
	parentSymbol.Position.Line = "0"
	// loop through potential parent symbols
	for _, s := range c.Symbols {
		// skip if the symbol is not a function
		if s.Kind != function && s.Kind != method {
			continue
		}
		isFurtherDown := parentSymbol.Position.Line < s.Position.Line
		isAboveRef := s.Position.Line < refLinePos
		if isFurtherDown && isAboveRef {
			*parentSymbol = *s
		}
	}
	return nil
}
