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

func CheckMapOps(lm, ln, create, add, delete *bool, mapName *string, nodeName *string, content *string) {

	operations := types.Operation{
		{Condition: *lm, Action: listMaps, Msg: "error listing maps"},
		{Condition: *ln, Action: func() error { return listNodes(mapName) }, Msg: "error listing nodes"},
		{Condition: *create, Action: func() error { return createMap(mapName) }, Msg: "error creating map"},
		{Condition: *delete, Action: func() error { return deleteMap(mapName) }, Msg: "error deleting map"},
		{Condition: *add, Action: func() error { return addContentToMap(mapName, content, nodeName) }, Msg: "error adding content to map"},
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

func addContentToMap(mapName, content, nodeName *string) error {
	projectPath := projectPath()

	if *mapName == "" || *content == "" {
		log.Fatal("Please provide a map name and content to add")
	}

	rMap, err := LoadMap(mapName)
	if err != nil {
		return err
	}

	if err := determineNodeName(nodeName, rMap.Nodes); err != nil {
		return fmt.Errorf("error determining node name: %v", err)
	}
	node, ok := rMap.Nodes[*nodeName]
	if !ok {
		node = types.NewNode(*nodeName, projectPath)
		rMap.Nodes[*nodeName] = node
	}

	paths, err := findContent(content)
	if err != nil {
		return err
	}

	for _, p := range paths {
		if err := addPath(p, projectPath, node.RootFolder); err != nil {
			return fmt.Errorf("error adding path: %v", err)
		}
	}

	if err := marshalAndWriteToFile(rMap, getMapPath(rMap.Name)); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

func addPath(p, projectPath string, rootFolder *types.Folder) error {
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
		relPath, err := filepath.Rel(projectPath, p)
		if err != nil {
			return fmt.Errorf("error getting relative path: %s, err: %v", p, err)
		}
		if err := addFileToFolder(p, relPath, rootFolder); err != nil {
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

func determineNodeName(nodeName *string, nodes map[string]*types.Node) error {
	if *nodeName == "" {
		l := len(nodes)
		switch l {
		case 0:
			log.Println("Zero nodes found in map")
			log.Println("Please provide a node name")
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

func addFileToFolder(absPath, relPath string, folder *types.Folder) error {
	folder = getRelatedFolder(relPath, folder)

	cacheEntry, _, err := getSymbols(absPath, false)
	if err != nil {
		return fmt.Errorf("error getting symbols: %v", err)
	}

	folderName, fileName := filepath.Dir(relPath), filepath.Base(relPath)
	file := types.NewFile(fileName, folderName)
	for _, s := range cacheEntry.Symbols {
		file.AddSymbol(sortRefsIntoHierarchy(*s, &folder.Refs, &file.Refs, folderName, fileName))
	}
	folder.AddFile(file)
	return nil
}

// sortRefsIntoHierarchy loops through the references of a symbol and moves them to the appropriate folder or file
// based on the path of the reference
// refs are removed from the symbol after they are moved
func sortRefsIntoHierarchy(s types.Symbol, folderRefs, fileRefs *[]types.SymbolRef, folderPath, fileName string) types.Symbol {
	refsToRemove := []string{}
	for key, r := range s.Refs {
		sRef := types.SymbolRef{
			Definition: s,
			Ref:        *r,
		}
		folderPageDiffer := filepath.Dir(r.Path) != folderPath
		fileNameDiffer := filepath.Base(r.Path) != fileName

		// if the reference is in a different folder or page, move it to the appropriate folder or file
		// and remove it from the symbol
		if folderPageDiffer || fileNameDiffer {
			refsToRemove = append(refsToRemove, key)

			// if the reference is in a different folder, move it to the folder
			if folderPageDiffer {
				*folderRefs = append(*folderRefs, sRef)
			}
			// if the reference is in a different file, move it to the file
			if fileNameDiffer {
				*fileRefs = append(*fileRefs, sRef)
			}
		}
	}
	for _, key := range refsToRemove {
		delete(s.Refs, key)
	}
	return s
}

// getRelatedFolder returns the folder related to the path
// recursively keys through the recursive folder data structure
// to find the folder related to the path
func getRelatedFolder(relPath string, folder *types.Folder) *types.Folder {
	if folder == nil {
		folder = types.NewFolder(relPath)
	}

	dirs := strings.Split(filepath.Dir(relPath), string(filepath.Separator))
	projectPath := projectPath()
	for _, d := range dirs {
		projectPath = filepath.Join(projectPath, d)

		if folder.SubFolders == nil {
			folder.SubFolders = make(map[string]*types.Folder)
		}

		if _, exists := folder.SubFolders[d]; !exists {
			folder.SubFolders[d] = types.NewFolder(projectPath)
		}
		folder = folder.SubFolders[d]
	}
	return folder
}

func addNodeToMap(mapName, nodeName *string) error {

	if *mapName == "" || *nodeName == "" {
		log.Fatal("Please provide a map name and a node name")
	}

	rMap, err := LoadMap(mapName)
	if err != nil {
		return err
	}
	rMap.Nodes[*nodeName] = types.NewNode(*nodeName, projectPath())

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
