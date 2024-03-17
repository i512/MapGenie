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

type TimeToString2 struct {
	Time   *Pkg
	Format string
	BaseMapStatement
	BaseFrag
}

func NewTimeToString2(base BaseMapStatement) *TimeToString2 {
	f := &TimeToString2{
		BaseMapStatement: base,
		Time:             &Pkg{Path: "time"},
		Format:           "RFC3339",
	}

	return f
}

func (f *TimeToString2) Lines() []string {
	w := writer().s("input.", f.InField, ".Format(", f.Time.LocalName, ".", f.Format, ")")

	return w.Lines()
}

func (f *TimeToString2) Deps(registry *DependencyRegistry) {
	registry.Pkg(f.Time)
}
