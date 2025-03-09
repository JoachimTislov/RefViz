package types

type Cache struct {
	Errors        []string                           `json:"errors,omitempty"`
	UnusedSymbols map[string]map[string]OrphanSymbol `json:"UnusedSymbols,omitempty"`
	Entries       map[string]CacheEntry              `json:"entries,omitempty"`
}

type OrphanSymbol struct {
	Dir      string `json:"dir,omitempty"`
	FileName string `json:"fileName,omitempty"`
	Location string `json:"location,omitempty"`
}

type CacheEntry struct {
	Name    string             `json:"name,omitempty"`
	ModTime int64              `json:"modTime,omitempty"`
	Symbols map[string]*Symbol `json:"symbols,omitempty"`
}
