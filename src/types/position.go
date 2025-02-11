package types

import "fmt"

type Position struct {
	Line      string `json:"line"`
	CharRange string `json:"charRange"`
}

func (p Position) GetPos() string {
	return fmt.Sprintf("%s:%s", p.Line, p.CharRange)
}
