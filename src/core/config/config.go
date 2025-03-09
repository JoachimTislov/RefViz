package config

import (
	"fmt"

	"github.com/JoachimTislov/RefViz/internal"
	"github.com/JoachimTislov/RefViz/internal/path"
)

var config = NewConfig()

func NewConfig() *Config {
	return &Config{
		InExt:   newSbMap(".go"),
		ExDirs:  newSbMap("node_modules", ".git"),
		ExFiles: newSbMap(),
	}
}

func newSbMap(args ...string) SbMap {
	var m SbMap
	if len(args) == 0 {
		return m
	}
	m = make(SbMap)
	for i := range args {
		m[args[i]] = true
	}
	return m
}

type Config struct {
	InExt   SbMap `json:"includedExtensions,omitempty"`
	ExDirs  SbMap `json:"excludedDirectories,omitempty"`
	ExFiles SbMap `json:"excludedFiles,omitempty"`
}

type SbMap map[string]bool

func Get() *Config {
	return config
}

func (c *Config) Save() error {
	if err := internal.MarshalAndWriteToFile(config, path.Config()); err != nil {
		return fmt.Errorf("error updating configurations: %v", err)
	}
	return nil
}
