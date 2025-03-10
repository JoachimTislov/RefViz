package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	refVizRootPath = "refVizProjectRoot"
	tempFolder     = "/refViz"
	// customPath is used to adjust the root path of the project
	// This only for development, TODO: remove later
	customPath = "/sample-code" // /sample-code/quickfeed
)

func Project() string {
	path := os.Getenv(refVizRootPath)
	if path == "" {
		if err := loadRoot(); err != nil {
			panic(fmt.Sprintf("error loading root path: %v", err))
		}
		path = os.Getenv(refVizRootPath)
	}
	if customPath != "" {
		path = filepath.Join(path, customPath)
	}
	return path
}

// loadRootPath loads the root path of the project and sets it
func loadRoot() error {
	root, err := getProjectRoot()
	if err != nil {
		return err
	}
	if err := setProjectRoot(root); err != nil {
		return err
	}
	return nil
}

func setProjectRoot(path string) error {
	return os.Setenv(refVizRootPath, path)
}

func GetMap(name string) string {
	if !strings.Contains(name, ".") {
		name = fmt.Sprintf("%s.json", name)
	}
	return filepath.Join(Map(), name)
}

func GetAbs(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path of file: %s, err: %v", path, err)
	}
	return absPath, nil
}

func getRoot(name string) string {
	return filepath.Join(Project(), name)
}

func Config() string {
	return getRoot(tmp("config.json"))
}

func Cache() string {
	return getRoot(tmp("cache.json"))
}

func GetTempFolder() string {
	return getRoot(tempFolder)
}

// tmp returns the path of the temporary folder
func tmp(name string) string {
	return filepath.Join(tempFolder, name)
}

func Map() string {
	return getRoot(tmp("maps"))
}

func Graphviz() string {
	return getRoot(tmp("graphviz"))
}

func DotFile(mapName *string) string {
	return filepath.Join(Graphviz(), fmt.Sprintf("%s.dot", *mapName))
}
