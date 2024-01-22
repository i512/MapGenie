package maps_unexported_fields_within_pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB(t *testing.T) {
	origin := A{u: 1}
	dest := MapAB(origin)

	assert.Equal(t, origin.u, dest.u)
}

func TestMapAC(t *testing.T) {
	origin := A{u: 1}
	MapAC(origin)
}
