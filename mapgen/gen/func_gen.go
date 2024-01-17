package gen

import (
	"bytes"
	"fmt"
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
		{{ . }}
	{{- end }}

	return {{ if .OutIsPtr }}&{{ end }}result
}`

type MapTemplateData struct {
	FuncName string
	InType   string
	InIsPtr  bool
	InputVar string
	OutType  string
	OutIsPtr bool
	Resolver *FileImports
	Mappings []string
	Maps     []MapExpression
}

type MapExpression interface {
	String(imports *FileImports) string
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
	data.Mappings = generateExpressions(data)

	t := template.Must(template.New("map").Parse(mapTemplate))
	funcSource := bytes.NewBuffer(nil)
	err := t.Execute(funcSource, data)
	if err != nil {
		panic(err)
	}

	fmt.Println(funcSource.String())
	file, err := parser.ParseFile(fset, "mapgenie_temp.go", "package main\n"+funcSource.String(), 0)
	if err != nil {
		panic(err)
	}
	fset.RemoveFile(fset.File(file.Pos()))

	return file.Decls[0].(*ast.FuncDecl)
}

func generateExpressions(data MapTemplateData) []string {
	results := make([]string, len(data.Maps))
	for i, m := range data.Maps {
		results[i] = m.String(data.Resolver)
	}

	return results
}
