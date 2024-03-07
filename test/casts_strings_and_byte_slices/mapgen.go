package casts_strings_and_byte_slices

// MapAB map this pls
func MapAB(input A) B {
	return B{A: []byte(input.A)}
}

// MapBA map this pls
func MapBA(input B) A {
	return A{A: string(input.A)}
}
