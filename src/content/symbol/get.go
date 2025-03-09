package symbol

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/types"
	"github.com/JoachimTislov/RefViz/lsp"
)

const (
	symbols = "symbols"
)

func GetOne(path, symbolName string, forceScan *bool) (*types.Symbol, error) {
	entry, _, err := GetMany(path, *forceScan)
	if err != nil {
		return nil, fmt.Errorf("error getting symbols: %v", err)
	}
	s, ok := entry.Symbols[symbolName]
	if !ok {
		return nil, fmt.Errorf("symbol not found: %s", symbolName)
	}
	return s, nil
}

func GetMany(filePath string, scanAgain bool) (*types.CacheEntry, bool, error) {

	f, err := os.Stat(filePath)
	if err != nil {
		return nil, false, fmt.Errorf("error getting file info: %s, err: %v", filePath, err)
	}

	entry, err := cache.Check(filePath)
	if err != nil {
		return nil, false, fmt.Errorf("error checking cache: %s, err: %v", filePath, err)
	}

	modTime := f.ModTime().Unix()
	shouldScan := entry.ModTime != modTime || scanAgain
	if shouldScan {

		log.Printf("\tScanning for symbols for file: %s\n", filePath)

		output, err := lsp.RunGopls(path.Project(), symbols, filePath)
		if err != nil {
			cache.LogError(fmt.Sprintf("gopls %s %s", symbols, filePath))
			return nil, false, nil
		}

		parseSymbols(string(output), filePath, &entry.Symbols)

		entry.Name = filepath.Base(filePath)
		entry.ModTime = modTime

		log.Printf("\tCaching symbols for file: %s\n", filePath)

		if err := cache.CacheEntry(entry, filePath); err != nil {
			return entry, false, fmt.Errorf("error caching symbols: %s, err: %v", filePath, err)
		}
	}
	return entry, shouldScan, nil
}
