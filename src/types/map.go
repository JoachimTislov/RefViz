package types

import "path/filepath"

func NewMap(name *string) *RMap {
	return &RMap{
		Name:  *name,
		Nodes: map[string]*Node{},
	}
}

func NewNode(name, path string) *Node {
	return &Node{
		Name:       name,
		RootFolder: NewFolder(path),
	}
}

func NewFolder(path string) *Folder {
	return &Folder{
		FolderName: filepath.Base(path),
		FolderPath: path,
		SubFolders: map[string]*Folder{},
	}
}

func (f *Folder) AddFile(file File) {
	f.Files = append(f.Files, file)
}

func (f *Folder) AddRef(ref SymbolRef) {
	f.Refs = append(f.Refs, ref)
}

func (f *Folder) MoveToSubFolder(name, path string) {
	f = f.SubFolders[name]
	if f == nil {
		f = NewFolder(path)
	}
}

func NewFile(name string, path string) File {
	return File{
		Name: name,
		Path: path,
	}
}

// Recursive data structure to store the project structure.
// Used for graphviz file generation
type RMap struct {
	Name  string           `json:"name"`
	Nodes map[string]*Node `json:"nodes"`
}

type Node struct {
	Name       string  `json:"name"`
	RootFolder *Folder `json:"rootFolder"`
}

type Folder struct {
	FolderName string             `json:"folderName"`
	FolderPath string             `json:"folderPath"`
	Refs       []SymbolRef        `json:"refs,omitempty"`
	Files      []File             `json:"files,omitempty"`
	SubFolders map[string]*Folder `json:"subFolders,omitempty"`
}

type File struct {
	Name    string      `json:"name"`
	Path    string      `json:"path"`
	Refs    []SymbolRef `json:"refs,omitempty"`
	Symbols []Symbol    `json:"symbols,omitempty"`
}

type SymbolRef struct {
	Definition Symbol `json:"definition"`
	Ref        Ref    `json:"reference"`
}
