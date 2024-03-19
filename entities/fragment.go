package entities

import (
	"go/types"
)

type Fragment interface {
	Deps(DepReg)
	Body() Writer
	Result() Writer
}

type DepReg interface {
	Var(...*Var)
	Pkg(...*Pkg)
	Type(...*Type)
	Register(...Fragment)
}

type Var struct {
	DesiredName string
	ContextName string
	Name        string // an available name is chosen at generation phase
}

type Pkg struct {
	Path      string
	LocalName string // an available name is chosen at generation phase
}

type Type struct {
	Type      types.Type
	LocalName string // an available name is chosen at generation phase
}

type Writer interface {
	Ln(...string) Writer
	Lnf(string, ...any) Writer
	Merge(Writer) Writer
	Indent(func(Writer)) Writer
	String() string
	Lines() []string
	PreLastLine() Writer
	LastLine() string
}
