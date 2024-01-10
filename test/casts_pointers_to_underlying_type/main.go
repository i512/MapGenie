package casts_pointers_to_underlying_type

type Int int

type A struct {
	PtrToPtr *int
	PtrToVal *int
	ValToPtr int
}

type B struct {
	PtrToPtr *Int
	PtrToVal Int
	ValToPtr *Int
}
