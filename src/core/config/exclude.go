package config

import (
	"log"
	"strings"
)

func Exclude(exDir, exFile, inExt, exSymbol, exFindRefsForSymbols, exFindRefsForSymbolsPrefix *string) {
	if *exDir != "" {
		if err := addExDirs(*exDir); err != nil {
			log.Fatalf("Error excluding directory: %v\n", err)
		}
	}
	if *exFile != "" {
		if err := addExFiles(*exFile); err != nil {
			log.Fatalf("Error excluding file: %v\n", err)
		}
	}
	if *inExt != "" {
		if err := addInExt(*inExt); err != nil {
			log.Fatalf("Error including extension: %v\n", err)
		}
	}
	if *exSymbol != "" {
		if err := addExSymbols(*exSymbol); err != nil {
			log.Fatalf("Error excluding symbol: %v\n", err)
		}
	}
	if *exFindRefsForSymbols != "" {
		if err := addExFindRefsForSymbols(*exFindRefsForSymbols); err != nil {
			log.Fatalf("Error excluding find refs for symbols: %v\n", err)
		}
	}
	if *exFindRefsForSymbolsPrefix != "" {
		if err := addExFindRefsForSymbolsPrefix(*exFindRefsForSymbolsPrefix); err != nil {
			log.Fatalf("Error excluding find refs for symbols prefix: %v\n", err)
		}
	}
}

func addExDirs(dir string) error {
	return exclude(&config.ExDirs, dir)
}

func addExFiles(file string) error {
	if config.ExFiles == nil {
		sMap := newSbMap(file)
		config.ExFiles = sMap
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
	return exclude(&config.ExFindRefsForSymbols.Name, symbol)
}

func addExFindRefsForSymbolsPrefix(prefix string) error {
	return exclude(&config.ExFindRefsForSymbols.Prefix, prefix)
}

func exclude(m *SbMap, item string) error {
	if (*m)[item] {
		log.Fatalf("Item %s already exists\n", item)
	}
	(*m)[item] = true
	return save()
}
