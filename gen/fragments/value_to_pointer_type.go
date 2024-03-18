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

func (f *ValueToPointerType) Lines() Writer {
	if f.CastWith == nil {
		return writer().Ln("&input.", f.InField)
	}

	return writer().
		Ln(f.Var.Name, " := ", f.CastWith.LocalName, "(input.", f.InField, ")").
		Ln("&", f.Var.Name)
}

func (f *ValueToPointerType) Deps(registry *DependencyRegistry) {
	registry.Type(f.CastWith)
	registry.Type(f.Type)
	registry.Var(f.Var)
}
