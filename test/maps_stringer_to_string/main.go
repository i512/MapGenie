package maps_stringer_to_string

import "fmt"

type A struct {
	Stringer       fmt.Stringer
	StructStringer StructStringer
	TypeStringer   TypeStringer
}

type B struct {
	Stringer       string
	StructStringer string
	TypeStringer   string
}

type TypeStringer int

func (s TypeStringer) String() string {
	return fmt.Sprint(s)
}

type StructStringer struct {
	V string
}

func (s StructStringer) String() string {
	return s.V
}
