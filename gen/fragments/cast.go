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

func (f *Cast) Lines() []string {
	w := writer()
	if f.outType != nil {
		w.s(f.outType.LocalName, "(input.", f.InField, ")")
	} else {
		w.s("input.", f.InField)
	}
	return w.Lines()
}

func (f *Cast) Deps(r *DependencyRegistry) {
	r.Type(f.outType)
}
