package types

import "fmt"

type Position struct {
	Line      int    `json:"line"`
	CharRange string `json:"charRange"`
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%s", p.Line, p.CharRange)
}
