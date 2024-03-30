package maps_imported_structs

import "github.com/i512/mapgenie/test/maps_imported_structs/imported/b"

// MapAB map this pls
func MapAB(input A) b.B {
	return b.B{Int: input.Int}
}

// MapBA map this pls
func MapBA(input b.B) A {
	return A{Int: input.Int}
}

// MapBB map this pls
func MapBB(input b.B) b.B {
	return b.B{Int: input.Int, UncommonIntB: input.UncommonIntB}
}
