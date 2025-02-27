package mappers

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/JoachimTislov/RefViz/internal"
	"github.com/JoachimTislov/RefViz/ops"
)

func CreateGraphvizFile(mapName *string) error {
	m, err := ops.LoadMap(mapName)
	if err != nil {
		return fmt.Errorf("error loading map: %v", err)
	}

	file, err := createDotFile(mapName)
	if err != nil {
		return fmt.Errorf("error creating dot file: %v", err)
	}
	defer file.Close()

	t, err := createTemplate(mapName)
	if err != nil {
		return fmt.Errorf("error creating template: %v", err)
	}

	err = t.Execute(file, m)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}
	return nil
}

func createDotFile(mapName *string) (*os.File, error) {
	file, err := os.Create(internal.DotFilePath(mapName))
	if err != nil {
		return nil, fmt.Errorf("error creating dot file: %v", err)
	}
	return file, nil
}

func createTemplate(mapName *string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"replace":   strings.ReplaceAll,
		"trimSpace": strings.TrimSpace,
		"arr": func(els ...any) any { // https://dev.to/moniquelive/passing-multiple-arguments-to-golang-templates-16h8
			return els
		},
		"debug": func(msg any) error {
			fmt.Println(msg)
			return nil
		},
	}
	return template.New(*mapName).Funcs(funcMap).Parse(tmpl)
}

// https://golang.org/pkg/text/template/
// recursive template with nested definitions
// Whitespace control: https://golang.org/pkg/text/template/#hdr-Text_and_spaces, its a bit tricky
const tmpl = `digraph {{.Name}} {
{{- range $name, $node := .Nodes}}
	subgraph {{$name}} {
		rankdir=TB;
		{{- template "graph" $node.RootFolder }}
	}
{{- end}}
	{{- define "refs"}}
		{{- $symbolRefs := index . 0}}
		{{- $folderName := index . 1}}
		{{- range $symbolRef := $symbolRefs}}
				{{$folderName}}_{{trimSpace $symbolRef.Definition.Name}} -> {{$symbolRef.Ref.FolderName}}_{{trimSpace $symbolRef.Ref.MethodName}};
		{{- end}}
	{{- end}}
{{- define "subgraph" -}}
	{{- if .SubFolders }}
		{{- range $d, $subfolder := .SubFolders -}}
			{{- template "graph" $subfolder }}
		{{- end}}
	{{- end}}
{{- end}}
{{- define "graph" -}}
		{{- $folderName := .FolderName}} {{- /* to avoid issues in files loop */ -}}
		{{- if .Files}}		
		subgraph cluster_{{replace .FolderName "-" "_"}} {
			label = "{{.FolderName}} (folder)";
			rankdir=TB;
			{{- range $file := .Files}}
			subgraph cluster_{{replace (replace $file.Name "." "_") "-" "_"}} {
				label = "{{$file.Name}}";
				labelloc="t";
				rankdir=TB;
				{{- range $symbol := $file.Symbols}}
				{{$folderName}}_{{trimSpace $symbol.Name}} [label = "{{trimSpace $symbol.Name}}, {{$symbol.Kind}}";shape = box;];
					{{- template "refs" (arr $symbol.Refs $folderName) -}}
				{{- end}}
			}
			{{- template "refs" (arr $file.Refs $folderName) -}}
			{{- end}}

			{{- if .SubFolders }}
				{{- template "subgraph" . -}}
			{{- end}}
		}
		{{- if .Refs }}
			{{- template "refs" (arr .Refs .FolderName) -}}
		{{- end}}
		{{- else}}
			{{- if .SubFolders }}
				{{- template "subgraph" . -}}
			{{- end}}
		{{- end}}
{{- end}}
}`
