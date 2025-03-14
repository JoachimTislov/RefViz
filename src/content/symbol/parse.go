package symbol

import (
	"fmt"
	"strings"

	"github.com/JoachimTislov/RefViz/content/position"
	"github.com/JoachimTislov/RefViz/core/config"
	"github.com/JoachimTislov/RefViz/internal/types"
)

const (
	method = "Method"
)

// parses the output of the gopls symbols command and extracts the name, kind, and position of each symbol
func parseSymbols(output, filePath string, s *map[string]*types.Symbol) {
	// retrieve values from last, to handle this specific case: uint64 | string Field 27:2-27:17
	// the usual case is: uint64 Field 27:2-27:17, name, kind, position
	for _, line := range strings.Split(output, "\n") {
		args := strings.Split(line, " ")
		l := len(args)
		if l < 3 {
			continue
		}
		name := strings.TrimSpace(strings.Join(args[:l-2], " ")) // name is everything except the last 2 elements

		// skip excluded symbols
		if config.NotValidSymbol(name) {
			continue
		}

		kind := args[l-2]
		// for methods, remove the receiver type
		// (*Service[K, V]).SendTo Method -> SendTo
		if kind == method && strings.Contains(name, ".") {
			name = strings.Split(name, ".")[1]
		}
		(*s)[name] = &types.Symbol{
			Name:     name,
			Kind:     kind,
			Path:     fmt.Sprintf("%s:%s", filePath, args[l-1]),
			FilePath: filePath,
			Position: position.Create(args[l-1]),
		}
	}
}
