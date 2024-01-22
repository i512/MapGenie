package fmt_aliased

import (
	fmt "bytes"
	"mapgenie/test/casts_number_types_to_and_from_string/structs"
)

type Unused fmt.Buffer

// MapAB map this pls
func MapAB(input structs.A) structs.B {
	var result structs.B
	return result
}

// MapBA map this pls
func MapBA(input structs.B) structs.A {
	var result structs.A
	return result
}
