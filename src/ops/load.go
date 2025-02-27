package ops

import (
	"fmt"
	"os"

	"github.com/JoachimTislov/RefViz/internal"
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
	root, err := internal.GetProjectRoot()
	if err != nil {
		return err
	}
	if err := internal.SetProjectPath(root); err != nil {
		return err
	}
	return nil
}

// initFolder initializes the project folder if it does not exist
func initFolder() error {
	folderPaths := []string{internal.GetTempFolderPath(), internal.MapPath(), internal.GraphvizPath()}
	for _, p := range folderPaths {
		if !internal.Exists(p) {
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
	if err := loadFile(internal.ConfigPath(), config); err != nil {
		return fmt.Errorf("error loading configurations: %v", err)
	}
	return nil
}

func loadCache() error {
	if err := loadFile(internal.CachePath(), cache); err != nil {
		return fmt.Errorf("error loading cache: %v", err)
	}
	return nil
}

func loadFile(path string, v any) error {
	if !internal.Exists(path) {
		if err := marshalAndWriteToFile(v, path); err != nil {
			return fmt.Errorf("error creating config file: %v", err)
		}
	} else {
		if err := getFile(path, v); err != nil {
			return fmt.Errorf("error getting config file: %v", err)
		}
	}
	return nil
}
