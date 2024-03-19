package fragments

import "mapgenie/entities"

type StringToTime struct {
	Time   *entities.Pkg
	Var    *entities.Var
	Format string
	BaseMapStatement
	BaseFrag
}

func NewStringToTime(base BaseMapStatement) *StringToTime {
	f := &StringToTime{
		BaseMapStatement: base,
		Var:              &entities.Var{DesiredName: base.OutField},
		Time:             &entities.Pkg{Path: "time"},
		Format:           "RFC3339",
	}

	return f
}

func (f *StringToTime) Body() entities.Writer {
	w := writer().Ln(f.Var.Name, ", _ := ", f.Time.LocalName, ".Parse(", f.Time.LocalName, ".", f.Format, ", input.", f.InField, ")")

	return w
}
func (f *StringToTime) Result() entities.Writer {
	return writer().Ln(f.Var.Name)
}

func (f *StringToTime) Deps(registry entities.DepReg) {
	registry.Pkg(f.Time)
	registry.Var(f.Var)
}
