package ops

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/JoachimTislov/RefViz/routines"
	"github.com/JoachimTislov/RefViz/types"
)

// GetFile reads the content of the file and unmarshals it into the given variable
func getFile(filePath string, v any) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("get content from cache error: %s", err)
	}
	if err := json.Unmarshal(bytes, &v); err != nil {
		return fmt.Errorf("unmarshaling error: %s", err)
	}
	return nil
}

func getContent(path string, scanAgain bool, everythingIsUpToDate *bool) func() error {
	return func() error {
		c, scannedForSymbols, err := getSymbols(path, scanAgain)
		if err != nil {
			return fmt.Errorf("error getting symbols: %s, err: %v", path, err)
		}

		var scannedForRefs bool
		var jobs []func() error
		for _, s := range c.Symbols {
			if !strings.HasPrefix(s.Name, "Test") && s.Name != "init" && s.Name != "main" && (len(s.Refs) == 0 && !s.ZeroRefs || scanAgain) {
				scannedForRefs = true
				if s.Refs == nil {
					s.Refs = make(map[string]*types.Ref)
				}
				jobs = append(jobs, getRefs(path, s, &s.Refs))
			}
		}
		var workers int
		l := len(jobs)
		if l < 3 {
			workers = l
		}
		routines.StartWork(workers, jobs)

		if scannedForSymbols {
			if scannedForRefs {
				log.Println("Found content for path: ", path)
				*everythingIsUpToDate = false
			} else {
				log.Println("No references to scan for in path: ", path)
			}
		}
		if scannedForRefs {
			log.Printf("Final caching for path: %s\n", path)

			if err := cacheEntry(c, path); err != nil {
				return fmt.Errorf("error caching symbols: %s, err: %v", path, err)
			}
		}
		return nil
	}
}
