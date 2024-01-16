package casts_pointer_types

// MapAB map this pls
func MapAB(input A) B {
	var result B
	if input.IntPtrToInt != nil {
		result.IntPtrToInt = *input.IntPtrToInt
	}
	if input.IntPtrTypeToIntType != nil {
		result.IntPtrTypeToIntType = Int(*input.IntPtrTypeToIntType)
	}
	return result
}

// MapBA map this pls
func MapBA(input B) A {
	var result A
	if input.IntToIntPtr != nil {
		result.IntToIntPtr = *input.IntToIntPtr
	}
	if input.IntTypeToIntPtrType != nil {
		result.IntTypeToIntPtrType = Int(*input.IntTypeToIntPtrType)
	}
	return result
}
