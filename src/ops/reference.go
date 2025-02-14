package ops

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/Project-Visualizer/lsp"
	"github.com/JoachimTislov/Project-Visualizer/types"
)

func getRefs(pathToSymbol string, refs *map[string]*types.Ref) error {
	output, err := lsp.RunGopls(references, pathToSymbol)
	if err != nil {
		command := fmt.Sprintf("gopls references %s", pathToSymbol)
		return fmt.Errorf("error when running gopls command: %s, err: %s", command, err)
	}
	if err := parseRefs(string(output), refs); err != nil {
		return fmt.Errorf("error parsing references: %s, err: %v", pathToSymbol, err)
	}
	return nil
}

func parseRefs(output string, refs *map[string]*types.Ref) error {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		args := strings.Split(line, ":")
		path := args[0]
		LinePos := args[1]
		refRelPath, err := filepath.Rel(projectPath(), path)
		if err != nil {
			fmt.Printf("error getting relative path: %s, err: %v", path, err)
		}
		fileName := filepath.Base(path)
		folderName := filepath.Base(filepath.Dir(path))

		parentSymbol := types.Symbol{}
		if err := getRelatedMethod(path, LinePos, &parentSymbol); err != nil {
			return fmt.Errorf("error getting related method: %s, err: %v", path, err)
		}
		key := fmt.Sprintf("%s:%s", path, parentSymbol.Position.String())
		(*refs)[key] = &types.Ref{
			Path:         refRelPath,
			FolderName:   folderName,
			FileName:     fileName,
			ParentSymbol: parentSymbol,
		}
	}
	return nil
}

// getRelatedMethod finds the closest method above the reference
func getRelatedMethod(path string, refLinePos string, parentSymbol *types.Symbol) error {
	symbols, err := getSymbols(path)
	if err != nil {
		return fmt.Errorf("error getting symbols: %s, err: %v", path, err)
	}
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols found in %s", path)
	}
	parentSymbol.Position.Line = "0"
	// loop through potential parent symbols
	for _, s := range symbols {
		// skip if the symbol is not a function
		if s.Kind != function && s.Kind != method {
			continue
		}
		isFurtherDown := parentSymbol.Position.Line < s.Position.Line
		isAboveRef := s.Position.Line < refLinePos
		if isFurtherDown && isAboveRef {
			*parentSymbol = *s
		}
	}
	return nil
}
