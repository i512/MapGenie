package fmt_included

import (
	"fmt"
	"mapgenie/test/casts_number_types_to_string/structs"
)

func print() {
	fmt.Print("hello!")
}

// MapAB map this pls
func MapAB(structs.A) structs.B {
	return structs.B{}
}

// MapBA map this pls
func MapBA(structs.B) structs.A {
	return structs.A{}
}
