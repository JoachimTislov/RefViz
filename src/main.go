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
	"github.com/JoachimTislov/RefViz/internal/ops"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/mappers"
)

func init() {
	if err := load.Defs(); err != nil {
		log.Fatal(err)
	}
}

const entities = "map or a node"

type flags struct {
	graphviz                   *string
	lm                         *bool
	ln                         *bool
	create                     *bool
	delete                     *bool
	mapName                    *string
	nodeName                   *string
	scan                       *bool
	findSemiOrphans            *bool
	unusedSymbol               *bool
	add                        *bool
	content                    *string
	display                    *bool
	forceScan                  *bool
	forceUpdate                *bool
	ask                        *bool
	exDir                      *string
	exFile                     *string
	inExt                      *string
	baseLinkToBranch           *string
	exSymbol                   *string
	exFindRefsForSymbols       *string
	exFindRefsForSymbolsPrefix *string
}

func handleFlags() flags {
	f := flags{
		graphviz:                   flag.String("graphviz", "", "generate graphviz file with the given map"),
		lm:                         flag.Bool("lm", false, "list maps"),
		ln:                         flag.Bool("ln", false, "list nodes"),
		create:                     flag.Bool("c", false, fmt.Sprintf("create (%s)", entities)),
		delete:                     flag.Bool("d", false, fmt.Sprintf("delete (%s)", entities)),
		mapName:                    flag.String("m", "", "map name"),
		nodeName:                   flag.String("n", "", "node name"),
		scan:                       flag.Bool("scan", false, "scan the project for content"),
		findSemiOrphans:            flag.Bool("scan.Orphans", false, "find semi-orphans"),
		unusedSymbol:               flag.Bool("scan.Unused", false, "find unused symbols"),
		add:                        flag.Bool("add", false, "add content to map"),
		content:                    flag.String("content", "", "content to scan, file or folder"),
		display:                    flag.Bool("display", false, "display the map"),
		forceScan:                  flag.Bool("fs", false, "force scan, ignores cache"),
		forceUpdate:                flag.Bool("fu", false, "force update map content"),
		ask:                        flag.Bool("a", false, "select content to add to map"),
		exDir:                      flag.String("ex.dir", "", "exclude directory"),
		exFile:                     flag.String("ex.file", "", "exclude file"),
		inExt:                      flag.String("in.ext", "", "include extension"),
		baseLinkToBranch:           flag.String("baseLink", "", "base link to branch of a github repository"),
		exSymbol:                   flag.String("ex.symbol", "", "exclude symbol"),
		exFindRefsForSymbols:       flag.String("ex.RefsForSymbols.Name", "", "exclude find refs for symbols"),
		exFindRefsForSymbolsPrefix: flag.String("ex.RefsForSymbols.Prefix", "", "exclude find refs for symbols prefix"),
	}
	flag.Parse()
	return f
}

func applyConfig(f flags) {
	if err := config.Exclude(f.exDir, f.exFile, f.inExt, f.exSymbol, f.exFindRefsForSymbols, f.exFindRefsForSymbolsPrefix); err != nil {
		log.Fatalf("Error processing exclusions: %v", err)
	}
	if err := config.SetBaseBranchLink(*f.baseLinkToBranch); err != nil {
		log.Fatalf("Error setting base branch link: %v", err)
	}
	config.Save() // Call Save after all config modifications
}

func performActions(f flags) {
	// Determine if map operations are to be performed
	ops.Check(f.lm, f.ln, f.create, f.add, f.delete, f.mapName, f.nodeName, f.content, f.forceScan, f.forceUpdate, f.ask)

	if *f.scan {
		if err := c.Scan(f.content, *f.forceScan, *f.ask); err != nil {
			log.Fatalf("Error scanning content: %v\n", err)
		}
	}
	if *f.findSemiOrphans {
		cache.QuickfeedFindSemiOrphans()
	}
	if *f.unusedSymbol {
		cache.UnusedSymbols()
	}
}

func handleOutput(f flags) {
	// Update the if graphviz flag is set or map name is provided and user wants to display the map
	if *f.mapName != "" && *f.display {
		if _, err := exec.LookPath("xdot"); err != nil {
			log.Println("xdot is not installed or not in PATH. Please install xdot to display graphs. (e.g., sudo apt-get install xdot)")
		} else {
			cmd := exec.Command("xdot", path.DotFile(f.mapName))
			if err := cmd.Start(); err != nil {
				// Log Fatalf here because if LookPath succeeds, xdot should be callable.
				// If cmd.Start() fails, it's likely a more serious issue with xdot execution itself.
				log.Fatalf("Error starting xdot: %v. Please ensure xdot is installed and configured correctly.", err)
			}
		}
	} else if *f.mapName == "" && *f.display {
		log.Fatal("Please provide a map name to display")
	}
	if *f.graphviz != "" {
		// Following can be written with any graphing library
		// Currently, the graph is visualized with graphviz
		// Extension: tintinweb.graphviz-interactive-preview, can display the graph in vscode
		if err := mappers.CreateGraphvizFile(f.graphviz); err != nil {
			log.Fatalf("error creating graphviz map: %v", err)
		}
	}
}

func main() {
	f := handleFlags()
	applyConfig(f)
	performActions(f)
	handleOutput(f)
}
