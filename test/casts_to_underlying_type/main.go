package casts_to_underlying_type

import "github.com/i512/mapgenie/test/casts_to_underlying_type/other"

type A struct {
	A string
	B byte
	C []byte
	D map[int]int
}

type B struct {
	A String
	B other.Byte
	C ByteSlice
	D IntMap
}

type String string
type ByteSlice []byte
type IntMap map[int]int
