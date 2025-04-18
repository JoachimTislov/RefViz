package types

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Packages map[string]*Pkg

func (m Packages) Add(name, dir string) *Pkg {
	key := name
	dirName := filepath.Base(dir)
	if name != filepath.Base(dir) && !strings.Contains(name, "_test") {
		key = fmt.Sprintf("%s-%s", name, dirName)
	}
	m[key] = &Pkg{
		Path:  dir,
		Files: make(map[string]*File2),
	}
	return m[key]
}

type Pkg struct {
	Path  string
	Files map[string]*File2
}

func (p *Pkg) AddFile(filePath string) *File2 {
	fileName := filepath.Base(filePath)
	p.Files[fileName] = &File2{
		Path:    filePath,
		Symbols: make(map[string]*Symbol),
	}
	return p.Files[fileName]
}

type File2 struct {
	Path    string
	Symbols map[string]*Symbol
}

func (f *File2) AddSymbol(name, filePath, kind string) {
	f.Symbols[name] = &Symbol{
		Name:     name,
		FilePath: filePath,
		Kind:     kind,
	}
}
