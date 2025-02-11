package main

import (
	"flag"
	"fmt"

	"github.com/JoachimTislov/Project-Visualizer/mappers"
	op "github.com/JoachimTislov/Project-Visualizer/operations"
)

/*
TODO: Make the error handling return the gopls log instead of the error status message..
TODO: implement libraries which finds references for typescript and react .tsx and .ts files
*/

func main() {
	graphviz := flag.String("graphviz", "", "generate graphviz file with the given map")
	list := flag.Bool("list", false, "list all maps")
	create := flag.String("create", "", "create map")
	delete := flag.String("delete", "", "delete map")
	scan := flag.Bool("scan", false, "scan the project for symbols")
	findRefs := flag.Bool("references", false, "when scanning, also find references")
	content := flag.String("content", "", "name of file or folder to scan, default is everything")
	flag.Parse()
	if *scan {
		var isDir bool
		if err := op.HandleContentInput(&isDir, content); err != nil {
			fmt.Printf("Error handling content input: %v\n", err)
			return
		}
		if err := op.Scan(&isDir, findRefs, content); err != nil {
			fmt.Printf("Error scanning content: %v\n", err)
			return
		}
	}
	if *list {
		if err := op.ListMaps(); err != nil {
			fmt.Println(err)
			return
		}
	}
	if *create != "" {
		if err := op.CreateMap(create); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Map created")
		return
	}
	if *delete != "" {
		if err := op.DeleteMap(delete); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Map deleted")
		return
	}
	if *graphviz != "" {
		// Following can be written with any graphing library
		// Currently, the graph is visualized with graphviz
		// Extension: tintinweb.graphviz-interactive-preview, can display the graph in vscode
		if err := mappers.CreateGraphvizFile(graphviz); err != nil {
			fmt.Println(err)
			return
		}
	}
}
