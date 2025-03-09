package graphMap

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/JoachimTislov/RefViz/internal/path"
)

func ListMaps() error {
	maps, err := getMaps()
	if err != nil {
		return fmt.Errorf("error getting maps: %v", err)
	}
	if len(maps) == 0 {
		log.Println("No maps found")
		return nil
	}
	log.Println("Your maps:")
	for _, m := range maps {
		log.Printf("\t%s\n", *m)
	}
	return nil
}

func ListNodes(maps ...*string) error {
	if len(maps) == 0 || *maps[0] == "" {
		allMaps, err := getMaps()
		if err != nil {
			return fmt.Errorf("error getting maps: %v", err)
		}
		maps = allMaps
	}
	for _, m := range maps {
		rMap, err := Load(m)
		if err != nil {
			return err
		}
		if len(rMap.Nodes) == 0 {
			log.Printf("Zero nodes found in map: %s\n", *m)
			continue
		}
		log.Printf("Nodes in map: %s\n", *m)
		for n := range rMap.Nodes {
			log.Printf("\t%s\n", n)
		}
	}
	return nil
}

func getMaps() ([]*string, error) {
	maps, err := os.ReadDir(path.Map())
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}
	var mapNames []*string
	for _, m := range maps {
		mapNames = append(mapNames, &strings.Split(m.Name(), ".")[0])
	}
	return mapNames, nil
}
