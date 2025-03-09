package graphMap

import (
	"fmt"
	"log"
	"os"

	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/prompt"
)

func Delete(name *string) error {

	if *name == "" {
		log.Fatal("Please provide a map name")
	}

	if prompt.Confirm(fmt.Sprintf("You are about to delete map %s", *name)) {
		if err := os.Remove(path.GetMap(*name)); err != nil {
			if os.IsNotExist(err) {
				log.Printf("Map: %s does not exist\n", *name)
				return nil
			} else {
				return fmt.Errorf("error deleting map: %v", err)
			}
		}
		log.Printf("Deleted map: %s \n", *name)
	} else {
		log.Printf("Cancelled deletion of map: %s\n", *name)
	}
	return nil
}
