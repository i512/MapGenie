package fragments

type StringerFrag struct {
	BaseFrag
	InField string
}

func (f *StringerFrag) Lines() []string {
	return []string{"input." + f.InField + ".String()"}
}
