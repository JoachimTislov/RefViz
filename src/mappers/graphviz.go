package mappers

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/JoachimTislov/RefViz/ops"
)

func CreateGraphvizFile(mapName *string) error {
	file, err := createDotFile(mapName)
	if err != nil {
		return fmt.Errorf("error creating dot file: %v", err)
	}
	defer file.Close()

	t, err := createTemplate(mapName)
	if err != nil {
		return fmt.Errorf("error creating template: %v", err)
	}

	err = t.Execute(file, mapName)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}
	return nil
}

func createDotFile(mapName *string) (*os.File, error) {
	file, err := os.Create(ops.DotFilePath(mapName))
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
	}
	return template.New(*mapName).Funcs(funcMap).Parse(tmpl)
}

// https://golang.org/pkg/text/template/
// recursive template with nested definitions
// Whitespace control: https://golang.org/pkg/text/template/#hdr-Text_and_spaces, its a bit tricky
const tmpl = `
{{- range $folderName, $folder := .Folder}}
digraph {{$folderName}} {
	rankdir=TB;
	{{- template "subgraph" $folder -}}
}
{{- end}}
	{{- define "refs"}}
		{{- $refs := index . 0}}
		{{- $folderName := index . 1}}
		{{- range $ref := $refs}}
			{{$folderName}}_{{trimSpace $ref.Source.MethodName}} -> {{$folderName}}_{{trimSpace $ref.Info.MethodName}};
		{{- end}}
	{{- end}}
{{- define "subgraph"}}
{{- range $folderName, $subfolder := .SubFolders.Folder}}
	subgraph cluster_{{replace $folderName "-" "_"}} {
		label = "{{$folderName}} (folder)";
		rankdir=TB;
		{{- range $file := $subfolder.Files}}
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
	}
	{{- template "refs" (arr $subfolder.Refs $folderName) -}}
	{{- template "subgraph" $subfolder -}}
{{- end}}
{{- end}}`
