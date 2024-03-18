package fragments

import "go/types"

type NumberToTime struct {
	Time *Pkg
	BaseMapStatement
	BaseFrag
}

func NewNumberToTime(base BaseMapStatement) *NumberToTime {
	f := &NumberToTime{
		BaseMapStatement: base,
		Time:             &Pkg{Path: "time"},
	}

	return f
}

func (f *NumberToTime) Lines() Writer {
	w := writer()

	if b, ok := f.In.(*types.Basic); ok && b.Kind() != types.Int64 {
		w.S(f.Time.LocalName, ".Unix(int64(input.", f.InField, "), 0).UTC()")
	} else {
		w.S(f.Time.LocalName, ".Unix(input.", f.InField, ", 0).UTC()")
	}

	return w
}

func (f *NumberToTime) Deps(registry *DependencyRegistry) {
	registry.Pkg(f.Time)
}
