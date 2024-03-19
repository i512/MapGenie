package fragments

import (
	"go/types"
	"mapgenie/entities"
)

type TimeToNumber struct {
	BaseMapStatement
	BaseFrag
}

func NewTimeToNumber(base BaseMapStatement) *TimeToNumber {
	f := &TimeToNumber{
		BaseMapStatement: base,
	}

	return f
}

func (f *TimeToNumber) Result() entities.Writer {
	w := writer()

	if b, ok := f.Out.(*types.Basic); ok && b.Kind() != types.Int64 {
		castWith := f.Out.String()
		w.Ln(castWith, "(input.", f.InField, ".Unix())")
	} else {
		w.Ln("input.", f.InField, ".Unix()")
	}

	return w
}
