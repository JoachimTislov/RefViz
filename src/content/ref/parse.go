package ref

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/RefViz/content/symbol"
	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/internal/types"
)

func parseRefs(output string, childSymbol *types.Symbol, relPath string) error {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		args := strings.Split(line, ":")
		path := args[0]
		LinePos := args[1]

		fileName := filepath.Base(path)
		folderName := filepath.Base(filepath.Dir(path))

		parentSymbol, hasParent, err := findParent(path, LinePos)
		if !hasParent {
			// Skip if no parent symbol is found
			cache.LogError(err.Error())
			continue
		}
		if err != nil {
			return fmt.Errorf("error getting related method: %s, err: %v", path, err)
		}
		key := fmt.Sprintf("%s:%s", path, parentSymbol.Name)
		childSymbol.Refs[key] = &types.Ref{
			Path:       line,
			FilePath:   path,
			FolderName: folderName,
			FileName:   fileName,
			MethodName: parentSymbol.Name,
		}
		parentSymbol.AddChildSymbol(childSymbol.Name, childSymbol.FilePath, relPath)
	}
	return nil
}

const (
	function = "Function"
	method   = "Method"
)

// findParent finds the closest method above the reference
func findParent(path string, refLinePos string) (*types.Symbol, bool, error) {
	c, _, err := symbol.GetMany(path, false)
	if err != nil {
		return nil, false, fmt.Errorf("error getting symbols: %s, err: %v", path, err)
	}
	if len(c.Symbols) == 0 {
		return nil, false, fmt.Errorf("zero symbols found in %s", path)
	}
	var parentSymbol *types.Symbol
	// loop through potential parent symbols
	for _, s := range c.Symbols {

		// TODO: filter out wanted parent symbols

		// Initialize parentSymbol with the first symbol
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
		return nil, false, fmt.Errorf("no parent symbol found for %s", path)
	}
	return parentSymbol, true, nil
}
