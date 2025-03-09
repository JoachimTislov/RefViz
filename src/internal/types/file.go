package types

type File struct {
	Name    string               `json:"name,omitempty"`
	Path    string               `json:"path,omitempty"`
	Refs    map[string]SymbolRef `json:"refs,omitempty"`
	Symbols map[string]symbol    `json:"symbols,omitempty"`
}

func newFile(name string, path string) *File {
	return &File{
		Name: name,
		Path: path,
	}
}

func (f *File) AddSymbols(folderRefs *map[string]SymbolRef, symbols *map[string]*Symbol, fullFolderPath, fileName *string, force *bool) {
	for _, s := range *symbols {
		f.addSymbol(s.sortRefsIntoHierarchy(folderRefs, &f.Refs, fullFolderPath, fileName, force), force)
	}
}

func (f *File) addSymbol(s Symbol, force *bool) {
	if f.Symbols == nil {
		f.Symbols = make(map[string]symbol)
	}
	if _, ok := f.Symbols[s.Name]; !ok || *force {
		f.Symbols[s.Name] = s.CreateSymbol()
	} else {
		//log.Printf("symbol: %s already exists in file: %s", s.Name, f.Name)
	}
}
