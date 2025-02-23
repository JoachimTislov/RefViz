package ops

import (
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/RefViz/types"
)

func CreateMap(name *string) error {
	mapPath := getMapPath(*name)
	act := "created"
	if exists(mapPath) {
		if !confirm(fmt.Sprintf("Map: %s already exists", *name)) {
			return nil
		}
		act = "overwritten"
	}
	if _, err := os.Create(mapPath); err != nil {
		return fmt.Errorf("error creating map: %v", err)
	}
	if err := marshalAndWriteToFile(types.NewMap(name), mapPath); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	fmt.Printf("Map %s %s\n", *name, act)
	return nil
}

func addContentToMap(mapName, content, nodeName *string) error {
	rMap, err := loadMap(mapName)
	if err != nil {
		return err
	}
	if *nodeName == "" {
		if len(rMap.Nodes) == 0 {
			fmt.Println("No nodes found in map")
		} else if len(rMap.Nodes) == 1 {
			for n := range rMap.Nodes {
				*nodeName = n
				break
			}
		} else {
			prompt := selectPrompt("Select node to add content to", maps.Keys(rMap.Nodes))
			_, *nodeName, err = prompt.Run()
			if err != nil {
				return fmt.Errorf("error selecting node: %v", err)
			}
		}
	}
	node := rMap.Nodes[*nodeName]
	if node == nil {
		node = types.NewNode(*nodeName, filepath.Base(projectPath()))
	}
	paths, err := findContent(content)
	if err != nil {
		return err
	}
	for _, p := range paths {
		e, valid := checkIfValid(p)
		if !valid {
			return fmt.Errorf("error: %s is not a valid entity, err: %v", p, err)
		}
		relPath, err := filepath.Rel(projectPath(), p)
		if err != nil {
			return fmt.Errorf("error getting relative path: %s, err: %v", p, err)
		}
		folder := getRelatedFolder(relPath, node.RootFolder)
		if e.IsDir() {
			currentRelPath := relPath
			if err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return fmt.Errorf("error walking through directory: %v", err)
				}

				if d.IsDir() {
					relPath, err := filepath.Rel(projectPath(), path)
					if err != nil || relPath == "." {
						return fmt.Errorf("error getting relative path: %s, err: %v", path, err)
					}
					if relPath != currentRelPath {
						folder = getRelatedFolder(relPath, node.RootFolder)
						currentRelPath = relPath
					}
					folder.MoveToSubFolder(d.Name(), relPath)
				} else {
					addFileToFolder(path, relPath, folder)
				}
				return nil
			}); err != nil {
				return fmt.Errorf("error walking through directory: %v", err)
			}
		} else {
			addFileToFolder(p, relPath, folder)
		}
	}
	return nil
}

func addFileToFolder(absPath, relPath string, folder *types.Folder) error {
	cacheEntry, _, err := getSymbols(absPath, false)
	if err != nil {
		return fmt.Errorf("error getting symbols: %v", err)
	}
	for _, s := range cacheEntry.Symbols {
		folderPath := filepath.Dir(relPath)
		fileName := filepath.Base(relPath)
		file := types.NewFile(fileName, folderPath)

		sortRefsIntoHierarchy(s, folder, &file, folderPath, fileName)

		file.Symbols = append(file.Symbols, *s)
		folder.Files = append(folder.Files, file)
	}
	return nil
}

// sortRefsIntoHierarchy loops through the references of a symbol and moves them to the appropriate folder or file
// based on the path of the reference
// refs are removed from the symbol after they are moved
func sortRefsIntoHierarchy(s *types.Symbol, folder *types.Folder, file *types.File, folderPath, fileName string) {
	refsToRemove := []string{}
	for key, r := range s.Refs {
		sRef := types.SymbolRef{
			Definition: *s,
			Ref:        *r,
		}
		folderPageDiffer := filepath.Dir(r.Path) != folderPath
		fileNameDiffer := filepath.Base(r.Path) != fileName

		if folderPageDiffer || fileNameDiffer {
			refsToRemove = append(refsToRemove, key)

			if folderPageDiffer {
				folder.Refs = append(folder.Refs, sRef)
			}
			if fileNameDiffer {
				file.Refs = append(file.Refs, sRef)
			}
		}
	}
	for _, key := range refsToRemove {
		delete(s.Refs, key)
	}
}

// getRelatedFolder returns the folder related to the path
// recursively keys through the recursive folder data structure
// to find the folder related to the path
func getRelatedFolder(relPath string, folder *types.Folder) *types.Folder {
	dirs := strings.Split(filepath.Dir(relPath), string(filepath.Separator))
	for _, d := range dirs {
		if d != "" {
			folder.MoveToSubFolder(d, relPath)
		}
	}
	return folder
}

func addNodeToMap(mapName, nodeName *string) error {
	rMap, err := loadMap(mapName)
	if err != nil {
		return err
	}
	rMap.Nodes[*nodeName] = types.NewNode(*nodeName, filepath.Base(projectPath()))

	if err := marshalAndWriteToFile(rMap, getMapPath(*mapName)); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

func loadMap(name *string) (*types.RMap, error) {
	rMap := types.NewMap(name)
	path := getMapPath(*name)
	if err := getFile(path, rMap); err != nil {
		return nil, fmt.Errorf("error loading map from file with path: %s, err: %v", err, path)
	}
	return rMap, nil
}

func DeleteMap(name *string) error {
	if confirm(fmt.Sprintf("You are about to delete map %s", *name)) {
		if err := os.Remove(getMapPath(*name)); err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Map: %s does not exist\n", *name)
				return nil
			} else {
				return fmt.Errorf("error deleting map: %v", err)
			}
		}
		fmt.Printf("Map %s deleted\n", *name)
	} else {
		fmt.Printf("Map %s not deleted\n", *name)
	}
	return nil
}

func getMaps() ([]*string, error) {
	maps, err := os.ReadDir(mapPath())
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}
	var mapNames []*string
	for _, m := range maps {
		mapNames = append(mapNames, &strings.Split(m.Name(), ".")[0])
	}
	return mapNames, nil
}

func ListMaps() error {
	maps, err := getMaps()
	if err != nil {
		return fmt.Errorf("error getting maps: %v", err)
	}
	if len(maps) == 0 {
		fmt.Println("No maps found")
		return nil
	}
	fmt.Println("Your maps:")
	for _, m := range maps {
		fmt.Printf("\t%s\n", *m)
	}
	return nil
}

func ListNodes(maps ...*string) error {
	if len(maps) == 0 {
		allMaps, err := getMaps()
		if err != nil {
			return fmt.Errorf("error getting maps: %v", err)
		}
		maps = allMaps
	}
	for _, m := range maps {
		rMap, err := loadMap(m)
		if err != nil {
			return err
		}
		fmt.Printf("Nodes in map: %s\n", *m)
		for n := range rMap.Nodes {
			fmt.Printf("\t%s\n", n)
		}
	}
	return nil
}
