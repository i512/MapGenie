package fragments

import "reflect"

type Cast struct {
	BaseFrag
	BaseMapStatement
	outType *Type
}

func NewCast(base BaseMapStatement) *Cast {
	f := &Cast{
		BaseMapStatement: base,
	}

	if !reflect.DeepEqual(base.In, base.Out) {
		f.outType = &Type{Type: base.Out}
	}

	return f
}

func (f *Cast) Lines() Writer {
	w := writer()
	if f.outType != nil {
		w.Ln(f.outType.LocalName, "(input.", f.InField, ")")
	} else {
		w.Ln("input.", f.InField)
	}
	return w
}

func (f *Cast) Deps(r *DependencyRegistry) {
	r.Type(f.outType)
}
