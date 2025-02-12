package op

import (
	"fmt"
	"os"
	"path/filepath"
)

// checks if the directory or file is valid
func isValid(isDir bool, content *string) error {
	exDirs, exFiles, inExt := getContentFilters()
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

// checks if the content is valid
func checkIfValid(content *string) bool {
	if c, err := os.Stat(*content); err != nil {
		return false
	} else {
		if err := isValid(c.IsDir(), content); err != nil {
			return false
		}
	}
	return true
}
