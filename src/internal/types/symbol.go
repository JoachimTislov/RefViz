package types

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/RefViz/core/config"
)

type Symbol struct {
	Name         string                  `json:"name,omitempty"`
	Kind         string                  `json:"kind,omitempty"`
	Position     Position                `json:"position,omitempty"`
	Path         string                  `json:"path,omitempty"`
	FilePath     string                  `json:"filePath,omitempty"`
	Children     []*Symbol               `json:"children,omitempty"`
	ChildSymbols map[string]*ChildSymbol `json:"childSymbols,omitempty"`
	Refs         map[string]*Ref         `json:"refs,omitempty"`
	// if true, the symbol has no references
	ZeroRefs bool `json:"zeroRefs,omitempty"`
}

type ChildSymbol struct {
	Key      string `json:"name,omitempty"`
	FilePath string `json:"filePath,omitempty"`
}

func (s *Symbol) newSymbolRef(ref *Ref) SymbolRef {
	return SymbolRef{
		Definition: definition{
			Name:     s.Name,
			Kind:     s.Kind,
			FilePath: s.FilePath,
		},
		Ref: *ref,
	}
}

// Removes child symbols and keeps only the related references
func (s *Symbol) Clean(methodName string) {
	s.ChildSymbols = nil
	var refs = make(map[string]*Ref)
	for key, ref := range s.Refs {
		if ref.MethodName == methodName {
			refs[key] = ref
		}
	}
	s.Refs = refs
}

func (s *Symbol) AddChildSymbol(name, filePath, relPath string) {
	if s.ChildSymbols == nil {
		s.ChildSymbols = make(map[string]*ChildSymbol)
	}
	s.ChildSymbols[name] = &ChildSymbol{
		Key:      relPath,
		FilePath: filePath,
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
func (s *Symbol) createDefinition() *definition {
	symbolRefs := make(map[string]SymbolRef)
	for _, ref := range s.Refs {
		symbolRefs[s.createSymbolMapKey(ref.FilePath, ref.MethodName)] = s.newSymbolRef(ref)
	}
	return &definition{
		Name:     s.Name,
		Kind:     s.Kind,
		FilePath: s.FilePath,
		Refs:     symbolRefs,
	}
}

func (s *Symbol) createSymbolMapKey(refPath, methodName string) string {
	return fmt.Sprintf("%s:%s_%s:%s", s.FilePath, s.Name, refPath, methodName)
}

func (s *Symbol) createGithubLink(symbol *Symbol) string {
	baseLink := config.GetBaseBranchLink()

	split2 := strings.Split(strings.Split(symbol.Path, "/quickfeed/")[1], ":")
	partialLink := split2[0] + "#L" + split2[1] + "-L" + strings.Split(split2[2], "-")[1]

	return fmt.Sprintf("[%s](%s%s), ", symbol.Name, baseLink, partialLink)
}
