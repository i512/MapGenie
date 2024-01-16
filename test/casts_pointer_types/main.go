package casts_pointer_types

type Int int
type IntPtr *int

type A struct {
	IntPtrToInt         IntPtr
	IntToIntPtr         int
	IntTypeToIntPtrType Int
	IntPtrTypeToIntType IntPtr
}

type B struct {
	IntPtrToInt         int
	IntToIntPtr         IntPtr
	IntTypeToIntPtrType IntPtr
	IntPtrTypeToIntType Int
}
