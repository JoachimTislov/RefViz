package config

import (
	"fmt"
	"log" // Added import for log.Fatalf
	"net/url" // Import for url.ParseRequestURI
	"strings"

	"github.com/JoachimTislov/RefViz/internal"
	"github.com/JoachimTislov/RefViz/internal/path"
)

var config = NewConfig()

const Name = "config.json"

func NewConfig() *Config {
	return &Config{
		InExt:                newSbMap(".go"),
		ExDirs:               newSbMap("node_modules", ".git"),
		ExFiles:              newSbMap(),
		ExSymbols:            newSbMap(),
		ExFindRefsForSymbols: nameAndPrefix{Name: newSbMap("init", "main"), Prefix: newSbMap("Test")},
		BaseBranchLink:       "https://github.com/quickfeed/quickfeed/tree/master/",
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

type nameAndPrefix struct {
	Name   SbMap `json:"name,omitempty"`
	Prefix SbMap `json:"prefix,omitempty"`
}

type Config struct {
	InExt                SbMap         `json:"includedExtensions,omitempty"`
	ExDirs               SbMap         `json:"excludedDirectories,omitempty"`
	ExFiles              SbMap         `json:"excludedFiles,omitempty"`
	ExSymbols            SbMap         `json:"excludedSymbols,omitempty"`
	ExFindRefsForSymbols nameAndPrefix `json:"excludedFindRefsForSymbolsPrefix,omitempty"`
	// BaseBranchLink is the link to a branch of a github repository
	// Used to create links to a custom branch
	// Example: https://github.com/quickfeed/quickfeed/tree/master/ is path to the master branch
	BaseBranchLink string `json:"baseBranchLink,omitempty"`
}

type SbMap map[string]bool

func Get() *Config {
	return config
}

func NotValidSymbol(name string) bool {
	return config.ExSymbols[name]
}

func FindRefsForSymbols(name string) bool {
	if config.ExFindRefsForSymbols.Name[name] {
		return false
	}
	for prefix := range config.ExFindRefsForSymbols.Prefix {
		if strings.HasPrefix(name, prefix) {
			return false
		}
	}
	return true
}

func GetBaseBranchLink() string {
	return config.BaseBranchLink
}

func SetBaseBranchLink(link string) error {
	if link == "" {
		return nil // No error if the link is empty, just don't set it.
	}
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("invalid base branch link URL: %w", err)
	}
	config.BaseBranchLink = link
	// Removed save() call here, will be called by Save()
	return nil
}

func save() error {
	return internal.MarshalAndWriteToFile(config, path.Tmp(Name))
}

func Save() {
	if err := save(); err != nil {
		log.Fatalf("Error saving config: %v", err)
	}
}
