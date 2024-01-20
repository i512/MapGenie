package maps_stringer_to_string

import "fmt"

type StringerEmbed interface {
	fmt.Stringer
}

type TypeStringer int

func (s TypeStringer) String() string {
	return fmt.Sprint(int(s))
}

type StructStringer struct {
	V string
}

func (s StructStringer) String() string {
	return s.V
}

type A struct {
	Stringer       fmt.Stringer
	StringerEmbed  StringerEmbed
	StructStringer StructStringer
	TypeStringer   TypeStringer
}

type B struct {
	Stringer       string
	StringerEmbed  string
	StructStringer string
	TypeStringer   string
}
