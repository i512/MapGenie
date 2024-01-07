package casts_to_underlying_type

import (
	"github.com/stretchr/testify/assert"
	"mapgenie/test/casts_to_underlying_type/other"
	"testing"
)

func TestMapAB(t *testing.T) {
	orig := A{A: "str", B: []byte("slice")}
	dest := MapAB(orig)

	assert.Equal(t, String("str"), dest.A)
	assert.Equal(t, other.ByteSlice("slice"), dest.B)
}

func TestMapBA(t *testing.T) {
	orig := B{A: String("str")}
	dest := MapBA(orig)

	assert.Equal(t, "str", dest.A)
	assert.Equal(t, []byte("slice"), dest.B)
}
