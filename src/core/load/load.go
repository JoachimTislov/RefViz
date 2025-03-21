package load

import (
	"fmt"

	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/core/config"
	"github.com/JoachimTislov/RefViz/internal"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/utils"
)

// loadConfig creates default config file if it does not exist
// If the file does exist it reads the file and unmarshals it into the config variable
func loadConfig(name string) error {
	if err := File(path.Tmp(name), config.Get()); err != nil {
		return fmt.Errorf("error loading configurations: %v", err)
	}
	return nil
}

func loadCache(name string) error {
	if err := File(path.Tmp(name), cache.Get()); err != nil {
		return fmt.Errorf("error loading cache: %v", err)
	}
	return nil
}

// loadFile creates a file if it does not exist
// If the file does exist it reads the file and unmarshals it into the v variable
func File(path string, v any) error {
	if !utils.Exists(path) {
		if err := internal.MarshalAndWriteToFile(v, path); err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
	} else {
		if err := internal.GetAndUnmarshalFile(path, v); err != nil {
			return fmt.Errorf("error getting file: %v", err)
		}
	}
	return nil
}
