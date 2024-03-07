package fragments

type AssignFrag struct {
	BaseFrag
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

func (f *AssignFrag) Deps(r *DependencyRegistry) {
	r.Register(f.Source)
}

func (f *AssignFrag) ResVar() *Var {
	return f.Result
}
