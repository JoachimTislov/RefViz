package main

import (
	"flag"
	"log"

	"github.com/JoachimTislov/Project-Visualizer/mappers"
	"github.com/JoachimTislov/Project-Visualizer/ops"
	"github.com/JoachimTislov/Project-Visualizer/types"
)

/*
TODO: Make the error handling return the gopls log instead of the error status message..
TODO: implement libraries which finds references for typescript
*/

func init() {
	if err := ops.LoadDefs(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	graphviz := flag.String("graphviz", "", "generate graphviz file with the given map")
	list := flag.Bool("list", false, "list all maps")
	create := flag.String("c", "", "create map")
	delete := flag.String("d", "", "delete map")
	scan := flag.String("scan", "scan", "scan the project (just -scan), to scan a file or folder provide a name (e.g. -scan=main.go)")
	findRefs := flag.Bool("refs", false, "include references in the scan")
	flag.Parse()

	checkOps(list, create, delete)
	if *scan != "scan" {
		if err := ops.Scan(findRefs, scan); err != nil {
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

func checkOps(list *bool, create, delete *string) {
	operations := types.Operation{
		{*list, ops.ListMaps, "error listing maps"},
		{*create != "", func() error { return ops.CreateMap(create) }, "error creating map"},
		{*delete != "", func() error { return ops.DeleteMap(delete) }, "error deleting map"},
	}
	for _, op := range operations {
		if op.Condition {
			if err := op.Action(); err != nil {
				log.Fatalf("%s: %v\n", op.Msg, err)
			}
		}
	}
}
