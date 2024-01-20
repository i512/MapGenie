package gen

type AssignStringer struct {
	BaseMapStatement
}

func NewAssignStringer(base BaseMapStatement) *AssignStringer {
	return &AssignStringer{BaseMapStatement: base}
}

func (c *AssignStringer) String(resolver *FileImports) string {
	sourceTemplate :=
		`result.{{ .OutField }} = input.{{ .InField }}.String()`

	return c.RunTemplate(c, sourceTemplate)
}
