package fragments

import "github.com/i512/mapgenie/entities"

type NumberToString struct {
	Fmt *entities.Pkg
	BaseMapStatement
	BaseFrag
}

func NewNumberToString(base BaseMapStatement) *NumberToString {
	f := &NumberToString{
		BaseMapStatement: base,
		Fmt:              &entities.Pkg{Path: "fmt"},
	}

	return f
}

func (f *NumberToString) Result() entities.Writer {
	return writer().Ln(f.Fmt.LocalName, ".Sprint(input.", f.InField, ")")
}

func (f *NumberToString) Deps(registry entities.DepReg) {
	registry.Pkg(f.Fmt)
}
