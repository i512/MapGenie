package casts_pointers_to_underlying_type

type Int int

type A struct {
	PtrToPtr *int
	PtrToVal *int
	ValToPtr int

	StrPtrToByte *string
	ByteToStrPtr []byte
	StrToBytePtr string
	BytePtrToStr *[]byte

	StrPtrToBytePtr *string
	BytePtrToStrPtr *[]byte
}

type B struct {
	PtrToPtr *Int
	PtrToVal Int
	ValToPtr *Int

	StrPtrToByte []byte
	ByteToStrPtr *string
	StrToBytePtr *[]byte
	BytePtrToStr string

	StrPtrToBytePtr *[]byte
	BytePtrToStrPtr *string
}
