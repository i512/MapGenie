package fragments

type OneLinerFrag struct {
	BaseFrag
	Line string
}

func (f *OneLinerFrag) Lines() []string {
	return []string{f.Line}
}

func (f *OneLinerFrag) VarSet(VarSet) {
}

func (f *OneLinerFrag) ResVar() *Var {
	return nil
}
