package casts_strings_and_byte_slices

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB(t *testing.T) {
	origin := A{A: "Hello"}
	dest := MapAB(origin)

	assert.Equal(t, []byte("Hello"), dest.A)
}

func TestMapBA(t *testing.T) {
	origin := B{A: []byte("Hello")}
	dest := MapBA(origin)

	assert.Equal(t, "Hello", dest.A)
}
