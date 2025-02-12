package types

func NewCache() *Cache {
	return &Cache{}
}

func NewCacheEntry(path string, modTime int64, symbols []*Symbol) *entry {
	return &entry{
		Path:    path,
		ModTime: modTime,
		Symbols: symbols,
	}
}

type Cache map[string]entry

type entry struct {
	Path    string    `json:"path"`
	ModTime int64     `json:"modTime"`
	Symbols []*Symbol `json:"symbols"`
}

type Symbol struct {
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Position Position `json:"position"`
	Refs     []*Ref   `json:"refs,omitempty"`
}

type Ref struct {
	Path       string `json:"path"`
	FolderName string `json:"folderName"`
	FileName   string `json:"fileName"`
	MethodName string `json:"methodName"`
}
