package types

type File struct {
	Name    string                 `json:"name,omitempty"`
	Path    string                 `json:"path,omitempty"`
	Refs    map[string]SymbolRef   `json:"refs,omitempty"`
	Symbols map[string]*definition `json:"symbols,omitempty"`
}

func newFile(name string, path string) *File {
	return &File{
		Name: name,
		Path: path,
	}
}

func (f *File) AddSymbols(folderRefs *map[string]SymbolRef, symbols *map[string]*Symbol, fullFolderPath, fileName *string, forceUpdate *bool) {
	for _, s := range *symbols {
		f.addSymbol(s.sortRefsIntoHierarchy(folderRefs, &f.Refs, fullFolderPath, fileName, forceUpdate), forceUpdate)
	}
}

func (f *File) addSymbol(s Symbol, forceUpdate *bool) {
	if f.Symbols == nil {
		f.Symbols = make(map[string]*definition)
	}
	if _, ok := f.Symbols[s.Name]; !ok || *forceUpdate {
		f.Symbols[s.Name] = s.createDefinition()
	} else {
		//log.Printf("symbol: %s already exists in file: %s", s.Name, f.Name)
	}
}
