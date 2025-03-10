package types

import (
	"fmt"
)

type SymbolRef struct {
	Definition definition `json:"definition"`
	Ref        Ref        `json:"reference"`
}

func (s *SymbolRef) createSymbolMapKey() string {
	return fmt.Sprintf("%s:%s_%s:%s", s.Definition.FilePath, s.Definition.Name, s.Ref.FilePath, s.Ref.MethodName)
}
