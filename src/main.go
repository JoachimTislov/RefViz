package main

import (
	"flag"
	"fmt"
	"log"

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
	scanAgain := flag.Bool("f", false, "force scan, ignore cache")
	flag.Parse()

	// Determine if map operations are to be performed
	ops.CheckMapOps(lm, ln, create, add, delete, mapName, nodeName, content)

	if *scan {
		if err := ops.Scan(content, scanAgain); err != nil {
			log.Fatalf("Error scanning content: %v\n", err)
		}
	}
	if *graphviz != "" {
		// Following can be written with any graphing library
		// Currently, the graph is visualized with graphviz
		// Extension: tintinweb.graphviz-interactive-preview, can display the graph in vscode
		if err := mappers.CreateGraphvizFile(graphviz); err != nil {
			log.Fatalf("error creating graphviz map: %v", err)
		}
	}
}
