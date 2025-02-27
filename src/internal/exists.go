package internal

import "os"

// exists checks if a file exists
func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
