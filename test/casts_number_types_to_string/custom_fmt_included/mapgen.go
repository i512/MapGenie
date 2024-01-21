package custom_fmt_included

import (
	"mapgenie/test/casts_number_types_to_string/custom_fmt_included/fmt"
	"mapgenie/test/casts_number_types_to_string/structs"
)

type Unused2 fmt.Type

// MapAB map this pls
func MapAB(structs.A) structs.B {
	return structs.B{}
}

// MapBA map this pls
func MapBA(structs.B) structs.A {
	return structs.A{}
}
