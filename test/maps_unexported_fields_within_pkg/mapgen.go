package maps_unexported_fields_within_pkg

import "mapgenie/test/maps_unexported_fields_within_pkg/c"

// MapAB map this pls
func MapAB(input A) B {
	var result B
	result.u = input.u

	return result
}

// MapAC map this pls
func MapAC(input A) c.C {
	var result c.C

	return result
}
