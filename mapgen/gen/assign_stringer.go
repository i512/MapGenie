package gen

type AssignStringer struct {
	BaseMapStatement
}

func NewAssignStringer(base BaseMapStatement) *AssignStringer {
	return &AssignStringer{BaseMapStatement: base}
}

func (c *AssignStringer) Generate(resolver *FileImports) (string, error) {
	sourceTemplate :=
		`result.{{ .OutField }} = input.{{ .InField }}.String()`

	return c.RunTemplate(c, sourceTemplate)
}
