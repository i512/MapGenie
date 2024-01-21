package casts_numbers_to_string

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB(t *testing.T) {
	orig := A{
		Int:    1,
		Int8:   2,
		Int16:  3,
		Int32:  4,
		Int64:  5,
		Uint:   6,
		Uint8:  7,
		Uint16: 8,
		Uint32: 9,
		Uint64: 10,
	}
	dest := MapAB(orig)

	assert.Equal(t, orig.Int, dest.Int)
	assert.Equal(t, orig.Uint, dest.Uint)
}
