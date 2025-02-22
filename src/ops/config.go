package ops

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JoachimTislov/RefViz/types"
)

func save() error {
	if err := marshalAndWriteToFile(config, configPath()); err != nil {
		return fmt.Errorf("error updating configurations: %v", err)
	}
	return nil
}

// checkPath checks if the project path is valid
// If the path is valid, it returns the absolute path
func checkPath(projectPath string) (string, error) {
	if f, err := os.Stat(projectPath); err == nil && f.IsDir() {
		absPath, err := getAbsPath(projectPath)
		if err != nil {
			return "", err
		}
		return absPath, nil
	} else if err == nil && !f.IsDir() {
		return "", fmt.Errorf("Project path: %s is not a directory\n", projectPath)
	}
	return "", fmt.Errorf("Project path: %s does not exist\n", projectPath)
}

func addExDirs(dirs ...string) error {
	return exclude(&config.ExDirs, dirs...)
}

func addExFiles(files ...string) error {
	return exclude(&config.ExDirs, files...)
}

func exclude(m *types.SbMap, items ...string) error {
	for i := range items {
		_, valid := checkIfValid(items[i])
		if valid {
			(*m)[items[i]] = true
		} else {
			return fmt.Errorf("invalid exclusion: %s", items[i])
		}
	}
	return save()
}

func projectPath() string {
	path := os.Getenv(refVizRootPath)
	if path == "" {
		panic("project path not set")
	}
	if customPath != "" {
		path = filepath.Join(path, customPath)
	}
	return path
}

func setProjectPath(path string) error {
	return os.Setenv(refVizRootPath, path)
}

func getContentFilters() (types.SbMap, types.SbMap, types.SbMap) {
	return config.ExDirs, config.ExFiles, config.InExt
}
