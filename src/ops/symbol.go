package ops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/Project-Visualizer/lsp"
	"github.com/JoachimTislov/Project-Visualizer/types"
)

func getSymbols(filePath string) (map[string]*types.Symbol, error) {
	cacheEntry, err := checkCache(filePath)
	if err != nil {
		return nil, err
	}
	f, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("error analyzing content: %s, err: %v", filePath, err)
	}
	if cacheEntry.ModTime != f.ModTime().Unix() {
		output, err := lsp.RunGopls(symbols, filePath)
		if err != nil {
			return nil, fmt.Errorf("error when running gopls command: %s, err: %s", symbols, err)
		}
		parseSymbols(string(output), &cacheEntry.Symbols)
	}
	return cacheEntry.Symbols, nil
}

func checkCache(filePath string) (types.CacheEntry, error) {
	var emptyCache types.CacheEntry
	if cache == nil {
		if err := loadCache(); err != nil {
			return emptyCache, err
		}
	}
	relPath, err := filepath.Rel(projectPath(), filePath)
	if err != nil {
		return emptyCache, fmt.Errorf("error getting relative path: %s, err: %v", filePath, err)
	}
	if entry, ok := (*cache)[relPath]; ok {
		return entry, nil
	}
	return types.NewCacheEntry("", 0, make(map[string]*types.Symbol)), nil
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
		fmt.Println(name, kind, createPosition(args[l-1]))
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
	fmt.Print(args[1])
	return types.Position{
		Line:      args2[0], // starting line position
		CharRange: fmt.Sprintf("%s-%s", args2[1], strings.Split(args[1], ":")[1]),
	}
}

func cacheSymbols(symbols map[string]*types.Symbol, path string) error {
	f, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error analyzing content: %s, err: %v", path, err)
	}
	relPath, err := filepath.Rel(projectPath(), path)
	if err != nil {
		return fmt.Errorf("error getting relative path: %s, err: %v", path, err)
	}
	if cache == nil {
		if err := loadCache(); err != nil {
			return err
		}
	}
	(*cache)[relPath] = types.NewCacheEntry(f.Name(), f.ModTime().Unix(), symbols)
	// updates the cache file
	// writefile creates the cache file if it does not exist
	if err := marshalAndWriteToFile(cache, cachePath()); err != nil {
		return err
	}
	return nil
}
