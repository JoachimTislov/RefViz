package ops

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
	path, err := findContent(content)
	if err != nil {
		return fmt.Errorf("error finding content: %s, err: %v", *content, err)
	}
	e, valid := checkIfValid(path)
	if !valid {
		return fmt.Errorf("error: %s is not a valid entity, err: %v", path, err)
	}
	if e.IsDir() {
		if err := filepath.WalkDir(path, walk(findRefs)); err != nil {
			return fmt.Errorf("error walking through directory: %s, err: %v", path, err)
		}
	}
	if err := getContent(path, findRefs); err != nil {
		return fmt.Errorf("error getting content: %s, err: %v", path, err)
	}
	return nil
}

// walk walks through the directory and scans the files
func walk(findRefs *bool) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through directory: %s, err: %v", path, err)
		}
		if !d.IsDir() && isValid(d.IsDir(), path) {
			if err := getContent(path, findRefs); err != nil {
				return fmt.Errorf("error getting content: %s, err: %v", path, err)
			}
		}
		return nil
	}
}

// getContent walks for the content root and attempts to find the content
func findContent(content *string) (string, error) {
	projectRootPath := projectPath()
	if *content == "" {
		return projectRootPath, nil
	}
	var contentPath string
	err := filepath.WalkDir(projectRootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking from path: %s, err: %v", path, err)
		}
		_, valid := checkIfValid(path)
		if !valid {
			return nil
		}
		if filepath.Base(path) == *content {
			contentPath = path
			return fmt.Errorf(found)
		}
		return nil
	})
	if err != nil {
		if err.Error() == found {
			return contentPath, nil
		}
	}
	return "", fmt.Errorf("error finding content: %s, err: %v", *content, err)
}
