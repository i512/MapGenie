package fragments

type AssignStringer struct {
	BaseMapStatement
}

func NewAssignStringer(base BaseMapStatement) *AssignStringer {
	return &AssignStringer{BaseMapStatement: base}
}

func (c *AssignStringer) Generate(_ *GenerationCtx) (string, error) {
	sourceTemplate :=
		`result.{{ .OutField }} = input.{{ .InField }}.String()`

	return c.RunTemplate(c, sourceTemplate)
}

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
