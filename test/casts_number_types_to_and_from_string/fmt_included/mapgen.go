package fmt_included

import (
	"fmt"
	"github.com/i512/mapgenie/test/casts_number_types_to_and_from_string/structs"
)

func print() {
	fmt.Print("hello!")
}

// MapAB map this pls
func MapAB(input structs.A) structs.B {
	return structs.B{}
}

// MapBA map this pls
func MapBA(input structs.B) structs.A {
	return structs.A{}
}
