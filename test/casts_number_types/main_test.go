package casts_number_types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB(t *testing.T) {
	origin := A{
		Int:     1<<63 - 1,
		Int8:    1<<7 - 1,
		Int16:   1<<15 - 1,
		Int32:   1<<31 - 1,
		Int64:   1<<63 - 1,
		Uint:    1<<64 - 1,
		Uint8:   1<<8 - 1,
		Uint16:  1<<16 - 1,
		Uint32:  1<<32 - 1,
		Uint64:  1<<64 - 1,
		Float32: 32,
		Float64: 64,
		Byte:    255,
		Rune:    1<<31 - 1,
	}
	dest := MapAB(origin)

	assert.Equal(t, origin.Int, dest.Int)
	assert.Equal(t, origin.Int8, dest.Int8)
	assert.Equal(t, origin.Int16, dest.Int16)
	assert.Equal(t, origin.Uint, dest.Uint)
}

func TestMapBA(t *testing.T) {
	origin := B{Int: 1}
	dest := MapBA(origin)

	assert.Equal(t, origin.Int, dest.Int)
}
