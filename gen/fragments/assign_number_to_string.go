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

type NumberToString struct {
	Fmt *Pkg
	BaseMapStatement
	BaseFrag
}

func NewNumberToString(base BaseMapStatement) *NumberToString {
	f := &NumberToString{
		BaseMapStatement: base,
		Fmt:              &Pkg{Path: "fmt"},
	}

	return f
}

func (f *NumberToString) Lines() []string {
	return writer().s(f.Fmt.LocalName, ".Sprint(input.", f.InField, ")").Lines()
}

func (f *NumberToString) Deps(registry *DependencyRegistry) {
	registry.Pkg(f.Fmt)
}
