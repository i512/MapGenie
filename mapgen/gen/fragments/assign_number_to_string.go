package fragments

type AssignNumberToString struct {
	BaseMapStatement
	FmtName string
}

func NewAssignNumberToString(base BaseMapStatement) *AssignNumberToString {
	return &AssignNumberToString{BaseMapStatement: base}
}

func (c *AssignNumberToString) Generate(g *GenerationCtx) (string, error) {
	c.FmtName = g.NameResolver.ResolvePkgImport("fmt")

	sourceTemplate :=
		`result.{{ .OutField }} = {{ .FmtName }}.Sprint(input.{{ .InField }})`

	return c.RunTemplate(c, sourceTemplate)
}
