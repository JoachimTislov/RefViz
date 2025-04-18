package ast

import (
	"github.com/JoachimTislov/RefViz/internal/types"
)

type ast interface {
	GetSymbols() []types.Symbol
	GetReferences() []types.Ref
}
