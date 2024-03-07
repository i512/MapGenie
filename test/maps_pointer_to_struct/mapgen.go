package maps_pointer_to_struct

// MapStructPtr map this pls
func MapStructPtr(input *A) B {
	if input == nil {
		return B{}
	}
	return B{Int: input.Int}
}

// MapToStructPtr map this pls
func MapToStructPtr(input A) *B {
	return &B{Int: input.Int}
}

// MapStructPtrToStructPtr map this pls
func MapStructPtrToStructPtr(input *A) *B {
	if input == nil {
		return &B{}
	}
	return &B{Int: input.Int}
}
