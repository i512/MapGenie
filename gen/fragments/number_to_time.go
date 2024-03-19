package fragments

import (
	"go/types"
	"mapgenie/entities"
)

type NumberToTime struct {
	Time *entities.Pkg
	BaseMapStatement
	BaseFrag
}

func NewNumberToTime(base BaseMapStatement) *NumberToTime {
	f := &NumberToTime{
		BaseMapStatement: base,
		Time:             &entities.Pkg{Path: "time"},
	}

	return f
}

func (f *NumberToTime) Result() entities.Writer {
	w := writer()

	if b, ok := f.In.(*types.Basic); ok && b.Kind() != types.Int64 {
		w.Ln(f.Time.LocalName, ".Unix(int64(input.", f.InField, "), 0).UTC()")
	} else {
		w.Ln(f.Time.LocalName, ".Unix(input.", f.InField, ", 0).UTC()")
	}

	return w
}

func (f *NumberToTime) Deps(registry entities.DepReg) {
	registry.Pkg(f.Time)
}
