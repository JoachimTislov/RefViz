package path

import (
	"path/filepath"
)

func QuickfeedRootMain() string {
	return filepath.Join(Quickfeed(), "main.go")
}

func Quickfeed() string {
	return filepath.Join(refvizStaticRoot(), "sample-code", "quickfeed")
}

func TestData() string {
	return filepath.Join(refvizStaticRoot(), "testData")
}

func refvizStaticRoot() string {
	root, err := getProjectRoot()
	if err != nil {
		panic("error getting user home directory")
	}
	return root
}
