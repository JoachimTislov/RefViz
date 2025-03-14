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

// Project loads the project path, and will panic if is fails
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
		return fmt.Errorf("error getting project root, err: %v", err)
	}
	if err := os.Setenv(refVizRootPath, root); err != nil {
		return fmt.Errorf("error setting env %s, err: %v", refVizRootPath, err)
	}
	return nil
}

func GetMap(name string) string {
	if !strings.Contains(name, ".") {
		name = fmt.Sprintf("%s.json", name)
	}
	return filepath.Join(Map(), name)
}

func Tmp(name string) string {
	return getRoot(tmp(name))
}

func Map() string {
	return getRoot(tmp("maps"))
}

func DotFile(mapName *string) string {
	return filepath.Join(getRoot(tmp("graphviz")), fmt.Sprintf("%s.dot", *mapName))
}

func getRoot(name string) string {
	return filepath.Join(Project(), name)
}

// tmp returns the path of the temporary folder
func tmp(name string) string {
	return filepath.Join(tempFolder, name)
}
