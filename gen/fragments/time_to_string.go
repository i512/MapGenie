package fragments

import "github.com/i512/mapgenie/entities"

type TimeToString struct {
	Time   *entities.Pkg
	Format string
	BaseMapStatement
	BaseFrag
}

func NewTimeToString(base BaseMapStatement) *TimeToString {
	f := &TimeToString{
		BaseMapStatement: base,
		Time:             &entities.Pkg{Path: "time"},
		Format:           "RFC3339",
	}

	return f
}

func (f *TimeToString) Result() entities.Writer {
	w := writer().Ln("input.", f.InField, ".Format(", f.Time.LocalName, ".", f.Format, ")")

	return w
}

func (f *TimeToString) Deps(registry entities.DepReg) {
	registry.Pkg(f.Time)
}
