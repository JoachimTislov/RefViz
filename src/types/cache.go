package types

import "sync"

func NewCache() *Cache {
	return &Cache{
		UnusedSymbols: make(map[string]map[string]UnusedSymbol),
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

func NewUnusedSymbol(dir, fileName, location string) UnusedSymbol {
	return UnusedSymbol{
		Dir:      dir,
		FileName: fileName,
		Location: location,
	}
}

func (c *Cache) AddEntry(relPath string, entry *CacheEntry) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.Entries[relPath] = *entry
}

func (c *Cache) GetEntry(relPath string) *CacheEntry {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	if entry, ok := c.Entries[relPath]; ok {
		return &entry
	}
	return &CacheEntry{Symbols: make(map[string]*Symbol)}
}

func (c *Cache) AddUnusedSymbol(relPath string, name string, symbol UnusedSymbol) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if c.UnusedSymbols[relPath] == nil {
		c.UnusedSymbols[relPath] = make(map[string]UnusedSymbol)
	}
	c.UnusedSymbols[relPath][name] = symbol
}

type Cache struct {
	UnusedSymbols map[string]map[string]UnusedSymbol `json:"UnusedSymbols,omitempty"`
	Entries       map[string]CacheEntry              `json:"entries,omitempty"`
	Mu            sync.RWMutex                       `json:"mu,omitempty"`
}

type UnusedSymbol struct {
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
	ZeroRefs bool            `json:"zeroRefs,omitempty"` // if true, the symbol has no references
}

type Ref struct {
	Path       string `json:"path"`
	FolderName string `json:"folderName"`
	FileName   string `json:"fileName"`
	MethodName string `json:"methodName"`
}
