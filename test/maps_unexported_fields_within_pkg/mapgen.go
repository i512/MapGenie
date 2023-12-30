package maps_unexported_fields_within_pkg

import "mopper/test/maps_unexported_fields_within_pkg/c"

// MapAB map this pls
func MapAB(A) B {
	return B{}
}

// MapAC map this pls
func MapAC(A) c.C {
	return c.C{}
}
