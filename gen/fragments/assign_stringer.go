package fragments

import "mapgenie/entities"

type StringerToString struct {
	BaseMapStatement
	BaseFrag
}

func NewStringerToString(base BaseMapStatement) *StringerToString {
	f := &StringerToString{
		BaseMapStatement: base,
	}

	return f
}

func (f *StringerToString) Result() entities.Writer {
	return writer().Ln("input.", f.InField, ".String()")
}
