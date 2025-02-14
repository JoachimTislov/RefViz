package types

func NewCache() *Cache {
	return &Cache{}
}

func NewCacheEntry(path string, modTime int64, symbols map[string]*Symbol) CacheEntry {
	return CacheEntry{
		Path:    path,
		ModTime: modTime,
		Symbols: symbols,
	}
}

type Cache map[string]CacheEntry

type CacheEntry struct {
	Path    string             `json:"path"`
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
