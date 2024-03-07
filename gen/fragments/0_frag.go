package fragments

import "go/types"

type VarSet map[*Var]struct{}
type PkgSet map[*Pkg]struct{}
type TypeSet map[*Type]struct{}

type Fragment interface {
	VarSet(set VarSet)
	PkgSet(PkgSet)
	TypeSet(TypeSet)

	ResVar() *Var

	Lines() []string
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

type BaseFrag struct{}

func (f BaseFrag) VarSet(VarSet)   {}
func (f BaseFrag) PkgSet(PkgSet)   {}
func (f BaseFrag) TypeSet(TypeSet) {}
func (f BaseFrag) ResVar() *Var {
	return nil
}
