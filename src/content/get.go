package content

import (
	"fmt"
	"log"

	"github.com/JoachimTislov/RefViz/content/ref"
	"github.com/JoachimTislov/RefViz/content/symbol"
	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/core/config"
	"github.com/JoachimTislov/RefViz/internal/routines"
	"github.com/JoachimTislov/RefViz/internal/types"
)

func Get(path string, scanAgain bool, everythingIsUpToDate *bool) func() error {
	return func() error {
		c, scannedForSymbols, err := symbol.GetMany(path, scanAgain)
		if err != nil {
			return fmt.Errorf("error getting symbols: %s, err: %v", path, err)
		}

		var scannedForRefs bool
		var jobs []func() error
		for _, s := range c.Symbols {
			if config.FindRefsForSymbols(s.Name) && len(s.Refs) == 0 && !s.ZeroRefs || scanAgain {
				scannedForRefs = true
				if s.Refs == nil {
					s.Refs = make(map[string]*types.Ref)
				}
				jobs = append(jobs, ref.Get(path, s))
			}
		}
		routines.StartWork(3, jobs)

		if scannedForRefs {
			if scannedForSymbols {
				log.Println("Found content for path: ", path)
				if everythingIsUpToDate != nil {
					*everythingIsUpToDate = false
				}
			} else {
				log.Println("No references to scan for in path: ", path)
			}

			log.Printf("Final caching for path: %s\n", path)

			if err := cache.CacheEntry(c, path); err != nil {
				return fmt.Errorf("error caching symbols: %s, err: %v", path, err)
			}
		}
		return nil
	}
}
