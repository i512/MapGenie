package maps_pointer_to_value

// MapAB map this pls
func MapAB(input A) B {
	var result B
	result.Int = &input.Int

	return result
}

