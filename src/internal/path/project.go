package path

import (
	"os"
	"path/filepath"
)

func QuickfeedRootMain() string {
	return filepath.Join(Quickfeed(), "main.go")
}

func Quickfeed() string {
	return filepath.Join(sampleCode(), "quickfeed")
}

func sampleCode() string {
	return filepath.Join(refVizStaticRoot(), "sample-code")
}

func TestData() string {
	return filepath.Join(refVizStaticRoot(), "testData")
}

const refViz = "RefViz"

func refVizStaticRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("error getting user home directory")
	}
	return filepath.Join(home, refViz)
}
