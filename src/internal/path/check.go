package path

import (
	"fmt"
	"os"
)

// checkPath checks if the project path is valid
// If the path is valid, it returns the absolute path
func Check(projectPath string) (string, error) {
	if f, err := os.Stat(projectPath); err == nil && f.IsDir() {
		absPath, err := GetAbs(projectPath)
		if err != nil {
			return "", err
		}
		return absPath, nil
	} else if err == nil && !f.IsDir() {
		return "", fmt.Errorf("Project path: %s is not a directory\n", projectPath)
	}
	return "", fmt.Errorf("Project path: %s does not exist\n", projectPath)
}
