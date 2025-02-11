package op

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/Project-Visualizer/types"
)

const (
	mapFolderPath     = "maps"
	rootFolderName    = "quickfeed"
	rootFolderRelPath = "../../../quickfeed"
)

// GetFile checks status of file and creates it if it doesn't exist, with empty content
// Reads the file and unmarshal it into the provided variable
func GetFile(filePath string, v any) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.WriteFile(filePath, []byte("{}"), 0o644)
	}
	if bytes, err := os.ReadFile(filePath); err != nil {
		return fmt.Errorf("get content from cache error: %s", err)
	} else {
		if err := json.Unmarshal(bytes, &v); err != nil {
			return fmt.Errorf("unmarshaling error: %s", err)
		}
	}
	return nil
}

func GetMapPath(name *string) string {
	return filepath.Join(mapFolderPath, fmt.Sprintf("%s.json", *name))
}

func getRelPath(filePath string) (string, error) {
	relPath, err := filepath.Rel(rootFolderRelPath, filePath)
	if err != nil {
		return "", fmt.Errorf("error getting relative path of file: %s, err: %v", filePath, err)
	}
	return relPath, nil
}

func getContent(content *string, scanForRefs *bool) error {
	absPath, err := filepath.Abs(*content)
	if err != nil {
		return fmt.Errorf("error getting absolute path of file: %s, err: %v", *content, err)
	}
	symbols, err := getSymbols(absPath)
	if err != nil {
		return fmt.Errorf("error getting symbols: %s, err: %v", *content, err)
	}
	if *scanForRefs {
		for _, s := range *symbols {
			if err := GetRefs(absPath, s.Position.GetPos(), s.Refs); err != nil {
				return fmt.Errorf("error getting references: %s, err: %v", *content, err)
			}
		}
	}
	if err := addSymbolsToFile(symbols, &absPath); err != nil {
		return fmt.Errorf("error adding symbols to file: %s, err: %v", *content, err)
	}
	return nil
}

// getRelatedMethod finds the closest method above the reference
func getRelatedMethod(symbols []types.Symbol, refParent *types.Symbol, refLinePos string) error {
	// loop through potential parent symbols
	for _, s := range symbols {
		// skip if the symbol is not a function
		if s.Kind != function && s.Kind != method {
			continue
		}
		isFurtherDown := refParent.Position.Line < s.Position.Line
		isAboveRef := s.Position.Line < refLinePos
		// if the new method is further down and above the reference, update the refParent
		if isFurtherDown && isAboveRef {
			*refParent = s
		}
	}
	return nil
}

// returns entry relative to last, of a string array with a given delimiter, i determines how many entries from the end
func getLastEntry(str string, delimiter string, i int) string {
	split := strings.Split(str, delimiter)
	return split[len(split)-1-i]
}
