package ops

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/RefViz/internal"
	"github.com/JoachimTislov/RefViz/lsp"
	"github.com/JoachimTislov/RefViz/types"
)

const (
	symbols = "symbols"
)

func getSymbol(path, symbolName string, forceScan *bool) (*types.Symbol, error) {
	entry, _, err := getSymbols(path, *forceScan)
	if err != nil {
		return nil, fmt.Errorf("error getting symbols: %v", err)
	}
	s, ok := entry.Symbols[symbolName]
	if !ok {
		return nil, fmt.Errorf("symbol not found: %s", symbolName)
	}
	return s, nil
}

func getSymbols(filePath string, scanAgain bool) (*types.CacheEntry, bool, error) {

	f, err := os.Stat(filePath)
	if err != nil {
		return nil, false, fmt.Errorf("error getting file info: %s, err: %v", filePath, err)
	}

	entry, err := checkCache(filePath)
	if err != nil {
		return nil, false, fmt.Errorf("error checking cache: %s, err: %v", filePath, err)
	}

	modTime := f.ModTime().Unix()
	shouldScan := entry.ModTime != modTime || scanAgain
	if shouldScan {

		log.Printf("\tScanning for symbols for file: %s\n", filePath)

		output, err := lsp.RunGopls(internal.ProjectPath(), symbols, filePath)
		if err != nil {
			cache.LogError(fmt.Sprintf("gopls %s %s", symbols, filePath))
			return nil, false, nil
		}

		parseSymbols(string(output), filePath, &entry.Symbols)

		entry.Name = filepath.Base(filePath)
		entry.ModTime = modTime

		log.Printf("\tCaching symbols for file: %s\n", filePath)

		if err := cacheEntry(entry, filePath); err != nil {
			return entry, false, fmt.Errorf("error caching symbols: %s, err: %v", filePath, err)
		}
	}
	return entry, shouldScan, nil
}

func checkCache(filePath string) (*types.CacheEntry, error) {
	relPath, err := filepath.Rel(internal.ProjectPath(), filePath)
	if err != nil {
		return nil, fmt.Errorf("error getting relative path: %s, err: %v", filePath, err)
	}
	return cache.GetEntry(relPath), nil
}

// parses the output of the gopls symbols command and extracts the name, kind, and position of each symbol
func parseSymbols(output, filePath string, s *map[string]*types.Symbol) {
	// retrieve values from last, to handle this specific case: uint64 | string Field 27:2-27:17
	// the usual case is: uint64 Field 27:2-27:17, name, kind, position
	for _, line := range strings.Split(output, "\n") {
		args := strings.Split(line, " ")
		l := len(args)
		if l < 3 {
			continue
		}
		name := strings.TrimSpace(strings.Join(args[:l-2], " ")) // name is everything except the last 2 elements
		kind := args[l-2]
		// for methods, remove the receiver type
		// (*Service[K, V]).SendTo Method -> SendTo
		if kind == method && strings.Contains(name, ".") {
			name = strings.Split(name, ".")[1]
		}
		(*s)[name] = &types.Symbol{
			Name:     name,
			Kind:     kind,
			Path:     fmt.Sprintf("%s:%s", filePath, args[l-1]),
			FilePath: filePath,
			Position: createPosition(args[l-1]),
		}
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
