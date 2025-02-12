package op

import (
	"fmt"
	"os"

	"github.com/JoachimTislov/Project-Visualizer/types"
)

func LoadDefs() error {
	if err := loadRootPath(); err != nil {
		return fmt.Errorf("error loading root path: %v", err)
	}
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

// loadRootPath loads the root path of the project and sets it
func loadRootPath() error {
	root, err := getProjectRoot()
	if err != nil {
		return err
	}
	if err := setProjectPath(root); err != nil {
		return err
	}
	return nil
}

// loadConfig creates default config file if it does not exist
// If the file does exist it reads the file and unmarshals it into the config variable
func loadConfig() error {
	path := configPath()
	if !exists(path) {
		if err := marshalAndWriteToFile(config, path); err != nil {
			return fmt.Errorf("error creating config file: %v", err)
		}
	} else {
		if err := GetFile(path, config); err != nil {
			return fmt.Errorf("error getting config file: %v", err)
		}
	}
	return nil
}

// initFolder initializes the project folder if it does not exist
func initFolder() error {
	path := getTempFolderPath()
	if !exists(path) {
		if err := os.Mkdir(path, 0755); err != nil {
			return fmt.Errorf("error creating project folder: %v", err)
		}
	}
	return nil
}

func loadCache() error {
	if exists(cachePath()) {
		if err := GetFile(cachePath(), &cache); err != nil {
			return fmt.Errorf("error getting cache file: %v", err)
		}
	} else {
		cache = types.NewCache()
		if err := marshalAndWriteToFile(cache, cachePath()); err != nil {
			return fmt.Errorf("error creating cache file: %v", err)
		}
	}
	return nil
}
