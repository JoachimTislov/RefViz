package ops

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/RefViz/lsp"
	"github.com/JoachimTislov/RefViz/types"
)

func getSymbols(filePath string, scanAgain bool) (*types.CacheEntry, error) {

	f, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %s, err: %v", filePath, err)
	}

	entry, err := checkCache(filePath)
	if err != nil {
		return nil, fmt.Errorf("error checking cache: %s, err: %v", filePath, err)
	}

	modTime := f.ModTime().Unix()
	if entry.ModTime != modTime || scanAgain {

		log.Printf("\tScanning for symbols for file: %s\n", filePath)

		output, err := lsp.RunGopls(projectPath(), symbols, filePath)
		if err != nil {
			return nil, fmt.Errorf("error when running gopls command: %s, err: %s", symbols, err)
		}
		parseSymbols(string(output), &entry.Symbols)

		entry.Name = filepath.Base(filePath)
		entry.ModTime = modTime

		log.Printf("\tCaching symbols for file: %s\n", filePath)

		if err := cacheEntry(entry, filePath); err != nil {
			return entry, fmt.Errorf("error caching symbols: %s, err: %v", filePath, err)
		}
	}
	return entry, nil
}

func checkCache(filePath string) (*types.CacheEntry, error) {
	relPath, err := filepath.Rel(projectPath(), filePath)
	if err != nil {
		return nil, fmt.Errorf("error getting relative path: %s, err: %v", filePath, err)
	}
	return cache.GetEntry(relPath), nil
}

// parses the output of the gopls symbols command and extracts the name, kind, and position of each symbol
func parseSymbols(output string, s *map[string]*types.Symbol) {
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

func cacheEntry(cacheEntry *types.CacheEntry, path string) error {
	relPath, err := filepath.Rel(projectPath(), path)
	if err != nil {
		return fmt.Errorf("error getting relative path: %s, err: %v", path, err)
	}
	cache.AddEntry(relPath, cacheEntry)
	// updates the cache file
	// writefile creates the cache file if it does not exist
	if err := marshalAndWriteToFile(cache, cachePath()); err != nil {
		return err
	}
	return nil
}
