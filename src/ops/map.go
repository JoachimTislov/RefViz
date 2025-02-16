package ops

import (
	"fmt"
	"os"
)

func CreateMap(name *string) error {
	if _, err := os.Create(getMapPath(name)); err != nil {
		return fmt.Errorf("error creating map: %v", err)
	}
	return nil
}

func DeleteMap(name *string) error {
	if err := os.Remove(getMapPath(name)); err != nil {
		return fmt.Errorf("error deleting map: %v", err)
	}
	return nil
}

func ListMaps() error {
	maps, err := os.ReadDir(mapPath())
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}
	if len(maps) == 0 {
		fmt.Println("No maps found")
		return nil
	}
	for _, m := range maps {
		fmt.Println(m.Name())
	}
	return nil
}
