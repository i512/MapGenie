package maps_pointer_to_value

// MapAB map this pls
func MapAB(input A) B {
	var Int int
	if input.Int != nil {
		Int = int(*input.Int)
	}
	return B{Int: Int}
}
