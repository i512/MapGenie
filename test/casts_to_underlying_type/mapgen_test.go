package casts_to_underlying_type

import (
	"github.com/i512/mapgenie/test/casts_to_underlying_type/other"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB(t *testing.T) {
	orig := A{
		A: "str",
		B: 'b',
		C: []byte("hello"),
		D: map[int]int{1: 2},
	}
	dest := MapAB(orig)

	assert.Equal(t, String("str"), dest.A)
	assert.Equal(t, other.Byte('b'), dest.B)
	assert.Equal(t, ByteSlice("hello"), dest.C)
	assert.Equal(t, IntMap{1: 2}, dest.D)
}

func TestMapBA(t *testing.T) {
	orig := B{
		A: String("str"),
		B: other.Byte('b'),
		C: ByteSlice("hello"),
		D: IntMap{1: 2},
	}
	dest := MapBA(orig)

	assert.Equal(t, "str", dest.A)
	assert.Equal(t, byte('b'), dest.B)
	assert.Equal(t, []byte("hello"), dest.C)
	assert.Equal(t, map[int]int{1: 2}, dest.D)
}
