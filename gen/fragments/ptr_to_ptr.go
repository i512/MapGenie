package fragments

import (
	"go/types"
	"reflect"
)

type PtrToPtr struct {
	CastWith *Type
	Var      *Var
	BaseMapStatement
	BaseFrag
}

func NewPtrToPtr(base BaseMapStatement) *PtrToPtr {
	f := &PtrToPtr{BaseMapStatement: base}

	outElemType := base.Out.(*types.Pointer).Elem()
	if !reflect.DeepEqual(base.In, outElemType) {
		f.CastWith = &Type{Type: outElemType}
		f.Var = &Var{DesiredName: base.OutField}
	}
	return f
}

func (f *PtrToPtr) Lines() Writer {
	w := writer()
	w.S("var ", f.Var.Name, " ", f.CastWith.LocalName)
	w.S("if input.", f.InField, " != nil {").Indent(func(w Writer) {
		if f.CastWith == nil {
			w.S(f.Var.Name, " = *input.", f.InField)
		} else {
			w.S(f.Var.Name, " = ", f.CastWith.LocalName, "(*input.", f.InField, ")")
		}
	}).S("}")
	w.S("&", f.Var.Name)

	return w
}

func (f *PtrToPtr) Deps(registry *DependencyRegistry) {
	registry.Type(f.CastWith)
	registry.Var(f.Var)
}
