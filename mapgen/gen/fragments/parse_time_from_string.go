package fragments

type ParseTimeFromString struct {
	BaseMapStatement
	Format   string
	TimeName string
}

func NewParseTimeFromString(base BaseMapStatement) *ParseTimeFromString {
	return &ParseTimeFromString{
		BaseMapStatement: base,
		Format:           "RFC3339",
	}
}

func (c *ParseTimeFromString) Generate(g *GenerationCtx) (string, error) {
	c.TimeName = g.NameResolver.ResolvePkgImport("time")

	sourceTemplate := `
{{ .OutField }}, err := {{ .TimeName }}.Parse({{ .TimeName }}.{{ .Format }}, input.{{ .InField }}) 
if err == nil {
	result.{{ .OutField }} = {{ .OutField }}
}
`

	return c.RunTemplate(c, sourceTemplate)
}
