package ops

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/JoachimTislov/Project-Visualizer/types"
)

// GetFile reads the content of the file and unmarshals it into the given variable
func GetFile(filePath string, v any) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("get content from cache error: %s", err)
	}
	if err := json.Unmarshal(bytes, &v); err != nil {
		return fmt.Errorf("unmarshaling error: %s", err)
	}
	return nil
}

func getContent(path string, scanForRefs *bool) error {
	symbols, err := getSymbols(path)
	if err != nil {
		return fmt.Errorf("error getting symbols: %s, err: %v", path, err)
	}
	if *scanForRefs {
		for _, s := range symbols {
			if s.Refs == nil {
				s.Refs = make(map[string]*types.Ref)
			}
			if err := getRefs(fmt.Sprintf("%s:%s", path, s.Position.String()), &s.Refs); err != nil {
				return fmt.Errorf("error getting references: %s, err: %v", path, err)
			}
		}
	}
	if err := cacheSymbols(symbols, path); err != nil {
		return fmt.Errorf("error caching symbols: %s, err: %v", path, err)
	}
	return nil
}
