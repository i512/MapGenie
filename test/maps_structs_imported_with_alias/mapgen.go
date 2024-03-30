package maps_imported_structs_with_alias

import b2 "github.com/i512/mapgenie/test/maps_structs_imported_with_alias/imported/b"

// MapAB map this pls
func MapAB(input A) b2.B {
	return b2.B{Int: input.Int}
}

// MapBA map this pls
func MapBA(input b2.B) A {
	return A{Int: input.Int}
}

// MapBB map this pls
func MapBB(input b2.B) b2.B {
	return b2.B{Int: input.Int, UncommonIntB: input.UncommonIntB}
}
