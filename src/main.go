package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/JoachimTislov/RefViz/mappers"
	"github.com/JoachimTislov/RefViz/ops"
)

/*
TODO: implement libraries which finds references for typescript
*/

func init() {
	if err := ops.LoadDefs(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	graphviz := flag.String("graphviz", "", "generate graphviz file with the given map")
	lm := flag.Bool("lm", false, "list maps")
	ln := flag.Bool("ln", false, "list nodes")
	create := flag.Bool("c", false, fmt.Sprintf("create (%s)", ops.Entities))
	delete := flag.Bool("d", false, fmt.Sprintf("delete (%s)", ops.Entities))
	mapName := flag.String("m", "", "map name")
	nodeName := flag.String("n", "", "node name")
	scan := flag.Bool("scan", false, "scan the project for content")
	add := flag.Bool("add", false, "add content to map")
	content := flag.String("content", "", "content to scan, file or folder")
	display := flag.Bool("display", false, "display the map")
	force := flag.Bool("f", false, "force scans or updates, ignores cache")
	ask := flag.Bool("a", false, "select content to add to map")
	flag.Parse()

	// Determine if map operations are to be performed
	ops.CheckMapOps(lm, ln, create, add, delete, mapName, nodeName, content, force, ask)

	if *scan {
		if err := ops.Scan(content, force, ask); err != nil {
			log.Fatalf("Error scanning content: %v\n", err)
		}
	}
	// Update the if graphviz flag is set or map name is provided and user wants to display the map
	if *mapName != "" && *display {
		graphviz = mapName
	} else if *mapName == "" && *display {
		log.Fatal("Please provide a map name to display")
	}
	if *graphviz != "" {
		// Following can be written with any graphing library
		// Currently, the graph is visualized with graphviz
		// Extension: tintinweb.graphviz-interactive-preview, can display the graph in vscode
		if err := mappers.CreateGraphvizFile(graphviz); err != nil {
			log.Fatalf("error creating graphviz map: %v", err)
		}

		cmd := exec.Command("xdot", ops.DotFilePath(graphviz))
		if err := cmd.Start(); err != nil {
			log.Fatalf("Please install xdot with: sudo apt-get install xdot, its used to display the graph")
		}
	}
}
