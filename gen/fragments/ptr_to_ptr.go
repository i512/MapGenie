package fragments

import (
	"go/types"
	"mapgenie/entities"
	"reflect"
)

type PtrToPtr struct {
	CastWith *entities.Type
	Var      *entities.Var
	BaseMapStatement
	BaseFrag
}

func NewPtrToPtr(base BaseMapStatement) *PtrToPtr {
	f := &PtrToPtr{BaseMapStatement: base}

	outElemType := base.Out.(*types.Pointer).Elem()
	if !reflect.DeepEqual(base.In, outElemType) {
		f.CastWith = &entities.Type{Type: outElemType}
		f.Var = &entities.Var{DesiredName: base.OutField}
	}
	return f
}

func (f *PtrToPtr) Body() entities.Writer {
	w := writer()
	w.Ln("var ", f.Var.Name, " ", f.CastWith.LocalName)
	w.Ln("if input.", f.InField, " != nil {").Indent(func(w entities.Writer) {
		if f.CastWith == nil {
			w.Ln(f.Var.Name, " = *input.", f.InField)
		} else {
			w.Ln(f.Var.Name, " = ", f.CastWith.LocalName, "(*input.", f.InField, ")")
		}
	}).Ln("}")

	return w
}

func (f *PtrToPtr) Result() entities.Writer {
	return writer().Ln("&", f.Var.Name)
}

func (f *PtrToPtr) Deps(registry entities.DepReg) {
	registry.Type(f.CastWith)
	registry.Var(f.Var)
}
