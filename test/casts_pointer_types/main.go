package casts_pointer_types

type Int int
type IntPtr *int

type A struct {
	IntPtrToInt         IntPtr
	IntToIntPtr         Int
	IntTypeToIntPtrType Int
	IntPtrTypeToIntType IntPtr
}

type B struct {
	IntPtrToInt         int
	IntToIntPtr         *int
	IntTypeToIntPtrType IntPtr
	IntPtrTypeToIntType Int
}
