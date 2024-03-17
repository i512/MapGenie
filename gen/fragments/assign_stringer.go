package fragments

type StringerToString struct {
	BaseMapStatement
	BaseFrag
}

func NewStringerToString(base BaseMapStatement) *StringerToString {
	f := &StringerToString{
		BaseMapStatement: base,
	}

	return f
}

func (f *StringerToString) Lines() []string {
	return writer().s("input.", f.InField, ".String()").Lines()
}
