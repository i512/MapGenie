package entities

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
)

type TargetFile struct {
	Pkg   *packages.Package
	Ast   *ast.File
	Fset  *token.FileSet
	Funcs []TargetFunc
}

func (f TargetFile) Name() string {
	return f.Fset.Position(f.Ast.Pos()).Filename
}

type TargetFunc struct {
	FuncDecl  *ast.FuncDecl
	In, Out   Argument
	Fragments map[string]Fragment
}

func (f TargetFunc) Name() string {
	return f.FuncDecl.Name.Name
}

type Argument struct {
	Named  *types.Named
	Struct *types.Struct
	IsPtr  bool
	Local  bool // Type is defined in target package
}

func (s Argument) FieldMap() map[string]types.Type {
	result := map[string]types.Type{}

	for i := 0; i < s.Struct.NumFields(); i++ {
		f := s.Struct.Field(i)
		result[f.Name()] = f.Type()
	}

	return result
}

type TypeNameResolver interface {
	ResolveTypeName(types.Type) string
}

type MapExpression interface {
	String(resolver TypeNameResolver) string
}
