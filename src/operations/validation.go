package op

import (
	"os"
	"path/filepath"
)

// checks if the directory or file is valid
func isValid(isDir bool, content string) bool {
	exDirs, exFiles, inExt := getContentFilters()
	name := filepath.Base(content)
	if isDir {
		if exDirs[name] {
			return false
		}
	} else {
		if exFiles[name] {
			return false
		}
		e := filepath.Ext(name)
		if !inExt[e] {
			return false
		}
		// bools are initialized to false, so no need to set it to false
	}
	return true
}

// checks if the content is valid
func checkIfValid(content string) (os.FileInfo, bool) {
	c, err := os.Stat(content)
	if err != nil {
		return c, false
	}
	return c, isValid(c.IsDir(), content)
}
