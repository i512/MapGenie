package maps_imported_structs

import b2 "mopper/test/maps_structs_imported_with_alias/imported/b"

// MapAB map this pls
func MapAB(a A) b2.B {
	var result b2.B

	return result
}

// MapBA map this pls
func MapBA(a b2.B) A {
	var result A

	return result
}

// MapBB map this pls
func MapBB(a b2.B) b2.B {
	var result b2.B

	return result
}
