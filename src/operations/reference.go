package op

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/Project-Visualizer/lsp"
	"github.com/JoachimTislov/Project-Visualizer/types"
)

func getRefs(filePath string, symbolPos string, refs []*types.Ref) error {
	pathToSymbol := fmt.Sprintf("%s:%s", filePath, symbolPos)
	output, err := lsp.RunGopls(references, pathToSymbol)
	if err != nil {
		return fmt.Errorf("error when running gopls command: %s, err: %s", references, err)
	}
	parseRefs(string(output), refs)
	return nil
}

func parseRefs(output string, refs []*types.Ref) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		path := strings.Split(line, ":")[0]
		refRelPath, err := filepath.Rel(projectPath(), path)
		if err != nil {
			fmt.Printf("error getting relative path: %s, err: %v", path, err)
		}
		fileName := filepath.Base(path)
		folderName := filepath.Base(filepath.Dir(path))
		ref := &types.Ref{Path: refRelPath, FolderName: folderName, FileName: fileName, MethodName: ""}
		refs = append(refs, ref)
	}
}
