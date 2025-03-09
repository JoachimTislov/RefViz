package ref

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/RefViz/internal/types"
)

func parseRefs(output string, refs *map[string]*types.Ref) error {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		args := strings.Split(line, ":")
		path := args[0]
		LinePos := args[1]

		fileName := filepath.Base(path)
		folderName := filepath.Base(filepath.Dir(path))

		parentSymbolName, err := getRelatedMethod(path, LinePos)
		if err != nil {
			return fmt.Errorf("error getting related method: %s, err: %v", path, err)
		}
		(*refs)[path] = &types.Ref{
			Path:       fmt.Sprintf("%s:%s:%s", path, args[1], args[2]),
			FilePath:   path,
			FolderName: folderName,
			FileName:   fileName,
			MethodName: *parentSymbolName,
		}
	}
	return nil
}
