package maps_local_structs

// MapAB map this pls
func MapAB(input A) B {
	return B{IntPtr: input.IntPtr, Str: input.Str, StrPtr: input.StrPtr, ByteSlice: input.ByteSlice, IntSlice: input.IntSlice, MapIntInt: input.MapIntInt, ChanInt: input.ChanInt, Int: input.Int}
}

// MapBA map this pls
func MapBA(input B) A {
	return A{StrPtr: input.StrPtr, ByteSlice: input.ByteSlice, IntSlice: input.IntSlice, MapIntInt: input.MapIntInt, ChanInt: input.ChanInt, Int: input.Int, IntPtr: input.IntPtr, Str: input.Str}
}

// MapAA map this pls
func MapAA(input A) A {
	return A{UncommonMapA: input.UncommonMapA, UncommonStrA: input.UncommonStrA, IntSlice: input.IntSlice, UncommonChanA: input.UncommonChanA, Str: input.Str, ByteSlice: input.ByteSlice, Int: input.Int, IntPtr: input.IntPtr, UncommonIntA: input.UncommonIntA, StrPtr: input.StrPtr, UncommonSliceA: input.UncommonSliceA, MapIntInt: input.MapIntInt, ChanInt: input.ChanInt}
}
