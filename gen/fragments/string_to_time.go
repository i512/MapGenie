package fragments

type StringToTime struct {
	Time   *Pkg
	Var    *Var
	Format string
	BaseMapStatement
	BaseFrag
}

func NewStringToTime(base BaseMapStatement) *StringToTime {
	f := &StringToTime{
		BaseMapStatement: base,
		Var:              &Var{DesiredName: base.OutField},
		Time:             &Pkg{Path: "time"},
		Format:           "RFC3339",
	}

	return f
}

func (f *StringToTime) Lines() []string {
	w := writer().s(f.Var.Name, ", _ := ", f.Time.LocalName, ".Parse(", f.Time.LocalName, ".", f.Format, ", input.", f.InField, ")")
	w.s(f.Var.Name)

	return w.Lines()
}

func (f *StringToTime) Deps(registry *DependencyRegistry) {
	registry.Pkg(f.Time)
	registry.Var(f.Var)
}
