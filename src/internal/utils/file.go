package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateFile creates a file, if the folder path does not exist it will create it.
func CreateFile(path string) (*os.File, error) {
	dirPath := filepath.Dir(path)
	if !Exists(dirPath) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return nil, fmt.Errorf("error creating directory: %v", err)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("error creating map: %v", err)
	}
	return file, nil
}
