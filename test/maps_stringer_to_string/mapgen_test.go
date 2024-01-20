package maps_stringer_to_string

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB(t *testing.T) {
	orig := A{
		Stringer:       bytes.NewBuffer([]byte("1")),
		StructStringer: StructStringer{V: "2"},
		TypeStringer:   3,
	}
	dest := MapAB(orig)

	assert.Equal(t, "1", dest.Stringer)
	assert.Equal(t, "2", dest.StructStringer)
	assert.Equal(t, "3", dest.TypeStringer)
}
