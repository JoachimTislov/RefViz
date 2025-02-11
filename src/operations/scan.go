package op

import (
	"fmt"
	"os"
	"path/filepath"
)

const fileMapPath = "fileMap.json"

func Scan(isDir *bool, scanForRefs *bool, content *string) error {
	// render file map
	if err := GetFile(fileMapPath, &cache); err != nil {
		return fmt.Errorf("error getting symbols from cache: %s, err: %v", fileMapPath, err)
	}
	if *isDir {
		if err := filepath.WalkDir(*content, func(path string, d os.DirEntry, err error) error {
			if err := isValid(d.IsDir(), &path); err != nil {
				return fmt.Errorf("error: %s is not a valid entity, err: %v", path, err)
			}
			if !d.IsDir() {
				if err := getContent(&path, scanForRefs); err != nil {
					return fmt.Errorf("error getting content: %s, err: %v", path, err)
				}
			}
			return nil
		}); err != nil {
			return fmt.Errorf("error walking through directory: %s, err: %v", *content, err)
		}
	} else {
		if err := getContent(content, scanForRefs); err != nil {
			return fmt.Errorf("error getting content: %s, err: %v", *content, err)
		}
	}
	return nil
}
