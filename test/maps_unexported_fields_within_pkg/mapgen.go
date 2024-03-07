package maps_unexported_fields_within_pkg

import "mapgenie/test/maps_unexported_fields_within_pkg/c"

// MapAB map this pls
func MapAB(input A) B {
	return B{u: input.u}
}

// MapAC map this pls
func MapAC(input A) c.C {
	return c.C{}
}
