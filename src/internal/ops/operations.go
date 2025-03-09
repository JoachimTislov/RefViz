package ops

import (
	"log"

	"github.com/JoachimTislov/RefViz/internal/graphMap"
)

type operation []struct {
	Condition bool
	Action    func() error
	Msg       string
}

func Check(lm, ln, create, add, delete *bool, mapName *string, nodeName *string, content *string, forceScan, forceUpdate, ask *bool) {

	operations := operation{
		{Condition: *lm, Action: graphMap.ListMaps, Msg: "error listing maps"},
		{Condition: *ln, Action: func() error { return graphMap.ListNodes(mapName) }, Msg: "error listing nodes"},
		{Condition: *create, Action: func() error { return graphMap.Create(mapName) }, Msg: "error creating map"},
		{Condition: *delete, Action: func() error { return graphMap.Delete(mapName) }, Msg: "error deleting map"},
		{Condition: *add, Action: func() error { return graphMap.AddContent(mapName, content, nodeName, forceScan, forceUpdate, ask) }, Msg: "error adding content to map"},
		{Condition: *nodeName != "" && !*add, Action: func() error { return graphMap.AddNode(mapName, nodeName) }, Msg: "error adding node to map"},
	}

	for _, op := range operations {
		if op.Condition {
			if err := op.Action(); err != nil {
				log.Fatalf("%s: %v\n", op.Msg, err)
			}
		}
	}
}
