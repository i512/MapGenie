package fragments

import "reflect"

type PtrToValue struct {
	BaseFrag
	BaseMapStatement
	CastWith *Type
	Result   *Var
}

func NewPtrToValue(base BaseMapStatement) *PtrToValue {
	f := &PtrToValue{
		BaseMapStatement: base,
		Result:           &Var{DesiredName: base.OutField},
	}

	if !reflect.DeepEqual(f.In, f.Out) {
		f.CastWith = &Type{Type: base.Out}
	}

	return f
}

func (f *PtrToValue) Lines() Writer {
	w := writer()
	w.S("var ", f.Result.Name, " ", f.CastWith.LocalName)
	w.S("if input.", f.InField, " != nil {").Indent(func(w Writer) {
		if f.CastWith == nil {
			w.S(f.Result.Name, " = *input.", f.InField)
		} else {
			w.S(f.Result.Name, " = ", f.CastWith.LocalName, "(*input.", f.InField, ")")
		}
	})
	w.S("}")

	return w.S(f.Result.Name)
}

func (f *PtrToValue) Deps(r *DependencyRegistry) {
	r.Type(f.CastWith)
	r.Var(f.Result)
}
