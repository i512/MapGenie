package gen

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"mapgenie/entities"
	"mapgenie/gen/fragments"
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
	Maps     []entities.Statement
}

type MapStatement interface {
	Generate(*fragments.GenerationCtx) (string, error)
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

func MapperFuncAst(ctx context.Context, fset *token.FileSet, data MapTemplateData) (*ast.FuncDecl, error) {
	mappings, err := generateExpressions(ctx, data)
	if err != nil {
		return nil, err
	}

	data.Mappings = mappings

	t := template.Must(template.New("map").Parse(mapTemplate))
	funcSource := bytes.NewBuffer(nil)
	err = t.Execute(funcSource, data)
	if err != nil {
		return nil, fmt.Errorf("func template generation: %w", err)
	}

	file, err := parser.ParseFile(fset, "mapgenie_temp.go", "package main\n"+funcSource.String(), 0)
	if err != nil {
		return nil, fmt.Errorf("parse generated fragment: %w", err)
	}
	fset.RemoveFile(fset.File(file.Pos()))

	return file.Decls[0].(*ast.FuncDecl), nil
}

func generateExpressions(ctx context.Context, data MapTemplateData) ([]string, error) {
	generationCtx := &fragments.GenerationCtx{
		Ctx:          ctx,
		NameResolver: data.Resolver,
	}

	results := make([]string, len(data.Maps))
	for i, m := range data.Maps {
		code, err := m.Generate(generationCtx)
		if err != nil {
			return nil, fmt.Errorf("%+v fragment generation: %w", m, err)
		}

		results[i] = code
	}

	return results, nil
}
