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

	paths, err := c.Find(content, ask)
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

	rMap.CreateMissingSymbols(path.Project())

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
		if err := addFileToFolder(p, rootFolder, forceScan, forceUpdate); err != nil {
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

func addFileToFolder(absPath string, folder *types.Folder, forceScan, forceUpdate *bool) error {

	folder, err := folder.GetRelatedFolder(absPath, path.Project())
	if err != nil {
		return fmt.Errorf("error updating to related folder: %v", err)
	}
	// honestly, just stupid, but it works
	// result of making a function specific for one case....
	c.Get(absPath, *forceScan, nil)()

	cacheEntry, _, err := symbol.GetMany(absPath, false)
	if err != nil {
		return fmt.Errorf("error getting symbols: %v", err)
	}
	////

	folderPath, fileName, err := getFolderPathAndFileName(absPath)
	if err != nil {
		return fmt.Errorf("error getting folder path and file name: %v", err)
	}
	file := folder.GetFile(&fileName, &folderPath)
	fullFolderPath := filepath.Join(path.Project(), folderPath)
	file.AddSymbols(&folder.Refs, &cacheEntry.Symbols, &fullFolderPath, &fileName, forceUpdate)
	folder.AddFile(file, forceUpdate)
	return nil
}

func getFolderPathAndFileName(absPath string) (string, string, error) {
	relPath, err := filepath.Rel(path.Project(), absPath)
	if err != nil {
		return "", "", fmt.Errorf("error getting relative path: %s, err: %v", absPath, err)
	}
	return filepath.Dir(relPath), filepath.Base(relPath), nil
}
