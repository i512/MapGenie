package casts_pointers_to_underlying_type

import (
	"github.com/stretchr/testify/assert"
	"mapgenie/pkg/ptr"
	"testing"
)

func TestMapAB(t *testing.T) {
	orig := A{
		PtrToPtr: ptr.Take(1),
		PtrToVal: ptr.Take(2),
		ValToPtr: 3,
	}

	dest := MapAB(orig)
	assert.Equal(t, ptr.Take(Int(1)), dest.PtrToPtr)
	assert.Equal(t, Int(2), dest.PtrToVal)
	assert.Equal(t, ptr.Take(Int(3)), dest.ValToPtr)
}
