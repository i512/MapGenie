package fragments

type OneLinerFrag struct {
	BaseFrag
	Line string
}

func (f *OneLinerFrag) Lines() []string {
	return []string{f.Line}
}
