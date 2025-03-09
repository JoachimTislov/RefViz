package types

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Folder struct {
	FolderName string               `json:"folderName"`
	FolderPath string               `json:"folderPath"`
	Refs       map[string]SymbolRef `json:"refs,omitempty"`
	Files      map[string]*File     `json:"files,omitempty"`
	SubFolders map[string]*Folder   `json:"subFolders,omitempty"`
}

func newFolder(path string) *Folder {
	return &Folder{
		FolderName: filepath.Base(path),
		FolderPath: path,
		Files:      make(map[string]*File),
		SubFolders: make(map[string]*Folder),
	}
}

func (f *Folder) GetFile(fileName, folderPath *string) *File {
	if f.Files == nil {
		f.Files = make(map[string]*File)
	}
	file, ok := f.Files[*fileName]
	if !ok {
		file = newFile(*fileName, *folderPath)
		f.Files[*fileName] = file
	}
	return file
}

func (f *Folder) AddFile(file *File, forceUpdate *bool) {
	if f.Files == nil {
		f.Files = make(map[string]*File)
	}
	if _, ok := f.Files[file.Name]; !ok || *forceUpdate {
		f.Files[file.Name] = file
	} else {
		//log.Printf("file: %s already exists in folder: %s", file.Name, f.FolderName)
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

func (f *Folder) createMissingSymbols(projectPath string) error {
	var refs []SymbolRef
	f.getRefs(&refs)
	for _, ref := range refs {
		folder, err := f.GetRelatedFolder(ref.Ref.FilePath, projectPath)
		if err != nil {
			return fmt.Errorf("error getting related folder: %v", err)
		}
		name := ref.Ref.MethodName
		file := folder.GetFile(&ref.Ref.FileName, &folder.FolderPath)
		if _, ok := file.Symbols[name]; !ok {
			if file.Symbols == nil {
				(*file).Symbols = make(map[string]symbol)
			}
			(*file).Symbols[name] = symbol{
				Name:     name,
				FilePath: ref.Ref.FilePath,
			}
		}
	}
	return nil
}

func (f *Folder) getRefs(refs *[]SymbolRef) {
	if f.SubFolders == nil {
		return
	}
	for _, folder := range f.SubFolders {
		for _, ref := range folder.Refs {
			*refs = append(*refs, ref)
		}
		for _, file := range folder.Files {
			for _, ref := range file.Refs {
				*refs = append(*refs, ref)
			}
		}
		folder.getRefs(refs)
	}
}
