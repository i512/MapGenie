package fragments

import (
	"mapgenie/entities"
	"reflect"
)

type Cast struct {
	BaseFrag
	BaseMapStatement
	outType *entities.Type
}

func NewCast(base BaseMapStatement) *Cast {
	f := &Cast{
		BaseMapStatement: base,
	}

	if !reflect.DeepEqual(base.In, base.Out) {
		f.outType = &entities.Type{Type: base.Out}
	}

	return f
}

func (f *Cast) Result() entities.Writer {
	w := writer()
	if f.outType != nil {
		w.Ln(f.outType.LocalName, "(input.", f.InField, ")")
	} else {
		w.Ln("input.", f.InField)
	}
	return w
}

func (f *Cast) Deps(r entities.DepReg) {
	r.Type(f.outType)
}
