package gen

type AssignNumberToString struct {
	BaseMapStatement
	FmtName string
}

func NewAssignNumberToString(base BaseMapStatement) *AssignNumberToString {
	return &AssignNumberToString{BaseMapStatement: base}
}

func (c *AssignNumberToString) Generate(resolver *FileImports) (string, error) {
	c.FmtName = resolver.ResolvePkgImport("fmt")

	sourceTemplate :=
		`result.{{ .OutField }} = {{ .FmtName }}.Sprint(input.{{ .InField }})`

	return c.RunTemplate(c, sourceTemplate)
}
