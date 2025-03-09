package cache

import (
	"slices"
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

func LogError(command string) {
	lock()
	defer unlock()

	if slices.Contains(cache.Errors, command) {
		return
	}
	cache.Errors = append(cache.Errors, command)
}

func AddEntry(relPath string, entry *types.CacheEntry) {
	lock()
	defer unlock()

	cache.Entries[relPath] = *entry
}

func GetEntry(relPath string) *types.CacheEntry {
	lock()
	defer unlock()

	if entry, ok := cache.Entries[relPath]; ok {
		return &entry
	}
	return newCacheEntry("", 0, make(map[string]*types.Symbol))
}

func save() error {
	lock()
	defer unlock()

	return internal.MarshalAndWriteToFile(cache, path.Cache())
}

func AddUnusedSymbol(relPath string, name string, symbol types.OrphanSymbol) {
	lock()
	defer unlock()

	if cache.UnusedSymbols[relPath] == nil {
		cache.UnusedSymbols[relPath] = make(map[string]types.OrphanSymbol)
	}
	cache.UnusedSymbols[relPath][name] = symbol
}
