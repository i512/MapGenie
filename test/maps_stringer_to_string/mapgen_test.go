package maps_stringer_to_string

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB(t *testing.T) {
	orig := A{
		Stringer:       bytes.NewBuffer([]byte("1")),
		StringerEmbed:  bytes.NewBuffer([]byte("2")),
		StructStringer: StructStringer{V: "3"},
		TypeStringer:   4,
	}
	dest := MapAB(orig)

	assert.Equal(t, "1", dest.Stringer)
	assert.Equal(t, "2", dest.StringerEmbed)
	assert.Equal(t, "3", dest.StructStringer)
	assert.Equal(t, "4", dest.TypeStringer)
}
