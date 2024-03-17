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

func (f *PtrToPtr) Lines() []string {
	w := writer()
	w.s("var ", f.Var.Name, " ", f.CastWith.LocalName)
	w.s("if input.", f.InField, " != nil {").ident(func(w *LineWriter) {
		if f.CastWith == nil {
			w.s(f.Var.Name, " = *input.", f.InField)
		} else {
			w.s(f.Var.Name, " = ", f.CastWith.LocalName, "(*input.", f.InField, ")")
		}
	}).s("}")
	w.s("&", f.Var.Name)

	return w.Lines()
}

func (f *PtrToPtr) Deps(registry *DependencyRegistry) {
	registry.Type(f.CastWith)
	registry.Var(f.Var)
}
