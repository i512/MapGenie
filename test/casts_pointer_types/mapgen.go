package casts_pointer_types

// MapAB map this pls
func MapAB(input A) B {
	return B{}
}

// MapBA map this pls
func MapBA(input B) A {
	var IntToIntPtr Int
	if input.IntToIntPtr != nil {
		IntToIntPtr = Int(*input.IntToIntPtr)
	}
	return A{IntToIntPtr: IntToIntPtr}
}
