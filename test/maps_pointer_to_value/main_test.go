package maps_imported_structs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB_PointerNotNil(t *testing.T) {
	v := 1
	origin := A{Int: &v}
	dest := MapAB(origin)

	assert.Equal(t, *origin.Int, dest.Int)
}

func TestMapAB_PointerNil(t *testing.T) {
	origin := A{Int: nil}
	dest := MapAB(origin)

	assert.Equal(t, 0, dest.Int)
}
