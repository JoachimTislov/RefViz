package internal

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	refVizRootPath = "refVizProjectRoot"
	tempFolder     = "/refViz"
	// customPath is used to adjust the root path of the project
	// This only for development, TODO: remove later
	customPath = "" // /sample-code/quickfeed
)

func GetMapPath(name string) string {
	if !strings.Contains(name, ".") {
		name = fmt.Sprintf("%s.json", name)
	}
	return filepath.Join(MapPath(), name)
}

func GetAbsPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path of file: %s, err: %v", path, err)
	}
	return absPath, nil
}

func ProjectPath() string {
	path := os.Getenv(refVizRootPath)
	if path == "" {
		panic("project path not set")
	}
	if customPath != "" {
		path = filepath.Join(path, customPath)
	}
	return path
}

func SetProjectPath(path string) error {
	return os.Setenv(refVizRootPath, path)
}

func getRootPath(name string) string {
	return filepath.Join(ProjectPath(), name)
}

func ConfigPath() string {
	return getRootPath(tmp("config.json"))
}

func CachePath() string {
	return getRootPath(tmp("cache.json"))
}

func GetTempFolderPath() string {
	return getRootPath(tempFolder)
}

// tmp returns the path of the temporary folder
func tmp(name string) string {
	return filepath.Join(tempFolder, name)
}

func MapPath() string {
	return getRootPath(tmp("maps"))
}

func GraphvizPath() string {
	return getRootPath(tmp("graphviz"))
}

func DotFilePath(mapName *string) string {
	return filepath.Join(GraphvizPath(), fmt.Sprintf("%s.dot", *mapName))
}

// getProjectRoot returns the root directory of the users project
// If the user is in a git project, it will return the root of the git repository
// Attempts to find the root of a project, if the user is not in a git repository
// TODO: This is problematic for non-go projects, solve this
func GetProjectRoot() (string, error) {
	if gitRoot, err := rootGitProject(); err == nil {
		return gitRoot, nil
	}
	root, err := getRoot()
	if err == nil {
		return root, nil
	}
	return "", fmt.Errorf("error getting project root: %v", err)
}

// rootGitProject returns the root directory of the git project
// Depends on the user being inside a git project
func rootGitProject() (string, error) {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(path)), nil
}

// getRoot attempts to find the root of a project
// Walks up the directory tree looking for a marker file
// If a marker is not found in the directory, it will walk up to the parent directories
// TODO: Add more markers for different types of projects
// TODO: Can use the content input of the user to determine the project root faster
func getRoot() (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current working directory: %v", err)
	}

	markers := []string{
		"go.mod",
	}
	found := "found marker"

	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if dirHasMarker(path, markers) {
				root = path
				return fmt.Errorf("%s", found)
			}
		}
		return nil
	}); err != nil {
		if err.Error() != found {
			return "", fmt.Errorf("error walking through directory: %v", err)
		}
		return root, nil
	}

	// Walk up the directory tree until a marker is found
	for {
		for _, marker := range markers {
			if _, err := os.Stat(filepath.Join(root, marker)); err == nil {
				return root, nil
			}
		}
		parent := filepath.Dir(root)
		if parent == root {
			// reached the root of the file system
			break
		}
		root = parent
	}

	return "", fmt.Errorf("error getting project root")
}

// dirHasMarker checks if a directory has a marker file
// Returns true if the directory has a marker file
func dirHasMarker(dir string, markers []string) bool {
	for i := range markers {
		if Exists(filepath.Join(dir, markers[i])) {
			return true
		}
	}
	return false
}
