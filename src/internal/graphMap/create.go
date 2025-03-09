package graphMap

import (
	"fmt"
	"log"
	"os"

	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/prompt"
	"github.com/JoachimTislov/RefViz/internal/types"
	"github.com/JoachimTislov/RefViz/internal/utils"
)

func Create(name *string) error {

	if *name == "" {
		log.Fatal("Please provide a map name")
	}

	mapPath := path.GetMap(*name)
	act := "created"
	if utils.Exists(mapPath) {
		if !prompt.Confirm(fmt.Sprintf("Map: %s already exists", *name)) {
			return nil
		}
		act = "overwritten"
	}
	if _, err := os.Create(mapPath); err != nil {
		return fmt.Errorf("error creating map: %v", err)
	}
	m := types.NewMap(name)
	if err := m.Save(mapPath); err != nil {
		return err
	}
	log.Printf("Map %s %s\n", *name, act)
	return nil
}
