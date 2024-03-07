package casts_number_types

// MapAB map this pls
func MapAB(input A) B {
	return B{Uint16: int(input.Uint16), Int8: int(input.Int8), Int32: int(input.Int32), Int: input.Int, Byte: int(input.Byte), Uint8: int(input.Uint8), Int16: int(input.Int16), Uint: int(input.Uint), Uint32: int(input.Uint32), Int64: int(input.Int64), Float32: int(input.Float32), Rune: int(input.Rune), Uint64: int(input.Uint64), Float64: int(input.Float64)}
}

// MapBA map this pls
func MapBA(input B) A {
	return A{Byte: byte(input.Byte), Int8: int8(input.Int8), Float64: float64(input.Float64), Rune: rune(input.Rune), Uint16: uint16(input.Uint16), Float32: float32(input.Float32), Uint32: uint32(input.Uint32), Int: input.Int, Int32: int32(input.Int32), Uint: uint(input.Uint), Uint8: uint8(input.Uint8), Uint64: uint64(input.Uint64), Int16: int16(input.Int16), Int64: int64(input.Int64)}
}
