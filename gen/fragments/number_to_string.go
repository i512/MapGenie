package fragments

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
