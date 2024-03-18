package fragments

type TimeToString struct {
	Time   *Pkg
	Format string
	BaseMapStatement
	BaseFrag
}

func NewTimeToString(base BaseMapStatement) *TimeToString {
	f := &TimeToString{
		BaseMapStatement: base,
		Time:             &Pkg{Path: "time"},
		Format:           "RFC3339",
	}

	return f
}

func (f *TimeToString) Lines() Writer {
	w := writer().S("input.", f.InField, ".Format(", f.Time.LocalName, ".", f.Format, ")")

	return w
}

func (f *TimeToString) Deps(registry *DependencyRegistry) {
	registry.Pkg(f.Time)
}
