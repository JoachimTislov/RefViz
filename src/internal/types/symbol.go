package types

import (
	"fmt"
	"log"
	"path/filepath"
)

type Symbol struct {
	Name     string          `json:"name,omitempty"`
	Kind     string          `json:"kind,omitempty"`
	Position Position        `json:"position,omitempty"`
	Path     string          `json:"path,omitempty"`
	FilePath string          `json:"filePath,omitempty"`
	Refs     map[string]*Ref `json:"refs,omitempty"`
	ZeroRefs bool            `json:"zeroRefs,omitempty"` // if true, the symbol has no references
}

func (s *Symbol) newSymbolRef(ref *Ref) SymbolRef {
	return SymbolRef{
		Definition: symbol{
			Name:     s.Name,
			Kind:     s.Kind,
			FilePath: s.FilePath,
		},
		Ref: *ref,
	}
}

func (s *Symbol) sortRefsIntoHierarchy(folderRefs, fileRefs *map[string]SymbolRef, folderPath, fileName *string, force *bool) Symbol {
	if folderRefs == nil || fileRefs == nil {
		log.Fatal("folderRefs or fileRefs is nil")
	}
	var refsToRemove []string
	for key, r := range s.Refs {
		sRef := s.newSymbolRef(r)

		folderPathDiffer := filepath.Dir(r.FilePath) != *folderPath
		fileNameDiffer := r.FileName != *fileName
		// if the reference is in a different folder or page, move it to the appropriate folder or file
		// and remove it from the symbol
		if folderPathDiffer || fileNameDiffer {
			refsToRemove = append(refsToRemove, key)
			refsPointer := fileRefs
			if folderPathDiffer {
				refsPointer = folderRefs
			}
			addEntryToMap(refsPointer, sRef.createSymbolMapKey(), sRef, force)
		}
	}
	for _, key := range refsToRemove {
		delete(s.Refs, key)
	}
	return *s
}

func addEntryToMap(m *map[string]SymbolRef, key string, sr SymbolRef, force *bool) {
	if *m == nil {
		*m = make(map[string]SymbolRef)
	}
	if _, ok := (*m)[key]; !ok || *force {
		(*m)[key] = sr
	} else {
		log.Printf("symbolRef already exists, definition name: %s, ref name: %s", sr.Definition.Name, sr.Ref.MethodName)
	}
}

// purely done to match the other refs slices type
// makes it easier to loop through later
func (s *Symbol) CreateSymbol() symbol {
	symbolRefs := make(map[string]SymbolRef)
	for _, ref := range s.Refs {
		symbolRefs[s.createSymbolMapKey(ref.FilePath, ref.MethodName)] = s.newSymbolRef(ref)
	}
	return symbol{
		Name:     s.Name,
		Kind:     s.Kind,
		FilePath: s.FilePath,
		Refs:     symbolRefs,
	}
}

func (s *Symbol) createSymbolMapKey(refPath, methodName string) string {
	return fmt.Sprintf("%s:%s_%s:%s", s.FilePath, s.Name, refPath, methodName)
}
