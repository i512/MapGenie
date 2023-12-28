package maps_pointer_to_struct

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapStructPtr(t *testing.T) {
	origin := A{Int: 1}
	dest := MapStructPtr(&origin)

	assert.IsType(t, B{}, dest)
	assert.Equal(t, origin.Int, dest.Int)
}

func TestMapToStructPointer(t *testing.T) {
	origin := A{Int: 1}
	dest := MapToStructPtr(origin)

	assert.IsType(t, &B{}, dest)
	assert.Equal(t, origin.Int, dest.Int)
}

func TestMapStructPtrToStructPtr(t *testing.T) {
	origin := A{Int: 1}
	dest := MapStructPtrToStructPtr(&origin)

	assert.IsType(t, &B{}, dest)
	assert.Equal(t, origin.Int, dest.Int)
}
