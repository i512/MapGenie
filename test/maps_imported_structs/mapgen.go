package maps_imported_structs

import "mopper/test/maps_imported_structs/imported/b"

// MapAB map this pls
func MapAB(input A) b.B {
	var result b.B
	result.Int = input.Int

	return result
}


// MapBA map this pls
func MapBA(input b.B) A {
	var result A
	result.Int = input.Int

	return result
}


// MapBB map this pls
func MapBB(input b.B) b.B {
	var result b.B
	result.UncommonIntB = input.UncommonIntB
	result.Int = input.Int

	return result
}

