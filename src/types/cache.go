package types

func NewCache() *Cache {
	return &Cache{
		UnusedSymbols: make(map[string][]unusedSymbol),
		Entries:       make(map[string]CacheEntry),
	}
}

func NewCacheEntry(name string, modTime int64, symbols map[string]*Symbol) CacheEntry {
	return CacheEntry{
		Name:    name,
		ModTime: modTime,
		Symbols: symbols,
	}
}

func NewUnusedSymbol(name, dir, fileName, location string) unusedSymbol {
	return unusedSymbol{
		Name:     name,
		Dir:      dir,
		FileName: fileName,
		Location: location,
	}
}

type Cache struct {
	UnusedSymbols map[string][]unusedSymbol `json:"unusedSymbols,omitempty"`
	Entries       map[string]CacheEntry     `json:"entries,omitempty"`
}

type unusedSymbol struct {
	Name     string `json:"name"`
	Dir      string `json:"dir"`
	FileName string `json:"fileName"`
	Location string `json:"location"`
}

type CacheEntry struct {
	Name    string             `json:"name"`
	ModTime int64              `json:"modTime"`
	Symbols map[string]*Symbol `json:"symbols"`
}

type Symbol struct {
	Name     string          `json:"name"`
	Kind     string          `json:"kind"`
	Position Position        `json:"position"`
	Refs     map[string]*Ref `json:"refs,omitempty"`
}

type Ref struct {
	Path         string `json:"path"`
	FolderName   string `json:"folderName"`
	FileName     string `json:"fileName"`
	ParentSymbol Symbol `json:"parentSymbol"`
}
