package maps_pointer_to_value

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB(t *testing.T) {
	origin := A{Int: 1}
	dest := MapAB(origin)

	assert.Equal(t, origin.Int, *dest.Int)
}
