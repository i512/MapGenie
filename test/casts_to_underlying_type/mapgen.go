package casts_to_underlying_type

import "mapgenie/test/casts_to_underlying_type/other"

// MapAB map this pls
func MapAB(input A) B {
	return B{C: ByteSlice(input.C), D: IntMap(input.D), A: String(input.A), B: other.Byte(input.B)}
}

// MapBA map this pls
func MapBA(input B) A {
	return A{C: []byte(input.C), D: map[int]int(input.D), A: string(input.A), B: byte(input.B)}
}
