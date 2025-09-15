package graphMap

import (
	"fmt"
	"log"

	"github.com/JoachimTislov/refviz/internal"
	"github.com/JoachimTislov/refviz/internal/path"
	"github.com/JoachimTislov/refviz/internal/types"
	"github.com/JoachimTislov/refviz/internal/utils"
)

func Load(name *string) (*types.RMap, error) {
	rMap := types.NewMap(name)
	path := path.GetMap(*name)
	if !utils.Exists(path) {
		log.Fatalf("Map: %s does not exist", *name)
	}
	if err := internal.GetAndUnmarshalFile(path, rMap); err != nil {
		return nil, fmt.Errorf("error loading map from file with path: %s, err: %v", err, path)
	}
	return rMap, nil
}
