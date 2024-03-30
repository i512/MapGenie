package casts_pointer_types

import (
	"github.com/i512/mapgenie/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB(t *testing.T) {
	orig := A{
		IntPtrToInt:         IntPtr(ptr.Take(1)),
		IntToIntPtr:         2,
		IntTypeToIntPtrType: Int(3),
		IntPtrTypeToIntType: IntPtr(ptr.Take(4)),
	}
	dest := MapAB(orig)

	assert.Equal(t, 1, dest.IntPtrToInt)
	assert.Equal(t, ptr.Take(2), dest.IntToIntPtr)
	assert.Equal(t, IntPtr(ptr.Take(3)), dest.IntTypeToIntPtrType)
	assert.Equal(t, Int(4), dest.IntPtrTypeToIntType)
}
