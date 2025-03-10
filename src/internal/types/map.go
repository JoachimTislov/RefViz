package types

import (
	"fmt"

	"github.com/JoachimTislov/RefViz/internal"
)

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

type definition struct {
	Name     string               `json:"name,omitempty"`
	Kind     string               `json:"kind,omitempty"`
	FilePath string               `json:"path,omitempty"`
	Refs     map[string]SymbolRef `json:"refs,omitempty"`
}

func NewMap(name *string) *RMap {
	return &RMap{
		Name:  *name,
		Nodes: map[string]*Node{},
	}
}

func (m *RMap) Save(path string) error {
	if err := internal.MarshalAndWriteToFile(m, path); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
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

func newNode(name, path string) *Node {
	return &Node{
		Name:       name,
		RootFolder: newFolder(path),
	}
}
