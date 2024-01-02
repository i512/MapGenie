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

	assert.Equal(t, 1<<63-1, dest.Int)
	assert.Equal(t, 1<<7-1, dest.Int8)
	assert.Equal(t, 1<<15-1, dest.Int16)
	assert.Equal(t, 1<<31-1, dest.Int32)
	assert.Equal(t, 1<<63-1, dest.Int64)
	assert.Equal(t, -1, dest.Uint)
	assert.Equal(t, 1<<8-1, dest.Uint8)
	assert.Equal(t, 1<<16-1, dest.Uint16)
	assert.Equal(t, 1<<32-1, dest.Uint32)
	assert.Equal(t, -1, dest.Uint64)
	assert.Equal(t, 32, dest.Float32)
	assert.Equal(t, 64, dest.Float64)
	assert.Equal(t, 255, dest.Byte)
	assert.Equal(t, 1<<31-1, dest.Rune)
}

func TestMapBA(t *testing.T) {
	origin := B{
		Int:     1,
		Int8:    2,
		Int16:   3,
		Int32:   4,
		Int64:   5,
		Uint:    6,
		Uint8:   7,
		Uint16:  8,
		Uint32:  9,
		Uint64:  10,
		Float32: 11,
		Float64: 12,
		Byte:    13,
		Rune:    14,
	}
	dest := MapBA(origin)

	assert.Equal(t, 1, dest.Int)
	assert.Equal(t, int8(2), dest.Int8)
	assert.Equal(t, int16(3), dest.Int16)
	assert.Equal(t, int32(4), dest.Int32)
	assert.Equal(t, int64(5), dest.Int64)
	assert.Equal(t, uint(6), dest.Uint)
	assert.Equal(t, uint8(7), dest.Uint8)
	assert.Equal(t, uint16(8), dest.Uint16)
	assert.Equal(t, uint32(9), dest.Uint32)
	assert.Equal(t, uint64(10), dest.Uint64)
	assert.Equal(t, float32(11), dest.Float32)
	assert.Equal(t, float64(12), dest.Float64)
	assert.Equal(t, byte(13), dest.Byte)
	assert.Equal(t, rune(14), dest.Rune)
}
