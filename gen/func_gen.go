package gen

import (
	"context"
	"fmt"
	"github.com/i512/mapgenie/entities"
	"github.com/i512/mapgenie/gen/fragments"
	"github.com/i512/mapgenie/pkg/log"
	"go/ast"
	"go/parser"
	"go/token"
)

type MapTemplateData struct {
	InType   string
	InputVar string
	OutType  string
	OutIsPtr bool
}

func FuncAst(ctx context.Context, tf entities.TargetFunc, fset *token.FileSet, imports *FileImports) (*ast.FuncDecl, error) {
	mapFunc := fragments.NewMapFunc(tf)
	registry := NewDependencyRegistry()
	registry.Register(mapFunc)
	registry.Setup(imports)

	return generateAst(ctx, fset, mapFunc.Result().String())
}

func generateAst(ctx context.Context, fset *token.FileSet, funcSource string) (*ast.FuncDecl, error) {
	log.Debugf(ctx, "Generated source:\n%s", funcSource)

	file, err := parser.ParseFile(fset, "mapgenie_temp.go", "package main\n"+funcSource, 0)
	if err != nil {
		return nil, fmt.Errorf("parse generated fragment: %w:\n%s", err, funcSource)
	}
	fset.RemoveFile(fset.File(file.Pos()))

	return file.Decls[0].(*ast.FuncDecl), nil
}
