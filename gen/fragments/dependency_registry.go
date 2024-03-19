package fragments

import "mapgenie/entities"

type VarSet map[*entities.Var]struct{}
type PkgSet map[*entities.Pkg]struct{}
type TypeSet map[*entities.Type]struct{}

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

func (r *DependencyRegistry) Var(vars ...*entities.Var) {
	for _, v := range vars {
		if v == nil {
			continue
		}
		r.VarSet[v] = struct{}{}
	}
}

func (r *DependencyRegistry) Type(types ...*entities.Type) {
	for _, t := range types {
		if t == nil {
			continue
		}
		r.TypeSet[t] = struct{}{}
	}
}

func (r *DependencyRegistry) Pkg(pkgs ...*entities.Pkg) {
	for _, p := range pkgs {
		if p == nil {
			continue
		}
		r.PkgSet[p] = struct{}{}
	}
}

func (r *DependencyRegistry) Register(fragments ...entities.Fragment) {
	for _, f := range fragments {
		f.Deps(r)
	}
}
