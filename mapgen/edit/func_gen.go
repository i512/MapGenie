package edit

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"text/template"
)

const mapTemplate = `func {{ .FuncName }}(input {{ .InTypeArg }}) {{ .OutTypeArg }} {
	var result {{ .OutType }}

	{{- if .InIsPtr }}
	if input == nil {
		return {{ if .OutIsPtr }}&{{ end }}result
	}
	{{ end }}

	{{- range .Mappings }}
		{{- if and .InPtr .OutPtr }}
			if input.{{ .InName }} != nil {
				{{ .OutName }} := 
					{{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
					*input.{{ .InName }}
					{{- if ne .CastWith "" }}){{- end }}

				result.{{ .OutName }} = &{{ .OutName }}
			}
		{{- else if .InPtr }}
			if input.{{ .InName }} != nil {
				result.{{ .OutName }} ={{" "}}
					{{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
					*input.{{ .InName }}
					{{- if ne .CastWith "" }}){{- end }}
			}
		{{- else }}
			{{- if and .OutPtr (ne .CastWith "")}}
				{{ .OutName }} := {{ .CastWith }}(input.{{ .InName }})
				result.{{ .OutName }} = &{{ .OutName }}
			{{- else }}
				result.{{ .OutName }} ={{" "}}
					{{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
					{{- if .OutPtr }}&{{- end }}
					{{- ""}}input.{{ .InName }}
					{{- if ne .CastWith "" }}){{- end }}
		   {{- end }}
		{{- end }}
	{{- end }}

	return {{ if .OutIsPtr }}&{{ end }}result
}`

type TemplateMapping struct {
	InName   string
	InPtr    bool
	OutName  string
	OutPtr   bool
	CastWith string
	Cast     bool
}

type MapTemplateData struct {
	FuncName string
	InType   string
	InIsPtr  bool
	InputVar string
	OutType  string
	OutIsPtr bool
	Mappings []TemplateMapping
}

func (d MapTemplateData) InTypeArg() string {
	if d.InIsPtr {
		return "*" + d.InType
	}

	return d.InType
}

func (d MapTemplateData) OutTypeArg() string {
	if d.OutIsPtr {
		return "*" + d.OutType
	}

	return d.OutType
}

func MapperFuncAst(fset *token.FileSet, data MapTemplateData) *ast.FuncDecl {
	t := template.Must(template.New("map").Parse(mapTemplate))
	funcSource := bytes.NewBuffer(nil)
	err := t.Execute(funcSource, data)
	if err != nil {
		panic(err)
	}

	file, err := parser.ParseFile(fset, "mapgenie_temp.go", "package main\n"+funcSource.String(), 0)
	if err != nil {
		panic(err)
	}
	fset.RemoveFile(fset.File(file.Pos()))

	return file.Decls[0].(*ast.FuncDecl)
}
