package fragments

type AssignFrag struct {
	Source Fragment
	Result *Var
}

func NewAssign(source Fragment, name string) *AssignFrag {
	return &AssignFrag{
		Source: source,
		Result: &Var{DesiredName: name},
	}
}

func (f *AssignFrag) Lines() []string {
	sourceLines := f.Source.Lines()

	if f.Source.ResVar() != nil {
		return append(sourceLines, f.Result.Name+" := "+f.Source.ResVar().Name)
	}

	return append(sourceLines[:len(sourceLines)-1], f.Result.Name+" := "+sourceLines[len(sourceLines)-1])
}

func (f *AssignFrag) VarSet(set VarSet) {
	set[f.Result] = struct{}{}
	f.Source.VarSet(set)
}

func (f *AssignFrag) ResVars() []*Var {
	return []*Var{f.Result}
}
