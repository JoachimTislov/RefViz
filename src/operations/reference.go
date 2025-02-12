package op

import (
	"fmt"
	"strings"

	"github.com/JoachimTislov/Project-Visualizer/lsp"
	"github.com/JoachimTislov/Project-Visualizer/types"
)

func GetRefs(filePath string, symbolPos string, refs *[]*types.Ref) error {
	pathToSymbol := fmt.Sprintf("%s:%s", filePath, symbolPos)
	output, err := lsp.RunGopls(references, pathToSymbol)
	if err != nil {
		return fmt.Errorf("error when running gopls command: %s, err: %s", references, err)
	}
	parseRefs(string(output), refs)
	return nil
}

func parseRefs(output string, refs *[]*types.Ref) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		// TODO: Is there a better way ? What library method can be used to parse this?
		filePath := strings.Split(line, ":")[0]
		fileName := getLastEntry(filePath, "/", 0)
		folderName := getLastEntry(filePath, "/", 1)
		ref := &types.Ref{Path: line, FolderName: folderName, FileName: fileName, MethodName: ""}
		*refs = append(*refs, ref)
	}
}

// returns entry relative to last, of a string array with a given delimiter, i determines how many entries from the end
func getLastEntry(str string, delimiter string, i int) string {
	split := strings.Split(str, delimiter)
	return split[len(split)-1-i]
}
