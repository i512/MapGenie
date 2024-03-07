package casts_pointers_to_underlying_type

// MapAB map this pls
func MapAB(input A) B {
	var PtrToVal Int
	if input.PtrToVal != nil {
		PtrToVal = Int(*input.PtrToVal)
	}
	var StrPtrToByte []byte
	if input.StrPtrToByte != nil {
		StrPtrToByte = []byte(*input.StrPtrToByte)
	}
	var BytePtrToStr string
	if input.BytePtrToStr != nil {
		BytePtrToStr = string(*input.BytePtrToStr)
	}
	return B{PtrToVal: PtrToVal, StrPtrToByte: StrPtrToByte, BytePtrToStr: BytePtrToStr}
}

// MapBA map this pls
func MapBA(input B) A {
	var ValToPtr int
	if input.ValToPtr != nil {
		ValToPtr = int(*input.ValToPtr)
	}
	var ByteToStrPtr []byte
	if input.ByteToStrPtr != nil {
		ByteToStrPtr = []byte(*input.ByteToStrPtr)
	}
	var StrToBytePtr string
	if input.StrToBytePtr != nil {
		StrToBytePtr = string(*input.StrToBytePtr)
	}
	return A{ValToPtr: ValToPtr, ByteToStrPtr: ByteToStrPtr, StrToBytePtr: StrToBytePtr}
}
