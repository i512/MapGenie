package casts_pointers_to_underlying_type

import (
	"github.com/stretchr/testify/assert"
	"mapgenie/pkg/ptr"
	"testing"
)

func TestMapABWithValues(t *testing.T) {
	orig := A{
		PtrToPtr: ptr.Take(1),
		PtrToVal: ptr.Take(2),
		ValToPtr: 3,

		StrPtrToByte: ptr.Take("3"),
		ByteToStrPtr: []byte("4"),
		StrToBytePtr: "5",
		BytePtrToStr: ptr.Take([]byte("6")),

		StrPtrToBytePtr: ptr.Take("7"),
		BytePtrToStrPtr: ptr.Take([]byte("8")),
	}

	dest := MapAB(orig)
	assert.Equal(t, ptr.Take(Int(1)), dest.PtrToPtr)
	assert.Equal(t, Int(2), dest.PtrToVal)
	assert.Equal(t, ptr.Take(Int(3)), dest.ValToPtr)

	assert.Equal(t, []byte("3"), dest.StrPtrToByte)
	assert.Equal(t, ptr.Take("4"), dest.ByteToStrPtr)
	assert.Equal(t, ptr.Take([]byte("5")), dest.StrToBytePtr)
	assert.Equal(t, "6", dest.BytePtrToStr)

	assert.Equal(t, ptr.Take([]byte("7")), dest.StrPtrToBytePtr)
	assert.Equal(t, ptr.Take("8"), dest.BytePtrToStrPtr)
}

func TestMapABWithEmptyValues(t *testing.T) {
	MapAB(A{})
	MapBA(B{})
}
