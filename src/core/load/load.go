package load

import (
	"fmt"
	"os"

	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/core/config"
	"github.com/JoachimTislov/RefViz/internal"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/utils"
)

func Defs() error {
	_ = path.Project() // load project path, will panic if is fails

	if err := initFolder(); err != nil {
		return fmt.Errorf("error initializing project folder: %v", err)
	}
	if err := loadConfig(); err != nil {
		return fmt.Errorf("error loading configurations: %v", err)
	}
	if err := loadCache(); err != nil {
		return fmt.Errorf("error loading cache: %v", err)
	}
	return nil
}

// initFolder initializes the project folder if it does not exist
func initFolder() error {
	folderPaths := []string{path.GetTempFolder(), path.Map(), path.Graphviz()}
	for _, p := range folderPaths {
		if !utils.Exists(p) {
			if err := os.Mkdir(p, 0755); err != nil {
				return fmt.Errorf("error creating project folder: %v", err)
			}
		}
	}
	return nil
}

// loadConfig creates default config file if it does not exist
// If the file does exist it reads the file and unmarshals it into the config variable
func loadConfig() error {
	if err := loadFile(path.Config(), config.Get()); err != nil {
		return fmt.Errorf("error loading configurations: %v", err)
	}
	return nil
}

func loadCache() error {
	if err := loadFile(path.Cache(), cache.Get()); err != nil {
		return fmt.Errorf("error loading cache: %v", err)
	}
	return nil
}

func loadFile(path string, v any) error {
	if !utils.Exists(path) {
		if err := internal.MarshalAndWriteToFile(v, path); err != nil {
			return fmt.Errorf("error creating config file: %v", err)
		}
	} else {
		if err := internal.GetAndUnmarshalFile(path, v); err != nil {
			return fmt.Errorf("error getting config file: %v", err)
		}
	}
	return nil
}
