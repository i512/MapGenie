package maps_imported_structs

import (
	b2 "mopper/test/maps_structs_imported_with_alias/imported/b"
)

type A struct {
	Int          int
	UncommonIntA int
	unexported   int
}

// MapAB map this pls
func MapAB(a A) b2.B {
	var result b2.B

	result.Int = a.Int

	return result
}

// MapBA map this pls
func MapBA(a b2.B) A {
	var result A

	result.Int = a.Int

	return result
}

// MapBB map this pls
func MapBB(a b2.B) b2.B {
	var result b2.B

	result.Int = a.Int
	result.UncommonIntB = a.UncommonIntB

	return result
}
