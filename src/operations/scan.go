package op

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Scan scans the content for symbols and references
// If the content is a directory, it scans all files in the directory
// If scanForRefs is true, it scans for references
// If the content is a file, it only scans the file
func Scan(findRefs *bool, content *string) error {
	e, err := os.Stat(getRootPath(*content))
	if err != nil {
		return fmt.Errorf("error: %s is not a valid entity, err: %v", *content, err)
	}
	if e.IsDir() {
		if err := filepath.WalkDir(*content, walk(findRefs)); err != nil {
			return fmt.Errorf("error walking through directory: %s, err: %v", *content, err)
		}
	}
	if err := getContent(content, findRefs); err != nil {
		return fmt.Errorf("error getting content: %s, err: %v", *content, err)
	}
	return nil
}

// walk walks through the directory and scans the files
func walk(findRefs *bool) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through directory: %s, err: %v", path, err)
		}
		if err := isValid(d.IsDir(), &path); err != nil {
			return fmt.Errorf("error: %s is not a valid entity, err: %v", path, err)
		}
		if !d.IsDir() {
			if err := getContent(&path, findRefs); err != nil {
				return fmt.Errorf("error getting content: %s, err: %v", path, err)
			}
		}
		return nil
	}
}
