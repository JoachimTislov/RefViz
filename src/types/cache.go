package types

type Cache map[string]entry

type entry struct {
	Path    string     `json:"path"` // relative path to the file
	ModTime int64      `json:"modTime"`
	Symbols *[]*Symbol `json:"symbols"`
}

type Symbol struct {
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Position Position `json:"position"`
	Refs     *[]*Ref  `json:"refs,omitempty"`
}

type Ref struct {
	Path       string `json:"path"`
	FolderName string `json:"folderName"`
	FileName   string `json:"fileName"`
	MethodName string `json:"methodName"`
}
