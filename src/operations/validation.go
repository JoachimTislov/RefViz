package op

import (
	"fmt"
	"path/filepath"
)

type sbMap map[string]bool

var (
	inExt   = sbMap{".go": true, ".ts": true, ".tsx": true} // supported extensions
	exDirs  = sbMap{"node_modules": true, "doc": true, ".git": true}
	exFiles = sbMap{}
)

// checks if the directory or file is valid
func isValid(isDir bool, content *string) error {
	if isDir {
		if exDirs[*content] {
			return fmt.Errorf("error: %s is an excluded directory", *content)
		}
	} else {
		if exFiles[*content] {
			return fmt.Errorf("error: %s is an excluded file", *content)
		}
		if !inExt[filepath.Ext(*content)] {
			return fmt.Errorf("error: File is not in a supported extension")
		}
		// bools are initialized to false, so no need to set it to false
	}
	return nil
}
