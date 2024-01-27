package gen

type TimeToString struct {
	BaseMapStatement
	CastWith string
	TimeName string
	Format   string
}

func NewTimeToString(base BaseMapStatement) *TimeToString {
	return &TimeToString{BaseMapStatement: base, Format: "RFC3339"}
}

func (c *TimeToString) String(resolver *FileImports) string {
	c.TimeName = resolver.ResolvePkgImport("time")
	sourceTemplate := `result.{{ .OutField }} = input.{{ .InField }}.Format({{ .TimeName }}.{{ .Format }})`

	return c.RunTemplate(c, sourceTemplate)
}
