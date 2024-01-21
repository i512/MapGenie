package fmt_aliased

import (
	fmt "bytes"
	"mapgenie/test/casts_number_types_to_string/structs"
)

type Unused fmt.Buffer

// MapAB map this pls
func MapAB(structs.A) structs.B {
	return structs.B{}
}

// MapBA map this pls
func MapBA(structs.B) structs.A {
	return structs.A{}
}
