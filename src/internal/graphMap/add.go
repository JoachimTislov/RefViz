package graphMap

import (
	"fmt"
	"io/fs"
	"log"
	"maps"
	"os"
	"path/filepath"

	c "github.com/JoachimTislov/RefViz/content"
	"github.com/JoachimTislov/RefViz/content/symbol"
	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/prompt"
	"github.com/JoachimTislov/RefViz/internal/types"
)

func AddNode(mapName, nodeName *string) error {

	if *mapName == "" || *nodeName == "" {
		log.Fatal("Please provide a map name and a node name")
	}

	rMap, err := Load(mapName)
	if err != nil {
		return err
	}
	rMap.AddNode(nodeName, path.Project())
	if err := rMap.Save(path.GetMap(rMap.Name)); err != nil {
		return fmt.Errorf("error saving map: %v", err)
	}
	log.Printf("Node: %s added to map: %s\n", *nodeName, *mapName)
	return nil
}

func AddContent(mapName, content, nodeName *string, forceScan, forceUpdate, ask *bool) error {

	if *mapName == "" || *content == "" {
		log.Fatal("Please provide a map name and content to add")
	}

	rMap, err := Load(mapName)
	if err != nil {
		return err
	}

	if err := determineNodeName(nodeName, rMap.Nodes, mapName); err != nil {
		return fmt.Errorf("error determining node name: %v", err)
	}
	node, err := rMap.GetOrCreateNode(nodeName, path.Project())
	if err != nil {
		return fmt.Errorf("error getting or creating node: %v", err)
	}

	paths, err := c.Find(content, *ask)
	if err != nil {
		return err
	}

	for _, p := range paths {
		if err := addPath(p, node.RootFolder, forceScan, forceUpdate); err != nil {
			return fmt.Errorf("error adding path: %v", err)
		}
	}

	// TODO: run in go routine and wait for symbol request and the go routine to finish
	// OR move getSymbol to different package

	for _, node := range rMap.Nodes {
		if err := node.RootFolder.CreateMissingSymbols(); err != nil {
			return fmt.Errorf("error creating missing symbols: %v", err)
		}
	}

	if err := rMap.Save(path.GetMap(rMap.Name)); err != nil {
		return fmt.Errorf("error saving map: %v", err)
	}

	log.Printf("Content: %s added to node: %s in map: %s\n", *content, *nodeName, *mapName)

	return nil
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
			prompt := prompt.SelectPrompt("Select node to add content to", maps.Keys(nodes))
			_, name, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("error selecting node: %v", err)
			}
			*nodeName = name
		}
	}
	return nil
}

func addPath(p string, rootFolder *types.Folder, forceScan, forceUpdate *bool) error {
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
		// honestly, just stupid, but it works
		// result of making a function specific for one case....
		// TODO: This can be removed, its just a precaution to confirm that file is scanned
		c.Get(p, *forceScan, nil)()

		cacheEntry, _, err := symbol.GetMany(p, false)
		if err != nil {
			return fmt.Errorf("error getting symbols: %v", err)
		}

		if err := addFileToFolder(cacheEntry.Symbols, p, rootFolder, forceScan, forceUpdate); err != nil {
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

func addFileToFolder(symbols map[string]*types.Symbol, absPath string, root *types.Folder, forceScan, forceUpdate *bool) error {
	projectPath := path.Project()
	folder, err := root.GetRelatedFolder(absPath, projectPath)
	if err != nil {
		return fmt.Errorf("error updating to related folder: %v", err)
	}

	// Add child symbols to the map as well
	for parentName, symbol := range symbols {
		for name, childSymbol := range symbol.ChildSymbols {

			childFromCache := cache.GetSymbol(childSymbol.Key, name)
			if symbol == nil {
				// Should not happen
				panic(fmt.Sprintf("Child symbol: %s not found in cache, key: %s", name, childSymbol.Key))
			}

			childFromCache.Clean(parentName)

			symbols := make(map[string]*types.Symbol)
			symbols[symbol.Name] = childFromCache

			addFileToFolder(symbols, childFromCache.FilePath, root, forceScan, forceUpdate)
		}
	}

	folderPath, fileName, err := getFolderPathAndFileName(absPath, projectPath)
	if err != nil {
		return fmt.Errorf("error getting folder path and file name: %v", err)
	}
	file := folder.GetFile(fileName, folderPath)
	fullFolderPath := filepath.Join(projectPath, folderPath)
	file.AddSymbols(&folder.Refs, &symbols, &fullFolderPath, &fileName, forceUpdate)
	folder.AddFile(file, forceUpdate)
	return nil
}

func getFolderPathAndFileName(absPath, projectPath string) (string, string, error) {
	relPath, err := filepath.Rel(projectPath, absPath)
	if err != nil {
		return "", "", fmt.Errorf("error getting relative path: %s, err: %v", absPath, err)
	}
	return filepath.Dir(relPath), filepath.Base(relPath), nil
}
