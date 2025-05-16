package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"

	c "github.com/JoachimTislov/RefViz/content"
	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/core/config"
	"github.com/JoachimTislov/RefViz/core/load"
	ast "github.com/JoachimTislov/RefViz/internal/ast/go"
	"github.com/JoachimTislov/RefViz/internal/mappers"
	"github.com/JoachimTislov/RefViz/internal/ops"
	"github.com/JoachimTislov/RefViz/internal/path"
)

func init() {
	if err := load.Defs(); err != nil {
		log.Fatal(err)
	}
}

const entities = "map or a node"

func main() {
	manager := ast.NewManager(path.Quickfeed())
	manager.GetSymbols()

	graphviz := flag.String("graphviz", "", "generate graphviz file with the given map")
	lm := flag.Bool("lm", false, "list maps")
	ln := flag.Bool("ln", false, "list nodes")
	create := flag.Bool("c", false, fmt.Sprintf("create (%s)", entities))
	delete := flag.Bool("d", false, fmt.Sprintf("delete (%s)", entities))
	mapName := flag.String("m", "", "map name")
	nodeName := flag.String("n", "", "node name")
	scan := flag.Bool("scan", false, "scan the project for content")
	findSemiOrphans := flag.Bool("scan.Orphans", false, "find semi-orphans")
	unusedSymbol := flag.Bool("scan.Unused", false, "find unused symbols")
	add := flag.Bool("add", false, "add content to map")
	content := flag.String("content", "", "content to scan, file or folder")
	display := flag.Bool("display", false, "display the map")
	forceScan := flag.Bool("fs", false, "force scan, ignores cache")
	forceUpdate := flag.Bool("fu", false, "force update map content")
	ask := flag.Bool("a", false, "select content to add to map")
	exDir := flag.String("ex.dir", "", "exclude directory")
	exFile := flag.String("ex.file", "", "exclude file")
	inExt := flag.String("in.ext", "", "include extension")
	baseLinkToBranch := flag.String("baseLink", "", "base link to branch of a github repository")
	exSymbol := flag.String("ex.symbol", "", "exclude symbol")
	exFindRefsForSymbols := flag.String("ex.RefsForSymbols.Name", "", "exclude find refs for symbols")
	exFindRefsForSymbolsPrefix := flag.String("ex.RefsForSymbols.Prefix", "", "exclude find refs for symbols prefix")
	flag.Parse()

	config.Exclude(exDir, exFile, inExt, exSymbol, exFindRefsForSymbols, exFindRefsForSymbolsPrefix)
	config.SetBaseBranchLink(*baseLinkToBranch)

	// Determine if map operations are to be performed
	ops.Check(lm, ln, create, add, delete, mapName, nodeName, content, forceScan, forceUpdate, ask)

	if *scan {
		if err := c.Scan(content, *forceScan, *ask); err != nil {
			log.Fatalf("Error scanning content: %v\n", err)
		}
	}
	if *findSemiOrphans {
		cache.QuickfeedFindSemiOrphans()
	}
	if *unusedSymbol {
		cache.UnusedSymbols()
	}
	// Update the if graphviz flag is set or map name is provided and user wants to display the map
	if *mapName != "" && *display {
		cmd := exec.Command("xdot", path.DotFile(mapName))
		if err := cmd.Start(); err != nil {
			log.Fatalf("Please install xdot with: sudo apt-get install xdot, its used to display the graph")
		}
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
	}
}
