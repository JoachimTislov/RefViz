package cache

import (
	"fmt"
	"path/filepath"

	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/types"
)

var cache = NewCache()

func Get() *types.Cache {
	return cache
}

func ResetCache() {
	cache = NewCache()
}

func Check(filePath string) (*types.CacheEntry, error) {
	relPath, err := filepath.Rel(path.Project(), filePath)
	if err != nil {
		return nil, fmt.Errorf("error getting relative path: %s, err: %v", filePath, err)
	}
	return GetEntry(relPath), nil
}

// cacheEntry caches the cache entry
func CacheEntry(cacheEntry *types.CacheEntry, filePath string) error {
	relPath, err := filepath.Rel(path.Project(), filePath)
	if err != nil {
		return fmt.Errorf("error getting relative path: %s, err: %v", filePath, err)
	}
	addEntry(relPath, cacheEntry)
	if err := save(); err != nil {
		return err
	}
	return nil
}

func NewCache() *types.Cache {
	return &types.Cache{
		Errors:        []string{},
		UnusedSymbols: make(map[string]map[string]types.OrphanSymbol),
		Entries:       make(map[string]types.CacheEntry),
	}
}

func NewCacheEntry(name string, modTime int64, symbols map[string]*types.Symbol) *types.CacheEntry {
	return &types.CacheEntry{
		Name:    name,
		ModTime: modTime,
		Symbols: symbols,
	}
}

func NewOrphan(dir, fileName, location string) types.OrphanSymbol {
	return types.OrphanSymbol{
		Dir:      dir,
		FileName: fileName,
		Location: location,
	}
}
