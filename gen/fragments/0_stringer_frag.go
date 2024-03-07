package fragments

type StringerFrag struct {
	InField string
}

func (f *StringerFrag) Lines() []string {
	return []string{"input." + f.InField + ".String()"}
}

func (f *StringerFrag) VarSet(set VarSet) {
}

func (f *StringerFrag) ResVar() *Var {
	return nil
}
