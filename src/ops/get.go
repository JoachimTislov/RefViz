package ops

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/JoachimTislov/RefViz/generics"
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

func getContent(path string, scanAgain bool) error {
	log.Println("Getting content for path: ", path)

	c, err := getSymbols(path, scanAgain)
	if err != nil {
		return fmt.Errorf("error getting symbols: %s, err: %v", path, err)
	}

	var scannedForRefs bool
	var jobs []func() error
	for _, s := range c.Symbols {
		if !strings.HasPrefix(s.Name, "Test") && s.Name != "init" && len(s.Refs) == 0 && !s.ZeroRefs || scanAgain {
			scannedForRefs = true
			if s.Refs == nil {
				s.Refs = make(map[string]*types.Ref)
			}
			jobs = append(jobs, generics.JobThreeArgs(getRefs, path, s, &s.Refs))
		}
	}
	routines.StartWork(5, jobs)

	if scannedForRefs {
		log.Printf("Final caching for path: %s\n", path)

		if err := cacheEntry(c, path); err != nil {
			return fmt.Errorf("error caching symbols: %s, err: %v", path, err)
		}
	} else {
		log.Println("No references to scan for in path: ", path)
	}
	return nil
}
