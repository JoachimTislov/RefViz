package cache

import (
	"fmt"
	"log"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/JoachimTislov/RefViz/internal"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/types"
)

var cacheMutex sync.Mutex

func lock() {
	cacheMutex.Lock()
}

func unlock() {
	cacheMutex.Unlock()
}

func Get() *types.Cache {
	lock()
	defer unlock()

	return cache
}

func Reset() {
	lock()
	defer unlock()

	cache = newCache()
}

func LogError(command string) {
	lock()
	defer unlock()

	if slices.Contains(cache.Errors, command) {
		return
	}
	cache.Errors = append(cache.Errors, command)
}

func addEntry(relPath string, entry *types.CacheEntry) {
	lock()
	defer unlock()

	cache.Entries[relPath] = *entry
}

func getEntry(relPath string) *types.CacheEntry {
	lock()
	defer unlock()

	if entry, ok := cache.Entries[relPath]; ok {
		return &entry
	}
	return NewCacheEntry("", 0, make(map[string]*types.Symbol))
}

func GetSymbol(relPath string, name string) *types.Symbol {
	lock()
	defer unlock()

	if entry, ok := cache.Entries[relPath]; ok {
		if symbol, ok := entry.Symbols[name]; ok {
			return symbol
		}
	}
	return nil
}

func GetEntries() map[string]types.CacheEntry {
	lock()
	defer unlock()

	return cache.Entries
}

func save() error {
	lock()
	defer unlock()

	return internal.MarshalAndWriteToFile(cache, path.Tmp(Name))
}

func AddUnusedSymbol(relPath string, name string, symbol types.OrphanSymbol) {
	lock()
	defer unlock()

	if cache.UnusedSymbols[relPath] == nil {
		cache.UnusedSymbols[relPath] = make(map[string]types.OrphanSymbol)
	}
	cache.UnusedSymbols[relPath][name] = symbol
}

func createGithubLink(symbol *types.Symbol, baseLink string) string {
	split2 := strings.Split(strings.Split(symbol.Path, "/quickfeed/")[1], ":")
	partialLink := split2[0] + "#L" + split2[1] + "-L" + strings.Split(split2[2], "-")[1]

	return fmt.Sprintf("[%s](%s%s), ", symbol.Name, baseLink, partialLink)
}

func UnusedSymbols() {
	for _, symbol := range cache.UnusedSymbols {
		for name, orphanSymbol := range symbol {
			if strings.Contains(orphanSymbol.Location, ".pb.go") || strings.Contains(orphanSymbol.Location, "kit") {
				continue
			}
			split2 := strings.Split(strings.Split(orphanSymbol.Location, "/quickfeed/")[1], ":")
			partialLink := split2[0] + "#L" + split2[1]
			fmt.Printf("[%s](%s%s), ", name, "https://github.com/quickfeed/quickfeed/tree/master/", partialLink)
		}
	}
}

func QuickfeedFindSemiOrphans() {
	log.Printf("\tFinding semi-orphans\n")

	for _, entry := range cache.Entries {
		for _, symbol := range entry.Symbols {
			basePath := filepath.Base(symbol.FilePath)

			qtest := filepath.Base(filepath.Dir(symbol.Path)) == "qtest"
			test_helper := basePath == "test_helper.go"
			github_mock_opts := basePath == "github_mock_opts.go"
			quickfeed_mock_opts := basePath == "quickfeed_mock_client.go"
			test := strings.HasSuffix(basePath, "_test.go")

			if !strings.Contains(symbol.Path, "kit") || qtest || len(symbol.Refs) == 0 || test_helper || test || github_mock_opts || quickfeed_mock_opts {
				continue
			}
			checkRefs(symbol)
		}
	}
}

func checkRefs(symbol *types.Symbol) {
	for _, ref := range symbol.Refs {
		if !strings.HasSuffix(ref.FileName, "_test.go") {
			return
		}
	}
	fmt.Printf(createGithubLink(symbol, "https://github.com/quickfeed/quickfeed/tree/master/"))
}
