package fragments

import (
	"mapgenie/entities"
	"reflect"
)

type PtrToValue struct {
	BaseFrag
	BaseMapStatement
	CastWith *entities.Type
	Var      *entities.Var
}

func NewPtrToValue(base BaseMapStatement) *PtrToValue {
	f := &PtrToValue{
		BaseMapStatement: base,
		Var:              &entities.Var{DesiredName: base.OutField},
	}

	if !reflect.DeepEqual(f.In, f.Out) {
		f.CastWith = &entities.Type{Type: base.Out}
	}

	return f
}

func (f *PtrToValue) Body() entities.Writer {
	w := writer()
	w.Ln("var ", f.Var.Name, " ", f.CastWith.LocalName)
	w.Ln("if input.", f.InField, " != nil {").Indent(func(w entities.Writer) {
		if f.CastWith == nil {
			w.Ln(f.Var.Name, " = *input.", f.InField)
		} else {
			w.Ln(f.Var.Name, " = ", f.CastWith.LocalName, "(*input.", f.InField, ")")
		}
	})
	w.Ln("}")

	return w
}

func (f *PtrToValue) Result() entities.Writer {
	return writer().Ln(f.Var.Name)
}

func (f *PtrToValue) Deps(r entities.DepReg) {
	r.Type(f.CastWith)
	r.Var(f.Var)
}
