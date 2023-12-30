package maps_local_structs

// MapAB map this pls
func MapAB(input A) B {
	var result B
	result.MapIntInt = input.MapIntInt
	result.Int = input.Int
	result.IntPtr = input.IntPtr
	result.Str = input.Str
	result.StrPtr = input.StrPtr
	result.ByteSlice = input.ByteSlice
	result.IntSlice = input.IntSlice
	result.ChanInt = input.ChanInt

	return result
}


// MapBA map this pls
func MapBA(input B) A {
	var result A
	result.Int = input.Int
	result.MapIntInt = input.MapIntInt
	result.IntPtr = input.IntPtr
	result.StrPtr = input.StrPtr
	result.IntSlice = input.IntSlice
	result.Str = input.Str
	result.ByteSlice = input.ByteSlice
	result.ChanInt = input.ChanInt

	return result
}


// MapAA map this pls
func MapAA(input A) A {
	var result A
	result.UncommonMapA = input.UncommonMapA
	result.UncommonChanA = input.UncommonChanA
	result.StrPtr = input.StrPtr
	result.MapIntInt = input.MapIntInt
	result.UncommonSliceA = input.UncommonSliceA
	result.ChanInt = input.ChanInt
	result.Int = input.Int
	result.UncommonIntA = input.UncommonIntA
	result.UncommonStrA = input.UncommonStrA
	result.ByteSlice = input.ByteSlice
	result.IntSlice = input.IntSlice
	result.IntPtr = input.IntPtr
	result.Str = input.Str

	return result
}

