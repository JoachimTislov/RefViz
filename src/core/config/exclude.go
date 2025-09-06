package config

import (
	"fmt" // Added import for fmt.Errorf
	"strings"
)

func Exclude(exDir, exFile, inExt, exSymbol, exFindRefsForSymbols, exFindRefsForSymbolsPrefix *string) error {
	if *exDir != "" {
		if err := addExDirs(*exDir); err != nil {
			return fmt.Errorf("error excluding directory: %w", err)
		}
	}
	if *exFile != "" {
		if err := addExFiles(*exFile); err != nil {
			return fmt.Errorf("error excluding file: %w", err)
		}
	}
	if *inExt != "" {
		if err := addInExt(*inExt); err != nil {
			return fmt.Errorf("error including extension: %w", err)
		}
	}
	if *exSymbol != "" {
		if err := addExSymbols(*exSymbol); err != nil {
			return fmt.Errorf("error excluding symbol: %w", err)
		}
	}
	if *exFindRefsForSymbols != "" {
		if err := addExFindRefsForSymbols(*exFindRefsForSymbols); err != nil {
			return fmt.Errorf("error excluding find refs for symbols: %w", err)
		}
	}
	if *exFindRefsForSymbolsPrefix != "" {
		if err := addExFindRefsForSymbolsPrefix(*exFindRefsForSymbolsPrefix); err != nil {
			return fmt.Errorf("error excluding find refs for symbols prefix: %w", err)
		}
	}
	return nil
}

func addExDirs(dir string) error {
	return exclude(&config.ExDirs, dir)
}

func addExFiles(file string) error {
	if config.ExFiles == nil {
		config.ExFiles = newSbMap()
	}
	return exclude(&config.ExFiles, file)
}

func addInExt(ext string) error {
	if !strings.Contains(ext, ".") {
		ext = "." + ext
	}
	return exclude(&config.InExt, ext)
}

func addExSymbols(symbol string) error {
	return exclude(&config.ExSymbols, symbol)
}

func addExFindRefsForSymbols(symbol string) error {
	if config.ExFindRefsForSymbols.Name == nil {
		config.ExFindRefsForSymbols.Name = newSbMap()
	}
	return exclude(&config.ExFindRefsForSymbols.Name, symbol)
}

func addExFindRefsForSymbolsPrefix(prefix string) error {
	if config.ExFindRefsForSymbols.Prefix == nil {
		config.ExFindRefsForSymbols.Prefix = newSbMap()
	}
	return exclude(&config.ExFindRefsForSymbols.Prefix, prefix)
}

func exclude(m *SbMap, item string) error {
	if (*m)[item] {
		return fmt.Errorf("item %s already exists in config", item) // Return an error instead of log.Fatalf
	}
	if *m == nil { // Ensure map is initialized
		*m = newSbMap()
	}
	(*m)[item] = true
	return nil // Removed save() call
}
