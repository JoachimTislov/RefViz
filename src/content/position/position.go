package position

import (
	"fmt"
	"strings"

	"github.com/JoachimTislov/RefViz/internal/types"
	"github.com/JoachimTislov/RefViz/internal/utils"
)

// Gets the line and character range position of the symbol
func Create(p string) types.Position {
	args := strings.Split(p, "-")
	args2 := strings.Split(args[0], ":")
	return types.Position{
		Line:      utils.Atoi(args2[0]), // starting line position
		CharRange: fmt.Sprintf("%s-%s", args2[1], strings.Split(args[1], ":")[1]),
	}
}
