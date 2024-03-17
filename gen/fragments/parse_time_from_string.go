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

type TimeFromString struct {
	Time   *Pkg
	Var    *Var
	Format string
	BaseMapStatement
	BaseFrag
}

func NewTimeFromString(base BaseMapStatement) *TimeFromString {
	f := &TimeFromString{
		BaseMapStatement: base,
		Var:              &Var{DesiredName: base.OutField},
		Time:             &Pkg{Path: "time"},
		Format:           "RFC3339",
	}

	return f
}

func (f *TimeFromString) Lines() []string {
	w := writer().s(f.Var.Name, ", _ := ", f.Time.LocalName, ".Parse(", f.Time.LocalName, ".", f.Format, ", input.", f.InField, ")")
	w.s(f.Var.Name)

	return w.Lines()
}

func (f *TimeFromString) Deps(registry *DependencyRegistry) {
	registry.Pkg(f.Time)
	registry.Var(f.Var)
}
