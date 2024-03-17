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

func (f *PtrToValue) Lines() []string {
	w := writer()
	w.s("var ", f.Result.Name, " ", f.CastWith.LocalName)
	w.s("if input.", f.InField, " != nil {").ident(func(w *LineWriter) {
		if f.CastWith == nil {
			w.s(f.Result.Name, " = *input.", f.InField)
		} else {
			w.s(f.Result.Name, " = ", f.CastWith.LocalName, "(*input.", f.InField, ")")
		}
	})
	w.s("}")

	w.s(f.Result.Name)
	return w.Lines()
}

func (f *PtrToValue) Deps(r *DependencyRegistry) {
	r.Type(f.CastWith)
	r.Var(f.Result)
}
