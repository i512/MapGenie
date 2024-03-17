package fragments

import "go/types"

type VarSet map[*Var]struct{}
type PkgSet map[*Pkg]struct{}
type TypeSet map[*Type]struct{}

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

func NewDependencyRegistry() *DependencyRegistry {
	return &DependencyRegistry{
		VarSet:  VarSet{},
		PkgSet:  PkgSet{},
		TypeSet: TypeSet{},
	}
}

type DependencyRegistry struct {
	VarSet  VarSet
	PkgSet  PkgSet
	TypeSet TypeSet
}

func (r *DependencyRegistry) Var(vars ...*Var) {
	for _, v := range vars {
		if v == nil {
			continue
		}
		r.VarSet[v] = struct{}{}
	}
}

func (r *DependencyRegistry) Type(types ...*Type) {
	for _, t := range types {
		if t == nil {
			continue
		}
		r.TypeSet[t] = struct{}{}
	}
}

func (r *DependencyRegistry) Pkg(pkgs ...*Pkg) {
	for _, p := range pkgs {
		if p == nil {
			continue
		}
		r.PkgSet[p] = struct{}{}
	}
}

func (r *DependencyRegistry) Register(fragments ...Fragment) {
	for _, f := range fragments {
		if f.ResVar() != nil {
			r.Var(f.ResVar())
		}

		f.Deps(r)
	}
}
