package maps_imported_structs

import (
	"github.com/stretchr/testify/assert"
	"mapgenie/test/maps_imported_structs/imported/b"
	"testing"
)

func TestMapAB(t *testing.T) {
	origin := A{Int: 1}
	dest := MapAB(origin)

	assert.Equal(t, origin.Int, dest.Int)
}

func TestMapBA(t *testing.T) {
	origin := b.B{Int: 1}
	dest := MapBA(origin)

	assert.Equal(t, origin.Int, dest.Int)
}

func TestMapBB(t *testing.T) {
	origin := b.B{Int: 1, UncommonIntB: 2}
	dest := MapBB(origin)

	assert.Equal(t, origin.Int, dest.Int)
	assert.Equal(t, origin.UncommonIntB, dest.UncommonIntB)
}
