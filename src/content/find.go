package content

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/JoachimTislov/RefViz/core/config"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/prompt"
)

// findContent walks for the content root and attempts to find the content
// Returns early if the content is an empty string, equal to scanning everything from project root
func Find(content *string, ask bool) ([]string, error) {
	var paths []string
	var err error
	projectRootPath := path.Project()
	switch {
	case *content == "":
		paths = append(paths, projectRootPath)
	case filepath.IsAbs(*content):
		paths = append(paths, *content)
	default:
		err = filepath.WalkDir(projectRootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error walking from path: %s, err: %v", path, err)
			}
			if filepath.Base(path) == *content && config.IsValid(d.IsDir(), path) {
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error finding content: %s, err: %v", *content, err)
		}
	}

	if len(paths) == 0 {
		log.Fatal("No content found")
	}
	if len(paths) == 1 {
		return paths, nil
	}

	if ask {
		paths, err = prompt.AskUser(paths, []string{})
		if err != nil {
			return nil, fmt.Errorf("error asking user: %v", err)
		}
	}

	return paths, nil
}
