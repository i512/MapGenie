package fragments

type PtrToValueFrag struct {
	BaseFrag
	BaseMapStatement
	CastWith *Type
	Result   *Var
}

func NewPtrToValueFrag(base BaseMapStatement) *PtrToValueFrag {
	return &PtrToValueFrag{
		BaseMapStatement: base,
		CastWith:         &Type{Type: base.Out},
		Result:           &Var{DesiredName: base.OutField},
	}
}

func (f *PtrToValueFrag) Lines() []string {
	return writer().
		s("var ", f.Result.Name, " ", f.CastWith.LocalName).
		s("if input.", f.InField, " != nil {").
		s(f.Result.Name, " = ").
		a((&CastFrag{Fragment: &OneLinerFrag{Line: "*input." + f.InField}, CastWith: f.CastWith.LocalName}).Lines()).
		s("}").
		Lines()
}

func (f *PtrToValueFrag) Deps(r *DependencyRegistry) {
	r.Type(f.CastWith)
}

func (f *PtrToValueFrag) ResVar() *Var {
	return f.Result
}