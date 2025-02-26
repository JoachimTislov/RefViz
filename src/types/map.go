package types

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func NewMap(name *string) *RMap {
	return &RMap{
		Name:  *name,
		Nodes: map[string]*Node{},
	}
}

func newNode(name, path string) *Node {
	return &Node{
		Name:       name,
		RootFolder: newFolder(path),
	}
}

func (m *RMap) GetOrCreateNode(nodeName *string, projectPath string) (*Node, error) {
	if m.Nodes == nil {
		m.Nodes = make(map[string]*Node)
	}
	node, ok := m.Nodes[*nodeName]
	if !ok {
		m.Nodes[*nodeName] = newNode(*nodeName, projectPath)
		node = m.Nodes[*nodeName]
	}
	return node, nil
}

func (m *RMap) AddNode(nodeName *string, projectPath string) {
	if m.Nodes == nil {
		m.Nodes = make(map[string]*Node)
	}
	m.Nodes[*nodeName] = newNode(*nodeName, projectPath)
}

func newFolder(path string) *Folder {
	return &Folder{
		FolderName: filepath.Base(path),
		FolderPath: path,
		Files:      make(map[string]File),
		SubFolders: make(map[string]*Folder),
	}
}
func (f *File) AddSymbols(folderRefs *map[string]SymbolRef, symbols *map[string]*Symbol, folderPath, fileName *string) {
	for _, s := range *symbols {
		f.AddSymbol(s.SortRefsIntoHierarchy(folderRefs, &f.Refs, folderPath, fileName))
	}
}

func (f *File) AddSymbol(s Symbol) {
	if f.Symbols == nil {
		f.Symbols = make(map[string]symbol)
	}
	if _, ok := f.Symbols[s.Name]; !ok {
		f.Symbols[s.Name] = s.createSymbol()
	} else {
		log.Printf("symbol: %s already exists in file: %s", s.Name, f.Name)
	}
}

func (s *Symbol) newSymbolRef(ref *Ref) SymbolRef {
	return SymbolRef{
		Definition: Symbol{
			Name: s.Name,
			Kind: s.Kind,
			Path: s.Path,
		},
		Ref: *ref,
	}
}

func (s *Symbol) SortRefsIntoHierarchy(folderRefs, fileRefs *map[string]SymbolRef, folderPath, fileName *string) Symbol {
	var refsToRemove []string
	for key, r := range s.Refs {
		sRef := s.newSymbolRef(r)

		folderPageDiffer := filepath.Dir(r.Path) != *folderPath
		fileNameDiffer := r.FileName != *fileName
		// if the reference is in a different folder or page, move it to the appropriate folder or file
		// and remove it from the symbol
		if folderPageDiffer || fileNameDiffer {
			refsToRemove = append(refsToRemove, key)
			key := fmt.Sprintf("%s_%s", sRef.Definition.Path, sRef.Ref.Path)
			if folderPageDiffer {
				addEntryToMap(folderRefs, key, sRef)
			}
			if fileNameDiffer {
				addEntryToMap(fileRefs, key, sRef)
			}
		}
	}
	for _, key := range refsToRemove {
		delete(s.Refs, key)
	}
	return *s
}

func addEntryToMap[T any](m *map[string]T, key string, entry T) {
	if m == nil {
		panic(fmt.Sprintf("map: %s is nil", reflect.TypeOf(m)))
	}
	if _, ok := (*m)[key]; !ok {
		(*m)[key] = entry
	} else {
		log.Printf("entry: %v already exists in map: %v, key: %s", entry, m, key)
	}
}

// purely done to match the other refs slices type
// makes it easier to loop through later
func (s Symbol) createSymbol() symbol {
	symbolRefs := make(map[string]SymbolRef)
	for _, ref := range s.Refs {
		key := fmt.Sprintf("%s_%s", s.Path, ref.Path)
		symbolRefs[key] = s.newSymbolRef(ref)
	}
	return symbol{
		Name: s.Name,
		Kind: s.Kind,
		Path: s.Path,
		Refs: symbolRefs,
	}
}

func newFile(name string, path string) File {
	return File{
		Name: name,
		Path: path,
	}
}

func (f *Folder) GetFile(fileName, folderPath *string) *File {
	if f.Files == nil {
		return nil
	}
	file, ok := f.Files[*fileName]
	if !ok {
		file = newFile(*fileName, *folderPath)
	}
	return &file
}

func (f *Folder) AddFile(file *File, force bool) {
	if f.Files == nil {
		f.Files = make(map[string]File)
	}
	if _, ok := f.Files[file.Name]; !ok || force {
		f.Files[file.Name] = *file
	} else {
		log.Printf("file: %s already exists in folder: %s", file.Name, f.FolderName)
	}
}

// The pointer complexity of this function is quite annoying
// Essentially, it gets the related folder based on the absolute path
// Updates the local pointer in the method and return the pointer to the related folder
// This does not override the original folder
// *f = *f.SubFolders[d] instead of f = f.SubFolders[d] will override the original folder
// The updated local pointer is therefore returned, and the original folder how the natural path of folders
func (f *Folder) GetRelatedFolder(absPath, projectPath string) (*Folder, error) {
	dirs, err := determineFolderPath(absPath, projectPath)
	if err != nil {
		return nil, err
	}
	for _, d := range *dirs {
		projectPath = filepath.Join(projectPath, d)

		if f.SubFolders == nil {
			f.SubFolders = make(map[string]*Folder)
		}

		if _, exists := f.SubFolders[d]; !exists {
			f.SubFolders[d] = newFolder(projectPath)
		}
		f = f.SubFolders[d]
	}
	return f, nil
}

func determineFolderPath(absPath, projectPath string) (*[]string, error) {
	relPath, err := filepath.Rel(projectPath, absPath)
	if err != nil {
		return nil, fmt.Errorf("error getting relative path: %s, err: %v", absPath, err)
	}
	dirs := []string{relPath}
	if strings.Contains(relPath, string(filepath.Separator)) {
		if f, err := os.Stat(absPath); err == nil && !f.IsDir() {
			dirs = strings.Split(filepath.Dir(relPath), string(filepath.Separator))
		} else if err != nil {
			return nil, fmt.Errorf("error getting directory name: %s, err: %v", absPath, err)
		}
	}
	return &dirs, nil
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
	FolderName string               `json:"folderName"`
	FolderPath string               `json:"folderPath"`
	Refs       map[string]SymbolRef `json:"refs,omitempty"`
	Files      map[string]File      `json:"files,omitempty"`
	SubFolders map[string]*Folder   `json:"subFolders,omitempty"`
}

type File struct {
	Name    string               `json:"name,omitempty"`
	Path    string               `json:"path,omitempty"`
	Refs    map[string]SymbolRef `json:"refs,omitempty"`
	Symbols map[string]symbol    `json:"symbols,omitempty"`
}

type symbol struct {
	Name string               `json:"name,omitempty"`
	Kind string               `json:"kind,omitempty"`
	Path string               `json:"path,omitempty"`
	Refs map[string]SymbolRef `json:"refs,omitempty"`
}

type SymbolRef struct {
	Definition Symbol `json:"definition"`
	Ref        Ref    `json:"reference"`
}
