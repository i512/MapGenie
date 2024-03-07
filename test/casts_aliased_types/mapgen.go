package casts_aliased_types

// MapAB map this pls
func MapAB(input A) B {
	return B{T: input.T}
}
