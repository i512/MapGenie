package maps_imported_structs

import "mopper/test/maps_imported_structs/imported/b"

// MapAB map this pls
func MapAB(a A) b.B {
	var result b.B

	return result
}

// MapBA map this pls
func MapBA(a b.B) A {
	var result A

	return result
}

// MapBB map this pls
func MapBB(a b.B) b.B {
	var result b.B

	return result
}
