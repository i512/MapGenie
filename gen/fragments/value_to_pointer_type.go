package fragments

import (
	"reflect"
)

type ValueToPointerType struct {
	CastWith *Type
	Type     *Type
	Var      *Var
	BaseMapStatement
	BaseFrag
}

func NewValueToPointerType(base BaseMapStatement) *ValueToPointerType {
	f := &ValueToPointerType{BaseMapStatement: base}
	f.Type = &Type{Type: f.Out}
	f.Var = &Var{DesiredName: base.OutField}

	inUnderlying := f.In.Underlying()
	if !reflect.DeepEqual(base.In, inUnderlying) {
		f.CastWith = &Type{Type: inUnderlying}
	}
	return f
}

func (f *ValueToPointerType) Lines() []string {
	if f.CastWith == nil {
		return writer().s("&input.", f.InField).Lines()
	}

	return writer().
		s(f.Var.Name, " := ", f.CastWith.LocalName, "(input.", f.InField, ")").
		s("&", f.Var.Name).Lines()
}

func (f *ValueToPointerType) Deps(registry *DependencyRegistry) {
	registry.Type(f.CastWith)
	registry.Type(f.Type)
	registry.Var(f.Var)
}
