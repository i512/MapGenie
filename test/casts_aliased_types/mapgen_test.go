package casts_aliased_types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB(t *testing.T) {
	orig := A{T: 1}
	dest := MapAB(orig)
	assert.Equal(t, 1, dest.T)
}
