package ops

import (
	"fmt"
	"io/fs"
	"log"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/JoachimTislov/RefViz/types"
)

func CheckMapOps(lm, ln, create, add, delete *bool, mapName *string, nodeName *string, content *string, force, ask *bool) {

	operations := types.Operation{
		{Condition: *lm, Action: listMaps, Msg: "error listing maps"},
		{Condition: *ln, Action: func() error { return listNodes(mapName) }, Msg: "error listing nodes"},
		{Condition: *create, Action: func() error { return createMap(mapName) }, Msg: "error creating map"},
		{Condition: *delete, Action: func() error { return deleteMap(mapName) }, Msg: "error deleting map"},
		{Condition: *add, Action: func() error { return addContentToMap(mapName, content, nodeName, force, ask) }, Msg: "error adding content to map"},
		{Condition: *nodeName != "" && !*add, Action: func() error { return addNodeToMap(mapName, nodeName) }, Msg: "error adding node to map"},
	}

	for _, op := range operations {
		if op.Condition {
			if err := op.Action(); err != nil {
				log.Fatalf("%s: %v\n", op.Msg, err)
			}
		}
	}
}

func createMap(name *string) error {

	if *name == "" {
		log.Fatal("Please provide a map name")
	}

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
	log.Printf("Map %s %s\n", *name, act)
	return nil
}

func addContentToMap(mapName, content, nodeName *string, force, ask *bool) error {

	if *mapName == "" || *content == "" {
		log.Fatal("Please provide a map name and content to add")
	}

	rMap, err := LoadMap(mapName)
	if err != nil {
		return err
	}

	if err := determineNodeName(nodeName, rMap.Nodes, mapName); err != nil {
		return fmt.Errorf("error determining node name: %v", err)
	}
	node, err := rMap.GetOrCreateNode(nodeName, projectPath())
	if err != nil {
		return fmt.Errorf("error getting or creating node: %v", err)
	}

	paths, err := findContent(content, ask)
	if err != nil {
		return err
	}

	for _, p := range paths {
		if err := addPath(p, node.RootFolder, force); err != nil {
			return fmt.Errorf("error adding path: %v", err)
		}
	}

	if err := marshalAndWriteToFile(rMap, getMapPath(rMap.Name)); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

func addPath(p string, rootFolder *types.Folder, force *bool) error {
	e, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("error analyzing path: %s, err: %v", p, err)
	}

	var subPaths []string
	if e.IsDir() {
		if err := retrieveContentInDir(p, &subPaths); err != nil {
			return fmt.Errorf("error retrieving content in directory: %v", err)
		}
	} else {
		subPaths = append(subPaths, p)
	}

	for _, p := range subPaths {
		if err := addFileToFolder(p, rootFolder, force); err != nil {
			return fmt.Errorf("error adding file to folder: %v", err)
		}
	}
	return nil
}

func retrieveContentInDir(p string, paths *[]string) error {
	return filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through directory: %v", err)
		}
		if !d.IsDir() {
			*paths = append(*paths, path)
		}
		return nil
	})
}

func determineNodeName(nodeName *string, nodes map[string]*types.Node, mapName *string) error {
	if *nodeName == "" {
		l := len(nodes)
		switch l {
		case 0:
			*nodeName = "node_" + *mapName
			break
		case 1:
			for n := range nodes {
				*nodeName = n
				break
			}
		default:
			prompt := selectPrompt("Select node to add content to", maps.Keys(nodes))
			_, name, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("error selecting node: %v", err)
			}
			*nodeName = name
		}
	}
	return nil
}

func addFileToFolder(absPath string, folder *types.Folder, force *bool) error {

	folder, err := folder.GetRelatedFolder(absPath, projectPath())
	if err != nil {
		return fmt.Errorf("error updating to related folder: %v", err)
	}
	// honestly, just stupid, but it works
	// result of making a function specific for one case....
	getContent(absPath, *force, nil)()

	cacheEntry, _, err := getSymbols(absPath, false)
	if err != nil {
		return fmt.Errorf("error getting symbols: %v", err)
	}
	////

	folderPath, fileName, err := getFolderPathAndFileName(absPath)
	if err != nil {
		return fmt.Errorf("error getting folder path and file name: %v", err)
	}
	file := folder.GetFile(&fileName, &folderPath)
	file.AddSymbols(&folder.Refs, &cacheEntry.Symbols, &folderPath, &fileName)
	folder.AddFile(file, *force)
	return nil
}

func addNodeToMap(mapName, nodeName *string) error {

	if *mapName == "" || *nodeName == "" {
		log.Fatal("Please provide a map name and a node name")
	}

	rMap, err := LoadMap(mapName)
	if err != nil {
		return err
	}
	rMap.AddNode(nodeName, projectPath())

	if err := marshalAndWriteToFile(rMap, getMapPath(rMap.Name)); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

func LoadMap(name *string) (*types.RMap, error) {
	rMap := types.NewMap(name)
	path := getMapPath(*name)
	if err := getFile(path, rMap); err != nil {
		return nil, fmt.Errorf("error loading map from file with path: %s, err: %v", err, path)
	}
	return rMap, nil
}

func deleteMap(name *string) error {

	if *name == "" {
		log.Fatal("Please provide a map name")
	}

	if confirm(fmt.Sprintf("You are about to delete map %s", *name)) {
		if err := os.Remove(getMapPath(*name)); err != nil {
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

func listMaps() error {
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

func listNodes(maps ...*string) error {
	if len(maps) == 0 || *maps[0] == "" {
		allMaps, err := getMaps()
		if err != nil {
			return fmt.Errorf("error getting maps: %v", err)
		}
		maps = allMaps
	}
	for _, m := range maps {
		rMap, err := LoadMap(m)
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
