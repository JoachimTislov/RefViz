package op

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/Project-Visualizer/lsp"
	"github.com/JoachimTislov/Project-Visualizer/types"
)

func getSymbols(filePath string) ([]*types.Symbol, error) {
	var s []*types.Symbol
	output, err := lsp.RunGopls(symbols, filePath)
	if err != nil {
		return nil, fmt.Errorf("error when running gopls command: %s, err: %s", symbols, err)
	}
	extractSymbols(string(output), &s)
	return s, nil
}

// parses the output of the gopls symbols command and extracts the name, kind, and position of each symbol
func extractSymbols(output string, s *[]*types.Symbol) {
	for _, line := range strings.Split(output, "\n") {
		args := strings.Split(line, " ")
		if len(args) < 3 {
			continue
		}
		name := args[0]
		kind := args[1]
		// for methods, remove the receiver type
		if kind == method && strings.Contains(name, ".") {
			name = strings.Split(name, ".")[1]
		}
		*s = append(*s, &types.Symbol{
			Name:     name,
			Kind:     kind,
			Position: createPosition(args[2]),
		})
	}
}

// Gets the line and character range position of the symbol
func createPosition(p string) types.Position {
	args := strings.Split(p, "-")
	args2 := strings.Split(args[0], ":")
	return types.Position{
		Line:      args2[0], // starting line position
		CharRange: fmt.Sprintf("%s-%s", args2[1], strings.Split(args[1], ":")[1]),
	}
}

func cacheSymbols(symbols []*types.Symbol, path string) error {
	f, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error analyzing content: %s, err: %v", path, err)
	}
	name := f.Name()
	relPath, err := filepath.Rel(projectPath(), path)
	if err != nil {
		return fmt.Errorf("error getting relative path: %s, err: %v", path, err)
	}
	if cache == nil {
		if err := loadCache(); err != nil {
			return fmt.Errorf("error loading cache: %v", err)
		}
	}
	(*cache)[name] = *types.NewCacheEntry(relPath, f.ModTime().Unix(), symbols)
	// updates the cache file
	// writefile creates the cache file if it does not exist
	if err := marshalAndWriteToFile(cache, cachePath()); err != nil {
		return err
	}
	return nil
}
