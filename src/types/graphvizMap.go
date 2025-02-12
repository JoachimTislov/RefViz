package types

func NewGraphvizMap(name *string) *GraphvizMap {
	return &GraphvizMap{
		Root: &folderMap{},
		Name: *name,
	}
}

// Recursive data structure to store the project structure.
// Used for graphviz file generation
type GraphvizMap struct {
	Root *folderMap `json:"root"`
	Name string     `json:"name"`
}

type folderMap map[string]*folder

type folder struct {
	FolderName string       `json:"folderName"`
	FolderPath string       `json:"folderPath"`
	Refs       []ref        `json:"refs,omitempty"`
	Files      []file       `json:"files,omitempty"`
	Errors     []GoplsError `json:"errors,omitempty"`
	SubFolders *folderMap   `json:"subFolders,omitempty"`
}

type file struct {
	Name    string `json:"name"`
	Path    string
	Refs    []ref    `json:"refs,omitempty"`
	Symbols []symbol `json:"symbols,omitempty"`
}

type symbol struct {
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Position Position `json:"position"`
	Refs     []ref    `json:"refs,omitempty"`
}

// source is there since the symbol is a reference to a symbol in another file
// Will result in duplicate data, but it's needed to keep track of the source
type ref struct {
	Source refInfo `json:"source"`
	Info   refInfo `json:"info"`
}

type refInfo struct {
	Path       string
	FolderName string `json:"folderName"`
	FileName   string `json:"fileName"`
	MethodName string `json:"methodName"`
}

// remove entries with zero files and sub folders
func (m *GraphvizMap) Clean() error {
	for _, key := range getKeysToDelete(m.Root) {
		delete(*m.Root, key)
	}
	return nil
}

// getKeysToDelete recursively finds all keys with zero files and sub folders
func getKeysToDelete(m *folderMap) []string {
	var keysToDelete []string
	for key, folder := range *m {
		if len(*folder.SubFolders) == 0 {
			if len(folder.Files) == 0 {
				keysToDelete = append(keysToDelete, key)
			}
		} else {
			getKeysToDelete(folder.SubFolders)
		}
	}
	return keysToDelete
}
