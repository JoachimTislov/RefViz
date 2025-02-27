package ops

import (
	"fmt"
	"path/filepath"

	"github.com/JoachimTislov/RefViz/internal"
	"github.com/JoachimTislov/RefViz/types"
)

func cacheEntry(cacheEntry *types.CacheEntry, path string) error {
	relPath, err := filepath.Rel(internal.ProjectPath(), path)
	if err != nil {
		return fmt.Errorf("error getting relative path: %s, err: %v", path, err)
	}
	cache.AddEntry(relPath, cacheEntry)
	// updates the cache file
	// writefile creates the cache file if it does not exist
	cache.Mu.Lock()
	defer cache.Mu.Unlock()
	if err := marshalAndWriteToFile(cache, internal.CachePath()); err != nil {
		return err
	}
	return nil
}
