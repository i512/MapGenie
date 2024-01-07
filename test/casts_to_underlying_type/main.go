package casts_to_underlying_type

import "mapgenie/test/casts_to_underlying_type/other"

type A struct {
	A string
	B []byte
	C []byte
	D []int
}

type B struct {
	A String
	B other.ByteSlice
	C ByteSlice
	D IntSlice
}

type String string
type ByteSlice []byte
type IntSlice []int
