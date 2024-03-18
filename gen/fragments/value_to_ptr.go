package fragments

import (
	"go/types"
	"reflect"
)

type ValueToPtr struct {
	CastWith *Type
	Var      *Var
	BaseMapStatement
	BaseFrag
}

func NewValueToPtr(base BaseMapStatement) *ValueToPtr {
	f := &ValueToPtr{BaseMapStatement: base}

	outElemType := base.Out.(*types.Pointer).Elem()
	if !reflect.DeepEqual(base.In, outElemType) {
		f.CastWith = &Type{Type: outElemType}
		f.Var = &Var{DesiredName: base.OutField}
	}
	return f
}

func (f *ValueToPtr) Lines() Writer {
	if f.CastWith == nil {
		return writer().S("&input.", f.InField)
	}

	return writer().
		S(f.Var.Name, " := ", f.CastWith.LocalName, "(input.", f.InField, ")").
		S("&", f.Var.Name)
}

func (f *ValueToPtr) Deps(registry *DependencyRegistry) {
	registry.Type(f.CastWith)
	registry.Var(f.Var)
}
