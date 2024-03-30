package gen

import (
	"fmt"
	"github.com/i512/mapgenie/entities"
)

type OrderedSet[V comparable] struct {
	slice []V
	set   map[V]struct{}
}

func (s *OrderedSet[V]) Add(v V) {
	if s.set == nil {
		s.set = make(map[V]struct{})
	}

	if _, ok := s.set[v]; ok {
		return
	}

	s.set[v] = struct{}{}
	s.slice = append(s.slice, v)
}

func (s *OrderedSet[V]) Slice() []V {
	return s.slice
}

func NewDependencyRegistry() *DependencyRegistry {
	return &DependencyRegistry{
		VarSet:  &OrderedSet[*entities.Var]{},
		PkgSet:  &OrderedSet[*entities.Pkg]{},
		TypeSet: &OrderedSet[*entities.Type]{},
	}
}

type DependencyRegistry struct {
	VarSet  *OrderedSet[*entities.Var]
	PkgSet  *OrderedSet[*entities.Pkg]
	TypeSet *OrderedSet[*entities.Type]
}

func (r *DependencyRegistry) Var(vars ...*entities.Var) {
	for _, v := range vars {
		if v == nil {
			continue
		}
		r.VarSet.Add(v)
	}
}

func (r *DependencyRegistry) Type(types ...*entities.Type) {
	for _, t := range types {
		if t == nil {
			continue
		}

		r.TypeSet.Add(t)
	}
}

func (r *DependencyRegistry) Pkg(pkgs ...*entities.Pkg) {
	for _, p := range pkgs {
		if p == nil {
			continue
		}
		r.PkgSet.Add(p)
	}
}

func (r *DependencyRegistry) Register(fragments ...entities.Fragment) {
	for _, f := range fragments {
		f.Deps(r)
	}
}

func (r *DependencyRegistry) Setup(imports *FileImports) {
	for _, t := range r.TypeSet.Slice() {
		t.LocalName = imports.ResolveTypeName(t.Type)
	}

	for _, pkg := range r.PkgSet.Slice() {
		pkg.LocalName = imports.ResolvePkgImport(pkg.Path)
	}

	locals := make(map[string]int)
	for _, v := range r.VarSet.Slice() {
		counter := locals[v.DesiredName]
		locals[v.DesiredName] += 1

		name := v.DesiredName
		if counter > 0 {
			name = fmt.Sprintf("%s%d", name, counter)
		}
		v.Name = name
	}
}
