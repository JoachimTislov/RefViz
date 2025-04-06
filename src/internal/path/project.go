package path

import (
	"path/filepath"
)

func QuickfeedRootMain() string {
	return filepath.Join(Quickfeed(), "main.go")
}

func Quickfeed() string {
	return filepath.Join(refVizStaticRoot(), "sample-code", "quickfeed")
}

func TestData() string {
	return filepath.Join(refVizStaticRoot(), "testData")
}

func refVizStaticRoot() string {
	root, err := getProjectRoot()
	if err != nil {
		panic("error getting user home directory")
	}
	return root
}
