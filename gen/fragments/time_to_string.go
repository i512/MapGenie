package fragments

type TimeToString struct {
	BaseMapStatement
	CastWith string
	TimeName string
	Format   string
}

func NewTimeToString(base BaseMapStatement) *TimeToString {
	return &TimeToString{BaseMapStatement: base, Format: "RFC3339"}
}

func (c *TimeToString) Generate(g *GenerationCtx) (string, error) {
	c.TimeName = g.NameResolver.ResolvePkgImport("time")
	sourceTemplate := `result.{{ .OutField }} = input.{{ .InField }}.Format({{ .TimeName }}.{{ .Format }})`

	return c.RunTemplate(c, sourceTemplate)
}
