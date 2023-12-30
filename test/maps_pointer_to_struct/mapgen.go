package maps_pointer_to_struct

// MapStructPtr map this pls
func MapStructPtr(input *A) B {
	var result B
	if input == nil {
		return result
	}
	
	result.Int = input.Int

	return result
}


// MapToStructPtr map this pls
func MapToStructPtr(input A) *B {
	var result B
	result.Int = input.Int

	return &result
}


// MapStructPtrToStructPtr map this pls
func MapStructPtrToStructPtr(input *A) *B {
	var result B
	if input == nil {
		return &result
	}
	
	result.Int = input.Int

	return &result
}

